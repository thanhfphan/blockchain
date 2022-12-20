package sender

import (
	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
)

type ExternalSender interface {
	Send(msg message.OutboundMessage, nodeIDs []ids.NodeID) []ids.NodeID
	Gossip(msg message.OutboundMessage, numPeersToSend int) []ids.NodeID
}
