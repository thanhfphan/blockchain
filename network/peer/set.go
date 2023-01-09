package peer

import (
	"errors"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/utils/sampler"
)

var _Set = (*set)(nil)

type Set interface {
	Add(peer Peer)
	GetByID(nodeID ids.NodeID) (Peer, bool)
	GetByIndex(i int) (Peer, error)
	Remove(nodeID ids.NodeID)
	Len() int
	Sample(n int, precondition func(Peer) bool) []Peer
}

type set struct {
	peersMap   map[ids.NodeID]int // nodeID -> peer's index in peersSlice
	peersSlice []Peer
}

func NewSet() Set {
	return &set{
		peersMap: make(map[ids.NodeID]int),
	}
}

func (s *set) Add(peer Peer) {
	nodeID := peer.ID()
	index, ok := s.peersMap[nodeID]
	if !ok {
		s.peersMap[nodeID] = len(s.peersSlice)
		s.peersSlice = append(s.peersSlice, peer)
	} else {
		s.peersSlice[index] = peer
	}
}
func (s *set) GetByID(nodeID ids.NodeID) (Peer, bool) {
	index, ok := s.peersMap[nodeID]
	if !ok {
		return nil, false
	}

	return s.peersSlice[index], true
}

func (s *set) GetByIndex(index int) (Peer, error) {
	if index < 0 || index >= len(s.peersSlice) {
		return nil, errors.New("out of range peersSlice")
	}

	return s.peersSlice[index], nil
}

func (s *set) Remove(nodeID ids.NodeID) {
	index, ok := s.peersMap[nodeID]
	if !ok {
		return
	}

	lastIndex := s.Len() - 1
	lastPeer := s.peersSlice[lastIndex]

	s.peersMap[lastPeer.ID()] = index
	s.peersSlice[index] = lastPeer
	delete(s.peersMap, nodeID)
	s.peersSlice[lastIndex] = nil
	s.peersSlice = s.peersSlice[:lastIndex]
}

func (s *set) Len() int {
	return len(s.peersSlice)
}

func (s *set) Sample(n int, precondition func(Peer) bool) []Peer {
	if n <= 0 {
		return nil
	}

	sampler := sampler.NewUniform()
	sampler.Initialize(uint64(len(s.peersSlice)))

	peers := make([]Peer, 0, n)
	for len(peers) < n {
		index, err := sampler.Next()
		if err != nil {
			break
		}
		peer := s.peersSlice[index]
		if !precondition(peer) {
			continue
		}

		peers = append(peers, peer)
	}

	return peers
}
