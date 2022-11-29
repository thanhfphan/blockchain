package runner

import (
	"os"

	"github.com/thanhfphan/blockchain/app"
	"github.com/thanhfphan/blockchain/app/process"
)

func Run() {
	nodeApp := process.NewApp()

	exitCode := app.Run(nodeApp)

	os.Exit(exitCode)
}
