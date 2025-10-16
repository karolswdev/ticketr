package ports

import "github.com/karolswdev/ticktr/internal/core/domain"

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

	// SearchTickets searches for tickets in Jira using JQL query
	SearchTickets(projectKey string, jql string) ([]domain.Ticket, error)
}
