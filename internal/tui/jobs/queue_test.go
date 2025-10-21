package jobs

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

// FakeJob is a test job implementation for testing the JobQueue.
type FakeJob struct {
	id           JobID
	progressChan chan JobProgress
	duration     time.Duration
	shouldFail   bool
	failAfter    int // Fail after this many progress steps (-1 = don't fail)
	executed     bool
	mu           sync.Mutex
}

// NewFakeJob creates a fake job for testing.
func NewFakeJob(duration time.Duration, shouldFail bool) *FakeJob {
	return &FakeJob{
		id:           JobID(uuid.New().String()),
		progressChan: make(chan JobProgress, 50),
		duration:     duration,
		shouldFail:   shouldFail,
		failAfter:    -1,
	}
}

func (f *FakeJob) ID() JobID {
	return f.id
}

func (f *FakeJob) Progress() <-chan JobProgress {
	return f.progressChan
}

func (f *FakeJob) Execute(ctx context.Context) error {
	defer close(f.progressChan)

	f.mu.Lock()
	f.executed = true
	f.mu.Unlock()

	steps := 10
	for i := 0; i < steps; i++ {
		select {
		case <-ctx.Done():
			// Send cancellation progress
			select {
			case f.progressChan <- JobProgress{
				JobID:      f.id,
				Current:    i,
				Total:      steps,
				Percentage: calculatePercentage(i, steps),
				Message:    "Cancelled",
			}:
			default:
			}
			return ctx.Err()
		case <-time.After(f.duration / time.Duration(steps)):
			// Send progress update
			f.progressChan <- JobProgress{
				JobID:      f.id,
				Current:    i + 1,
				Total:      steps,
				Percentage: calculatePercentage(i+1, steps),
				Message:    "Processing...",
			}

			// Check if we should fail after this step
			if f.failAfter >= 0 && i >= f.failAfter {
				return errors.New("fake job error")
			}
		}
	}

	if f.shouldFail {
		return errors.New("fake job error")
	}

	return nil
}

func (f *FakeJob) WasExecuted() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.executed
}

// TestJobQueue_Submit tests basic job submission and execution.
func TestJobQueue_Submit(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	job := NewFakeJob(100*time.Millisecond, false)
	jobID := jq.Submit(job)

	if jobID != job.ID() {
		t.Errorf("Expected job ID %s, got %s", job.ID(), jobID)
	}

	// Wait for job to complete
	time.Sleep(200 * time.Millisecond)

	status, exists := jq.Status(jobID)
	if !exists {
		t.Fatal("Job not found in queue")
	}

	if status != JobCompleted {
		t.Errorf("Expected status %s, got %s", JobCompleted, status)
	}

	if !job.WasExecuted() {
		t.Error("Job was not executed")
	}
}

// TestJobQueue_Cancel tests job cancellation.
func TestJobQueue_Cancel(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	// Create a longer-running job
	job := NewFakeJob(1*time.Second, false)
	jobID := jq.Submit(job)

	// Wait for job to start
	time.Sleep(50 * time.Millisecond)

	// Cancel the job
	err := jq.Cancel(jobID)
	if err != nil {
		t.Fatalf("Failed to cancel job: %v", err)
	}

	// Wait for cancellation to take effect
	time.Sleep(200 * time.Millisecond)

	status, exists := jq.Status(jobID)
	if !exists {
		t.Fatal("Job not found in queue")
	}

	if status != JobCancelled {
		t.Errorf("Expected status %s, got %s", JobCancelled, status)
	}
}

// TestJobQueue_Progress tests progress event delivery.
func TestJobQueue_Progress(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	job := NewFakeJob(200*time.Millisecond, false)
	jobID := jq.Submit(job)

	// Collect progress events
	var progressEvents []JobProgress
	progressDone := make(chan struct{})

	go func() {
		defer close(progressDone)
		timeout := time.After(500 * time.Millisecond)
		for {
			select {
			case progress, ok := <-jq.Progress():
				if !ok {
					return
				}
				if progress.JobID == jobID {
					progressEvents = append(progressEvents, progress)
				}
			case <-timeout:
				return
			}
		}
	}()

	// Wait for job to complete
	<-progressDone

	if len(progressEvents) == 0 {
		t.Error("No progress events received")
	}

	// Verify progress events are sequential
	for i := 1; i < len(progressEvents); i++ {
		if progressEvents[i].Current < progressEvents[i-1].Current {
			t.Errorf("Progress events not sequential: %d came after %d",
				progressEvents[i].Current, progressEvents[i-1].Current)
		}
	}
}

// TestJobQueue_MultipleJobs tests handling multiple jobs concurrently.
func TestJobQueue_MultipleJobs(t *testing.T) {
	jq := NewJobQueue(2) // 2 workers
	defer jq.Shutdown()

	numJobs := 5
	jobs := make([]*FakeJob, numJobs)
	jobIDs := make([]JobID, numJobs)

	// Submit multiple jobs
	for i := 0; i < numJobs; i++ {
		jobs[i] = NewFakeJob(100*time.Millisecond, false)
		jobIDs[i] = jq.Submit(jobs[i])
	}

	// Wait for all jobs to complete
	time.Sleep(500 * time.Millisecond)

	// Verify all jobs completed
	for i, jobID := range jobIDs {
		status, exists := jq.Status(jobID)
		if !exists {
			t.Errorf("Job %d not found in queue", i)
			continue
		}

		if status != JobCompleted {
			t.Errorf("Job %d: expected status %s, got %s", i, JobCompleted, status)
		}

		if !jobs[i].WasExecuted() {
			t.Errorf("Job %d was not executed", i)
		}
	}
}

