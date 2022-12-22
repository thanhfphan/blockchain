package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
	"github.com/thanhfphan/blockchain/network/peer"
	"github.com/thanhfphan/blockchain/snow/networking/sender"
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
}

type network struct {
	config          *Config
	peerConfig      *peer.Config
	peersLock       sync.RWMutex
	connectingPeers peer.Set
	connectedPeers  peer.Set
	listener        net.Listener

	serverUpgrader peer.Upgrader
	clientUpgrader peer.Upgrader

	closeOnce        sync.Once
	onCloseCtx       context.Context
	onCloseCtxCancel func()
}

func New(config *Config, msgCreator message.Creator, listener net.Listener) (Network, error) {

	onCloseCtx, cancel := context.WithCancel(context.Background())
	peerConfig := &peer.Config{
		MessageCreator: msgCreator,
		PongTimeout:    config.PongTimeout,

		//Beacons //TODO
	}

	n := &network{
		config:          config,
		peerConfig:      peerConfig,
		listener:        listener,
		connectingPeers: peer.NewSet(),
		connectedPeers:  peer.NewSet(),

		onCloseCtx:       onCloseCtx,
		onCloseCtxCancel: cancel,

		clientUpgrader: peer.NewClientUpgrader(config.TLSConfig),
		serverUpgrader: peer.NewServerUpgrader(config.TLSConfig),
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
		fmt.Println("shutting down the p2p networking")
		if err := n.listener.Close(); err != nil {
			fmt.Printf("close the network listener err=%v\n", err)
		}

		n.peersLock.Lock()
		defer n.peersLock.Unlock()

		n.onCloseCtxCancel()

		for i := 0; i < n.connectingPeers.Len(); i++ {
			peer, err := n.connectingPeers.GetByIndex(i)
			if err != nil {
				fmt.Printf("get peer by index=%d failed\n", i)
			} else {
				peer.StartClose()
			}
		}

		for i := 0; i < n.connectedPeers.Len(); i++ {
			peer, err := n.connectedPeers.GetByIndex(i)
			if err != nil {
				fmt.Printf("get peer by index=%d failed\n", i)
			} else {
				peer.StartClose()
			}
		}
	})
}

// Dispatch start to accepting connections from other nodes to connect to this node
func (n *network) Dispatch() error {

	for {
		if n.onCloseCtx.Err() != nil {
			break
		}

		conn, err := n.listener.Accept()
		if err != nil {
			fmt.Printf("listen connection err=%v\n", err)
			time.Sleep(time.Millisecond)
			continue
		}

		go func() {
			if err := n.upgrade(conn, n.serverUpgrader); err != nil {
				fmt.Printf("perr upgrade err=%v\n", err)
			}
		}()
	}

	n.StartClose()
	//TODO: wait peer close

	return nil
}

func (n *network) upgrade(conn net.Conn, upgrader peer.Upgrader) error {
	nodeID, tlsConn, cert, err := upgrader.Upgrader(conn)
	if err != nil {
		_ = conn.Close()
		fmt.Printf("upgrade connection failed %v\n", err)
		return err
	}

	if err := tlsConn.SetReadDeadline(time.Time{}); err != nil {
		_ = tlsConn.Close()
		fmt.Printf("set readDeadLine failed %v\n", err)
		return err
	}

	if nodeID == n.config.MyNodeID {
		_ = tlsConn.Close()
		fmt.Printf("dropping connection to myslef\n")
		return nil
	}

	n.peersLock.Lock()
	defer n.peersLock.Unlock()

	if _, isConnecting := n.connectingPeers.GetByID(nodeID); isConnecting {
		fmt.Printf("dropping connection because already connecting to peer=%s", nodeID.String())
		return nil
	}

	if _, isConnected := n.connectedPeers.GetByID(nodeID); isConnected {
		fmt.Printf("dropping connection because already connected to peer=%s", nodeID.String())
		return nil
	}

	peer := peer.Start(n.peerConfig, conn, cert, nodeID, peer.NewBlockingQueue(maxMessageInQueue))
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

func (n *network) samplePeers(numberPeersToSample int) []peer.Peer {
	return n.connectedPeers.Sample(numberPeersToSample, preconditionSample)
}

func preconditionSample(p peer.Peer) bool {
	return true
}
