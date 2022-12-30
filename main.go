package main

import (
	"fmt"
	"strconv"

	"github.com/luisedmc/blockgo/blockchain"
)

func main() {
	chain := blockchain.InitBlockchain()

	chain.AddBlock("First Block after Genesis.")
	chain.AddBlock("Second Block after Genesis.")
	chain.AddBlock("Third Block after Genesis.")

	// Displaying the Blockchain
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		fmt.Println()

		pow := blockchain.NewProof(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
