package wallet

import "github.com/btcsuite/btcutil/base58"

// Base58Encode encodes to Base58 encoding
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

// Base58Decode decodes a Base58 encoded data
func Base58Decode(input []byte) []byte {
	decode := base58.Decode(string(input))

	return decode
}
