package blockchain

type Blockchain struct {
	Blocks []*Block
}

// AddBlock adds a block to the blockchain
func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	newBlock := CreateBlock(data, prevBlock.Hash)

	chain.Blocks = append(chain.Blocks, newBlock)
}

// Genesis is the first block of the Blockchain
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// InitBlockchain starts the Blockchain
func InitBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{Genesis()},
	}
}
