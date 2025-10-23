package jira

import (
	"context"
	"fmt"
	"os"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// JiraAdapter implements JiraPort using andygrunwald/go-jira library
type JiraAdapter struct {
	client        *jira.Client
	projectKey    string
	storyType     string
	subTaskType   string
	fieldMappings map[string]interface{} // Maps human-readable names to Jira field IDs
}

// NewJiraAdapterFromConfig creates a new Jira adapter from workspace configuration
func NewJiraAdapterFromConfig(config *domain.WorkspaceConfig, fieldMappings map[string]interface{}) (ports.JiraPort, error) {
	if config == nil {
		return nil, fmt.Errorf("[jira] workspace configuration is required")
	}

	// Validate configuration
	if config.JiraURL == "" {
		return nil, fmt.Errorf("[jira] Jira URL is required in workspace configuration")
	}
	if config.Username == "" {
		return nil, fmt.Errorf("[jira] username is required in workspace configuration")
	}
	if config.APIToken == "" {
		return nil, fmt.Errorf("[jira] API token is required in workspace configuration")
	}
	if config.ProjectKey == "" {
		return nil, fmt.Errorf("[jira] project key is required in workspace configuration")
	}

	// Create authenticated transport
	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.APIToken,
	}

	// Create Jira client
	client, err := jira.NewClient(tp.Client(), strings.TrimRight(config.JiraURL, "/"))
	if err != nil {
		return nil, fmt.Errorf("[jira] failed to create Jira client: %w", err)
	}

	// Get issue types from environment with sensible defaults
	storyType := os.Getenv("JIRA_STORY_TYPE")
	if storyType == "" {
		storyType = "Task" // Default to Task which is more common
	}

	subTaskType := os.Getenv("JIRA_SUBTASK_TYPE")
	if subTaskType == "" {
		subTaskType = "Sub-task" // Standard JIRA subtask type
	}

	// If no field mappings provided, use defaults
	if fieldMappings == nil {
		fieldMappings = getDefaultFieldMappings()
	}

	return &JiraAdapter{
		client:        client,
		projectKey:    config.ProjectKey,
		storyType:     storyType,
		subTaskType:   subTaskType,
		fieldMappings: fieldMappings,
	}, nil
}

