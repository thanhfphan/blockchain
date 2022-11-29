package node

import (
	"fmt"
	"net"

	"github.com/thanhfphan/blockchain/network"

	"github.com/labstack/gommon/log"
)

type Node struct {
	Net network.Network
}

func (n *Node) Initialize() error {

	return nil
}

func (n *Node) Start() error {

	return nil
}

func (n *Node) Stop() error {

	return nil
}

func (n *Node) Shutdown() {

}

func (n *Node) ExitCode() int {
	return 0
}

func (n *Node) initNetworking() error {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 3012))
	if err != nil {
		return err
	}

	n.Net, err = network.New(listener)
	if err != nil {
		log.Warnf("init network failed %v", err)
		return err
	}

	return nil
}
