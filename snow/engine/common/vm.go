package common

import "context"

type VM interface {
	Initialize(ctx context.Context) error
	// Shutdown is called when node is shutting down
	Shutdown(context.Context) error
}
