// Package jira provides an example implementation of JiraPort using andygrunwald/go-jira
// This demonstrates how simple the migration would be
package jira

import (
	"context"
	"fmt"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// JiraAdapterV2 implements JiraPort using andygrunwald/go-jira library
type JiraAdapterV2 struct {
	client        *jira.Client
	projectKey    string
	storyType     string
	subTaskType   string
	fieldMappings map[string]string // Human name → Jira field ID
}

// NewJiraAdapterV2FromConfig creates a new adapter from workspace configuration
func NewJiraAdapterV2FromConfig(config *domain.WorkspaceConfig, fieldMappings map[string]string) (ports.JiraPort, error) {
	if config == nil {
		return nil, fmt.Errorf("workspace configuration is required")
	}

	// Validate configuration
	if config.JiraURL == "" || config.Username == "" || config.APIToken == "" || config.ProjectKey == "" {
		return nil, fmt.Errorf("missing required configuration fields")
	}

	// Create authenticated transport
	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.APIToken,
	}

	// Create Jira client
	client, err := jira.NewClient(tp.Client(), strings.TrimRight(config.JiraURL, "/"))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	// Default field mappings if not provided
	if fieldMappings == nil {
		fieldMappings = map[string]string{
			"Story Points": "customfield_10010",
			"Sprint":       "customfield_10020",
		}
	}

	return &JiraAdapterV2{
		client:        client,
		projectKey:    config.ProjectKey,
		storyType:     "Task", // Could come from config
		subTaskType:   "Sub-task",
		fieldMappings: fieldMappings,
	}, nil
}

// Authenticate verifies the connection to Jira
func (j *JiraAdapterV2) Authenticate() error {
	// Use the /myself endpoint to verify authentication
	_, resp, err := j.client.User.GetSelf()
	if err != nil {
		return fmt.Errorf("authentication failed (status %d): %w", resp.StatusCode, err)
	}
	return nil
}

// SearchTickets searches for tickets using JQL with progress callbacks
func (j *JiraAdapterV2) SearchTickets(ctx context.Context, projectKey string, jql string, progressCallback ports.JiraProgressCallback) ([]domain.Ticket, error) {
	// Construct full JQL
	fullJQL := fmt.Sprintf("project = %s", projectKey)
	if jql != "" {
		fullJQL = fmt.Sprintf("%s AND %s", fullJQL, jql)
	}

	// Report initial connection
	if progressCallback != nil {
		progressCallback(0, 0, "Connecting to Jira...")
	}

	// Search with library (handles pagination automatically if needed)
	searchOptions := &jira.SearchOptions{
		MaxResults: 100, // Can paginate if more needed
		StartAt:    0,
		Fields:     []string{"summary", "description", "issuetype", "status", "parent", "subtasks"},
	}

	// TODO: Add context support via custom http.Client if needed
	issues, resp, err := j.client.Issue.Search(fullJQL, searchOptions)
	if err != nil {
		return nil, fmt.Errorf("search failed (status %d): %w", resp.StatusCode, err)
	}

	// Convert to domain tickets with progress reporting
	tickets := make([]domain.Ticket, len(issues))
	for i, issue := range issues {
		tickets[i] = j.convertToDomainTicket(&issue)

		if progressCallback != nil {
			progressCallback(i+1, len(issues), fmt.Sprintf("Processing %d/%d tickets", i+1, len(issues)))
		}
	}

	return tickets, nil
}

// CreateTicket creates a new ticket in Jira
func (j *JiraAdapterV2) CreateTicket(ticket domain.Ticket) (string, error) {
	// Build Jira issue from domain ticket
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Project: jira.Project{
				Key: j.projectKey,
			},
			Type: jira.IssueType{
				Name: j.storyType,
			},
			Summary:     ticket.Title,
			Description: j.buildDescription(ticket.Description, ticket.AcceptanceCriteria),
			Unknowns:    make(map[string]interface{}),
		},
	}

	// Map custom fields
	for fieldName, fieldValue := range ticket.CustomFields {
		if jiraFieldID, exists := j.fieldMappings[fieldName]; exists {
			issue.Fields.Unknowns[jiraFieldID] = j.convertFieldValue(fieldName, fieldValue)
		}
	}

	// Create issue
	createdIssue, resp, err := j.client.Issue.Create(issue)
	if err != nil {
		return "", fmt.Errorf("create failed (status %d): %w", resp.StatusCode, err)
	}

	return createdIssue.Key, nil
}

// UpdateTicket updates an existing ticket
func (j *JiraAdapterV2) UpdateTicket(ticket domain.Ticket) error {
	if ticket.JiraID == "" {
		return fmt.Errorf("ticket does not have a Jira ID")
	}

	// Build update fields
	fields := map[string]interface{}{
		"summary":     ticket.Title,
		"description": j.buildDescription(ticket.Description, ticket.AcceptanceCriteria),
	}

	// Add custom fields
	for fieldName, fieldValue := range ticket.CustomFields {
		if jiraFieldID, exists := j.fieldMappings[fieldName]; exists {
			fields[jiraFieldID] = j.convertFieldValue(fieldName, fieldValue)
		}
	}

	// Update issue
	resp, err := j.client.Issue.UpdateIssue(ticket.JiraID, fields)
	if err != nil {
		return fmt.Errorf("update failed (status %d): %w", resp.StatusCode, err)
	}

	return nil
}

