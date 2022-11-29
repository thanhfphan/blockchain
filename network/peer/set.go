package peer

var _Set = (*set)(nil)

type Set interface {
	Add(peer Peer)
	GetByID(nodeID int) (Peer, bool)
	Remove(nodeID int)
	Len() int
}

type set struct {
	peersMap   map[int]int // nodeID -> peer's index in peersSlice
	peersSlice []Peer
}

func NewSet() Set {
	return &set{
		peersMap: make(map[int]int),
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

func (s *set) GetByID(nodeID int) (Peer, bool) {
	index, ok := s.peersMap[nodeID]
	if !ok {
		return nil, false
	}

	return s.peersSlice[index], true
}

func (s *set) Remove(nodeID int) {
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
