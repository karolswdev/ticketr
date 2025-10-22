package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockJiraClient for testing
type mockJiraClient struct {
	createTicketFunc func(ticket domain.Ticket) (string, error)
	createTaskFunc   func(task domain.Task, parentID string) (string, error)
	ticketCounter    int
	taskCounter      int
}

func (m *mockJiraClient) Authenticate() error {
	return nil
}

func (m *mockJiraClient) CreateTicket(ticket domain.Ticket) (string, error) {
	if m.createTicketFunc != nil {
		return m.createTicketFunc(ticket)
	}
	m.ticketCounter++
	return "TICKET-" + string(rune('0'+m.ticketCounter)), nil
}

func (m *mockJiraClient) UpdateTicket(ticket domain.Ticket) error {
	return nil
}

func (m *mockJiraClient) CreateTask(task domain.Task, parentID string) (string, error) {
	if m.createTaskFunc != nil {
		return m.createTaskFunc(task, parentID)
	}
	m.taskCounter++
	return "TASK-" + string(rune('0'+m.taskCounter)), nil
}

func (m *mockJiraClient) UpdateTask(task domain.Task) error {
	return nil
}

func (m *mockJiraClient) GetProjectIssueTypes() (map[string][]string, error) {
	return nil, nil
}

func (m *mockJiraClient) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *mockJiraClient) SearchTickets(ctx context.Context, projectKey string, jql string, progressCallback ports.JiraProgressCallback) ([]domain.Ticket, error) {
	return nil, nil
}

func setupTestTemplateService(t *testing.T) (*TemplateService, *PathResolver, string) {
	// Create temporary directory for templates
	tempDir := t.TempDir()

	// Create mock path resolver
	pathResolver, err := NewPathResolverWithOptions(
		"ticketr-test",
		func(key string) string { return "" },
		func() (string, error) { return tempDir, nil },
	)
	require.NoError(t, err)

	// Create templates directory
	err = pathResolver.EnsureDirectory(pathResolver.TemplatesDir())
	require.NoError(t, err)

	// Create mock Jira client
	jiraClient := &mockJiraClient{}

	// Create service
	service := NewTemplateService(jiraClient, pathResolver)

	return service, pathResolver, tempDir
}

func TestApplyTemplate_EpicOnly(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	yaml := `
name: epic-only
structure:
  epic:
    title: "Epic: {{.Name}}"
    description: "Epic for {{.Name}}"
`
	tmpl, err := templates.Parse([]byte(yaml))
	require.NoError(t, err)

	vars := map[string]string{"Name": "Test Feature"}

	result, err := service.ApplyTemplate(context.Background(), tmpl, vars)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEmpty(t, result.EpicID)
	assert.Equal(t, 0, len(result.StoryIDs))
	assert.Equal(t, 0, len(result.TaskIDs))
	assert.Equal(t, 1, result.TotalCreated)
}

func TestApplyTemplate_StoriesOnly(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	yaml := `
name: stories-only
structure:
  stories:
    - title: "Story 1: {{.Name}}"
      description: "First story"
    - title: "Story 2: {{.Name}}"
      description: "Second story"
`
	tmpl, err := templates.Parse([]byte(yaml))
	require.NoError(t, err)

	vars := map[string]string{"Name": "Feature"}

	result, err := service.ApplyTemplate(context.Background(), tmpl, vars)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Empty(t, result.EpicID)
	assert.Equal(t, 2, len(result.StoryIDs))
	assert.Equal(t, 0, len(result.TaskIDs))
	assert.Equal(t, 2, result.TotalCreated)
}

func TestApplyTemplate_WithTasks(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	yaml := `
name: with-tasks
structure:
  epic:
    title: "Epic: {{.Name}}"
    tasks:
      - "Epic task 1"
      - "Epic task 2"
  stories:
    - title: "Story: {{.Name}}"
      tasks:
        - "Story task 1"
        - "Story task 2"
`
	tmpl, err := templates.Parse([]byte(yaml))
	require.NoError(t, err)

	vars := map[string]string{"Name": "Feature"}

	result, err := service.ApplyTemplate(context.Background(), tmpl, vars)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEmpty(t, result.EpicID)
	assert.Equal(t, 1, len(result.StoryIDs))
	assert.Equal(t, 4, len(result.TaskIDs)) // 2 epic tasks + 2 story tasks
	assert.Equal(t, 6, result.TotalCreated) // 1 epic + 1 story + 4 tasks
}

