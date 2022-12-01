package runner

import (
	"fmt"
	"os"

	"github.com/thanhfphan/blockchain/app"
	"github.com/thanhfphan/blockchain/app/process"
	"golang.org/x/term"
)

func Run() {
	nodeApp := process.NewApp()

	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println(process.Header)
	}

	exitCode := app.Run(nodeApp)
	os.Exit(exitCode)
}
