package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// newKeyPair creates and returns ECDSA private and public keys.
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// NewWallet creates and returns a new Wallet.
func NewWallet() *Wallet {
	private, public := newKeyPair()

	wallet := Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}

	return &wallet
}
