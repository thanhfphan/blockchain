package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhfphan/blockchain/config"
	"github.com/thanhfphan/blockchain/network"
	"github.com/thanhfphan/blockchain/network/dialer"
	"github.com/thanhfphan/blockchain/node"
	"github.com/thanhfphan/blockchain/snow/networking/router"
	"github.com/thanhfphan/blockchain/staking"
	"github.com/thanhfphan/blockchain/utils/ips"
)

func main() {
	n := &node.Node{}

	cfg := buildConfig()
	err := n.Initialize(cfg)
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
	cfg.IPConfig = node.IPConfig{
		IPPort: ips.NewDynamicIPPort(net.ParseIP("127.0.0.1"), uint16(4444)),
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
	}

	cfg.ConsensusRouter = &router.ChainRouter{}
	return cfg
}
