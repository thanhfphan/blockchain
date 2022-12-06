package common

type BootstrapableEngine interface {
	Bootstrapable
	Engine
}

type Bootstrapable interface {
	Clear() error
}
