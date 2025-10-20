package services

import (
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// MockAliasRepository is a mock implementation of ports.AliasRepository for testing.
type MockAliasRepository struct {
	aliases map[string]*domain.JQLAlias
}

func NewMockAliasRepository() *MockAliasRepository {
	return &MockAliasRepository{
		aliases: make(map[string]*domain.JQLAlias),
	}
}

func (m *MockAliasRepository) Create(alias *domain.JQLAlias) error {
	key := alias.Name + ":" + alias.WorkspaceID
	if _, exists := m.aliases[key]; exists {
		return ports.ErrAliasExists
	}
	m.aliases[key] = alias
	return nil
}

func (m *MockAliasRepository) Get(id string) (*domain.JQLAlias, error) {
	for _, alias := range m.aliases {
		if alias.ID == id {
			return alias, nil
		}
	}
	return nil, ports.ErrAliasNotFound
}

func (m *MockAliasRepository) GetByName(name string, workspaceID string) (*domain.JQLAlias, error) {
	// First try workspace-specific alias
	key := name + ":" + workspaceID
	if alias, exists := m.aliases[key]; exists {
		return alias, nil
	}

	// Then try global alias (empty workspaceID)
	globalKey := name + ":"
	if alias, exists := m.aliases[globalKey]; exists {
		return alias, nil
	}

	return nil, ports.ErrAliasNotFound
}

func (m *MockAliasRepository) List(workspaceID string) ([]*domain.JQLAlias, error) {
	result := []*domain.JQLAlias{}
	for _, alias := range m.aliases {
		if alias.WorkspaceID == workspaceID || alias.WorkspaceID == "" {
			result = append(result, alias)
		}
	}
	return result, nil
}

func (m *MockAliasRepository) ListAll() ([]*domain.JQLAlias, error) {
	result := []*domain.JQLAlias{}
	for _, alias := range m.aliases {
		result = append(result, alias)
	}
	return result, nil
}

func (m *MockAliasRepository) Update(alias *domain.JQLAlias) error {
	for key, existing := range m.aliases {
		if existing.ID == alias.ID {
			if existing.IsPredefined {
				return ports.ErrCannotModifyPredefined
			}
			m.aliases[key] = alias
			return nil
		}
	}
	return ports.ErrAliasNotFound
}

func (m *MockAliasRepository) Delete(id string) error {
	for key, alias := range m.aliases {
		if alias.ID == id {
			if alias.IsPredefined {
				return ports.ErrCannotModifyPredefined
			}
			delete(m.aliases, key)
			return nil
		}
	}
	return ports.ErrAliasNotFound
}

func (m *MockAliasRepository) DeleteByName(name string, workspaceID string) error {
	key := name + ":" + workspaceID
	if alias, exists := m.aliases[key]; exists {
		if alias.IsPredefined {
			return ports.ErrCannotModifyPredefined
		}
		delete(m.aliases, key)
		return nil
	}
	return ports.ErrAliasNotFound
}

func TestAliasService_Create(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	tests := []struct {
		name        string
		aliasName   string
		jql         string
		description string
		wantErr     bool
		errType     error
	}{
		{
			name:        "valid alias",
			aliasName:   "my-bugs",
			jql:         "assignee = currentUser() AND type = Bug",
			description: "My bug tickets",
			wantErr:     false,
		},
		{
			name:      "duplicate alias",
			aliasName: "my-bugs",
			jql:       "different JQL",
			wantErr:   true,
			errType:   ports.ErrAliasExists,
		},
		{
			name:      "reserved name",
			aliasName: "mine",
			jql:       "some JQL",
			wantErr:   true,
		},
		{
			name:      "invalid name",
			aliasName: "my alias",
			jql:       "some JQL",
			wantErr:   true,
		},
		{
			name:      "empty JQL",
			aliasName: "test",
			jql:       "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Create(tt.aliasName, tt.jql, tt.description, workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliasService_Get(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create a test alias
	service.Create("my-bugs", "assignee = currentUser() AND type = Bug", "My bugs", workspaceID)

	tests := []struct {
		name      string
		aliasName string
		wantErr   bool
	}{
		{
			name:      "get existing alias",
			aliasName: "my-bugs",
			wantErr:   false,
		},
		{
			name:      "get predefined alias",
			aliasName: "mine",
			wantErr:   false,
		},
		{
			name:      "get non-existent alias",
			aliasName: "does-not-exist",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alias, err := service.Get(tt.aliasName, workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && alias == nil {
				t.Error("Get() returned nil alias without error")
			}
		})
	}
}

func TestAliasService_List(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create test aliases
	service.Create("my-bugs", "assignee = currentUser() AND type = Bug", "My bugs", workspaceID)
	service.Create("high-priority", "priority = High", "High priority", workspaceID)

	aliases, err := service.List(workspaceID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// Should include predefined aliases + user aliases
	expectedMin := len(domain.PredefinedAliases) + 2
	if len(aliases) < expectedMin {
		t.Errorf("List() returned %d aliases, want at least %d", len(aliases), expectedMin)
	}

	// Check that predefined aliases are included
	hasMine := false
	for _, alias := range aliases {
		if alias.Name == "mine" {
			hasMine = true
			break
		}
	}
	if !hasMine {
		t.Error("List() did not include predefined 'mine' alias")
	}
}

func TestAliasService_Delete(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create a test alias
	service.Create("my-bugs", "assignee = currentUser() AND type = Bug", "My bugs", workspaceID)

	tests := []struct {
		name      string
		aliasName string
		wantErr   bool
		errType   error
	}{
		{
			name:      "delete existing alias",
			aliasName: "my-bugs",
			wantErr:   false,
		},
		{
			name:      "delete predefined alias",
			aliasName: "mine",
			wantErr:   true,
			errType:   ports.ErrCannotModifyPredefined,
		},
		{
			name:      "delete non-existent alias",
			aliasName: "does-not-exist",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(tt.aliasName, workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliasService_Update(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create a test alias
	service.Create("my-bugs", "assignee = currentUser() AND type = Bug", "My bugs", workspaceID)

	tests := []struct {
		name        string
		aliasName   string
		newJQL      string
		description string
		wantErr     bool
	}{
		{
			name:        "update existing alias",
			aliasName:   "my-bugs",
			newJQL:      "assignee = currentUser() AND type = Bug AND status != Done",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:      "update predefined alias",
			aliasName: "mine",
			newJQL:    "different JQL",
			wantErr:   true,
		},
		{
			name:      "update non-existent alias",
			aliasName: "does-not-exist",
			newJQL:    "some JQL",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Update(tt.aliasName, tt.newJQL, tt.description, workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliasService_ExpandAlias(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create test aliases
	service.Create("my-bugs", "assignee = currentUser() AND type = Bug", "", workspaceID)

	tests := []struct {
		name        string
		aliasName   string
		wantJQL     string
		wantErr     bool
		containsJQL string
	}{
		{
			name:      "expand simple alias",
			aliasName: "my-bugs",
			wantJQL:   "assignee = currentUser() AND type = Bug",
			wantErr:   false,
		},
		{
			name:        "expand predefined alias",
			aliasName:   "mine",
			containsJQL: "currentUser()",
			wantErr:     false,
		},
		{
			name:      "expand non-existent alias",
			aliasName: "does-not-exist",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jql, err := service.ExpandAlias(tt.aliasName, workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if tt.wantJQL != "" && jql != tt.wantJQL {
					t.Errorf("ExpandAlias() = %s, want %s", jql, tt.wantJQL)
				}
				if tt.containsJQL != "" && !containsSubstring(jql, tt.containsJQL) {
					t.Errorf("ExpandAlias() = %s, want to contain %s", jql, tt.containsJQL)
				}
			}
		})
	}
}

func TestAliasService_ValidateJQL(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)

	tests := []struct {
		name    string
		jql     string
		wantErr bool
	}{
		{
			name:    "valid JQL",
			jql:     "assignee = currentUser()",
			wantErr: false,
		},
		{
			name:    "empty JQL",
			jql:     "",
			wantErr: true,
		},
		{
			name:    "JQL too long",
			jql:     string(make([]byte, 2001)),
			wantErr: true,
		},
		{
			name:    "balanced parentheses",
			jql:     "status IN (Open, InProgress)",
			wantErr: false,
		},
		{
			name:    "unbalanced parentheses - missing close",
			jql:     "status IN (Open, InProgress",
			wantErr: true,
		},
		{
			name:    "unbalanced parentheses - missing open",
			jql:     "status IN Open, InProgress)",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateJQL(tt.jql)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliasService_ExpandAlias_Recursive(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create a chain of aliases: base -> intermediate -> @base
	service.Create("base", "assignee = currentUser()", "", workspaceID)
	service.Create("intermediate", "@base AND type = Bug", "", workspaceID)
	service.Create("extended", "@intermediate AND priority = High", "", workspaceID)

	tests := []struct {
		name        string
		aliasName   string
		wantContain []string
		wantErr     bool
	}{
		{
			name:        "expand two-level nested alias",
			aliasName:   "intermediate",
			wantContain: []string{"assignee = currentUser()", "type = Bug"},
			wantErr:     false,
		},
		{
			name:        "expand three-level nested alias",
			aliasName:   "extended",
			wantContain: []string{"assignee = currentUser()", "type = Bug", "priority = High"},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jql, err := service.ExpandAlias(tt.aliasName, workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for _, substr := range tt.wantContain {
					if !containsSubstring(jql, substr) {
						t.Errorf("ExpandAlias() = %s, want to contain %s", jql, substr)
					}
				}
			}
		})
	}
}

func TestAliasService_ExpandAlias_CircularReference(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create circular references: a -> b -> c -> a
	service.Create("circular-a", "@circular-b", "", workspaceID)
	service.Create("circular-b", "@circular-c", "", workspaceID)
	service.Create("circular-c", "@circular-a", "", workspaceID)

	// Test circular reference detection
	_, err := service.ExpandAlias("circular-a", workspaceID)
	if err == nil {
		t.Error("ExpandAlias() expected error for circular reference, got nil")
	}
	if err != nil && err != ports.ErrCircularReference {
		// The error might be wrapped, so check if it contains the circular reference error
		if !containsSubstring(err.Error(), "circular") {
			t.Errorf("ExpandAlias() error = %v, want circular reference error", err)
		}
	}
}

func TestAliasService_ExpandAlias_SelfReference(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create self-referencing alias
	service.Create("self-ref", "@self-ref", "", workspaceID)

	// Test self-reference detection
	_, err := service.ExpandAlias("self-ref", workspaceID)
	if err == nil {
		t.Error("ExpandAlias() expected error for self-reference, got nil")
	}
}

func TestAliasService_ExpandAlias_MultipleReferences(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create base aliases
	service.Create("my-tickets", "assignee = currentUser()", "", workspaceID)
	service.Create("high-priority", "priority = High", "", workspaceID)

	// Create alias that references multiple other aliases
	service.Create("urgent-mine", "@my-tickets AND @high-priority", "", workspaceID)

	jql, err := service.ExpandAlias("urgent-mine", workspaceID)
	if err != nil {
		t.Fatalf("ExpandAlias() error = %v", err)
	}

	// Should contain both expanded references
	if !containsSubstring(jql, "assignee = currentUser()") {
		t.Errorf("ExpandAlias() = %s, want to contain 'assignee = currentUser()'", jql)
	}
	if !containsSubstring(jql, "priority = High") {
		t.Errorf("ExpandAlias() = %s, want to contain 'priority = High'", jql)
	}
}

func TestAliasService_ExpandAlias_WithPredefined(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)
	workspaceID := "test-workspace"

	// Create alias that references a predefined alias
	service.Create("my-bugs", "@mine AND type = Bug", "", workspaceID)

	jql, err := service.ExpandAlias("my-bugs", workspaceID)
	if err != nil {
		t.Fatalf("ExpandAlias() error = %v", err)
	}

	// Should expand the predefined alias
	if !containsSubstring(jql, "currentUser()") {
		t.Errorf("ExpandAlias() = %s, want to contain 'currentUser()'", jql)
	}
	if !containsSubstring(jql, "type = Bug") {
		t.Errorf("ExpandAlias() = %s, want to contain 'type = Bug'", jql)
	}
}

func TestAliasService_Create_GlobalAlias(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)

	// Create global alias (empty workspaceID)
	err := service.Create("global-test", "priority = High", "Global alias", "")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify it can be retrieved from any workspace
	alias1, err := service.Get("global-test", "workspace-1")
	if err != nil {
		t.Errorf("Get() from workspace-1 error = %v", err)
	}
	if alias1 == nil || alias1.WorkspaceID != "" {
		t.Error("Expected global alias (empty WorkspaceID)")
	}

	alias2, err := service.Get("global-test", "workspace-2")
	if err != nil {
		t.Errorf("Get() from workspace-2 error = %v", err)
	}
	if alias2 == nil || alias2.WorkspaceID != "" {
		t.Error("Expected global alias (empty WorkspaceID)")
	}
}

func TestAliasService_WorkspaceIsolation(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)

	// Create alias in workspace-1
	err := service.Create("test-alias", "status = Open", "", "workspace-1")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Should be accessible from workspace-1
	alias1, err := service.Get("test-alias", "workspace-1")
	if err != nil {
		t.Errorf("Get() from workspace-1 error = %v", err)
	}
	if alias1 == nil {
		t.Error("Expected alias to be found in workspace-1")
	}

	// Should NOT be accessible from workspace-2 (unless we modify the mock to support this)
	// This test depends on proper workspace isolation in the repository
	alias2, err := service.Get("test-alias", "workspace-2")
	if err == nil && alias2 != nil && alias2.WorkspaceID != "" {
		t.Error("Expected alias to NOT be found in workspace-2 (workspace isolation)")
	}
}

func TestAliasService_Create_LongJQL(t *testing.T) {
	repo := NewMockAliasRepository()
	service := NewAliasService(repo)

	// Test with max valid length (2000 chars)
	longJQL := make([]byte, 2000)
	for i := range longJQL {
		longJQL[i] = 'a'
	}

	err := service.Create("long-jql", string(longJQL), "", "test-workspace")
	if err != nil {
		t.Errorf("Create() with 2000 char JQL error = %v, expected success", err)
	}

	// Test with too long JQL (2001 chars)
	tooLongJQL := make([]byte, 2001)
	for i := range tooLongJQL {
		tooLongJQL[i] = 'a'
	}

	err = service.Create("too-long-jql", string(tooLongJQL), "", "test-workspace")
	if err == nil {
		t.Error("Create() with 2001 char JQL expected error, got nil")
	}
}

// Helper function to check if string contains substring
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
