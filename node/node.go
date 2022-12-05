package node

import (
	"fmt"
	"net"

	"github.com/thanhfphan/blockchain/chains"
	"github.com/thanhfphan/blockchain/network"

	"github.com/labstack/gommon/log"
)

type Node struct {
	Net           network.Network
	ChainsManager chains.Manager
	Config        *Config
}

func (n *Node) Initialize(config *Config) error {
	n.Config = config
	if err := n.initNetworking(); err != nil {
		log.Warnf("initNetworking err %v", err)
		return err
	}
	if err := n.initChainManager(); err != nil {
		log.Warnf("init chain manager err %v", err)
		return err
	}

	n.initChains()

	return nil
}

func (n *Node) Start() error {

	return nil
}

func (n *Node) Stop() error {

	return nil
}

func (n *Node) Shutdown() {

}

func (n *Node) ExitCode() int {
	return 0
}

func (n *Node) initNetworking() error {

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
