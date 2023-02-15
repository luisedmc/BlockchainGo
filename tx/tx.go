package tx

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

// Transaction holds the transaction inputs and outputs
type Transaction struct {
	ID      []byte
	Inputs  []TXInput
	Outputs []TXOutput
}

// SetTransactionID sets the transaction ID
func (tx *Transaction) SetTransactionID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	if err := encode.Encode(tx); err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())

	tx.ID = hash[:]
}

// NewCoinBaseTX creates a new coinbase transaction
func NewCoinBaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to: '%s'.", to)
	}

	// Hard coding the transaction
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
	tx.SetTransactionID()

	return &tx
}

// IsCoinBase checks if a transaction is coinbase
func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Output == -1
}
