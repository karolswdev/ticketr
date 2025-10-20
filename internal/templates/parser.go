package templates

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Template represents a parsed template file
type Template struct {
	Name      string            `yaml:"name"`
	Structure TemplateStructure `yaml:"structure"`
}

// TemplateStructure defines the ticket hierarchy
type TemplateStructure struct {
	Epic    *TicketTemplate  `yaml:"epic,omitempty"`
	Stories []TicketTemplate `yaml:"stories,omitempty"`
}

// TicketTemplate defines a single ticket or story
type TicketTemplate struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description,omitempty"`
	Tasks       []string `yaml:"tasks,omitempty"`
}

// Parse parses a YAML template file
func Parse(yamlContent []byte) (*Template, error) {
	if len(yamlContent) == 0 {
		return nil, fmt.Errorf("empty template content")
	}

	var tmpl Template
	if err := yaml.Unmarshal(yamlContent, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if err := tmpl.Validate(); err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}

	return &tmpl, nil
}

// Validate validates the template structure
func (t *Template) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("template name is required")
	}

	// Template must have at least an epic or stories
	if t.Structure.Epic == nil && len(t.Structure.Stories) == 0 {
		return fmt.Errorf("template must define at least an epic or stories")
	}

	// Validate epic if present
	if t.Structure.Epic != nil {
		if err := t.Structure.Epic.Validate(); err != nil {
			return fmt.Errorf("invalid epic: %w", err)
		}
	}

	// Validate each story
	for i, story := range t.Structure.Stories {
		if err := story.Validate(); err != nil {
			return fmt.Errorf("invalid story at index %d: %w", i, err)
		}
	}

	return nil
}

// Validate validates a single ticket template
func (tt *TicketTemplate) Validate() error {
	if tt.Title == "" {
		return fmt.Errorf("title is required")
	}

	// Title must not be excessively long
	if len(tt.Title) > 255 {
		return fmt.Errorf("title exceeds maximum length of 255 characters")
	}

	// Validate tasks are not empty strings
	for i, task := range tt.Tasks {
		trimmed := strings.TrimSpace(task)
		if trimmed == "" {
			return fmt.Errorf("task at index %d is empty", i)
		}
	}

	return nil
}

// HasEpic returns true if the template defines an epic
func (t *Template) HasEpic() bool {
	return t.Structure.Epic != nil
}

// HasStories returns true if the template defines stories
func (t *Template) HasStories() bool {
	return len(t.Structure.Stories) > 0
}

// StoryCount returns the number of stories in the template
func (t *Template) StoryCount() int {
	return len(t.Structure.Stories)
}

// TotalTaskCount returns the total number of tasks across all stories
func (t *Template) TotalTaskCount() int {
	count := 0
	if t.Structure.Epic != nil {
		count += len(t.Structure.Epic.Tasks)
	}
	for _, story := range t.Structure.Stories {
		count += len(story.Tasks)
	}
	return count
}
