package peer

import (
	"bufio"
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
)

var (
	_              Peer = (*peer)(nil)
	PeerBufferSize      = 16
)

type Peer interface {
	ID() ids.NodeID
	Send(ctx context.Context, message message.OutboundMessage) bool
	StartClose()
}

type peer struct {
	*Config
	id           ids.NodeID
	conn         net.Conn
	cert         *x509.Certificate
	messageQueue MessageQueue

	numExecuting        int64
	startClosingOnce    sync.Once
	onClosingCtx        context.Context
	conClosingCtxCancel func()
	// handle err when close
	onClose chan struct{}
}

func Start(
	config *Config,
	conn net.Conn,
	cert *x509.Certificate,
	nodeID ids.NodeID,
	msgQueue MessageQueue) Peer {
	onClosingCtx, onClosingCtxCancel := context.WithCancel(context.Background())
	p := &peer{
		Config:              config,
		id:                  nodeID,
		conn:                conn,
		cert:                cert,
		messageQueue:        msgQueue,
		onClosingCtx:        onClosingCtx,
		conClosingCtxCancel: onClosingCtxCancel,
		// onClose:             make(chan struct{}),
		numExecuting: 3,
	}

	// number goroutine = numExecuting
	go p.readMessages()
	go p.writeMessages()
	go p.sendNetworkMessages()

	return p
}

func (p *peer) ID() ids.NodeID {
	return p.id
}

func (p *peer) StartClose() {
	p.startClosingOnce.Do(func() {
		if err := p.conn.Close(); err != nil {
			fmt.Printf("failed to close connection node=%s", p.ID().String())
		}

		p.messageQueue.Close()
		p.conClosingCtxCancel()
	})
}

func (p *peer) Send(ctx context.Context, msg message.OutboundMessage) bool {
	return p.messageQueue.Push(ctx, msg)
}

func (p *peer) readMessages() {
	reader := bufio.NewReader(p.conn)

	for {
		msgBytes, err := io.ReadAll(reader)
		if err != nil {
			fmt.Printf("reading message in nodeId=%d failed %v\n", p.id, err)
			return
		}

		if err := p.onClosingCtx.Err(); err != nil {
			return
		}

		msg, err := p.MessageCreator.Parse(msgBytes, p.id)
		if err != nil {
			fmt.Printf("failed to parse message %v, err=%v", string(msgBytes), err)
			continue
		}

		p.handle(msg)
	}
}

func (p *peer) handle(msg message.InboundMessage) {
	fmt.Printf("handle msg=%s\n", msg)
	switch msg.Op() {
	case message.PingOp:
		p.handlePing(msg)
		return
	case message.PongOp:
		p.handlePong(msg)
		return
	}

	// TODO: handle in consensus level
	fmt.Printf("receive unknown message %v\n", msg)
}

func (p *peer) writeMessages() {
	writer := bufio.NewWriter(p.conn)

	for {
		msg, ok := p.messageQueue.PopNow()
		if ok {
			p.writeMessage(writer, msg)
			continue
		}

		if err := writer.Flush(); err != nil {
			fmt.Printf("failed to flush writer %v\n", err)
			return
		}

		msg, ok = p.messageQueue.Pop()
		if !ok {
			return
		}

		p.writeMessage(writer, msg)
	}
}

func (p *peer) handlePing(m message.InboundMessage) {

	msg, err := p.MessageCreator.Pong()
	if err != nil {
		fmt.Printf("create pong message failed %v\n", err)
		return
	}
	p.Send(p.onClosingCtx, msg)
}

func (p *peer) handlePong(msg message.InboundMessage) {
	fmt.Println("receive pong message")
}

func (p *peer) writeMessage(writer *bufio.Writer, msg message.OutboundMessage) {
	msgBytes := msg.Bytes()

	if err := p.conn.SetWriteDeadline(time.Now().Add(p.PongTimeout)); err != nil {
		fmt.Printf("set writeDeadLine failed %v\n", err)
		return
	}
	if _, err := io.Copy(writer, bytes.NewReader(msgBytes)); err != nil {
		fmt.Printf("error writing message %v\n", err)
		return
	}
}

func (p *peer) close() {
	if atomic.AddInt64(&p.numExecuting, -1) != 0 {
		return
	}

	p.Network.Disconnected(p.id)
	// close(p.onClose)
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
		case <-sendPingsTicker.C:
			pingMessage, err := p.Config.MessageCreator.Ping()
			if err != nil {
				fmt.Printf("create PING message failed %v\n", err)
				return
			}
			p.Send(p.onClosingCtx, pingMessage)
		case <-p.onClosingCtx.Done():
			return
		}
	}
}
