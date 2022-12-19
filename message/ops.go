package message

import (
	"errors"
	"fmt"
)

type Op byte

const (
	PingOp Op = iota
	PongOp
	HelloOp
)

var (
	errUnknownMessageType = errors.New("unknown message type")
)

func (op Op) String() string {
	switch op {
	case PingOp:
		return "ping"
	case PongOp:
		return "pong"
	case HelloOp:
		return "hello"
	default:
		return "unknown"
	}
}

func Unwrap(m *Message) (interface{}, error) {
	// This func is prepare when use protobuf
	return m, nil
}

func ToOp(m *Message) (Op, error) {
	switch m.Type {
	case MessageTypePing:
		return PingOp, nil
	case MessageTypePong:
		return PongOp, nil
	case MessageTypeHello:
		return HelloOp, nil
	default:
		return 0, fmt.Errorf("%w: %T", errUnknownMessageType, m)
	}
}
