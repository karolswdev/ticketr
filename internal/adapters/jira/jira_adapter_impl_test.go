package jira

import (
	"context"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewJiraAdapterFromConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *domain.WorkspaceConfig
		wantErr     bool
		errContains string
	}{
		{
			name:        "nil config",
			config:      nil,
			wantErr:     true,
			errContains: "workspace configuration is required",
		},
		{
			name: "missing JiraURL",
			config: &domain.WorkspaceConfig{
				Username:   "test@example.com",
				APIToken:   "token123",
				ProjectKey: "TEST",
			},
			wantErr:     true,
			errContains: "Jira URL is required",
		},
		{
			name: "missing Username",
			config: &domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				APIToken:   "token123",
				ProjectKey: "TEST",
			},
			wantErr:     true,
			errContains: "username is required",
		},
		{
			name: "missing APIToken",
			config: &domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				Username:   "test@example.com",
				ProjectKey: "TEST",
			},
			wantErr:     true,
			errContains: "API token is required",
		},
		{
			name: "missing ProjectKey",
			config: &domain.WorkspaceConfig{
				JiraURL:  "https://test.atlassian.net",
				Username: "test@example.com",
				APIToken: "token123",
			},
			wantErr:     true,
			errContains: "project key is required",
		},
		{
			name: "valid config",
			config: &domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net",
				Username:   "test@example.com",
				APIToken:   "token123",
				ProjectKey: "TEST",
			},
			wantErr: false,
		},
		{
			name: "valid config with trailing slash",
			config: &domain.WorkspaceConfig{
				JiraURL:    "https://test.atlassian.net/",
				Username:   "test@example.com",
				APIToken:   "token123",
				ProjectKey: "TEST",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, err := NewJiraAdapterFromConfig(tt.config, nil)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, adapter)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, adapter)

				// Verify adapter is properly initialized
				v2Adapter, ok := adapter.(*JiraAdapter)
				assert.True(t, ok)
				assert.NotNil(t, v2Adapter.client)
				assert.Equal(t, tt.config.ProjectKey, v2Adapter.projectKey)
				assert.NotNil(t, v2Adapter.fieldMappings)
			}
		})
	}
}

