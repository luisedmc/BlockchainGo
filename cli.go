package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// CommandLine struct holds the Blockchain receiving CLI commands
type CommandLine struct {
	Blockchain *Blockchain
}

// PrintUsage prints all usable commands
func (cli *CommandLine) printUsage() {
	fmt.Println("\tCommand Line Usage: ")

	fmt.Println("\tadd -block BLOCK_DATA | Adds a new Block to the Blockchain.")
	fmt.Println("\tprint | Prints all Blocks in the Blockchain.")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.Blockchain.AddBlock(data)
	fmt.Println("New Block added in the Blockchain!")
}

// printChain prints all Blocks stored in the Blockchain Database
func (cli *CommandLine) printChain() {
	bci := cli.Blockchain.Iterator()

	// Iterating through the Blocks on the Blockchain
	for {
		block := bci.Next()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("Proof of Work: %s\n", strings.Title(strconv.FormatBool(pow.Validate())))
		fmt.Println()

		// Break on Genesis Block
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// RunCLI runs the Command Line Interface
func (cli *CommandLine) RunCLI() {
	cli.validateArgs()

	// Commands
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	// Subcommands
	addBlockData := addBlockCmd.String("block", "", "Block Data")

	// Checking commands
	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	// Checking subcommands
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
