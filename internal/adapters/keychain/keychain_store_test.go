package keychain

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/zalando/go-keyring"
)

// init disables the mock backend for go-keyring to test against the real OS keychain
// in integration tests. For CI environments without keychain access, tests will be skipped.
func init() {
	// We'll handle keychain availability per test
}

// cleanupCredential is a helper that removes a credential from the keychain.
// It's used in test cleanup to ensure tests don't leave artifacts.
func cleanupCredential(serviceName, keychainID string) {
	_ = keyring.Delete(serviceName, keychainID)
}

// TestKeychainStore_Store tests storing credentials in the keychain.
func TestKeychainStore_Store(t *testing.T) {
	if !isKeychainAvailable(t) {
		t.Skip("Keychain not available in this environment")
	}

	store := NewKeychainStore()

	tests := []struct {
		name        string
		workspaceID string
		config      domain.WorkspaceConfig
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid credentials",
			workspaceID: "test-workspace-1",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				ProjectKey: "TEST",
				Username:   "user@test.com",
				APIToken:   "test-token-123",
			},
			wantErr: false,
		},
		{
			name:        "empty workspace ID",
			workspaceID: "",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				ProjectKey: "TEST",
				Username:   "user@test.com",
				APIToken:   "test-token-123",
			},
			wantErr:     true,
			errContains: "workspace ID cannot be empty",
		},
		{
			name:        "invalid config - missing Jira URL",
			workspaceID: "test-workspace-2",
			config: domain.WorkspaceConfig{
				JiraURL:    "",
				ProjectKey: "TEST",
				Username:   "user@test.com",
				APIToken:   "test-token-123",
			},
			wantErr:     true,
			errContains: "Jira URL is required",
		},
		{
			name:        "invalid config - missing project key",
			workspaceID: "test-workspace-3",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				ProjectKey: "",
				Username:   "user@test.com",
				APIToken:   "test-token-123",
			},
			wantErr:     true,
			errContains: "project key is required",
		},
		{
			name:        "invalid config - missing username",
			workspaceID: "test-workspace-4",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				ProjectKey: "TEST",
				Username:   "",
				APIToken:   "test-token-123",
			},
			wantErr:     true,
			errContains: "username is required",
		},
		{
			name:        "invalid config - missing API token",
			workspaceID: "test-workspace-5",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				ProjectKey: "TEST",
				Username:   "user@test.com",
				APIToken:   "",
			},
			wantErr:     true,
			errContains: "API token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cleanup before and after test
			t.Cleanup(func() {
				cleanupCredential(serviceName, tt.workspaceID)
			})
			cleanupCredential(serviceName, tt.workspaceID)

			ref, err := store.Store(tt.workspaceID, tt.config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Store() expected error containing %q, got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Store() error = %q, want error containing %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Store() unexpected error = %v", err)
				return
			}

			// Verify the reference
			if ref.KeychainID != tt.workspaceID {
				t.Errorf("Store() ref.KeychainID = %q, want %q", ref.KeychainID, tt.workspaceID)
			}
			if ref.ServiceID != serviceName {
				t.Errorf("Store() ref.ServiceID = %q, want %q", ref.ServiceID, serviceName)
			}
		})
	}
}

