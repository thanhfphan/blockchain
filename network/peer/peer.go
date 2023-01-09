package peer

import (
	"bufio"
	"context"
	"crypto/x509"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
	"github.com/thanhfphan/blockchain/utils"
	"github.com/thanhfphan/blockchain/utils/constants"
	"github.com/thanhfphan/blockchain/utils/ips"
	"github.com/thanhfphan/blockchain/utils/logging"
)

var (
	_              Peer = (*peer)(nil)
	PeerBufferSize      = 16
)

type Peer interface {
	ID() ids.NodeID
	Cert() *x509.Certificate
	IP() *ips.SignedIP

	Send(ctx context.Context, message message.OutboundMessage) bool
	StartClose()
	StartSendPeerList()
}

type peer struct {
	*Config
	log  logging.Logger
	id   ids.NodeID
	conn net.Conn
	cert *x509.Certificate
	ip   *ips.SignedIP

	messageQueue  MessageQueue
	gossipTracker GossipTracker

	finishedHandshake   utils.AtomicBool
	gotHello            utils.AtomicBool
	numExecuting        int64
	startClosingOnce    sync.Once
	onClosingCtx        context.Context
	conClosingCtxCancel func()
	// handle err when close
	onClose chan struct{}

	// peerListChan signals that use to send list peer to this peer
	peerListChan chan struct{}
}

func Start(
	config *Config,
	log logging.Logger,
	conn net.Conn,
	cert *x509.Certificate,
	nodeID ids.NodeID,
	msgQueue MessageQueue,
	gossipTracker GossipTracker,
) Peer {
	onClosingCtx, onClosingCtxCancel := context.WithCancel(context.Background())
	p := &peer{
		Config:              config,
		log:                 log,
		id:                  nodeID,
		conn:                conn,
		cert:                cert,
		onClosingCtx:        onClosingCtx,
		conClosingCtxCancel: onClosingCtxCancel,
		// onClose:             make(chan struct{}),
		numExecuting: 3,
		peerListChan: make(chan struct{}, 1), //only send once

		messageQueue:  msgQueue,
		gossipTracker: gossipTracker,
	}

	p.gossipTracker.StartTrackingPeer(nodeID)

	// number goroutine = numExecuting
	go p.readMessages()
	go p.writeMessages()
	go p.sendNetworkMessages()

	return p
}

func (p *peer) ID() ids.NodeID {
	return p.id
}

func (p *peer) IP() *ips.SignedIP {
	return p.ip
}

func (p *peer) Cert() *x509.Certificate {
	return p.cert
}

func (p *peer) StartClose() {
	p.startClosingOnce.Do(func() {
		p.log.Verbof("Closing peer %s", p.ID().String())
		if err := p.conn.Close(); err != nil {
			p.log.Verbof("Failed to close connection node=%s", p.ID().String())
		}

		p.messageQueue.Close()
		p.conClosingCtxCancel()
	})
}

func (p *peer) Send(ctx context.Context, msg message.OutboundMessage) bool {
	return p.messageQueue.Push(ctx, msg)
}

func (p *peer) StartSendPeerList() {
	select {
	case p.peerListChan <- struct{}{}:
	default:
	}
}

func (p *peer) readMessages() {
	defer func() {
		p.StartClose()
		p.close()
	}()

	reader := bufio.NewReaderSize(p.conn, p.Config.ReadBufferSize)
	msgLenBytes := make([]byte, constants.IntLen)
	for {
		if err := p.conn.SetReadDeadline(p.nextTimeout()); err != nil {
			p.log.Verbof("error setting the connection read timeout - readMessages")
			return
		}

		// read message length
		if _, err := io.ReadFull(reader, msgLenBytes); err != nil {
			p.log.Verbof("Reading message length in nodeId=%s failed %v\n", p.id.String(), err)
			return
		}

		msgLen, err := readMsgLen(msgLenBytes, constants.DefaultMaxMessageSize)
		if err != nil {
			p.log.Verbof("Parse message length failed - readMsgLen")
			return
		}

		if err := p.onClosingCtx.Err(); err != nil {
			return
		}

		if err := p.conn.SetReadDeadline(p.nextTimeout()); err != nil {
			p.log.Verbof("error setting the connection read timeout - readMessages")
			return
		}

		msgBytes := make([]byte, msgLen)
		// read the message
		if _, err := io.ReadFull(reader, msgBytes); err != nil {
			p.log.Verbof("Reading message in nodeId=%d failed %v\n", p.id, err)
			return
		}

		msg, err := p.MessageCreator.Parse(msgBytes, p.id)
		if err != nil {
			p.log.Verbof("Failed to parse message %v, err=%v", string(msgBytes), err)
			continue
		}

		p.handle(msg)
	}
}

