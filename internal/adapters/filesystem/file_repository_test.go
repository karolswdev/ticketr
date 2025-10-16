package filesystem

import (
	"os"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// Test Case TC-1.1: Parser_ParseInput_ValidFile_ReturnsCorrectStoryCount
func TestParser_ParseInput_ValidFile_ReturnsCorrectStoryCount(t *testing.T) {
	// Arrange: Create a string representing a valid Markdown input with two distinct stories
	markdownContent := `# TICKET: User Authentication
## Description
As a user, I want to be able to log in to the system.

## Fields
Type: Story
Project: PROJ

## Acceptance Criteria
- User can enter username and password
- System validates credentials

## Tasks
- Implement login form
- Add validation logic

# TICKET: User Profile Management

## Description
As a user, I want to manage my profile information.

## Fields
Type: Story  
Project: PROJ

## Acceptance Criteria
- User can view profile
- User can edit profile

## Tasks
- Create profile page
- Add edit functionality
`

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_stories_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to temp file
	if _, err := tmpFile.WriteString(markdownContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Act: Pass the file to the parser
	repo := NewFileRepository()
	tickets, err := repo.GetTickets(tmpFile.Name())

	// Assert: The parser returns a slice containing exactly two Ticket objects
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(tickets) != 2 {
		t.Errorf("Expected 2 tickets, got %d", len(tickets))
	}
}

// Test Case TC-1.2: Parser_ParseInput_TaskWithDetails_CorrectlyPopulatesTaskFields
func TestParser_ParseInput_TaskWithDetails_CorrectlyPopulatesTaskFields(t *testing.T) {
	// Arrange: Create a Markdown string for a single story with one task that has nested Description and Acceptance Criteria
	markdownContent := `# TICKET: Task with Details
## Description
Story description here.

## Fields
Type: Story
Project: PROJ

## Tasks
- Implement feature
  ## Description
  This is a detailed description of the task that needs to be implemented
  
  ## Acceptance Criteria
  - The feature should work correctly
  - All tests should pass
`

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_task_details_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to temp file
	if _, err := tmpFile.WriteString(markdownContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Act: Parse the string
	repo := NewFileRepository()
	tickets, err := repo.GetTickets(tmpFile.Name())

	// Assert: The resulting Task object has its Description and AcceptanceCriteria fields correctly populated
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}
	if len(tickets[0].Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tickets[0].Tasks))
	}

	task := tickets[0].Tasks[0]
	if task.Description != "This is a detailed description of the task that needs to be implemented" {
		t.Errorf("Expected task description to be populated, got: %q", task.Description)
	}
	if len(task.AcceptanceCriteria) != 2 {
		t.Errorf("Expected 2 acceptance criteria, got %d", len(task.AcceptanceCriteria))
	}
	if len(task.AcceptanceCriteria) > 0 && task.AcceptanceCriteria[0] != "The feature should work correctly" {
		t.Errorf("First AC incorrect: %q", task.AcceptanceCriteria[0])
	}
}

