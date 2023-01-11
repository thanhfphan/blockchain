package snowman

import (
	"github.com/thanhfphan/blockchain/snow/consensus/snowman"
	"github.com/thanhfphan/blockchain/snow/engine/snowman/block"
)

type Config struct {
	Consensus snowman.Consensus
	VM        block.ChainVM
}
