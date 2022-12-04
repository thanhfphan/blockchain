package snowman

import (
	"context"

	"github.com/thanhfphan/blockchain/snow/consensus/snowman"
)

type Getter interface {
	GetBlock(ctx context.Context, blockID int) (snowman.Block, error)
}

type Parser interface {
	ParseBlock(ctx context.Context, blockBytes []byte) (snowman.Block, error)
}
