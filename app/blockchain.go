package app

import "sync"

type Blockchain struct {
	Blocks []*Block
	Mutex  *sync.Mutex
}

func CreateBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{NewGenesisBlock()},
		Mutex:  &sync.Mutex{},
	}
}
