package domain

import (
	"fmt"
	"regexp"
	"time"
)

// CredentialProfile represents a reusable credential configuration for Jira access.
// Profiles can be shared across multiple workspaces to reduce credential duplication.
type CredentialProfile struct {
	ID          string
	Name        string
	JiraURL     string
	Username    string
	KeychainRef CredentialRef
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CredentialProfileInput contains the data needed to create a credential profile.
type CredentialProfileInput struct {
	Name     string
	JiraURL  string
	Username string
	APIToken string
}

var (
	// credentialProfileNamePattern validates credential profile names (alphanumeric, hyphens, underscores, spaces).
	credentialProfileNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_\- ]+$`)
)

// Validate checks if the credential profile has valid values.
func (cp *CredentialProfile) Validate() error {
	if cp.ID == "" {
		return fmt.Errorf("credential profile ID cannot be empty")
	}

	if cp.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if !credentialProfileNamePattern.MatchString(cp.Name) {
		return fmt.Errorf("name must contain only alphanumeric characters, hyphens, underscores, and spaces")
	}

	if len(cp.Name) > 64 {
		return fmt.Errorf("name must be 64 characters or less")
	}

	if cp.JiraURL == "" {
		return fmt.Errorf("jira_url cannot be empty")
	}

	if cp.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if cp.KeychainRef.KeychainID == "" || cp.KeychainRef.ServiceID == "" {
		return fmt.Errorf("keychain_ref cannot be empty")
	}

	return nil
}

// ValidateCredentialProfileInput checks if a credential profile input is valid.
func ValidateCredentialProfileInput(input CredentialProfileInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}

	if !credentialProfileNamePattern.MatchString(input.Name) {
		return fmt.Errorf("name must contain only alphanumeric characters, hyphens, underscores, and spaces")
	}

	if len(input.Name) > 64 {
		return fmt.Errorf("name must be 64 characters or less")
	}

	if input.JiraURL == "" {
		return fmt.Errorf("Jira URL is required")
	}

	if input.Username == "" {
		return fmt.Errorf("username is required")
	}

	if input.APIToken == "" {
		return fmt.Errorf("API token is required")
	}

	return nil
}

// ValidateCredentialProfileName checks if a credential profile name is valid.
func ValidateCredentialProfileName(name string) error {
	if name == "" {
		return fmt.Errorf("credential profile name cannot be empty")
	}

	if !credentialProfileNamePattern.MatchString(name) {
		return fmt.Errorf("credential profile name must contain only alphanumeric characters, hyphens, underscores, and spaces")
	}

	if len(name) > 64 {
		return fmt.Errorf("credential profile name must be 64 characters or less")
	}

	return nil
}

// Touch updates the UpdatedAt timestamp to the current time.
func (cp *CredentialProfile) Touch() {
	cp.UpdatedAt = time.Now()
}
