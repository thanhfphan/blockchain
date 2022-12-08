package ids

import (
	"fmt"

	"github.com/thanhfphan/blockchain/utils/ips"
)

type NodeID int

func (id NodeID) String() string {
	return fmt.Sprintf("NodeID-%d", id)
}

func NodeIDFromIPPort(ip ips.IPPort) (NodeID, error) {
	return NodeID(ip.Port), nil
}

func (id NodeID) Equal(other NodeID) bool {
	return id == other
}
