package keychain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// FileStore implements the ports.CredentialStore interface using file-based storage.
// This is a fallback mechanism when OS keychain is unavailable (e.g., locked keyring on Linux).
//
// SECURITY WARNING: Credentials are stored as JSON files with 0600 permissions.
// While better than plaintext in version control, this is less secure than OS keychain.
// Use only when keychain is unavailable.
type FileStore struct {
	storePath string
	mu        sync.RWMutex
}

// NewFileStore creates a new file-based credential store.
// Credentials are stored in ~/.config/ticketr/credentials/
func NewFileStore() (*FileStore, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	storePath := filepath.Join(configDir, "ticketr", "credentials")

	// Create the credentials directory if it doesn't exist
	if err := os.MkdirAll(storePath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create credentials directory: %w", err)
	}

	return &FileStore{
		storePath: storePath,
	}, nil
}

// Store saves credentials to a file.
func (f *FileStore) Store(workspaceID string, config domain.WorkspaceConfig) (domain.CredentialRef, error) {
	if workspaceID == "" {
		return domain.CredentialRef{}, fmt.Errorf("workspace ID cannot be empty")
	}

	// Validate the configuration before storing
	if err := domain.ValidateConfig(config); err != nil {
		return domain.CredentialRef{}, fmt.Errorf("invalid configuration: %w", err)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// Create file path (using workspace ID as filename)
	filePath := filepath.Join(f.storePath, workspaceID+".json")

	// Marshal credentials to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return domain.CredentialRef{}, fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Write to file with restrictive permissions (owner read/write only)
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return domain.CredentialRef{}, fmt.Errorf("failed to write credentials file: %w", err)
	}

	// Return a reference
	return domain.CredentialRef{
		KeychainID: workspaceID,
		ServiceID:  "ticketr-file",
	}, nil
}

// Retrieve reads credentials from a file.
func (f *FileStore) Retrieve(ref domain.CredentialRef) (*domain.WorkspaceConfig, error) {
	if ref.KeychainID == "" {
		return nil, fmt.Errorf("keychain ID cannot be empty")
	}

	if ref.ServiceID != "ticketr-file" && ref.ServiceID != "ticketr" {
		return nil, fmt.Errorf("invalid service ID: expected 'ticketr-file' or 'ticketr', got %q", ref.ServiceID)
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	// Construct file path
	filePath := filepath.Join(f.storePath, ref.KeychainID+".json")

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("credentials not found for workspace ID %q", ref.KeychainID)
		}
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	// Unmarshal
	var config domain.WorkspaceConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return &config, nil
}

// Delete removes a credentials file.
func (f *FileStore) Delete(ref domain.CredentialRef) error {
	if ref.KeychainID == "" {
		return fmt.Errorf("keychain ID cannot be empty")
	}

	if ref.ServiceID != "ticketr-file" && ref.ServiceID != "ticketr" {
		return fmt.Errorf("invalid service ID: expected 'ticketr-file' or 'ticketr', got %q", ref.ServiceID)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// Construct file path
	filePath := filepath.Join(f.storePath, ref.KeychainID+".json")

	// Remove file (idempotent - no error if doesn't exist)
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete credentials file: %w", err)
	}

	return nil
}

// List returns all credential references stored in files.
func (f *FileStore) List() ([]domain.CredentialRef, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Read directory
	entries, err := os.ReadDir(f.storePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.CredentialRef{}, nil
		}
		return nil, fmt.Errorf("failed to read credentials directory: %w", err)
	}

	// Build list of references
	refs := make([]domain.CredentialRef, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Extract workspace ID from filename (remove .json extension)
		name := entry.Name()
		if filepath.Ext(name) != ".json" {
			continue
		}

		workspaceID := name[:len(name)-5] // Remove ".json"

		refs = append(refs, domain.CredentialRef{
			KeychainID: workspaceID,
			ServiceID:  "ticketr-file",
		})
	}

	return refs, nil
}
