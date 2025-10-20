package ports

import (
	"context"

	"github.com/karolswdev/ticktr/internal/templates"
)

// TemplateService handles template-based ticket creation
type TemplateService interface {
	// ApplyTemplate applies a template with provided variables to create tickets
	// Returns the IDs of created tickets or an error if creation fails
	ApplyTemplate(ctx context.Context, template *templates.Template, vars map[string]string) (*TemplateResult, error)

	// ListTemplates lists available templates from the templates directory
	// Returns a slice of template metadata or an error
	ListTemplates(ctx context.Context) ([]*TemplateMetadata, error)

	// LoadTemplate loads a template from the specified file path
	// Returns the parsed template or an error if loading/parsing fails
	LoadTemplate(ctx context.Context, path string) (*templates.Template, error)

	// ValidateTemplate validates a template without applying it
	// Returns nil if valid, error otherwise
	ValidateTemplate(ctx context.Context, template *templates.Template) error
}

// TemplateResult contains the IDs of tickets created by applying a template
type TemplateResult struct {
	// EpicID is the Jira ID of the created epic (empty if no epic in template)
	EpicID string

	// StoryIDs are the Jira IDs of created stories
	StoryIDs []string

	// TaskIDs are the Jira IDs of created tasks (includes tasks under epic and stories)
	TaskIDs []string

	// TotalCreated is the total number of tickets created
	TotalCreated int
}

// TemplateMetadata contains metadata about a template file
type TemplateMetadata struct {
	// Name is the template name from the YAML
	Name string

	// Path is the absolute file path to the template
	Path string

	// Description is a brief description (extracted from first line of description or title)
	Description string

	// HasEpic indicates if the template defines an epic
	HasEpic bool

	// StoryCount is the number of stories in the template
	StoryCount int

	// TaskCount is the total number of tasks across all stories
	TaskCount int
}
