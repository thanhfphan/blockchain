package peer

import (
	"crypto"

	"github.com/thanhfphan/blockchain/utils/ips"
	"github.com/thanhfphan/blockchain/utils/timer"
)

type IPSigner struct {
	ip     ips.DynamicIPPort
	signer crypto.Signer
	clock  timer.Clock
}

func NewIPSigner(ip ips.DynamicIPPort, signer crypto.Signer) *IPSigner {
	return &IPSigner{
		ip:     ip,
		signer: signer,
	}
}

func (s *IPSigner) GetSignedIP() (*ips.SignedIP, error) {
	unsigned := ips.UnsignedIP{
		IP:        s.ip.IPPort(),
		Timestamp: uint64(s.clock.Now().Unix()),
	}

	signedIP, err := unsigned.Sign(s.signer)
	if err != nil {
		return nil, err
	}

	return signedIP, nil
}
