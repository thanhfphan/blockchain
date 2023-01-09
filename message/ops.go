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
	PeerListOp
	PeerListAckOp
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
	case PeerListOp:
		return "peer-list"
	case PeerListAckOp:
		return "peer-list-ack"
	default:
		return "unknown"
	}
}

func ToOp(m *Message) (Op, error) {
	switch m.Type {
	case MessageTypePing:
		return PingOp, nil
	case MessageTypePong:
		return PongOp, nil
	case MessageTypeHello:
		return HelloOp, nil
	case MessageTypePeerList:
		return PeerListOp, nil
	case MessageTypePeerAckList:
		return PeerListAckOp, nil
	default:
		return 0, fmt.Errorf("%w: %T", errUnknownMessageType, m)
	}
}
