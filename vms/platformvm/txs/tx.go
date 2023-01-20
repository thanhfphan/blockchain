package txs

import (
	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/vms/components/xxx"
)

type Tx struct {
	id    ids.ID
	bytes []byte

	Unsigned UnsignedTx
}

func (tx *Tx) Bytes() []byte {
	return tx.bytes
}

func (tx *Tx) ID() ids.ID {
	return tx.id
}

func (tx *Tx) UTXOs() []*xxx.UTXO {
	// TODO: implement
	return nil
}

// TODO: implement
func (tx *Tx) Sign() error {
	return nil
}
