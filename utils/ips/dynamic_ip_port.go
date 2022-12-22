package ips

import (
	"net"
	"sync"
)

var _ DynamicIPPort = (*dynamicIPPort)(nil)

type DynamicIPPort interface {
	IPPort() IPPort
	SetIP(ip net.IP)
}

type dynamicIPPort struct {
	lock   sync.RWMutex
	ipPort IPPort
}

func NewDynamicIPPort(ip net.IP, port uint16) DynamicIPPort {
	return &dynamicIPPort{
		ipPort: IPPort{
			IP:   ip,
			Port: port,
		},
	}
}

func (i *dynamicIPPort) IPPort() IPPort {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.ipPort
}

func (i *dynamicIPPort) SetIP(ip net.IP) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.ipPort.IP = ip
}
