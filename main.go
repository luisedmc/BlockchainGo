package main

import "fmt"

func main() {
	// Blockchain
	chain := InitBlockchain()

	chain.AddBlock("First Block after Genesis Block.")
	chain.AddBlock("Second Block after Genesis Block.")
	chain.AddBlock("Third Block after Genesis Block.")

	// Printing Blocks in the Blockchain
	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n\n", block.Hash)
	}

	// Transaction

	// Proof of Work

	// Wallet

}
