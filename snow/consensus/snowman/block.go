package snowman

import (
	"time"

	"github.com/thanhfphan/blockchain/snow/consensus/snowball"
)

type Block interface {
	Decidable
	Parent() int
	Byte() []byte
	Height() int
	Timestamp() time.Time
}

type snowmanBlock struct {
	block    Block
	children map[int]Block
	sb       snowball.Consensus
}

func (n *snowmanBlock) AddChild(child Block) {
	childID := child.ID()

	if n.sb == nil {
		n.sb = &snowball.Tree{}
		n.sb.Initialize(childID)
		n.children = make(map[int]Block)
	} else {
		n.sb.Add(childID)
	}

	n.children[childID] = child
}

func (n *snowmanBlock) Accepted() bool {
	if n.block == nil {
		return false
	}

	return n.block.Status() == Accepted
}
