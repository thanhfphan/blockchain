package message

import (
	"github.com/thanhfphan/blockchain/utils/ips"
)

var _ OutboundMsgBuilder = (*outboundMsgBuilder)(nil)

type OutboundMsgBuilder interface {
	Hello(
		helloTime uint64,
		ip ips.IPPort,
		signTime uint64,
		sig []byte,
	) (OutboundMessage, error)
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

func (mb *outboundMsgBuilder) Hello(
	helloTime uint64,
	ip ips.IPPort,
	signedTime uint64,
	signature []byte,

) (OutboundMessage, error) {
	return mb.builder.createOutbound(&Message{
		Type: MessageTypeHello,
		Message: &MessageHello{
			HelloTime:  helloTime,
			IPAddress:  ip.IP.To16(),
			IPPort:     ip.Port,
			Signature:  signature,
			SignedTime: signedTime,
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
			IPAddress: p.IPPort.IP.To16(),
			IPPort:    p.IPPort.Port,
			Signature: p.Signature,
			TxID:      p.TxID[:],
			Timestamp: p.TimeStamp,
		}
	}
	return mb.builder.createOutbound(&Message{
		Type: MessageTypePeerList,
		Message: &MessagePeerList{
			PeerList: peerList,
		},
	})
}
