package chains

import (
	"github.com/thanhfphan/blockchain/ids"
	smEngine "github.com/thanhfphan/blockchain/snow/engine/snowman"
	"github.com/thanhfphan/blockchain/snow/networking/handler"
)

type Chainparameters struct {
	// ID of the chain
	ID          ids.ID
	GenesisData []byte
	// The ID of  the VM this chain is running
	VMID ids.ID
}

type chain struct {
	name    string
	handler handler.Handler
	engine  smEngine.Engine
}
