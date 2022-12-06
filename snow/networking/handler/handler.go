package handler

import (
	"context"
	"fmt"

	"github.com/thanhfphan/blockchain/snow/engine/common"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	Start(ctx context.Context)
	Push(ctx context.Context, msg string)
	Len() int

	SetConsensus(engine common.Engine)
	GetConsensus() common.Engine

	SetBootstrapper(engine common.BootstrapableEngine)
	GetBootstrapper() common.Engine
}

type handler struct {
	messageQueue MessageQueue
	engine       common.Engine
	bootstrapper common.BootstrapableEngine
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

func (h *handler) SetConsensus(engine common.Engine) {
	h.engine = engine
}

func (h *handler) GetConsensus() common.Engine {
	return h.engine
}

func (h *handler) SetBootstrapper(engine common.BootstrapableEngine) {
	h.bootstrapper = engine
}

func (h *handler) GetBootstrapper() common.Engine {
	return h.bootstrapper
}

func (h *handler) Start(ctx context.Context) {
	err := h.bootstrapper.Start(ctx, 0)
	if err != nil {
		fmt.Printf("start bootstrapper error %v\n", err)
		return
	}
}
