package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"runtime"

	"github.com/dgraph-io/badger/v3"
)

// Blockchain interacts with the Database
type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

// InitBlockchain initializes the Blockchain Database
func InitBlockchain(address string) *Blockchain {
	var lastHash []byte

	if DBExists() {
		fmt.Println("Blockchain already exists.")
		runtime.Goexit()
	}

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
		cbtx := NewCoinBaseTX(address, genesisData)

		genesis := NewGenesis(cbtx)
		fmt.Println("Genesis created.")

		err = txn.Set(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err

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

// FindUnspentTransactions finds all transactions outputs that weren't referenced in any inputs
func (chain *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXOs []Transaction
	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIDx, out := range tx.Outputs {
				// Checking if an output was already spent
				if spentTXOs[txID] != nil {
					for _, spentOutIDx := range spentTXOs[txID] {
						if spentOutIDx == outIDx {
							continue Outputs
						}
					}
				}

				// Getting all unlocked transactions and append it to unspent transactions
				if out.CanUnlockOutput(address) {
					unspentTXOs = append(unspentTXOs, *tx)
				}
			}

			if !tx.IsCoinBase() {
				for _, txin := range tx.Inputs {
					if txin.CanUnlockInput(address) {
						txinID := hex.EncodeToString(txin.ID)
						spentTXOs[txinID] = append(spentTXOs[txinID], txin.Output)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTXOs
}

// Iterator retuns a Blockchain iterat
func (chain *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{
		CurrentHash: chain.LastHash,
		Database:    chain.Database,
	}

	return bci
}
