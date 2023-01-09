package peer

import (
	"crypto"
	"time"

	"github.com/thanhfphan/blockchain/utils/ips"
)

type IPSigner struct {
	ip     ips.DynamicIPPort
	signer crypto.Signer
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
		Timestamp: uint64(time.Now().Unix()), //FIXME: replace with clock mockable
	}

	signedIP, err := unsigned.Sign(s.signer)
	if err != nil {
		return nil, err
	}

	return signedIP, nil
}
