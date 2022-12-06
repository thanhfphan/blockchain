package node

import (
	"fmt"
	"net"
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/thanhfphan/blockchain/chains"
	"github.com/thanhfphan/blockchain/network"
	"github.com/thanhfphan/blockchain/utils"
)

type Node struct {
	Net           network.Network
	ChainsManager chains.Manager
	Config        *Config

	shuttingDown         utils.AtomicBool
	shuttingDownExitCode utils.AtomicInterface
	shutdownOnce         sync.Once
}

func (n *Node) Initialize(config *Config) error {
	fmt.Println("initializing node")
	n.Config = config
	if err := n.initNetworking(); err != nil {
		fmt.Printf("initNetworking err %v\n", err)
		return err
	}
	if err := n.initChainManager(); err != nil {
		fmt.Printf("init chain manager err %v\n", err)
		return err
	}

	n.initChains()

	return nil
}

func (n *Node) ExitCode() int {
	if exitCode, ok := n.shuttingDownExitCode.GetValue().(int); ok {
		return exitCode
	}
	return 0
}

func (n *Node) Shutdown(exitCode int) {
	if !n.shuttingDown.GetValue() {
		n.shuttingDownExitCode.SetValue(exitCode)
	}

	n.shuttingDown.SetValue(true)
	n.shutdownOnce.Do(n.shutdown)
}

func (n *Node) shutdown() {
	fmt.Println("shutting down node")

	fmt.Println("finished shutdown")
}

func (n *Node) initNetworking() error {
	fmt.Println("initializing networking")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 3012))
	if err != nil {
		return err
	}

	n.Net, err = network.New(listener)
	if err != nil {
		log.Warnf("init network failed %v", err)
		return err
	}

	return nil
}

func (n *Node) initChainManager() error {
	fmt.Println("initializing chain manager")
	err := n.Config.ConsensusRouter.Initialize()
	if err != nil {
		return err
	}
	n.ChainsManager = chains.New(chains.ManagerConfig{
		Net:    n.Net,
		Router: n.Config.ConsensusRouter,
	})
	return nil
}

func (n *Node) initChains() {
	fmt.Println("init chains")
	n.ChainsManager.StartChainCreator(&chains.Chainparameters{
		ID:          100,
		GenesisData: []byte("genesis data"),
	})
}
