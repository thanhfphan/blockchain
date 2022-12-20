package network

import (
	"crypto/tls"

	"github.com/thanhfphan/blockchain/ids"
)

type Config struct {
	TLSConfig *tls.Config
	MyNodeID  ids.NodeID
}
