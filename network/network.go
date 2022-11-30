package network

import (
	"context"
	"net"

	"github.com/thanhfphan/blockchain/network/peer"
	"github.com/thanhfphan/blockchain/snow/networking/sender"
)

var _ Network = (*network)(nil)

type Network interface {
	sender.ExternalSender
	peer.Network
}

type network struct {
	connectingPeers peer.Set
	connectedPeers  peer.Set
}

func New(listener net.Listener) (Network, error) {
	n := &network{
		connectingPeers: peer.NewSet(),
		connectedPeers:  peer.NewSet(),
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
