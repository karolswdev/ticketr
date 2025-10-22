package commands

import (
	"fmt"
	"sort"
	"strings"
)

// Category represents a command category for grouping.
type Category string

const (
	CategoryView   Category = "View"
	CategoryEdit   Category = "Edit"
	CategorySync   Category = "Sync"
	CategoryNav    Category = "Navigation"
	CategorySystem Category = "System"
)

// Command represents a single executable command with metadata.
type Command struct {
	Name        string
	Description string
	Keybinding  string // e.g., "p", "Ctrl+P", "F2"
	Category    Category
	Handler     func() error
}

// Registry manages all available commands.
type Registry struct {
	commands map[string]*Command
}

// NewRegistry creates a new command registry.
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]*Command),
	}
}

// Register adds a command to the registry.
func (r *Registry) Register(cmd *Command) error {
	if cmd.Name == "" {
		return fmt.Errorf("command name is required")
	}
	if cmd.Handler == nil {
		return fmt.Errorf("command handler is required")
	}
	r.commands[cmd.Name] = cmd
	return nil
}

// Get retrieves a command by name.
func (r *Registry) Get(name string) (*Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// All returns all registered commands.
func (r *Registry) All() []*Command {
	cmds := make([]*Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}

	// Sort by category, then by name
	sort.Slice(cmds, func(i, j int) bool {
		if cmds[i].Category != cmds[j].Category {
			return cmds[i].Category < cmds[j].Category
		}
		return cmds[i].Name < cmds[j].Name
	})

	return cmds
}

// ByCategory returns commands grouped by category.
func (r *Registry) ByCategory() map[Category][]*Command {
	result := make(map[Category][]*Command)

	for _, cmd := range r.commands {
		result[cmd.Category] = append(result[cmd.Category], cmd)
	}

	// Sort commands within each category
	for cat := range result {
		sort.Slice(result[cat], func(i, j int) bool {
			return result[cat][i].Name < result[cat][j].Name
		})
	}

	return result
}

// Search performs fuzzy search on command names and descriptions.
func (r *Registry) Search(query string) []*Command {
	if query == "" {
		return r.All()
	}

	query = strings.ToLower(query)
	var results []*Command

	for _, cmd := range r.commands {
		// Simple fuzzy matching: check if query is substring of name or description
		name := strings.ToLower(cmd.Name)
		desc := strings.ToLower(cmd.Description)

		if strings.Contains(name, query) || strings.Contains(desc, query) {
			results = append(results, cmd)
		}
	}

	// Sort results by relevance (name matches first, then description matches)
	sort.Slice(results, func(i, j int) bool {
		iNameMatch := strings.Contains(strings.ToLower(results[i].Name), query)
		jNameMatch := strings.Contains(strings.ToLower(results[j].Name), query)

		if iNameMatch != jNameMatch {
			return iNameMatch // Name matches come first
		}

		return results[i].Name < results[j].Name
	})

	return results
}

// Execute executes a command by name.
func (r *Registry) Execute(name string) error {
	cmd, exists := r.Get(name)
	if !exists {
		return fmt.Errorf("command not found: %s", name)
	}
	return cmd.Handler()
}

// Count returns the total number of registered commands.
func (r *Registry) Count() int {
	return len(r.commands)
}
