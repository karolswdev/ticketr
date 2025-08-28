package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

func TestStateManager_BidirectionalHashTracking(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "test.state")
	
	// Create a state manager
	sm := NewStateManager(stateFile)
	
	// Create a test ticket
	ticket := domain.Ticket{
		JiraID:      "TEST-123",
		Title:       "Test Ticket",
		Description: "Original description",
	}
	
	// Test 1: New ticket should have no stored state
	state, exists := sm.GetStoredState(ticket.JiraID)
	if exists {
		t.Error("Expected no stored state for new ticket")
	}
	
	// Test 2: UpdateHash should set both local and remote to same hash
	sm.UpdateHash(ticket)
	state, exists = sm.GetStoredState(ticket.JiraID)
	if !exists {
		t.Fatal("Expected stored state after UpdateHash")
	}
	if state.LocalHash != state.RemoteHash {
		t.Error("Expected LocalHash and RemoteHash to be equal after UpdateHash")
	}
	if state.LocalHash == "" {
		t.Error("Expected non-empty hash")
	}
	
	// Test 3: Save and Load persistence
	err := sm.Save()
	if err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}
	
	// Create new manager and load
	sm2 := NewStateManager(stateFile)
	err = sm2.Load()
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}
	
	state2, exists := sm2.GetStoredState(ticket.JiraID)
	if !exists {
		t.Fatal("Expected persisted state after load")
	}
	if state2.LocalHash != state.LocalHash || state2.RemoteHash != state.RemoteHash {
		t.Error("State not properly persisted")
	}
	
	// Test 4: UpdateLocalHash only updates local
	ticket.Description = "Modified locally"
	originalRemote := state.RemoteHash
	sm2.UpdateLocalHash(ticket)
	state3, _ := sm2.GetStoredState(ticket.JiraID)
	if state3.LocalHash == originalRemote {
		t.Error("Expected LocalHash to change after UpdateLocalHash")
	}
	if state3.RemoteHash != originalRemote {
		t.Error("RemoteHash should not change with UpdateLocalHash")
	}
	
	// Test 5: UpdateRemoteHash only updates remote
	newRemoteHash := "remote-hash-from-jira"
	sm2.UpdateRemoteHash(ticket.JiraID, newRemoteHash)
	state4, _ := sm2.GetStoredState(ticket.JiraID)
	if state4.RemoteHash != newRemoteHash {
		t.Error("Expected RemoteHash to be updated")
	}
	if state4.LocalHash != state3.LocalHash {
		t.Error("LocalHash should not change with UpdateRemoteHash")
	}
}

func TestStateManager_ConflictDetection(t *testing.T) {
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "test.state")
	sm := NewStateManager(stateFile)
	
	ticket := domain.Ticket{
		JiraID:      "TEST-456",
		Title:       "Conflict Test",
		Description: "Original",
	}
	
	// Initially no conflict (no stored state)
	if sm.DetectConflict(ticket) {
		t.Error("Should not detect conflict for new ticket")
	}
	
	// Set initial state (synced)
	sm.UpdateHash(ticket)
	if sm.DetectConflict(ticket) {
		t.Error("Should not detect conflict when synced")
	}
	
	// Simulate remote change only
	sm.UpdateRemoteHash(ticket.JiraID, "different-remote-hash")
	if sm.DetectConflict(ticket) {
		t.Error("Should not detect conflict when only remote changed")
	}
	
	// Now change local too - this creates a conflict
	ticket.Description = "Modified locally"
	if !sm.DetectConflict(ticket) {
		t.Error("Should detect conflict when both local and remote changed")
	}
	
	// Test IsRemoteChanged
	if !sm.IsRemoteChanged(ticket.JiraID, "new-remote-hash") {
		t.Error("Should detect remote change")
	}
	
	state, _ := sm.GetStoredState(ticket.JiraID)
	if sm.IsRemoteChanged(ticket.JiraID, state.RemoteHash) {
		t.Error("Should not detect change when hash matches")
	}
}

func TestStateManager_BackwardCompatibility(t *testing.T) {
	// Test that old state files with simple string hashes can still be loaded
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "old.state")
	
	// Create an old-format state file
	oldState := `{"TEST-789": "simple-hash-string"}`
	err := os.WriteFile(stateFile, []byte(oldState), 0644)
	if err != nil {
		t.Fatalf("Failed to write old state file: %v", err)
	}
	
	// Try to load with new manager - should handle gracefully
	sm := NewStateManager(stateFile)
	err = sm.Load()
	// We expect this to fail with the new structure, which is acceptable
	// as we're making a breaking change in v2.0
	if err == nil {
		t.Error("Expected error when loading old format, as v2.0 is a breaking change")
	}
}