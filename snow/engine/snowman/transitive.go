package snowman

import (
	"context"
	"fmt"
)

var _ Engine = (*Transitive)(nil)

func New(config Config) (Engine, error) {
	return newTransitive(config)
}

type Transitive struct {
	Config
}

func newTransitive(config Config) (Engine, error) {
	t := &Transitive{
		Config: config,
	}

	return t, nil
}

func (t *Transitive) Start(ctx context.Context, id int) error {

	err := t.Consensus.Initialize()
	if err != nil {
		fmt.Printf("error inti consensus %v\n", err)
		return err
	}
	return nil
}
