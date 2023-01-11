package snow

import (
	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/utils/logging"
)

type Context struct {
	ChainID ids.ID
	NodeID  ids.NodeID

	Log logging.Logger
}

type ConsensusContext struct {
	*Context
}
