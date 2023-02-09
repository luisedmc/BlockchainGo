package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block struct holds the block data in the blockchain
type Block struct {
	Data      []byte
	Hash      []byte
	PrevHash  []byte
	Timestamp int64
	Nonce     int
}

// CreateBlock creates a new block
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		Data:      []byte(data),
		Hash:      []byte{},
		PrevHash:  prevHash,
		Timestamp: time.Now().Unix(),
		Nonce:     0,
	}
	// Deriving hash from Proof of Work
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates the first block in the blockchain
func NewGenesisBlock() *Block {
	return CreateBlock("Genesis Block", []byte{})
}

// Serialize serializes a block to be stored in the database
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// Deserialize derializes data and retuns a block
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
