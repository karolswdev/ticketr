package services

import (
	"errors"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// Test helper to create a ticket with specific fields
func createTestTicket(title, desc, jiraID string, customFields map[string]string, ac []string, tasks []domain.Task) *domain.Ticket {
	return &domain.Ticket{
		Title:              title,
		Description:        desc,
		JiraID:             jiraID,
		CustomFields:       customFields,
		AcceptanceCriteria: ac,
		Tasks:              tasks,
	}
}

func TestNewSyncStrategy(t *testing.T) {
	tests := []struct {
		name         string
		strategyName string
		wantType     string
		wantErr      bool
	}{
		{
			name:         "local-wins strategy",
			strategyName: StrategyLocalWins,
			wantType:     StrategyLocalWins,
			wantErr:      false,
		},
		{
			name:         "remote-wins strategy",
			strategyName: StrategyRemoteWins,
			wantType:     StrategyRemoteWins,
			wantErr:      false,
		},
		{
			name:         "three-way-merge strategy",
			strategyName: StrategyThreeWayMerge,
			wantType:     StrategyThreeWayMerge,
			wantErr:      false,
		},
		{
			name:         "unknown strategy",
			strategyName: "unknown-strategy",
			wantType:     "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := NewSyncStrategy(tt.strategyName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewSyncStrategy() expected error, got nil")
				}
				if !errors.Is(err, ErrUnknownStrategy) {
					t.Errorf("NewSyncStrategy() error = %v, want ErrUnknownStrategy", err)
				}
				return
			}

			if err != nil {
				t.Errorf("NewSyncStrategy() unexpected error = %v", err)
				return
			}

			if strategy == nil {
				t.Errorf("NewSyncStrategy() returned nil strategy")
				return
			}

			if strategy.Name() != tt.wantType {
				t.Errorf("NewSyncStrategy() strategy name = %v, want %v", strategy.Name(), tt.wantType)
			}
		})
	}
}

func TestLocalWinsStrategy_ShouldSync(t *testing.T) {
	strategy := &LocalWinsStrategy{}

	tests := []struct {
		name             string
		localHash        string
		remoteHash       string
		storedLocalHash  string
		storedRemoteHash string
		wantSync         bool
	}{
		{
			name:             "remote changed",
			localHash:        "local-hash-1",
			remoteHash:       "remote-hash-2",
			storedLocalHash:  "local-hash-1",
			storedRemoteHash: "remote-hash-1",
			wantSync:         true,
		},
		{
			name:             "remote unchanged",
			localHash:        "local-hash-2",
			remoteHash:       "remote-hash-1",
			storedLocalHash:  "local-hash-1",
			storedRemoteHash: "remote-hash-1",
			wantSync:         false,
		},
		{
			name:             "both changed",
			localHash:        "local-hash-2",
			remoteHash:       "remote-hash-2",
			storedLocalHash:  "local-hash-1",
			storedRemoteHash: "remote-hash-1",
			wantSync:         true,
		},
		{
			name:             "both unchanged",
			localHash:        "hash-1",
			remoteHash:       "hash-1",
			storedLocalHash:  "hash-1",
			storedRemoteHash: "hash-1",
			wantSync:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strategy.ShouldSync(tt.localHash, tt.remoteHash, tt.storedLocalHash, tt.storedRemoteHash)
			if got != tt.wantSync {
				t.Errorf("LocalWinsStrategy.ShouldSync() = %v, want %v", got, tt.wantSync)
			}
		})
	}
}

