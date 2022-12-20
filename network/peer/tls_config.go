package peer

import "crypto/tls"

func TLSConfig(cert tls.Certificate) *tls.Config {
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAnyClientCert,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS13,
	}
}
