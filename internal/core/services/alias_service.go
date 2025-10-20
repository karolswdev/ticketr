package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// AliasService provides business logic for managing JQL aliases.
// It handles alias expansion, validation, and cycle detection.
// Optimized with caching for predefined aliases and memoization for expansion.
type AliasService struct {
	repo             ports.AliasRepository
	predefinedCache  map[string]*domain.JQLAlias // Cache for predefined aliases
	expansionCache   map[string]string           // Cache for expanded JQL
	expansionCacheMu sync.RWMutex                // Protect expansion cache
	cacheMu          sync.RWMutex                // Protect predefined cache
}

// NewAliasService creates a new AliasService instance.
func NewAliasService(repo ports.AliasRepository) *AliasService {
	s := &AliasService{
		repo:            repo,
		predefinedCache: make(map[string]*domain.JQLAlias),
		expansionCache:  make(map[string]string),
	}

	// Pre-populate predefined aliases cache
	s.initPredefinedCache()

	return s
}

// initPredefinedCache initializes the cache with predefined aliases.
func (s *AliasService) initPredefinedCache() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	for name, alias := range domain.PredefinedAliases {
		aliasCopy := alias // Make a copy
		s.predefinedCache[name] = &aliasCopy
	}
}

// Create creates a new JQL alias.
// It validates the alias and checks for duplicates before creating.
func (s *AliasService) Create(name, jql, description, workspaceID string) error {
	// Validate alias name
	if err := domain.ValidateAliasName(name); err != nil {
		return fmt.Errorf("invalid alias name '%s': %w. Alias names must contain only letters, numbers, hyphens, and underscores", name, err)
	}

	// Check if this is a predefined alias name
	if domain.IsPredefinedAlias(name) {
		return fmt.Errorf("cannot create alias '%s': name is reserved for predefined system alias. Choose a different name or use 'ticketr alias list' to see reserved names", name)
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
		return fmt.Errorf("alias validation failed for '%s': %w. Check JQL syntax and field values", name, err)
	}

	// Persist alias
	if err := s.repo.Create(alias); err != nil {
		return fmt.Errorf("failed to save alias '%s' to database: %w", name, err)
	}

	// Invalidate expansion cache since new alias may affect expansions
	s.invalidateExpansionCache()

	return nil
}

// invalidateExpansionCache clears the expansion cache.
func (s *AliasService) invalidateExpansionCache() {
	s.expansionCacheMu.Lock()
	defer s.expansionCacheMu.Unlock()
	s.expansionCache = make(map[string]string)
}

// Get retrieves an alias by name.
// It first checks for user-defined aliases, then falls back to cached predefined aliases.
func (s *AliasService) Get(name, workspaceID string) (*domain.JQLAlias, error) {
	// Try to get user-defined alias first
	alias, err := s.repo.GetByName(name, workspaceID)
	if err == nil {
		return alias, nil
	}

	// If not found, check predefined aliases cache (faster than domain lookup)
	s.cacheMu.RLock()
	if predefined, ok := s.predefinedCache[name]; ok {
		s.cacheMu.RUnlock()
		return predefined, nil
	}
	s.cacheMu.RUnlock()

	return nil, ports.ErrAliasNotFound
}

// List returns all aliases for a workspace (including predefined ones).
func (s *AliasService) List(workspaceID string) ([]domain.JQLAlias, error) {
	// Get user-defined aliases
	userAliases, err := s.repo.List(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list aliases from database for workspace '%s': %w", workspaceID, err)
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
		return fmt.Errorf("failed to delete alias '%s': %w. Verify alias exists in current workspace", name, err)
	}

	// Invalidate expansion cache since alias was removed
	s.invalidateExpansionCache()

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
		return fmt.Errorf("alias '%s' not found in current workspace: %w. Use 'ticketr alias list' to see available aliases", name, err)
	}

	// Update fields
	alias.JQL = jql
	alias.Description = description
	alias.UpdatedAt = time.Now()

	// Validate updated alias
	if err := alias.Validate(); err != nil {
		return fmt.Errorf("updated alias validation failed for '%s': %w. Check JQL syntax", name, err)
	}

	// Persist changes
	if err := s.repo.Update(alias); err != nil {
		return fmt.Errorf("failed to save updated alias '%s' to database: %w", name, err)
	}

	// Invalidate expansion cache since alias content changed
	s.invalidateExpansionCache()

	return nil
}

