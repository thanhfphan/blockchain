package snowman

import "context"

type Handler interface {
}

type Engine interface {
	Handler

	Start(ctx context.Context, reqID int) error
}

func New() (Engine, error) {
	return nil, nil
}
