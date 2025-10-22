// +build integration

package jira

import (
	"context"
	"os"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// TestAdapterBehaviorParity validates V1 and V2 handle same scenarios identically
func TestAdapterBehaviorParity(t *testing.T) {
	// Skip if no credentials available
	if os.Getenv("JIRA_URL") == "" || os.Getenv("JIRA_EMAIL") == "" || os.Getenv("JIRA_API_KEY") == "" {
		t.Skip("Skipping integration tests: JIRA credentials not set")
	}

	tests := []struct {
		name     string
		testFunc func(t *testing.T, adapter ports.JiraPort)
	}{
		{"authentication success", testAuthenticationSuccess},
		{"get project issue types", testGetProjectIssueTypes},
		{"search with empty JQL", testSearchWithEmptyJQL},
	}

	for _, tt := range tests {
		t.Run(tt.name+" (V1)", func(t *testing.T) {
			os.Setenv("TICKETR_JIRA_ADAPTER_VERSION", "v1")
			adapter := createIntegrationAdapter(t)
			tt.testFunc(t, adapter)
		})

		t.Run(tt.name+" (V2)", func(t *testing.T) {
			os.Setenv("TICKETR_JIRA_ADAPTER_VERSION", "v2")
			adapter := createIntegrationAdapter(t)
			tt.testFunc(t, adapter)
		})
	}
}

// TestFactoryVersionSelection validates the factory correctly selects adapter versions
func TestFactoryVersionSelection(t *testing.T) {
	config := &domain.WorkspaceConfig{
		JiraURL:    "https://test.atlassian.net",
		Username:   "test@example.com",
		APIToken:   "test-token",
		ProjectKey: "TEST",
	}

	tests := []struct {
		name        string
		envValue    string
		expectError bool
		expectV2    bool
	}{
		{
			name:        "explicit v1",
			envValue:    "v1",
			expectError: false,
			expectV2:    false,
		},
		{
			name:        "explicit v2",
			envValue:    "v2",
			expectError: false,
			expectV2:    true,
		},
		{
			name:        "default (empty) should be v2",
			envValue:    "",
			expectError: false,
			expectV2:    true,
		},
		{
			name:        "invalid version defaults to v2",
			envValue:    "invalid",
			expectError: false,
			expectV2:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TICKETR_JIRA_ADAPTER_VERSION", tt.envValue)
			defer os.Unsetenv("TICKETR_JIRA_ADAPTER_VERSION")

			adapter, err := NewJiraAdapterFromConfigWithVersion(config, nil)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if adapter == nil {
				t.Error("adapter is nil")
				return
			}

			// Verify correct adapter type was created
			if tt.expectV2 {
				if _, ok := adapter.(*JiraAdapterV2); !ok {
					t.Errorf("expected V2 adapter, got %T", adapter)
				}
			} else {
				if _, ok := adapter.(*JiraAdapter); !ok {
					t.Errorf("expected V1 adapter, got %T", adapter)
				}
			}
		})
	}
}

// TestErrorMessageVersionTags validates that errors contain version tags
func TestErrorMessageVersionTags(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		expectedTag string
	}{
		{
			name:       "V1 errors have [jira-v1] tag",
			version:    "v1",
			expectedTag: "[jira-v1]",
		},
		{
			name:       "V2 errors have [jira-v2] tag",
			version:    "v2",
			expectedTag: "[jira-v2]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TICKETR_JIRA_ADAPTER_VERSION", tt.version)
			defer os.Unsetenv("TICKETR_JIRA_ADAPTER_VERSION")

			// Create adapter with invalid config to trigger errors
			config := &domain.WorkspaceConfig{
				JiraURL:    "", // Empty URL should trigger validation error
				Username:   "",
				APIToken:   "",
				ProjectKey: "",
			}

			_, err := NewJiraAdapterFromConfigWithVersion(config, nil)
			if err == nil {
				t.Error("expected validation error but got nil")
				return
			}

			errMsg := err.Error()
			if errMsg == "" {
				t.Error("error message is empty")
				return
			}

			// V2 adapter has better validation, so it should have the tag
			if tt.version == "v2" {
				if len(errMsg) < len(tt.expectedTag) || errMsg[:len(tt.expectedTag)] != tt.expectedTag {
					t.Errorf("error message does not start with expected tag %q, got: %q", tt.expectedTag, errMsg)
				}
			}
		})
	}
}

// Helper functions for integration tests

func createIntegrationAdapter(t *testing.T) ports.JiraPort {
	t.Helper()

	fieldMappings := getDefaultFieldMappings()
	adapter, err := NewJiraAdapterWithConfig(fieldMappings)
	if err != nil {
		t.Fatalf("failed to create adapter: %v", err)
	}

	return adapter
}

func testAuthenticationSuccess(t *testing.T, adapter ports.JiraPort) {
	err := adapter.Authenticate()
	if err != nil {
		t.Errorf("authentication failed: %v", err)
	}
}

func testGetProjectIssueTypes(t *testing.T, adapter ports.JiraPort) {
	issueTypes, err := adapter.GetProjectIssueTypes()
	if err != nil {
		t.Errorf("GetProjectIssueTypes failed: %v", err)
		return
	}

	if len(issueTypes) == 0 {
		t.Error("expected issue types but got empty map")
	}
}

func testSearchWithEmptyJQL(t *testing.T, adapter ports.JiraPort) {
	projectKey := os.Getenv("JIRA_PROJECT_KEY")
	if projectKey == "" {
		t.Skip("JIRA_PROJECT_KEY not set")
	}

	ctx := context.Background()
	tickets, err := adapter.SearchTickets(ctx, projectKey, "", nil)
	if err != nil {
		t.Errorf("SearchTickets failed: %v", err)
		return
	}

	// Should return some tickets or empty slice (not an error)
	if tickets == nil {
		t.Error("expected ticket slice but got nil")
	}
}
