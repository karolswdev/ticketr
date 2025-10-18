package domain

import (
	"fmt"
	"regexp"
	"time"
)

// Workspace represents a Jira project workspace configuration.
// Each workspace maintains independent tickets, sync state, and credentials.
type Workspace struct {
	ID          string
	Name        string
	JiraURL     string
	ProjectKey  string
	Credentials CredentialRef
	IsDefault   bool
	LastUsed    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CredentialRef holds a reference to credentials stored in the OS keychain.
// This ensures credentials are never stored in the database or logs.
type CredentialRef struct {
	KeychainID string
	ServiceID  string
}

// WorkspaceConfig contains the configuration needed to create a workspace.
type WorkspaceConfig struct {
	JiraURL    string
	ProjectKey string
	Username   string
	APIToken   string
}

var (
	// workspaceNamePattern validates workspace names (alphanumeric, hyphens, underscores).
	workspaceNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// Validate checks if the workspace has valid values.
func (w *Workspace) Validate() error {
	if w.ID == "" {
		return fmt.Errorf("workspace ID cannot be empty")
	}

	if w.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if !workspaceNamePattern.MatchString(w.Name) {
		return fmt.Errorf("name must contain only alphanumeric characters, hyphens, and underscores")
	}

	if len(w.Name) > 64 {
		return fmt.Errorf("name must be 64 characters or less")
	}

	if w.JiraURL == "" {
		return fmt.Errorf("jira_url cannot be empty")
	}

	if w.ProjectKey == "" {
		return fmt.Errorf("project_key cannot be empty")
	}

	return nil
}

// ValidateConfig checks if a workspace configuration is valid.
func ValidateConfig(config WorkspaceConfig) error {
	if config.JiraURL == "" {
		return fmt.Errorf("Jira URL is required")
	}

	if config.ProjectKey == "" {
		return fmt.Errorf("project key is required")
	}

	if config.Username == "" {
		return fmt.Errorf("username is required")
	}

	if config.APIToken == "" {
		return fmt.Errorf("API token is required")
	}

	return nil
}

// ValidateWorkspaceName checks if a workspace name is valid.
func ValidateWorkspaceName(name string) error {
	if name == "" {
		return fmt.Errorf("workspace name cannot be empty")
	}

	if !workspaceNamePattern.MatchString(name) {
		return fmt.Errorf("workspace name must contain only alphanumeric characters, hyphens, and underscores")
	}

	if len(name) > 64 {
		return fmt.Errorf("workspace name must be 64 characters or less")
	}

	return nil
}

// Touch updates the LastUsed timestamp to the current time.
func (w *Workspace) Touch() {
	w.LastUsed = time.Now()
}

// SetDefault marks this workspace as the default.
func (w *Workspace) SetDefault(isDefault bool) {
	w.IsDefault = isDefault
}