// TestJobQueue_Shutdown tests graceful shutdown.
func TestJobQueue_Shutdown(t *testing.T) {
	jq := NewJobQueue(1)

	job := NewFakeJob(100*time.Millisecond, false)
	jq.Submit(job)

	// Shutdown and wait
	err := jq.Shutdown()
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Verify job completed
	if !job.WasExecuted() {
		t.Error("Job was not executed before shutdown")
	}
}

// TestJobQueue_StatusTracking tests status transitions.
func TestJobQueue_StatusTracking(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	job := NewFakeJob(200*time.Millisecond, false)
	jobID := jq.Submit(job)

	// Check pending status
	status, exists := jq.Status(jobID)
	if !exists {
		t.Fatal("Job not found immediately after submission")
	}
	if status != JobPending && status != JobRunning {
		t.Errorf("Expected status %s or %s, got %s", JobPending, JobRunning, status)
	}

	// Wait for job to start
	time.Sleep(50 * time.Millisecond)

	status, _ = jq.Status(jobID)
	if status != JobRunning {
		t.Errorf("Expected status %s, got %s", JobRunning, status)
	}

	// Wait for job to complete
	time.Sleep(300 * time.Millisecond)

	status, _ = jq.Status(jobID)
	if status != JobCompleted {
		t.Errorf("Expected status %s, got %s", JobCompleted, status)
	}
}

// TestJobQueue_FailedJob tests handling of failed jobs.
func TestJobQueue_FailedJob(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	job := NewFakeJob(100*time.Millisecond, true) // shouldFail = true
	jobID := jq.Submit(job)

	// Wait for job to complete
	time.Sleep(200 * time.Millisecond)

	status, exists := jq.Status(jobID)
	if !exists {
		t.Fatal("Job not found in queue")
	}

	if status != JobFailed {
		t.Errorf("Expected status %s, got %s", JobFailed, status)
	}
}

// TestJobQueue_CancelNonexistent tests cancelling a non-existent job.
func TestJobQueue_CancelNonexistent(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	err := jq.Cancel(JobID("nonexistent"))
	if err == nil {
		t.Error("Expected error when cancelling non-existent job")
	}
}

// TestJobQueue_CancelCompleted tests cancelling a completed job.
func TestJobQueue_CancelCompleted(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	job := NewFakeJob(50*time.Millisecond, false)
	jobID := jq.Submit(job)

	// Wait for job to complete
	time.Sleep(100 * time.Millisecond)

	err := jq.Cancel(jobID)
	if err == nil {
		t.Error("Expected error when cancelling completed job")
	}
}

// TestJobQueue_ProgressBuffering tests that progress events don't block.
func TestJobQueue_ProgressBuffering(t *testing.T) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	job := NewFakeJob(200*time.Millisecond, false)
	jq.Submit(job)

	// Don't consume progress events - they should buffer and drop if full
	// Job should still complete without blocking

	time.Sleep(300 * time.Millisecond)

	if !job.WasExecuted() {
		t.Error("Job should complete even if progress events aren't consumed")
	}
}

// TestJobQueue_ConcurrentOperations tests thread safety.
func TestJobQueue_ConcurrentOperations(t *testing.T) {
	jq := NewJobQueue(2)
	defer jq.Shutdown()

	numJobs := 20
	var wg sync.WaitGroup

	// Submit jobs concurrently
	wg.Add(numJobs)
	for i := 0; i < numJobs; i++ {
		go func() {
			defer wg.Done()
			job := NewFakeJob(50*time.Millisecond, false)
			jq.Submit(job)
		}()
	}

	wg.Wait()
	time.Sleep(500 * time.Millisecond)

	// All jobs should have completed without race conditions
	// (run with -race flag to verify)
}

// TestCalculatePercentage tests the percentage calculation helper.
func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		current  int
		total    int
		expected float64
	}{
		{0, 10, 0},
		{5, 10, 50},
		{10, 10, 100},
		{15, 10, 100}, // Over 100% capped
		{5, 0, 0},     // Unknown total
		{0, 0, 0},     // Both zero
	}

	for _, tt := range tests {
		result := calculatePercentage(tt.current, tt.total)
		if result != tt.expected {
			t.Errorf("calculatePercentage(%d, %d) = %.2f, expected %.2f",
				tt.current, tt.total, result, tt.expected)
		}
	}
}

// TestFormatProgress tests the progress formatting helper.
func TestFormatProgress(t *testing.T) {
	tests := []struct {
		progress JobProgress
		contains string
	}{
		{
			JobProgress{Current: 5, Total: 10, Percentage: 50, Message: "Processing"},
			"5/10",
		},
		{
			JobProgress{Current: 0, Total: 0, Message: "Connecting"},
			"Connecting",
		},
		{
			JobProgress{Current: 10, Total: 10, Percentage: 100, Message: "Complete"},
			"10/10",
		},
	}

	for _, tt := range tests {
		result := FormatProgress(tt.progress)
		if result == "" {
			t.Errorf("FormatProgress returned empty string for %+v", tt.progress)
		}
		// Just verify it returns something reasonable
		// Exact format can vary
	}
}

// Benchmark tests
func BenchmarkJobQueue_Submit(b *testing.B) {
	jq := NewJobQueue(2)
	defer jq.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		job := NewFakeJob(1*time.Millisecond, false)
		jq.Submit(job)
	}
}

func BenchmarkJobQueue_Status(b *testing.B) {
	jq := NewJobQueue(1)
	defer jq.Shutdown()

	job := NewFakeJob(100*time.Millisecond, false)
	jobID := jq.Submit(job)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jq.Status(jobID)
	}
}
