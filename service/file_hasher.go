package service

import (
	"goflet/util/hash"
	"log"
)

// startAsyncTask starts a new async task and sends the result to the channel
func startAsyncTask(algo func(string) (string, error), path string, channel chan string) {
	result, err := algo(path)
	if err != nil {
		log.Printf("Error hashing file: %s", err.Error())
		channel <- ""
		return
	}
	channel <- result
}

// HashFile returns the hash of the file
func HashFile(path string) FileHash {
	sha1 := make(chan string)
	sha256 := make(chan string)
	md5 := make(chan string)

	go startAsyncTask(hash.FileSha1, path, sha1)
	go startAsyncTask(hash.FileSha256, path, sha256)
	go startAsyncTask(hash.FileMd5, path, md5)

	return FileHash{
		HashSha1:   <-sha1,
		HashSha256: <-sha256,
		HashMd5:    <-md5,
	}
}
