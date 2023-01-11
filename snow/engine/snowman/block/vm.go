package block

import (
	"context"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/snow/consensus/snowman"
	"github.com/thanhfphan/blockchain/snow/engine/common"
)

type ChainVM interface {
	common.VM
	Getter
	Parse

	BuildBlock(context.Context) (snowman.Block, error)
	SetPreference(ctx context.Context, bloclID ids.ID) error
	LastAccepted(context.Context) (ids.ID, error)
}

type Getter interface {
	GetBlock(ctx context.Context, blockID ids.ID) (snowman.Block, error)
}

type Parse interface {
	ParseBlock(ctx context.Context, blockBytes []byte) (snowman.Block, error)
}
