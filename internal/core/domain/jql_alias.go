package domain

import (
	"fmt"
	"regexp"
	"time"
)

// JQLAlias represents a named, reusable JQL query.
// Aliases can be predefined (system) or user-defined (custom).
type JQLAlias struct {
	ID           string
	Name         string
	JQL          string
	Description  string
	IsPredefined bool   // true for system aliases, false for user aliases
	WorkspaceID  string // empty string for global aliases
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var (
	// aliasNamePattern validates alias names (alphanumeric, hyphens, underscores).
	aliasNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	// PredefinedAliases contains built-in JQL aliases available to all users.
	PredefinedAliases = map[string]JQLAlias{
		"mine": {
			Name:         "mine",
			JQL:          "assignee = currentUser()",
			Description:  "Tickets assigned to the current user",
			IsPredefined: true,
		},
		"sprint": {
			Name:         "sprint",
			JQL:          "sprint in openSprints()",
			Description:  "Tickets in currently active sprints",
			IsPredefined: true,
		},
		"blocked": {
			Name:         "blocked",
			JQL:          "status = Blocked OR labels in (blocked)",
			Description:  "Tickets that are blocked or have blocked label",
			IsPredefined: true,
		},
	}
)

// Validate checks if the alias has valid values.
func (a *JQLAlias) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("alias name cannot be empty")
	}

	if !aliasNamePattern.MatchString(a.Name) {
		return fmt.Errorf("alias name must contain only alphanumeric characters, hyphens, and underscores")
	}

	if len(a.Name) > 64 {
		return fmt.Errorf("alias name must be 64 characters or less")
	}

	if a.JQL == "" {
		return fmt.Errorf("JQL query cannot be empty")
	}

	if len(a.JQL) > 2000 {
		return fmt.Errorf("JQL query must be 2000 characters or less")
	}

	// Predefined aliases cannot be associated with a workspace
	if a.IsPredefined && a.WorkspaceID != "" {
		return fmt.Errorf("predefined aliases cannot be workspace-specific")
	}

	return nil
}

// ValidateAliasName checks if an alias name is valid.
func ValidateAliasName(name string) error {
	if name == "" {
		return fmt.Errorf("alias name cannot be empty")
	}

	if !aliasNamePattern.MatchString(name) {
		return fmt.Errorf("alias name must contain only alphanumeric characters, hyphens, and underscores")
	}

	if len(name) > 64 {
		return fmt.Errorf("alias name must be 64 characters or less")
	}

	return nil
}

// IsPredefinedAlias checks if an alias name is a predefined system alias.
func IsPredefinedAlias(name string) bool {
	_, exists := PredefinedAliases[name]
	return exists
}

// GetPredefinedAlias retrieves a predefined alias by name.
// Returns nil if the alias doesn't exist.
func GetPredefinedAlias(name string) *JQLAlias {
	if alias, exists := PredefinedAliases[name]; exists {
		// Return a copy to prevent modification of the predefined alias
		aliasCopy := alias
		return &aliasCopy
	}
	return nil
}
