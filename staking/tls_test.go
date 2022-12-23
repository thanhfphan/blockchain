package staking

import (
	"crypto"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thanhfphan/blockchain/utils/hashing"
)

func Test_NewTLSCert(t *testing.T) {
	cert, err := NewTLSCert()
	require.NoError(t, err)

	msg := "hom nay la thu 7"
	msgHash := hashing.ComputeHash256([]byte(msg))

	sig, err := cert.PrivateKey.(crypto.Signer).Sign(rand.Reader, msgHash, crypto.SHA256)
	require.NoError(t, err)

	err = cert.Leaf.CheckSignature(cert.Leaf.SignatureAlgorithm, []byte(msg), sig)
	require.NoError(t, err)
}
