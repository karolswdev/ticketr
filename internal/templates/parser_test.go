package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_ValidTemplate(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		wantName string
		wantEpic bool
		wantLen  int
	}{
		{
			name: "template with epic and stories",
			yaml: `
name: feature
structure:
  epic:
    title: "Feature: {{.Name}}"
    description: "Epic description"
  stories:
    - title: "Story 1"
      description: "Description 1"
    - title: "Story 2"
      description: "Description 2"
`,
			wantName: "feature",
			wantEpic: true,
			wantLen:  2,
		},
		{
			name: "template with only epic",
			yaml: `
name: epic-only
structure:
  epic:
    title: "Epic Title"
    description: "Epic description"
`,
			wantName: "epic-only",
			wantEpic: true,
			wantLen:  0,
		},
		{
			name: "template with only stories",
			yaml: `
name: stories-only
structure:
  stories:
    - title: "Story 1"
    - title: "Story 2"
    - title: "Story 3"
`,
			wantName: "stories-only",
			wantEpic: false,
			wantLen:  3,
		},
		{
			name: "template with tasks",
			yaml: `
name: with-tasks
structure:
  stories:
    - title: "Story with tasks"
      tasks:
        - "Task 1"
        - "Task 2"
        - "Task 3"
`,
			wantName: "with-tasks",
			wantEpic: false,
			wantLen:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := Parse([]byte(tt.yaml))
			require.NoError(t, err)
			assert.Equal(t, tt.wantName, tmpl.Name)
			assert.Equal(t, tt.wantEpic, tmpl.HasEpic())
			assert.Equal(t, tt.wantLen, tmpl.StoryCount())
		})
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr string
	}{
		{
			name:    "empty content",
			yaml:    "",
			wantErr: "empty template content",
		},
		{
			name: "malformed YAML",
			yaml: `
name: test
structure:
  epic:
    title: "Missing quote
`,
			wantErr: "failed to parse YAML",
		},
		{
			name: "invalid structure",
			yaml: `
name: test
structure: "not a structure"
`,
			wantErr: "failed to parse YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.yaml))
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestValidate_MissingName(t *testing.T) {
	yaml := `
structure:
  epic:
    title: "Epic Title"
`
	_, err := Parse([]byte(yaml))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "template name is required")
}

func TestValidate_EmptyStructure(t *testing.T) {
	yaml := `
name: empty
structure: {}
`
	_, err := Parse([]byte(yaml))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "template must define at least an epic or stories")
}

func TestValidate_MissingTitle(t *testing.T) {
	yaml := `
name: missing-title
structure:
  epic:
    description: "Epic without title"
`
	_, err := Parse([]byte(yaml))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

func TestValidate_TitleTooLong(t *testing.T) {
	longTitle := make([]byte, 300)
	for i := range longTitle {
		longTitle[i] = 'A'
	}

	yaml := "name: long-title\nstructure:\n  epic:\n    title: \"" + string(longTitle) + "\""
	_, err := Parse([]byte(yaml))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum length")
}

func TestValidate_EmptyTask(t *testing.T) {
	yaml := `
name: empty-task
structure:
  stories:
    - title: "Story"
      tasks:
        - "Valid task"
        - ""
        - "Another valid task"
`
	_, err := Parse([]byte(yaml))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task at index 1 is empty")
}

func TestValidate_WhitespaceTask(t *testing.T) {
	yaml := `
name: whitespace-task
structure:
  stories:
    - title: "Story"
      tasks:
        - "Valid task"
        - "   "
`
	_, err := Parse([]byte(yaml))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task at index 1 is empty")
}

func TestTemplate_HasEpic(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want bool
	}{
		{
			name: "has epic",
			yaml: `
name: test
structure:
  epic:
    title: "Epic"
`,
			want: true,
		},
		{
			name: "no epic",
			yaml: `
name: test
structure:
  stories:
    - title: "Story"
`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := Parse([]byte(tt.yaml))
			require.NoError(t, err)
			assert.Equal(t, tt.want, tmpl.HasEpic())
		})
	}
}

func TestTemplate_StoryCount(t *testing.T) {
	yaml := `
name: test
structure:
  stories:
    - title: "Story 1"
    - title: "Story 2"
    - title: "Story 3"
`
	tmpl, err := Parse([]byte(yaml))
	require.NoError(t, err)
	assert.Equal(t, 3, tmpl.StoryCount())
}

func TestTemplate_TotalTaskCount(t *testing.T) {
	yaml := `
name: test
structure:
  epic:
    title: "Epic"
    tasks:
      - "Epic task 1"
      - "Epic task 2"
  stories:
    - title: "Story 1"
      tasks:
        - "Story 1 task 1"
        - "Story 1 task 2"
    - title: "Story 2"
      tasks:
        - "Story 2 task 1"
`
	tmpl, err := Parse([]byte(yaml))
	require.NoError(t, err)
	// 2 epic tasks + 2 story1 tasks + 1 story2 task = 5
	assert.Equal(t, 5, tmpl.TotalTaskCount())
}

func TestTemplate_ComplexHierarchy(t *testing.T) {
	yaml := `
name: feature-complete
structure:
  epic:
    title: "Feature: {{.Name}}"
    description: |
      As a {{.Actor}}
      I want {{.Goal}}
      So that {{.Benefit}}
    tasks:
      - "Research and design"
      - "Create technical spec"
  stories:
    - title: "Frontend: {{.Name}}"
      description: "Implement {{.Name}} in the UI"
      tasks:
        - "Component implementation"
        - "Unit tests"
        - "Integration tests"
    - title: "Backend: {{.Name}}"
      description: "Implement {{.Name}} API"
      tasks:
        - "API endpoint"
        - "Database schema"
        - "Unit tests"
    - title: "Documentation: {{.Name}}"
      description: "Document {{.Name}} feature"
      tasks:
        - "API documentation"
        - "User guide"
`
	tmpl, err := Parse([]byte(yaml))
	require.NoError(t, err)

	assert.Equal(t, "feature-complete", tmpl.Name)
	assert.True(t, tmpl.HasEpic())
	assert.True(t, tmpl.HasStories())
	assert.Equal(t, 3, tmpl.StoryCount())
	assert.Equal(t, 10, tmpl.TotalTaskCount()) // 2 + 3 + 3 + 2

	// Verify epic
	require.NotNil(t, tmpl.Structure.Epic)
	assert.Equal(t, "Feature: {{.Name}}", tmpl.Structure.Epic.Title)
	assert.Contains(t, tmpl.Structure.Epic.Description, "{{.Actor}}")
	assert.Len(t, tmpl.Structure.Epic.Tasks, 2)

	// Verify stories
	assert.Equal(t, "Frontend: {{.Name}}", tmpl.Structure.Stories[0].Title)
	assert.Len(t, tmpl.Structure.Stories[0].Tasks, 3)
	assert.Equal(t, "Backend: {{.Name}}", tmpl.Structure.Stories[1].Title)
	assert.Len(t, tmpl.Structure.Stories[1].Tasks, 3)
	assert.Equal(t, "Documentation: {{.Name}}", tmpl.Structure.Stories[2].Title)
	assert.Len(t, tmpl.Structure.Stories[2].Tasks, 2)
}