// ExpandAlias expands an alias name to its JQL query.
// It handles recursive expansion and detects circular references.
// Uses memoization to cache expansion results for better performance.
func (s *AliasService) ExpandAlias(name, workspaceID string) (string, error) {
	// Check expansion cache first
	cacheKey := workspaceID + ":" + name
	s.expansionCacheMu.RLock()
	if expanded, ok := s.expansionCache[cacheKey]; ok {
		s.expansionCacheMu.RUnlock()
		return expanded, nil
	}
	s.expansionCacheMu.RUnlock()

	// Track visited aliases to detect cycles
	visited := make(map[string]bool)
	expanded, err := s.expandAliasRecursive(name, workspaceID, visited)
	if err != nil {
		return "", err
	}

	// Cache the result for future lookups
	s.expansionCacheMu.Lock()
	s.expansionCache[cacheKey] = expanded
	s.expansionCacheMu.Unlock()

	return expanded, nil
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
		return "", fmt.Errorf("alias '%s' not found: %w. Use 'ticketr alias list' to see available aliases", name, err)
	}

	// Check if the JQL contains other alias references
	// Alias references are in the format: @alias_name
	jql := alias.JQL
	if !strings.Contains(jql, "@") {
		// No alias references, return as-is
		return jql, nil
	}

	// Expand any alias references in the JQL efficiently
	// Use strings.Builder for better performance with multiple replacements
	var result strings.Builder
	result.Grow(len(jql) * 2) // Pre-allocate assuming some expansion

	i := 0
	for i < len(jql) {
		if jql[i] == '@' {
			// Find the end of the alias name
			start := i + 1
			end := start
			for end < len(jql) && isAliasNameChar(jql[end]) {
				end++
			}

			if end > start {
				aliasName := jql[start:end]

				// Recursively expand this alias
				expandedAlias, err := s.expandAliasRecursive(aliasName, workspaceID, visited)
				if err != nil {
					return "", fmt.Errorf("failed to expand alias reference '@%s' in JQL: %w. Check that referenced alias exists", aliasName, err)
				}

				// Write the expanded alias wrapped in parentheses
				result.WriteString("(")
				result.WriteString(expandedAlias)
				result.WriteString(")")

				i = end
				continue
			}
		}

		// Regular character, just copy it
		result.WriteByte(jql[i])
		i++
	}

	return result.String(), nil
}

// isAliasNameChar returns true if the character is valid in an alias name.
func isAliasNameChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') || c == '_' || c == '-'
}

// ValidateJQL performs basic validation on a JQL query.
// This is a simple validation to catch obvious errors.
func (s *AliasService) ValidateJQL(jql string) error {
	if jql == "" {
		return fmt.Errorf("JQL query cannot be empty. Provide a valid Jira Query Language expression")
	}

	if len(jql) > 2000 {
		return fmt.Errorf("JQL query too long (%d characters, max 2000). Simplify query or use alias references", len(jql))
	}

	// Basic syntax check: balanced parentheses
	depth := 0
	for _, ch := range jql {
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
			if depth < 0 {
				return fmt.Errorf("unbalanced parentheses in JQL query: found closing ')' without matching opening '('")
			}
		}
	}
	if depth != 0 {
		return fmt.Errorf("unbalanced parentheses in JQL query: %d unclosed '(' parentheses. Check query syntax", depth)
	}

	return nil
}
