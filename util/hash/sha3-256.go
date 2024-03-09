package hash

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

// StringSha3New256 returns the sha3-256 hash of the string
func StringSha3New256(data string) string {
	hasher := sha3.New256()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}
