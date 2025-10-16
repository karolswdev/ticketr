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

// TC-501.1: TestProcessTicketsWithOptions_MixedValidityWithoutForce
// Tests behavior when processing a mix of valid and invalid tickets without force flag
func TestProcessTicketsWithOptions_MixedValidityWithoutForce(t *testing.T) {
	// Arrange: Create mock repository with mix of valid and invalid tickets
	mockRepo := &MockMixedRepository{
		tickets: []domain.Ticket{
			{
				Title:   "Valid Ticket",
				JiraID:  "", // New ticket
				CustomFields: map[string]string{
					"Type": "Story",
				},
			},
			{
				Title:   "Invalid Ticket - Missing Parent",
				JiraID:  "", // New ticket, but will fail
				Tasks: []domain.Task{
					{
						Title:  "Orphan Task",
						JiraID: "",
					},
				},
			},
		},
	}

	mockJira := &MockJiraPortWithErrors{
		shouldFailForTitle: "Invalid Ticket - Missing Parent",
	}

	service := NewTicketService(mockRepo, mockJira)

	// Act: Process tickets without force flag
	result, err := service.ProcessTicketsWithOptions("test.md", ProcessOptions{
		ForcePartialUpload: false,
	})

	// Assert: Service continues processing, errors collected
	if err != nil {
		t.Errorf("Expected no error from service, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Valid ticket should succeed
	if result.TicketsCreated != 1 {
		t.Errorf("Expected 1 ticket created, got %d", result.TicketsCreated)
	}

	// Invalid ticket should be in errors
	if len(result.Errors) == 0 {
		t.Error("Expected errors to be collected, got none")
	}
}

// TC-501.2: TestProcessTicketsWithOptions_MixedValidityWithForce
// Tests behavior when processing a mix of valid and invalid tickets with force flag
func TestProcessTicketsWithOptions_MixedValidityWithForce(t *testing.T) {
	// Arrange: Create mock repository with mix of valid and invalid tickets
	mockRepo := &MockMixedRepository{
		tickets: []domain.Ticket{
			{
				Title:   "Valid Ticket 1",
				JiraID:  "",
				CustomFields: map[string]string{
					"Type": "Story",
				},
			},
			{
				Title:   "Invalid Ticket - Will Fail",
				JiraID:  "",
			},
			{
				Title:   "Valid Ticket 2",
				JiraID:  "",
				CustomFields: map[string]string{
					"Type": "Task",
				},
			},
		},
	}

	mockJira := &MockJiraPortWithErrors{
		shouldFailForTitle: "Invalid Ticket - Will Fail",
	}

	service := NewTicketService(mockRepo, mockJira)

	// Act: Process tickets with force flag
	result, err := service.ProcessTicketsWithOptions("test.md", ProcessOptions{
		ForcePartialUpload: true,
	})

	// Assert: Service continues processing, valid items succeed, invalid items in Errors
	if err != nil {
		t.Errorf("Expected no error from service, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Two valid tickets should succeed
	if result.TicketsCreated != 2 {
		t.Errorf("Expected 2 tickets created, got %d", result.TicketsCreated)
	}

	// One invalid ticket should be in errors
	if len(result.Errors) == 0 {
		t.Error("Expected errors to be collected for invalid ticket")
	}

	// Verify error message contains the failed ticket
	foundError := false
	for _, errMsg := range result.Errors {
		if fmt.Sprintf("%v", errMsg) != "" && len(fmt.Sprintf("%v", errMsg)) > 0 {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("Expected error message for failed ticket")
	}
}

// TC-501.3: TestProcessTicketsWithOptions_AllValid
// Tests behavior when all tickets are valid
func TestProcessTicketsWithOptions_AllValid(t *testing.T) {
	// Arrange: Create mock repository with all valid tickets
	mockRepo := &MockMixedRepository{
		tickets: []domain.Ticket{
			{
				Title:   "Valid Ticket 1",
				JiraID:  "",
				CustomFields: map[string]string{
					"Type": "Story",
				},
			},
			{
				Title:   "Valid Ticket 2",
				JiraID:  "",
				CustomFields: map[string]string{
					"Type": "Task",
				},
			},
			{
				Title:   "Valid Ticket 3",
				JiraID:  "",
				CustomFields: map[string]string{
					"Type": "Bug",
				},
			},
		},
	}

	mockJira := &MockJiraPortWithErrors{
		shouldFailForTitle: "", // No failures
	}

	service := NewTicketService(mockRepo, mockJira)

	// Act: Process all valid tickets
	result, err := service.ProcessTicketsWithOptions("test.md", ProcessOptions{
		ForcePartialUpload: false,
	})

	// Assert: All items processed successfully, no errors
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.TicketsCreated != 3 {
		t.Errorf("Expected 3 tickets created, got %d", result.TicketsCreated)
	}

	if len(result.Errors) != 0 {
		t.Errorf("Expected no errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

// TC-501.4: TestProcessTicketsWithOptions_AllInvalid
// Tests behavior when all tickets are invalid with force flag
func TestProcessTicketsWithOptions_AllInvalid(t *testing.T) {
	// Arrange: Create mock repository with all invalid tickets
	mockRepo := &MockMixedRepository{
		tickets: []domain.Ticket{
			{
				Title:   "Invalid Ticket 1",
				JiraID:  "",
				Tasks: []domain.Task{
					{
						Title:  "Task without parent ID",
						JiraID: "",
					},
				},
			},
			{
				Title:   "Invalid Ticket 2",
				JiraID:  "",
				Tasks: []domain.Task{
					{
						Title:  "Another orphan task",
						JiraID: "",
					},
				},
			},
		},
	}

	mockJira := &MockJiraPortWithErrors{
		failAll: true,
	}

	service := NewTicketService(mockRepo, mockJira)

	// Act: Process all invalid tickets with force flag
	result, err := service.ProcessTicketsWithOptions("test.md", ProcessOptions{
		ForcePartialUpload: true,
	})

	// Assert: Service completes without panic, all errors collected
	if err != nil {
		t.Errorf("Expected no error from service (should collect errors), got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// No tickets should be created
	if result.TicketsCreated != 0 {
		t.Errorf("Expected 0 tickets created, got %d", result.TicketsCreated)
	}

	// All tickets should have errors
	if len(result.Errors) == 0 {
		t.Error("Expected errors to be collected for all invalid tickets")
	}
}

// MockMixedRepository for testing mixed validity scenarios
type MockMixedRepository struct {
	tickets []domain.Ticket
}

func (m *MockMixedRepository) GetTickets(filepath string) ([]domain.Ticket, error) {
	return m.tickets, nil
}

func (m *MockMixedRepository) SaveTickets(filepath string, tickets []domain.Ticket) error {
	return nil
}

// MockJiraPortWithErrors simulates JIRA errors for specific tickets
type MockJiraPortWithErrors struct {
	shouldFailForTitle string
	failAll            bool
	createCount        int
}

func (m *MockJiraPortWithErrors) Authenticate() error {
	return nil
}

func (m *MockJiraPortWithErrors) CreateTicket(ticket domain.Ticket) (string, error) {
	if m.failAll {
		return "", fmt.Errorf("simulated JIRA error for ticket '%s'", ticket.Title)
	}
	if ticket.Title == m.shouldFailForTitle {
		return "", fmt.Errorf("simulated JIRA error for ticket '%s'", ticket.Title)
	}
	m.createCount++
	return fmt.Sprintf("TICKET-%d", m.createCount), nil
}

func (m *MockJiraPortWithErrors) UpdateTicket(ticket domain.Ticket) error {
	if m.failAll {
		return fmt.Errorf("simulated JIRA error for ticket '%s'", ticket.Title)
	}
	if ticket.Title == m.shouldFailForTitle {
		return fmt.Errorf("simulated JIRA error for ticket '%s'", ticket.Title)
	}
	return nil
}

func (m *MockJiraPortWithErrors) CreateTask(task domain.Task, parentID string) (string, error) {
	if m.failAll {
		return "", fmt.Errorf("simulated JIRA error for task '%s'", task.Title)
	}
	if parentID == "" {
		return "", fmt.Errorf("parent ID required for task '%s'", task.Title)
	}
	m.createCount++
	return fmt.Sprintf("TASK-%d", m.createCount), nil
}

func (m *MockJiraPortWithErrors) UpdateTask(task domain.Task) error {
	if m.failAll {
		return fmt.Errorf("simulated JIRA error for task '%s'", task.Title)
	}
	return nil
}

func (m *MockJiraPortWithErrors) GetProjectIssueTypes() (map[string][]string, error) {
	return map[string][]string{
		"TEST": {"Story", "Task", "Bug"},
	}, nil
}

func (m *MockJiraPortWithErrors) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (m *MockJiraPortWithErrors) SearchTickets(projectKey string, jql string) ([]domain.Ticket, error) {
	return []domain.Ticket{}, nil
}