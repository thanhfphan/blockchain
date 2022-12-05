package snowman

import "context"

type Decidable interface {
	ID() int
	Accept(context.Context) error
	Reject(context.Context) error
	Status() Status
}
