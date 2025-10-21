package jobs

import (
	"context"
	"fmt"
)

// JobID is a unique identifier for a job.
type JobID string

// JobStatus represents the current state of a job.
type JobStatus string

const (
	// JobPending indicates the job is queued but not yet started.
	JobPending JobStatus = "pending"
	// JobRunning indicates the job is currently executing.
	JobRunning JobStatus = "running"
	// JobCompleted indicates the job finished successfully.
	JobCompleted JobStatus = "completed"
	// JobFailed indicates the job finished with an error.
	JobFailed JobStatus = "failed"
	// JobCancelled indicates the job was cancelled by the user.
	JobCancelled JobStatus = "cancelled"
)

// String returns the string representation of a JobStatus.
func (s JobStatus) String() string {
	return string(s)
}

// JobProgress represents a progress update from a running job.
type JobProgress struct {
	// JobID identifies which job this progress update belongs to.
	JobID JobID
	// Current is the number of items processed so far.
	Current int
	// Total is the total number of items to process (0 if unknown).
	Total int
	// Percentage is the calculated completion percentage (0-100).
	Percentage float64
	// Message is a human-readable status message.
	Message string
}

// Job represents an asynchronous operation that can report progress and be cancelled.
type Job interface {
	// Execute runs the job with the given context.
	// The job must respect context cancellation and return quickly when ctx.Done() is closed.
	// Returns an error if the job fails, or nil on success.
	Execute(ctx context.Context) error

	// ID returns the unique identifier for this job.
	ID() JobID

	// Progress returns a read-only channel that emits progress updates.
	// The channel is closed when the job completes (successfully or with error).
	Progress() <-chan JobProgress
}

// calculatePercentage computes the completion percentage from current and total counts.
// Returns 0 if total is 0 or unknown.
func calculatePercentage(current, total int) float64 {
	if total <= 0 {
		return 0
	}
	percentage := (float64(current) / float64(total)) * 100
	if percentage > 100 {
		return 100
	}
	return percentage
}

// FormatProgress creates a human-readable progress string.
func FormatProgress(progress JobProgress) string {
	if progress.Total > 0 {
		return fmt.Sprintf("%s: %d/%d (%.0f%%)",
			progress.Message,
			progress.Current,
			progress.Total,
			progress.Percentage)
	}
	if progress.Message != "" {
		return progress.Message
	}
	return "Processing..."
}