// TestKeychainStore_Retrieve tests retrieving credentials from the keychain.
func TestKeychainStore_Retrieve(t *testing.T) {
	if !isKeychainAvailable(t) {
		t.Skip("Keychain not available in this environment")
	}

	store := NewKeychainStore()

	// Setup: Store a credential
	workspaceID := "test-retrieve-workspace"
	expectedConfig := domain.WorkspaceConfig{
		JiraURL:    "https://retrieve-test.atlassian.net",
		ProjectKey: "RETR",
		Username:   "retrieve@test.com",
		APIToken:   "retrieve-token-456",
	}

	t.Cleanup(func() {
		cleanupCredential(serviceName, workspaceID)
	})
	cleanupCredential(serviceName, workspaceID)

	ref, err := store.Store(workspaceID, expectedConfig)
	if err != nil {
		t.Fatalf("Setup: Store() error = %v", err)
	}

	tests := []struct {
		name        string
		ref         domain.CredentialRef
		wantConfig  *domain.WorkspaceConfig
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid reference",
			ref:        ref,
			wantConfig: &expectedConfig,
			wantErr:    false,
		},
		{
			name: "empty keychain ID",
			ref: domain.CredentialRef{
				KeychainID: "",
				ServiceID:  serviceName,
			},
			wantErr:     true,
			errContains: "keychain ID cannot be empty",
		},
		{
			name: "invalid service ID",
			ref: domain.CredentialRef{
				KeychainID: workspaceID,
				ServiceID:  "wrong-service",
			},
			wantErr:     true,
			errContains: "invalid service ID",
		},
		{
			name: "non-existent credential",
			ref: domain.CredentialRef{
				KeychainID: "non-existent-workspace",
				ServiceID:  serviceName,
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := store.Retrieve(tt.ref)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Retrieve() expected error containing %q, got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Retrieve() error = %q, want error containing %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Retrieve() unexpected error = %v", err)
				return
			}

			if config == nil {
				t.Fatal("Retrieve() returned nil config")
			}

			// Verify the retrieved config matches the expected config
			if config.JiraURL != tt.wantConfig.JiraURL {
				t.Errorf("Retrieve() config.JiraURL = %q, want %q", config.JiraURL, tt.wantConfig.JiraURL)
			}
			if config.ProjectKey != tt.wantConfig.ProjectKey {
				t.Errorf("Retrieve() config.ProjectKey = %q, want %q", config.ProjectKey, tt.wantConfig.ProjectKey)
			}
			if config.Username != tt.wantConfig.Username {
				t.Errorf("Retrieve() config.Username = %q, want %q", config.Username, tt.wantConfig.Username)
			}
			if config.APIToken != tt.wantConfig.APIToken {
				t.Errorf("Retrieve() config.APIToken = %q, want %q", config.APIToken, tt.wantConfig.APIToken)
			}
		})
	}
}

// TestKeychainStore_Delete tests deleting credentials from the keychain.
func TestKeychainStore_Delete(t *testing.T) {
	if !isKeychainAvailable(t) {
		t.Skip("Keychain not available in this environment")
	}

	store := NewKeychainStore()

	// Setup: Store a credential
	workspaceID := "test-delete-workspace"
	config := domain.WorkspaceConfig{
		JiraURL:    "https://delete-test.atlassian.net",
		ProjectKey: "DEL",
		Username:   "delete@test.com",
		APIToken:   "delete-token-789",
	}

	t.Cleanup(func() {
		cleanupCredential(serviceName, workspaceID)
	})
	cleanupCredential(serviceName, workspaceID)

	ref, err := store.Store(workspaceID, config)
	if err != nil {
		t.Fatalf("Setup: Store() error = %v", err)
	}

	tests := []struct {
		name        string
		ref         domain.CredentialRef
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid reference",
			ref:     ref,
			wantErr: false,
		},
		{
			name: "empty keychain ID",
			ref: domain.CredentialRef{
				KeychainID: "",
				ServiceID:  serviceName,
			},
			wantErr:     true,
			errContains: "keychain ID cannot be empty",
		},
		{
			name: "invalid service ID",
			ref: domain.CredentialRef{
				KeychainID: workspaceID,
				ServiceID:  "wrong-service",
			},
			wantErr:     true,
			errContains: "invalid service ID",
		},
		{
			name: "non-existent credential (idempotent)",
			ref: domain.CredentialRef{
				KeychainID: "non-existent-delete-workspace",
				ServiceID:  serviceName,
			},
			wantErr: false, // Delete is idempotent
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Delete(tt.ref)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Delete() expected error containing %q, got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Delete() error = %q, want error containing %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Delete() unexpected error = %v", err)
			}
		})
	}

	// Verify the credential was actually deleted
	t.Run("verify deletion", func(t *testing.T) {
		_, err := store.Retrieve(ref)
		if err == nil {
			t.Error("Retrieve() after Delete() should return error, got nil")
		}
	})
}

