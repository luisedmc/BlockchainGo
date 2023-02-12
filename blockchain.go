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

// CreateBlockchain creates the Blockchain Database
func CreateBlockchain(address string) *Blockchain {
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
		log.Panic(err)
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

// ContinueBlockchain keeps the Blockchain running
func ContinueBlockchain(address string) *Blockchain {
	var lastHash []byte

	if !DBExists() {
		fmt.Println("ERROR: Blockchain not found!")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Panic(err)
		}

		err = item.Value(func(val []byte) error {
			lastHash = append(lastHash, val...)

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		return err
	})

	chain := Blockchain{
		LastHash: lastHash,
		Database: db,
	}

	return &chain
}

// AddBlock adds a new Block to the Blockchain
func (chain *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	// To add a new Block, we get the last block hash from the Database to use it to mine a new Block hash
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

	newBlock := CreateBlock(transactions, lastHash)

	// Updating the Database with the newBlock
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
				if out.CanBeUnlockedOutput(address) {
					unspentTXOs = append(unspentTXOs, *tx)
				}
			}

			if !tx.IsCoinBase() {
				for _, txin := range tx.Inputs {
					if txin.CanBeUnlockedInput(address) {
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

// FindUTXO returns only the unspent transactions outputs
func (chain *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	// Iterate through every unspent transaction and getting only the outputs
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlockedOutput(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs finds all unspent outputs and ensure that they store enough value to make a transaction
func (chain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTXs := chain.FindUnspentTransactions(address)
	accumulated := 0

	// Iterate through all unspent transactions and accumulate their values
Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlockedOutput(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

// NewTransaction creates a new transaction
func NewTransaction(from, to string, amount int, chain *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enought funds!")
	}

	// List of inputs
	for txID, outs := range validOutputs {
		txID, err := hex.DecodeString(txID)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{
				ID:        txID,
				Output:    out,
				Signature: from,
			}
			inputs = append(inputs, input)
		}
	}

	// List of outputs
	outputs = append(outputs, TXOutput{amount, to})

	// Amount stored in account > Amount sent in the transaction
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	// Create a new transaction
	tx := Transaction{
		ID:      nil,
		Inputs:  inputs,
		Outputs: outputs,
	}
	tx.SetTransactionID()

	return &tx
}

// Iterator retuns a Blockchain iterat
func (chain *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{
		CurrentHash: chain.LastHash,
		Database:    chain.Database,
	}

	return bci
}
