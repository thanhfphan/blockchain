package chains

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/message"
	"github.com/thanhfphan/blockchain/network"
	"github.com/thanhfphan/blockchain/snow"
	smConsensus "github.com/thanhfphan/blockchain/snow/consensus/snowman"
	"github.com/thanhfphan/blockchain/snow/engine/common"
	smEngine "github.com/thanhfphan/blockchain/snow/engine/snowman"
	"github.com/thanhfphan/blockchain/snow/engine/snowman/block"
	smBootstrap "github.com/thanhfphan/blockchain/snow/engine/snowman/bootstrap"
	"github.com/thanhfphan/blockchain/snow/networking/handler"
	"github.com/thanhfphan/blockchain/snow/networking/router"
	"github.com/thanhfphan/blockchain/utils/logging"
	"github.com/thanhfphan/blockchain/vms"
)

var (
	_ Manager = (*manager)(nil)
)

// Manager manager the chains run on the nodes
type Manager interface {
	Router() router.Router

	StartChainCreator(parameters *Chainparameters)
	Shutdown()
}

type ManagerConfig struct {
	Net         network.Network
	Router      router.Router
	StakingCert tls.Certificate // use to sign block
	MsgCreator  message.OutboundMsgBuilder
	NodeID      ids.NodeID
	// callback to shutdown node has an error
	ShutdownNodeFunc func(exitCode int)
	VMManager        vms.Manager
}

type manager struct {
	ManagerConfig
	// TODO: might want to use factory to create new logger instead reused logger of the node
	log        logging.Logger
	chainsLock sync.Mutex
	chains     map[ids.ID]handler.Handler
}

func New(log logging.Logger, config ManagerConfig) Manager {
	return &manager{
		log:           log,
		ManagerConfig: config,
		chains:        make(map[ids.ID]handler.Handler),
	}
}

func (m *manager) Router() router.Router {
	return m.ManagerConfig.Router
}

func (m *manager) StartChainCreator(parameters *Chainparameters) {
	m.createChain(parameters)
}

func (m *manager) Shutdown() {
	m.log.Infof("Shutting down chain manager")
}

func (m *manager) createChain(params *Chainparameters) {
	m.log.Infof("Creating a chain: chainID=%s\n, VMID=%s", params.ID.String(), params.VMID.String())
	chain, err := m.buildChain(params)
	if err != nil {
		m.log.Fatalf("Error creating chain: chainID=%s\n, VMI=%s", params.ID.String(), params.ID.String())
		m.ShutdownNodeFunc(1)
		return
	}

	m.chainsLock.Lock()
	m.chains[params.ID] = chain.handler
	m.chainsLock.Unlock()

	m.ManagerConfig.Router.AddChain(context.Background(), chain.handler)

	chain.handler.Start(context.TODO())
}

func (m *manager) buildChain(params *Chainparameters) (*chain, error) {
	vmFactory, err := m.VMManager.GetFactory(params.VMID)
	if err != nil {
		return nil, fmt.Errorf("get VMFactory failed %v", err)
	}

	ctx := &snow.ConsensusContext{
		Context: &snow.Context{
			ChainID: params.ID,
			NodeID:  m.NodeID,
			Log:     m.log, // TODO: create other logger
		},
	}

	vm, err := vmFactory.New(ctx.Context)
	if err != nil {
		return nil, fmt.Errorf("factory new failed %v", err)
	}

	var chain *chain
	// we might want to create custom VM here
	// for now: just hard code
	switch vm := vm.(type) {
	case block.ChainVM:
		chain, err = m.createSnowmanChain(vm, params.GenesisData)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Not support custom chain yet")
	}

	return chain, nil
}

func (m *manager) createSnowmanChain(vm block.ChainVM, genesisData []byte) (*chain, error) {

	if err := vm.Initialize(context.TODO()); err != nil {
		return nil, err
	}

	handler, err := handler.New()
	if err != nil {
		return nil, err
	}

	//FIXME: Remove harcode
	commonCfg := common.Config{
		SampleK: 1,
		Alpha:   1,
	}

	bootstrapperConfig := smBootstrap.Config{
		Config: commonCfg,
	}

	bootstrapper, err := smBootstrap.New(context.TODO(), bootstrapperConfig)
	if err != nil {
		return nil, err
	}
	handler.SetBootstrapper(bootstrapper)

	consensus := &smConsensus.Topological{}
	engineConfig := smEngine.Config{
		Consensus: consensus,
	}

	engine, err := smEngine.New(engineConfig)
	if err != nil {
		return nil, err
	}
	handler.SetConsensus(engine)
	return &chain{
		engine:  engine,
		handler: handler,
	}, nil
}
