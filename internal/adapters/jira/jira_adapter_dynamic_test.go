package jira

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

    "github.com/karolswdev/ticketr/internal/core/domain"
)

// MockRoundTripper is a mock for testing HTTP requests
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
	LastRequest   *http.Request
	LastBody      []byte
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.LastRequest = req
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		m.LastBody = body
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
	return m.RoundTripFunc(req)
}

func TestJiraAdapter_CreateTicket_DynamicPayload(t *testing.T) {
	// Create field mappings
	fieldMappings := map[string]interface{}{
		"Story Points": map[string]interface{}{
			"id":   "customfield_10010",
			"type": "number",
		},
		"Sprint": "customfield_10020",
		"Labels": "labels",
	}

	// Create a mock round tripper
	mockTransport := &MockRoundTripper{}
	capturedPayload := ""
	
	mockTransport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
		// Capture the request body
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			capturedPayload = string(body)
		}
		
		// Return a successful response
		responseBody := `{"id":"10000","key":"TEST-123","self":"https://test.atlassian.net/rest/api/2/issue/10000"}`
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(responseBody)),
		}, nil
	}

	// Create HTTP client with mock transport
	httpClient := &http.Client{
		Transport: mockTransport,
	}

	// Create adapter with custom field mappings
	adapter := &JiraAdapter{
		baseURL:       "https://test.atlassian.net",
		email:         "test@example.com",
		apiKey:        "test-key",
		projectKey:    "TEST",
		storyType:     "Task",
		client:        httpClient,
		fieldMappings: fieldMappings,
	}

	// Create a ticket with custom fields
	ticket := domain.Ticket{
		Title:       "Test Ticket",
		Description: "Test Description",
		CustomFields: map[string]string{
			"Story Points": "5",
			"Sprint":       "Sprint 23",
			"Labels":       "backend, api",
		},
	}

	// Call CreateTicket
	key, err := adapter.CreateTicket(ticket)
	if err != nil {
		t.Fatalf("CreateTicket failed: %v", err)
	}

	if key != "TEST-123" {
		t.Errorf("Expected key TEST-123, got %s", key)
	}

	// Parse the captured payload
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(capturedPayload), &payload); err != nil {
		t.Fatalf("Failed to parse captured payload: %v", err)
	}

	// Check that the payload contains the correct field mappings
	fields, ok := payload["fields"].(map[string]interface{})
	if !ok {
		t.Fatal("Payload does not contain 'fields'")
	}

	// Verify Story Points is mapped to customfield_10010 as a number
	if storyPoints, exists := fields["customfield_10010"]; !exists {
		t.Error("customfield_10010 (Story Points) not found in payload")
	} else {
		// Check that it's a number, not a string
		if _, ok := storyPoints.(float64); !ok {
			t.Errorf("Story Points should be a number, got %T: %v", storyPoints, storyPoints)
		} else if storyPoints.(float64) != 5 {
			t.Errorf("Expected Story Points to be 5, got %v", storyPoints)
		}
	}

	// Verify Sprint is mapped to customfield_10020
	if sprint, exists := fields["customfield_10020"]; !exists {
		t.Error("customfield_10020 (Sprint) not found in payload")
	} else if sprint != "Sprint 23" {
		t.Errorf("Expected Sprint to be 'Sprint 23', got %v", sprint)
	}

	// Verify Labels are mapped correctly as an array
	if labels, exists := fields["labels"]; !exists {
		t.Error("labels not found in payload")
	} else {
		if labelArray, ok := labels.([]interface{}); !ok {
			t.Errorf("Labels should be an array, got %T", labels)
		} else if len(labelArray) != 2 {
			t.Errorf("Expected 2 labels, got %d", len(labelArray))
		}
	}

	// Verify standard fields
	if summary, exists := fields["summary"]; !exists || summary != "Test Ticket" {
		t.Errorf("Expected summary 'Test Ticket', got %v", summary)
	}

	if description, exists := fields["description"]; !exists || description != "Test Description" {
		t.Errorf("Expected description 'Test Description', got %v", description)
	}
}
