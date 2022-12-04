package chains

import (
	"fmt"
	"sync"

	"github.com/thanhfphan/blockchain/network"
	smEngine "github.com/thanhfphan/blockchain/snow/engine/snowman"
	"github.com/thanhfphan/blockchain/snow/networking/handler"
	"github.com/thanhfphan/blockchain/snow/networking/router"
)

var (
	_ Manager = (*manager)(nil)
)

// Manager manager the chains
// It can:
//   - Create a chain
type Manager interface {
	Router() router.Router

	StartChainCreator(parameters *Chainparameters)
	Shutdown()
}

type ManagerConfig struct {
	Net    network.Network
	Router router.Router
}

type manager struct {
	ManagerConfig
	chainsLock sync.Mutex
	chains     map[int64]handler.Handler
}

type Chainparameters struct {
	ID          int64
	GenesisData []byte
}

type chain struct {
	Name    string
	Handler handler.Handler
	Engine  smEngine.Engine
}

func New(config ManagerConfig) Manager {
	return &manager{
		ManagerConfig: config,
		chains:        make(map[int64]handler.Handler),
	}
}

func (m *manager) Router() router.Router {
	return m.ManagerConfig.Router
}

func (m *manager) StartChainCreator(parameters *Chainparameters) {
	m.createChain(parameters)
}

func (m *manager) Shutdown() {
	fmt.Println("shutting down chain manager")
	//TODO: notify
}

func (m *manager) createChain(params *Chainparameters) {
	fmt.Printf("Creating a chain: chainId=%d\n", params.ID)
	chain, err := m.buildChain(params)
	if err != nil {
		fmt.Printf("error creating chain: chainID=%d\n", params.ID)
		return
	}

	m.chainsLock.Lock()
	m.chains[params.ID] = chain.Handler
	m.chainsLock.Unlock()

	// TODO: notify other if needed
}

func (m *manager) buildChain(params *Chainparameters) (*chain, error) {

	chain, err := m.createSnowmanChain()
	if err != nil {
		return nil, err
	}

	return chain, nil
}

func (m *manager) createSnowmanChain() (*chain, error) {

	handler, err := handler.New()
	if err != nil {
		return nil, err
	}

	engine, err := smEngine.New()
	if err != nil {
		return nil, err
	}

	return &chain{
		Engine:  engine,
		Handler: handler,
	}, nil
}
