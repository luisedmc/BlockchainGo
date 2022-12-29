package main

import (
	"bytes"
	"crypto/sha256"
)

// DeriveHash derives a hash from the block's data
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{
		b.Data,
		b.PrevHash,
	}, []byte{})
	hash := sha256.Sum256(info)

	b.Hash = hash[:]
}

func main() {

}
