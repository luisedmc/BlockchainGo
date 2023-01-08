package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"log"
	"os"
)

const walletFileLocation = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

// SaveFile saves the content of the wallets to a file.
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)

	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFileLocation, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}

}

// LoadFile checks if the wallet file exists and loads its content into the wallets.
func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFileLocation); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := os.ReadFile(walletFileLocation)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets

	return nil
}
