package jira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// JiraAdapter implements the JiraPort interface for Jira API integration
type JiraAdapter struct {
	baseURL       string
	email         string
	apiKey        string
	projectKey    string
	storyType     string
	subTaskType   string
	client        *http.Client
}

// NewJiraAdapter creates a new instance of JiraAdapter using environment variables
func NewJiraAdapter() (ports.JiraPort, error) {
	baseURL := os.Getenv("JIRA_URL")
	email := os.Getenv("JIRA_EMAIL")
	apiKey := os.Getenv("JIRA_API_KEY")
	projectKey := os.Getenv("JIRA_PROJECT_KEY")

	if baseURL == "" || email == "" || apiKey == "" || projectKey == "" {
		return nil, fmt.Errorf("missing required environment variables: JIRA_URL, JIRA_EMAIL, JIRA_API_KEY, JIRA_PROJECT_KEY")
	}

	// Get issue types from environment, with sensible defaults
	storyType := os.Getenv("JIRA_STORY_TYPE")
	if storyType == "" {
		storyType = "Task" // Default to Task which is more common
	}
	
	subTaskType := os.Getenv("JIRA_SUBTASK_TYPE")
	if subTaskType == "" {
		subTaskType = "Sub-task" // Standard JIRA subtask type
	}

	// Ensure base URL doesn't have trailing slash
	baseURL = strings.TrimRight(baseURL, "/")

	return &JiraAdapter{
		baseURL:      baseURL,
		email:        email,
		apiKey:       apiKey,
		projectKey:   projectKey,
		storyType:    storyType,
		subTaskType:  subTaskType,
		client:       &http.Client{},
	}, nil
}

