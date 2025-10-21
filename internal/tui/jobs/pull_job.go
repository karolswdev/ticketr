package jobs

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/karolswdev/ticktr/internal/core/services"
)

// PullJob wraps a PullService operation to implement the Job interface.
// It provides async execution with progress reporting and cancellation support.
type PullJob struct {
	id           JobID
	pullService  *services.PullService
	filePath     string
	options      services.PullOptions
	progressChan chan JobProgress

	// Result storage (protected by mu)
	mu     sync.Mutex
	result *services.PullResult
	err    error
}

// NewPullJob creates a new pull job that will execute PullService.Pull().
func NewPullJob(pullService *services.PullService, filePath string, options services.PullOptions) *PullJob {
	return &PullJob{
		id:           JobID(uuid.New().String()),
		pullService:  pullService,
		filePath:     filePath,
		options:      options,
		progressChan: make(chan JobProgress, 50), // Buffer for smooth updates
	}
}

// ID returns the unique identifier for this job.
func (pj *PullJob) ID() JobID {
	return pj.id
}

// Progress returns the read-only progress channel.
func (pj *PullJob) Progress() <-chan JobProgress {
	return pj.progressChan
}

// Execute runs the pull operation with context cancellation support.
// It wraps PullService.Pull() and converts progress callbacks to JobProgress events.
func (pj *PullJob) Execute(ctx context.Context) error {
	defer close(pj.progressChan)

	// Track if we've been cancelled
	cancelled := false
	cancelMu := sync.Mutex{}

	// Set up progress callback to emit JobProgress events
	pj.options.ProgressCallback = func(current, total int, message string) {
		// Check if cancelled
		cancelMu.Lock()
		if cancelled {
			cancelMu.Unlock()
			return
		}
		cancelMu.Unlock()

		// Create progress event
		progress := JobProgress{
			JobID:      pj.id,
			Current:    current,
			Total:      total,
			Percentage: calculatePercentage(current, total),
			Message:    message,
		}

		// Send progress update (non-blocking if cancelled)
		select {
		case pj.progressChan <- progress:
		case <-ctx.Done():
			cancelMu.Lock()
			cancelled = true
			cancelMu.Unlock()
			return
		}
	}

	// Run the pull operation in a goroutine so we can cancel it
	resultChan := make(chan struct {
		result *services.PullResult
		err    error
	}, 1)

	go func() {
		result, err := pj.pullService.Pull(pj.filePath, pj.options)
		resultChan <- struct {
			result *services.PullResult
			err    error
		}{result: result, err: err}
	}()

	// Wait for either completion or cancellation
	select {
	case res := <-resultChan:
		// Pull completed
		pj.mu.Lock()
		pj.result = res.result
		pj.err = res.err
		pj.mu.Unlock()

		return res.err

	case <-ctx.Done():
		// Job cancelled
		cancelMu.Lock()
		cancelled = true
		cancelMu.Unlock()

		// Note: The PullService.Pull() call will continue in the background
		// since it doesn't support context cancellation yet.
		// This is a known limitation documented in the architecture doc.

		// Send cancellation progress
		select {
		case pj.progressChan <- JobProgress{
			JobID:      pj.id,
			Current:    0,
			Total:      0,
			Percentage: 0,
			Message:    "Cancelling...",
		}:
		default:
		}

		return ctx.Err()
	}
}

// Result returns the pull result if the job has completed.
// Returns nil if the job hasn't finished or was cancelled.
func (pj *PullJob) Result() *services.PullResult {
	pj.mu.Lock()
	defer pj.mu.Unlock()
	return pj.result
}

// Error returns the error if the job failed.
// Returns nil if the job succeeded or hasn't finished.
func (pj *PullJob) Error() error {
	pj.mu.Lock()
	defer pj.mu.Unlock()
	return pj.err
}

// FormatResult returns a human-readable summary of the pull result.
func (pj *PullJob) FormatResult() string {
	pj.mu.Lock()
	result := pj.result
	err := pj.err
	pj.mu.Unlock()

	if err != nil {
		return fmt.Sprintf("Pull failed: %v", err)
	}

	if result == nil {
		return "No result available"
	}

	parts := []string{}
	if result.TicketsPulled > 0 {
		parts = append(parts, fmt.Sprintf("%d new", result.TicketsPulled))
	}
	if result.TicketsUpdated > 0 {
		parts = append(parts, fmt.Sprintf("%d updated", result.TicketsUpdated))
	}
	if result.TicketsSkipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", result.TicketsSkipped))
	}

	if len(parts) == 0 {
		return "No changes"
	}

	msg := "Pull complete: "
	for i, part := range parts {
		if i > 0 {
			msg += ", "
		}
		msg += part
	}

	// Add conflict/error info
	if len(result.Conflicts) > 0 {
		msg += fmt.Sprintf(" (%d conflicts)", len(result.Conflicts))
	}
	if len(result.Errors) > 0 {
		msg += fmt.Sprintf(" (%d errors)", len(result.Errors))
	}

	return msg
}