func (p *peer) writeMessages() {
	defer func() {
		p.StartClose()
		p.close()
	}()

	signedIP, err := p.IPSigner.GetSignedIP()
	if err != nil {
		p.log.Errorf("get signed IP failed %v", err)
		return
	}

	msg, err := p.MessageCreator.Hello(
		uint64(time.Now().Unix()),
		signedIP.IP.IP,
		signedIP.IP.Timestamp,
		signedIP.Signature,
	)
	if err != nil {
		p.log.Errorf("Create msg hello failed %v\n", err)
		return
	}

	writer := bufio.NewWriterSize(p.conn, p.Config.WriteBufferSize)
	p.writeMessage(writer, msg)

	for {
		msg, ok := p.messageQueue.PopNow()
		if ok {
			p.writeMessage(writer, msg)
			continue
		}

		if err := writer.Flush(); err != nil {
			p.log.Verbof("Failed to flush writer %v\n", err)
			return
		}

		msg, ok = p.messageQueue.Pop()
		if !ok {
			return
		}

		p.writeMessage(writer, msg)
	}
}

func (p *peer) writeMessage(writer *bufio.Writer, msg message.OutboundMessage) {
	msgBytes := msg.Bytes()

	if err := p.conn.SetWriteDeadline(p.nextTimeout()); err != nil {
		p.log.Verbof("Set writeDeadLine failed %v\n", err)
		return
	}

	msgLen := uint32(len(msgBytes))
	msgLenBytes, err := writeMsgLen(msgLen, constants.DefaultMaxMessageSize)
	if err != nil {
		p.log.Verbof("writeMsgLen got err %v\n", err)
		return
	}

	var buf net.Buffers = [][]byte{msgLenBytes[:], msgBytes}
	if _, err := io.CopyN(writer, &buf, int64(constants.IntLen+msgLen)); err != nil {
		p.log.Verbof("error writing message %v\n", err)
		return
	}
}

func (p *peer) sendNetworkMessages() {
	sendPingsTicker := time.NewTicker(p.PingFrequency)
	defer func() {
		sendPingsTicker.Stop()
		p.StartClose()
		p.close()
	}()

	for {
		select {
		case <-p.peerListChan:
			peerIPs, err := p.Network.Peers(p.id)
			if err != nil {
				p.log.Errorf("get peers failed %v", err)
				return
			}
			if len(peerIPs) == 0 {
				p.log.Debugf("skipping peerListChan cause found nobody")
				continue
			}

			msg, err := p.MessageCreator.PeerList(peerIPs)
			if err != nil {
				p.log.Errorf("create PeerList message failed %v", err)
				continue
			}

			if !p.Send(p.onClosingCtx, msg) {
				p.log.Debugf("failed to send peer list message")
				continue
			}

		case <-sendPingsTicker.C:
			pingMessage, err := p.MessageCreator.Ping()
			if err != nil {
				p.log.Errorf("Create PING message failed %v\n", err)
				return
			}
			p.Send(p.onClosingCtx, pingMessage)
		case <-p.onClosingCtx.Done():
			return
		}
	}
}

func (p *peer) handle(msg message.InboundMessage) {
	switch msg.Op() {
	case message.PingOp:
		p.handlePing(msg)
		return
	case message.PongOp:
		p.handlePong(msg)
		return
	case message.HelloOp:
		p.handleHello(msg)
		return
	case message.PeerListOp:
		p.handlePeerList(msg)
		return
	}

	if !p.finishedHandshake.GetValue() {
		p.log.Warnf("dropping message because not finish handshake yet")
		return
	}

	// TODO: handle in consensus level
	p.log.Verbof("receive unknown message %v\n", msg)
}

func (p *peer) close() {
	if atomic.AddInt64(&p.numExecuting, -1) != 0 {
		return
	}

	p.gossipTracker.StopTrackingPeer(p.id)
	p.Network.Disconnected(p.id)
	// close(p.onClose)
}

func (p *peer) nextTimeout() time.Time {
	// FIXME: replace time.Now with clock can be mock for test
	return time.Now().Add(p.PongTimeout)
}