// CreateTask creates a new subtask under a parent issue
func (j *JiraAdapterV2) CreateTask(task domain.Task, parentID string) (string, error) {
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Project: jira.Project{
				Key: j.projectKey,
			},
			Type: jira.IssueType{
				Name: j.subTaskType,
			},
			Parent: &jira.Parent{
				Key: parentID,
			},
			Summary:     task.Title,
			Description: j.buildDescription(task.Description, task.AcceptanceCriteria),
			Unknowns:    make(map[string]interface{}),
		},
	}

	// Map custom fields
	for fieldName, fieldValue := range task.CustomFields {
		if jiraFieldID, exists := j.fieldMappings[fieldName]; exists {
			issue.Fields.Unknowns[jiraFieldID] = j.convertFieldValue(fieldName, fieldValue)
		}
	}

	createdIssue, resp, err := j.client.Issue.Create(issue)
	if err != nil {
		return "", fmt.Errorf("create subtask failed (status %d): %w", resp.StatusCode, err)
	}

	return createdIssue.Key, nil
}

// UpdateTask updates an existing subtask
func (j *JiraAdapterV2) UpdateTask(task domain.Task) error {
	if task.JiraID == "" {
		return fmt.Errorf("task does not have a Jira ID")
	}

	fields := map[string]interface{}{
		"summary":     task.Title,
		"description": j.buildDescription(task.Description, task.AcceptanceCriteria),
	}

	for fieldName, fieldValue := range task.CustomFields {
		if jiraFieldID, exists := j.fieldMappings[fieldName]; exists {
			fields[jiraFieldID] = j.convertFieldValue(fieldName, fieldValue)
		}
	}

	resp, err := j.client.Issue.UpdateIssue(task.JiraID, fields)
	if err != nil {
		return fmt.Errorf("update subtask failed (status %d): %w", resp.StatusCode, err)
	}

	return nil
}

// GetProjectIssueTypes fetches available issue types for the configured project
func (j *JiraAdapterV2) GetProjectIssueTypes() (map[string][]string, error) {
	project, resp, err := j.client.Project.Get(j.projectKey)
	if err != nil {
		return nil, fmt.Errorf("get project failed (status %d): %w", resp.StatusCode, err)
	}

	result := map[string][]string{
		"project": {project.Name},
		"key":     {project.Key},
	}

	issueTypes := make([]string, 0, len(project.IssueTypes))
	for _, issueType := range project.IssueTypes {
		if issueType.Subtask {
			issueTypes = append(issueTypes, fmt.Sprintf("%s (subtask)", issueType.Name))
		} else {
			issueTypes = append(issueTypes, issueType.Name)
		}
	}
	result["issueTypes"] = issueTypes

	return result, nil
}

// GetIssueTypeFields fetches field requirements for a specific issue type
func (j *JiraAdapterV2) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	// Use the create metadata endpoint
	// Library provides metaProject, metaIssueType structures
	// This is a placeholder - actual implementation would use client.Issue.GetCreateMeta()
	return nil, fmt.Errorf("not yet implemented - use client.Issue.GetCreateMeta()")
}

// Helper functions

func (j *JiraAdapterV2) buildDescription(description string, acceptanceCriteria []string) string {
	if len(acceptanceCriteria) == 0 {
		return description
	}

	desc := description + "\n\nh3. Acceptance Criteria\n"
	for _, ac := range acceptanceCriteria {
		desc += fmt.Sprintf("* %s\n", ac)
	}
	return desc
}

func (j *JiraAdapterV2) convertFieldValue(fieldName, value string) interface{} {
	// Handle known numeric fields
	if fieldName == "Story Points" {
		var num float64
		if _, err := fmt.Sscanf(value, "%f", &num); err == nil {
			return num
		}
	}

	// Default to string
	return value
}

func (j *JiraAdapterV2) convertToDomainTicket(issue *jira.Issue) domain.Ticket {
	ticket := domain.Ticket{
		JiraID:       issue.Key,
		Title:        issue.Fields.Summary,
		CustomFields: make(map[string]string),
	}

	// Parse description and acceptance criteria
	if issue.Fields.Description != "" {
		parts := strings.Split(issue.Fields.Description, "h3. Acceptance Criteria")
		ticket.Description = strings.TrimSpace(parts[0])

		if len(parts) > 1 {
			acLines := strings.Split(parts[1], "\n")
			for _, line := range acLines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "* ") {
					ticket.AcceptanceCriteria = append(ticket.AcceptanceCriteria, strings.TrimPrefix(line, "* "))
				}
			}
		}
	}

	// Map custom fields back to human-readable names
	for humanName, jiraFieldID := range j.fieldMappings {
		if value, exists := issue.Fields.Unknowns[jiraFieldID]; exists {
			switch v := value.(type) {
			case string:
				ticket.CustomFields[humanName] = v
			case float64:
				ticket.CustomFields[humanName] = fmt.Sprintf("%g", v)
			}
		}
	}

	// Convert subtasks
	// Note: In v2, we would fetch subtasks separately as the library
	// doesn't automatically include full subtask details
	// For now, we'll leave tasks empty and fetch them on-demand
	ticket.Tasks = []domain.Task{}

	return ticket
}

func (j *JiraAdapterV2) convertToDomainTask(issue *jira.Issue) domain.Task {
	task := domain.Task{
		JiraID:       issue.Key,
		Title:        issue.Fields.Summary,
		CustomFields: make(map[string]string),
	}

	// Parse description and acceptance criteria
	if issue.Fields.Description != "" {
		parts := strings.Split(issue.Fields.Description, "h3. Acceptance Criteria")
		task.Description = strings.TrimSpace(parts[0])

		if len(parts) > 1 {
			acLines := strings.Split(parts[1], "\n")
			for _, line := range acLines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "* ") {
					task.AcceptanceCriteria = append(task.AcceptanceCriteria, strings.TrimPrefix(line, "* "))
				}
			}
		}
	}

	return task
}
