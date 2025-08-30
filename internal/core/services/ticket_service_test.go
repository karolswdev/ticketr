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

// Test Case TC-302.1: TestTicketService_DryRunMode
func TestTicketService_DryRunMode(t *testing.T) {
	// Arrange: Create mock repository with test tickets
	mockRepo := &MockDryRunRepository{
		tickets: []domain.Ticket{
			{
				Title:       "New Ticket",
				Description: "This is a new ticket",
				JiraID:      "", // No JIRA ID - should be created
				Tasks: []domain.Task{
					{
						Title:       "New Task",
						Description: "This is a new task",
						JiraID:      "", // No JIRA ID - should be created
					},
				},
			},
			{
				Title:       "Existing Ticket",
				Description: "This ticket exists in JIRA",
				JiraID:      "PROJ-100", // Has JIRA ID - should be updated
				Tasks: []domain.Task{
					{
						Title:       "Existing Task",
						Description: "This task exists in JIRA",
						JiraID:      "PROJ-101", // Has JIRA ID - should be updated
					},
					{
						Title:       "New Task in Existing Ticket",
						Description: "This is a new task",
						JiraID:      "", // No JIRA ID - should be created
					},
				},
			},
		},
	}
	
	mockJira := &MockDryRunJiraPort{}
	service := NewTicketService(mockRepo, mockJira)
	
	// Act: Process tickets with DryRun enabled
	result, err := service.ProcessTicketsWithOptions("test.md", ProcessOptions{
		DryRun: true,
	})
	
	// Assert: No errors and correct counts
	if err != nil {
		t.Errorf("Expected no error in dry-run mode, got: %v", err)
	}
	
	if result.TicketsCreated != 1 {
		t.Errorf("Expected 1 ticket to be marked for creation, got %d", result.TicketsCreated)
	}
	
	if result.TicketsUpdated != 1 {
		t.Errorf("Expected 1 ticket to be marked for update, got %d", result.TicketsUpdated)
	}
	
	if result.TasksCreated != 2 {
		t.Errorf("Expected 2 tasks to be marked for creation, got %d", result.TasksCreated)
	}
	
	if result.TasksUpdated != 1 {
		t.Errorf("Expected 1 task to be marked for update, got %d", result.TasksUpdated)
	}
	
	// Assert: JIRA methods should NOT be called in dry-run mode
	if mockJira.createTicketCalled {
		t.Error("CreateTicket should not be called in dry-run mode")
	}
	
	if mockJira.updateTicketCalled {
		t.Error("UpdateTicket should not be called in dry-run mode")
	}
	
	if mockJira.createTaskCalled {
		t.Error("CreateTask should not be called in dry-run mode")
	}
	
	if mockJira.updateTaskCalled {
		t.Error("UpdateTask should not be called in dry-run mode")
	}
	
	// Assert: Repository save should NOT be called in dry-run mode
	if mockRepo.saveCalled {
		t.Error("SaveTickets should not be called in dry-run mode")
	}
}

// Test Case TC-302.2: TestTicketService_NormalModeAfterDryRun
func TestTicketService_NormalModeAfterDryRun(t *testing.T) {
	// This test ensures that normal mode works correctly (actually calls JIRA)
	mockRepo := &MockDryRunRepository{
		tickets: []domain.Ticket{
			{
				Title:       "New Ticket",
				Description: "This is a new ticket",
				JiraID:      "",
			},
		},
	}
	
	mockJira := &MockDryRunJiraPort{}
	service := NewTicketService(mockRepo, mockJira)
	
	// Act: Process tickets WITHOUT DryRun (normal mode)
	result, err := service.ProcessTicketsWithOptions("test.md", ProcessOptions{
		DryRun: false,
	})
	
	// Assert: No errors and JIRA should be called
	if err != nil {
		t.Errorf("Expected no error in normal mode, got: %v", err)
	}
	
	if result.TicketsCreated != 1 {
		t.Errorf("Expected 1 ticket to be created, got %d", result.TicketsCreated)
	}
	
	// Assert: JIRA methods SHOULD be called in normal mode
	if !mockJira.createTicketCalled {
		t.Error("CreateTicket should be called in normal mode")
	}
	
	// Assert: Repository save SHOULD be called in normal mode
	if !mockRepo.saveCalled {
		t.Error("SaveTickets should be called in normal mode")
	}
}

// MockDryRunRepository for testing dry-run functionality
type MockDryRunRepository struct {
	tickets    []domain.Ticket
	saveCalled bool
}

func (m *MockDryRunRepository) GetTickets(filepath string) ([]domain.Ticket, error) {
	return m.tickets, nil
}

func (m *MockDryRunRepository) SaveTickets(filepath string, tickets []domain.Ticket) error {
	m.saveCalled = true
	return nil
}

// MockDryRunJiraPort tracks which methods were called
type MockDryRunJiraPort struct {
	createTicketCalled bool
	updateTicketCalled bool
	createTaskCalled   bool
	updateTaskCalled   bool
}

func (m *MockDryRunJiraPort) Authenticate() error {
	return nil
}

func (m *MockDryRunJiraPort) CreateTicket(ticket domain.Ticket) (string, error) {
	m.createTicketCalled = true
	return "PROJ-200", nil
}

func (m *MockDryRunJiraPort) UpdateTicket(ticket domain.Ticket) error {
	m.updateTicketCalled = true
	return nil
}

func (m *MockDryRunJiraPort) CreateTask(task domain.Task, parentID string) (string, error) {
	m.createTaskCalled = true
	return "PROJ-201", nil
}

func (m *MockDryRunJiraPort) UpdateTask(task domain.Task) error {
	m.updateTaskCalled = true
	return nil
}

func (m *MockDryRunJiraPort) GetProjectIssueTypes() (map[string][]string, error) {
	return nil, nil
}

func (m *MockDryRunJiraPort) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *MockDryRunJiraPort) SearchTickets(projectKey string, jql string) ([]domain.Ticket, error) {
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