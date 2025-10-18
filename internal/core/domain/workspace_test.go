package domain

import (
	"strings"
	"testing"
	"time"
)

// TestWorkspace_Validation tests workspace field validation
func TestWorkspace_Validation(t *testing.T) {
	tests := []struct {
		name      string
		workspace Workspace
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid workspace",
			workspace: Workspace{
				ID:         "ws-1",
				Name:       "backend",
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "BACK",
				IsDefault:  false,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			workspace: Workspace{
				ID:         "ws-2",
				Name:       "",
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "BACK",
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "empty jira url",
			workspace: Workspace{
				ID:         "ws-3",
				Name:       "backend",
				JiraURL:    "",
				ProjectKey: "BACK",
			},
			wantErr: true,
			errMsg:  "jira_url cannot be empty",
		},
		{
			name: "empty project key",
			workspace: Workspace{
				ID:         "ws-4",
				Name:       "backend",
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "",
			},
			wantErr: true,
			errMsg:  "project_key cannot be empty",
		},
		{
			name: "name with special characters",
			workspace: Workspace{
				ID:         "ws-6",
				Name:       "backend@123!",
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "BACK",
			},
			wantErr: true,
			errMsg:  "name must contain only alphanumeric characters, hyphens, and underscores",
		},
		{
			name: "valid name with hyphens and underscores",
			workspace: Workspace{
				ID:         "ws-7",
				Name:       "backend_team-1",
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "BACK",
			},
			wantErr: false,
		},
		{
			name: "project key with numbers",
			workspace: Workspace{
				ID:         "ws-9",
				Name:       "backend",
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "BACK123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.workspace.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Workspace.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Workspace.Validate() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

// TestWorkspace_NameValidation tests specific name validation rules
func TestWorkspace_NameValidation(t *testing.T) {
	tests := []struct {
		name    string
		wsName  string
		wantErr bool
	}{
		{"alphanumeric only", "backend123", false},
		{"with hyphen", "backend-team", false},
		{"with underscore", "backend_team", false},
		{"mixed valid chars", "backend_team-v2", false},
		{"starting with number", "1backend", false},
		{"with space", "backend team", true},
		{"with dot", "backend.team", true},
		{"with slash", "backend/team", true},
		{"with at sign", "backend@team", true},
		{"empty string", "", true},
		{"too long", strings.Repeat("a", 256), true},
		{"single character", "b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := Workspace{
				ID:         "test-id",
				Name:       tt.wsName,
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "TEST",
			}
			err := ws.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() with name %q error = %v, wantErr %v", tt.wsName, err, tt.wantErr)
			}
		})
	}
}

// TestWorkspace_DefaultBehavior tests default workspace constraints
func TestWorkspace_DefaultBehavior(t *testing.T) {
	// Test that only one workspace can be default is enforced at repository level
	// Here we just test the field behavior
	ws1 := Workspace{
		ID:         "ws-1",
		Name:       "workspace1",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "WS1",
		IsDefault:  true,
	}

	ws2 := Workspace{
		ID:         "ws-2",
		Name:       "workspace2",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "WS2",
		IsDefault:  false,
	}

	if !ws1.IsDefault {
		t.Error("Expected ws1.IsDefault to be true")
	}

	if ws2.IsDefault {
		t.Error("Expected ws2.IsDefault to be false")
	}
}

// TestWorkspace_LastUsedTracking tests last used timestamp tracking
func TestWorkspace_LastUsedTracking(t *testing.T) {
	ws := Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
	}

	// Initially, LastUsed should be zero
	if !ws.LastUsed.IsZero() {
		t.Error("Expected LastUsed to be zero initially")
	}

	// Update last used
	now := time.Now()
	ws.Touch()

	if ws.LastUsed.IsZero() {
		t.Error("Expected LastUsed to be set after Touch()")
	}

	// Should be within 1 second of now
	diff := ws.LastUsed.Sub(now)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("Expected LastUsed to be close to now, got diff %v", diff)
	}
}

// TestWorkspace_CreatedAtUpdatedAt tests timestamp management
func TestWorkspace_CreatedAtUpdatedAt(t *testing.T) {
	ws := Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
	}

	// Set creation time
	ws.CreatedAt = time.Now()
	ws.UpdatedAt = ws.CreatedAt

	if ws.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if ws.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	if ws.UpdatedAt.Before(ws.CreatedAt) {
		t.Error("UpdatedAt should not be before CreatedAt")
	}

	// Simulate update
	time.Sleep(10 * time.Millisecond)
	ws.UpdatedAt = time.Now()

	if !ws.UpdatedAt.After(ws.CreatedAt) {
		t.Error("Expected UpdatedAt to be after CreatedAt after update")
	}
}

// TestWorkspace_CredentialRef tests credential reference handling
func TestWorkspace_CredentialRef(t *testing.T) {
	ws := Workspace{
		ID:          "ws-1",
		Name:        "backend",
		JiraURL:     "https://company.atlassian.net",
		ProjectKey:  "BACK",
		Credentials: CredentialRef{KeychainID: "keychain-backend", ServiceID: "ticketr"},
	}

	if ws.Credentials.KeychainID == "" {
		t.Error("Expected Credentials.KeychainID to be set")
	}

	if ws.Credentials.ServiceID == "" {
		t.Error("Expected Credentials.ServiceID to be set")
	}
}

// TestWorkspace_SetDefault tests the SetDefault method
func TestWorkspace_SetDefault(t *testing.T) {
	ws := Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		IsDefault:  false,
	}

	ws.SetDefault(true)
	if !ws.IsDefault {
		t.Error("Expected workspace to be default after SetDefault(true)")
	}

	ws.SetDefault(false)
	if ws.IsDefault {
		t.Error("Expected workspace to not be default after SetDefault(false)")
	}
}

// TestWorkspace_Touch tests the Touch method
func TestWorkspace_Touch(t *testing.T) {
	ws := Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
	}

	before := ws.LastUsed
	time.Sleep(10 * time.Millisecond)
	ws.Touch()

	if !ws.LastUsed.After(before) {
		t.Error("Expected LastUsed to be updated after Touch()")
	}
}

// BenchmarkWorkspace_Validate benchmarks validation performance
func BenchmarkWorkspace_Validate(b *testing.B) {
	ws := Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ws.Validate()
	}
}
