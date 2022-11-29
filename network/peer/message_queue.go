package peer

import "context"

type MessageQueue interface {
	Push(ctx context.Context, msg string) bool
	Pop() (string, bool)
}