func TestLocalWinsStrategy_ResolveConflict(t *testing.T) {
	strategy := &LocalWinsStrategy{}

	local := createTestTicket("Local Title", "Local Description", "PROJ-123", nil, nil, nil)
	remote := createTestTicket("Remote Title", "Remote Description", "PROJ-123", nil, nil, nil)

	tests := []struct {
		name    string
		local   *domain.Ticket
		remote  *domain.Ticket
		wantErr bool
	}{
		{
			name:    "both tickets valid - returns local",
			local:   local,
			remote:  remote,
			wantErr: false,
		},
		{
			name:    "nil local ticket",
			local:   nil,
			remote:  remote,
			wantErr: true,
		},
		{
			name:    "nil remote ticket",
			local:   local,
			remote:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := strategy.ResolveConflict(tt.local, tt.remote)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LocalWinsStrategy.ResolveConflict() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("LocalWinsStrategy.ResolveConflict() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Errorf("LocalWinsStrategy.ResolveConflict() returned nil result")
				return
			}

			// Verify it returns local ticket unchanged
			if result.Title != tt.local.Title {
				t.Errorf("LocalWinsStrategy.ResolveConflict() title = %v, want %v", result.Title, tt.local.Title)
			}
			if result.Description != tt.local.Description {
				t.Errorf("LocalWinsStrategy.ResolveConflict() description = %v, want %v", result.Description, tt.local.Description)
			}
		})
	}
}

func TestRemoteWinsStrategy_ShouldSync(t *testing.T) {
	strategy := &RemoteWinsStrategy{}

	tests := []struct {
		name             string
		localHash        string
		remoteHash       string
		storedLocalHash  string
		storedRemoteHash string
		wantSync         bool
	}{
		{
			name:             "remote changed",
			localHash:        "local-hash-1",
			remoteHash:       "remote-hash-2",
			storedLocalHash:  "local-hash-1",
			storedRemoteHash: "remote-hash-1",
			wantSync:         true,
		},
		{
			name:             "remote unchanged",
			localHash:        "local-hash-2",
			remoteHash:       "remote-hash-1",
			storedLocalHash:  "local-hash-1",
			storedRemoteHash: "remote-hash-1",
			wantSync:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strategy.ShouldSync(tt.localHash, tt.remoteHash, tt.storedLocalHash, tt.storedRemoteHash)
			if got != tt.wantSync {
				t.Errorf("RemoteWinsStrategy.ShouldSync() = %v, want %v", got, tt.wantSync)
			}
		})
	}
}

func TestRemoteWinsStrategy_ResolveConflict(t *testing.T) {
	strategy := &RemoteWinsStrategy{}

	local := createTestTicket("Local Title", "Local Description", "PROJ-123", nil, nil, nil)
	remote := createTestTicket("Remote Title", "Remote Description", "PROJ-123", nil, nil, nil)

	tests := []struct {
		name    string
		local   *domain.Ticket
		remote  *domain.Ticket
		wantErr bool
	}{
		{
			name:    "both tickets valid - returns remote",
			local:   local,
			remote:  remote,
			wantErr: false,
		},
		{
			name:    "nil local ticket",
			local:   nil,
			remote:  remote,
			wantErr: true,
		},
		{
			name:    "nil remote ticket",
			local:   local,
			remote:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := strategy.ResolveConflict(tt.local, tt.remote)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemoteWinsStrategy.ResolveConflict() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("RemoteWinsStrategy.ResolveConflict() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Errorf("RemoteWinsStrategy.ResolveConflict() returned nil result")
				return
			}

			// Verify it returns remote ticket unchanged
			if result.Title != tt.remote.Title {
				t.Errorf("RemoteWinsStrategy.ResolveConflict() title = %v, want %v", result.Title, tt.remote.Title)
			}
			if result.Description != tt.remote.Description {
				t.Errorf("RemoteWinsStrategy.ResolveConflict() description = %v, want %v", result.Description, tt.remote.Description)
			}
		})
	}
}

func TestThreeWayMergeStrategy_ShouldSync(t *testing.T) {
	strategy := &ThreeWayMergeStrategy{}

	tests := []struct {
		name             string
		localHash        string
		remoteHash       string
		storedLocalHash  string
		storedRemoteHash string
		wantSync         bool
	}{
		{
			name:             "remote changed",
			localHash:        "local-hash-1",
			remoteHash:       "remote-hash-2",
			storedLocalHash:  "local-hash-1",
			storedRemoteHash: "remote-hash-1",
			wantSync:         true,
		},
		{
			name:             "remote unchanged",
			localHash:        "local-hash-2",
			remoteHash:       "remote-hash-1",
			storedLocalHash:  "local-hash-1",
			storedRemoteHash: "remote-hash-1",
			wantSync:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strategy.ShouldSync(tt.localHash, tt.remoteHash, tt.storedLocalHash, tt.storedRemoteHash)
			if got != tt.wantSync {
				t.Errorf("ThreeWayMergeStrategy.ShouldSync() = %v, want %v", got, tt.wantSync)
			}
		})
	}
}

