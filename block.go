package main

// Block represents each 'item' in the blockchain
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}
