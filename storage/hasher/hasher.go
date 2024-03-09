package hasher

import (
	"goflet/storage"
	"goflet/storage/model"
	"goflet/util/hash"
	"goflet/worker"
	"log"
	"path/filepath"
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
			args := job.Args.([1]string)
			return updateFileHash(args[0])
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

// hashFile returns the hash of the file
func hashFile(path string) model.FileHash {
	sha1 := make(chan string)
	sha256 := make(chan string)
	md5 := make(chan string)

	// Append the file-append to the path
	path = filepath.Join(path, model.FileAppend)

	go startAsyncTask(hash.FileSha1, path, sha1)
	go startAsyncTask(hash.FileSha256, path, sha256)
	go startAsyncTask(hash.FileMd5, path, md5)

	return model.FileHash{
		HashSha1:   <-sha1,
		HashSha256: <-sha256,
		HashMd5:    <-md5,
	}
}

// updateFileHash updates the hash of the file
func updateFileHash(fsPath string) error {
	fileHash := hashFile(fsPath)
	err := storage.UpdateFileMeta(fsPath, model.FileMeta{
		Hash: fileHash,
	})
	if err != nil {
		log.Printf("Error updating file meta: %s", err.Error())
		return err
	}
	return nil
}

// HashFileAsync updates the hash of the file asynchronously
func HashFileAsync(fsPath string) {
	hashTaskPool.JobChain <- worker.Job{
		RetryCount: 0,
		Args:       [1]string{fsPath},
	}
}
