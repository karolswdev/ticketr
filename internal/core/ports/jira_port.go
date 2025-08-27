package ports

import "github.com/karolswdev/ticktr/internal/core/domain"

// JiraPort defines the interface for Jira integration operations
type JiraPort interface {
	// Authenticate verifies the connection to Jira with the provided credentials
	Authenticate() error
	
	// CreateStory creates a new story in Jira and returns the created issue key
	CreateStory(story domain.Story) (string, error)
	
	// CreateTask creates a new sub-task in Jira under the specified parent story
	CreateTask(task domain.Task, parentID string) (string, error)
}