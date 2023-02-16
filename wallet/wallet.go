package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

// Wallet holds the Public Key and Private Key, to encrypt and decrypt transaction data, respectively
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewKeyPair returns a new Public Key and Private Key
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	// Generating Private Key with ecdsa algorithm
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	// Public Key derives Private Key
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, publicKey
}

// CreateWallet creates a new Wallet
func CreateWallet() *Wallet {
	private, public := NewKeyPair()

	wallet := Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}

	return &wallet
}
