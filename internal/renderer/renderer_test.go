package renderer

import (
	"os"
	"strings"
	"testing"

	"github.com/karolswdev/ticketr/internal/core/domain"
)

// Test Case TC-206.1: TestRenderer_GeneratesCorrectMarkdown
func TestRenderer_GeneratesCorrectMarkdown(t *testing.T) {
	// Arrange: Create a domain.Ticket object with a title, description, and custom fields
	ticket := domain.Ticket{
		Title:       "Implement User Authentication",
		Description: "As a developer, I want to implement secure user authentication so that users can safely access the application.",
		CustomFields: map[string]string{
			"Priority":     "High",
			"Story Points": "5",
			"Sprint":       "Sprint 23",
		},
		AcceptanceCriteria: []string{
			"Users can register with email and password",
			"Passwords are securely hashed",
			"Session management is implemented",
		},
		JiraID: "PROJ-123",
		Tasks: []domain.Task{
			{
				Title:       "Set up authentication database schema",
				Description: "Create necessary database tables for user authentication",
				CustomFields: map[string]string{
					"Assignee": "john.doe",
				},
				JiraID: "PROJ-124",
			},
			{
				Title:       "Implement password hashing service",
				Description: "Create a service that securely hashes passwords using bcrypt",
				JiraID:      "PROJ-125",
			},
		},
	}

	// Act: Pass the ticket to a new renderer.Render function
	renderer := NewRenderer(nil)
	result := renderer.Render(ticket)

	// Assert: The returned string is a well-formed # TICKET: block
	// Check for main ticket structure
	if !strings.Contains(result, "# TICKET: [PROJ-123] Implement User Authentication") {
		t.Errorf("Result does not contain expected ticket header with JIRA ID")
	}

	if !strings.Contains(result, "## Description") {
		t.Errorf("Result does not contain Description section")
	}

	if !strings.Contains(result, "As a developer, I want to implement secure user authentication") {
		t.Errorf("Result does not contain expected description text")
	}

	if !strings.Contains(result, "## Acceptance Criteria") {
		t.Errorf("Result does not contain Acceptance Criteria section")
	}

	if !strings.Contains(result, "- Users can register with email and password") {
		t.Errorf("Result does not contain expected acceptance criteria")
	}

	if !strings.Contains(result, "## Fields") {
		t.Errorf("Result does not contain Fields section")
	}

	if !strings.Contains(result, "- Priority: High") {
		t.Errorf("Result does not contain Priority field")
	}

	if !strings.Contains(result, "- Story Points: 5") {
		t.Errorf("Result does not contain Story Points field")
	}

	if !strings.Contains(result, "## Tasks") {
		t.Errorf("Result does not contain Tasks section")
	}

	if !strings.Contains(result, "- [PROJ-124] Set up authentication database schema") {
		t.Errorf("Result does not contain first task with JIRA ID")
	}

	if !strings.Contains(result, "- [PROJ-125] Implement password hashing service") {
		t.Errorf("Result does not contain second task with JIRA ID")
	}

	// Optional: Compare with golden file if it exists
	goldenFile := "testdata/rendered_ticket.md"
	if _, err := os.Stat(goldenFile); err == nil {
		golden, err := os.ReadFile(goldenFile)
		if err != nil {
			t.Fatalf("Failed to read golden file: %v", err)
		}

		expectedContent := string(golden)
		if result != expectedContent {
			t.Errorf("Rendered output does not match golden file.\nExpected:\n%s\nGot:\n%s", expectedContent, result)
		}
	} else {
		// If golden file doesn't exist, log the output for manual verification
		t.Logf("Rendered markdown:\n%s", result)
	}
}
