package domain

import (
	"testing"
	"time"
)

func TestJQLAlias_Validate(t *testing.T) {
	tests := []struct {
		name    string
		alias   JQLAlias
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid alias",
			alias: JQLAlias{
				ID:          "test-id",
				Name:        "my-alias",
				JQL:         "assignee = currentUser()",
				Description: "My tickets",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			alias: JQLAlias{
				ID:   "test-id",
				Name: "",
				JQL:  "assignee = currentUser()",
			},
			wantErr: true,
			errMsg:  "alias name cannot be empty",
		},
		{
			name: "invalid name with spaces",
			alias: JQLAlias{
				ID:   "test-id",
				Name: "my alias",
				JQL:  "assignee = currentUser()",
			},
			wantErr: true,
			errMsg:  "alias name must contain only alphanumeric characters",
		},
		{
			name: "name too long",
			alias: JQLAlias{
				ID:   "test-id",
				Name: "a123456789012345678901234567890123456789012345678901234567890123456789",
				JQL:  "assignee = currentUser()",
			},
			wantErr: true,
			errMsg:  "alias name must be 64 characters or less",
		},
		{
			name: "empty JQL",
			alias: JQLAlias{
				ID:   "test-id",
				Name: "my-alias",
				JQL:  "",
			},
			wantErr: true,
			errMsg:  "JQL query cannot be empty",
		},
		{
			name: "JQL too long",
			alias: JQLAlias{
				ID:   "test-id",
				Name: "my-alias",
				JQL:  string(make([]byte, 2001)),
			},
			wantErr: true,
			errMsg:  "JQL query must be 2000 characters or less",
		},
		{
			name: "predefined alias with workspace ID",
			alias: JQLAlias{
				ID:           "test-id",
				Name:         "my-alias",
				JQL:          "assignee = currentUser()",
				IsPredefined: true,
				WorkspaceID:  "workspace-id",
			},
			wantErr: true,
			errMsg:  "predefined aliases cannot be workspace-specific",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.alias.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidateAliasName(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		wantErr bool
	}{
		{"valid simple name", "mine", false},
		{"valid with hyphen", "my-alias", false},
		{"valid with underscore", "my_alias", false},
		{"valid mixed", "my-alias_123", false},
		{"empty name", "", true},
		{"name with space", "my alias", true},
		{"name with special char", "my@alias", true},
		{"name too long", string(make([]byte, 65)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAliasName(tt.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAliasName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsPredefinedAlias(t *testing.T) {
	tests := []struct {
		name     string
		alias    string
		expected bool
	}{
		{"mine is predefined", "mine", true},
		{"sprint is predefined", "sprint", true},
		{"blocked is predefined", "blocked", true},
		{"custom is not predefined", "custom", false},
		{"my-bugs is not predefined", "my-bugs", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPredefinedAlias(tt.alias)
			if result != tt.expected {
				t.Errorf("IsPredefinedAlias(%s) = %v, want %v", tt.alias, result, tt.expected)
			}
		})
	}
}

func TestGetPredefinedAlias(t *testing.T) {
	tests := []struct {
		name      string
		aliasName string
		wantNil   bool
		wantJQL   string
	}{
		{
			name:      "get mine alias",
			aliasName: "mine",
			wantNil:   false,
			wantJQL:   "assignee = currentUser()",
		},
		{
			name:      "get sprint alias",
			aliasName: "sprint",
			wantNil:   false,
			wantJQL:   "sprint in openSprints()",
		},
		{
			name:      "get blocked alias",
			aliasName: "blocked",
			wantNil:   false,
			wantJQL:   "status = Blocked OR labels in (blocked)",
		},
		{
			name:      "get non-existent alias",
			aliasName: "custom",
			wantNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPredefinedAlias(tt.aliasName)
			if (result == nil) != tt.wantNil {
				t.Errorf("GetPredefinedAlias(%s) nil = %v, want %v", tt.aliasName, result == nil, tt.wantNil)
				return
			}
			if result != nil && result.JQL != tt.wantJQL {
				t.Errorf("GetPredefinedAlias(%s).JQL = %s, want %s", tt.aliasName, result.JQL, tt.wantJQL)
			}
		})
	}
}

func TestPredefinedAliases(t *testing.T) {
	// Test that all predefined aliases are valid
	for name, alias := range PredefinedAliases {
		t.Run("validate_"+name, func(t *testing.T) {
			if err := alias.Validate(); err != nil {
				t.Errorf("Predefined alias %s is invalid: %v", name, err)
			}
		})
	}

	// Test that expected aliases exist
	expectedAliases := []string{"mine", "sprint", "blocked"}
	for _, name := range expectedAliases {
		t.Run("exists_"+name, func(t *testing.T) {
			if _, exists := PredefinedAliases[name]; !exists {
				t.Errorf("Expected predefined alias %s not found", name)
			}
		})
	}
}

func TestJQLAlias_CreationTimestamps(t *testing.T) {
	now := time.Now()
	alias := JQLAlias{
		ID:          "test-id",
		Name:        "test",
		JQL:         "assignee = currentUser()",
		Description: "Test alias",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if alias.CreatedAt != now {
		t.Errorf("CreatedAt = %v, want %v", alias.CreatedAt, now)
	}
	if alias.UpdatedAt != now {
		t.Errorf("UpdatedAt = %v, want %v", alias.UpdatedAt, now)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr))
}
