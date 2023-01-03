package peer

import (
	"bufio"
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
	"github.com/thanhfphan/blockchain/utils/constants"
	"github.com/thanhfphan/blockchain/utils/wrappers"
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
	reader := bufio.NewReaderSize(p.conn, p.Config.ReadBufferSize)
	msgLenBytes := make([]byte, wrappers.IntLen)
	for {
		// TODO: implement ping-pong
		// if err := p.conn.SetReadDeadline(time.Now().Add(p.PongTimeout)); err != nil {
		// 	fmt.Println("error setting the connection read timeout - readMessages")
		// 	return
		// }

		// read message length
		if _, err := io.ReadFull(reader, msgLenBytes); err != nil {
			fmt.Printf("reading message length in nodeId=%s failed %v\n", p.id.String(), err)
			return
		}

		msgLen, err := readMsgLen(msgLenBytes, constants.DefaultMaxMessageSize)
		if err != nil {
			fmt.Println("parse message length failed - readMsgLen")
			return
		}

		if err := p.onClosingCtx.Err(); err != nil {
			return
		}

		msgBytes := make([]byte, msgLen)
		// read the message
		if _, err := io.ReadFull(reader, msgBytes); err != nil {
			fmt.Printf("reading message in nodeId=%d failed %v\n", p.id, err)
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

func (p *peer) writeMessages() {
	writer := bufio.NewWriterSize(p.conn, p.Config.WriteBufferSize)

	msg, err := p.MessageCreator.Hello()
	if err != nil {
		fmt.Printf("create msg hello failed %v\n", err)
		return
	}

	p.writeMessage(writer, msg)

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

func (p *peer) writeMessage(writer *bufio.Writer, msg message.OutboundMessage) {
	msgBytes := msg.Bytes()

	if err := p.conn.SetWriteDeadline(time.Now().Add(p.PongTimeout)); err != nil {
		fmt.Printf("set writeDeadLine failed %v\n", err)
		return
	}

	msgLen := uint32(len(msgBytes))
	msgLenBytes, err := writeMsgLen(msgLen, constants.DefaultMaxMessageSize)
	if err != nil {
		fmt.Printf("writeMsgLen got err %v\n", err)
		return
	}

	var buf net.Buffers = [][]byte{msgLenBytes[:], msgBytes}
	if _, err := io.CopyN(writer, &buf, int64(wrappers.IntLen+msgLen)); err != nil {
		fmt.Printf("error writing message %v\n", err)
		return
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
