package domain

type Ticket struct {
	Title              string
	Description        string
	CustomFields       map[string]string
	AcceptanceCriteria []string
	JiraID             string
	Tasks              []Task
	SourceLine         int
}

type Task struct {
	Title              string
	Description        string
	CustomFields       map[string]string // Task-specific overrides
	AcceptanceCriteria []string
	JiraID             string
	SourceLine         int
}

// Legacy types for backward compatibility - will be removed after full migration
type Story struct {
	Title              string
	Description        string
	AcceptanceCriteria []string
	JiraID             string
	Tasks              []Task
}