// TestKeychainStore_List tests listing all credential references.
func TestKeychainStore_List(t *testing.T) {
	if !isKeychainAvailable(t) {
		t.Skip("Keychain not available in this environment")
	}

	store := NewKeychainStore()

	// Cleanup any existing credentials from previous tests
	t.Cleanup(func() {
		// Clean up the list itself
		cleanupCredential(serviceName, listKeychainID)
		// Clean up test credentials
		cleanupCredential(serviceName, "list-test-1")
		cleanupCredential(serviceName, "list-test-2")
		cleanupCredential(serviceName, "list-test-3")
	})

	// Start with a clean slate
	cleanupCredential(serviceName, listKeychainID)
	cleanupCredential(serviceName, "list-test-1")
	cleanupCredential(serviceName, "list-test-2")
	cleanupCredential(serviceName, "list-test-3")

	t.Run("empty list initially", func(t *testing.T) {
		refs, err := store.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if len(refs) != 0 {
			t.Errorf("List() returned %d credentials, want 0", len(refs))
		}
	})

	// Store multiple credentials
	workspaceIDs := []string{"list-test-1", "list-test-2", "list-test-3"}
	storedRefs := make([]domain.CredentialRef, 0, len(workspaceIDs))

	for i, id := range workspaceIDs {
		config := domain.WorkspaceConfig{
			JiraURL:    fmt.Sprintf("https://list-test-%d.atlassian.net", i+1),
			ProjectKey: fmt.Sprintf("LST%d", i+1),
			Username:   fmt.Sprintf("user%d@test.com", i+1),
			APIToken:   fmt.Sprintf("token-%d", i+1),
		}

		ref, err := store.Store(id, config)
		if err != nil {
			t.Fatalf("Store() for %s error = %v", id, err)
		}
		storedRefs = append(storedRefs, ref)
	}

	t.Run("list all credentials", func(t *testing.T) {
		refs, err := store.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(refs) != len(workspaceIDs) {
			t.Errorf("List() returned %d credentials, want %d", len(refs), len(workspaceIDs))
		}

		// Verify all expected IDs are present
		foundIDs := make(map[string]bool)
		for _, ref := range refs {
			foundIDs[ref.KeychainID] = true
			if ref.ServiceID != serviceName {
				t.Errorf("List() ref.ServiceID = %q, want %q", ref.ServiceID, serviceName)
			}
		}

		for _, expectedID := range workspaceIDs {
			if !foundIDs[expectedID] {
				t.Errorf("List() missing expected workspace ID %q", expectedID)
			}
		}
	})

	// Delete one credential and verify the list updates
	t.Run("list after deletion", func(t *testing.T) {
		err := store.Delete(storedRefs[1]) // Delete the second credential
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		refs, err := store.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(refs) != len(workspaceIDs)-1 {
			t.Errorf("List() after delete returned %d credentials, want %d", len(refs), len(workspaceIDs)-1)
		}

		// Verify the deleted ID is not in the list
		for _, ref := range refs {
			if ref.KeychainID == storedRefs[1].KeychainID {
				t.Errorf("List() still contains deleted credential %q", ref.KeychainID)
			}
		}
	})
}

