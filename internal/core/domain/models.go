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
