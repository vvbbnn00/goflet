package hash

import (
	"crypto/md5"
	"encoding/hex"
	"goflet/storage/model"
	"io"
	"os"
)

// FileMd5 returns the md5 hash of the file
func FileMd5(path string) (string, error) {
	reader, err := os.OpenFile(path, os.O_RDONLY, model.FilePerm)
	if err != nil {
		return "", err
	}
	hasher := md5.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return "", err
	}
	_ = reader.Close()
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
