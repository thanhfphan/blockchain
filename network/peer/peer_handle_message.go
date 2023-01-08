package peer

import (
	"encoding/json"

	"github.com/thanhfphan/blockchain/message"
)

func (p *peer) handlePing(_ message.InboundMessage) {
	msg, err := p.MessageCreator.Pong(p.ID().String())
	if err != nil {
		p.log.Errorf("Create pong message failed %v\n", err)
		return
	}

	p.Send(p.onClosingCtx, msg)
}

func (p *peer) handlePong(m message.InboundMessage) {
	// FIXME: use protobuf
	msgRaw, err := json.Marshal(m.Message().Message)
	if err != nil {
		p.log.Errorf("marshal pong message failed")
		return
	}
	var msg message.MessagePong
	if err = json.Unmarshal(msgRaw, &msg); err != nil {
		p.log.Errorf("unmarshal pong message failed %v", err)
		return
	}

	p.log.Debugf("receive pong from %s", p.ID().String())
}

func (p *peer) handleHello(m message.InboundMessage) {
	// FIXME: use protobuf
	msgRaw, err := json.Marshal(m.Message().Message)
	if err != nil {
		p.log.Errorf("marshal hello message failed")
		return
	}
	var msg message.MessageHello
	if err = json.Unmarshal(msgRaw, &msg); err != nil {
		p.log.Errorf("unmarshal hello message failed %v", err)
		return
	}
	p.log.Debugf("receive hello from %s", p.ID().String())
}

func (p *peer) handlePeerList(m message.InboundMessage) {
	// FIXME: use protobuf
	msgRaw, err := json.Marshal(m.Message().Message)
	if err != nil {
		p.log.Errorf("marshal peer list message failed")
		return
	}
	var msg message.MessagePeerList
	if err = json.Unmarshal(msgRaw, &msg); err != nil {
		p.log.Errorf("unmarshal peer list message failed %v", err)
		return
	}
	p.log.Debugf("receive peer listfrom %s", p.id.String())
}
