package node

import (
	"fmt"
	"net"
	"sync"

	"github.com/thanhfphan/blockchain/chains"
	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
	"github.com/thanhfphan/blockchain/network"
	"github.com/thanhfphan/blockchain/network/dialer"
	"github.com/thanhfphan/blockchain/network/peer"
	"github.com/thanhfphan/blockchain/utils"
	"github.com/thanhfphan/blockchain/utils/constants"
	"github.com/thanhfphan/blockchain/utils/logging"
)

type Node struct {
	ID            ids.NodeID
	Net           network.Network
	ChainsManager chains.Manager
	Config        *Config
	log           logging.Logger

	shuttingDown         utils.AtomicBool
	shuttingDownExitCode utils.AtomicInterface
	shutdownOnce         sync.Once

	doneShuttingDown sync.WaitGroup
}

func (n *Node) Initialize(config *Config, log logging.Logger) error {
	n.doneShuttingDown.Add(1)
	n.Config = config
	n.log = log
	n.ID = ids.NodeIDFromCert(n.Config.StakingTLSCert.Leaf)

	n.log.Infof("Starting %s\n", n.ID.String())

	// TODO: init beacons

	if err := n.initNetworking(); err != nil {
		n.log.Errorf("Initialize networking error=%v\n", err)
		return err
	}
	if err := n.initChainManager(); err != nil {
		n.log.Errorf("Initialize chain manager error=%v\n", err)
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

	for i, peerIP := range n.Config.BootstrapIPs {
		n.Net.ManuallyTrack(n.Config.BootstrapIDs[i], peerIP)
	}

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
	n.log.Infof("Shutting down %s", n.ID.String())

	n.Net.StartClose()

	n.doneShuttingDown.Done()

	n.log.Infof("Finished shutdown %s", n.ID.String())
}

func (n *Node) initNetworking() error {
	curIPPort := n.Config.IPPort.IPPort()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", curIPPort.Port))
	if err != nil {
		return err
	}

	n.log.Infof("Initializing networking at: %s\n", curIPPort.String())

	msgCreator, err := message.NewCreator()
	if err != nil {
		return err
	}

	n.Config.NetworkConfig.TLSConfig = peer.TLSConfig(n.Config.StakingTLSCert)
	n.Config.NetworkConfig.MyNodeID = n.ID
	n.Config.NetworkConfig.IPPort = n.Config.IPPort

	n.Net, err = network.New(
		&n.Config.NetworkConfig,
		n.log,
		msgCreator,
		listener,
		dialer.NewDialer(constants.NetworkType, n.Config.NetworkConfig.DialerConfig),
	)
	if err != nil {
		n.log.Errorf("Initialize networking error=%v\n", err)
		return err
	}

	return nil
}

func (n *Node) initChainManager() error {
	n.log.Infof("Initializing chain manager")
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
	n.log.Infof("Initializing chains")
	n.ChainsManager.StartChainCreator(&chains.Chainparameters{
		ID:          100,
		GenesisData: []byte("genesis data"),
	})
}
