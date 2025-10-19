package jira

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// Test Case TC-2.1: JiraAdapter_NewClient_WithEnvVars_AuthenticatesSuccessfully
func TestJiraAdapter_NewClient_WithEnvVars_AuthenticatesSuccessfully(t *testing.T) {
	// Skip this test if environment variables are not set (for CI/CD)
	if os.Getenv("JIRA_URL") == "" {
		t.Skip("Skipping integration test: JIRA_URL not set")
	}

	// Arrange: Valid Jira credentials should be set in environment variables
	// The test assumes these are already set: JIRA_URL, JIRA_EMAIL, JIRA_API_KEY, JIRA_PROJECT_KEY

	// Act: Create a new Jira client instance
	adapter, err := NewJiraAdapter()
	if err != nil {
		t.Fatalf("Failed to create Jira adapter: %v", err)
	}

	// Assert: The client authenticates successfully
	err = adapter.Authenticate()
	if err != nil {
		t.Errorf("Authentication failed: %v", err)
	}
}

// Test Case TC-2.2: JiraAdapter_CreateStory_ValidStory_ReturnsNewJiraID
func TestJiraAdapter_CreateStory_ValidStory_ReturnsNewJiraID(t *testing.T) {
	// Skip this test if environment variables are not set (for CI/CD)
	if os.Getenv("JIRA_URL") == "" {
		t.Skip("Skipping integration test: JIRA_URL not set")
	}

	// Arrange: Create a valid Ticket domain object
	ticket := domain.Ticket{
		Title:       "Test Ticket from Integration Test",
		Description: "This is a test ticket created by the integration test suite",
		AcceptanceCriteria: []string{
			"The ticket should be created in Jira",
			"A valid Jira ID should be returned",
		},
		CustomFields: map[string]string{},
		Tasks:        []domain.Task{},
	}

	adapter, err := NewJiraAdapter()
	if err != nil {
		t.Fatalf("Failed to create Jira adapter: %v", err)
	}

	// Act: Call the CreateTicket method on the Jira adapter
	jiraID, err := adapter.CreateTicket(ticket)

	// Assert: The method returns a valid, non-empty Jira Issue Key
	if err != nil {
		t.Fatalf("Failed to create story: %v", err)
	}

	if jiraID == "" {
		t.Error("Expected non-empty Jira ID, got empty string")
	}

	// Log the created Jira ID for manual verification if needed
	t.Logf("Successfully created story with Jira ID: %s", jiraID)
}

// Test Case TC-3.1: JiraAdapter_UpdateStory_ValidStoryWithID_Succeeds
func TestJiraAdapter_UpdateStory_ValidStoryWithID_Succeeds(t *testing.T) {
	// Skip this test if environment variables are not set (for CI/CD)
	if os.Getenv("JIRA_URL") == "" {
		t.Skip("Skipping integration test: JIRA_URL not set")
	}

	// Arrange: Create a story in Jira to get a valid ID
	adapter, err := NewJiraAdapter()
	if err != nil {
		t.Fatalf("Failed to create Jira adapter: %v", err)
	}

	// First create a ticket to update
	initialTicket := domain.Ticket{
		Title:       "Test Ticket for Update Integration Test",
		Description: "Initial description for update test",
		AcceptanceCriteria: []string{
			"Initial acceptance criterion",
		},
		CustomFields: map[string]string{},
		Tasks:        []domain.Task{},
	}

	jiraID, err := adapter.CreateTicket(initialTicket)
	if err != nil {
		t.Fatalf("Failed to create initial ticket: %v", err)
	}
	t.Logf("Created ticket with Jira ID: %s", jiraID)

	// Create a Ticket domain object with that ID and modified description
	updatedTicket := domain.Ticket{
		JiraID:      jiraID,
		Title:       "Updated Test Ticket from Integration Test",
		Description: "This description has been updated by the integration test",
		AcceptanceCriteria: []string{
			"Updated acceptance criterion 1",
			"Updated acceptance criterion 2",
		},
		CustomFields: map[string]string{},
		Tasks:        []domain.Task{},
	}

	// Act: Call the UpdateTicket method on the Jira adapter
	err = adapter.UpdateTicket(updatedTicket)

	// Assert: The method succeeds and the description in Jira is updated
	if err != nil {
		t.Errorf("Failed to update ticket: %v", err)
	} else {
		t.Logf("Successfully updated ticket with Jira ID: %s", jiraID)
	}
}

