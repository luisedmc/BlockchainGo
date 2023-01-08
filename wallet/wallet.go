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

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
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

// GetAddress returns the wallet address.
func (w Wallet) GetAddress() []byte {
	pubHash := PublicHashKey(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checkSum := checksum(versionedHash)

	fullHash := append(versionedHash, checkSum...)
	address := Base58Encode(fullHash)

	fmt.Printf("Public Key: %x\n", w.PublicKey)
	fmt.Printf("Public Hash: %x\n", pubHash)
	fmt.Printf("Address: %x\n", address)

	return address
}

// PublicHashKey hashes public key using SHA-256 and RIPEMD-160.
func PublicHashKey(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD := hasher.Sum(nil)

	return publicRIPEMD
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

// checksum returns the sum of the first 4 bytes of the SHA-256 hash of the payload.
func checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checkSumLength]
}
