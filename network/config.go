package network

import (
	"crypto"
	"crypto/tls"
	"time"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/network/dialer"
	"github.com/thanhfphan/blockchain/utils/ips"
)

type Config struct {
	TLSConfig           *tls.Config
	TLSSignIPKey        crypto.Signer
	MyNodeID            ids.NodeID
	IPPort              ips.DynamicIPPort
	DialerConfig        dialer.Config
	PeerReadBufferSize  int
	PeerWriteBufferSize int

	TimeoutConfig
	PeerGossipConfig
}

type TimeoutConfig struct {
	PongTimeout   time.Duration
	PingFrequency time.Duration
}

type PeerGossipConfig struct {
	PeerListNumberValidator uint
	PeerListGossipFrequency time.Duration
	PeerListGossipSize      uint
}
