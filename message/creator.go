package message

type Creator interface {
	OutboundMsgBuilder
	InboundMsgBuilder
}

type creator struct {
	OutboundMsgBuilder
	InboundMsgBuilder
}

func NewCreator() (Creator, error) {
	builder, err := newMsgBuilder()
	if err != nil {
		return nil, err
	}
	return &creator{
		InboundMsgBuilder:  newInboundBuilder(builder),
		OutboundMsgBuilder: newOutboundBuilder(builder),
	}, nil
}
