package app

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	DbFile              = "thenetwork.db"
	BlocksBucket        = "blocks"
	GenesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(DbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func NewBlockchain(address string) *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

func CreateBlockchain(address string) *Blockchain {
	return &Blockchain{
		Tip: nil,
		DB:  nil,
	}
}
