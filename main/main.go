package main

import (
	"github.com/thanhfphan/blockchain/app/runner"
	"github.com/thanhfphan/blockchain/node"
	"github.com/thanhfphan/blockchain/snow/networking/router"
)

func main() {
	nodeConfig := node.Config{
		ConsensusRouter: &router.ChainRouter{},
	}
	runner.Run(nodeConfig)
}