func TestThreeWayMergeStrategy_ResolveConflict(t *testing.T) {
	strategy := &ThreeWayMergeStrategy{}

	tests := []struct {
		name        string
		local       *domain.Ticket
		remote      *domain.Ticket
		wantTitle   string
		wantDesc    string
		wantErr     bool
		errContains string
	}{
		{
			name:      "no changes - both identical",
			local:     createTestTicket("Same Title", "Same Description", "PROJ-123", nil, nil, nil),
			remote:    createTestTicket("Same Title", "Same Description", "PROJ-123", nil, nil, nil),
			wantTitle: "Same Title",
			wantDesc:  "Same Description",
			wantErr:   false,
		},
		{
			name:      "local changed title, remote unchanged",
			local:     createTestTicket("New Local Title", "", "PROJ-123", nil, nil, nil),
			remote:    createTestTicket("", "", "PROJ-123", nil, nil, nil),
			wantTitle: "New Local Title",
			wantDesc:  "",
			wantErr:   false,
		},
		{
			name:      "remote changed title, local unchanged",
			local:     createTestTicket("", "", "PROJ-123", nil, nil, nil),
			remote:    createTestTicket("New Remote Title", "", "PROJ-123", nil, nil, nil),
			wantTitle: "New Remote Title",
			wantDesc:  "",
			wantErr:   false,
		},
		{
			name:        "both changed title differently - conflict",
			local:       createTestTicket("Local Title", "", "PROJ-123", nil, nil, nil),
			remote:      createTestTicket("Remote Title", "", "PROJ-123", nil, nil, nil),
			wantErr:     true,
			errContains: "Title",
		},
		{
			name:      "local changed description, remote unchanged",
			local:     createTestTicket("", "New Local Description", "PROJ-123", nil, nil, nil),
			remote:    createTestTicket("", "", "PROJ-123", nil, nil, nil),
			wantTitle: "",
			wantDesc:  "New Local Description",
			wantErr:   false,
		},
		{
			name:      "remote changed description, local unchanged",
			local:     createTestTicket("", "", "PROJ-123", nil, nil, nil),
			remote:    createTestTicket("", "New Remote Description", "PROJ-123", nil, nil, nil),
			wantTitle: "",
			wantDesc:  "New Remote Description",
			wantErr:   false,
		},
		{
			name:        "both changed description differently - conflict",
			local:       createTestTicket("", "Local Description", "PROJ-123", nil, nil, nil),
			remote:      createTestTicket("", "Remote Description", "PROJ-123", nil, nil, nil),
			wantErr:     true,
			errContains: "Description",
		},
		{
			name:      "compatible changes - different fields",
			local:     createTestTicket("Local Title", "", "PROJ-123", nil, nil, nil),
			remote:    createTestTicket("", "Remote Description", "PROJ-123", nil, nil, nil),
			wantTitle: "Local Title",
			wantDesc:  "Remote Description",
			wantErr:   false,
		},
		{
			name:    "nil local ticket",
			local:   nil,
			remote:  createTestTicket("Title", "Desc", "PROJ-123", nil, nil, nil),
			wantErr: true,
		},
		{
			name:    "nil remote ticket",
			local:   createTestTicket("Title", "Desc", "PROJ-123", nil, nil, nil),
			remote:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := strategy.ResolveConflict(tt.local, tt.remote)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ThreeWayMergeStrategy.ResolveConflict() expected error, got nil")
					return
				}
				if tt.errContains != "" && !errors.Is(err, ErrConflictUnresolvable) {
					t.Errorf("ThreeWayMergeStrategy.ResolveConflict() error should be ErrConflictUnresolvable, got %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("ThreeWayMergeStrategy.ResolveConflict() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Errorf("ThreeWayMergeStrategy.ResolveConflict() returned nil result")
				return
			}

			if result.Title != tt.wantTitle {
				t.Errorf("ThreeWayMergeStrategy.ResolveConflict() title = %v, want %v", result.Title, tt.wantTitle)
			}
			if result.Description != tt.wantDesc {
				t.Errorf("ThreeWayMergeStrategy.ResolveConflict() description = %v, want %v", result.Description, tt.wantDesc)
			}
		})
	}
}

