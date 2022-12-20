package ids

import (
	"crypto/x509"

	"github.com/thanhfphan/blockchain/utils/hashing"
)

const (
	NodeIDPrefix = "NodeID-"
)

type NodeID ID

func (id NodeID) Bytes() []byte {
	return id[:]
}

func (id NodeID) String() string {
	return NodeID(id).PrefixedString(NodeIDPrefix)
}

func (id NodeID) PrefixedString(prefix string) string {
	return prefix + id.String()
}

func ToNodeID(bytes []byte) (NodeID, error) {
	nodeID, err := ToID(bytes)
	return NodeID(nodeID), err
}

func NodeIDFromCert(cert *x509.Certificate) NodeID {
	return hashing.ComputeHash256Array(cert.Raw)
}

func NodeIDFromString(idStr string) (NodeID, error) {
	id, err := IDFromPrefixedString(idStr, NodeIDPrefix)
	if err != nil {
		return NodeID{}, err
	}

	return NodeID(id), nil
}
