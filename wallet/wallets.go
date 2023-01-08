package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

const walletFileLocation = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

// CreateWallets creates a Wallets object and fills it with the wallets from the file.
func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err
}

// GetWallet returns the wallet for the given address.
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// AddWallet creates a new wallet and adds it to the wallets.
func (ws *Wallets) AddWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	ws.Wallets[address] = wallet

	return address
}

// GetAllAdresses returns all the addresses in the wallets.
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
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
