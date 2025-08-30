package jira

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

    "github.com/karolswdev/ticketr/internal/core/domain"
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
		Tasks: []domain.Task{},
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
		Tasks: []domain.Task{},
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
		Tasks: []domain.Task{},
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
	expectedURL := "https://test.atlassian.net/rest/api/2/search"
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
