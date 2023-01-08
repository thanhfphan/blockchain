package ips

import "crypto/x509"

type ClaimedIPPort struct {
	Cert   *x509.Certificate
	IPPort IPPort
	// [Cert]'s signature used to ensure that this IPPort was claimed by the peer
	Signature []byte
}
