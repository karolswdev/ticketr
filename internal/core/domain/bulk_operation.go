package domain

import (
	"fmt"
	"regexp"
)

// BulkOperationAction represents the valid actions that can be performed in bulk.
type BulkOperationAction string

const (
	// BulkActionUpdate modifies field values on multiple tickets.
	BulkActionUpdate BulkOperationAction = "update"
	// BulkActionMove transitions multiple tickets to a new status.
	BulkActionMove BulkOperationAction = "move"
	// BulkActionDelete removes multiple tickets from the system.
	BulkActionDelete BulkOperationAction = "delete"
)

// jiraIDRegex validates Jira ticket ID format: PROJECT-123
// Project key must be uppercase letters, followed by hyphen, followed by digits.
var jiraIDRegex = regexp.MustCompile(`^[A-Z]+-\d+$`)

const (
	// MinTicketCount is the minimum number of tickets required for a bulk operation.
	MinTicketCount = 1
	// MaxTicketCount is the maximum number of tickets allowed in a single bulk operation.
	MaxTicketCount = 100
)

// BulkOperation represents a bulk action to be performed on multiple tickets.
// It encapsulates the action type, target tickets, and any changes to be applied.
type BulkOperation struct {
	// Action specifies the bulk operation to perform (update, move, delete).
	Action BulkOperationAction

	// TicketIDs contains the Jira ticket IDs (e.g., "PROJ-123") to operate on.
	// Must contain between 1 and 100 ticket IDs.
	TicketIDs []string

	// Changes contains the field modifications to apply (for update operations).
	// Key is the field name, value is the new value.
	// Required for "update" action, ignored for other actions.
	Changes map[string]interface{}
}

// BulkOperationResult represents the outcome of a bulk operation.
// It tracks which tickets succeeded, failed, and any error messages.
type BulkOperationResult struct {
	// SuccessCount is the number of tickets successfully processed.
	SuccessCount int

	// FailureCount is the number of tickets that failed to process.
	FailureCount int

	// Errors contains any errors encountered during processing.
	// Map key is the ticket ID, value is the error message.
	Errors map[string]string

	// SuccessfulTickets contains the IDs of tickets that were processed successfully.
	SuccessfulTickets []string

	// FailedTickets contains the IDs of tickets that failed to process.
	FailedTickets []string
}

// NewBulkOperation creates a new BulkOperation with the specified action and ticket IDs.
// Changes can be nil for non-update operations.
func NewBulkOperation(action BulkOperationAction, ticketIDs []string, changes map[string]interface{}) *BulkOperation {
	return &BulkOperation{
		Action:    action,
		TicketIDs: ticketIDs,
		Changes:   changes,
	}
}

// validateJiraID ensures ticket ID matches Jira format: PROJECT-123
// This prevents JQL injection attacks through malformed ticket IDs.
func validateJiraID(id string) error {
	if !jiraIDRegex.MatchString(id) {
		return fmt.Errorf("invalid Jira ID format: %s (expected format: PROJECT-123)", id)
	}
	return nil
}

// Validate checks if the BulkOperation has valid values according to business rules.
func (bo *BulkOperation) Validate() error {
	// Validate action
	if err := bo.validateAction(); err != nil {
		return err
	}

	// Validate ticket IDs
	if err := bo.validateTicketIDs(); err != nil {
		return err
	}

	// Validate changes (action-specific)
	if err := bo.validateChanges(); err != nil {
		return err
	}

	return nil
}

// validateAction ensures the action is a valid operation type.
func (bo *BulkOperation) validateAction() error {
	switch bo.Action {
	case BulkActionUpdate, BulkActionMove, BulkActionDelete:
		return nil
	default:
		return fmt.Errorf("invalid action: must be one of 'update', 'move', or 'delete'")
	}
}

// validateTicketIDs ensures the ticket ID list meets requirements.
func (bo *BulkOperation) validateTicketIDs() error {
	// Check count constraints
	if len(bo.TicketIDs) < MinTicketCount {
		return fmt.Errorf("ticket_ids cannot be empty: must contain at least %d ticket", MinTicketCount)
	}

	if len(bo.TicketIDs) > MaxTicketCount {
		return fmt.Errorf("ticket_ids cannot exceed %d tickets: found %d", MaxTicketCount, len(bo.TicketIDs))
	}

	// Check for empty strings and validate format
	for i, id := range bo.TicketIDs {
		if id == "" {
			return fmt.Errorf("ticket_ids[%d] cannot be an empty string", i)
		}

		// Validate Jira ID format to prevent JQL injection
		if err := validateJiraID(id); err != nil {
			return fmt.Errorf("ticket_ids[%d]: %w", i, err)
		}
	}

	return nil
}

// validateChanges ensures changes are provided when required (update action).
func (bo *BulkOperation) validateChanges() error {
	if bo.Action == BulkActionUpdate {
		if bo.Changes == nil || len(bo.Changes) == 0 {
			return fmt.Errorf("changes are required for 'update' action")
		}
	}

	return nil
}

// IsValidAction checks if the given action string is valid.
func IsValidAction(action string) bool {
	switch BulkOperationAction(action) {
	case BulkActionUpdate, BulkActionMove, BulkActionDelete:
		return true
	default:
		return false
	}
}

// TicketCount returns the number of tickets in this bulk operation.
func (bo *BulkOperation) TicketCount() int {
	return len(bo.TicketIDs)
}

// NewBulkOperationResult creates a new result object with initialized maps and slices.
func NewBulkOperationResult() *BulkOperationResult {
	return &BulkOperationResult{
		Errors:            make(map[string]string),
		SuccessfulTickets: make([]string, 0),
		FailedTickets:     make([]string, 0),
	}
}
