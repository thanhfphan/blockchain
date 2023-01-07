package message

import (
	"fmt"

	"github.com/thanhfphan/blockchain/ids"
)

var _ OutboundMsgBuilder = (*outboundMsgBuilder)(nil)

type OutboundMsgBuilder interface {
	Hello(nodeID ids.NodeID) (OutboundMessage, error)
	Ping() (OutboundMessage, error)
	Pong(msg string) (OutboundMessage, error)
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