func TestJiraAdapter_BuildDescription(t *testing.T) {
	adapter := &JiraAdapter{}

	tests := []struct {
		name               string
		description        string
		acceptanceCriteria []string
		want               string
	}{
		{
			name:               "no acceptance criteria",
			description:        "Simple description",
			acceptanceCriteria: nil,
			want:               "Simple description",
		},
		{
			name:               "with acceptance criteria",
			description:        "Description text",
			acceptanceCriteria: []string{"AC1", "AC2"},
			want:               "Description text\n\nh3. Acceptance Criteria\n* AC1\n* AC2\n",
		},
		{
			name:               "empty acceptance criteria slice",
			description:        "Description",
			acceptanceCriteria: []string{},
			want:               "Description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.buildDescription(tt.description, tt.acceptanceCriteria)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJiraAdapter_GetJiraFieldID(t *testing.T) {
	adapter := &JiraAdapter{
		fieldMappings: map[string]interface{}{
			"Story Points": map[string]interface{}{
				"id":   "customfield_10010",
				"type": "number",
			},
			"Sprint": "customfield_10020",
			"Labels": "labels",
		},
	}

	tests := []struct {
		name      string
		humanName string
		want      string
	}{
		{
			name:      "complex mapping",
			humanName: "Story Points",
			want:      "customfield_10010",
		},
		{
			name:      "simple mapping",
			humanName: "Sprint",
			want:      "customfield_10020",
		},
		{
			name:      "standard field",
			humanName: "Labels",
			want:      "labels",
		},
		{
			name:      "unknown field",
			humanName: "Unknown",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.getJiraFieldID(tt.humanName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJiraAdapter_ConvertFieldValue(t *testing.T) {
	adapter := &JiraAdapter{
		fieldMappings: map[string]interface{}{
			"Story Points": map[string]interface{}{
				"id":   "customfield_10010",
				"type": "number",
			},
			"Labels": "labels",
			"Text":   "customfield_10030",
		},
	}

	tests := []struct {
		name      string
		fieldName string
		value     string
		want      interface{}
	}{
		{
			name:      "number field - valid",
			fieldName: "Story Points",
			value:     "5",
			want:      float64(5),
		},
		{
			name:      "number field - decimal",
			fieldName: "Story Points",
			value:     "3.5",
			want:      float64(3.5),
		},
		{
			name:      "number field - invalid",
			fieldName: "Story Points",
			value:     "abc",
			want:      "abc",
		},
		{
			name:      "labels field - single",
			fieldName: "Labels",
			value:     "bug",
			want:      []string{"bug"},
		},
		{
			name:      "labels field - multiple",
			fieldName: "Labels",
			value:     "bug, feature, urgent",
			want:      []string{"bug", "feature", "urgent"},
		},
		{
			name:      "labels field - empty",
			fieldName: "Labels",
			value:     "",
			want:      []string{},
		},
		{
			name:      "text field",
			fieldName: "Text",
			value:     "some text",
			want:      "some text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.convertFieldValue(tt.fieldName, tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJiraAdapter_CreateReverseFieldMapping(t *testing.T) {
	adapter := &JiraAdapter{
		fieldMappings: map[string]interface{}{
			"Story Points": map[string]interface{}{
				"id":   "customfield_10010",
				"type": "number",
			},
			"Sprint": "customfield_10020",
			"Labels": "labels",
		},
	}

	reverse := adapter.createReverseFieldMapping()

	expected := map[string]string{
		"customfield_10010": "Story Points",
		"customfield_10020": "Sprint",
		"labels":            "Labels",
	}

	assert.Equal(t, expected, reverse)
}

func TestJiraAdapter_FormatFieldValue(t *testing.T) {
	adapter := &JiraAdapter{}

	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "string value",
			value: "test",
			want:  "test",
		},
		{
			name:  "float64 value",
			value: float64(5),
			want:  "5",
		},
		{
			name:  "float64 decimal",
			value: float64(3.5),
			want:  "3.5",
		},
		{
			name:  "array of strings",
			value: []interface{}{"a", "b", "c"},
			want:  "a, b, c",
		},
		{
			name: "array of objects with name",
			value: []interface{}{
				map[string]interface{}{"name": "Bug"},
				map[string]interface{}{"name": "Feature"},
			},
			want: "Bug, Feature",
		},
		{
			name: "object with name",
			value: map[string]interface{}{
				"name": "John Doe",
			},
			want: "John Doe",
		},
		{
			name: "object with displayName",
			value: map[string]interface{}{
				"displayName": "Jane Smith",
			},
			want: "Jane Smith",
		},
		{
			name:  "empty array",
			value: []interface{}{},
			want:  "",
		},
		{
			name:  "nil value",
			value: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.formatFieldValue(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJiraAdapter_SearchTickets_ContextCancellation(t *testing.T) {
	// This test verifies that SearchTickets respects context cancellation
	// We can't test against real Jira, but we can verify the context handling

	adapter := &JiraAdapter{
		projectKey: "TEST",
		fieldMappings: map[string]interface{}{
			"Sprint": "customfield_10020",
		},
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := adapter.SearchTickets(ctx, "TEST", "", nil)

	// Should return context.Canceled error
	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestJiraAdapter_ConvertToDomainTicket(t *testing.T) {
	adapter := &JiraAdapter{
		fieldMappings: map[string]interface{}{
			"Story Points": "customfield_10010",
			"Sprint":       "customfield_10020",
		},
	}

	// Note: We can't import jira types directly in tests without mocking,
	// but we can test the reverse mapping and formatting logic
	reverse := adapter.createReverseFieldMapping()
	assert.Equal(t, "Story Points", reverse["customfield_10010"])
	assert.Equal(t, "Sprint", reverse["customfield_10020"])
}

func TestGetDefaultFieldMappings(t *testing.T) {
	mappings := getDefaultFieldMappings()

	// Verify standard fields exist
	assert.Contains(t, mappings, "Type")
	assert.Contains(t, mappings, "Project")
	assert.Contains(t, mappings, "Summary")
	assert.Contains(t, mappings, "Description")

	// Verify custom fields have proper structure
	storyPointsMapping, ok := mappings["Story Points"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "customfield_10010", storyPointsMapping["id"])
	assert.Equal(t, "number", storyPointsMapping["type"])

	// Verify simple mappings
	assert.Equal(t, "customfield_10020", mappings["Sprint"])
}
