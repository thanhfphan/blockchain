package handler

import "context"

var _ Handler = (*handler)(nil)

type Handler interface {
	Push(ctx context.Context, msg string)
	Len() int
}

type handler struct {
	messageQueue MessageQueue
}

func New() (Handler, error) {
	h := &handler{}

	return h, nil
}

func (h *handler) Push(ctx context.Context, msg string) {
	h.messageQueue.Push(ctx, msg)
}

func (h *handler) Len() int {
	return h.messageQueue.Len()
}
