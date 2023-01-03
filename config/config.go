package config

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/spf13/viper"
	"github.com/thanhfphan/blockchain/genesis"
	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/network"
	"github.com/thanhfphan/blockchain/network/dialer"
	"github.com/thanhfphan/blockchain/node"
	"github.com/thanhfphan/blockchain/snow/networking/router"
	"github.com/thanhfphan/blockchain/staking"
	"github.com/thanhfphan/blockchain/utils/ips"
	"github.com/thanhfphan/blockchain/utils/units"
)

func GetNodeConfig(v *viper.Viper) (node.Config, error) {
	cfg := node.Config{}

	var err error

	cfg.IPConfig, err = getIPConfig(v)
	if err != nil {
		return node.Config{}, fmt.Errorf("get ipConfig err=%v", err)
	}

	cfg.StakingConfig, err = getCertConfig(v)
	if err != nil {
		return node.Config{}, fmt.Errorf("get certConfig err=%v", err)
	}

	cfg.BootstrapConfig, err = getBootstrapConfig()
	if err != nil {
		return node.Config{}, fmt.Errorf("get boostrapConfig err=%v", err)
	}

	cfg.NetworkConfig, err = getNetworkConfig(v)
	if err != nil {
		return node.Config{}, fmt.Errorf("get networkConfig err=%v", err)
	}

	cfg.ConsensusRouter = &router.ChainRouter{}
	return cfg, nil
}

func getIPConfig(v *viper.Viper) (node.IPConfig, error) {

	publicIP := v.GetString(PublicIPKey)
	ip := net.ParseIP(publicIP)
	port := v.GetUint(StakingPortKey)

	return node.IPConfig{
		// use dynamic cause the machine's IP might change
		// TODO: add a job to detect ip change and update to this one
		IPPort: ips.NewDynamicIPPort(ip, uint16(port)),
	}, nil
}

func getTLSCert(v *viper.Viper) (tls.Certificate, error) {
	cert, err := staking.NewTLSCert()
	if err != nil {
		return tls.Certificate{}, err
	}
	return *cert, nil
}

func getCertConfig(v *viper.Viper) (node.StakingConfig, error) {
	cfg := node.StakingConfig{}

	cert, err := getTLSCert(v)
	if err != nil {
		return node.StakingConfig{}, err
	}
	cfg.StakingTLSCert = cert

	return cfg, nil
}

func getBootstrapConfig() (node.BootstrapConfig, error) {
	config := node.BootstrapConfig{}

	genesisNodes := genesis.GetGenesisNodes()
	for _, item := range genesisNodes {
		addr, err := ips.ToIPPort(item.IP)
		if err != nil {
			return node.BootstrapConfig{}, err
		}
		config.BootstrapIPs = append(config.BootstrapIPs, addr)

		nodeID, err := ids.NodeIDFromString(item.NodeID)
		if err != nil {
			return node.BootstrapConfig{}, err
		}
		config.BootstrapIDs = append(config.BootstrapIDs, nodeID)
	}

	return config, nil
}

func getNetworkConfig(v *viper.Viper) (network.Config, error) {
	pingFrequency := v.GetDuration(NetworkPingFrequencyKey)
	pingTimeout := v.GetDuration(NetworkPingTimeoutKey)
	dialerTimeout := v.GetDuration(NetworkDialerTimeoutKey)
	config := network.Config{
		TimeoutConfig: network.TimeoutConfig{
			PongTimeout:   pingTimeout,
			PingFrequency: pingFrequency,
		},
		DialerConfig: dialer.Config{
			ConnectionTimeout: dialerTimeout,
		},
		PeerReadBufferSize:  8 * units.KiB,
		PeerWriteBufferSize: 8 * units.KiB,
	}

	return config, nil
}
