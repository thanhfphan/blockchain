package dialer

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/thanhfphan/blockchain/utils/ips"
)

var _ Dialer = (*dialer)(nil)

type Dialer interface {
	Dial(ctx context.Context, ip ips.IPPort) (net.Conn, error)
}

type Config struct {
	ConnectionTimeout time.Duration
}

type dialer struct {
	dialer  net.Dialer
	network string
}

func NewDialer(network string, cfg Config) Dialer {
	return &dialer{
		dialer:  net.Dialer{Timeout: cfg.ConnectionTimeout},
		network: network,
	}
}

func (d *dialer) Dial(ctx context.Context, ip ips.IPPort) (net.Conn, error) {
	fmt.Printf("dialing to ip: %s\n", ip.String())

	conn, err := d.dialer.DialContext(ctx, d.network, ip.String())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
