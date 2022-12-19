package message

import (
	"encoding/json"

	"github.com/thanhfphan/blockchain/ids"
)

var (
	_ InboundMessage  = (*inboundMessage)(nil)
	_ OutboundMessage = (*outboundMessage)(nil)
)

// ******************** Inbound stuffs **********************
type InboundMessage interface {
	NodeID() ids.NodeID
	Op() Op
	Message() any
}

type inboundMessage struct {
	nodeID  ids.NodeID
	op      Op
	message any
}

func (i *inboundMessage) NodeID() ids.NodeID {
	return i.nodeID
}

func (i *inboundMessage) Op() Op {
	return i.op
}

func (i *inboundMessage) Message() any {
	return i.message
}

// ********************* Outbound stuffs **********************
type OutboundMessage interface {
	Op() Op
	Bytes() []byte
}

type outboundMessage struct {
	op    Op
	bytes []byte
}

func (i *outboundMessage) Op() Op {
	return i.op
}

func (i *outboundMessage) Bytes() []byte {
	return i.bytes
}

// ********************* msgBuilder stuffs ********************
type msgBuilder struct {
}

func newMsgBuilder() (*msgBuilder, error) {
	mb := &msgBuilder{}

	return mb, nil
}

func (mb *msgBuilder) marsharl(msg *Message) ([]byte, error) {
	rawBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return rawBytes, nil
}

func (mb *msgBuilder) unmarshal(bytes []byte) (*Message, error) {
	m := new(Message)
	if err := json.Unmarshal(bytes, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (mb *msgBuilder) parseInbound(bytes []byte, nodeID ids.NodeID) (*inboundMessage, error) {
	m, err := mb.unmarshal(bytes)
	if err != nil {
		return nil, err
	}

	op, err := ToOp(m)
	if err != nil {
		return nil, err
	}

	msg, err := Unwrap(m)
	if err != nil {
		return nil, err
	}

	return &inboundMessage{
		nodeID:  nodeID,
		op:      op,
		message: msg,
	}, nil
}

func (mb *msgBuilder) createOutbound(m *Message) (*outboundMessage, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	op, err := ToOp(m)
	if err != nil {
		return nil, err
	}

	return &outboundMessage{
		op:    op,
		bytes: b,
	}, nil
}
