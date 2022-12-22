package node

import (
	"crypto/tls"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/network"
	"github.com/thanhfphan/blockchain/snow/networking/router"
	"github.com/thanhfphan/blockchain/utils/ips"
)

type Config struct {
	ConsensusRouter router.Router
	IPConfig
	StakingConfig
	BootstrapConfig
	NetworkConfig network.Config
}

type IPConfig struct {
	IPPort ips.DynamicIPPort
}

type StakingConfig struct {
	StakingTLSCert tls.Certificate
	//TODO:
	// StakingSigningKey
}

type BootstrapConfig struct {
	BoostrapIDs []ids.NodeID
	BoostrapIPs []ips.IPPort
}
