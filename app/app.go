package app

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type IApp interface {
	Start() error

	Stop() error

	ExitCode() (int, error)
}

func Run(app IApp) int {
	if err := app.Start(); err != nil {
		return 1
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	signal.Notify(signals, syscall.SIGTERM)

	var eg errgroup.Group
	eg.Go(func() error {
		for range signals {
			return app.Stop()
		}

		return nil
	})

	exitCode, err := app.ExitCode()
	signal.Stop(signals)
	close(signals)
	if eg.Wait() != nil || err != nil {
		return 1
	}

	return exitCode
}
