package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/luisedmc/BlockchainGo/blockchain"
	"github.com/luisedmc/BlockchainGo/tx"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Command Line Usage:")
	fmt.Println()
	fmt.Println("createBlockchain -address ADDRESS | Creates a Blockchain.")
	fmt.Println("printChain | Prints all Blocks in the Blockchain.")
	fmt.Println("getBalance -address ADDRESS | Gets address balance.")
	fmt.Println("send -from FROM -to TO -amount AMOUNT | Send a amount from one address to another.")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.CreateBlockchain(address)
	chain.Database.Close()

	fmt.Println("Blockchain created!")
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()

	iter := chain.Iterator()

	// Iterating through the Blocks on the Blockchain
	for {
		block := iter.Next()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Block Hash: %x\n", block.Hash)

		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("Proof of Work: %s\n", strings.Title(strconv.FormatBool(pow.Validate())))
		fmt.Println()

		// Break on Genesis Block
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Address: %s\nBalance: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()

	transaction := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*tx.Transaction{transaction})

	fmt.Printf("The amount of %d has been sent successfully.\nFrom: %s\tTo: %s", amount, from, to)
}

func (cli *CommandLine) RunCLI() {
	cli.validateArgs()

	// Commands
	createBlockchainCmd := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "Wallet address")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "Miner address")
	sendFrom := sendCmd.String("from", "", "Sender address")
	sendTo := sendCmd.String("to", "", "Receiver address")
	sendAmount := sendCmd.Int("amount", 0, "Sending amount")

	// Checking commands
	switch os.Args[1] {
	case "createBlockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}

		cli.createBlockchain(*createBlockchainAddress)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}

		cli.getBalance(*getBalanceAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == 0 {
			sendCmd.Usage()
			runtime.Goexit()

		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}