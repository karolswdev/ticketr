package services

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/karolswdev/ticktr/internal/core/domain"
)

// Test Case TC-301.1: TestTicketService_RejectsLegacyStoryFormat
func TestTicketService_RejectsLegacyStoryFormat(t *testing.T) {
	// Arrange: Create a Markdown file containing the old # STORY: format
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "legacy_story.md")
	
	legacyContent := `# STORY: Old Format Story

## Description
This uses the old format

## Acceptance Criteria
- Should be rejected

## Tasks
- Old task format`
	
	err := os.WriteFile(testFile, []byte(legacyContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Create mock repository that returns error for legacy format
	mockRepo := &MockLegacyRepository{}
	mockJira := &MockJiraPortForLegacy{}
	
	// Act: Pass this file to the ticket_service
	service := NewTicketService(mockRepo, mockJira)
	result, err := service.ProcessTicketsWithOptions(testFile, ProcessOptions{})
	
	// Assert: The service returns an error and the ProcessResult indicates zero tickets were processed
	if err == nil {
		t.Error("Expected error for legacy STORY format, but got nil")
	}
	
	if result != nil && result.TicketsCreated > 0 {
		t.Errorf("Expected zero tickets processed, but got %d created", result.TicketsCreated)
	}
	
	if mockJira.createCalled {
		t.Error("JIRA adapter should not be called for legacy format")
	}
}

// MockLegacyRepository rejects legacy format
type MockLegacyRepository struct{}

func (m *MockLegacyRepository) GetTickets(filepath string) ([]domain.Ticket, error) {
	// Read the file to check format
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	
	// Check for legacy STORY format
	contentStr := string(content)
	if len(contentStr) >= 8 && contentStr[:8] == "# STORY:" {
		return nil, fmt.Errorf("legacy STORY format is no longer supported, use # TICKET: instead")
	}
	
	return []domain.Ticket{}, nil
}

func (m *MockLegacyRepository) SaveTickets(filepath string, tickets []domain.Ticket) error {
	return nil
}

// MockJiraPortForLegacy tracks if methods were called
type MockJiraPortForLegacy struct {
	createCalled bool
	updateCalled bool
}

func (m *MockJiraPortForLegacy) Authenticate() error {
	return nil
}

func (m *MockJiraPortForLegacy) CreateTask(task domain.Task, parentID string) (string, error) {
	return "TASK-123", nil
}

func (m *MockJiraPortForLegacy) UpdateTask(task domain.Task) error {
	return nil
}

func (m *MockJiraPortForLegacy) GetProjectIssueTypes() (map[string][]string, error) {
	return nil, nil
}

func (m *MockJiraPortForLegacy) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *MockJiraPortForLegacy) CreateTicket(ticket domain.Ticket) (string, error) {
	m.createCalled = true
	return "TICKET-123", nil
}

func (m *MockJiraPortForLegacy) UpdateTicket(ticket domain.Ticket) error {
	m.updateCalled = true
	return nil
}

func (m *MockJiraPortForLegacy) SearchTickets(projectKey string, jql string) ([]domain.Ticket, error) {
	return []domain.Ticket{}, nil
}

// Original test
func TestTicketService_CalculateFinalFields(t *testing.T) {
	service := NewTicketService(nil, nil)
	
	parent := domain.Ticket{
		CustomFields: map[string]string{
			"Priority": "High",
			"Sprint": "10",
		},
	}
	
	task := domain.Task{
		CustomFields: map[string]string{
			"Priority": "Low",
		},
	}
	
	result := service.calculateFinalFields(parent, task)
	
	// Assert: Priority should be overridden to "Low", Sprint should be inherited as "10"
	if result["Priority"] != "Low" {
		t.Errorf("Expected Priority to be 'Low', got '%s'", result["Priority"])
	}
	
	if result["Sprint"] != "10" {
		t.Errorf("Expected Sprint to be '10', got '%s'", result["Sprint"])
	}
}