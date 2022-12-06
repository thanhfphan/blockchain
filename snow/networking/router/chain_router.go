package router

import (
	"context"
	"fmt"

	"github.com/thanhfphan/blockchain/snow/networking/handler"
)

type ChainRouter struct {
	chains map[int]handler.Handler
}

func (c *ChainRouter) Initialize() error {
	c.chains = make(map[int]handler.Handler)
	return nil
}

func (c *ChainRouter) HandleInbound(ctx context.Context, msg string) {
	// TODO: currently we assume have 1 chain
	chain, ok := c.chains[0]
	if !ok {
		fmt.Println("not found chain when HandleInbound in ChainRouter")
		return
	}

	chain.Push(ctx, msg)
}
