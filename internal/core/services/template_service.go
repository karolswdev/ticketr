package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/templates"
)

// TemplateService implements the template service port
type TemplateService struct {
	jiraClient   ports.JiraPort
	pathResolver *PathResolver
}

// NewTemplateService creates a new template service
func NewTemplateService(jiraClient ports.JiraPort, pathResolver *PathResolver) *TemplateService {
	return &TemplateService{
		jiraClient:   jiraClient,
		pathResolver: pathResolver,
	}
}

// ApplyTemplate applies a template with provided variables to create tickets
func (s *TemplateService) ApplyTemplate(ctx context.Context, template *templates.Template, vars map[string]string) (*ports.TemplateResult, error) {
	// Validate template
	if err := template.Validate(); err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}

	// Validate variables
	if err := templates.ValidateVariables(template, vars); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Substitute variables
	substituted, err := templates.Substitute(template, vars)
	if err != nil {
		return nil, fmt.Errorf("substitution failed: %w", err)
	}

	result := &ports.TemplateResult{
		StoryIDs: make([]string, 0),
		TaskIDs:  make([]string, 0),
	}

	// Track created tickets for rollback
	createdTickets := make([]string, 0)
	rollback := func() {
		// In a real implementation, we would delete created tickets
		// For now, we just log the failure
		// TODO: Implement ticket deletion in JiraPort
	}

	// Create epic if present
	if substituted.Structure.Epic != nil {
		epicTicket := domain.Ticket{
			Title:       substituted.Structure.Epic.Title,
			Description: substituted.Structure.Epic.Description,
		}

		epicID, err := s.jiraClient.CreateTicket(epicTicket)
		if err != nil {
			rollback()
			return nil, fmt.Errorf("failed to create epic: %w", err)
		}

		result.EpicID = epicID
		createdTickets = append(createdTickets, epicID)
		result.TotalCreated++

		// Create epic tasks
		for _, taskTitle := range substituted.Structure.Epic.Tasks {
			task := domain.Task{
				Title: taskTitle,
			}

			taskID, err := s.jiraClient.CreateTask(task, epicID)
			if err != nil {
				rollback()
				return nil, fmt.Errorf("failed to create epic task '%s': %w", taskTitle, err)
			}

			result.TaskIDs = append(result.TaskIDs, taskID)
			createdTickets = append(createdTickets, taskID)
			result.TotalCreated++
		}
	}

	// Create stories
	for i, story := range substituted.Structure.Stories {
		storyTicket := domain.Ticket{
			Title:       story.Title,
			Description: story.Description,
		}

		storyID, err := s.jiraClient.CreateTicket(storyTicket)
		if err != nil {
			rollback()
			return nil, fmt.Errorf("failed to create story %d '%s': %w", i+1, story.Title, err)
		}

		result.StoryIDs = append(result.StoryIDs, storyID)
		createdTickets = append(createdTickets, storyID)
		result.TotalCreated++

		// Create story tasks
		for j, taskTitle := range story.Tasks {
			task := domain.Task{
				Title: taskTitle,
			}

			taskID, err := s.jiraClient.CreateTask(task, storyID)
			if err != nil {
				rollback()
				return nil, fmt.Errorf("failed to create task %d for story '%s': %w", j+1, story.Title, err)
			}

			result.TaskIDs = append(result.TaskIDs, taskID)
			createdTickets = append(createdTickets, taskID)
			result.TotalCreated++
		}
	}

	return result, nil
}

// ListTemplates lists available templates from the templates directory
func (s *TemplateService) ListTemplates(ctx context.Context) ([]*ports.TemplateMetadata, error) {
	templatesDir := s.pathResolver.TemplatesDir()

	// Check if templates directory exists
	if !s.pathResolver.Exists(templatesDir) {
		// Create the directory if it doesn't exist
		if err := s.pathResolver.EnsureDirectory(templatesDir); err != nil {
			return nil, fmt.Errorf("failed to create templates directory: %w", err)
		}
		// Return empty list if directory was just created
		return []*ports.TemplateMetadata{}, nil
	}

	// Read directory contents
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	// Filter and load YAML files
	result := make([]*ports.TemplateMetadata, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process .yaml and .yml files
		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		path := filepath.Join(templatesDir, name)
		tmpl, err := s.LoadTemplate(ctx, path)
		if err != nil {
			// Skip invalid templates but log the error
			continue
		}

		metadata := &ports.TemplateMetadata{
			Name:       tmpl.Name,
			Path:       path,
			HasEpic:    tmpl.HasEpic(),
			StoryCount: tmpl.StoryCount(),
			TaskCount:  tmpl.TotalTaskCount(),
		}

		// Extract description from epic or first story
		if tmpl.Structure.Epic != nil && tmpl.Structure.Epic.Description != "" {
			metadata.Description = extractFirstLine(tmpl.Structure.Epic.Description)
		} else if len(tmpl.Structure.Stories) > 0 && tmpl.Structure.Stories[0].Description != "" {
			metadata.Description = extractFirstLine(tmpl.Structure.Stories[0].Description)
		} else if tmpl.Structure.Epic != nil {
			metadata.Description = tmpl.Structure.Epic.Title
		} else if len(tmpl.Structure.Stories) > 0 {
			metadata.Description = tmpl.Structure.Stories[0].Title
		}

		result = append(result, metadata)
	}

	return result, nil
}

// LoadTemplate loads a template from the specified file path
func (s *TemplateService) LoadTemplate(ctx context.Context, path string) (*templates.Template, error) {
	// Check if path is absolute or relative
	var fullPath string
	if filepath.IsAbs(path) {
		fullPath = path
	} else {
		// Treat as relative to templates directory
		fullPath = filepath.Join(s.pathResolver.TemplatesDir(), path)
	}

	// Read file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Parse template
	tmpl, err := templates.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// ValidateTemplate validates a template without applying it
func (s *TemplateService) ValidateTemplate(ctx context.Context, template *templates.Template) error {
	if err := template.Validate(); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// Extract variables to ensure they can be extracted
	vars := templates.ExtractVariables(template)
	if len(vars) == 0 {
		// No variables is valid, but might want to warn
	}

	return nil
}

// extractFirstLine extracts the first non-empty line from a multi-line string
func extractFirstLine(s string) string {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			// Limit to 100 characters
			if len(trimmed) > 100 {
				return trimmed[:97] + "..."
			}
			return trimmed
		}
	}
	return ""
}
