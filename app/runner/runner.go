package runner

import (
	"fmt"
	"os"

	"github.com/thanhfphan/blockchain/app"
	"github.com/thanhfphan/blockchain/app/process"
	"github.com/thanhfphan/blockchain/node"
	"golang.org/x/term"
)

func Run(nodeConfig node.Config) {
	nodeApp := process.NewApp(nodeConfig)

	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println(process.Header)
	}

	exitCode := app.Run(nodeApp)
	os.Exit(exitCode)
}
