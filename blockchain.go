package main

// Blockchain struct holds all the blocks in the blockchain
type Blockchain struct {
	Blocks []*Block
}

// AddBlock adds a new block to the blockchain
func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

// InitBlockchain initializes the blockchain
func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
