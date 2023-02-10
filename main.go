package main

import "os"

func main() {
	defer os.Exit(0)

	// Blockchain
	bc := InitBlockchain()
	defer bc.Database.Close()

	// CLI
	cli := CommandLine{
		Blockchain: bc,
	}
	cli.RunCLI()

	// Block

	// Proof of Work

	// Transaction

	// Wallet

}
