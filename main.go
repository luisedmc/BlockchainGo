package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// Blockchain
	chain := InitBlockchain()

	chain.AddBlock("First Block after Genesis Block.")
	chain.AddBlock("Second Block after Genesis Block.")
	chain.AddBlock("Third Block after Genesis Block.")

	// Printing Blocks in the Blockchain
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n\n", block.Hash)

		// Proof of Work
		pow := NewProofOfWork(block)
		fmt.Printf("Proof of Work: %s\n", strings.Title(strconv.FormatBool(pow.Validate())))
		fmt.Println()
	}

	// Transaction

	// Wallet

}
