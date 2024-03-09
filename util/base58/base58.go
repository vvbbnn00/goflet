// Package base58 provides base58 encoding and decoding functions.
package base58

import (
	"bytes"
	"math/big"
)

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Encode encodes a byte slice to a base58 encoded string
func Encode(input []byte) string {
	result := make([]byte, 0, len(input)*136/100)
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(base58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, base58Alphabet[mod.Int64()])
	}

	// Add leading zeros
	for _, b := range input {
		if b != 0 {
			break
		}
		result = append(result, base58Alphabet[0])
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// Decode decodes a base58 encoded string
func Decode(input string) ([]byte, error) {
	result := big.NewInt(0)
	for _, c := range input {
		result.Mul(result, big.NewInt(int64(len(base58Alphabet))))
		result.Add(result, big.NewInt(int64(bytes.IndexByte([]byte(base58Alphabet), byte(c)))))
	}

	decoded := result.Bytes()

	// Add leading zeros
	for i := 0; i < len(input) && input[i] == base58Alphabet[0]; i++ {
		decoded = append([]byte{0}, decoded...)
	}

	return decoded, nil
}
