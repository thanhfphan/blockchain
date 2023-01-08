package message

import (
	"fmt"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/utils/ips"
)

var _ OutboundMsgBuilder = (*outboundMsgBuilder)(nil)

type OutboundMsgBuilder interface {
	Hello(nodeID ids.NodeID) (OutboundMessage, error)
	Ping() (OutboundMessage, error)
	Pong(msg string) (OutboundMessage, error)
	PeerList(peerIPs []ips.ClaimedIPPort) (OutboundMessage, error)
}

type outboundMsgBuilder struct {
	builder *msgBuilder
}

func newOutboundBuilder(builder *msgBuilder) OutboundMsgBuilder {
	return &outboundMsgBuilder{
		builder: builder,
	}
}

func (mb *outboundMsgBuilder) Hello(nodeID ids.NodeID) (OutboundMessage, error) {
	return mb.builder.createOutbound(&Message{
		Type: MessageTypeHello,
		Message: &MessageHello{
			Message: fmt.Sprintf("hello %s", nodeID.String()),
		},
	})
}

func (mb *outboundMsgBuilder) Ping() (OutboundMessage, error) {
	return mb.builder.createOutbound(&Message{
		Type: MessageTypePing,
		Message: &MessagePing{
			Message: "ping",
		},
	})
}

func (mb *outboundMsgBuilder) Pong(msg string) (OutboundMessage, error) {
	return mb.builder.createOutbound(&Message{
		Type: MessageTypePong,
		Message: &MessagePong{
			Message: msg,
		},
	})
}
func (mb *outboundMsgBuilder) PeerList(peerIPs []ips.ClaimedIPPort) (OutboundMessage, error) {
	peerList := make([]*PeerList, len(peerIPs))
	for i, p := range peerIPs {
		peerList[i] = &PeerList{
			Cert:      p.Cert.Raw,
			IP:        p.IPPort.IP.To16(),
			Port:      p.IPPort.Port,
			Signature: p.Signature,
		}
	}
	return mb.builder.createOutbound(&Message{
		Type: MessageTypePeerList,
		Message: &MessagePeerList{
			PeerList: peerList,
		},
	})
}
