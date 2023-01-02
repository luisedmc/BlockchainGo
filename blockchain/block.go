package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Block represents each 'item' in the blockchain.
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// Genesis is the first block of the Blockchain.
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// DeriveHash derives a hash from the block's data.
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{
		b.Data,
		b.PrevHash,
	}, []byte{})
	hash := sha256.Sum256(info)

	b.Hash = hash[:]
}

// CreateBlock creates and returns a new Block.
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		Hash:     []byte{},
		Data:     []byte(data),
		PrevHash: prevHash,
	}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

// Serialize serializes the block.
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	HandleErrors(err)

	return result.Bytes()
}

// Deserialize deserializes the block.
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	HandleErrors(err)

	return &block
}