// getAuthHeader returns the base64 encoded authentication header value
func (j *JiraAdapter) getAuthHeader() string {
	auth := fmt.Sprintf("%s:%s", j.email, j.apiKey)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Authenticate verifies the connection to Jira with the provided credentials
func (j *JiraAdapter) Authenticate() error {
	// Use the myself endpoint to verify authentication
	url := fmt.Sprintf("%s/rest/api/2/myself", j.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.getAuthHeader()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// CreateStory creates a new story in Jira and returns the created issue key
func (j *JiraAdapter) CreateStory(story domain.Story) (string, error) {
	// Build the description with acceptance criteria
	description := story.Description
	if len(story.AcceptanceCriteria) > 0 {
		description += "\n\nh3. Acceptance Criteria\n"
		for _, ac := range story.AcceptanceCriteria {
			description += fmt.Sprintf("* %s\n", ac)
		}
	}

	// Create the request payload
	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{
				"key": j.projectKey,
			},
			"summary":     story.Title,
			"description": description,
			"issuetype": map[string]string{
				"name": j.storyType,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/2/issue", j.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.getAuthHeader()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create story with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response to get the issue key
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	key, ok := result["key"].(string)
	if !ok {
		return "", fmt.Errorf("response did not contain issue key")
	}

	return key, nil
}

// CreateTask creates a new sub-task in Jira under the specified parent story
func (j *JiraAdapter) CreateTask(task domain.Task, parentID string) (string, error) {
	// Build the description with acceptance criteria
	description := task.Description
	if len(task.AcceptanceCriteria) > 0 {
		description += "\n\nh3. Acceptance Criteria\n"
		for _, ac := range task.AcceptanceCriteria {
			description += fmt.Sprintf("* %s\n", ac)
		}
	}

	// Create the request payload
	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{
				"key": j.projectKey,
			},
			"summary":     task.Title,
			"description": description,
			"issuetype": map[string]string{
				"name": j.subTaskType,
			},
			"parent": map[string]string{
				"key": parentID,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/2/issue", j.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.getAuthHeader()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create task with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response to get the issue key
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	key, ok := result["key"].(string)
	if !ok {
		return "", fmt.Errorf("response did not contain issue key")
	}

	return key, nil
}

// UpdateStory updates an existing story in Jira
func (j *JiraAdapter) UpdateStory(story domain.Story) error {
	if story.JiraID == "" {
		return fmt.Errorf("story does not have a Jira ID")
	}

	// Build the description with acceptance criteria
	description := story.Description
	if len(story.AcceptanceCriteria) > 0 {
		description += "\n\nh3. Acceptance Criteria\n"
		for _, ac := range story.AcceptanceCriteria {
			description += fmt.Sprintf("* %s\n", ac)
		}
	}

	// Create the request payload - only update fields that can change
	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"summary":     story.Title,
			"description": description,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/2/issue/%s", j.baseURL, story.JiraID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.getAuthHeader()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update story with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetProjectIssueTypes fetches available issue types for the configured project
func (j *JiraAdapter) GetProjectIssueTypes() (map[string][]string, error) {
	url := fmt.Sprintf("%s/rest/api/2/project/%s", j.baseURL, j.projectKey)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.getAuthHeader()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get project with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response
	var project map[string]interface{}
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := make(map[string][]string)
	
	// Get project name
	if name, ok := project["name"].(string); ok {
		result["project"] = []string{name}
	}
	
	// Get project key
	if key, ok := project["key"].(string); ok {
		result["key"] = []string{key}
	}
	
	// Get issue types
	issueTypes := []string{}
	if types, ok := project["issueTypes"].([]interface{}); ok {
		for _, t := range types {
			if typeMap, ok := t.(map[string]interface{}); ok {
				if name, ok := typeMap["name"].(string); ok {
					// Check if it's a subtask type
					if subtask, ok := typeMap["subtask"].(bool); ok && subtask {
						issueTypes = append(issueTypes, fmt.Sprintf("%s (subtask)", name))
					} else {
						issueTypes = append(issueTypes, name)
					}
				}
			}
		}
	}
	result["issueTypes"] = issueTypes
	
	return result, nil
}

// GetIssueTypeFields fetches field requirements for a specific issue type
func (j *JiraAdapter) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	// Use the createmeta endpoint to get field information
	url := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&expand=projects.issuetypes.fields", 
		j.baseURL, j.projectKey)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.getAuthHeader()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get createmeta with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response
	var createMeta map[string]interface{}
	if err := json.Unmarshal(body, &createMeta); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := make(map[string]interface{})
	
	// Navigate through the response structure
	projects, ok := createMeta["projects"].([]interface{})
	if !ok || len(projects) == 0 {
		return nil, fmt.Errorf("no projects found in response")
	}
	
	project := projects[0].(map[string]interface{})
	issueTypes, ok := project["issuetypes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no issue types found in project")
	}
	
	// Find the requested issue type
	var targetIssueType map[string]interface{}
	for _, it := range issueTypes {
		issueType := it.(map[string]interface{})
		if name, ok := issueType["name"].(string); ok && name == issueTypeName {
			targetIssueType = issueType
			break
		}
	}
	
	if targetIssueType == nil {
		// List available issue types
		availableTypes := []string{}
		for _, it := range issueTypes {
			if issueType, ok := it.(map[string]interface{}); ok {
				if name, ok := issueType["name"].(string); ok {
					availableTypes = append(availableTypes, name)
				}
			}
		}
		return nil, fmt.Errorf("issue type '%s' not found. Available types: %v", issueTypeName, availableTypes)
	}
	
	// Extract field information
	fields, ok := targetIssueType["fields"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no fields found for issue type")
	}
	
	// Process fields to extract relevant information
	fieldInfo := []map[string]interface{}{}
	for fieldKey, fieldData := range fields {
		field := fieldData.(map[string]interface{})
		
		info := map[string]interface{}{
			"key": fieldKey,
		}
		
		if name, ok := field["name"].(string); ok {
			info["name"] = name
		}
		
		if required, ok := field["required"].(bool); ok {
			info["required"] = required
		}
		
		if schema, ok := field["schema"].(map[string]interface{}); ok {
			if fieldType, ok := schema["type"].(string); ok {
				info["type"] = fieldType
			}
			if items, ok := schema["items"].(string); ok {
				info["items"] = items
			}
		}
		
		if allowedValues, ok := field["allowedValues"].([]interface{}); ok && len(allowedValues) > 0 {
			values := []string{}
			for _, v := range allowedValues {
				if val, ok := v.(map[string]interface{}); ok {
					if name, ok := val["name"].(string); ok {
						values = append(values, name)
					} else if value, ok := val["value"].(string); ok {
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
	
	result["issueType"] = issueTypeName
	result["fields"] = fieldInfo
	
	return result, nil
}

// UpdateTask updates an existing task in Jira
func (j *JiraAdapter) UpdateTask(task domain.Task) error {
	if task.JiraID == "" {
		return fmt.Errorf("task does not have a Jira ID")
	}

	// Build the description with acceptance criteria
	description := task.Description
	if len(task.AcceptanceCriteria) > 0 {
		description += "\n\nh3. Acceptance Criteria\n"
		for _, ac := range task.AcceptanceCriteria {
			description += fmt.Sprintf("* %s\n", ac)
		}
	}

	// Create the request payload - only update fields that can change
	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"summary":     task.Title,
			"description": description,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/2/issue/%s", j.baseURL, task.JiraID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.getAuthHeader()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update task with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
