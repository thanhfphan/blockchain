package network

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
	"github.com/thanhfphan/blockchain/network/dialer"
	"github.com/thanhfphan/blockchain/network/peer"
	"github.com/thanhfphan/blockchain/snow/networking/sender"
	"github.com/thanhfphan/blockchain/utils/ips"
	"github.com/thanhfphan/blockchain/utils/logging"
	"github.com/thanhfphan/blockchain/utils/sampler"
)

var _ Network = (*network)(nil)

const (
	maxMessageInQueue = 1024
)

type Network interface {
	sender.ExternalSender
	peer.Network

	// Should only be called once, run until error occur or network closed
	Dispatch() error
	StartClose()
	ManuallyTrack(nodeID ids.NodeID, ip ips.IPPort)
}

type network struct {
	config          *Config
	log             logging.Logger
	peerConfig      *peer.Config
	peersLock       sync.RWMutex
	connectingPeers peer.Set
	connectedPeers  peer.Set
	listener        net.Listener
	dialer          dialer.Dialer

	serverUpgrader peer.Upgrader
	clientUpgrader peer.Upgrader

	closeOnce        sync.Once
	onCloseCtx       context.Context
	onCloseCtxCancel func()

	gossiptracker peer.GossipTracker
}

func New(
	config *Config,
	log logging.Logger,
	msgCreator message.Creator,
	listener net.Listener,
	dialer dialer.Dialer,
	gossipTracker peer.GossipTracker,
) (Network, error) {

	onCloseCtx, cancel := context.WithCancel(context.Background())
	peerConfig := &peer.Config{
		MessageCreator:  msgCreator,
		PongTimeout:     config.PongTimeout,
		PingFrequency:   config.PingFrequency,
		ReadBufferSize:  config.PeerReadBufferSize,
		WriteBufferSize: config.PeerWriteBufferSize,

		IPSigner: peer.NewIPSigner(config.IPPort, config.TLSSignIPKey),
	}

	n := &network{
		log:             log,
		config:          config,
		peerConfig:      peerConfig,
		listener:        listener,
		connectingPeers: peer.NewSet(),
		connectedPeers:  peer.NewSet(),
		dialer:          dialer,

		onCloseCtx:       onCloseCtx,
		onCloseCtxCancel: cancel,

		clientUpgrader: peer.NewClientUpgrader(config.TLSConfig),
		serverUpgrader: peer.NewServerUpgrader(config.TLSConfig),

		gossiptracker: gossipTracker,
	}

	n.peerConfig.Network = n

	return n, nil
}

func (n *network) Send(msg message.OutboundMessage, nodeIDs []ids.NodeID) []ids.NodeID {
	peers := n.getPeeers(nodeIDs)
	return n.send(msg, peers)
}

func (n *network) Gossip(msg message.OutboundMessage, numerPeersToSend int) []ids.NodeID {
	peers := n.samplePeers(numerPeersToSend)
	return n.send(msg, peers)
}

func (n *network) Pong(nodeID ids.NodeID) (string, error) {
	return "", nil
}

// Call after handshake
func (n *network) Connected(nodeID ids.NodeID) {
	n.peersLock.Lock()
	peer, ok := n.connectingPeers.GetByID(nodeID)
	if !ok {
		n.peersLock.Unlock()
		return
	}

	n.connectingPeers.Remove(nodeID)
	n.connectedPeers.Add(peer)
	n.peersLock.Unlock()

	//TODO: send message to others
}

func (n *network) Disconnected(nodeID ids.NodeID) {
	n.peersLock.RLock()
	_, isConnecting := n.connectingPeers.GetByID(nodeID)
	_, isConnected := n.connectedPeers.GetByID(nodeID)
	n.peersLock.RUnlock()

	if isConnecting {
		n.peersLock.Lock()
		n.connectingPeers.Remove(nodeID)
		n.peersLock.Unlock()
	}

	if isConnected {
		n.peersLock.Lock()
		n.connectedPeers.Remove(nodeID)
		n.peersLock.Unlock()
	}
}

func (n *network) StartClose() {
	n.closeOnce.Do(func() {
		n.log.Verbof("Shutting down the p2p networking")
		if err := n.listener.Close(); err != nil {
			n.log.Errorf("Close the network listener err=%v\n", err)
		}

		n.peersLock.Lock()
		defer n.peersLock.Unlock()

		n.onCloseCtxCancel()

		for i := 0; i < n.connectingPeers.Len(); i++ {
			peer, err := n.connectingPeers.GetByIndex(i)
			if err != nil {
				n.log.Warnf("Get peer by index=%d failed\n", i)
			} else {
				peer.StartClose()
			}
		}

		for i := 0; i < n.connectedPeers.Len(); i++ {
			peer, err := n.connectedPeers.GetByIndex(i)
			if err != nil {
				n.log.Warnf("Get peer by index=%d failed\n", i)
			} else {
				peer.StartClose()
			}
		}
	})
}

