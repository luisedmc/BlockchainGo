package wallet

import (
	"fmt"
)

const walletFile = "./tmp/wallets.data"

// Wallets holds multiple Wallets, keyed by the address
type Wallets struct {
	Wallets map[string]*Wallet
}

// CreateWallets creates Wallets and fill it's content from an existent file
func CreateWallets() *Wallets {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	return &wallets
}

// GetWallet gets a specific Wallet by address
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// GetAllAddresses gets all Wallets addresses
func (ws Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws *Wallets) AddWallet() string {
	wallet := CreateWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}
