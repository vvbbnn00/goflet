// Package hash provides a simple interface to hash strings and files using various algorithms.
package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"

	"golang.org/x/crypto/sha3"
)

// prepareHasher prepares the hasher for the algorithm
func prepareHasher(alg Algorithm) (hasher interface{}) {
	switch alg {
	case Md5:
		hasher = md5.New()
	case Sha1:
		hasher = sha1.New()
	case Sha256:
		hasher = sha256.New()
	case Sha3New256:
		hasher = sha3.New256()
	}
	return
}

// hashString hashes the string
func hashString(alg Algorithm, data string) string {
	hasher := prepareHasher(alg).(hash.Hash)
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

// hashFile hashes the file
func hashFile(alg Algorithm, path string) (string, error) {
	hasher := prepareHasher(alg).(hash.Hash)
	fs, err := getFs(path)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(hasher, fs); err != nil {
		return "", err
	}
	_ = fs.Close()
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// FileMd5 returns the md5 hash of the file
func FileMd5(path string) (string, error) {
	return hashFile(Md5, path)
}

// FileSha1 returns the sha1 hash of the file
func FileSha1(path string) (string, error) {
	return hashFile(Sha1, path)
}

// FileSha256 returns the sha256 hash of the file
func FileSha256(path string) (string, error) {
	return hashFile(Sha256, path)
}

// FileSha3New256 returns the sha3-256 hash of the file
func FileSha3New256(path string) (string, error) {
	return hashFile(Sha3New256, path)
}

// StringMd5 returns the md5 hash of the string
func StringMd5(data string) string {
	return hashString(Md5, data)
}

// StringSha1 returns the sha1 hash of the string
func StringSha1(data string) string {
	return hashString(Sha1, data)
}

// StringSha256 returns the sha256 hash of the string
func StringSha256(data string) string {
	return hashString(Sha256, data)
}

// StringSha3New256 returns the sha3-256 hash of the string
func StringSha3New256(data string) string {
	return hashString(Sha3New256, data)
}