// Test Case TC-1.3: Parser_ParseInput_WithAndWithoutJiraKeys_CorrectlyPopulatesIDs
func TestParser_ParseInput_WithAndWithoutJiraKeys_CorrectlyPopulatesIDs(t *testing.T) {
	// Arrange: Create a Markdown string with stories and tasks with and without Jira keys
	markdownContent := `# TICKET: [PROJ-123] Story with Jira Key
## Description
This story has a Jira key.

## Fields
Type: Story
Project: PROJ

## Tasks
- [PROJ-124] Task with Jira key
- Task without Jira key

# TICKET: Story without Jira Key
## Description
This story has no Jira key.

## Fields
Type: Story
Project: PROJ

## Tasks
- Another task without key
`

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_jira_keys_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to temp file
	if _, err := tmpFile.WriteString(markdownContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Act: Parse the string
	repo := NewFileRepository()
	tickets, err := repo.GetTickets(tmpFile.Name())

	// Assert: The JiraID field is correctly populated for items with keys and empty for those without
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(tickets) != 2 {
		t.Fatalf("Expected 2 tickets, got %d", len(tickets))
	}

	// Check first ticket (has Jira key)
	if tickets[0].JiraID != "PROJ-123" {
		t.Errorf("Expected first ticket JiraID to be 'PROJ-123', got: %q", tickets[0].JiraID)
	}
	if tickets[0].Title != "Story with Jira Key" {
		t.Errorf("Expected first ticket title to be 'Story with Jira Key', got: %q", tickets[0].Title)
	}

	// Check first ticket's tasks
	if len(tickets[0].Tasks) != 2 {
		t.Fatalf("Expected 2 tasks in first ticket, got %d", len(tickets[0].Tasks))
	}
	if tickets[0].Tasks[0].JiraID != "PROJ-124" {
		t.Errorf("Expected first task JiraID to be 'PROJ-124', got: %q", tickets[0].Tasks[0].JiraID)
	}
	if tickets[0].Tasks[1].JiraID != "" {
		t.Errorf("Expected second task JiraID to be empty, got: %q", tickets[0].Tasks[1].JiraID)
	}

	// Check second ticket (no Jira key)
	if tickets[1].JiraID != "" {
		t.Errorf("Expected second ticket JiraID to be empty, got: %q", tickets[1].JiraID)
	}
	if tickets[1].Title != "Story without Jira Key" {
		t.Errorf("Expected second ticket title to be 'Story without Jira Key', got: %q", tickets[1].Title)
	}
}

// Test Case TC-1.4: Parser_ParseInput_MalformedStoryHeading_ReturnsErrorAndNoStories
func TestParser_ParseInput_MalformedStoryHeading_ReturnsErrorAndNoStories(t *testing.T) {
	// Arrange: Create a Markdown string where a story heading is malformed
	markdownContent := `## STORY: This is malformed (should be # not ##)
## Description
This should fail parsing.

## Tasks
- Some task
`

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_malformed_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to temp file
	if _, err := tmpFile.WriteString(markdownContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Act: Parse the string
	repo := NewFileRepository()
	tickets, err := repo.GetTickets(tmpFile.Name())

	// Assert: The parser returns no error but an empty slice (no valid tickets found)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(tickets) != 0 {
		t.Errorf("Expected 0 tickets for malformed input, got %d", len(tickets))
	}
}

// TestFileRepository_GetTickets_NonExistentFile tests error handling for missing files
func TestFileRepository_GetTickets_NonExistentFile(t *testing.T) {
	repo := NewFileRepository()

	_, err := repo.GetTickets("/nonexistent/path/to/file.md")

	if err == nil {
		t.Fatal("Expected error for nonexistent file, got nil")
	}
}

// TestFileRepository_GetTickets_EmptyFile tests handling of empty files
func TestFileRepository_GetTickets_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_empty_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	repo := NewFileRepository()
	tickets, err := repo.GetTickets(tmpFile.Name())

	if err != nil {
		t.Fatalf("Expected no error for empty file, got: %v", err)
	}
	if len(tickets) != 0 {
		t.Errorf("Expected 0 tickets for empty file, got %d", len(tickets))
	}
}

