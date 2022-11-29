package network

import (
	"net"

	"github.com/thanhfphan/blockchain/network/peer"
)

type Network interface {
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