// Authenticate verifies the connection to Jira
func (j *JiraAdapter) Authenticate() error {
	// Use the /myself endpoint to verify authentication
	_, resp, err := j.client.User.GetSelf()
	if err != nil {
		if resp != nil {
			return fmt.Errorf("[jira] authentication failed with status %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("[jira] authentication failed: %w", err)
	}
	return nil
}

// SearchTickets searches for tickets using JQL with context support and progress callbacks
func (j *JiraAdapter) SearchTickets(ctx context.Context, projectKey string, jql string, progressCallback ports.JiraProgressCallback) ([]domain.Ticket, error) {
	// Construct full JQL
	fullJQL := fmt.Sprintf(`project = "%s"`, projectKey)
	if jql != "" {
		fullJQL = fmt.Sprintf(`%s AND %s`, fullJQL, jql)
	}

	// Report initial connection
	if progressCallback != nil {
		progressCallback(0, 0, "Connecting to Jira...")
	}

	// Build fields list based on field mappings
	fields := []string{"key", "summary", "description", "issuetype", "parent", "subtasks"}
	for _, mapping := range j.fieldMappings {
		switch m := mapping.(type) {
		case string:
			if m != "summary" && m != "description" && m != "issuetype" && m != "project" {
				fields = append(fields, m)
			}
		case map[string]interface{}:
			if id, ok := m["id"].(string); ok {
				fields = append(fields, id)
			}
		}
	}

	// Use pagination to fetch tickets in batches
	const pageSize = 50
	allTickets := make([]domain.Ticket, 0)
	startAt := 0
	total := -1

	// Paginate through all results
	for {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Search with library
		searchOptions := &jira.SearchOptions{
			MaxResults: pageSize,
			StartAt:    startAt,
			Fields:     fields,
		}

		// Execute search (library doesn't support context directly, but we check above)
		issues, resp, err := j.client.Issue.Search(fullJQL, searchOptions)
		if err != nil {
			if resp != nil {
				return nil, fmt.Errorf("[jira] search failed with status %d: %w", resp.StatusCode, err)
			}
			return nil, fmt.Errorf("[jira] search failed: %w", err)
		}

		// Get total from response (library populates this)
		if total == -1 {
			total = resp.Total
		}

		// Convert issues to domain tickets
		for _, issue := range issues {
			ticket := j.convertToDomainTicket(&issue)
			allTickets = append(allTickets, ticket)
		}

		// Report progress
		if progressCallback != nil {
			currentCount := len(allTickets)
			if total > 0 {
				progressCallback(currentCount, total, fmt.Sprintf("Fetched %d/%d tickets", currentCount, total))
			} else {
				progressCallback(currentCount, currentCount, fmt.Sprintf("Fetched %d tickets", currentCount))
			}
		}

		// Check if done
		if len(issues) == 0 || len(allTickets) >= total {
			break
		}

		startAt += pageSize
	}

	// Report subtask fetching phase
	if progressCallback != nil && len(allTickets) > 0 {
		progressCallback(0, len(allTickets), "Fetching subtasks...")
	}

	// Fetch subtasks for each parent ticket
	for i := range allTickets {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		subtasks, err := j.fetchSubtasks(ctx, allTickets[i].JiraID)
		if err != nil {
			// Non-fatal error, continue processing
			continue
		}
		allTickets[i].Tasks = subtasks

		// Report progress periodically
		if progressCallback != nil && (i%10 == 0 || i == len(allTickets)-1) {
			progressCallback(i+1, len(allTickets), fmt.Sprintf("Processing subtasks %d/%d", i+1, len(allTickets)))
		}
	}

	return allTickets, nil
}

// fetchSubtasks fetches all subtasks for a given parent issue
func (j *JiraAdapter) fetchSubtasks(ctx context.Context, parentKey string) ([]domain.Task, error) {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Construct JQL for subtasks
	jql := fmt.Sprintf(`parent = "%s"`, parentKey)

	// Build fields list
	fields := []string{"key", "summary", "description", "issuetype"}
	for _, mapping := range j.fieldMappings {
		switch m := mapping.(type) {
		case string:
			if m != "summary" && m != "description" && m != "issuetype" && m != "project" {
				fields = append(fields, m)
			}
		case map[string]interface{}:
			if id, ok := m["id"].(string); ok {
				fields = append(fields, id)
			}
		}
	}

	searchOptions := &jira.SearchOptions{
		MaxResults: 100,
		Fields:     fields,
	}

	issues, resp, err := j.client.Issue.Search(jql, searchOptions)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("[jira] subtask search failed with status %d: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("[jira] subtask search failed: %w", err)
	}

	tasks := make([]domain.Task, 0, len(issues))
	for _, issue := range issues {
		task := j.convertToDomainTask(&issue)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// CreateTicket creates a new ticket in Jira
func (j *JiraAdapter) CreateTicket(ticket domain.Ticket) (string, error) {
	// Build description with acceptance criteria
	description := j.buildDescription(ticket.Description, ticket.AcceptanceCriteria)

	// Build Jira issue
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Project: jira.Project{
				Key: j.projectKey,
			},
			Type: jira.IssueType{
				Name: j.storyType,
			},
			Summary:     ticket.Title,
			Description: description,
			Unknowns:    make(map[string]interface{}),
		},
	}

	// Map custom fields
	j.applyCustomFields(issue.Fields, ticket.CustomFields)

	// Create issue
	createdIssue, resp, err := j.client.Issue.Create(issue)
	if err != nil {
		if resp != nil {
			return "", fmt.Errorf("[jira] create failed with status %d: %w", resp.StatusCode, err)
		}
		return "", fmt.Errorf("[jira] create failed: %w", err)
	}

	return createdIssue.Key, nil
}

// UpdateTicket updates an existing ticket in Jira
func (j *JiraAdapter) UpdateTicket(ticket domain.Ticket) error {
	if ticket.JiraID == "" {
		return fmt.Errorf("[jira] ticket does not have a Jira ID")
	}

	// Build description with acceptance criteria
	description := j.buildDescription(ticket.Description, ticket.AcceptanceCriteria)

	// Build update fields
	fields := map[string]interface{}{
		"summary":     ticket.Title,
		"description": description,
	}

	// Add custom fields
	for fieldName, fieldValue := range ticket.CustomFields {
		if jiraField := j.getJiraFieldID(fieldName); jiraField != "" {
			// Skip fields that shouldn't be updated
			if jiraField == "project" || jiraField == "issuetype" {
				continue
			}
			fields[jiraField] = j.convertFieldValue(fieldName, fieldValue)
		}
	}

	// Get the issue first
	issue, resp, err := j.client.Issue.Get(ticket.JiraID, nil)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("[jira] get issue failed with status %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("[jira] get issue failed: %w", err)
	}

	// Update the fields
	issue.Fields.Summary = ticket.Title
	issue.Fields.Description = description

	// Apply custom fields to Unknowns
	for k, v := range fields {
		if k != "summary" && k != "description" {
			issue.Fields.Unknowns[k] = v
		}
	}

	// Update issue
	_, resp, err = j.client.Issue.Update(issue)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("[jira] update failed with status %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("[jira] update failed: %w", err)
	}

	return nil
}

// CreateTask creates a new subtask under a parent issue
func (j *JiraAdapter) CreateTask(task domain.Task, parentID string) (string, error) {
	// Build description with acceptance criteria
	description := j.buildDescription(task.Description, task.AcceptanceCriteria)

	// Build Jira issue
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
			Description: description,
			Unknowns:    make(map[string]interface{}),
		},
	}

	// Map custom fields
	j.applyCustomFields(issue.Fields, task.CustomFields)

	// Create issue
	createdIssue, resp, err := j.client.Issue.Create(issue)
	if err != nil {
		if resp != nil {
			return "", fmt.Errorf("[jira] create subtask failed with status %d: %w", resp.StatusCode, err)
		}
		return "", fmt.Errorf("[jira] create subtask failed: %w", err)
	}

	return createdIssue.Key, nil
}

