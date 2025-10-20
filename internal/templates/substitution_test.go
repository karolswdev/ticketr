package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want []string
	}{
		{
			name: "single variable in title",
			yaml: `
name: test
structure:
  epic:
    title: "Feature: {{.Name}}"
`,
			want: []string{"Name"},
		},
		{
			name: "multiple variables in description",
			yaml: `
name: test
structure:
  epic:
    title: "Feature"
    description: |
      As a {{.Actor}}
      I want {{.Goal}}
      So that {{.Benefit}}
`,
			want: []string{"Actor", "Benefit", "Goal"}, // Sorted alphabetically
		},
		{
			name: "variables in stories",
			yaml: `
name: test
structure:
  stories:
    - title: "Frontend: {{.Name}}"
      description: "Implement {{.Feature}}"
    - title: "Backend: {{.Name}}"
      description: "API for {{.Feature}}"
`,
			want: []string{"Feature", "Name"},
		},
		{
			name: "duplicate variables",
			yaml: `
name: test
structure:
  epic:
    title: "{{.Name}}"
  stories:
    - title: "{{.Name}} - Part 1"
    - title: "{{.Name}} - Part 2"
`,
			want: []string{"Name"}, // Should only appear once
		},
		{
			name: "variables in tasks",
			yaml: `
name: test
structure:
  stories:
    - title: "Story"
      tasks:
        - "Implement {{.Component}}"
        - "Test {{.Component}}"
`,
			want: []string{"Component"},
		},
		{
			name: "no variables",
			yaml: `
name: test
structure:
  epic:
    title: "Static Title"
    description: "Static Description"
`,
			want: []string{},
		},
		{
			name: "mixed variables across all fields",
			yaml: `
name: test
structure:
  epic:
    title: "Epic: {{.Name}}"
    description: "For {{.Actor}}"
    tasks:
      - "Research {{.Topic}}"
  stories:
    - title: "Story: {{.Name}}"
      description: "Sprint {{.Sprint}}"
      tasks:
        - "Implement {{.Feature}}"
`,
			want: []string{"Actor", "Feature", "Name", "Sprint", "Topic"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := Parse([]byte(tt.yaml))
			require.NoError(t, err)

			got := ExtractVariables(tmpl)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractVariables_InvalidSyntax(t *testing.T) {
	// Variables without dot should not be extracted
	yaml := `
name: test
structure:
  epic:
    title: "{{Name}} vs {{.Name}}"
`
	tmpl, err := Parse([]byte(yaml))
	require.NoError(t, err)

	got := ExtractVariables(tmpl)
	// Only {{.Name}} should be extracted, not {{Name}}
	assert.Equal(t, []string{"Name"}, got)
}

func TestSubstitute_Success(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		vars      map[string]string
		checkFunc func(*testing.T, *Template)
	}{
		{
			name: "simple substitution",
			yaml: `
name: test
structure:
  epic:
    title: "Feature: {{.Name}}"
`,
			vars: map[string]string{"Name": "Authentication"},
			checkFunc: func(t *testing.T, tmpl *Template) {
				assert.Equal(t, "Feature: Authentication", tmpl.Structure.Epic.Title)
			},
		},
		{
			name: "multiple variables",
			yaml: `
name: test
structure:
  epic:
    title: "Feature: {{.Name}}"
    description: |
      As a {{.Actor}}
      I want {{.Goal}}
      So that {{.Benefit}}
`,
			vars: map[string]string{
				"Name":    "Login",
				"Actor":   "user",
				"Goal":    "to access my account",
				"Benefit": "I can view my data",
			},
			checkFunc: func(t *testing.T, tmpl *Template) {
				assert.Equal(t, "Feature: Login", tmpl.Structure.Epic.Title)
				assert.Contains(t, tmpl.Structure.Epic.Description, "As a user")
				assert.Contains(t, tmpl.Structure.Epic.Description, "I want to access my account")
				assert.Contains(t, tmpl.Structure.Epic.Description, "So that I can view my data")
			},
		},
		{
			name: "substitution in stories",
			yaml: `
name: test
structure:
  stories:
    - title: "Frontend: {{.Name}}"
      description: "Implement {{.Component}}"
    - title: "Backend: {{.Name}}"
      description: "API for {{.Component}}"
`,
			vars: map[string]string{
				"Name":      "User Profile",
				"Component": "profile page",
			},
			checkFunc: func(t *testing.T, tmpl *Template) {
				assert.Equal(t, "Frontend: User Profile", tmpl.Structure.Stories[0].Title)
				assert.Equal(t, "Implement profile page", tmpl.Structure.Stories[0].Description)
				assert.Equal(t, "Backend: User Profile", tmpl.Structure.Stories[1].Title)
				assert.Equal(t, "API for profile page", tmpl.Structure.Stories[1].Description)
			},
		},
		{
			name: "substitution in tasks",
			yaml: `
name: test
structure:
  stories:
    - title: "Story"
      tasks:
        - "Implement {{.Component}}"
        - "Test {{.Component}}"
        - "Document {{.Component}}"
`,
			vars: map[string]string{
				"Component": "authentication",
			},
			checkFunc: func(t *testing.T, tmpl *Template) {
				assert.Equal(t, "Implement authentication", tmpl.Structure.Stories[0].Tasks[0])
				assert.Equal(t, "Test authentication", tmpl.Structure.Stories[0].Tasks[1])
				assert.Equal(t, "Document authentication", tmpl.Structure.Stories[0].Tasks[2])
			},
		},
		{
			name: "no variables to substitute",
			yaml: `
name: test
structure:
  epic:
    title: "Static Title"
    description: "Static Description"
`,
			vars: map[string]string{},
			checkFunc: func(t *testing.T, tmpl *Template) {
				assert.Equal(t, "Static Title", tmpl.Structure.Epic.Title)
				assert.Equal(t, "Static Description", tmpl.Structure.Epic.Description)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original, err := Parse([]byte(tt.yaml))
			require.NoError(t, err)

			result, err := Substitute(original, tt.vars)
			require.NoError(t, err)
			require.NotNil(t, result)

			tt.checkFunc(t, result)

			// Verify original template is unchanged
			assert.NotEqual(t, original, result, "Substitute should create a new template")
		})
	}
}

func TestSubstitute_MissingVariable(t *testing.T) {
	yaml := `
name: test
structure:
  epic:
    title: "Feature: {{.Name}}"
    description: "For {{.Actor}}"
`
	tmpl, err := Parse([]byte(yaml))
	require.NoError(t, err)

	// Only provide one of two required variables
	vars := map[string]string{
		"Name": "Authentication",
	}

	_, err = Substitute(tmpl, vars)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing required variable: Actor")
}

func TestSubstitute_AllMissingVariables(t *testing.T) {
	yaml := `
name: test
structure:
  epic:
    title: "{{.Name}}"
`
	tmpl, err := Parse([]byte(yaml))
	require.NoError(t, err)

	_, err = Substitute(tmpl, map[string]string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing required variable")
}

func TestValidateVariables_Success(t *testing.T) {
	yaml := `
name: test
structure:
  epic:
    title: "{{.Name}}"
    description: "{{.Description}}"
`
	tmpl, err := Parse([]byte(yaml))
	require.NoError(t, err)

	vars := map[string]string{
		"Name":        "Feature",
		"Description": "Feature description",
	}

	err = ValidateVariables(tmpl, vars)
	assert.NoError(t, err)
}

func TestValidateVariables_MissingVariable(t *testing.T) {
	yaml := `
name: test
structure:
  epic:
    title: "{{.Name}}"
    description: "{{.Description}}"
`
	tmpl, err := Parse([]byte(yaml))
	require.NoError(t, err)

	vars := map[string]string{
		"Name": "Feature",
		// Missing Description
	}

	err = ValidateVariables(tmpl, vars)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing required variables")
	assert.Contains(t, err.Error(), "Description")
}

func TestValidateVariables_EmptyVariable(t *testing.T) {
	yaml := `
name: test
structure:
  epic:
    title: "{{.Name}}"
`
	tmpl, err := Parse([]byte(yaml))
	require.NoError(t, err)

	vars := map[string]string{
		"Name": "", // Empty value
	}

	err = ValidateVariables(tmpl, vars)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing required variables")
}

func TestSubstituteString(t *testing.T) {
	tests := []struct {
		name string
		str  string
		vars map[string]string
		want string
	}{
		{
			name: "simple substitution",
			str:  "Hello {{.Name}}",
			vars: map[string]string{"Name": "World"},
			want: "Hello World",
		},
		{
			name: "multiple variables",
			str:  "{{.First}} {{.Last}}",
			vars: map[string]string{"First": "John", "Last": "Doe"},
			want: "John Doe",
		},
		{
			name: "no variables",
			str:  "Static string",
			vars: map[string]string{},
			want: "Static string",
		},
		{
			name: "empty string",
			str:  "",
			vars: map[string]string{"Name": "Value"},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := substituteString(tt.str, tt.vars)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSubstitute_DeepCopy(t *testing.T) {
	yaml := `
name: test
structure:
  epic:
    title: "{{.Name}}"
    tasks:
      - "Task 1"
  stories:
    - title: "Story 1"
      tasks:
        - "Story task 1"
`
	original, err := Parse([]byte(yaml))
	require.NoError(t, err)

	vars := map[string]string{"Name": "Feature"}
	result, err := Substitute(original, vars)
	require.NoError(t, err)

	// Verify original is unchanged
	assert.Equal(t, "{{.Name}}", original.Structure.Epic.Title)
	assert.Equal(t, "Feature", result.Structure.Epic.Title)

	// Verify deep copy by modifying result
	result.Structure.Epic.Tasks[0] = "Modified task"
	assert.Equal(t, "Task 1", original.Structure.Epic.Tasks[0])
	assert.Equal(t, "Modified task", result.Structure.Epic.Tasks[0])
}
