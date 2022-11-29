package peer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
)

var (
	_              Peer = (*peer)(nil)
	PeerBufferSize      = 16
)

type Peer interface {
	ID() int
	Send(ctx context.Context, message string) bool
}

type peer struct {
	id           int
	conn         net.Conn
	messageQueue MessageQueue
}

func Start(conn net.Conn, msgQueue MessageQueue) Peer {

	p := &peer{
		conn:         conn,
		messageQueue: msgQueue,
	}

	go p.readMessages()
	go p.writeMessages()

	return p
}

func (p *peer) ID() int {
	return p.id
}

func (p *peer) Send(ctx context.Context, msg string) bool {
	return p.messageQueue.Push(ctx, msg)
}

func (p *peer) readMessages() {
	reader := bufio.NewReaderSize(p.conn, PeerBufferSize)

	for {
		msgBytes, err := io.ReadAll(reader)
		if err != nil {
			fmt.Printf("reading message in nodeId=%d failed %v\n", p.id, err)
			return
		}

		p.handle(string(msgBytes))
	}
}

func (p *peer) handle(msg string) {
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
