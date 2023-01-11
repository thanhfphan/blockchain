package vms

import (
	"context"
	"fmt"

	"github.com/thanhfphan/blockchain/ids"
	"github.com/thanhfphan/blockchain/snow"
)

type Factory interface {
	New(*snow.Context) (interface{}, error)
}

type Manager interface {
	GetFactory(vmID ids.ID) (Factory, error)
	RegisterFactory(ctx context.Context, vmID ids.ID, factory Factory) error
}

type manager struct {
	factories map[ids.ID]Factory
}

func NewManager() Manager {
	return &manager{
		factories: make(map[ids.ID]Factory),
	}
}

func (m *manager) GetFactory(vmID ids.ID) (Factory, error) {
	if factory, ok := m.factories[vmID]; ok {
		return factory, nil
	}

	return nil, fmt.Errorf("GetFactory vmId=%s not found", vmID.String())
}

func (m *manager) RegisterFactory(ctx context.Context, vmID ids.ID, factory Factory) error {
	if _, ok := m.factories[vmID]; ok {
		return fmt.Errorf("RegisterFactory failed cause vmID=%s already existed", vmID.String())
	}

	m.factories[vmID] = factory

	return nil
}
