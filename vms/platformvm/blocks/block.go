package blocks

import (
	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/vms/platformvm/txs"
)

type Block interface {
	ID() ids.ID
	Parent() ids.ID
	Bytes() []byte
	Height() uint64

	Txs() []txs.Tx
}
