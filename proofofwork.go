package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const maxNonce = math.MaxInt64

// targetBits define the defficulty of mining a block
const targetBits = 16

// ProofOfWork struct holds a block and target, that is the requirement for mining a block
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProofOfWork creates and returns a new proof of work
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{
		Block:  b,
		Target: target,
	}

	return pow
}

// InitData merges the block fields with target and nonce
func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.Data,
			pow.Block.PrevHash,
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run runs a proof of work
func (pow *ProofOfWork) Run() ([]byte, int) {
	var intHash big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.Block.Data)

	for nonce < maxNonce {
		// Prepare data
		data := pow.PrepareData(nonce)

		// Hash it with SHA-256 hash
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		// Convert the hash to a big integer
		intHash.SetBytes(hash[:])

		// 4. Compare the integer with the target
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println()

	return hash[:], nonce
}

// Validate validates a proof of work and returns the validation answer
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.PrepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
