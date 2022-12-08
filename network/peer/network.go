package peer

import "github.com/thanhfphan/blockchain/ids"

type Network interface {
	Pong(id ids.NodeID) (string, error)

	Connected(id ids.NodeID)
	Disconnected(ids ids.NodeID)
}
