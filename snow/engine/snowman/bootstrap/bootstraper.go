package bootstrap

import (
	"context"
	"fmt"

	"github.com/thanhfphan/blockchain/snow/engine/common"
)

var (
	_ common.BootstrapableEngine = (*bootstrapper)(nil)
)

type bootstrapper struct {
	Config
	common.Bootstrapper
}

func New(ctx context.Context, config Config) (common.BootstrapableEngine, error) {
	b := &bootstrapper{
		Config:       config,
		Bootstrapper: common.NewCommonBootstrapper(config.Config),
	}
	return b, nil
}

func (b *bootstrapper) Start(ctx context.Context, id int) error {
	fmt.Println("starting bootstrapper in engine.snowman.bootstrap")
	return b.Startup(ctx)
}

func (b *bootstrapper) Clear() error {
	return nil
}
