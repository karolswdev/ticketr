package ports

import (
	"context"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// JiraProgressCallback is called to report progress during ticket fetching operations
// current: number of items processed so far
// total: total number of items to process
// message: human-readable status message
type JiraProgressCallback func(current, total int, message string)

// JiraPort defines the interface for Jira integration operations
type JiraPort interface {
	// Authenticate verifies the connection to Jira with the provided credentials
	Authenticate() error

	// CreateTask creates a new sub-task in Jira under the specified parent
	CreateTask(task domain.Task, parentID string) (string, error)

	// UpdateTask updates an existing task in Jira
	UpdateTask(task domain.Task) error

	// GetProjectIssueTypes fetches available issue types for the configured project
	GetProjectIssueTypes() (map[string][]string, error)

	// GetIssueTypeFields fetches field requirements for a specific issue type
	GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error)

	// CreateTicket creates a new ticket in Jira with dynamic field mapping
	CreateTicket(ticket domain.Ticket) (string, error)

	// UpdateTicket updates an existing ticket in Jira with dynamic field mapping
	UpdateTicket(ticket domain.Ticket) error

	// SearchTickets searches for tickets in Jira using JQL query with context support for cancellation.
	// progressCallback can be nil if progress reporting is not needed.
	SearchTickets(ctx context.Context, projectKey string, jql string, progressCallback JiraProgressCallback) ([]domain.Ticket, error)
}
