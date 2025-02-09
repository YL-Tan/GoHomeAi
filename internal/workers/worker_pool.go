package workers

import (
	"fmt"
	"time"
)

const (
	NumWorkers   = 3
	JobQueueSize = 100	// Maximum pending jobs in the queue
)

type Job struct {
	ID      int
	Message string
}

type WorkerPool struct {
	jobQueue chan Job
	quit     chan struct{}
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		jobQueue: make(chan Job, JobQueueSize),
		quit:     make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < NumWorkers; i++ {
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	fmt.Printf("Worker %d started\n", id)
	for {
		select {
		case job := <-wp.jobQueue:
			fmt.Printf("Worker %d processing job: %d - %s\n", id, job.ID, job.Message)
			time.Sleep(2 * time.Second)
		case <-wp.quit:
			fmt.Printf("Worker %d shutting down\n", id)
			return
		}
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobQueue <- job
}

func (wp *WorkerPool) Stop() {
	close(wp.quit)
	close(wp.jobQueue)
}
