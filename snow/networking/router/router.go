package router

import (
	"context"

	"github.com/thanhfphan/blockchain/snow/networking/handler"
)

type Router interface {
	ExternalHandler
	InternalHandler
	Initialize() error
	AddChain(ctx context.Context, chain handler.Handler)
}

type ExternalHandler interface {
}

type InternalHandler interface {
}
