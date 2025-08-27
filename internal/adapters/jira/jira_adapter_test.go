package jira

import (
	"os"
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

	// Arrange: Create a valid Story domain object
	story := domain.Story{
		Title:       "Test Story from Integration Test",
		Description: "This is a test story created by the integration test suite",
		AcceptanceCriteria: []string{
			"The story should be created in Jira",
			"A valid Jira ID should be returned",
		},
		Tasks: []domain.Task{},
	}

	adapter, err := NewJiraAdapter()
	if err != nil {
		t.Fatalf("Failed to create Jira adapter: %v", err)
	}

	// Act: Call the CreateStory method on the Jira adapter
	jiraID, err := adapter.CreateStory(story)
	
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

	// First create a story to update
	initialStory := domain.Story{
		Title:       "Test Story for Update Integration Test",
		Description: "Initial description for update test",
		AcceptanceCriteria: []string{
			"Initial acceptance criterion",
		},
		Tasks: []domain.Task{},
	}

	jiraID, err := adapter.CreateStory(initialStory)
	if err != nil {
		t.Fatalf("Failed to create initial story: %v", err)
	}
	t.Logf("Created story with Jira ID: %s", jiraID)

	// Create a Story domain object with that ID and modified description
	updatedStory := domain.Story{
		JiraID:      jiraID,
		Title:       "Updated Test Story from Integration Test",
		Description: "This description has been updated by the integration test",
		AcceptanceCriteria: []string{
			"Updated acceptance criterion 1",
			"Updated acceptance criterion 2",
		},
		Tasks: []domain.Task{},
	}

	// Act: Call the UpdateStory method on the Jira adapter
	err = adapter.UpdateStory(updatedStory)

	// Assert: The method succeeds and the description in Jira is updated
	if err != nil {
		t.Errorf("Failed to update story: %v", err)
	} else {
		t.Logf("Successfully updated story with Jira ID: %s", jiraID)
	}
}