func TestThreeWayMergeStrategy_AcceptanceCriteria(t *testing.T) {
	strategy := &ThreeWayMergeStrategy{}

	tests := []struct {
		name    string
		local   []string
		remote  []string
		want    []string
		wantErr bool
	}{
		{
			name:    "both empty",
			local:   []string{},
			remote:  []string{},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "local has criteria, remote empty",
			local:   []string{"Local AC 1", "Local AC 2"},
			remote:  []string{},
			want:    []string{"Local AC 1", "Local AC 2"},
			wantErr: false,
		},
		{
			name:    "remote has criteria, local empty",
			local:   []string{},
			remote:  []string{"Remote AC 1", "Remote AC 2"},
			want:    []string{"Remote AC 1", "Remote AC 2"},
			wantErr: false,
		},
		{
			name:    "both have different criteria - conflict",
			local:   []string{"Local AC"},
			remote:  []string{"Remote AC"},
			wantErr: true,
		},
		{
			name:    "both have same criteria",
			local:   []string{"Same AC"},
			remote:  []string{"Same AC"},
			want:    []string{"Same AC"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			local := createTestTicket("", "", "PROJ-123", nil, tt.local, nil)
			remote := createTestTicket("", "", "PROJ-123", nil, tt.remote, nil)

			result, err := strategy.ResolveConflict(local, remote)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for conflicting acceptance criteria, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error = %v", err)
				return
			}

			if len(result.AcceptanceCriteria) != len(tt.want) {
				t.Errorf("AcceptanceCriteria length = %v, want %v", len(result.AcceptanceCriteria), len(tt.want))
				return
			}

			for i := range tt.want {
				if result.AcceptanceCriteria[i] != tt.want[i] {
					t.Errorf("AcceptanceCriteria[%d] = %v, want %v", i, result.AcceptanceCriteria[i], tt.want[i])
				}
			}
		})
	}
}

func TestThreeWayMergeStrategy_CustomFields(t *testing.T) {
	strategy := &ThreeWayMergeStrategy{}

	tests := []struct {
		name    string
		local   map[string]string
		remote  map[string]string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "both empty",
			local:   map[string]string{},
			remote:  map[string]string{},
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name:    "local has field, remote empty",
			local:   map[string]string{"field1": "local-value"},
			remote:  map[string]string{},
			want:    map[string]string{"field1": "local-value"},
			wantErr: false,
		},
		{
			name:    "remote has field, local empty",
			local:   map[string]string{},
			remote:  map[string]string{"field1": "remote-value"},
			want:    map[string]string{"field1": "remote-value"},
			wantErr: false,
		},
		{
			name:    "different fields - merge",
			local:   map[string]string{"field1": "local-value"},
			remote:  map[string]string{"field2": "remote-value"},
			want:    map[string]string{"field1": "local-value", "field2": "remote-value"},
			wantErr: false,
		},
		{
			name:    "same field, different values - conflict",
			local:   map[string]string{"field1": "local-value"},
			remote:  map[string]string{"field1": "remote-value"},
			wantErr: true,
		},
		{
			name:    "same field, same value - no conflict",
			local:   map[string]string{"field1": "same-value"},
			remote:  map[string]string{"field1": "same-value"},
			want:    map[string]string{"field1": "same-value"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			local := createTestTicket("", "", "PROJ-123", tt.local, nil, nil)
			remote := createTestTicket("", "", "PROJ-123", tt.remote, nil, nil)

			result, err := strategy.ResolveConflict(local, remote)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for conflicting custom fields, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error = %v", err)
				return
			}

			if len(result.CustomFields) != len(tt.want) {
				t.Errorf("CustomFields length = %v, want %v", len(result.CustomFields), len(tt.want))
			}

			for k, wantVal := range tt.want {
				gotVal, exists := result.CustomFields[k]
				if !exists {
					t.Errorf("CustomFields missing key %v", k)
				} else if gotVal != wantVal {
					t.Errorf("CustomFields[%v] = %v, want %v", k, gotVal, wantVal)
				}
			}
		})
	}
}

