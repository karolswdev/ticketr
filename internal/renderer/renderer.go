package renderer

import (
	"fmt"
	"strings"

	"github.com/karolswdev/ticketr/internal/core/domain"
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

    writeTitle(&sb, ticket)
    appendCustomFields(&sb, ticket.CustomFields)
    appendDescription(&sb, ticket.Description)
    appendAcceptance(&sb, ticket.AcceptanceCriteria)
    appendTasks(&sb, ticket.Tasks)

    return sb.String()
}

func writeTitle(sb *strings.Builder, t domain.Ticket) {
    if t.JiraID != "" {
        sb.WriteString(fmt.Sprintf("# TICKET: [%s] %s\n\n", t.JiraID, t.Title))
        return
    }
    sb.WriteString(fmt.Sprintf("# TICKET: %s\n\n", t.Title))
}

func appendCustomFields(sb *strings.Builder, fields map[string]string) {
    has := false
    for name, val := range fields {
        if name != "Type" && name != "Parent" && val != "" {
            if !has { sb.WriteString("## Fields\n"); has = true }
            sb.WriteString(fmt.Sprintf("- %s: %s\n", name, val))
        }
    }
    if has { sb.WriteString("\n") }
}

func appendDescription(sb *strings.Builder, desc string) {
    if desc == "" { return }
    sb.WriteString("## Description\n")
    sb.WriteString(desc)
    sb.WriteString("\n\n")
}

func appendAcceptance(sb *strings.Builder, ac []string) {
    if len(ac) == 0 { return }
    sb.WriteString("## Acceptance Criteria\n")
    for _, c := range ac {
        sb.WriteString(fmt.Sprintf("- %s\n", c))
    }
    sb.WriteString("\n")
}

func appendTasks(sb *strings.Builder, tasks []domain.Task) {
    if len(tasks) == 0 { return }
    sb.WriteString("## Tasks\n")
    for _, t := range tasks {
        if t.JiraID != "" {
            sb.WriteString(fmt.Sprintf("- [%s] %s\n", t.JiraID, t.Title))
        } else {
            sb.WriteString(fmt.Sprintf("- %s\n", t.Title))
        }
        for name, val := range t.CustomFields {
            if val != "" {
                sb.WriteString(fmt.Sprintf("  - %s: %s\n", name, val))
            }
        }
        if t.Description != "" {
            for _, line := range strings.Split(t.Description, "\n") {
                if line != "" { sb.WriteString(fmt.Sprintf("  %s\n", line)) }
            }
        }
        if len(t.AcceptanceCriteria) > 0 {
            sb.WriteString("  ### Acceptance Criteria\n")
            for _, c := range t.AcceptanceCriteria {
                sb.WriteString(fmt.Sprintf("  - %s\n", c))
            }
        }
    }
    sb.WriteString("\n")
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
