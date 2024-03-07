package service

import (
	"goflet/util/hash"
	"goflet/worker"
	"log"
)

const (
	hashTaskMaxWorkers = 4     // The maximum number of workers for the hash task
	hashTaskBufferSize = 10000 // The buffer size for the hash task
)

var hashTaskPool *worker.Pool = worker.NewPool(hashTaskMaxWorkers, hashTaskBufferSize, workerFactory) // The pool of workers for the hash task

func init() {
	// Start the hash task pool
	hashTaskPool.Start()
}

// workerFactory creates a new worker
func workerFactory() worker.Worker {
	return worker.Worker{
		JobName: "HashTask",
		Do: func(job worker.Job) error {
			args := job.Args.([2]string)
			return updateFileHash(args[0], args[1])
		},
	}
}

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

// updateFileHash updates the hash of the file
func updateFileHash(path string, decodedPath string) error {
	fileHash := HashFile(path)
	err := UpdateFileMeta(decodedPath, FileMeta{
		Hash: fileHash,
	})
	if err != nil {
		log.Printf("Error updating file meta: %s", err.Error())
		return err
	}
	return nil
}

// HashFileAsync updates the hash of the file asynchronously
func HashFileAsync(path string, decodedPath string) {
	hashTaskPool.JobChain <- worker.Job{
		RetryCount: 0,
		Args:       [2]string{path, decodedPath},
	}
}