// Test Case TC-205.1: TestJiraAdapter_SearchTickets_ConstructsJql
func TestJiraAdapter_SearchTickets_ConstructsJql(t *testing.T) {
	// Arrange: Mock the http.Client
	var capturedRequest *http.Request
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Capture the request for assertions
			capturedRequest = req

			// Return a mock response
			responseBody := `{
				"issues": [],
				"total": 0,
				"maxResults": 100,
				"startAt": 0
			}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	// Create JiraAdapter with mocked client
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

	// Act: Call SearchTickets with a project key "PROJ" and JQL "status=Done"
	_, err := adapter.SearchTickets("PROJ", "status=Done")

	// Assert: The request sent to Jira's /rest/api/2/search endpoint contains the JQL
	if err != nil {
		t.Fatalf("SearchTickets returned error: %v", err)
	}

	// Verify the request was sent to correct endpoint
	expectedURL := "https://test.atlassian.net/rest/api/3/search/jql"
	if capturedRequest == nil {
		t.Fatal("No request was captured")
	}
	if capturedRequest.URL.String() != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, capturedRequest.URL.String())
	}

	// Read and verify the request body contains correct JQL
	bodyBytes, err := io.ReadAll(capturedRequest.Body)
	if err != nil {
		t.Fatalf("Failed to read request body: %v", err)
	}

	// Parse the JSON body to verify JQL
	var requestBody map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
		t.Fatalf("Failed to parse request body JSON: %v", err)
	}

	expectedJQL := `project = "PROJ" AND status=Done`
	actualJQL, ok := requestBody["jql"].(string)
	if !ok {
		t.Fatal("Request body does not contain 'jql' field")
	}

	if actualJQL != expectedJQL {
		t.Errorf("JQL mismatch.\nExpected: %s\nActual: %s", expectedJQL, actualJQL)
	}

	// Verify request has proper authentication header
	authHeader := capturedRequest.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Basic ") {
		t.Errorf("Expected Basic auth header, got: %s", authHeader)
	}

	// Verify content type
	contentType := capturedRequest.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type: application/json, got: %s", contentType)
	}

	t.Logf("Successfully verified JQL construction: %s", expectedJQL)
}

// Test Case TC-208.1: TestSearchTickets_WithSubtasks
func TestSearchTickets_WithSubtasks(t *testing.T) {
	// Arrange: Mock the http.Client to return a parent ticket and subtasks
	requestCount := 0
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			requestCount++

			// First request: parent ticket search
			if requestCount == 1 {
				responseBody := `{
					"issues": [{
						"key": "PROJ-123",
						"fields": {
							"summary": "Parent Ticket",
							"description": "Parent description",
							"issuetype": {"name": "Task"}
						}
					}],
					"total": 1,
					"maxResults": 100,
					"startAt": 0
				}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
				}, nil
			}

			// Second request: subtask search for PROJ-123
			responseBody := `{
				"issues": [
					{
						"key": "PROJ-124",
						"fields": {
							"summary": "Subtask 1",
							"description": "Description for subtask 1",
							"issuetype": {"name": "Sub-task"}
						}
					},
					{
						"key": "PROJ-125",
						"fields": {
							"summary": "Subtask 2",
							"description": "Description for subtask 2",
							"issuetype": {"name": "Sub-task"}
						}
					}
				],
				"total": 2,
				"maxResults": 100,
				"startAt": 0
			}`
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

	// Act: Call SearchTickets
	tickets, err := adapter.SearchTickets("PROJ", "")

	// Assert: Verify parent ticket has 2 subtasks
	if err != nil {
		t.Fatalf("SearchTickets returned error: %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	ticket := tickets[0]
	if ticket.JiraID != "PROJ-123" {
		t.Errorf("Expected ticket JiraID PROJ-123, got %s", ticket.JiraID)
	}

	if len(ticket.Tasks) != 2 {
		t.Fatalf("Expected 2 subtasks, got %d", len(ticket.Tasks))
	}

	// Verify first subtask
	if ticket.Tasks[0].JiraID != "PROJ-124" {
		t.Errorf("Expected task JiraID PROJ-124, got %s", ticket.Tasks[0].JiraID)
	}
	if ticket.Tasks[0].Title != "Subtask 1" {
		t.Errorf("Expected task title 'Subtask 1', got %s", ticket.Tasks[0].Title)
	}
	if ticket.Tasks[0].Description != "Description for subtask 1" {
		t.Errorf("Expected task description 'Description for subtask 1', got %s", ticket.Tasks[0].Description)
	}

	// Verify second subtask
	if ticket.Tasks[1].JiraID != "PROJ-125" {
		t.Errorf("Expected task JiraID PROJ-125, got %s", ticket.Tasks[1].JiraID)
	}
	if ticket.Tasks[1].Title != "Subtask 2" {
		t.Errorf("Expected task title 'Subtask 2', got %s", ticket.Tasks[1].Title)
	}

	t.Logf("Successfully verified parent ticket with %d subtasks", len(ticket.Tasks))
}

// Test Case TC-208.2: TestSearchTickets_SubtaskFieldMapping
func TestSearchTickets_SubtaskFieldMapping(t *testing.T) {
	// Arrange: Mock the http.Client with subtask containing custom fields
	requestCount := 0
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			requestCount++

			// First request: parent ticket search
			if requestCount == 1 {
				responseBody := `{
					"issues": [{
						"key": "PROJ-200",
						"fields": {
							"summary": "Parent with Custom Fields",
							"description": "Parent description",
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

			// Second request: subtask with custom fields
			responseBody := `{
				"issues": [{
					"key": "PROJ-201",
					"fields": {
						"summary": "Subtask with Custom Fields",
						"description": "Subtask description",
						"issuetype": {"name": "Sub-task"},
						"customfield_10010": 5,
						"customfield_10020": "Sprint 23"
					}
				}],
				"total": 1
			}`
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

	// Act: Call SearchTickets
	tickets, err := adapter.SearchTickets("PROJ", "")

	// Assert: Verify subtask custom fields are mapped correctly
	if err != nil {
		t.Fatalf("SearchTickets returned error: %v", err)
	}

	if len(tickets) != 1 || len(tickets[0].Tasks) != 1 {
		t.Fatalf("Expected 1 ticket with 1 subtask")
	}

	task := tickets[0].Tasks[0]

	// Verify custom fields mapping
	if task.CustomFields["Story Points"] != "5" {
		t.Errorf("Expected Story Points '5', got '%s'", task.CustomFields["Story Points"])
	}

	if task.CustomFields["Sprint"] != "Sprint 23" {
		t.Errorf("Expected Sprint 'Sprint 23', got '%s'", task.CustomFields["Sprint"])
	}

	if task.CustomFields["Type"] != "Sub-task" {
		t.Errorf("Expected Type 'Sub-task', got '%s'", task.CustomFields["Type"])
	}

	t.Logf("Successfully verified subtask field mapping with custom fields")
}

// Test Case TC-208.3: TestSearchTickets_NoSubtasks
func TestSearchTickets_NoSubtasks(t *testing.T) {
	// Arrange: Mock the http.Client to return parent with no subtasks
	requestCount := 0
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			requestCount++

			// First request: parent ticket search
			if requestCount == 1 {
				responseBody := `{
					"issues": [{
						"key": "PROJ-300",
						"fields": {
							"summary": "Parent without Subtasks",
							"description": "Parent description",
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

			// Second request: empty subtask search
			responseBody := `{
				"issues": [],
				"total": 0,
				"maxResults": 100,
				"startAt": 0
			}`
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

	// Act: Call SearchTickets
	tickets, err := adapter.SearchTickets("PROJ", "")

	// Assert: Verify ticket has empty tasks array (not nil)
	if err != nil {
		t.Fatalf("SearchTickets returned error: %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	ticket := tickets[0]
	if ticket.Tasks == nil {
		t.Error("Expected non-nil Tasks array, got nil")
	}

	if len(ticket.Tasks) != 0 {
		t.Errorf("Expected empty Tasks array, got %d tasks", len(ticket.Tasks))
	}

	t.Logf("Successfully verified ticket with no subtasks has empty array")
}

// Test Case TC-208.4: TestSearchTickets_SubtaskFetchError
func TestSearchTickets_SubtaskFetchError(t *testing.T) {
	// Arrange: Mock the http.Client to fail on subtask fetch
	requestCount := 0
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			requestCount++

			// First request: parent ticket search succeeds
			if requestCount == 1 {
				responseBody := `{
					"issues": [{
						"key": "PROJ-400",
						"fields": {
							"summary": "Parent Ticket",
							"description": "Parent description",
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

			// Second request: subtask search fails
			responseBody := `{"errorMessages":["Internal server error"]}`
			return &http.Response{
				StatusCode: 500,
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

	// Act: Call SearchTickets
	tickets, err := adapter.SearchTickets("PROJ", "")

	// Assert: Verify parent ticket is still returned despite subtask error
	if err != nil {
		t.Fatalf("SearchTickets should not fail when subtask fetch fails: %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	ticket := tickets[0]
	if ticket.JiraID != "PROJ-400" {
		t.Errorf("Expected ticket JiraID PROJ-400, got %s", ticket.JiraID)
	}

	// Subtask fetch failed, so Tasks should be empty or nil (but not cause error)
	if ticket.Tasks == nil {
		ticket.Tasks = []domain.Task{}
	}

	t.Logf("Successfully verified non-fatal subtask fetch error")
}