// TestKeychainStore_StoreRetrieveRoundTrip tests the full lifecycle of storing and retrieving credentials.
func TestKeychainStore_StoreRetrieveRoundTrip(t *testing.T) {
	if !isKeychainAvailable(t) {
		t.Skip("Keychain not available in this environment")
	}

	store := NewKeychainStore()

	workspaceID := "roundtrip-test-workspace"
	originalConfig := domain.WorkspaceConfig{
		JiraURL:    "https://roundtrip.atlassian.net",
		ProjectKey: "ROUND",
		Username:   "roundtrip@test.com",
		APIToken:   "roundtrip-token-secret-12345",
	}

	t.Cleanup(func() {
		cleanupCredential(serviceName, workspaceID)
	})
	cleanupCredential(serviceName, workspaceID)

	// Store the credentials
	ref, err := store.Store(workspaceID, originalConfig)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	// Retrieve the credentials
	retrievedConfig, err := store.Retrieve(ref)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}

	// Verify the round-trip integrity
	if retrievedConfig.JiraURL != originalConfig.JiraURL {
		t.Errorf("Round-trip JiraURL = %q, want %q", retrievedConfig.JiraURL, originalConfig.JiraURL)
	}
	if retrievedConfig.ProjectKey != originalConfig.ProjectKey {
		t.Errorf("Round-trip ProjectKey = %q, want %q", retrievedConfig.ProjectKey, originalConfig.ProjectKey)
	}
	if retrievedConfig.Username != originalConfig.Username {
		t.Errorf("Round-trip Username = %q, want %q", retrievedConfig.Username, originalConfig.Username)
	}
	if retrievedConfig.APIToken != originalConfig.APIToken {
		t.Errorf("Round-trip APIToken = %q, want %q", retrievedConfig.APIToken, originalConfig.APIToken)
	}

	// Test updating a credential (store with same ID)
	updatedConfig := domain.WorkspaceConfig{
		JiraURL:    "https://updated-roundtrip.atlassian.net",
		ProjectKey: "UPDT",
		Username:   "updated@test.com",
		APIToken:   "updated-token-999",
	}

	ref2, err := store.Store(workspaceID, updatedConfig)
	if err != nil {
		t.Fatalf("Store() update error = %v", err)
	}

	// Verify the reference is the same
	if ref2.KeychainID != ref.KeychainID {
		t.Errorf("Updated ref.KeychainID = %q, want %q", ref2.KeychainID, ref.KeychainID)
	}

	// Retrieve the updated credentials
	retrievedUpdated, err := store.Retrieve(ref2)
	if err != nil {
		t.Fatalf("Retrieve() after update error = %v", err)
	}

	// Verify the updated values
	if retrievedUpdated.JiraURL != updatedConfig.JiraURL {
		t.Errorf("Updated JiraURL = %q, want %q", retrievedUpdated.JiraURL, updatedConfig.JiraURL)
	}
	if retrievedUpdated.APIToken != updatedConfig.APIToken {
		t.Errorf("Updated APIToken = %q, want %q", retrievedUpdated.APIToken, updatedConfig.APIToken)
	}

	// Delete and verify
	err = store.Delete(ref)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = store.Retrieve(ref)
	if err == nil {
		t.Error("Retrieve() after Delete() should return error, got nil")
	}
}

// TestKeychainStore_ConcurrentAccess tests thread-safety of the keychain store.
func TestKeychainStore_ConcurrentAccess(t *testing.T) {
	if !isKeychainAvailable(t) {
		t.Skip("Keychain not available in this environment")
	}

	store := NewKeychainStore()

	// Cleanup
	workspaceIDs := make([]string, 20)
	for i := 0; i < 20; i++ {
		workspaceIDs[i] = fmt.Sprintf("concurrent-test-%d", i)
	}

	t.Cleanup(func() {
		for _, id := range workspaceIDs {
			cleanupCredential(serviceName, id)
		}
		cleanupCredential(serviceName, listKeychainID)
	})

	for _, id := range workspaceIDs {
		cleanupCredential(serviceName, id)
	}
	cleanupCredential(serviceName, listKeychainID)

	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	// Concurrent stores
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			config := domain.WorkspaceConfig{
				JiraURL:    fmt.Sprintf("https://concurrent-%d.atlassian.net", idx),
				ProjectKey: fmt.Sprintf("CON%d", idx),
				Username:   fmt.Sprintf("user%d@test.com", idx),
				APIToken:   fmt.Sprintf("token-%d", idx),
			}

			_, err := store.Store(workspaceIDs[idx], config)
			if err != nil {
				errChan <- fmt.Errorf("concurrent Store() for workspace %d: %w", idx, err)
			}
		}(i)
	}

	// Concurrent lists
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := store.List()
			if err != nil {
				errChan <- fmt.Errorf("concurrent List(): %w", err)
			}
		}()
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Error(err)
	}

	// Verify all credentials were stored
	refs, err := store.List()
	if err != nil {
		t.Fatalf("Final List() error = %v", err)
	}

	if len(refs) != 20 {
		t.Errorf("After concurrent operations, List() returned %d credentials, want 20", len(refs))
	}
}

