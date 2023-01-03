package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

const (
	dbPath = "./tmp/blocks"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

// AddBlock adds a new Block to the Blockchain.
func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleErrors(err)

		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)

			return nil
		})

		return err
	})

	HandleErrors(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleErrors(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})

	HandleErrors(err)
}

// InitBlockchain creates a new Blockchain with genesis Block.
func InitBlockchain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	// Opening Database
	db, err := badger.Open(opts)
	HandleErrors(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing Blockchain found...")

			genesis := Genesis()
			fmt.Println("Genesis Proved")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleErrors(err)

			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			HandleErrors(err)

			err = item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)

				return nil
			})
			HandleErrors(err)

			return err
		}
	})

	HandleErrors(err)

	blockchain := Blockchain{
		LastHash: lastHash,
		Database: db,
	}

	return &blockchain
}
