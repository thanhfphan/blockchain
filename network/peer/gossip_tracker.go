package peer

import (
	"sync"

	"github.com/thanhfphan/blockchain/ids"
)

var _ GossipTracker = (*gossipTracker)(nil)

type GossipTracker interface {
	StartTrackingPeer(peerID ids.NodeID) bool
	StopTrackingPeer(peerID ids.NodeID) bool
	AddValidator(v ids.ValidatorID) bool
	RemoveValidator(id ids.NodeID) bool

	// AddKnown add [txIDs] to the list Tx known by [peerID]
	// Return list known Tx
	AddKnown(peerID ids.NodeID, txIDs []ids.ID) ([]ids.ID, bool)
	// GetUnknown returns the list validator that [peerID] doesn't know about.
	GetUnknown(peerID ids.NodeID) ([]ids.ValidatorID, bool)
}

type gossipTracker struct {
	lock sync.RWMutex
	// FIXME: for simple we user slice. Might change to bitset later
	trackerValidators map[ids.NodeID]map[ids.ID]bool
	validators        []ids.ValidatorID
	txIDsToNodeIDs    map[ids.ID]ids.NodeID
}

func NewGossipTracker() (GossipTracker, error) {
	return &gossipTracker{
		trackerValidators: make(map[ids.NodeID]map[ids.ID]bool),
		txIDsToNodeIDs:    make(map[ids.ID]ids.NodeID),
	}, nil
}
func (g *gossipTracker) StartTrackingPeer(peerID ids.NodeID) bool {
	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.trackerValidators[peerID]; ok {
		return false
	}

	g.trackerValidators[peerID] = make(map[ids.ID]bool)

	return true
}

func (g *gossipTracker) StopTrackingPeer(peerID ids.NodeID) bool {
	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.trackerValidators[peerID]; !ok {
		return false
	}

	delete(g.trackerValidators, peerID)

	return true
}

func (g *gossipTracker) AddValidator(v ids.ValidatorID) bool {
	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.txIDsToNodeIDs[v.TxID]; ok {
		return false
	}

	g.txIDsToNodeIDs[v.TxID] = v.NodeID
	g.validators = append(g.validators, v)

	return true
}

func (g *gossipTracker) RemoveValidator(id ids.NodeID) bool {
	g.lock.Lock()
	defer g.lock.Unlock()

	txNeedRemove := []ids.ID{}
	newV := []ids.ValidatorID{}
	for _, item := range g.validators {
		if item.NodeID != id {
			newV = append(newV, item)
		} else {
			txNeedRemove = append(txNeedRemove, item.TxID)
		}
	}

	if len(newV) != len(g.validators) {
		g.validators = newV
	}
	for _, item := range txNeedRemove {
		delete(g.txIDsToNodeIDs, item)
	}

	return true
}

func (g *gossipTracker) AddKnown(peerID ids.NodeID, txIDs []ids.ID) ([]ids.ID, bool) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.trackerValidators[peerID]; !ok {
		return nil, false
	}

	knowTxIDss := []ids.ID{}
	for _, tx := range txIDs {
		_, ok := g.txIDsToNodeIDs[tx]
		if !ok {
			continue
		}

		g.trackerValidators[peerID][tx] = true

		knowTxIDss = append(knowTxIDss, tx)
	}

	return knowTxIDss, true
}

func (g *gossipTracker) GetUnknown(peerID ids.NodeID) ([]ids.ValidatorID, bool) {
	g.lock.Lock()
	defer g.lock.Unlock()

	mapValidators, ok := g.trackerValidators[peerID]
	if !ok {
		return nil, false
	}

	result := []ids.ValidatorID{}
	for _, item := range g.validators {
		if _, ok := mapValidators[item.TxID]; !ok {
			result = append(result, item)
		}
	}

	return result, true
}
