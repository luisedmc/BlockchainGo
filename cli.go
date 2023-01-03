package main

import (
	"flag"
	"fmt"
	"github.com/luisedmc/blockgo/blockchain"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
	blockchain *blockchain.Blockchain
}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage: ")
	fmt.Println("add -block BLOCK_DATA - add a block to the blockchain.")
	fmt.Println("print - Print all blocks in the chain.")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) AddBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Block added !")
}

func (cli *CommandLine) PrintChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		fmt.Println()

		pow := blockchain.NewProof(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	addBlockData := addBlockCmd.String("block", "", "Block Data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.HandleErrors(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleErrors(err)
	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}

		cli.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}
