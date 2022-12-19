package message

import "time"

var _ OutboundMsgBuilder = (*outboundMsgBuilder)(nil)

type OutboundMsgBuilder interface {
	Hello() (OutboundMessage, error)
	Ping() (OutboundMessage, error)
	Pong() (OutboundMessage, error)
}

type outboundMsgBuilder struct {
	builder *msgBuilder
}

func newOutboundBuilder(builder *msgBuilder) OutboundMsgBuilder {
	return &outboundMsgBuilder{
		builder: builder,
	}
}

func (mb *outboundMsgBuilder) Hello() (OutboundMessage, error) {
	return mb.builder.createOutbound(&Message{
		Type: MessageTypeHello,
		Message: &MessageHello{
			MyTime:  uint64(time.Now().Unix()),
			Message: "hello there",
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

func (mb *outboundMsgBuilder) Pong() (OutboundMessage, error) {
	return mb.builder.createOutbound(&Message{
		Type: MessageTypePong,
		Message: &MessagePong{
			Message: "pong",
		},
	})
}
