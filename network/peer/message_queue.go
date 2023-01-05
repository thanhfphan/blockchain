package peer

import (
	"context"
	"fmt"
	"sync"

	"github.com/thanhfphan/blockchain/message"
	"github.com/thanhfphan/blockchain/utils/logging"
)

type MessageQueue interface {
	Push(ctx context.Context, msg message.OutboundMessage) bool
	// Pop block until a message is available and then returns the message
	Pop() (message.OutboundMessage, bool)
	// PopNow will return is there is no messages in queue
	PopNow() (message.OutboundMessage, bool)
	Close()
}

type blockingMessageQueue struct {
	closeOnce sync.Once
	mutex     sync.RWMutex
	closing   chan struct{}

	queue chan message.OutboundMessage
	log   logging.Logger
}

func NewBlockingQueue(bufferSize int, log logging.Logger) MessageQueue {
	return &blockingMessageQueue{
		closing: make(chan struct{}),
		queue:   make(chan message.OutboundMessage, bufferSize),
		log:     log,
	}
}

func (q *blockingMessageQueue) logMessage(msg string) {
	if q.log != nil {
		q.log.Verbof(msg)
	}
}

func (q *blockingMessageQueue) Push(ctx context.Context, msg message.OutboundMessage) bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	ctxDone := ctx.Done()
	select {
	case <-q.closing:
		q.logMessage("Dropping message cause channel close\n")
		return false
	case <-ctxDone:
		q.logMessage("Dropping message cause cancelled context\n")
		return false
	default:
	}

	select {
	case q.queue <- msg:
		return true
	case <-ctxDone:
		q.logMessage("dropping message cause cancelled context\n")
		return false
	case <-q.closing:
		q.logMessage("dropping message cause channel close\n")
		return false
	}
}

// Pop will wait if there is no messages in queue
func (q *blockingMessageQueue) Pop() (message.OutboundMessage, bool) {
	select {
	case msg := <-q.queue:
		return msg, true
	case <-q.closing:
		return nil, false
	}
}

// PopNow will return if there no messaages in queue
func (q *blockingMessageQueue) PopNow() (message.OutboundMessage, bool) {
	select {
	case msg := <-q.queue:
		return msg, true
	default:
		return nil, false
	}
}

func (q *blockingMessageQueue) Close() {
	q.closeOnce.Do(func() {
		close(q.closing)

		q.mutex.Lock()
		defer q.mutex.Unlock()

		for {
			select {
			case msg := <-q.queue:
				q.logMessage(fmt.Sprintf("Dropping message when closing MessageQueue OP %v\n", msg.Op()))
			default:
				return
			}
		}
	})
}
