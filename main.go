package main

import (
	"github.com/luisedmc/blockgo/blockchain"
	"os"
)

func main() {
	defer os.Exit(0)

	chain := blockchain.InitBlockchain()
	defer chain.Database.Close()

	cli := CommandLine{
		blockchain: chain,
	}

	cli.Run()
}
