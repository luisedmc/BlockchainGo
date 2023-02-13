package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"

	"github.com/luisedmc/BlockchainGo/tx"
)

// Block struct holds the Block data in the blockchain
type Block struct {
	Hash         []byte
	PrevHash     []byte
	Transactions []*tx.Transaction
	Timestamp    int64
	Nonce        int
}

// CreateBlock creates a new Block
func CreateBlock(txs []*tx.Transaction, prevHash []byte) *Block {
	block := &Block{
		Hash:         []byte{},
		PrevHash:     prevHash,
		Transactions: txs,
		Timestamp:    time.Now().Unix(),
		Nonce:        0,
	}
	// Deriving hash from Proof of Work
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates the first Block in the Blockchain
func NewGenesis(coinbase *tx.Transaction) *Block {
	return CreateBlock([]*tx.Transaction{coinbase}, []byte{})
}

// HashTransations creates a unique hash for transactions
func (b *Block) HashTransactions() []byte {
	var TXHashes [][]byte
	var TXHash [32]byte

	for _, tx := range b.Transactions {
		TXHashes = append(TXHashes, tx.ID)
	}

	TXHash = sha256.Sum256(bytes.Join(TXHashes, []byte{}))

	return TXHash[:]
}

// Serialize serializes a Block to be stored in the Database
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// Deserialize deserializes Block data
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
