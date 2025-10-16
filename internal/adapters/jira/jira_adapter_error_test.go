package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// TestJiraAdapter_CreateTicket_APIError tests handling of API errors during ticket creation
func TestJiraAdapter_CreateTicket_APIError(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"errorMessages":["Permission denied"]}`
			return &http.Response{
				StatusCode: 403,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	ticket := domain.Ticket{
		Title:       "Test Ticket",
		Description: "Test Description",
	}

	_, err := adapter.CreateTicket(ticket)

	if err == nil {
		t.Fatal("Expected error for 403 response, got nil")
	}

	if !strings.Contains(err.Error(), "403") && !strings.Contains(err.Error(), "Permission denied") {
		t.Errorf("Error message should mention permission issue, got: %v", err)
	}
}

// TestJiraAdapter_CreateTicket_MalformedResponse tests handling of malformed JSON responses
func TestJiraAdapter_CreateTicket_MalformedResponse(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"key": "TEST-123", invalid json here}`
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	ticket := domain.Ticket{
		Title:       "Test Ticket",
		Description: "Test Description",
	}

	_, err := adapter.CreateTicket(ticket)

	if err == nil {
		t.Fatal("Expected error for malformed JSON, got nil")
	}
}

// TestJiraAdapter_CreateTicket_NetworkError tests handling of network errors
func TestJiraAdapter_CreateTicket_NetworkError(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network timeout")
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	ticket := domain.Ticket{
		Title:       "Test Ticket",
		Description: "Test Description",
	}

	_, err := adapter.CreateTicket(ticket)

	if err == nil {
		t.Fatal("Expected error for network timeout, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("Error message should mention timeout, got: %v", err)
	}
}

// TestJiraAdapter_UpdateTicket_APIError tests handling of API errors during ticket update
func TestJiraAdapter_UpdateTicket_APIError(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"errorMessages":["Ticket not found"]}`
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	ticket := domain.Ticket{
		JiraID:      "PROJ-999",
		Title:       "Test Ticket",
		Description: "Test Description",
	}

	err := adapter.UpdateTicket(ticket)

	if err == nil {
		t.Fatal("Expected error for 404 response, got nil")
	}

	if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
		t.Errorf("Error message should mention not found, got: %v", err)
	}
}

// TestJiraAdapter_UpdateTicket_EmptyJiraID tests handling of missing Jira ID
func TestJiraAdapter_UpdateTicket_EmptyJiraID(t *testing.T) {
	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		client:        &http.Client{},
		fieldMappings: getDefaultFieldMappings(),
	}

	ticket := domain.Ticket{
		JiraID:      "", // Empty Jira ID
		Title:       "Test Ticket",
		Description: "Test Description",
	}

	err := adapter.UpdateTicket(ticket)

	if err == nil {
		t.Fatal("Expected error for empty Jira ID, got nil")
	}
}

// TestJiraAdapter_SearchTickets_APIError tests handling of API errors during search
func TestJiraAdapter_SearchTickets_APIError(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"errorMessages":["Invalid JQL query"]}`
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		subTaskType:   "Sub-task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	_, err := adapter.SearchTickets("PROJ", "invalid jql $$")

	if err == nil {
		t.Fatal("Expected error for invalid JQL, got nil")
	}

	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Error message should mention 400 status, got: %v", err)
	}
}

// TestJiraAdapter_SearchTickets_EmptyResponse tests handling of empty search results
func TestJiraAdapter_SearchTickets_EmptyResponse(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"issues": [], "total": 0}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		subTaskType:   "Sub-task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	tickets, err := adapter.SearchTickets("PROJ", "status = Done")

	if err != nil {
		t.Fatalf("Expected no error for empty results, got: %v", err)
	}

	if len(tickets) != 0 {
		t.Errorf("Expected 0 tickets, got %d", len(tickets))
	}
}

// TestJiraAdapter_SearchTickets_MalformedJiraResponse tests handling of malformed Jira responses
func TestJiraAdapter_SearchTickets_MalformedJiraResponse(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"issues": [{"key": "PROJ-1", "fields": null}]}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		subTaskType:   "Sub-task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	// Should handle null fields gracefully
	tickets, err := adapter.SearchTickets("PROJ", "")

	if err != nil {
		t.Fatalf("Expected no error for null fields, got: %v", err)
	}

	// Should still return tickets even with null fields
	if len(tickets) == 0 {
		t.Error("Expected at least 1 ticket despite null fields")
	}
}

