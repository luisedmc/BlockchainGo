package wallet

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

// Base58 alphabet don't contains the characters 0 (zero), O (capital letter o), I (capital letter i) and l (lowercase letter L) to avoid mistakes and errors.

// Base58Encode encodes to Base58 string
func Base58Encode(input []byte) []byte {
	encoded := base58.CheckEncode(input, version)

	return []byte(encoded)
}

// Base58Decode decodes a Base58 string
func Base58Decode(input []byte) []byte {
	decoded, version, err := base58.CheckDecode(string(input))

	fmt.Println("Version Byte: ", version)

	if err != nil {
		panic(err)
	}

	return decoded
}
