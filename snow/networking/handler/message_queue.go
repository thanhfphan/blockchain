package handler

import (
	"context"
	"sync"
)

type MessageQueue interface {
	Push(ctx context.Context, msg string)
	Pop() (string, error)
	Len() int
	Shutdown()
}

type messageQueue struct {
	unprocessMessages []string
	cond              *sync.Cond
}

func NewMessageQueue() MessageQueue {
	return &messageQueue{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (m *messageQueue) Push(ctx context.Context, msg string) {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()

	m.unprocessMessages = append(m.unprocessMessages, msg)

	m.cond.Signal()
}

func (m *messageQueue) Pop() (string, error) {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()

	for {
		if m.Len() != 0 {
			break
		}

		m.cond.Wait()
	}

	message := m.unprocessMessages[0]
	if cap(m.unprocessMessages) == 1 {
		m.unprocessMessages = nil
	} else {
		m.unprocessMessages = m.unprocessMessages[1:]
	}

	return message, nil
}

func (m *messageQueue) Len() int {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()

	return len(m.unprocessMessages)
}

func (m *messageQueue) Shutdown() {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()

	m.cond.Broadcast()
}
