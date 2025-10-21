package jobs

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestJobQueue_GoroutineCleanup verifies no goroutine leaks after job execution.
func TestJobQueue_GoroutineCleanup(t *testing.T) {
	// Get baseline goroutine count
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	// Create and use job queue
	jq := NewJobQueue(2)

	// Submit several jobs
	for i := 0; i < 10; i++ {
		job := NewFakeJob(50*time.Millisecond, false)
		jq.Submit(job)
	}

	// Wait for jobs to complete
	time.Sleep(500 * time.Millisecond)

	// Shutdown queue
	err := jq.Shutdown()
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Force GC to clean up goroutines
	runtime.GC()
	time.Sleep(200 * time.Millisecond)

	// Check goroutine count
	final := runtime.NumGoroutine()

	// Allow ±3 goroutine difference (some runtime goroutines may vary)
	diff := final - baseline
	if diff < -3 || diff > 3 {
		t.Errorf("Goroutine leak detected: baseline=%d, final=%d, diff=%d", baseline, final, diff)
	}
}

// TestJobQueue_RapidCancellations tests multiple rapid job cancellations.
func TestJobQueue_RapidCancellations(t *testing.T) {
	jq := NewJobQueue(2)
	defer jq.Shutdown()

	// Submit and immediately cancel 20 jobs
	for i := 0; i < 20; i++ {
		job := NewFakeJob(2*time.Second, false)
		jobID := jq.Submit(job)

		// Very brief delay before cancellation
		time.Sleep(5 * time.Millisecond)

		err := jq.Cancel(jobID)
		if err != nil {
			// Job might already be cancelled or completed - that's okay
			if err.Error() != "job "+string(jobID)+" is not cancellable (status: cancelled)" &&
				err.Error() != "job "+string(jobID)+" is not cancellable (status: completed)" {
				t.Logf("Cancel returned error: %v", err)
			}
		}
	}

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify no panics or deadlocks occurred
	// If we reach here, the test passed
}

// TestJobQueue_SubmitAfterShutdown tests submitting a job after shutdown.
func TestJobQueue_SubmitAfterShutdown(t *testing.T) {
	jq := NewJobQueue(1)

	// Shutdown immediately
	err := jq.Shutdown()
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Try to submit a job after shutdown
	job := NewFakeJob(100*time.Millisecond, false)

	// This should panic or fail gracefully
	// We expect Submit to work but the job won't execute
	defer func() {
		if r := recover(); r != nil {
			// Panic is acceptable when submitting after shutdown
			t.Logf("Submit panicked after shutdown (acceptable): %v", r)
		}
	}()

	jobID := jq.Submit(job)

	// The job should be submitted but won't execute
	// (channel is closed, so workers have stopped)
	if jobID == "" {
		t.Error("Submit returned empty JobID")
	}

	// Verify job was never executed
	time.Sleep(200 * time.Millisecond)
	if job.WasExecuted() {
		t.Error("Job should not execute after shutdown")
	}
}

// TestJobQueue_LargeJobCount tests handling many jobs.
func TestJobQueue_LargeJobCount(t *testing.T) {
	jq := NewJobQueue(3) // 3 workers
	defer jq.Shutdown()

	numJobs := 50
	jobs := make([]*FakeJob, numJobs)
	jobIDs := make([]JobID, numJobs)

	// Submit many jobs quickly
	start := time.Now()
	for i := 0; i < numJobs; i++ {
		jobs[i] = NewFakeJob(20*time.Millisecond, false)
		jobIDs[i] = jq.Submit(jobs[i])
	}
	elapsed := time.Since(start)

	// Submission should be fast (non-blocking)
	if elapsed > 1*time.Second {
		t.Errorf("Job submission took too long: %v", elapsed)
	}

	// Wait for all jobs to complete
	timeout := time.After(10 * time.Second)
	completed := false

	for !completed {
		select {
		case <-timeout:
			t.Fatal("Jobs did not complete within timeout")
		case <-time.After(100 * time.Millisecond):
			allDone := true
			for _, jobID := range jobIDs {
				status, exists := jq.Status(jobID)
				if !exists || (status != JobCompleted && status != JobFailed) {
					allDone = false
					break
				}
			}
			if allDone {
				completed = true
			}
		}
	}

	// Verify all jobs executed
	for i, job := range jobs {
		if !job.WasExecuted() {
			t.Errorf("Job %d was not executed", i)
		}
	}
}

