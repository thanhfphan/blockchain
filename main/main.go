package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/thanhfphan/blockchain/app/runner"
	"github.com/thanhfphan/blockchain/config"
)

func main() {
	fs := config.BuildFlagSet()
	v, err := config.BuildViper(fs, os.Args[1:])
	if errors.Is(err, pflag.ErrHelp) {
		os.Exit(0)
	}
	if err != nil {
		fmt.Printf("config flag failed: %s\n", err)
		os.Exit(1)
	}

	cfg, err := config.GetNodeConfig(v)
	if err != nil {
		fmt.Printf("get nodeConfig with err=%v", err)
		os.Exit(1)
	}

	runner.Run(cfg)
}
