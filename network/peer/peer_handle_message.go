package peer

import (
	"encoding/json"
	"net"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
	"github.com/thanhfphan/blockchain/utils/ips"
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
	if p.gotHello.GetValue() {
		p.log.Warnf("duplicate hello message")
		return
	}

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

	p.ip = &ips.SignedIP{
		IP: ips.UnsignedIP{
			IP: ips.IPPort{
				IP:   net.IP(msg.IPAddress),
				Port: msg.IPPort,
			},
			Timestamp: msg.HelloTime,
		},
		Signature: msg.Signature,
	}

	if err := p.ip.Verify(p.cert); err != nil {
		p.log.Debugf("verify signature when do hello failed %v", err)
		p.StartClose()
		return
	}

	p.gotHello.SetValue(true)

	otherPeers, err := p.Network.Peers(p.id)
	if err != nil {
		p.log.Errorf("get peers to gossip failed %v", err)
		return
	}

	peerListMsg, err := p.MessageCreator.PeerList(otherPeers)
	if err != nil {
		p.log.Errorf("create peer ist message failed %v", err)
		return
	}

	if !p.Send(p.onClosingCtx, peerListMsg) {
		p.log.Errorf("send peer list message failed %v", err)
		return
	}
}

func (p *peer) handlePeerList(m message.InboundMessage) {
	// FIXME: use protobuf
	if !p.finishedHandshake.GetValue() {
		if !p.gotHello.GetValue() {
			p.log.Debugf("not saw any hello from the peer")
			return
		}
		p.Network.Connected(p.id)
		p.finishedHandshake.SetValue(true)

	}

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

	discoveredTxIDs := make([]ids.ID, 0, len(msg.PeerList))
	for _, item := range msg.PeerList {
		if len(item.TxID) > 0 {
			txID, err := ids.ToID(item.TxID)
			if err != nil {
				p.log.Errorf("receive peerlist with invalid data %v", err)
				p.StartClose()
				return
			}
			discoveredTxIDs = append(discoveredTxIDs, txID)
		}
	}

	trackedTxIDs, ok := p.gossipTracker.AddKnown(p.id, discoveredTxIDs)
	if !ok {
		p.log.Errorf("add known peer failed")
		return
	}

	if len(trackedTxIDs) == 0 {
		p.log.Debugf("skipping peerListAck cause there were no trackedIDs")
		return
	}

	ackMsg, err := p.MessageCreator.PeerListAck(trackedTxIDs)
	if err != nil {
		p.log.Errorf("create peerListAck message failed %v", err)
		return
	}

	if !p.Send(p.onClosingCtx, ackMsg) {
		p.log.Debugf("send peerListAck message failed")
	}

	p.log.Debugf("receive peer listfrom %v", msg.PeerList)
}

func (p *peer) handlePeerListAck(m message.InboundMessage) {
	// FIXME: use protobuf
	msgRaw, err := json.Marshal(m.Message().Message)
	if err != nil {
		p.log.Errorf("marshal peer list ack message failed")
		return
	}
	var msg message.MessagePeerListAck
	if err = json.Unmarshal(msgRaw, &msg); err != nil {
		p.log.Errorf("unmarshal peer list message ack failed %v", err)
		return
	}

	txIDs := make([]ids.ID, len(msg.TxIDs))
	for i, item := range msg.TxIDs {
		txID, err := ids.ToID(item)
		if err != nil {
			p.log.Errorf("parse ID peer list ack failed %v", err)
			p.StartClose()
			return
		}

		txIDs[i] = txID
	}

	if _, ok := p.gossipTracker.AddKnown(p.id, txIDs); !ok {
		p.log.Warnf("update known peers failed from peerListAck")
	}

	p.log.Debugf("handle peerListAck with txIDs=%v", txIDs)
}
