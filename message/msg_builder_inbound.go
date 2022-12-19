package message

import "github.com/thanhfphan/blockchain/ids"

var _ InboundMsgBuilder = (*inboundMsgBuilder)(nil)

type InboundMsgBuilder interface {
	Parse(bytes []byte, nodeID ids.NodeID) (InboundMessage, error)
}

type inboundMsgBuilder struct {
	builder *msgBuilder
}

func newInboundBuilder(builder *msgBuilder) InboundMsgBuilder {
	return &inboundMsgBuilder{
		builder: builder,
	}
}

func (i *inboundMsgBuilder) Parse(bytes []byte, nodeID ids.NodeID) (InboundMessage, error) {
	return i.builder.parseInbound(bytes, nodeID)
}
