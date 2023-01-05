package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dgraph-io/badger/v3"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

// AddBlock adds a new Block to the Blockchain.
func (chain *Blockchain) AddBlock(transactions []*Transaction) {
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

	newBlock := CreateBlock(transactions, lastHash)

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
func InitBlockchain(address string) *Blockchain {
	var lastHash []byte

	if DBExists() {
		fmt.Println("Blockchain already exists.")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	// Opening Database
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinBaseTX(address, genesisData)
		genesis := Genesis(cbtx)

		fmt.Println("Genesis created!")

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

	blockchain := Blockchain{
		LastHash: lastHash,
		Database: db,
	}

	return &blockchain
}

// ContinueBlockchain continues an existing Blockchain.
func ContinueBlockchain(address string) *Blockchain {
	if !DBExists() {
		fmt.Println("No existing Blockchain.")
		runtime.Goexit()
	}

	var lastHash []byte

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
			lastHash = append([]byte{}, val...)
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

// FindUnspentTransactions finds and returns all unspent transactions.
func (chain *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTxOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTxOs[txID] != nil {
					for _, spentOut := range spentTxOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if !tx.IsCoinBase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTxOs[inTxID] = append(spentTxOs[inTxID], in.Output)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

// FindUTXO finds and returns all unspent transaction outputs.
func (chain *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs finds and returns unspent outputs to be used as inputs in a new transaction.
func (chain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTransactions := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTransactions {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
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

// DBExists checks if a Database exists.
func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
