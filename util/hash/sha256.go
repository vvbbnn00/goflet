package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// FileSha256 returns the sha256 hash of the file
func FileSha256(path string) (string, error) {
	reader, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	hasher := sha256.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return "", err
	}
	_ = reader.Close()
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// StringSha256 returns the sha256 hash of the string
func StringSha256(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}
