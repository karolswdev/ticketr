// +build integration

package jira

import (
	"context"
	"os"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// TestJiraAdapterIntegration validates the Jira adapter against a real Jira instance
func TestJiraAdapterIntegration(t *testing.T) {
	// Skip if no workspace config available
	config := getIntegrationConfig(t)
	if config == nil {
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
		t.Run(tt.name, func(t *testing.T) {
			adapter := createIntegrationAdapter(t, config)
			tt.testFunc(t, adapter)
		})
	}
}

// TestAdapterValidation validates adapter configuration requirements
func TestAdapterValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *domain.WorkspaceConfig
		expectError bool
		errorText   string
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			errorText:   "workspace configuration is required",
		},
		{
			name: "missing URL",
			config: &domain.WorkspaceConfig{
				Username:   "test@example.com",
				APIToken:   "token",
				ProjectKey: "TEST",
			},
			expectError: true,
			errorText:   "Jira URL is required",
		},
		{
			name: "missing username",
			config: &domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				APIToken:   "token",
				ProjectKey: "TEST",
			},
			expectError: true,
			errorText:   "username is required",
		},
		{
			name: "missing token",
			config: &domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				Username:   "test@example.com",
				ProjectKey: "TEST",
			},
			expectError: true,
			errorText:   "API token is required",
		},
		{
			name: "missing project key",
			config: &domain.WorkspaceConfig{
				JiraURL:  "https://test.atlassian.net",
				Username: "test@example.com",
				APIToken: "token",
			},
			expectError: true,
			errorText:   "project key is required",
		},
		{
			name: "valid config",
			config: &domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				Username:   "test@example.com",
				APIToken:   "token",
				ProjectKey: "TEST",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, err := NewJiraAdapterFromConfig(tt.config, nil)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q but got nil", tt.errorText)
					return
				}
				if tt.errorText != "" {
					errMsg := err.Error()
					if len(errMsg) == 0 || !contains(errMsg, tt.errorText) {
						t.Errorf("expected error containing %q, got: %q", tt.errorText, errMsg)
					}
				}
				if adapter != nil {
					t.Error("expected nil adapter on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if adapter == nil {
					t.Error("expected adapter but got nil")
				}
			}
		})
	}
}

// Helper functions for integration tests

func getIntegrationConfig(t *testing.T) *domain.WorkspaceConfig {
	t.Helper()

	url := os.Getenv("JIRA_URL")
	email := os.Getenv("JIRA_EMAIL")
	token := os.Getenv("JIRA_API_KEY")
	project := os.Getenv("JIRA_PROJECT_KEY")

	if url == "" || email == "" || token == "" || project == "" {
		return nil
	}

	return &domain.WorkspaceConfig{
		JiraURL:    url,
		Username:   email,
		APIToken:   token,
		ProjectKey: project,
	}
}

func createIntegrationAdapter(t *testing.T, config *domain.WorkspaceConfig) ports.JiraPort {
	t.Helper()

	fieldMappings := getDefaultFieldMappings()
	adapter, err := NewJiraAdapterFromConfig(config, fieldMappings)
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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
