package wallet

import (
	"github.com/btcsuite/btcutil/base58"
)

// Base58 alphabet don't contains the characters 0 (zero), O (capital letter o), I (capital letter i) and l (lowercase letter L) to avoid mistakes and errors.

// Base58Encode encodes to Base58 string
func Base58Encode(input []byte) []byte {
	encoded := base58.Encode(input)

	return []byte(encoded)
}

// Base58Decode decodes a Base58 string
func Base58Decode(input []byte) []byte {
	decoded := base58.Decode(string(input))

	return decoded

}
