package snowman

import (
	"time"
)

type Block interface {
	Parent() int
	Byte() []byte
	Height() int
	Timestamp() time.Time
}

type snowmanBlock struct {
	block    Block
	children map[int]Block
}

func (n *snowmanBlock) AddChild(child Block) {
}

func (n *snowmanBlock) Accepted() bool {

	return true
}
