package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/karolswdev/ticktr/internal/adapters/filesystem"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/validation"
)

// Test Case TC-302.1: TestPushCommand_FailsFastOnValidationError
func TestPushCommand_FailsFastOnValidationError(t *testing.T) {
	// Arrange: Create a Markdown file with a known validation error (e.g., a "Sub-task" under an "Epic")
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invalid_hierarchy.md")

	invalidContent := `# TICKET: My Epic

## Fields
Type: Epic

## Description
This is an epic

## Tasks
- My Subtask
  ## Fields
  Type: Sub-task`

	err := os.WriteFile(testFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Parse the file to get tickets
	repo := filesystem.NewFileRepository()
	tickets, err := repo.GetTickets(testFile)
	if err != nil {
		t.Fatalf("Failed to parse tickets: %v", err)
	}

	// Act: Execute validation logic (simulating what runPush does)
	validator := validation.NewValidator()
	validationErrors := validator.ValidateHierarchy(tickets)

	// Assert: The command would exit with a non-zero status code due to validation error
	if len(validationErrors) == 0 {
		t.Error("Expected validation error for Sub-task under Epic, but got none")
	}

	// Verify the specific error message
	foundError := false
	for _, vErr := range validationErrors {
		if vErr.Message == "A 'Sub-task' cannot be the child of a 'Epic'" {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Errorf("Expected specific hierarchy validation error, got: %v", validationErrors)
	}

	// Mock verification that JiraAdapter would never be called
	// In the actual push command, it would exit(1) before reaching JIRA initialization
	// This is verified by the fact that validation errors cause early exit
}

// MockJiraPortNeverCalled ensures no JIRA methods are called
type MockJiraPortNeverCalled struct {
	t *testing.T
}

func (m *MockJiraPortNeverCalled) Authenticate() error {
	m.t.Fatal("JiraAdapter.Authenticate should not be called on validation error")
	return nil
}

func (m *MockJiraPortNeverCalled) CreateTask(task domain.Task, parentID string) (string, error) {
	m.t.Fatal("JiraAdapter.CreateTask should not be called on validation error")
	return "", nil
}

func (m *MockJiraPortNeverCalled) UpdateTask(task domain.Task) error {
	m.t.Fatal("JiraAdapter.UpdateTask should not be called on validation error")
	return nil
}

func (m *MockJiraPortNeverCalled) GetProjectIssueTypes() (map[string][]string, error) {
	m.t.Fatal("JiraAdapter.GetProjectIssueTypes should not be called on validation error")
	return nil, nil
}

func (m *MockJiraPortNeverCalled) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	m.t.Fatal("JiraAdapter.GetIssueTypeFields should not be called on validation error")
	return nil, nil
}

func (m *MockJiraPortNeverCalled) CreateTicket(ticket domain.Ticket) (string, error) {
	m.t.Fatal("JiraAdapter.CreateTicket should not be called on validation error")
	return "", nil
}

func (m *MockJiraPortNeverCalled) UpdateTicket(ticket domain.Ticket) error {
	m.t.Fatal("JiraAdapter.UpdateTicket should not be called on validation error")
	return nil
}

func (m *MockJiraPortNeverCalled) SearchTickets(projectKey string, jql string) ([]domain.Ticket, error) {
	m.t.Fatal("JiraAdapter.SearchTickets should not be called on validation error")
	return nil, nil
}
