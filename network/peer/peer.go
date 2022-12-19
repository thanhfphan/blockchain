package peer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
)

var (
	_              Peer = (*peer)(nil)
	PeerBufferSize      = 16
)

type Peer interface {
	ID() ids.NodeID
	Send(ctx context.Context, message string) bool
	StartClose()
}

type peer struct {
	*Config
	id           ids.NodeID
	conn         net.Conn
	messageQueue MessageQueue

	startClosingOnce    sync.Once
	onClosingCtx        context.Context
	conClosingCtxCancel func()
}

func Start(config *Config, conn net.Conn, nodeID ids.NodeID, msgQueue MessageQueue) Peer {
	onClosingCtx, onClosingCtxCancel := context.WithCancel(context.Background())
	p := &peer{
		Config:              config,
		id:                  nodeID,
		conn:                conn,
		messageQueue:        msgQueue,
		onClosingCtx:        onClosingCtx,
		conClosingCtxCancel: onClosingCtxCancel,
	}

	go p.readMessages()
	go p.writeMessages()

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

func (p *peer) Send(ctx context.Context, msg string) bool {
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
}

func (p *peer) writeMessages() {
	writer := bufio.NewWriter(p.conn)

	for {
		msg, ok := p.messageQueue.Pop()
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

func (p *peer) writeMessage(writer *bufio.Writer, msg string) {
	if _, err := io.Copy(writer, strings.NewReader(msg)); err != nil {
		fmt.Printf("error writing message %v\n", err)
		return
	}
}
