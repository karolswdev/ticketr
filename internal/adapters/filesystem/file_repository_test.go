package filesystem

import (
	"os"
	"testing"
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
	stories, err := repo.GetStories(tmpFile.Name())

	// Assert: The parser returns a slice containing exactly two Story objects
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(stories) != 2 {
		t.Errorf("Expected 2 stories, got %d", len(stories))
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
	stories, err := repo.GetStories(tmpFile.Name())

	// Assert: The resulting Task object has its Description and AcceptanceCriteria fields correctly populated
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(stories) != 1 {
		t.Fatalf("Expected 1 story, got %d", len(stories))
	}
	if len(stories[0].Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(stories[0].Tasks))
	}

	task := stories[0].Tasks[0]
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
	stories, err := repo.GetStories(tmpFile.Name())

	// Assert: The JiraID field is correctly populated for items with keys and empty for those without
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(stories) != 2 {
		t.Fatalf("Expected 2 stories, got %d", len(stories))
	}

	// Check first story (has Jira key)
	if stories[0].JiraID != "PROJ-123" {
		t.Errorf("Expected first story JiraID to be 'PROJ-123', got: %q", stories[0].JiraID)
	}
	if stories[0].Title != "Story with Jira Key" {
		t.Errorf("Expected first story title to be 'Story with Jira Key', got: %q", stories[0].Title)
	}

	// Check first story's tasks
	if len(stories[0].Tasks) != 2 {
		t.Fatalf("Expected 2 tasks in first story, got %d", len(stories[0].Tasks))
	}
	if stories[0].Tasks[0].JiraID != "PROJ-124" {
		t.Errorf("Expected first task JiraID to be 'PROJ-124', got: %q", stories[0].Tasks[0].JiraID)
	}
	if stories[0].Tasks[1].JiraID != "" {
		t.Errorf("Expected second task JiraID to be empty, got: %q", stories[0].Tasks[1].JiraID)
	}

	// Check second story (no Jira key)
	if stories[1].JiraID != "" {
		t.Errorf("Expected second story JiraID to be empty, got: %q", stories[1].JiraID)
	}
	if stories[1].Title != "Story without Jira Key" {
		t.Errorf("Expected second story title to be 'Story without Jira Key', got: %q", stories[1].Title)
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
	stories, err := repo.GetStories(tmpFile.Name())

	// Assert: The parser returns no error but an empty slice (no valid tickets found)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(stories) != 0 {
		t.Errorf("Expected 0 stories for malformed input, got %d", len(stories))
	}
}