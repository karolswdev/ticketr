package jobs

import (
	"context"
	"fmt"
	"sync"
)

// JobQueue manages the execution of async jobs with a worker pool.
// It provides progress aggregation, cancellation, and status tracking.
type JobQueue struct {
	// Configuration
	workers int

	// Channels
	jobChan      chan Job         // Buffered channel for incoming jobs
	progressChan chan JobProgress // Buffered channel for all progress events
	doneChan     chan struct{}    // Signals shutdown

	// State (protected by mu)
	mu       sync.Mutex
	contexts map[JobID]context.CancelFunc // Active job contexts for cancellation
	statuses map[JobID]JobStatus          // Job status tracking

	// Worker coordination
	wg sync.WaitGroup // Tracks active workers
}

// NewJobQueue creates a new job queue with the specified number of workers.
// workerCount determines how many jobs can run concurrently (recommended: 1-3).
func NewJobQueue(workerCount int) *JobQueue {
	if workerCount < 1 {
		workerCount = 1
	}

	jq := &JobQueue{
		workers:      workerCount,
		jobChan:      make(chan Job, 10),          // Buffer for job submissions
		progressChan: make(chan JobProgress, 100), // Large buffer to prevent blocking
		doneChan:     make(chan struct{}),
		contexts:     make(map[JobID]context.CancelFunc),
		statuses:     make(map[JobID]JobStatus),
	}

	// Start worker pool
	jq.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go jq.worker()
	}

	return jq
}

// Submit adds a job to the queue and returns its ID.
// The job will be executed by an available worker.
func (jq *JobQueue) Submit(job Job) JobID {
	jobID := job.ID()

	// Mark job as pending
	jq.mu.Lock()
	jq.statuses[jobID] = JobPending
	jq.mu.Unlock()

	// Send to worker pool (non-blocking due to buffer)
	jq.jobChan <- job

	return jobID
}

// Cancel attempts to cancel a job by ID.
// Returns an error if the job doesn't exist or is not cancellable.
func (jq *JobQueue) Cancel(jobID JobID) error {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	// Check if job exists
	status, exists := jq.statuses[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	// Check if job is cancellable
	if status != JobPending && status != JobRunning {
		return fmt.Errorf("job %s is not cancellable (status: %s)", jobID, status)
	}

	// Get cancel function
	cancelFunc, hasContext := jq.contexts[jobID]
	if !hasContext {
		// Job is pending but not started yet
		// Mark it as cancelled (worker will skip it)
		jq.statuses[jobID] = JobCancelled
		return nil
	}

	// Cancel the context
	cancelFunc()

	return nil
}

// Status returns the current status of a job.
// Returns JobStatus and true if the job exists, or empty status and false if not found.
func (jq *JobQueue) Status(jobID JobID) (JobStatus, bool) {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	status, exists := jq.statuses[jobID]
	return status, exists
}

// Progress returns a read-only channel that receives progress updates from all jobs.
// The channel remains open for the lifetime of the JobQueue.
func (jq *JobQueue) Progress() <-chan JobProgress {
	return jq.progressChan
}

// Shutdown gracefully shuts down the job queue.
// It waits for all active jobs to complete and then closes all channels.
// Returns an error if shutdown fails.
func (jq *JobQueue) Shutdown() error {
	// Signal shutdown
	close(jq.doneChan)

	// Close job channel (no more jobs accepted)
	close(jq.jobChan)

	// Wait for all workers to finish
	jq.wg.Wait()

	// Close progress channel
	close(jq.progressChan)

	return nil
}

// worker is the main worker goroutine that processes jobs from the queue.
func (jq *JobQueue) worker() {
	defer jq.wg.Done()

	for job := range jq.jobChan {
		jq.executeJob(job)
	}
}

// executeJob runs a single job with context and progress tracking.
func (jq *JobQueue) executeJob(job Job) {
	jobID := job.ID()

	// Check if job was cancelled before starting
	jq.mu.Lock()
	if jq.statuses[jobID] == JobCancelled {
		jq.mu.Unlock()
		return
	}
	jq.mu.Unlock()

	// Create cancellable context for this job
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cleanup

	// Register context and mark as running
	jq.mu.Lock()
	jq.contexts[jobID] = cancel
	jq.statuses[jobID] = JobRunning
	jq.mu.Unlock()

	// Start progress forwarding in separate goroutine
	progressDone := make(chan struct{})
	go jq.forwardProgress(job, progressDone)

	// Execute the job
	err := job.Execute(ctx)

	// Wait for progress forwarding to complete
	<-progressDone

	// Update status based on result
	jq.mu.Lock()
	if ctx.Err() == context.Canceled {
		// Job was cancelled
		jq.statuses[jobID] = JobCancelled
	} else if err != nil {
		// Job failed
		jq.statuses[jobID] = JobFailed
	} else {
		// Job completed successfully
		jq.statuses[jobID] = JobCompleted
	}
	// Remove context (no longer needed)
	delete(jq.contexts, jobID)
	jq.mu.Unlock()
}

// forwardProgress forwards progress events from a job's channel to the global progress channel.
// It uses a non-blocking send to prevent slow UI updates from blocking job execution.
func (jq *JobQueue) forwardProgress(job Job, done chan struct{}) {
	defer close(done)

	for progress := range job.Progress() {
		select {
		case jq.progressChan <- progress:
			// Progress sent successfully
		default:
			// Progress channel full - drop this update
			// This prevents blocking if UI is slow to consume updates
		}
	}
}
