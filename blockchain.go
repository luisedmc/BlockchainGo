package main

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v3"
)

// Blockchain interacts with the Database
type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

// Path to save blocks in the Database
const (
	dbPath = "./tmp/blocks"
)

// InitBlockchain initializes the Blockchain Database
func InitBlockchain() *Blockchain {
	var lastHash []byte

	// Defining options for the Database
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	// Opening BadgerDB with options
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		// Checking if already exists a blockchain in the database by trying to find a last hash (lh)
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("Blockchain not found in the database.")

			genesis := NewGenesisBlock()
			fmt.Println("Genesis created.")

			// Saving Genesis Block
			err = txn.Set(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Fatal(err)
			}

			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash

			return err

		} else {
			// Retrieving data
			fmt.Println("Blockchain found in the database.")

			item, err := txn.Get([]byte("lh"))
			if err != nil {
				log.Fatal(err)
			}

			err = item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)

				return nil
			})
			if err != nil {
				log.Fatal(err)
			}

			return err
		}
	})
	if err != nil {
		log.Panic(err)
	}

	// Defining the Blockchain
	blockchain := Blockchain{
		LastHash: lastHash,
		Database: db,
	}

	return &blockchain
}

// AddBlock adds a new Block to the Blockchain
func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// To add a new Block, we get the last block hash from the Database to use it to mine a new Block hash
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Fatal(err)
		}

		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)

			return nil
		})

		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	newBlock := CreateBlock(data, lastHash)

	// Updating the Database with the newBlock
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Fatal(err)
		}

		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}
