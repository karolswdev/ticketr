package templates

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"text/template"
)

var variableRegex = regexp.MustCompile(`\{\{\.(\w+)\}\}`)

// ExtractVariables extracts all unique variables from a template
func ExtractVariables(tmpl *Template) []string {
	vars := make(map[string]bool)

	// Extract from epic
	if tmpl.Structure.Epic != nil {
		extractFromString(tmpl.Structure.Epic.Title, vars)
		extractFromString(tmpl.Structure.Epic.Description, vars)
		for _, task := range tmpl.Structure.Epic.Tasks {
			extractFromString(task, vars)
		}
	}

	// Extract from stories
	for _, story := range tmpl.Structure.Stories {
		extractFromString(story.Title, vars)
		extractFromString(story.Description, vars)
		for _, task := range story.Tasks {
			extractFromString(task, vars)
		}
	}

	// Convert to sorted slice for deterministic output
	result := make([]string, 0, len(vars))
	for v := range vars {
		result = append(result, v)
	}
	sort.Strings(result)

	return result
}

// extractFromString extracts variables from a single string
func extractFromString(s string, vars map[string]bool) {
	matches := variableRegex.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		if len(match) > 1 {
			vars[match[1]] = true
		}
	}
}

// Substitute replaces variables in a template with provided values
func Substitute(tmpl *Template, vars map[string]string) (*Template, error) {
	// Validate all required variables are provided
	requiredVars := ExtractVariables(tmpl)
	for _, v := range requiredVars {
		if _, ok := vars[v]; !ok {
			return nil, fmt.Errorf("missing required variable: %s", v)
		}
	}

	// Create a deep copy of the template
	result := &Template{
		Name: tmpl.Name,
		Structure: TemplateStructure{
			Stories: make([]TicketTemplate, len(tmpl.Structure.Stories)),
		},
	}

	// Substitute in epic
	if tmpl.Structure.Epic != nil {
		epic := &TicketTemplate{
			Title:       substituteString(tmpl.Structure.Epic.Title, vars),
			Description: substituteString(tmpl.Structure.Epic.Description, vars),
			Tasks:       make([]string, len(tmpl.Structure.Epic.Tasks)),
		}
		for i, task := range tmpl.Structure.Epic.Tasks {
			epic.Tasks[i] = substituteString(task, vars)
		}
		result.Structure.Epic = epic
	}

	// Substitute in stories
	for i, story := range tmpl.Structure.Stories {
		newStory := TicketTemplate{
			Title:       substituteString(story.Title, vars),
			Description: substituteString(story.Description, vars),
			Tasks:       make([]string, len(story.Tasks)),
		}
		for j, task := range story.Tasks {
			newStory.Tasks[j] = substituteString(task, vars)
		}
		result.Structure.Stories[i] = newStory
	}

	return result, nil
}

// substituteString replaces variables in a single string
func substituteString(s string, vars map[string]string) string {
	if s == "" {
		return s
	}

	tmpl, err := template.New("").Parse(s)
	if err != nil {
		// Return original string if parsing fails
		return s
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		// Return original string if execution fails
		return s
	}

	return buf.String()
}

// ValidateVariables checks that all required variables are provided
func ValidateVariables(tmpl *Template, vars map[string]string) error {
	requiredVars := ExtractVariables(tmpl)
	missing := []string{}

	for _, v := range requiredVars {
		if val, ok := vars[v]; !ok || val == "" {
			missing = append(missing, v)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required variables: %v", missing)
	}

	return nil
}
