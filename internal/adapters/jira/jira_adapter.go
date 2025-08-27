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
