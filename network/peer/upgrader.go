package peer

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"

	"github.com/thanhfphan/blockchain/ids"
)

var (
	_ Upgrader = (*tlsServerUpgrader)(nil)
	_ Upgrader = (*tlsClientUpgrader)(nil)
)

type Upgrader interface {
	Upgrader(net.Conn) (ids.NodeID, net.Conn, *x509.Certificate, error)
}

func connToIDAndCert(conn *tls.Conn) (ids.NodeID, net.Conn, *x509.Certificate, error) {
	if err := conn.Handshake(); err != nil {
		return ids.NodeID{}, nil, nil, err
	}

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return ids.NodeID{}, nil, nil, errors.New("peer does not have certificate")
	}

	peerCert := state.PeerCertificates[0]
	return ids.NodeIDFromCert(peerCert), conn, peerCert, nil
}

// tlsServerUpgrader
type tlsServerUpgrader struct {
	config *tls.Config
}

func NewServerUpgrader(config *tls.Config) Upgrader {
	return tlsServerUpgrader{
		config: config,
	}
}

func (t tlsServerUpgrader) Upgrader(conn net.Conn) (ids.NodeID, net.Conn, *x509.Certificate, error) {
	return connToIDAndCert(tls.Server(conn, t.config))
}

// tlsClientUpgrader
type tlsClientUpgrader struct {
	config *tls.Config
}

func NewClientUpgrader(config *tls.Config) Upgrader {
	return tlsClientUpgrader{
		config: config,
	}
}

func (t tlsClientUpgrader) Upgrader(conn net.Conn) (ids.NodeID, net.Conn, *x509.Certificate, error) {
	return connToIDAndCert(tls.Client(conn, t.config))
}
