package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/luisedmc/blockgo/blockchain"
)

type CommandLine struct{}

// PrintUsage prints command line usage.
func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage: ")
	fmt.Println("getbalance -address ADDRESS - get balance for address")
	fmt.Println("createblockchain -address ADDRESS - creates a blockchain")
	fmt.Println("printchain - Print all blocks in the chain.")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - Send amount to")
}

// ValidateArgs validates the command line arguments.
func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

// PrintChain prints the blocks in the chain.
func (cli *CommandLine) PrintChain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
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

// CreateBlockchain creates a new blockchain DB and adds a genesis block.
func (cli *CommandLine) CreateBlockchain(address string) {
	chain := blockchain.InitBlockchain(address)
	defer chain.Database.Close()

	fmt.Println("Finished!")
}

// GetBalance returns the balance of a specific address.
func (cli *CommandLine) GetBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// Send sends an amount from one address to another.
func (cli *CommandLine) Send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)

	chain.AddBlock([]*blockchain.Transaction{tx})

	fmt.Println("Success!")
}

// Run parses command line arguments and processes commands.
func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	// Create subcommands.
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	// Create flags for subcommands.
	getBalanceAddress := getBalanceCmd.String("address", "", "wallet address")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "miner address")
	sendFrom := sendCmd.String("from", "", "wallet address")
	sendTo := sendCmd.String("to", "", "wallet address")
	sendAmount := sendCmd.Int("amount", 0, "amount to send")

	// Switch on the first argument to determine the command.
	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}

		cli.GetBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}

		cli.CreateBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.Send(*sendFrom, *sendTo, *sendAmount)
	}
}
