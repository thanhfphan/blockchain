package network

import (
	"crypto/tls"
	"time"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/utils/ips"
)

type Config struct {
	TLSConfig *tls.Config
	MyNodeID  ids.NodeID
	IPPort    ips.DynamicIPPort
	TimeoutConfig
}

type TimeoutConfig struct {
	PongTimeout time.Duration
}
