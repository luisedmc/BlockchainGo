package blockchain

import (
	"bytes"
	"crypto/sha256"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

// DeriveHash derives a hash from the block's data
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{
		b.Data,
		b.PrevHash,
	}, []byte{})
	hash := sha256.Sum256(info)

	b.Hash = hash[:]
}

// CreateBlock creates and returns a new Block
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		Hash:     []byte{},
		Data:     []byte(data),
		PrevHash: prevHash,
	}
	block.DeriveHash()

	return block
}