// TestFileRepository_GetTickets_InvalidMarkdown tests handling of malformed markdown
func TestFileRepository_GetTickets_InvalidMarkdown(t *testing.T) {
	invalidContent := `This is not valid markdown for tickets
Just some random text without proper structure
# Not a ticket heading
Random content here
`

	tmpFile, err := os.CreateTemp("", "test_invalid_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(invalidContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	repo := NewFileRepository()
	tickets, err := repo.GetTickets(tmpFile.Name())

	// Should not error but return empty slice
	if err != nil {
		t.Fatalf("Expected no error for invalid markdown, got: %v", err)
	}
	if len(tickets) != 0 {
		t.Errorf("Expected 0 tickets for invalid markdown, got %d", len(tickets))
	}
}

// TestFileRepository_SaveTickets_ValidTickets tests saving tickets to a file
func TestFileRepository_SaveTickets_ValidTickets(t *testing.T) {
	tmpDir := t.TempDir()
	filepath := tmpDir + "/test_tickets.md"

	repo := NewFileRepository()
	tickets := []domain.Ticket{
		{
			JiraID:      "TEST-001",
			Title:       "Test Ticket",
			Description: "This is a test ticket",
			CustomFields: map[string]string{
				"Type":    "Story",
				"Project": "TEST",
			},
			AcceptanceCriteria: []string{
				"Should save correctly",
				"Should be readable",
			},
			Tasks: []domain.Task{
				{
					JiraID:      "TEST-002",
					Title:       "Test Task",
					Description: "This is a test task",
				},
			},
		},
	}

	err := repo.SaveTickets(filepath, tickets)
	if err != nil {
		t.Fatalf("Expected no error saving tickets, got: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		t.Fatal("Expected file to be created, but it doesn't exist")
	}

	// Verify we can read it back
	readTickets, err := repo.GetTickets(filepath)
	if err != nil {
		t.Fatalf("Expected no error reading tickets, got: %v", err)
	}

	if len(readTickets) != 1 {
		t.Errorf("Expected 1 ticket, got %d", len(readTickets))
	}

	if readTickets[0].JiraID != "TEST-001" {
		t.Errorf("Expected JiraID TEST-001, got %s", readTickets[0].JiraID)
	}
}

// TestFileRepository_SaveTickets_InvalidPath tests error handling for invalid paths
func TestFileRepository_SaveTickets_InvalidPath(t *testing.T) {
	repo := NewFileRepository()
	tickets := []domain.Ticket{
		{
			Title: "Test",
		},
	}

	// Try to save to invalid path
	err := repo.SaveTickets("/invalid/nonexistent/path/file.md", tickets)

	if err == nil {
		t.Fatal("Expected error for invalid path, got nil")
	}
}

// TestFileRepository_SaveTickets_EmptyTickets tests saving an empty ticket list
func TestFileRepository_SaveTickets_EmptyTickets(t *testing.T) {
	tmpDir := t.TempDir()
	filepath := tmpDir + "/empty_tickets.md"

	repo := NewFileRepository()
	tickets := []domain.Ticket{}

	err := repo.SaveTickets(filepath, tickets)
	if err != nil {
		t.Fatalf("Expected no error saving empty tickets, got: %v", err)
	}

	// Verify file exists but is empty
	content, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Expected to read file, got error: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty file, got %d bytes", len(content))
	}
}

// TestFileRepository_SaveTickets_RoundTrip tests save and load consistency
func TestFileRepository_SaveTickets_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	filepath := tmpDir + "/roundtrip.md"

	repo := NewFileRepository()

	original := []domain.Ticket{
		{
			JiraID:      "RT-100",
			Title:       "Round Trip Test",
			Description: "Testing save and load consistency",
			CustomFields: map[string]string{
				"Type":     "Bug",
				"Priority": "High",
				"Project":  "RT",
			},
			AcceptanceCriteria: []string{
				"Data persists correctly",
				"All fields are preserved",
			},
			Tasks: []domain.Task{
				{
					JiraID:      "RT-101",
					Title:       "Subtask 1",
					Description: "First subtask",
					CustomFields: map[string]string{
						"Status": "Done",
					},
					AcceptanceCriteria: []string{
						"Subtask AC 1",
					},
				},
			},
		},
	}

	// Save tickets
	err := repo.SaveTickets(filepath, original)
	if err != nil {
		t.Fatalf("Failed to save tickets: %v", err)
	}

	// Load tickets back
	loaded, err := repo.GetTickets(filepath)
	if err != nil {
		t.Fatalf("Failed to load tickets: %v", err)
	}

	// Verify structure
	if len(loaded) != len(original) {
		t.Fatalf("Expected %d tickets, got %d", len(original), len(loaded))
	}

	// Verify ticket data
	if loaded[0].JiraID != original[0].JiraID {
		t.Errorf("JiraID mismatch: expected %s, got %s", original[0].JiraID, loaded[0].JiraID)
	}

	if loaded[0].Title != original[0].Title {
		t.Errorf("Title mismatch: expected %s, got %s", original[0].Title, loaded[0].Title)
	}

	if loaded[0].Description != original[0].Description {
		t.Errorf("Description mismatch")
	}

	if len(loaded[0].Tasks) != len(original[0].Tasks) {
		t.Errorf("Tasks count mismatch: expected %d, got %d", len(original[0].Tasks), len(loaded[0].Tasks))
	}

	if len(loaded[0].AcceptanceCriteria) != len(original[0].AcceptanceCriteria) {
		t.Errorf("AC count mismatch: expected %d, got %d", len(original[0].AcceptanceCriteria), len(loaded[0].AcceptanceCriteria))
	}
}

// TestFileRepository_GetTickets_PermissionDenied tests handling of permission errors
func TestFileRepository_GetTickets_PermissionDenied(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()
	filepath := tmpDir + "/noperm.md"

	// Create a file and remove read permissions
	if err := os.WriteFile(filepath, []byte("test"), 0000); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	repo := NewFileRepository()
	_, err := repo.GetTickets(filepath)

	if err == nil {
		t.Error("Expected error for permission denied, got nil")
	}
}

// TestFileRepository_SaveTickets_MultipleTickets tests saving multiple tickets with various configurations
func TestFileRepository_SaveTickets_MultipleTickets(t *testing.T) {
	tmpDir := t.TempDir()
	filepath := tmpDir + "/multiple.md"

	repo := NewFileRepository()

	tickets := []domain.Ticket{
		{
			JiraID:      "MULT-1",
			Title:       "First Ticket",
			Description: "First description",
		},
		{
			JiraID:      "MULT-2",
			Title:       "Second Ticket",
			Description: "Second description",
			Tasks: []domain.Task{
				{Title: "Task 1"},
				{Title: "Task 2"},
			},
		},
		{
			Title:       "Third Ticket without Jira ID",
			Description: "Third description",
		},
	}

	err := repo.SaveTickets(filepath, tickets)
	if err != nil {
		t.Fatalf("Failed to save multiple tickets: %v", err)
	}

	loaded, err := repo.GetTickets(filepath)
	if err != nil {
		t.Fatalf("Failed to load multiple tickets: %v", err)
	}

	if len(loaded) != 3 {
		t.Errorf("Expected 3 tickets, got %d", len(loaded))
	}

	// Verify second ticket has tasks
	if len(loaded[1].Tasks) != 2 {
		t.Errorf("Expected 2 tasks in second ticket, got %d", len(loaded[1].Tasks))
	}

	// Verify third ticket has no Jira ID
	if loaded[2].JiraID != "" {
		t.Errorf("Expected third ticket to have no Jira ID, got %s", loaded[2].JiraID)
	}
}
