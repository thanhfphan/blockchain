package xxx

import "github.com/thanhfphan/blockchain/ids"

type UTXOID struct {
	id ids.ID

	TxID        ids.ID
	OutputIndex uint32
}
