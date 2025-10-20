package ports

import (
	"context"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// BulkOperationProgressCallback is invoked after each ticket operation.
// It provides real-time progress updates for CLI/TUI display.
// Parameters:
//   - ticketID: The Jira ticket ID that was processed
//   - success: Whether the operation succeeded for this ticket
//   - err: Error details if success is false, nil otherwise
type BulkOperationProgressCallback func(ticketID string, success bool, err error)

// BulkOperationService defines the interface for performing bulk operations on tickets.
// Implementations handle the coordination of bulk actions across multiple tickets,
// including validation, execution, and result aggregation.
type BulkOperationService interface {
	// ExecuteOperation performs a bulk operation on multiple tickets.
	// It validates the operation, executes it across all target tickets,
	// and returns a detailed result including success/failure counts.
	//
	// The operation respects context cancellation and will stop processing
	// if the context is cancelled, returning partial results.
	//
	// The progress callback (if provided) is invoked after each ticket operation.
	// Pass nil if progress updates are not needed (e.g., in tests).
	//
	// For update and move operations, it attempts best-effort rollback on partial
	// failure. Delete operations cannot be rolled back.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - op: The bulk operation to execute (already validated via op.Validate())
	//   - progress: Optional callback for progress updates (nil-safe)
	//
	// Returns:
	//   - *domain.BulkOperationResult: Detailed results including success/failure counts and errors
	//   - error: Returns error only for operation-level failures (not individual ticket failures)
	//
	// Individual ticket failures are captured in the result's Errors map,
	// while operation-level errors (e.g., network failure, invalid credentials) are returned directly.
	ExecuteOperation(
		ctx context.Context,
		op *domain.BulkOperation,
		progress BulkOperationProgressCallback,
	) (*domain.BulkOperationResult, error)
}
