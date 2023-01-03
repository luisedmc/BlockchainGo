package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

// Transaction represents a transaction.
type Transaction struct {
	ID      []byte
	Inputs  []TXInput
	Outputs []TXOutput
}

type TXInput struct {
	ID        []byte
	Output    int
	Signature string
}

type TXOutput struct {
	Value  int
	PubKey string
}

// SetID sets ID of a transaction.
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())

	tx.ID = hash[:]
}

// IsCoinBase checks if transaction is coin base.
func CoinBaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TXInput{
		ID:        []byte{},
		Output:    -1,
		Signature: data,
	}

	txout := TXOutput{
		Value:  100,
		PubKey: to,
	}

	tx := Transaction{
		ID:      nil,
		Inputs:  []TXInput{txin},
		Outputs: []TXOutput{txout},
	}

	tx.SetID()

	return &tx
}
