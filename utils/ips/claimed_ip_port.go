package ips

import (
	"crypto/x509"

	"github.com/thanhfphan/blockchain/ids"
)

type ClaimedIPPort struct {
	Cert   *x509.Certificate
	IPPort IPPort
	// [Cert]'s signature used to ensure that this IPPort was claimed by the peer
	Signature []byte
	TxID      ids.ID
	TimeStamp uint64
}
