package state

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// StateManager manages the state file for tracking ticket changes
type StateManager struct {
	stateFilePath string
	state         map[string]string // Maps ticket ID to content hash
}

// NewStateManager creates a new state manager instance
func NewStateManager(stateFilePath string) *StateManager {
	if stateFilePath == "" {
		stateFilePath = ".ticketr.state"
	}
	
	return &StateManager{
		stateFilePath: stateFilePath,
		state:         make(map[string]string),
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
	for k, v := range ticket.CustomFields {
		io.WriteString(h, k)
		io.WriteString(h, v)
	}
	
	// Include tasks
	for _, task := range ticket.Tasks {
		io.WriteString(h, task.Title)
		io.WriteString(h, task.Description)
		for _, ac := range task.AcceptanceCriteria {
			io.WriteString(h, ac)
		}
		for k, v := range task.CustomFields {
			io.WriteString(h, k)
			io.WriteString(h, v)
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
	storedHash, exists := sm.state[ticket.JiraID]
	
	// If we don't have a stored hash, consider it changed
	if !exists {
		return true
	}
	
	return currentHash != storedHash
}

// UpdateHash updates the stored hash for a ticket
func (sm *StateManager) UpdateHash(ticket domain.Ticket) {
	if ticket.JiraID != "" {
		sm.state[ticket.JiraID] = sm.CalculateHash(ticket)
	}
}

// GetStoredHash returns the stored hash for a ticket ID
func (sm *StateManager) GetStoredHash(ticketID string) (string, bool) {
	hash, exists := sm.state[ticketID]
	return hash, exists
}

// SetStoredHash sets the hash for a ticket ID (useful for testing)
func (sm *StateManager) SetStoredHash(ticketID string, hash string) {
	sm.state[ticketID] = hash
}