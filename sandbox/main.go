package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhfphan/blockchain/config"
	"github.com/thanhfphan/blockchain/genesis"
	"github.com/thanhfphan/blockchain/network"
	"github.com/thanhfphan/blockchain/network/dialer"
	"github.com/thanhfphan/blockchain/node"
	"github.com/thanhfphan/blockchain/snow/networking/router"
	"github.com/thanhfphan/blockchain/staking"
	"github.com/thanhfphan/blockchain/utils/ips"
	"github.com/thanhfphan/blockchain/utils/logging"
	"github.com/thanhfphan/blockchain/utils/units"
)

func main() {
	n := &node.Node{}
	cfg := buildConfig()

	logFactory := logging.NewFactory(cfg.LoggingConfig)
	log, err := logFactory.Make("sandbox")
	if err != nil {
		logFactory.Close()
		os.Exit(1)
	}

	err = n.Initialize(cfg, log)
	if err != nil {
		panic(err)
	}

	go func() {
		err := n.Dispatch()
		fmt.Printf("Dispath err=%v", err)
	}()

	waitToSignal()
	n.Shutdown(0)
}

func waitToSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()
	<-done
}

func buildConfig() *node.Config {
	cfg := &node.Config{}
	ipPort, err := ips.ToIPPort(genesis.GetGenesisNodes()[0].IP)
	if err != nil {
		panic(err)
	}
	cfg.IPConfig = node.IPConfig{
		IPPort: ips.NewDynamicIPPort(ipPort.IP, ipPort.Port),
	}

	cert, err := staking.LoadTLSCertFromFiles("./nodes/key1.key", "./nodes/cert1.crt")
	if err != nil {
		panic(err)
	}
	cfg.StakingTLSCert = *cert

	cfg.NetworkConfig = network.Config{
		TimeoutConfig: network.TimeoutConfig{
			PongTimeout:   config.DefaultPingTimeout,
			PingFrequency: config.DefaultPingFrequency,
		},
		DialerConfig: dialer.Config{
			ConnectionTimeout: config.DefaultDialerTimeout,
		},
		PeerReadBufferSize:  8 * units.KiB,
		PeerWriteBufferSize: 8 * units.KiB,
		PeerGossipConfig: network.PeerGossipConfig{
			PeerListNumberValidator: config.DefaultNumberValidatorOfPeer,
			PeerListGossipFrequency: config.DefaultPeerListFrequency,
			PeerListGossipSize:      config.DefaultPeerListGossipSize,
		},
	}

	cfg.LoggingConfig = logging.Config{
		LogLevel:   logging.Verbo,
		LoggerName: "sandbox",
	}

	cfg.ConsensusRouter = &router.ChainRouter{}

	return cfg
}
