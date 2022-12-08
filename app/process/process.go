package process

import (
	"fmt"
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/thanhfphan/blockchain/app"
	"github.com/thanhfphan/blockchain/node"
)

const (
	Header = `__________.__                 __   _________ .__           .__        
\______   \  |   ____   ____ |  | _\_   ___ \|  |__ _____  |__| ____  
 |    |  _/  |  /  _ \_/ ___\|  |/ /    \  \/|  |  \\__  \ |  |/    \ 
 |    |   \  |_(  <_> )  \___|    <\     \___|   Y  \/ __ \|  |   |  \
 |______  /____/\____/ \___  >__|_ \\______  /___|  (____  /__|___|  /
        \/                 \/     \/       \/     \/     \/        \/ `
)

var (
	_ app.IApp = (*Process)(nil)
)

type Process struct {
	node   *node.Node
	exitWG sync.WaitGroup
	config node.Config
}

func NewApp(config node.Config) app.IApp {
	return &Process{
		node:   &node.Node{},
		config: config,
	}
}

func (p *Process) Start() error {
	if err := p.node.Initialize(&p.config); err != nil {
		log.Warnf("init node failed %v", err)
		return err
	}

	p.exitWG.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("caught panic", r)
			}
			p.exitWG.Done()
		}()

		err := p.node.Dispatch()
		fmt.Printf("dispatch returned err=%v\n", err)
	}()

	return nil
}

func (p *Process) Stop() error {
	p.node.Shutdown(0)
	return nil
}

func (p *Process) ExitCode() (int, error) {
	p.exitWG.Wait()
	return p.node.ExitCode(), nil
}
