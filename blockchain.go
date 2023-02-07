package main

// Blockchain struct holds all the blocks in the blockchain
type Blockchain struct {
	blocks []*Block
}

// InitBlockchain initializes the blockchain
func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

// Genesis creates the first block in the blockchain
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// AddBlock adds a new block to the blockchain
func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}
