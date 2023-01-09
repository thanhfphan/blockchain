package ips

import "crypto/x509"

type SignedIP struct {
	IP        UnsignedIP
	Signature []byte
}

func (ip *SignedIP) Verify(cert *x509.Certificate) error {
	return cert.CheckSignature(cert.SignatureAlgorithm, ip.IP.bytes(), ip.Signature)
}