// TestKeychainStore_SpecialCharacters tests handling of special characters in credentials.
func TestKeychainStore_SpecialCharacters(t *testing.T) {
	if !isKeychainAvailable(t) {
		t.Skip("Keychain not available in this environment")
	}

	store := NewKeychainStore()

	workspaceID := "special-chars-test"
	config := domain.WorkspaceConfig{
		JiraURL:    "https://test.atlassian.net",
		ProjectKey: "SPEC",
		Username:   "user+tag@test.com",
		APIToken:   `token-with-"quotes"-and-\slashes-and-新字符-🔐`,
	}

	t.Cleanup(func() {
		cleanupCredential(serviceName, workspaceID)
	})
	cleanupCredential(serviceName, workspaceID)

	// Store
	ref, err := store.Store(workspaceID, config)
	if err != nil {
		t.Fatalf("Store() with special chars error = %v", err)
	}

	// Retrieve
	retrievedConfig, err := store.Retrieve(ref)
	if err != nil {
		t.Fatalf("Retrieve() with special chars error = %v", err)
	}

	// Verify special characters are preserved
	if retrievedConfig.Username != config.Username {
		t.Errorf("Special chars Username = %q, want %q", retrievedConfig.Username, config.Username)
	}
	if retrievedConfig.APIToken != config.APIToken {
		t.Errorf("Special chars APIToken = %q, want %q", retrievedConfig.APIToken, config.APIToken)
	}
}

// isKeychainAvailable checks if the OS keychain is available in the test environment.
// This helps gracefully skip tests in CI/CD environments without keychain access.
func isKeychainAvailable(t *testing.T) bool {
	t.Helper()

	// Try a simple keychain operation to test availability
	testKey := "ticketr-test-availability-check"
	err := keyring.Set(serviceName, testKey, "test-value")
	if err != nil {
		t.Logf("Keychain not available: %v", err)
		return false
	}

	// Clean up the test entry
	_ = keyring.Delete(serviceName, testKey)
	return true
}

// BenchmarkKeychainStore_Store benchmarks storing credentials.
func BenchmarkKeychainStore_Store(b *testing.B) {
	store := NewKeychainStore()

	config := domain.WorkspaceConfig{
		JiraURL:    "https://benchmark.atlassian.net",
		ProjectKey: "BENCH",
		Username:   "benchmark@test.com",
		APIToken:   "benchmark-token",
	}

	b.Cleanup(func() {
		for i := 0; i < b.N; i++ {
			cleanupCredential(serviceName, fmt.Sprintf("benchmark-store-%d", i))
		}
		cleanupCredential(serviceName, listKeychainID)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Store(fmt.Sprintf("benchmark-store-%d", i), config)
	}
}

// BenchmarkKeychainStore_Retrieve benchmarks retrieving credentials.
func BenchmarkKeychainStore_Retrieve(b *testing.B) {
	store := NewKeychainStore()

	config := domain.WorkspaceConfig{
		JiraURL:    "https://benchmark.atlassian.net",
		ProjectKey: "BENCH",
		Username:   "benchmark@test.com",
		APIToken:   "benchmark-token",
	}

	workspaceID := "benchmark-retrieve"
	ref, err := store.Store(workspaceID, config)
	if err != nil {
		b.Fatalf("Setup: Store() error = %v", err)
	}

	b.Cleanup(func() {
		cleanupCredential(serviceName, workspaceID)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Retrieve(ref)
	}
}

// BenchmarkKeychainStore_List benchmarks listing credentials.
func BenchmarkKeychainStore_List(b *testing.B) {
	store := NewKeychainStore()

	// Setup: Store 10 credentials
	config := domain.WorkspaceConfig{
		JiraURL:    "https://benchmark.atlassian.net",
		ProjectKey: "BENCH",
		Username:   "benchmark@test.com",
		APIToken:   "benchmark-token",
	}

	workspaceIDs := make([]string, 10)
	for i := 0; i < 10; i++ {
		workspaceIDs[i] = fmt.Sprintf("benchmark-list-%d", i)
		_, _ = store.Store(workspaceIDs[i], config)
	}

	b.Cleanup(func() {
		for _, id := range workspaceIDs {
			cleanupCredential(serviceName, id)
		}
		cleanupCredential(serviceName, listKeychainID)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.List()
	}
}
