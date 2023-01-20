package txs

import (
	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/utils/set"
)

type UnsignedTx interface {
	Bytes() []byte
	InputIDs() set.Set[ids.ID]

	//Outputs []
}
