package services

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/state"
)

// Test Case TC-303.1: TestPullService_DetectsConflictState
func TestPullService_DetectsConflictState(t *testing.T) {
	// Arrange: Create a pull_service and a StateManager
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "test.state")
	stateManager := state.NewStateManager(stateFile)
	
	// Pre-populate the state file with {"TICKET-1": {"local_hash": "A", "remote_hash": "B"}}
	stateManager.SetStoredState("TICKET-1", state.TicketState{
		LocalHash:  "A",
		RemoteHash: "B",
	})
	
	// Prepare a local Markdown file whose TICKET-1 content hashes to "C"
	localTicket := domain.Ticket{
		JiraID:      "TICKET-1",
		Title:       "Local Version",
		Description: "This is the local version that hashes to C",
	}
	// Calculate actual hash for local ticket (will be different from "A")
	localHash := stateManager.CalculateHash(localTicket)
	// This simulates that local has changed from "A" to something else
	
	// Mock repository with the local ticket
	mockRepo := &MockRepositoryForPull{
		tickets: []domain.Ticket{localTicket},
	}
	
	// Mock a Jira response for TICKET-1 that hashes to "D" (different from "B")
	remoteTicket := domain.Ticket{
		JiraID:      "TICKET-1",
		Title:       "Remote Version",
		Description: "This is the remote version that hashes to D",
	}
	
	mockJira := &MockJiraPortForPull{
		searchResult: []domain.Ticket{remoteTicket},
	}
	
	// Create the pull service
	pullService := NewPullService(mockJira, mockRepo, stateManager)
	
	// Act: Run the pull service
	result, err := pullService.Pull("test.md", PullOptions{
		ProjectKey: "TEST",
	})
	
	// Assert: The service returns a specific ErrConflictDetected error for TICKET-1
	if err == nil {
		t.Fatal("Expected conflict error, got nil")
	}
	
	if !errors.Is(err, ErrConflictDetected) {
		t.Errorf("Expected ErrConflictDetected, got: %v", err)
	}
	
	if result == nil {
		t.Fatal("Expected result even with error")
	}
	
	if len(result.Conflicts) != 1 {
		t.Errorf("Expected 1 conflict, got %d", len(result.Conflicts))
	}
	
	if len(result.Conflicts) > 0 && result.Conflicts[0] != "TICKET-1" {
		t.Errorf("Expected conflict for TICKET-1, got %v", result.Conflicts[0])
	}
	
	// Verify that the local hash was different from stored local hash "A"
	if localHash == "A" {
		t.Error("Local ticket should have a different hash than stored 'A'")
	}
	
	// Verify that remote hash was different from stored remote hash "B" 
	remoteHash := stateManager.CalculateHash(remoteTicket)
	if remoteHash == "B" {
		t.Error("Remote ticket should have a different hash than stored 'B'")
	}
}

// Mock implementations for testing
type MockRepositoryForPull struct {
	tickets      []domain.Ticket
	saveTickets  []domain.Ticket
	saveError    error
}

func (m *MockRepositoryForPull) GetTickets(filePath string) ([]domain.Ticket, error) {
	return m.tickets, nil
}

func (m *MockRepositoryForPull) SaveTickets(filePath string, tickets []domain.Ticket) error {
	m.saveTickets = tickets
	return m.saveError
}

type MockJiraPortForPull struct {
	searchResult []domain.Ticket
	searchError  error
}

func (m *MockJiraPortForPull) Authenticate() error {
	return nil
}

func (m *MockJiraPortForPull) CreateTask(task domain.Task, parentID string) (string, error) {
	return "", nil
}

func (m *MockJiraPortForPull) UpdateTask(task domain.Task) error {
	return nil
}

func (m *MockJiraPortForPull) GetProjectIssueTypes() (map[string][]string, error) {
	return nil, nil
}

func (m *MockJiraPortForPull) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *MockJiraPortForPull) CreateTicket(ticket domain.Ticket) (string, error) {
	return "", nil
}

func (m *MockJiraPortForPull) UpdateTicket(ticket domain.Ticket) error {
	return nil
}

func (m *MockJiraPortForPull) SearchTickets(projectKey string, jql string) ([]domain.Ticket, error) {
	return m.searchResult, m.searchError
}