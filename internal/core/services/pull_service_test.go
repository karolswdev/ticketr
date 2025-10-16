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

// Test Case TC-303.2: TestPullService_ConflictResolvedWithForce
func TestPullService_ConflictResolvedWithForce(t *testing.T) {
	// Arrange: Create a pull_service and a StateManager with a conflict scenario
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "test.state")
	outputFile := filepath.Join(tmpDir, "output.md")
	stateManager := state.NewStateManager(stateFile)

	// Pre-populate the state file with both local and remote hashes stored
	// This simulates a ticket that was previously synced
	stateManager.SetStoredState("TICKET-1", state.TicketState{
		LocalHash:  "originalLocalHash",
		RemoteHash: "originalRemoteHash",
	})

	// Prepare a local ticket that has changed (different from stored state)
	localTicket := domain.Ticket{
		JiraID:      "TICKET-1",
		Title:       "Local Version - Modified",
		Description: "This is the local version that has been modified locally",
	}
	localHash := stateManager.CalculateHash(localTicket)

	// Mock repository with the changed local ticket
	mockRepo := &MockRepositoryForPull{
		tickets: []domain.Ticket{localTicket},
	}

	// Mock a remote ticket that has also changed (different from stored state)
	remoteTicket := domain.Ticket{
		JiraID:      "TICKET-1",
		Title:       "Remote Version - Modified",
		Description: "This is the remote version that has been modified in JIRA",
	}
	remoteHash := stateManager.CalculateHash(remoteTicket)

	mockJira := &MockJiraPortForPull{
		searchResult: []domain.Ticket{remoteTicket},
	}

	// Verify both have actually changed from stored state
	if localHash == "originalLocalHash" {
		t.Error("Test setup error: local ticket should have different hash than stored")
	}
	if remoteHash == "originalRemoteHash" {
		t.Error("Test setup error: remote ticket should have different hash than stored")
	}

	// Create the pull service
	pullService := NewPullService(mockJira, mockRepo, stateManager)

	// Act: Run the pull service with Force: true
	result, err := pullService.Pull(outputFile, PullOptions{
		ProjectKey: "TEST",
		Force:      true,
	})

	// Assert: Should NOT return error (conflict resolved by force)
	if err != nil {
		t.Fatalf("Expected no error with Force=true, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Should have 1 conflict detected
	if len(result.Conflicts) != 1 {
		t.Errorf("Expected 1 conflict detected, got %d", len(result.Conflicts))
	}

	if len(result.Conflicts) > 0 && result.Conflicts[0] != "TICKET-1" {
		t.Errorf("Expected conflict for TICKET-1, got %v", result.Conflicts[0])
	}

	// Should have 1 ticket updated (remote overwrote local)
	if result.TicketsUpdated != 1 {
		t.Errorf("Expected 1 ticket updated, got %d", result.TicketsUpdated)
	}

	if result.TicketsSkipped != 0 {
		t.Errorf("Expected 0 tickets skipped with Force=true, got %d", result.TicketsSkipped)
	}

	// Verify the remote ticket was saved (overwrote local)
	if len(mockRepo.saveTickets) != 1 {
		t.Fatalf("Expected 1 ticket saved, got %d", len(mockRepo.saveTickets))
	}

	savedTicket := mockRepo.saveTickets[0]
	if savedTicket.Title != "Remote Version - Modified" {
		t.Errorf("Expected remote ticket title to be saved, got: %s", savedTicket.Title)
	}
	if savedTicket.Description != "This is the remote version that has been modified in JIRA" {
		t.Errorf("Expected remote ticket description to be saved, got: %s", savedTicket.Description)
	}
}