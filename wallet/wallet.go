package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checkSumLength = 4
	version        = byte(0x00)
)

// Wallet holds the Public Key and Private Key, to encrypt and decrypt transaction data, respectively
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// CreateWallet creates a new Wallet
func CreateWallet() *Wallet {
	private, public := newKeyPair()

	wallet := Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}

	return &wallet
}

// PublicKeyHash hashes twice the Public Key using RIPEMD160 and SHA-256 algorithms
func PublicKeyHash(publicKey []byte) []byte {
	publicSHA256 := sha256.Sum256(publicKey)

	hasherRIPEMD160 := ripemd160.New()

	_, err := hasherRIPEMD160.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := hasherRIPEMD160.Sum(nil)

	return publicRIPEMD160
}

// NewKeyPair returns a new Public Key and Private Key
func newKeyPair() (ecdsa.PrivateKey, []byte) {
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

// checkSum is used to verify an address
func checkSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checkSumLength]
}
