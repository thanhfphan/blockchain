package ids

import (
	"fmt"
	"sync"
)

type Aliaser interface {
	Alias(id ID, alias string) error
	RemoveAlias(id ID)

	Lookup(alias string) (ID, error)
	Aliases(id ID) ([]string, error)
}

type aliaser struct {
	lock    sync.RWMutex
	aliases map[ID][]string
	dealias map[string]ID
}

func NewAliaser() Aliaser {
	return &aliaser{
		aliases: make(map[ID][]string),
		dealias: make(map[string]ID),
	}
}

func (a *aliaser) Alias(id ID, alias string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, exist := a.dealias[alias]; exist {
		return fmt.Errorf("%s alias already used", alias)
	}

	a.dealias[alias] = id
	a.aliases[id] = append(a.aliases[id], alias)
	return nil
}

func (a *aliaser) RemoveAlias(id ID) {
	a.lock.Lock()
	defer a.lock.Unlock()

	aliases := a.aliases[id]
	delete(a.aliases, id)
	for _, item := range aliases {
		delete(a.dealias, item)
	}
}

func (a *aliaser) Lookup(alias string) (ID, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if id, ok := a.dealias[alias]; ok {
		return id, nil
	}

	return IDEmpty, fmt.Errorf("not found ID with alias=%s", alias)
}

func (a *aliaser) Aliases(id ID) ([]string, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.aliases[id], nil
}
