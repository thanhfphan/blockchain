package ips

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

type IPPort struct {
	IP   net.IP `json:"ip"`
	Port uint16 `json:"port"`
}

func (ip IPPort) Equal(other IPPort) bool {
	return ip.Port == other.Port && ip.IP.Equal(other.IP)
}

func (ip IPPort) String() string {
	return net.JoinHostPort(ip.IP.String(), fmt.Sprintf("%d", ip.Port))
}

func ToIPPort(str string) (IPPort, error) {
	host, portStr, err := net.SplitHostPort(str)
	if err != nil {
		return IPPort{}, nil
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return IPPort{}, nil
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return IPPort{}, errors.New("cant part ip from host")
	}

	return IPPort{
		IP:   ip,
		Port: uint16(port),
	}, nil
}
