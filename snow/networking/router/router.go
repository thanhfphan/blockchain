package router

type Router interface {
	ExternalHandler
	InternalHandler
	Initialize() error
	// AddChain(ctx context.Context, chain handler)
}

type ExternalHandler interface {
}

type InternalHandler interface {
}
