package services

import (
	"context"
	"fmt"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// ticketSnapshot stores the original state of a ticket for rollback purposes.
type ticketSnapshot struct {
	ticketID string
	ticket   domain.Ticket
}

// BulkOperationServiceImpl implements the BulkOperationService interface.
// It orchestrates bulk operations on tickets, handling validation, execution,
// progress tracking, and transaction rollback (best effort).
type BulkOperationServiceImpl struct {
	jiraAdapter ports.JiraPort
}

// NewBulkOperationService creates a new BulkOperationService instance.
func NewBulkOperationService(jiraAdapter ports.JiraPort) ports.BulkOperationService {
	return &BulkOperationServiceImpl{
		jiraAdapter: jiraAdapter,
	}
}

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
// Individual ticket failures are captured in the result's Errors map,
// while operation-level errors (e.g., invalid operation) are returned directly.
func (s *BulkOperationServiceImpl) ExecuteOperation(
	ctx context.Context,
	op *domain.BulkOperation,
	progress ports.BulkOperationProgressCallback,
) (*domain.BulkOperationResult, error) {
	// Validate the bulk operation
	if err := op.Validate(); err != nil {
		return nil, fmt.Errorf("invalid bulk operation: %w", err)
	}

	// Initialize result
	result := domain.NewBulkOperationResult()

	// Execute based on operation type
	switch op.Action {
	case domain.BulkActionUpdate:
		return s.executeUpdate(ctx, op, result, progress)
	case domain.BulkActionMove:
		return s.executeMove(ctx, op, result, progress)
	case domain.BulkActionDelete:
		return s.executeDelete(ctx, op, result, progress)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", op.Action)
	}
}

// executeUpdate performs bulk update operations on tickets.
// It stores original ticket state for rollback on partial failure.
func (s *BulkOperationServiceImpl) executeUpdate(
	ctx context.Context,
	op *domain.BulkOperation,
	result *domain.BulkOperationResult,
	progress ports.BulkOperationProgressCallback,
) (*domain.BulkOperationResult, error) {
	snapshots := make([]ticketSnapshot, 0, len(op.TicketIDs))

	// Process each ticket
	for _, ticketID := range op.TicketIDs {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return result, fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		// Fetch current ticket state (for rollback)
		tickets, err := s.jiraAdapter.SearchTickets("", fmt.Sprintf(`key = "%s"`, ticketID))
		if err != nil {
			s.recordFailure(result, ticketID, fmt.Errorf("failed to fetch ticket for backup: %w", err))

			// Invoke callback (nil-safe)
			if progress != nil {
				progress(ticketID, false, fmt.Errorf("failed to fetch ticket for backup: %w", err))
			}
			continue
		}
		if len(tickets) == 0 {
			s.recordFailure(result, ticketID, fmt.Errorf("ticket not found: %s", ticketID))

			// Invoke callback (nil-safe)
			if progress != nil {
				progress(ticketID, false, fmt.Errorf("ticket not found: %s", ticketID))
			}
			continue
		}

		// Store snapshot for rollback
		snapshots = append(snapshots, ticketSnapshot{
			ticketID: ticketID,
			ticket:   tickets[0],
		})

		// Apply changes to ticket
		updatedTicket := tickets[0]
		if updatedTicket.CustomFields == nil {
			updatedTicket.CustomFields = make(map[string]string)
		}
		for fieldName, fieldValue := range op.Changes {
			if strValue, ok := fieldValue.(string); ok {
				updatedTicket.CustomFields[fieldName] = strValue
			} else {
				updatedTicket.CustomFields[fieldName] = fmt.Sprintf("%v", fieldValue)
			}
		}

		// Update ticket in Jira
		if err := s.jiraAdapter.UpdateTicket(updatedTicket); err != nil {
			s.recordFailure(result, ticketID, fmt.Errorf("failed to update ticket: %w", err))

			// Invoke callback (nil-safe)
			if progress != nil {
				progress(ticketID, false, err)
			}
			continue
		}

		s.recordSuccess(result, ticketID)

		// Invoke callback (nil-safe)
		if progress != nil {
			progress(ticketID, true, nil)
		}
	}

	// If there were failures and some successes, attempt rollback
	if result.FailureCount > 0 && result.SuccessCount > 0 {
		s.rollbackUpdates(ctx, snapshots, result.SuccessfulTickets)
		return result, fmt.Errorf("partial failure: %d of %d tickets failed (rollback attempted)", result.FailureCount, len(op.TicketIDs))
	}

	// If all failed, return error
	if result.FailureCount == len(op.TicketIDs) {
		return result, fmt.Errorf("all tickets failed to update")
	}

	return result, nil
}

// executeMove performs bulk move operations on tickets.
// This updates the parent field (if supported by the Jira adapter).
func (s *BulkOperationServiceImpl) executeMove(
	ctx context.Context,
	op *domain.BulkOperation,
	result *domain.BulkOperationResult,
	progress ports.BulkOperationProgressCallback,
) (*domain.BulkOperationResult, error) {
	// For move operations, we need a "parent" field in Changes
	parentKey, ok := op.Changes["parent"]
	if !ok || parentKey == nil {
		return nil, fmt.Errorf("move operation requires 'parent' field in changes")
	}

	parentKeyStr, ok := parentKey.(string)
	if !ok {
		return nil, fmt.Errorf("parent field must be a string")
	}

	snapshots := make([]ticketSnapshot, 0, len(op.TicketIDs))

	// Process each ticket
	for _, ticketID := range op.TicketIDs {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return result, fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		// Fetch current ticket state (for rollback)
		tickets, err := s.jiraAdapter.SearchTickets("", fmt.Sprintf(`key = "%s"`, ticketID))
		if err != nil {
			s.recordFailure(result, ticketID, fmt.Errorf("failed to fetch ticket for backup: %w", err))

			// Invoke callback (nil-safe)
			if progress != nil {
				progress(ticketID, false, fmt.Errorf("failed to fetch ticket for backup: %w", err))
			}
			continue
		}
		if len(tickets) == 0 {
			s.recordFailure(result, ticketID, fmt.Errorf("ticket not found: %s", ticketID))

			// Invoke callback (nil-safe)
			if progress != nil {
				progress(ticketID, false, fmt.Errorf("ticket not found: %s", ticketID))
			}
			continue
		}

		// Store snapshot for rollback
		snapshots = append(snapshots, ticketSnapshot{
			ticketID: ticketID,
			ticket:   tickets[0],
		})

		// Update parent field
		updatedTicket := tickets[0]
		if updatedTicket.CustomFields == nil {
			updatedTicket.CustomFields = make(map[string]string)
		}
		updatedTicket.CustomFields["Parent"] = parentKeyStr

		// Update ticket in Jira
		if err := s.jiraAdapter.UpdateTicket(updatedTicket); err != nil {
			s.recordFailure(result, ticketID, fmt.Errorf("failed to move ticket: %w", err))

			// Invoke callback (nil-safe)
			if progress != nil {
				progress(ticketID, false, err)
			}
			continue
		}

		s.recordSuccess(result, ticketID)

		// Invoke callback (nil-safe)
		if progress != nil {
			progress(ticketID, true, nil)
		}
	}

	// If there were failures and some successes, attempt rollback
	if result.FailureCount > 0 && result.SuccessCount > 0 {
		s.rollbackUpdates(ctx, snapshots, result.SuccessfulTickets)
		return result, fmt.Errorf("partial failure: %d of %d tickets failed (rollback attempted)", result.FailureCount, len(op.TicketIDs))
	}

	// If all failed, return error
	if result.FailureCount == len(op.TicketIDs) {
		return result, fmt.Errorf("all tickets failed to move")
	}

	return result, nil
}

// executeDelete performs bulk delete operations on tickets.
// WARNING: Delete operations CANNOT be rolled back. Once deleted, tickets are gone.
func (s *BulkOperationServiceImpl) executeDelete(
	ctx context.Context,
	op *domain.BulkOperation,
	result *domain.BulkOperationResult,
	progress ports.BulkOperationProgressCallback,
) (*domain.BulkOperationResult, error) {
	// Process each ticket
	for _, ticketID := range op.TicketIDs {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return result, fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		// Note: Jira adapter doesn't currently have a DeleteTicket method
		// This is a limitation that should be documented
		deleteErr := fmt.Errorf("delete operation not supported by Jira adapter")
		s.recordFailure(result, ticketID, deleteErr)

		// Invoke callback (nil-safe)
		if progress != nil {
			progress(ticketID, false, deleteErr)
		}
	}

	// All deletes failed due to missing adapter support
	if result.FailureCount == len(op.TicketIDs) {
		return result, fmt.Errorf("delete operation not supported: Jira adapter does not implement DeleteTicket method")
	}

	return result, nil
}

// rollbackUpdates attempts to restore tickets to their original state.
// This is a best-effort operation - failures are logged but don't cascade errors.
func (s *BulkOperationServiceImpl) rollbackUpdates(ctx context.Context, snapshots []ticketSnapshot, successfulTickets []string) {
	// Create a map of successful tickets for quick lookup
	successMap := make(map[string]bool)
	for _, ticketID := range successfulTickets {
		successMap[ticketID] = true
	}

	// Attempt to restore each successfully updated ticket
	for _, snapshot := range snapshots {
		// Only rollback tickets that were successful
		if !successMap[snapshot.ticketID] {
			continue
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Attempt to restore original state
		if err := s.jiraAdapter.UpdateTicket(snapshot.ticket); err != nil {
			// Log error but continue (best effort)
			// In production, this should use proper logging
			// For now, we silently continue
			continue
		}
	}
}

// recordSuccess updates the result to reflect a successful operation on a ticket.
func (s *BulkOperationServiceImpl) recordSuccess(result *domain.BulkOperationResult, ticketID string) {
	result.SuccessCount++
	result.SuccessfulTickets = append(result.SuccessfulTickets, ticketID)
}

// recordFailure updates the result to reflect a failed operation on a ticket.
func (s *BulkOperationServiceImpl) recordFailure(result *domain.BulkOperationResult, ticketID string, err error) {
	result.FailureCount++
	result.FailedTickets = append(result.FailedTickets, ticketID)
	result.Errors[ticketID] = err.Error()
}