// TestJiraAdapter_Authenticate_InvalidCredentials tests authentication failure
func TestJiraAdapter_Authenticate_InvalidCredentials(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"errorMessages":["Invalid credentials"]}`
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "wrong@example.com",
		apiKey:        "wrong-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	err := adapter.Authenticate()

	if err == nil {
		t.Fatal("Expected error for invalid credentials, got nil")
	}

	if !strings.Contains(err.Error(), "401") && !strings.Contains(err.Error(), "credentials") {
		t.Errorf("Error message should mention authentication failure, got: %v", err)
	}
}

// TestJiraAdapter_FieldMapping_MissingFields tests handling of missing custom fields
func TestJiraAdapter_FieldMapping_MissingFields(t *testing.T) {
	requestCount := 0
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			requestCount++

			if requestCount == 1 {
				// Response without some expected fields
				responseBody := `{
					"issues": [{
						"key": "PROJ-500",
						"fields": {
							"summary": "Ticket with Missing Fields",
							"issuetype": {"name": "Task"}
						}
					}],
					"total": 1
				}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
				}, nil
			}

			// No subtasks
			responseBody := `{"issues": [], "total": 0}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		subTaskType:   "Sub-task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	tickets, err := adapter.SearchTickets("PROJ", "")

	if err != nil {
		t.Fatalf("Expected no error for missing fields, got: %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	// Missing description should result in empty string, not error
	if tickets[0].Description != "" {
		t.Logf("Description present (OK): %s", tickets[0].Description)
	}
}

// TestJiraAdapter_RateLimiting tests handling of rate limiting responses
func TestJiraAdapter_RateLimiting(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"errorMessages":["Rate limit exceeded"]}`
			return &http.Response{
				StatusCode: 429,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
				Header: http.Header{
					"Retry-After": []string{"60"},
				},
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	ticket := domain.Ticket{
		Title:       "Test Ticket",
		Description: "Test Description",
	}

	_, err := adapter.CreateTicket(ticket)

	if err == nil {
		t.Fatal("Expected error for rate limiting, got nil")
	}

	if !strings.Contains(err.Error(), "429") && !strings.Contains(err.Error(), "rate") {
		t.Errorf("Error message should mention rate limiting, got: %v", err)
	}
}

// TestJiraAdapter_CreateTicket_SpecialCharacters tests handling of special characters in ticket data
func TestJiraAdapter_CreateTicket_SpecialCharacters(t *testing.T) {
	var capturedPayload string
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.Body != nil {
				body, _ := io.ReadAll(req.Body)
				capturedPayload = string(body)
			}

			responseBody := `{"id":"10000","key":"TEST-999","self":"https://test.atlassian.net/rest/api/2/issue/10000"}`
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-api-key",
		projectKey:    "PROJ",
		storyType:     "Task",
		client:        &http.Client{Transport: mockTransport},
		fieldMappings: getDefaultFieldMappings(),
	}

	ticket := domain.Ticket{
		Title:       "Test with \"quotes\" and 'apostrophes'",
		Description: "Description with special chars: <>&\n\tTabbed and newlined",
		AcceptanceCriteria: []string{
			"AC with \\backslashes\\",
			"AC with Unicode: 你好",
		},
	}

	key, err := adapter.CreateTicket(ticket)

	if err != nil {
		t.Fatalf("Expected no error for special characters, got: %v", err)
	}

	if key != "TEST-999" {
		t.Errorf("Expected key TEST-999, got %s", key)
	}

	// Verify payload is valid JSON despite special characters
	if capturedPayload == "" {
		t.Fatal("No payload was captured")
	}

	// Should be able to parse as JSON
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(capturedPayload), &payload); err != nil {
		t.Errorf("Payload should be valid JSON: %v", err)
	} else {
		t.Logf("Successfully parsed payload with %d fields", len(payload))
	}

	t.Logf("Successfully handled special characters in payload")
}

// TestNewJiraAdapter_MissingEnvVars tests error handling when environment variables are missing
func TestNewJiraAdapter_MissingEnvVars(t *testing.T) {
	// Save current env vars
	originalURL := os.Getenv("JIRA_URL")
	originalEmail := os.Getenv("JIRA_EMAIL")
	originalKey := os.Getenv("JIRA_API_KEY")
	originalProject := os.Getenv("JIRA_PROJECT_KEY")

	// Unset env vars
	os.Unsetenv("JIRA_URL")
	os.Unsetenv("JIRA_EMAIL")
	os.Unsetenv("JIRA_API_KEY")
	os.Unsetenv("JIRA_PROJECT_KEY")

	defer func() {
		// Restore env vars
		os.Setenv("JIRA_URL", originalURL)
		os.Setenv("JIRA_EMAIL", originalEmail)
		os.Setenv("JIRA_API_KEY", originalKey)
		os.Setenv("JIRA_PROJECT_KEY", originalProject)
	}()

	_, err := NewJiraAdapter()

	if err == nil {
		t.Fatal("Expected error for missing env vars, got nil")
	}

	if !strings.Contains(err.Error(), "JIRA") {
		t.Errorf("Error message should mention missing JIRA config, got: %v", err)
	}
}
