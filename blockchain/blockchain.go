package blockchain

import (
	"fmt"
	"log"

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
		if err != nil {
			log.Panic(err)
		}

		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)

			return nil
		})

		return err
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})

	if err != nil {
		log.Panic(err)
	}
}

// InitBlockchain creates a new Blockchain with genesis Block.
func InitBlockchain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	// Opening Database
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing Blockchain found...")

			genesis := Genesis()
			fmt.Println("Genesis Proved")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				log.Panic(err)
			}

			err = item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)

				return nil
			})
			if err != nil {
				log.Panic(err)
			}

			return err
		}
	})
	if err != nil {
		log.Panic(err)
	}

	blockchain := Blockchain{
		LastHash: lastHash,
		Database: db,
	}

	return &blockchain
}

// Iterator iterates over the Blockchain.
func (chain *Blockchain) Iterator() *Blockchain {
	iterator := &Blockchain{
		LastHash: chain.LastHash,
		Database: chain.Database,
	}

	return iterator
}

func (iter *Blockchain) Next() *Block {
	var block *Block
	var blockData []byte

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.LastHash)
		if err != nil {
			log.Panic(err)
		}

		err = item.Value(func(val []byte) error {
			blockData = append([]byte{}, val...)

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		block = Deserialize(blockData)

		return err
	})
	if err != nil {
		log.Panic(err)
	}

	iter.LastHash = block.PrevHash
	return block
}
