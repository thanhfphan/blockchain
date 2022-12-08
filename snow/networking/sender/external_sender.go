package sender

import "github.com/thanhfphan/blockchain/ids"

type ExternalSender interface {
	Send(msg string, nodeIDs []ids.NodeID) []ids.NodeID
	Gossip(msg string, numPeersToSend int) []ids.NodeID
}
