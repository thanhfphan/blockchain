package platformvm

import (
	"context"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/snow/consensus/snowman"
	"github.com/thanhfphan/blockchain/snow/engine/snowman/block"
)

var (
	_ block.ChainVM = (*VM)(nil)
)

type VM struct {
}

func (v *VM) Initialize(ctx context.Context) error {
	return nil
}

func (v *VM) Shutdown(context.Context) error {

	return nil
}

func (v *VM) BuildBlock(context.Context) (snowman.Block, error) {

	return nil, nil
}

func (v *VM) SetPreference(ctx context.Context, bloclID ids.ID) error {

	return nil
}

func (v *VM) LastAccepted(context.Context) (ids.ID, error) {

	return ids.IDEmpty, nil
}

func (v *VM) GetBlock(ctx context.Context, blockID ids.ID) (snowman.Block, error) {
	return nil, nil
}

func (v *VM) ParseBlock(ctx context.Context, blockBytes []byte) (snowman.Block, error) {
	return nil, nil
}