// UpdateTask updates an existing subtask
func (j *JiraAdapter) UpdateTask(task domain.Task) error {
	if task.JiraID == "" {
		return fmt.Errorf("[jira] task does not have a Jira ID")
	}

	// Build description with acceptance criteria
	description := j.buildDescription(task.Description, task.AcceptanceCriteria)

	// Get the issue first
	issue, resp, err := j.client.Issue.Get(task.JiraID, nil)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("[jira] get issue failed with status %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("[jira] get issue failed: %w", err)
	}

	// Update the fields
	issue.Fields.Summary = task.Title
	issue.Fields.Description = description

	// Apply custom fields
	for fieldName, fieldValue := range task.CustomFields {
		if jiraField := j.getJiraFieldID(fieldName); jiraField != "" {
			// Skip fields that shouldn't be updated
			if jiraField == "project" || jiraField == "issuetype" || jiraField == "parent" {
				continue
			}
			issue.Fields.Unknowns[jiraField] = j.convertFieldValue(fieldName, fieldValue)
		}
	}

	// Update issue
	_, resp, err = j.client.Issue.Update(issue)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("[jira] update subtask failed with status %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("[jira] update subtask failed: %w", err)
	}

	return nil
}

// GetProjectIssueTypes fetches available issue types for the configured project
func (j *JiraAdapter) GetProjectIssueTypes() (map[string][]string, error) {
	project, resp, err := j.client.Project.Get(j.projectKey)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("[jira] get project failed with status %d: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("[jira] get project failed: %w", err)
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
func (j *JiraAdapter) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	// Use the create metadata endpoint with options to avoid deprecated v2 API
	// The expand parameter tells Jira to include detailed field information
	options := &jira.GetQueryOptions{
		ProjectKeys: j.projectKey,
		Expand:      "projects.issuetypes.fields",
	}

	createMeta, resp, err := j.client.Issue.GetCreateMetaWithOptions(options)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("[jira] get create metadata failed with status %d: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("[jira] get create metadata failed: %w", err)
	}

	// Find our project using the helper method
	targetProject := createMeta.GetProjectWithKey(j.projectKey)
	if targetProject == nil {
		return nil, fmt.Errorf("[jira] project %s not found in metadata", j.projectKey)
	}

	// Find the issue type
	var targetIssueType *jira.MetaIssueType
	for i := range targetProject.IssueTypes {
		if targetProject.IssueTypes[i].Name == issueTypeName {
			targetIssueType = targetProject.IssueTypes[i]
			break
		}
	}

	if targetIssueType == nil {
		availableTypes := make([]string, 0, len(targetProject.IssueTypes))
		for _, it := range targetProject.IssueTypes {
			availableTypes = append(availableTypes, it.Name)
		}
		return nil, fmt.Errorf("[jira] issue type '%s' not found. Available types: %v", issueTypeName, availableTypes)
	}

	// Extract field information
	fieldInfo := make([]map[string]interface{}, 0)
	for fieldKey, fieldValue := range targetIssueType.Fields {
		// The library uses interface{} for fields, need to cast
		fieldMap, ok := fieldValue.(map[string]interface{})
		if !ok {
			continue
		}

		info := map[string]interface{}{
			"key": fieldKey,
		}

		if name, ok := fieldMap["name"].(string); ok {
			info["name"] = name
		}

		if required, ok := fieldMap["required"].(bool); ok {
			info["required"] = required
		}

		if schema, ok := fieldMap["schema"].(map[string]interface{}); ok {
			if fieldType, ok := schema["type"].(string); ok {
				info["type"] = fieldType
			}
			if items, ok := schema["items"].(string); ok {
				info["items"] = items
			}
		}

		if allowedValues, ok := fieldMap["allowedValues"].([]interface{}); ok && len(allowedValues) > 0 {
			values := make([]string, 0, len(allowedValues))
			for _, v := range allowedValues {
				if valMap, ok := v.(map[string]interface{}); ok {
					if name, ok := valMap["name"].(string); ok {
						values = append(values, name)
					} else if value, ok := valMap["value"].(string); ok {
						values = append(values, value)
					}
				}
			}
			if len(values) > 0 {
				info["allowedValues"] = values
			}
		}

		fieldInfo = append(fieldInfo, info)
	}

	return map[string]interface{}{
		"issueType": issueTypeName,
		"fields":    fieldInfo,
	}, nil
}

// Helper functions

func (j *JiraAdapter) buildDescription(description string, acceptanceCriteria []string) string {
	if len(acceptanceCriteria) == 0 {
		return description
	}

	desc := description + "\n\nh3. Acceptance Criteria\n"
	for _, ac := range acceptanceCriteria {
		desc += fmt.Sprintf("* %s\n", ac)
	}
	return desc
}

func (j *JiraAdapter) getJiraFieldID(humanName string) string {
	if mapping, exists := j.fieldMappings[humanName]; exists {
		switch m := mapping.(type) {
		case string:
			return m
		case map[string]interface{}:
			if id, ok := m["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

func (j *JiraAdapter) convertFieldValue(fieldName, value string) interface{} {
	// Check if this is a number field
	if mapping, exists := j.fieldMappings[fieldName]; exists {
		if m, ok := mapping.(map[string]interface{}); ok {
			if fieldType, ok := m["type"].(string); ok && fieldType == "number" {
				var num float64
				if _, err := fmt.Sscanf(value, "%f", &num); err == nil {
					return num
				}
			}
		}
	}

	// Check if this is a known array field
	jiraField := j.getJiraFieldID(fieldName)
	if jiraField == "labels" || jiraField == "components" {
		if value == "" {
			return []string{}
		}
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}

	return value
}

func (j *JiraAdapter) applyCustomFields(fields *jira.IssueFields, customFields map[string]string) {
	for fieldName, fieldValue := range customFields {
		jiraField := j.getJiraFieldID(fieldName)
		if jiraField == "" {
			continue
		}

		// Skip standard fields handled separately
		if jiraField == "project" || jiraField == "issuetype" || jiraField == "summary" || jiraField == "description" {
			continue
		}

		fields.Unknowns[jiraField] = j.convertFieldValue(fieldName, fieldValue)
	}
}

func (j *JiraAdapter) convertToDomainTicket(issue *jira.Issue) domain.Ticket {
	ticket := domain.Ticket{
		JiraID:       issue.Key,
		Title:        issue.Fields.Summary,
		CustomFields: make(map[string]string),
		Tasks:        []domain.Task{}, // Will be populated by fetchSubtasks
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

	// Get issue type
	if issue.Fields.Type.Name != "" {
		ticket.CustomFields["Type"] = issue.Fields.Type.Name
	}

	// Get parent if it's a subtask
	if issue.Fields.Parent != nil && issue.Fields.Parent.Key != "" {
		ticket.CustomFields["Parent"] = issue.Fields.Parent.Key
	}

	// Map custom fields back to human-readable names
	reverseMapping := j.createReverseFieldMapping()
	for jiraField, humanName := range reverseMapping {
		if value, exists := issue.Fields.Unknowns[jiraField]; exists {
			ticket.CustomFields[humanName] = j.formatFieldValue(value)
		}
	}

	return ticket
}

func (j *JiraAdapter) convertToDomainTask(issue *jira.Issue) domain.Task {
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

	// Get issue type
	if issue.Fields.Type.Name != "" {
		task.CustomFields["Type"] = issue.Fields.Type.Name
	}

	// Map custom fields back to human-readable names
	reverseMapping := j.createReverseFieldMapping()
	for jiraField, humanName := range reverseMapping {
		if value, exists := issue.Fields.Unknowns[jiraField]; exists {
			task.CustomFields[humanName] = j.formatFieldValue(value)
		}
	}

	return task
}

func (j *JiraAdapter) createReverseFieldMapping() map[string]string {
	reverse := make(map[string]string)
	for humanName, mapping := range j.fieldMappings {
		switch m := mapping.(type) {
		case string:
			reverse[m] = humanName
		case map[string]interface{}:
			if id, ok := m["id"].(string); ok {
				reverse[id] = humanName
			}
		}
	}
	return reverse
}

func (j *JiraAdapter) formatFieldValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%g", v)
	case []interface{}:
		values := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				values = append(values, str)
			} else if obj, ok := item.(map[string]interface{}); ok {
				if name, ok := obj["name"].(string); ok {
					values = append(values, name)
				}
			}
		}
		if len(values) > 0 {
			return strings.Join(values, ", ")
		}
	case map[string]interface{}:
		if name, ok := v["name"].(string); ok {
			return name
		} else if displayName, ok := v["displayName"].(string); ok {
			return displayName
		}
	}
	return ""
}

// getDefaultFieldMappings returns default field mappings for Jira
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
