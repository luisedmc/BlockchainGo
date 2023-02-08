package main

import (
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
