package keychain

import (
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// TestKeychainStore_ValidationLogic tests input validation without keychain access.
func TestKeychainStore_ValidationLogic(t *testing.T) {
	store := NewKeychainStore()

	t.Run("service name is set correctly", func(t *testing.T) {
		if store.serviceName != "ticketr" {
			t.Errorf("serviceName = %q, want %q", store.serviceName, "ticketr")
		}
	})

	t.Run("new instance is not nil", func(t *testing.T) {
		if store == nil {
			t.Error("NewKeychainStore() returned nil")
		}
	})
}

// TestIsNotFoundError tests the error detection logic.
func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected bool
	}{
		{
			name:     "nil error",
			errMsg:   "",
			expected: false,
		},
		{
			name:     "macOS not found",
			errMsg:   "errSecItemNotFound: The specified item could not be found in the keychain",
			expected: true,
		},
		{
			name:     "Windows not found",
			errMsg:   "Element not found in the credential manager",
			expected: true,
		},
		{
			name:     "Linux not found",
			errMsg:   "The specified item could not be found in the keyring",
			expected: true,
		},
		{
			name:     "generic not found",
			errMsg:   "credential not found",
			expected: true,
		},
		{
			name:     "case insensitive",
			errMsg:   "Item NOT FOUND",
			expected: true,
		},
		{
			name:     "does not exist",
			errMsg:   "the credential does not exist",
			expected: true,
		},
		{
			name:     "no such item",
			errMsg:   "no such credential",
			expected: true,
		},
		{
			name:     "could not find",
			errMsg:   "could not find the requested credential",
			expected: true,
		},
		{
			name:     "unrelated error",
			errMsg:   "failed to connect to keychain service",
			expected: false,
		},
		{
			name:     "permission denied",
			errMsg:   "permission denied",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.errMsg != "" {
				err = &mockError{msg: tt.errMsg}
			}

			result := isNotFoundError(err)
			if result != tt.expected {
				t.Errorf("isNotFoundError(%q) = %v, want %v", tt.errMsg, result, tt.expected)
			}
		})
	}
}

// TestCredentialRef_Structure tests the CredentialRef structure.
func TestCredentialRef_Structure(t *testing.T) {
	ref := domain.CredentialRef{
		KeychainID: "test-workspace",
		ServiceID:  "ticketr",
	}

	if ref.KeychainID != "test-workspace" {
		t.Errorf("KeychainID = %q, want %q", ref.KeychainID, "test-workspace")
	}

	if ref.ServiceID != "ticketr" {
		t.Errorf("ServiceID = %q, want %q", ref.ServiceID, "ticketr")
	}
}

// TestWorkspaceConfig_Validation tests domain validation logic.
func TestWorkspaceConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      domain.WorkspaceConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				ProjectKey: "TEST",
				Username:   "user@test.com",
				APIToken:   "token123",
			},
			wantErr: false,
		},
		{
			name: "missing Jira URL",
			config: domain.WorkspaceConfig{
				JiraURL:    "",
				ProjectKey: "TEST",
				Username:   "user@test.com",
				APIToken:   "token123",
			},
			wantErr:     true,
			errContains: "Jira URL",
		},
		{
			name: "missing project key",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				ProjectKey: "",
				Username:   "user@test.com",
				APIToken:   "token123",
			},
			wantErr:     true,
			errContains: "project key",
		},
		{
			name: "missing username",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				ProjectKey: "TEST",
				Username:   "",
				APIToken:   "token123",
			},
			wantErr:     true,
			errContains: "username",
		},
		{
			name: "missing API token",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				ProjectKey: "TEST",
				Username:   "user@test.com",
				APIToken:   "",
			},
			wantErr:     true,
			errContains: "API token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidateConfig(tt.config)

			if tt.wantErr {
				if err == nil {
					t.Error("ValidateConfig() expected error, got nil")
					return
				}
				if tt.errContains != "" {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("ValidateConfig() error = %q, want error containing %q", err.Error(), tt.errContains)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateConfig() unexpected error = %v", err)
				}
			}
		})
	}
}

// mockError is a simple error implementation for testing.
type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) > 0 && len(s) > 0 && s != substr && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
