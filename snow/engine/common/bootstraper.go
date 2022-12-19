package common

import "context"

var (
	_ Bootstrapper = (*bootstrapper)(nil)
)

type Bootstrapper interface {
	Startup(ctx context.Context) error
}

type bootstrapper struct {
	Config
}

func NewCommonBootstrapper(config Config) Bootstrapper {
	return &bootstrapper{
		Config: config,
	}
}

func (b *bootstrapper) Startup(ctx context.Context) error {
	// TODO: get nodes by make http request server, then use it to bootstrap here
	return nil
}