func TestApplyTemplate_MissingVariable(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	yaml := `
name: test
structure:
  epic:
    title: "{{.Name}}"
    description: "{{.Description}}"
`
	tmpl, err := templates.Parse([]byte(yaml))
	require.NoError(t, err)

	// Only provide one variable
	vars := map[string]string{"Name": "Feature"}

	_, err = service.ApplyTemplate(context.Background(), tmpl, vars)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestApplyTemplate_InvalidTemplate(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	// Create invalid template (empty title)
	tmpl := &templates.Template{
		Name: "invalid",
		Structure: templates.TemplateStructure{
			Epic: &templates.TicketTemplate{
				Title: "", // Invalid: empty title
			},
		},
	}

	vars := map[string]string{}

	_, err := service.ApplyTemplate(context.Background(), tmpl, vars)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid template")
}

func TestListTemplates_EmptyDirectory(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	templates, err := service.ListTemplates(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, len(templates))
}

func TestListTemplates_MultipleTemplates(t *testing.T) {
	service, pathResolver, _ := setupTestTemplateService(t)

	// Create test templates
	template1 := `
name: template1
structure:
  epic:
    title: "Epic 1"
    description: "First template"
`
	template2 := `
name: template2
structure:
  stories:
    - title: "Story 1"
    - title: "Story 2"
`

	templatesDir := pathResolver.TemplatesDir()
	err := os.WriteFile(filepath.Join(templatesDir, "template1.yaml"), []byte(template1), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(templatesDir, "template2.yml"), []byte(template2), 0644)
	require.NoError(t, err)

	// Create a non-YAML file that should be ignored
	err = os.WriteFile(filepath.Join(templatesDir, "readme.txt"), []byte("Ignore me"), 0644)
	require.NoError(t, err)

	templates, err := service.ListTemplates(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, len(templates))

	// Check first template
	var tmpl1, tmpl2 *ports.TemplateMetadata
	for _, tmpl := range templates {
		if tmpl.Name == "template1" {
			tmpl1 = tmpl
		} else if tmpl.Name == "template2" {
			tmpl2 = tmpl
		}
	}

	require.NotNil(t, tmpl1)
	assert.Equal(t, "template1", tmpl1.Name)
	assert.True(t, tmpl1.HasEpic)
	assert.Equal(t, 0, tmpl1.StoryCount)
	assert.Equal(t, "First template", tmpl1.Description)

	require.NotNil(t, tmpl2)
	assert.Equal(t, "template2", tmpl2.Name)
	assert.False(t, tmpl2.HasEpic)
	assert.Equal(t, 2, tmpl2.StoryCount)
}

func TestListTemplates_InvalidTemplateSkipped(t *testing.T) {
	service, pathResolver, _ := setupTestTemplateService(t)

	// Create valid template
	validTemplate := `
name: valid
structure:
  epic:
    title: "Valid Epic"
`
	// Create invalid template (missing name)
	invalidTemplate := `
structure:
  epic:
    title: "No Name"
`

	templatesDir := pathResolver.TemplatesDir()
	err := os.WriteFile(filepath.Join(templatesDir, "valid.yaml"), []byte(validTemplate), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(templatesDir, "invalid.yaml"), []byte(invalidTemplate), 0644)
	require.NoError(t, err)

	// Should only return valid template
	templates, err := service.ListTemplates(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, len(templates))
	assert.Equal(t, "valid", templates[0].Name)
}

func TestLoadTemplate_AbsolutePath(t *testing.T) {
	service, pathResolver, _ := setupTestTemplateService(t)

	templateYAML := `
name: test-template
structure:
  epic:
    title: "Test Epic"
`
	templatesDir := pathResolver.TemplatesDir()
	templatePath := filepath.Join(templatesDir, "test.yaml")
	err := os.WriteFile(templatePath, []byte(templateYAML), 0644)
	require.NoError(t, err)

	tmpl, err := service.LoadTemplate(context.Background(), templatePath)
	require.NoError(t, err)
	assert.Equal(t, "test-template", tmpl.Name)
}

func TestLoadTemplate_RelativePath(t *testing.T) {
	service, pathResolver, _ := setupTestTemplateService(t)

	templateYAML := `
name: test-template
structure:
  epic:
    title: "Test Epic"
`
	templatesDir := pathResolver.TemplatesDir()
	err := os.WriteFile(filepath.Join(templatesDir, "test.yaml"), []byte(templateYAML), 0644)
	require.NoError(t, err)

	// Load using relative path
	tmpl, err := service.LoadTemplate(context.Background(), "test.yaml")
	require.NoError(t, err)
	assert.Equal(t, "test-template", tmpl.Name)
}

func TestLoadTemplate_FileNotFound(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	_, err := service.LoadTemplate(context.Background(), "nonexistent.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read template file")
}

func TestLoadTemplate_InvalidYAML(t *testing.T) {
	service, pathResolver, _ := setupTestTemplateService(t)

	invalidYAML := "{ this is not valid yaml"
	templatesDir := pathResolver.TemplatesDir()
	templatePath := filepath.Join(templatesDir, "invalid.yaml")
	err := os.WriteFile(templatePath, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	_, err = service.LoadTemplate(context.Background(), templatePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse template")
}

func TestValidateTemplate_Valid(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	yaml := `
name: valid
structure:
  epic:
    title: "Valid Epic"
`
	tmpl, err := templates.Parse([]byte(yaml))
	require.NoError(t, err)

	err = service.ValidateTemplate(context.Background(), tmpl)
	assert.NoError(t, err)
}

func TestValidateTemplate_Invalid(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	// Create invalid template directly
	tmpl := &templates.Template{
		Name:      "", // Invalid: missing name
		Structure: templates.TemplateStructure{},
	}

	err := service.ValidateTemplate(context.Background(), tmpl)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid template")
}

func TestExtractFirstLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single line",
			input: "Single line description",
			want:  "Single line description",
		},
		{
			name:  "multi-line",
			input: "First line\nSecond line\nThird line",
			want:  "First line",
		},
		{
			name:  "with leading whitespace",
			input: "\n\n  First line after whitespace",
			want:  "First line after whitespace",
		},
		{
			name:  "very long line",
			input: "This is a very long line that exceeds the maximum length limit and should be truncated to avoid overflow",
			want:  "This is a very long line that exceeds the maximum length limit and should be truncated to avoid o...",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only whitespace",
			input: "   \n\n   ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFirstLine(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestApplyTemplate_ComplexHierarchy(t *testing.T) {
	service, _, _ := setupTestTemplateService(t)

	yaml := `
name: feature-complete
structure:
  epic:
    title: "Feature: {{.Name}}"
    description: "For {{.Actor}}"
    tasks:
      - "Research {{.Name}}"
      - "Design {{.Name}}"
  stories:
    - title: "Frontend: {{.Name}}"
      description: "UI for {{.Name}}"
      tasks:
        - "Component"
        - "Tests"
    - title: "Backend: {{.Name}}"
      description: "API for {{.Name}}"
      tasks:
        - "Endpoint"
        - "Schema"
        - "Tests"
`
	tmpl, err := templates.Parse([]byte(yaml))
	require.NoError(t, err)

	vars := map[string]string{
		"Name":  "Authentication",
		"Actor": "user",
	}

	result, err := service.ApplyTemplate(context.Background(), tmpl, vars)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEmpty(t, result.EpicID)
	assert.Equal(t, 2, len(result.StoryIDs))
	assert.Equal(t, 7, len(result.TaskIDs))  // 2 epic + 2 frontend + 3 backend
	assert.Equal(t, 10, result.TotalCreated) // 1 epic + 2 stories + 7 tasks
}

func TestListTemplates_DescriptionExtraction(t *testing.T) {
	service, pathResolver, _ := setupTestTemplateService(t)

	// Template with epic description
	template1 := `
name: epic-desc
structure:
  epic:
    title: "Epic Title"
    description: "Epic description line"
`
	// Template with story description
	template2 := `
name: story-desc
structure:
  stories:
    - title: "Story Title"
      description: "Story description line"
`
	// Template with only titles
	template3 := `
name: no-desc
structure:
  epic:
    title: "Just a title"
`

	templatesDir := pathResolver.TemplatesDir()
	os.WriteFile(filepath.Join(templatesDir, "t1.yaml"), []byte(template1), 0644)
	os.WriteFile(filepath.Join(templatesDir, "t2.yaml"), []byte(template2), 0644)
	os.WriteFile(filepath.Join(templatesDir, "t3.yaml"), []byte(template3), 0644)

	templates, err := service.ListTemplates(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 3, len(templates))

	// Find templates by name
	templateMap := make(map[string]*ports.TemplateMetadata)
	for _, tmpl := range templates {
		templateMap[tmpl.Name] = tmpl
	}

	assert.Equal(t, "Epic description line", templateMap["epic-desc"].Description)
	assert.Equal(t, "Story description line", templateMap["story-desc"].Description)
	assert.Equal(t, "Just a title", templateMap["no-desc"].Description)
}
