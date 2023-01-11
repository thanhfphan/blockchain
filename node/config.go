package node

import (
	"crypto/tls"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/network"
	"github.com/thanhfphan/blockchain/snow/networking/router"
	"github.com/thanhfphan/blockchain/utils/ips"
	"github.com/thanhfphan/blockchain/utils/logging"
)

type Config struct {
	IPConfig
	StakingConfig
	BootstrapConfig

	ConsensusRouter router.Router
	NetworkConfig   network.Config
	LoggingConfig   logging.Config
}

type IPConfig struct {
	IPPort ips.DynamicIPPort
}

type StakingConfig struct {
	StakingTLSCert tls.Certificate
	// StakingSigningKey *bls
}

type BootstrapConfig struct {
	BootstrapIDs []ids.NodeID
	BootstrapIPs []ips.IPPort
}
