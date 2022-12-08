package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/thanhfphan/blockchain/network/peer"
	"github.com/thanhfphan/blockchain/snow/networking/sender"
)

var _ Network = (*network)(nil)

type Network interface {
	sender.ExternalSender
	peer.Network

	// Should only be called once, run until error occur or network closed
	Dispatch() error
	StartClose()
}

type network struct {
	peersLock       sync.RWMutex
	connectingPeers peer.Set
	connectedPeers  peer.Set
	listener        net.Listener

	closeOnce        sync.Once
	onCloseCtx       context.Context
	onCloseCtxCancel func()
}

func New(listener net.Listener) (Network, error) {

	onCloseCtx, cancel := context.WithCancel(context.Background())
	n := &network{
		listener:         listener,
		connectingPeers:  peer.NewSet(),
		connectedPeers:   peer.NewSet(),
		onCloseCtx:       onCloseCtx,
		onCloseCtxCancel: cancel,
	}

	return n, nil
}

func (n *network) Send(msg string, nodeIDs []int) []int {
	peers := n.getPeeers(nodeIDs)
	return n.send(msg, peers)
}

func (n *network) Gossip(msg string, numerPeersToSend int) []int {
	peers := n.samplePeers(numerPeersToSend)
	return n.send(msg, peers)
}

func (n *network) Pong(nodeID int) (string, error) {
	return "", nil
}

// Call after handshake
func (n *network) Connected(nodeID int) {
	n.peersLock.Lock()
	peer, ok := n.connectingPeers.GetByID(nodeID)
	if !ok {
		n.peersLock.Unlock()
		return
	}

	n.connectingPeers.Remove(nodeID)
	n.connectedPeers.Add(peer)
	n.peersLock.Unlock()
}

func (n *network) Disconnected(nodeID int) {
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
	})
}

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

		remoteAddr := conn.RemoteAddr().String()
		fmt.Printf("Hello: %s\n", remoteAddr)

	}

	return nil
}

func (n *network) getPeeers(nodeIDs []int) []peer.Peer {
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

func (n *network) send(msg string, peers []peer.Peer) []int {
	sendTo := []int{}

	for _, peer := range peers {
		ok := peer.Send(context.Background(), msg)
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
