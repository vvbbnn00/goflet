package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// FileSha1 returns the sha1 hash of the file
func FileSha1(path string) (string, error) {
	reader, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	hasher := sha1.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return "", err
	}
	_ = reader.Close()
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
