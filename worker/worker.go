package worker

import (
	"log"
	"sync"
	"time"
)

const (
	maxJobRetries   = 3 // Maximum number of retries for a job
	NoPoolSizeLimit = 0 // No limit for the job queue size
)

const retryDelay = 1 * time.Second // Delay before retrying a job

// Job represents a job to be executed by a worker
type Job struct {
	RetryCount int         // The number of retries
	Args       interface{} // The arguments for the job
}

// Worker represents a worker that executes jobs
type Worker struct {
	JobName string // The name of the job, to display in logs
	Do      func(Job) error
}

// Pool represents a pool of workers
type Pool struct {
	WorkerCount   int
	JobChain      chan Job
	CancelChain   chan struct{}
	WorkerFactory func() Worker // Factory function to create workers
	workers       []Worker
	wg            sync.WaitGroup
}

// work starts the worker and listens for jobs
func (w *Worker) work(jobChain chan Job, cancelChain <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-cancelChain:
			log.Printf("[Worker] Worker %s cancelled", w.JobName)
			return
		case job := <-jobChain:
			// Check if the job has exceeded the maximum number of retries
			if job.RetryCount > maxJobRetries {
				log.Printf("[Worker] Job %s(%v) failed after %d retries", w.JobName, job.Args, maxJobRetries)
				continue
			}

			err := w.Do(job)
			if err != nil {
				job.RetryCount++
				backoffDuration := time.Duration(job.RetryCount) * retryDelay // Exponential backoff

				log.Printf("[Worker] Failed to execute job %s(%v): %s, retry %d, wait %d sec",
					w.JobName,
					job.Args,
					err.Error(),
					job.RetryCount,
					backoffDuration/time.Second) // Log the error

				// Before sleep, check if the job has been cancelled
				select {
				case <-cancelChain:
					log.Printf("[Worker] Worker %s cancelled", w.JobName)
					return
				default:
				}
				time.Sleep(backoffDuration) // Wait before retrying

				log.Printf("[Worker] Retrying job %s, retry %d", w.JobName, job.RetryCount)

				// Retry the job
				go func() {
					// check whether the job chain is closed
					if jobChain == nil {
						return
					}
					select {
					case jobChain <- job:
					case <-cancelChain:
						log.Printf("[Worker] Worker %s cancelled", w.JobName)
					}
				}()
			}
		}
	}
}

// NewPool creates a new pool of workers
func NewPool(workerCount int, jobQueueSize int, workerFactory func() Worker) *Pool {
	if workerCount <= 0 {
		panic("Invalid worker count")
	}
	if jobQueueSize < 0 {
		panic("Invalid job queue size")
	}
	jobChain := make(chan Job, jobQueueSize)
	return &Pool{
		WorkerCount:   workerCount,
		JobChain:      jobChain,
		CancelChain:   make(chan struct{}, workerCount),
		WorkerFactory: workerFactory,
	}
}

// Start starts the pool of workers
func (p *Pool) Start() {
	for i := 0; i < p.WorkerCount; i++ {
		w := p.WorkerFactory()
		p.workers = append(p.workers, w)
		p.wg.Add(1)
		go w.work(p.JobChain, p.CancelChain, &p.wg)
	}
}

// Stop stops the pool of workers
func (p *Pool) Stop() {
	// Close the cancel channel to signal all workers to stop
	close(p.CancelChain)

	p.wg.Wait() // wait for all workers to finish

	// Drain the job chain
	for len(p.JobChain) > 0 {
		<-p.JobChain
	}

	close(p.JobChain) // close the job chain
}

// AddJob adds a new job to the pool
func (p *Pool) AddJob(job Job) {
	p.JobChain <- job
}
