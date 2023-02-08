package main

import (
	"bytes"
	"math/big"
)

// targetBits define the defficulty of mining a block
const targetBits = 16

// ProofOfWork struct holds a block and target, that is the requirement for mining a block
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork creates and returns a new proof of work
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{
		block:  b,
		target: target,
	}

	return pow
}

// InitData merges the block fields with target and nonce
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}
