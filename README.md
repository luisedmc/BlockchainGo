## BlockchainGo

BlockchainGo is a simple implementation of the Blockchain technology in Go language. This project was build to learn more about the Blockchain technology and its features. The code was written with some comments in a way that can maybe help understand the most complicated parts. The code is also very direct, I didn't try to reinvent the wheel, just implement some of the Blockchain features.

I will be adding more features to this project in the future. See [Roadmap](#roadmap) section for more details.

#### Disclaimer

Like i said before, this project was built for learning purposes. It is not a production ready code. It is not secure or optimized and probably was not written in the best way.

## Table of Contents

- [BlockchainGo](#blockchaingo)
  - [Disclaimer](#disclaimer)
- [Table of Contents](#table-of-contents)
- [Blockchain](#blockchain)
  - [What's Blockchain?](#whats-blockchain)
- [Features](#features)
- [Usage](#usage)
  - [Commands](#commands)
    - [Example](#example)
- [Roadmap](#roadmap)

## Blockchain

### What's Blockchain?

<img src="/images/rick_and_morty.jpg" alt="Rick and Morty explaning Blockchain" width=380 height=380>

## Features

At the current state, the Blockchain has the following features:

- Create a new Blockchain and store it in the database
- Print the Blockchain in the terminal
- Transactions
- Get the balance of an address

## Usage

Clone this repository and change to the project directory:

```
git clone https://github.com/luisedmc/BlockchainGo && cd BlockchainGo
```

In the project directory, just `go run .` and then the commands listed below.

You can also build the binary file with `go build` and then run it using `./BlockchainGo`

### Commands

| Command                                 | Description                                                       |
| --------------------------------------- | ----------------------------------------------------------------- |
| `createBlockchain -address ADDRESS`     | Create a new Blockchain with the given address as the first miner |
| `printChain`                            | Print all the blocks of the blockchain                            |
| `getBalance -address ADDRESS`           | Get the balance of the given address                              |
| `send -from FROM -to TO -amount AMOUNT` | Send the given amount of coins from one address to another        |

#### Example

`go run . createBlockchain -address "luisedmc"`
or
`./BlockchainGo createBlockchain -address "luisedmc"`

## Roadmap

Features that I would like to implement in the future:

- [ ] Wallets
- [ ] Proof of Stake
- [ ] Ethereum
