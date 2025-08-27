package domain

// Story represents a user story in the project management system.
type Story struct {
	// Title is the main heading of the story
	Title string
	// Description provides detailed information about the story
	Description string
	// AcceptanceCriteria defines the conditions that must be met for the story to be considered complete
	AcceptanceCriteria []string
	// JiraID is the unique identifier from Jira (e.g., "PROJ-123"), empty if not yet created in Jira
	JiraID string
	// Tasks is a collection of tasks that need to be completed for this story
	Tasks []Task
}

// Task represents an individual task within a story.
type Task struct {
	// Title is the main heading of the task
	Title string
	// Description provides detailed information about the task
	Description string
	// AcceptanceCriteria defines the conditions that must be met for the task to be considered complete
	AcceptanceCriteria []string
	// JiraID is the unique identifier from Jira (e.g., "PROJ-124"), empty if not yet created in Jira
	JiraID string
}