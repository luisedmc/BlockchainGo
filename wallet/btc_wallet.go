package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checkSumLength = 4
	version        = byte(0x00)
)

// Wallet holds the Private Key and Public Key, to sign transactions and verify the signature, respectively
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// CreateWallet creates a new Wallet by generating a new Private and Public key
func CreateWallet() *Wallet {
	private, public := generateKeyPair()

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

// Address returns the Wallet address
func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	address := Base58Encode(pubHash)

	fmt.Printf("pub hash: %x\n", pubHash)
	fmt.Printf("pub key: %x\n", w.PublicKey)
	fmt.Printf("address: %x\n", address)

	return address
}

// generateKeyPair generates a new Public and Private key
func generateKeyPair() (ecdsa.PrivateKey, []byte) {
	// secp256k1 curve also knows as P-256
	curve := elliptic.P256()

	// Generating Private Key with ECDSA algorithm
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	// Private Key derives Public Key
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, publicKey
}
