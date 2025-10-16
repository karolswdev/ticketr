package renderer

import (
	"fmt"
	"strings"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// Renderer handles conversion of tickets to Markdown format
type Renderer struct {
	fieldMappings map[string]interface{}
}

// NewRenderer creates a new Renderer instance
func NewRenderer(fieldMappings map[string]interface{}) *Renderer {
	if fieldMappings == nil {
		fieldMappings = getDefaultFieldMappings()
	}
	return &Renderer{
		fieldMappings: fieldMappings,
	}
}

// getDefaultFieldMappings returns default field mappings
func getDefaultFieldMappings() map[string]interface{} {
	return map[string]interface{}{
		"Type":        "issuetype",
		"Project":     "project",
		"Summary":     "summary",
		"Description": "description",
		"Assignee":    "assignee",
		"Reporter":    "reporter",
		"Priority":    "priority",
		"Labels":      "labels",
		"Components":  "components",
		"Fix Version": "fixVersions",
		"Sprint":      "customfield_10020",
		"Story Points": map[string]interface{}{
			"id":   "customfield_10010",
			"type": "number",
		},
	}
}

// Render converts a domain.Ticket to Markdown format
func (r *Renderer) Render(ticket domain.Ticket) string {
	var sb strings.Builder

	// Title with JIRA ID if present
	if ticket.JiraID != "" {
		sb.WriteString(fmt.Sprintf("# TICKET: [%s] %s\n", ticket.JiraID, ticket.Title))
	} else {
		sb.WriteString(fmt.Sprintf("# TICKET: %s\n", ticket.Title))
	}
	sb.WriteString("\n")

	// Custom fields section (excluding Type which is handled differently in some cases)
	hasCustomFields := false
	for fieldName, fieldValue := range ticket.CustomFields {
		if fieldName != "Type" && fieldName != "Parent" && fieldValue != "" {
			if !hasCustomFields {
				sb.WriteString("## Fields\n")
				hasCustomFields = true
			}
			sb.WriteString(fmt.Sprintf("- %s: %s\n", fieldName, fieldValue))
		}
	}
	if hasCustomFields {
		sb.WriteString("\n")
	}

	// Description section
	if ticket.Description != "" {
		sb.WriteString("## Description\n")
		sb.WriteString(ticket.Description)
		sb.WriteString("\n\n")
	}

	// Acceptance Criteria section
	if len(ticket.AcceptanceCriteria) > 0 {
		sb.WriteString("## Acceptance Criteria\n")
		for _, criterion := range ticket.AcceptanceCriteria {
			sb.WriteString(fmt.Sprintf("- %s\n", criterion))
		}
		sb.WriteString("\n")
	}

	// Tasks section
	if len(ticket.Tasks) > 0 {
		sb.WriteString("## Tasks\n")
		for _, task := range ticket.Tasks {
			if task.JiraID != "" {
				sb.WriteString(fmt.Sprintf("- [%s] %s\n", task.JiraID, task.Title))
			} else {
				sb.WriteString(fmt.Sprintf("- %s\n", task.Title))
			}

			// Task custom fields (indented)
			for fieldName, fieldValue := range task.CustomFields {
				if fieldValue != "" {
					sb.WriteString(fmt.Sprintf("  - %s: %s\n", fieldName, fieldValue))
				}
			}

			// Task description (indented)
			if task.Description != "" {
				lines := strings.Split(task.Description, "\n")
				for _, line := range lines {
					if line != "" {
						sb.WriteString(fmt.Sprintf("  %s\n", line))
					}
				}
			}

			// Task acceptance criteria (indented)
			if len(task.AcceptanceCriteria) > 0 {
				sb.WriteString("  ### Acceptance Criteria\n")
				for _, criterion := range task.AcceptanceCriteria {
					sb.WriteString(fmt.Sprintf("  - %s\n", criterion))
				}
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderMultiple renders multiple tickets to a single Markdown document
func (r *Renderer) RenderMultiple(tickets []domain.Ticket) string {
	var sb strings.Builder

	for i, ticket := range tickets {
		sb.WriteString(r.Render(ticket))

		// Add separator between tickets except for the last one
		if i < len(tickets)-1 {
			sb.WriteString("---\n\n")
		}
	}

	return sb.String()
}
