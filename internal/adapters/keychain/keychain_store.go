package keychain

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/zalando/go-keyring"
)

const (
	// serviceName is the identifier used for storing credentials in the OS keychain.
	serviceName = "ticketr"

	// listKeychainID is a special keychain ID used to store the list of all credential references.
	listKeychainID = "ticketr-credential-list"
)

// KeychainStore implements the ports.CredentialStore interface using OS keychain integration.
// It provides secure credential storage using:
// - macOS Keychain
// - Windows Credential Manager
// - Linux Secret Service API
//
// Thread-safety is provided by the underlying go-keyring library, which uses
// OS-level synchronization mechanisms.
type KeychainStore struct {
	serviceName string
	mu          sync.RWMutex
}

// NewKeychainStore creates a new instance of KeychainStore.
// The serviceName parameter identifies this application in the OS keychain.
func NewKeychainStore() *KeychainStore {
	return &KeychainStore{
		serviceName: serviceName,
	}
}

// Store saves credentials securely in the OS keychain and returns a reference.
// The workspaceID is used as the keychain account identifier.
//
// Security: Credentials are stored using OS-level encryption and never logged.
func (k *KeychainStore) Store(workspaceID string, config domain.WorkspaceConfig) (domain.CredentialRef, error) {
	if workspaceID == "" {
		return domain.CredentialRef{}, fmt.Errorf("workspace ID cannot be empty")
	}

	// Validate the configuration before storing
	if err := domain.ValidateConfig(config); err != nil {
		return domain.CredentialRef{}, fmt.Errorf("invalid configuration: %w", err)
	}

	// Marshal the config to JSON for storage
	credentialJSON, err := json.Marshal(config)
	if err != nil {
		return domain.CredentialRef{}, fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Store in the OS keychain
	// Note: go-keyring handles OS-level encryption automatically
	err = keyring.Set(k.serviceName, workspaceID, string(credentialJSON))
	if err != nil {
		return domain.CredentialRef{}, fmt.Errorf("failed to store credentials in keychain: %w", err)
	}

	// Update the credential list
	if err := k.addToList(workspaceID); err != nil {
		// If adding to list fails, try to clean up the credential we just stored
		_ = keyring.Delete(k.serviceName, workspaceID)
		return domain.CredentialRef{}, fmt.Errorf("failed to update credential list: %w", err)
	}

	// Return a reference to the stored credentials
	ref := domain.CredentialRef{
		KeychainID: workspaceID,
		ServiceID:  k.serviceName,
	}

	return ref, nil
}

// Retrieve fetches credentials from the OS keychain using a reference.
// Returns an error if the credentials are not found or cannot be decrypted.
func (k *KeychainStore) Retrieve(ref domain.CredentialRef) (*domain.WorkspaceConfig, error) {
	if ref.KeychainID == "" {
		return nil, fmt.Errorf("keychain ID cannot be empty")
	}

	if ref.ServiceID != k.serviceName {
		return nil, fmt.Errorf("invalid service ID: expected %q, got %q", k.serviceName, ref.ServiceID)
	}

	// Retrieve from the OS keychain
	credentialJSON, err := keyring.Get(k.serviceName, ref.KeychainID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("credentials not found for workspace ID %q", ref.KeychainID)
		}
		return nil, fmt.Errorf("failed to retrieve credentials from keychain: %w", err)
	}

	// Unmarshal the JSON
	var config domain.WorkspaceConfig
	if err := json.Unmarshal([]byte(credentialJSON), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return &config, nil
}

// Delete removes credentials from the OS keychain.
// This operation is idempotent - deleting non-existent credentials returns no error.
func (k *KeychainStore) Delete(ref domain.CredentialRef) error {
	if ref.KeychainID == "" {
		return fmt.Errorf("keychain ID cannot be empty")
	}

	if ref.ServiceID != k.serviceName {
		return fmt.Errorf("invalid service ID: expected %q, got %q", k.serviceName, ref.ServiceID)
	}

	// Delete from the OS keychain
	err := keyring.Delete(k.serviceName, ref.KeychainID)
	if err != nil && !isNotFoundError(err) {
		return fmt.Errorf("failed to delete credentials from keychain: %w", err)
	}

	// Remove from the credential list (ignore errors as the credential is already deleted)
	_ = k.removeFromList(ref.KeychainID)

	return nil
}

// List returns all credential references stored in the keychain.
// This is useful for auditing and cleanup operations.
func (k *KeychainStore) List() ([]domain.CredentialRef, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	credentialIDs, err := k.loadCredentialIDsLocked()
	if err != nil {
		return nil, err
	}

	// Build the list of credential references
	refs := make([]domain.CredentialRef, 0, len(credentialIDs))
	for _, id := range credentialIDs {
		refs = append(refs, domain.CredentialRef{
			KeychainID: id,
			ServiceID:  k.serviceName,
		})
	}

	return refs, nil
}

// addToList adds a credential ID to the internal list of stored credentials.
func (k *KeychainStore) addToList(credentialID string) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	// Get the current list
	ids, err := k.loadCredentialIDsLocked()
	if err != nil {
		return err
	}

	// Check if the ID already exists
	for _, id := range ids {
		if id == credentialID {
			// Already in the list, no need to add
			return nil
		}
	}

	// Add the new ID
	ids = append(ids, credentialID)

	return k.storeCredentialIDsLocked(ids)
}

// removeFromList removes a credential ID from the internal list of stored credentials.
func (k *KeychainStore) removeFromList(credentialID string) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	// Get the current list
	ids, err := k.loadCredentialIDsLocked()
	if err != nil {
		return err
	}

	// Filter out the ID to remove
	filtered := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != credentialID {
			filtered = append(filtered, id)
		}
	}

	// If the list is now empty, delete the list entry entirely
	if len(filtered) == 0 {
		err = keyring.Delete(k.serviceName, listKeychainID)
		if err != nil && !isNotFoundError(err) {
			return fmt.Errorf("failed to delete credential list: %w", err)
		}
		return nil
	}

	return k.storeCredentialIDsLocked(filtered)
}

func (k *KeychainStore) loadCredentialIDsLocked() ([]string, error) {
	// Retrieve the list of credential IDs from the special list entry
	listJSON, err := keyring.Get(k.serviceName, listKeychainID)
	if err != nil {
		if isNotFoundError(err) {
			// No credentials stored yet
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to retrieve credential list from keychain: %w", err)
	}

	// Unmarshal the list
	var credentialIDs []string
	if err := json.Unmarshal([]byte(listJSON), &credentialIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credential list: %w", err)
	}

	return credentialIDs, nil
}

func (k *KeychainStore) storeCredentialIDsLocked(ids []string) error {
	// Marshal and store the updated list
	listJSON, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("failed to marshal credential list: %w", err)
	}

	err = keyring.Set(k.serviceName, listKeychainID, string(listJSON))
	if err != nil {
		return fmt.Errorf("failed to store credential list: %w", err)
	}

	return nil
}

// isNotFoundError checks if an error from go-keyring indicates that an item was not found.
// This is necessary because go-keyring doesn't have a specific error type for not found.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()

	// Check for common "not found" error messages across different platforms
	notFoundPatterns := []string{
		"not found",
		"could not find",
		"does not exist",
		"no such",
		"errSecItemNotFound", // macOS
		"Element not found",  // Windows
		"The specified item could not be found in the keyring", // Linux
	}

	for _, pattern := range notFoundPatterns {
		if strings.Contains(strings.ToLower(errMsg), strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}
