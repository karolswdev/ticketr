package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// AliasService provides business logic for managing JQL aliases.
// It handles alias expansion, validation, and cycle detection.
type AliasService struct {
	repo ports.AliasRepository
}

// NewAliasService creates a new AliasService instance.
func NewAliasService(repo ports.AliasRepository) *AliasService {
	return &AliasService{
		repo: repo,
	}
}

// Create creates a new JQL alias.
// It validates the alias and checks for duplicates before creating.
func (s *AliasService) Create(name, jql, description, workspaceID string) error {
	// Validate alias name
	if err := domain.ValidateAliasName(name); err != nil {
		return fmt.Errorf("invalid alias name: %w", err)
	}

	// Check if this is a predefined alias name
	if domain.IsPredefinedAlias(name) {
		return fmt.Errorf("cannot create alias '%s': this name is reserved for a predefined alias", name)
	}

	// Check if alias already exists
	existing, err := s.repo.GetByName(name, workspaceID)
	if err == nil && existing != nil {
		return ports.ErrAliasExists
	}

	// Generate unique ID
	id := uuid.New().String()

	// Create alias
	now := time.Now()
	alias := &domain.JQLAlias{
		ID:           id,
		Name:         name,
		JQL:          jql,
		Description:  description,
		IsPredefined: false,
		WorkspaceID:  workspaceID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Validate alias
	if err := alias.Validate(); err != nil {
		return fmt.Errorf("alias validation failed: %w", err)
	}

	// Persist alias
	if err := s.repo.Create(alias); err != nil {
		return fmt.Errorf("failed to create alias: %w", err)
	}

	return nil
}

// Get retrieves an alias by name.
// It first checks for user-defined aliases, then falls back to predefined aliases.
func (s *AliasService) Get(name, workspaceID string) (*domain.JQLAlias, error) {
	// Try to get user-defined alias first
	alias, err := s.repo.GetByName(name, workspaceID)
	if err == nil {
		return alias, nil
	}

	// If not found, check predefined aliases
	if predefined := domain.GetPredefinedAlias(name); predefined != nil {
		return predefined, nil
	}

	return nil, ports.ErrAliasNotFound
}

// List returns all aliases for a workspace (including predefined ones).
func (s *AliasService) List(workspaceID string) ([]domain.JQLAlias, error) {
	// Get user-defined aliases
	userAliases, err := s.repo.List(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list aliases: %w", err)
	}

	// Combine with predefined aliases
	result := make([]domain.JQLAlias, 0, len(userAliases)+len(domain.PredefinedAliases))

	// Add predefined aliases first
	for _, predefined := range domain.PredefinedAliases {
		result = append(result, predefined)
	}

	// Add user-defined aliases (these override predefined ones with same name)
	for _, alias := range userAliases {
		result = append(result, *alias)
	}

	return result, nil
}

// Delete removes an alias by name.
// Predefined aliases cannot be deleted.
func (s *AliasService) Delete(name, workspaceID string) error {
	// Check if this is a predefined alias
	if domain.IsPredefinedAlias(name) {
		return ports.ErrCannotModifyPredefined
	}

	// Delete the alias
	if err := s.repo.DeleteByName(name, workspaceID); err != nil {
		return fmt.Errorf("failed to delete alias: %w", err)
	}

	return nil
}

// Update updates an existing alias.
// Predefined aliases cannot be updated.
func (s *AliasService) Update(name, jql, description, workspaceID string) error {
	// Check if this is a predefined alias
	if domain.IsPredefinedAlias(name) {
		return ports.ErrCannotModifyPredefined
	}

	// Get existing alias
	alias, err := s.repo.GetByName(name, workspaceID)
	if err != nil {
		return fmt.Errorf("alias not found: %w", err)
	}

	// Update fields
	alias.JQL = jql
	alias.Description = description
	alias.UpdatedAt = time.Now()

	// Validate updated alias
	if err := alias.Validate(); err != nil {
		return fmt.Errorf("alias validation failed: %w", err)
	}

	// Persist changes
	if err := s.repo.Update(alias); err != nil {
		return fmt.Errorf("failed to update alias: %w", err)
	}

	return nil
}

// ExpandAlias expands an alias name to its JQL query.
// It handles recursive expansion and detects circular references.
func (s *AliasService) ExpandAlias(name, workspaceID string) (string, error) {
	// Track visited aliases to detect cycles
	visited := make(map[string]bool)
	return s.expandAliasRecursive(name, workspaceID, visited)
}

// expandAliasRecursive is the recursive implementation of ExpandAlias.
func (s *AliasService) expandAliasRecursive(name, workspaceID string, visited map[string]bool) (string, error) {
	// Check for circular reference
	if visited[name] {
		return "", ports.ErrCircularReference
	}
	visited[name] = true

	// Get the alias
	alias, err := s.Get(name, workspaceID)
	if err != nil {
		return "", fmt.Errorf("alias '%s' not found: %w", name, err)
	}

	// Check if the JQL contains other alias references
	// Alias references are in the format: @alias_name
	jql := alias.JQL
	if !strings.Contains(jql, "@") {
		// No alias references, return as-is
		return jql, nil
	}

	// Expand any alias references in the JQL
	// This is a simple implementation that looks for @word patterns
	expanded := jql
	words := strings.Fields(jql)
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			// This is an alias reference
			aliasName := strings.TrimPrefix(word, "@")
			// Remove any trailing punctuation
			aliasName = strings.TrimRight(aliasName, ".,;:)\"")

			// Recursively expand this alias
			expandedAlias, err := s.expandAliasRecursive(aliasName, workspaceID, visited)
			if err != nil {
				return "", fmt.Errorf("failed to expand alias reference '@%s': %w", aliasName, err)
			}

			// Replace the alias reference with its expanded JQL
			expanded = strings.Replace(expanded, "@"+aliasName, "("+expandedAlias+")", 1)
		}
	}

	return expanded, nil
}

// ValidateJQL performs basic validation on a JQL query.
// This is a simple validation to catch obvious errors.
func (s *AliasService) ValidateJQL(jql string) error {
	if jql == "" {
		return fmt.Errorf("JQL query cannot be empty")
	}

	if len(jql) > 2000 {
		return fmt.Errorf("JQL query too long (max 2000 characters)")
	}

	// Basic syntax check: balanced parentheses
	depth := 0
	for _, ch := range jql {
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
			if depth < 0 {
				return fmt.Errorf("unbalanced parentheses in JQL query")
			}
		}
	}
	if depth != 0 {
		return fmt.Errorf("unbalanced parentheses in JQL query")
	}

	return nil
}
