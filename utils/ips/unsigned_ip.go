package ips

import (
	"crypto"
	"crypto/rand"
	"encoding/binary"

	"github.com/thanhfphan/blockchain/utils/constants"
	"github.com/thanhfphan/blockchain/utils/hashing"
)

type UnsignedIP struct {
	IP        IPPort
	Timestamp uint64
}

func (ip *UnsignedIP) Sign(signer crypto.Signer) (*SignedIP, error) {
	sig, err := signer.Sign(rand.Reader, hashing.ComputeHash256(ip.bytes()), crypto.SHA256)
	if err != nil {
		return nil, err
	}

	return &SignedIP{
		IP:        *ip,
		Signature: sig,
	}, nil
}

func (ip *UnsignedIP) bytes() []byte {
	bytes := make([]byte, constants.IPLen+constants.LongLen)
	offset := 0
	ipBytes := ip.IP.IP.To16()

	copy(bytes[offset:], ipBytes)
	offset += len(ipBytes)
	binary.BigEndian.PutUint64(bytes[offset:], ip.Timestamp)
	offset += constants.LongLen

	return bytes
}
