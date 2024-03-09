package worker

import (
	"errors"
	"testing"
	"time"
)

// mockDoFunction simulates a job function that can fail and succeed based on input.
func mockDoFunction(job Job) error {
	if job.Args.(int)%2 == 0 {
		return errors.New("mock error")
	}
	return nil
}

// TestWorkerPoolStop ensures that the pool stops all workers gracefully.
func TestWorkerPoolStop(t *testing.T) {
	workerFactory := func() Worker {
		return Worker{
			JobName: "TestJob",
			Do:      mockDoFunction,
		}
	}

	pool := NewPool(100, 1000, workerFactory)
	pool.Start()

	// Add some jobs to the pool.
	for i := 0; i < 1000; i++ {
		pool.AddJob(Job{Args: i})
	}

	t.Log("All jobs added to the pool")

	time.Sleep(3 * time.Second)

	// Stop the pool.
	pool.Stop()

	// Check if the job channel is closed.
	job, ok := <-pool.JobChain
	if ok {
		print(job.RetryCount, job.Args)
		t.Error("Job channel is not closed")
	}

	// Check if the cancel channel is closed.

	_, ok = <-pool.CancelChain
	if ok {
		t.Error("Cancel channel is not closed")
	}
}