func (n *network) ManuallyTrack(nodeID ids.NodeID, ip ips.IPPort) {
	n.peersLock.Lock()
	defer n.peersLock.Unlock()

	if _, isConnected := n.connectedPeers.GetByID(nodeID); isConnected {
		n.log.Warnf("%s already connected\n", nodeID)
		return
	}

	n.dial(n.onCloseCtx, nodeID, ip)
}

// Dispatch start to accepting connections from other nodes to connect to this node
func (n *network) Dispatch() error {

	for {
		if n.onCloseCtx.Err() != nil {
			break
		}

		conn, err := n.listener.Accept()
		if err != nil {
			n.log.Debugf("Listen tcp connection has error=%v\n", err)
			time.Sleep(time.Millisecond)
			continue
		}

		go func() {
			if err := n.upgrade(conn, n.serverUpgrader); err != nil {
				n.log.Errorf("Upgrade peer has error=%v\n", err)
			}
		}()
	}

	n.StartClose()
	//TODO: wait peer close

	return nil
}
func (n *network) Peers(peerID ids.NodeID) ([]ips.ClaimedIPPort, error) {
	unknown, ok := n.gossiptracker.GetUnknown(peerID)
	if !ok {
		n.log.Debugf("not found any unknown of peer %s to gossip to", peerID.String())
		return nil, nil
	}

	s := sampler.NewUniform()
	if err := s.Initialize(uint64(len(unknown))); err != nil {
		return nil, err
	}

	validators := make([]ips.ClaimedIPPort, 0, int(n.config.PeerListNumberValidator))

	for i := 0; i < len(unknown) && len(validators) < int(n.config.PeerListNumberValidator); i++ {
		draw, err := s.Next()
		if err != nil {
			return nil, err
		}

		v := unknown[draw]
		n.peersLock.RLock()
		p, ok := n.connectedPeers.GetByID(v.NodeID)
		n.peersLock.RUnlock()
		if !ok {
			n.log.Debugf("not found %s in list connected peers", v.NodeID)
			continue
		}
		peerIP := p.IP()
		validators = append(validators, ips.ClaimedIPPort{
			Cert:      p.Cert(),
			IPPort:    peerIP.IP.IP,
			Signature: peerIP.Signature,
			TimeStamp: peerIP.IP.Timestamp,
			TxID:      v.TxID,
		})
	}

	return validators, nil
}

func (n *network) upgrade(conn net.Conn, upgrader peer.Upgrader) error {
	// TODO: set timeout Upgrade

	nodeID, tlsConn, cert, err := upgrader.Upgrade(conn)
	if err != nil {
		_ = conn.Close()
		n.log.Errorf("Upgrade connection failed %v\n", err)
		return err
	}

	if err := tlsConn.SetReadDeadline(time.Time{}); err != nil {
		_ = tlsConn.Close()
		n.log.Errorf("Set readDeadLine failed %v\n", err)
		return err
	}

	if nodeID.String() == n.config.MyNodeID.String() {
		_ = tlsConn.Close()
		n.log.Verbof("Dropping connection to myslef\n")
		return nil
	}

	n.peersLock.Lock()
	defer n.peersLock.Unlock()

	if _, isConnecting := n.connectingPeers.GetByID(nodeID); isConnecting {
		n.log.Warnf("Dropping connection because already connecting to peer=%s", nodeID.String())
		return nil
	}

	if _, isConnected := n.connectedPeers.GetByID(nodeID); isConnected {
		n.log.Warnf("Dropping connection because already connected to peer=%s", nodeID.String())
		return nil
	}

	msgQueue := peer.NewBlockingQueue(maxMessageInQueue, n.log)
	peer := peer.Start(n.peerConfig, n.log, tlsConn, cert, nodeID, msgQueue, n.gossiptracker)
	n.connectingPeers.Add(peer)
	return nil
}

func (n *network) getPeeers(nodeIDs []ids.NodeID) []peer.Peer {
	peers := []peer.Peer{}

	for _, nID := range nodeIDs {
		peer, ok := n.connectedPeers.GetByID(nID)
		if !ok {
			continue
		}

		peers = append(peers, peer)
	}

	return peers
}

func (n *network) send(msg message.OutboundMessage, peers []peer.Peer) []ids.NodeID {
	sendTo := []ids.NodeID{}

	for _, peer := range peers {
		ok := peer.Send(n.onCloseCtx, msg)
		if ok {
			sendTo = append(sendTo, peer.ID())
		}
	}

	return sendTo
}

func (n *network) dial(ctx context.Context, nodeID ids.NodeID, ip ips.IPPort) {
	go func() {
		conn, err := n.dialer.Dial(ctx, ip)
		if err != nil {
			n.log.Errorf("Dial to IP: %s failed %v\n", ip.IP.String(), err)
			return
		}

		err = n.upgrade(conn, n.clientUpgrader)
		if err != nil {
			n.log.Errorf("Upgrade peer failed %v", err)
			return
		}

		return
	}()
}

func (n *network) samplePeers(numberPeersToSample int) []peer.Peer {
	return n.connectedPeers.Sample(numberPeersToSample, preconditionSample)
}

func preconditionSample(p peer.Peer) bool {
	return true
}
