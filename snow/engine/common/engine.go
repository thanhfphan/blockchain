package common

import "context"

type Engine interface {
	Start(ctx context.Context, id int) error
}
