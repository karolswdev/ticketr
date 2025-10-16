package state

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// TicketState represents the state of a ticket with bidirectional hashes
type TicketState struct {
	LocalHash  string `json:"local_hash"`
	RemoteHash string `json:"remote_hash"`
}

// StateManager manages the state file for tracking ticket changes
type StateManager struct {
	stateFilePath string
	state         map[string]TicketState // Maps ticket ID to bidirectional state
}

// NewStateManager creates a new state manager instance
func NewStateManager(stateFilePath string) *StateManager {
	if stateFilePath == "" {
		stateFilePath = ".ticketr.state"
	}
	
	return &StateManager{
		stateFilePath: stateFilePath,
		state:         make(map[string]TicketState),
	}
}

// Load reads the state file from disk
func (sm *StateManager) Load() error {
	// If state file doesn't exist, that's okay - we start with empty state
	if _, err := os.Stat(sm.stateFilePath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(sm.stateFilePath)
	if err != nil {
		return fmt.Errorf("failed to open state file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&sm.state); err != nil {
		return fmt.Errorf("failed to decode state file: %w", err)
	}

	return nil
}

// Save writes the current state to disk
func (sm *StateManager) Save() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(sm.stateFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	file, err := os.Create(sm.stateFilePath)
	if err != nil {
		return fmt.Errorf("failed to create state file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(sm.state); err != nil {
		return fmt.Errorf("failed to encode state: %w", err)
	}

	return nil
}

// CalculateHash computes the SHA256 hash of a ticket's content
func (sm *StateManager) CalculateHash(ticket domain.Ticket) string {
	h := sha256.New()
	
	// Include all relevant fields in the hash
	io.WriteString(h, ticket.Title)
	io.WriteString(h, ticket.Description)
	
	// Include acceptance criteria
	for _, ac := range ticket.AcceptanceCriteria {
		io.WriteString(h, ac)
	}
	
	// Include custom fields in a deterministic order
	// Extract and sort custom field keys for deterministic hashing
	// (Milestone 4: Go map iteration is non-deterministic)
	customFieldKeys := make([]string, 0, len(ticket.CustomFields))
	for key := range ticket.CustomFields {
		customFieldKeys = append(customFieldKeys, key)
	}
	sort.Strings(customFieldKeys)

	// Iterate in sorted order
	for _, key := range customFieldKeys {
		value := ticket.CustomFields[key]
		io.WriteString(h, key)
		io.WriteString(h, value)
	}
	
	// Include tasks
	for _, task := range ticket.Tasks {
		io.WriteString(h, task.Title)
		io.WriteString(h, task.Description)
		for _, ac := range task.AcceptanceCriteria {
			io.WriteString(h, ac)
		}

		// Extract and sort task custom field keys for deterministic hashing
		// (Milestone 4: Go map iteration is non-deterministic)
		taskCustomFieldKeys := make([]string, 0, len(task.CustomFields))
		for key := range task.CustomFields {
			taskCustomFieldKeys = append(taskCustomFieldKeys, key)
		}
		sort.Strings(taskCustomFieldKeys)

		// Iterate in sorted order
		for _, key := range taskCustomFieldKeys {
			value := task.CustomFields[key]
			io.WriteString(h, key)
			io.WriteString(h, value)
		}
	}
	
	return fmt.Sprintf("%x", h.Sum(nil))
}

// HasChanged checks if a ticket has changed since last push
func (sm *StateManager) HasChanged(ticket domain.Ticket) bool {
	if ticket.JiraID == "" {
		// New tickets always need to be pushed
		return true
	}
	
	currentHash := sm.CalculateHash(ticket)
	storedState, exists := sm.state[ticket.JiraID]
	
	// If we don't have a stored state, consider it changed
	if !exists {
		return true
	}
	
	return currentHash != storedState.LocalHash
}

// UpdateHash updates the stored hash for a ticket (updates both local and remote)
func (sm *StateManager) UpdateHash(ticket domain.Ticket) {
	if ticket.JiraID != "" {
		hash := sm.CalculateHash(ticket)
		sm.state[ticket.JiraID] = TicketState{
			LocalHash:  hash,
			RemoteHash: hash,
		}
	}
}

// UpdateLocalHash updates only the local hash for a ticket
func (sm *StateManager) UpdateLocalHash(ticket domain.Ticket) {
	if ticket.JiraID != "" {
		state := sm.state[ticket.JiraID]
		state.LocalHash = sm.CalculateHash(ticket)
		sm.state[ticket.JiraID] = state
	}
}

// UpdateRemoteHash updates only the remote hash for a ticket
func (sm *StateManager) UpdateRemoteHash(ticketID string, hash string) {
	state := sm.state[ticketID]
	state.RemoteHash = hash
	sm.state[ticketID] = state
}

// GetStoredState returns the stored state for a ticket ID
func (sm *StateManager) GetStoredState(ticketID string) (TicketState, bool) {
	state, exists := sm.state[ticketID]
	return state, exists
}

// SetStoredState sets the state for a ticket ID (useful for testing)
func (sm *StateManager) SetStoredState(ticketID string, state TicketState) {
	sm.state[ticketID] = state
}

// DetectConflict checks if there's a conflict (both local and remote changed)
func (sm *StateManager) DetectConflict(ticket domain.Ticket) bool {
	if ticket.JiraID == "" {
		return false
	}
	
	currentHash := sm.CalculateHash(ticket)
	storedState, exists := sm.state[ticket.JiraID]
	
	if !exists {
		return false
	}
	
	// Conflict occurs when both local and remote have changed
	localChanged := currentHash != storedState.LocalHash
	remoteChanged := storedState.RemoteHash != storedState.LocalHash
	
	return localChanged && remoteChanged
}

// IsRemoteChanged checks if only the remote has changed
func (sm *StateManager) IsRemoteChanged(ticketID string, remoteHash string) bool {
	storedState, exists := sm.state[ticketID]
	if !exists {
		return true // Consider new remote content as changed
	}
	return remoteHash != storedState.RemoteHash
}