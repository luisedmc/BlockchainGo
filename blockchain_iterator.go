package main

import (
	"log"

	"github.com/dgraph-io/badger/v3"
)

// BlockchainIterator holds the Current Hash of the Block being iterated and Database to access the Blocks
type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Next will return the next Block to be iterated
func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	var blockData []byte

	err := bci.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(bci.CurrentHash)
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

	bci.CurrentHash = block.PrevHash

	return block
}
