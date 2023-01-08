package peer

import (
	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/utils/ips"
)

type Network interface {
	Connected(id ids.NodeID)
	Disconnected(ids ids.NodeID)
	// Peers return peers that [PeerID] might not know about.
	Peers(peerIDs ids.NodeID) ([]ips.ClaimedIPPort, error)
}
