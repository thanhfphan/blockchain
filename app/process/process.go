package process

import (
	"github.com/labstack/gommon/log"
	"github.com/thanhfphan/blockchain/app"
	"github.com/thanhfphan/blockchain/node"
)

var (
	_ app.IApp = (*Process)(nil)
)

type Process struct {
	node *node.Node
}

func NewApp() app.IApp {
	return &Process{
		node: &node.Node{},
	}
}

func (p *Process) Start() error {

	if err := p.node.Initialize(); err != nil {
		log.Warnf("init node failed %v", err)
		return err
	}
	return nil
}

func (p *Process) Stop() error {
	p.node.Shutdown()
	return nil
}

func (p *Process) ExitCode() (int, error) {

	return p.node.ExitCode(), nil
}
