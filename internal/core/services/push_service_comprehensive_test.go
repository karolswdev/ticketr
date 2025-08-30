package services

import (
	"errors"
	"testing"

    "github.com/karolswdev/ticketr/internal/core/domain"
    "github.com/karolswdev/ticketr/internal/state"
)

// Test Case TC-304.1: TestPushService_ProcessesAllAndReportsFailures
func TestPushService_ProcessesAllAndReportsFailures(t *testing.T) {
	// Arrange: Create a Markdown file with three tickets
	ticket1 := domain.Ticket{
		Title:       "Ticket 1",
		Description: "Success",
		JiraID:      "TEST-1",
	}
	
	ticket2 := domain.Ticket{
		Title:       "Ticket 2", 
		Description: "Will fail",
		JiraID:      "TEST-2",
	}
	
	ticket3 := domain.Ticket{
		Title:       "Ticket 3",
		Description: "Success",
		JiraID:      "TEST-3",
	}
	
	tickets := []domain.Ticket{ticket1, ticket2, ticket3}
	
	// Mock repository
	mockRepo := &MockRepositoryComprehensive{
		tickets: tickets,
	}
	
	// Mock the JiraAdapter to succeed on ticket 1 and 3, but fail on ticket 2
	mockJira := &MockJiraPortComprehensive{
		failOnTicketID: "TEST-2",
	}
	
	// Create state manager
	tmpDir := t.TempDir()
	stateFile := tmpDir + "/test.state"
	stateManager := state.NewStateManager(stateFile)
	
	// Mark all tickets as changed by not having them in state
	
	// Create push service
	pushService := NewPushService(mockRepo, mockJira, stateManager)
	
	// Act: Run the push service without the --force flag
	result, err := pushService.PushTickets("test.md", ProcessOptions{
		ForcePartialUpload: false,
	})
	
	// Assert: The ProcessResult contains 2 successes and 1 failure
	if result == nil {
		t.Fatal("Expected result even with errors")
	}
	
	if result.TicketsUpdated != 2 {
		t.Errorf("Expected 2 tickets updated, got %d", result.TicketsUpdated)
	}
	
	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}
	
	// The service itself returns an error
	if err == nil {
		t.Error("Expected error to be returned when tickets fail")
	}
	
	// The mock confirms that API calls were attempted for all three tickets
	if mockJira.UpdateCallCount != 3 {
		t.Errorf("Expected UpdateTicket to be called 3 times (for all tickets), got %d", mockJira.UpdateCallCount)
	}
}

// Mock implementations for comprehensive test
type MockRepositoryComprehensive struct {
	tickets []domain.Ticket
}

func (m *MockRepositoryComprehensive) GetTickets(filePath string) ([]domain.Ticket, error) {
	return m.tickets, nil
}

func (m *MockRepositoryComprehensive) SaveTickets(filePath string, tickets []domain.Ticket) error {
	return nil
}

type MockJiraPortComprehensive struct {
	UpdateCallCount int
	failOnTicketID  string
}

func (m *MockJiraPortComprehensive) Authenticate() error {
	return nil
}

func (m *MockJiraPortComprehensive) CreateTask(task domain.Task, parentID string) (string, error) {
	return "", nil
}

func (m *MockJiraPortComprehensive) UpdateTask(task domain.Task) error {
	return nil
}

func (m *MockJiraPortComprehensive) GetProjectIssueTypes() (map[string][]string, error) {
	return nil, nil
}

func (m *MockJiraPortComprehensive) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *MockJiraPortComprehensive) CreateTicket(ticket domain.Ticket) (string, error) {
	return "", nil
}

func (m *MockJiraPortComprehensive) UpdateTicket(ticket domain.Ticket) error {
	m.UpdateCallCount++
	if ticket.JiraID == m.failOnTicketID {
		return errors.New("simulated failure for TEST-2")
	}
	return nil
}

func (m *MockJiraPortComprehensive) SearchTickets(projectKey string, jql string) ([]domain.Ticket, error) {
	return nil, nil
}
