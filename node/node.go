package node

import (
	"fmt"
	"net"
	"sync"

	"github.com/thanhfphan/blockchain/chains"
	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
	"github.com/thanhfphan/blockchain/network"
	"github.com/thanhfphan/blockchain/network/peer"
	"github.com/thanhfphan/blockchain/utils"
)

type Node struct {
	ID            ids.NodeID
	Net           network.Network
	ChainsManager chains.Manager
	Config        *Config

	shuttingDown         utils.AtomicBool
	shuttingDownExitCode utils.AtomicInterface
	shutdownOnce         sync.Once

	doneShuttingDown sync.WaitGroup
}

func (n *Node) Initialize(config *Config) error {
	n.doneShuttingDown.Add(1)
	n.Config = config
	n.ID = ids.NodeIDFromCert(n.Config.StakingTLSCert.Leaf)

	fmt.Printf("initializing node=%s\n", n.ID.String())

	// TODO: init beacons

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

func (n *Node) Dispatch() error {

	n.Net.Dispatch()
	n.Shutdown(1)

	n.doneShuttingDown.Wait()

	return nil
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

	n.Net.StartClose()

	n.doneShuttingDown.Done()

	fmt.Println("finished shutdown")
}

func (n *Node) initNetworking() error {
	curIPPort := n.Config.IPPort.IPPort()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", curIPPort.Port))
	if err != nil {
		return err
	}

	fmt.Printf("initializing networking at: %s\n", curIPPort.String())

	msgCreator, err := message.NewCreator()
	if err != nil {
		return err
	}

	n.Config.NetworkConfig.TLSConfig = peer.TLSConfig(n.Config.StakingTLSCert)
	n.Config.NetworkConfig.MyNodeID = n.ID
	n.Config.NetworkConfig.IPPort = n.Config.IPPort

	n.Net, err = network.New(&n.Config.NetworkConfig, msgCreator, listener)
	if err != nil {
		fmt.Printf("init network failed %v\n", err)
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
