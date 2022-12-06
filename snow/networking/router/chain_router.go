package router

import (
	"context"
	"fmt"
	"sync"

	"github.com/thanhfphan/blockchain/snow/networking/handler"
	"github.com/thanhfphan/blockchain/utils/constants"
)

type ChainRouter struct {
	chains map[int]handler.Handler
	lock   sync.Mutex
}

func (cr *ChainRouter) Initialize() error {
	cr.chains = make(map[int]handler.Handler)
	return nil
}

func (cr *ChainRouter) HandleInbound(ctx context.Context, msg string) {
	// TODO: currently we assume have 1 chain
	chain, ok := cr.chains[constants.FixChainID]
	if !ok {
		fmt.Println("not found chain when HandleInbound in ChainRouter")
		return
	}

	chain.Push(ctx, msg)
}

func (cr *ChainRouter) Shutdown(ctx context.Context) {

}

func (cr *ChainRouter) AddChain(ctx context.Context, chain handler.Handler) {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	cr.chains[constants.FixChainID] = chain
}