func TestThreeWayMergeStrategy_Tasks(t *testing.T) {
	strategy := &ThreeWayMergeStrategy{}

	tests := []struct {
		name    string
		local   []domain.Task
		remote  []domain.Task
		want    []domain.Task
		wantErr bool
	}{
		{
			name:    "both empty",
			local:   []domain.Task{},
			remote:  []domain.Task{},
			want:    []domain.Task{},
			wantErr: false,
		},
		{
			name: "local has task, remote empty",
			local: []domain.Task{
				{JiraID: "PROJ-1", Title: "Local Task"},
			},
			remote: []domain.Task{},
			want: []domain.Task{
				{JiraID: "PROJ-1", Title: "Local Task"},
			},
			wantErr: false,
		},
		{
			name:  "remote has task, local empty",
			local: []domain.Task{},
			remote: []domain.Task{
				{JiraID: "PROJ-1", Title: "Remote Task"},
			},
			want: []domain.Task{
				{JiraID: "PROJ-1", Title: "Remote Task"},
			},
			wantErr: false,
		},
		{
			name: "same task, no conflict",
			local: []domain.Task{
				{JiraID: "PROJ-1", Title: "Same Task"},
			},
			remote: []domain.Task{
				{JiraID: "PROJ-1", Title: "Same Task"},
			},
			want: []domain.Task{
				{JiraID: "PROJ-1", Title: "Same Task"},
			},
			wantErr: false,
		},
		{
			name: "same task ID, different titles - conflict",
			local: []domain.Task{
				{JiraID: "PROJ-1", Title: "Local Task"},
			},
			remote: []domain.Task{
				{JiraID: "PROJ-1", Title: "Remote Task"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			local := createTestTicket("", "", "PROJ-123", nil, nil, tt.local)
			remote := createTestTicket("", "", "PROJ-123", nil, nil, tt.remote)

			result, err := strategy.ResolveConflict(local, remote)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for conflicting tasks, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error = %v", err)
				return
			}

			if len(result.Tasks) != len(tt.want) {
				t.Errorf("Tasks length = %v, want %v", len(result.Tasks), len(tt.want))
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("stringSlicesEqual", func(t *testing.T) {
		tests := []struct {
			name string
			a    []string
			b    []string
			want bool
		}{
			{"both nil", nil, nil, true},
			{"both empty", []string{}, []string{}, true},
			{"equal slices", []string{"a", "b"}, []string{"a", "b"}, true},
			{"different lengths", []string{"a"}, []string{"a", "b"}, false},
			{"different values", []string{"a", "b"}, []string{"a", "c"}, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := stringSlicesEqual(tt.a, tt.b)
				if got != tt.want {
					t.Errorf("stringSlicesEqual() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("mapsEqual", func(t *testing.T) {
		tests := []struct {
			name string
			a    map[string]string
			b    map[string]string
			want bool
		}{
			{"both nil", nil, nil, true},
			{"both empty", map[string]string{}, map[string]string{}, true},
			{"equal maps", map[string]string{"k": "v"}, map[string]string{"k": "v"}, true},
			{"different sizes", map[string]string{"k": "v"}, map[string]string{"k": "v", "k2": "v2"}, false},
			{"different values", map[string]string{"k": "v1"}, map[string]string{"k": "v2"}, false},
			{"different keys", map[string]string{"k1": "v"}, map[string]string{"k2": "v"}, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := mapsEqual(tt.a, tt.b)
				if got != tt.want {
					t.Errorf("mapsEqual() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
