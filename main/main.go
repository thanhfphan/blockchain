package main

import (
	"fmt"

	"github.com/thanhfphan/blockchain/app/runner"
	"github.com/thanhfphan/blockchain/config"
)

func main() {
	cfg, err := config.GetNodeConfig()
	if err != nil {
		fmt.Printf("get nodeConfig with err=%v", err)
		return
	}

	runner.Run(cfg)
}