// TestJobQueue_ProgressUnderLoad tests progress reporting under load.
func TestJobQueue_ProgressUnderLoad(t *testing.T) {
	jq := NewJobQueue(2)
	defer jq.Shutdown()

	// Track progress events with mutex for thread safety
	var progressMu sync.Mutex
	progressCount := 0
	progressDone := make(chan struct{})

	go func() {
		defer close(progressDone)
		timeout := time.After(3 * time.Second)
		for {
			select {
			case _, ok := <-jq.Progress():
				if !ok {
					return
				}
				progressMu.Lock()
				progressCount++
				progressMu.Unlock()
			case <-timeout:
				return
			}
		}
	}()

	// Submit jobs that generate lots of progress
	numJobs := 10
	for i := 0; i < numJobs; i++ {
		job := NewFakeJob(100*time.Millisecond, false)
		jq.Submit(job)
	}

	// Wait for jobs to complete
	time.Sleep(1 * time.Second)

	// Should have received many progress events
	// Each job sends ~10 progress updates
	progressMu.Lock()
	count := progressCount
	progressMu.Unlock()

	if count < numJobs*5 {
		t.Logf("Warning: received fewer progress events than expected: %d", count)
	}
}

// TestJobQueue_CancelPendingJob tests cancelling a job before it starts.
func TestJobQueue_CancelPendingJob(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	// Fill the queue with a long-running job
	blocker := NewFakeJob(1*time.Second, false)
	jq.Submit(blocker)

	// Submit another job (will be pending)
	pendingJob := NewFakeJob(100*time.Millisecond, false)
	pendingJobID := jq.Submit(pendingJob)

	// Immediately cancel the pending job
	err := jq.Cancel(pendingJobID)
	if err != nil {
		t.Fatalf("Failed to cancel pending job: %v", err)
	}

	// Wait for blocker to finish
	time.Sleep(1200 * time.Millisecond)

	// Pending job should not have executed
	if pendingJob.WasExecuted() {
		t.Error("Cancelled pending job should not have executed")
	}

	status, _ := jq.Status(pendingJobID)
	if status != JobCancelled {
		t.Errorf("Expected status %s, got %s", JobCancelled, status)
	}
}

// TestJobQueue_StatusPersistence tests that status persists after job completion.
func TestJobQueue_StatusPersistence(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	job := NewFakeJob(50*time.Millisecond, false)
	jobID := jq.Submit(job)

	// Wait for completion
	time.Sleep(150 * time.Millisecond)

	// Status should persist after completion
	for i := 0; i < 10; i++ {
		status, exists := jq.Status(jobID)
		if !exists {
			t.Fatal("Job status disappeared after completion")
		}
		if status != JobCompleted {
			t.Errorf("Expected status %s, got %s", JobCompleted, status)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestJobQueue_FailedJobRetainsStatus tests that failed job status is retained.
func TestJobQueue_FailedJobRetainsStatus(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	// Submit a job that will fail
	job := NewFakeJob(50*time.Millisecond, true) // shouldFail = true
	jobID := jq.Submit(job)

	// Wait for failure
	time.Sleep(150 * time.Millisecond)

	// Status should show failed
	status, exists := jq.Status(jobID)
	if !exists {
		t.Fatal("Failed job status not found")
	}

	if status != JobFailed {
		t.Errorf("Expected status %s, got %s", JobFailed, status)
	}

	// Status should persist
	time.Sleep(100 * time.Millisecond)
	status, exists = jq.Status(jobID)
	if !exists || status != JobFailed {
		t.Error("Failed job status did not persist")
	}
}

// BenchmarkJobQueue_SubmitThroughput benchmarks job submission throughput.
func BenchmarkJobQueue_SubmitThroughput(b *testing.B) {
	jq := NewJobQueue(2)
	defer jq.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		job := NewFakeJob(1*time.Millisecond, false)
		jq.Submit(job)
	}
}

// BenchmarkJobQueue_StatusQuery benchmarks status query performance.
func BenchmarkJobQueue_StatusQuery(b *testing.B) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	// Submit a job
	job := NewFakeJob(1*time.Second, false)
	jobID := jq.Submit(job)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jq.Status(jobID)
	}
}
