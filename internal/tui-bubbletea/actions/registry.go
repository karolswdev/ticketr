package actions

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Registry manages all registered actions
type Registry struct {
	actions   map[ActionID]*Action
	byContext map[Context][]*Action
	byKey     map[string][]*Action
	mutex     sync.RWMutex
}

// NewRegistry creates a new action registry
func NewRegistry() *Registry {
	return &Registry{
		actions:   make(map[ActionID]*Action),
		byContext: make(map[Context][]*Action),
		byKey:     make(map[string][]*Action),
	}
}

// Register adds an action to the registry
func (r *Registry) Register(action *Action) error {
	if action.ID == "" {
		return fmt.Errorf("action ID is required")
	}
	if action.Execute == nil {
		return fmt.Errorf("action execute function is required")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check for duplicate ID
	if _, exists := r.actions[action.ID]; exists {
		return fmt.Errorf("action ID already registered: %s", action.ID)
	}

	// Store action
	r.actions[action.ID] = action

	// Index by context
	for _, ctx := range action.Contexts {
		r.byContext[ctx] = append(r.byContext[ctx], action)
	}

	// Index by keybinding
	for _, keyPattern := range action.Keybindings {
		keyStr := keyPattern.String()
		r.byKey[keyStr] = append(r.byKey[keyStr], action)
	}

	return nil
}

// Unregister removes an action from the registry (for plugins)
func (r *Registry) Unregister(id ActionID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	action, exists := r.actions[id]
	if !exists {
		return fmt.Errorf("action not found: %s", id)
	}

	// Remove from main map
	delete(r.actions, id)

	// Remove from context index
	for _, ctx := range action.Contexts {
		r.removeFromSlice(r.byContext[ctx], action)
	}

	// Remove from key index
	for _, keyPattern := range action.Keybindings {
		keyStr := keyPattern.String()
		r.removeFromSlice(r.byKey[keyStr], action)
	}

	return nil
}

// Get retrieves an action by ID
func (r *Registry) Get(id ActionID) (*Action, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	action, exists := r.actions[id]
	return action, exists
}

// ActionsForContext returns all actions available in a context
func (r *Registry) ActionsForContext(ctx Context, actx *ActionContext) []*Action {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Get actions for this context + global actions
	contextActions := r.byContext[ctx]
	globalActions := r.byContext[ContextGlobal]

	// Combine and filter by predicate
	var available []*Action
	for _, action := range append(contextActions, globalActions...) {
		if action.Predicate == nil || action.Predicate(actx) {
			available = append(available, action)
		}
	}

	return available
}

// ActionsForKey returns actions bound to a specific key
func (r *Registry) ActionsForKey(key string, ctx Context, actx *ActionContext) []*Action {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	keyActions := r.byKey[key]

	// Filter by context and predicate
	var available []*Action
	for _, action := range keyActions {
		// Check context
		if !r.actionMatchesContext(action, ctx) {
			continue
		}

		// Check predicate
		if action.Predicate != nil && !action.Predicate(actx) {
			continue
		}

		available = append(available, action)
	}

	return available
}

// All returns all registered actions
func (r *Registry) All() []*Action {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	actions := make([]*Action, 0, len(r.actions))
	for _, action := range r.actions {
		actions = append(actions, action)
	}

	return actions
}

// Search performs fuzzy search on actions
func (r *Registry) Search(query string, actx *ActionContext) []*Action {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	query = strings.ToLower(query)
	var results []*Action

	for _, action := range r.actions {
		// Check if action is available
		if action.Predicate != nil && !action.Predicate(actx) {
			continue
		}

		// Check if should show in UI
		if action.ShowInUI != nil && !action.ShowInUI(actx) {
			continue
		}

		// Fuzzy match on name, description, tags
		name := strings.ToLower(action.Name)
		desc := strings.ToLower(action.Description)

		if strings.Contains(name, query) || strings.Contains(desc, query) {
			results = append(results, action)
			continue
		}

		// Check tags
		for _, tag := range action.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, action)
				break
			}
		}
	}

	// Sort by relevance
	sort.Slice(results, func(i, j int) bool {
		iNameMatch := strings.Contains(strings.ToLower(results[i].Name), query)
		jNameMatch := strings.Contains(strings.ToLower(results[j].Name), query)

		if iNameMatch != jNameMatch {
			return iNameMatch
		}

		return results[i].Name < results[j].Name
	})

	return results
}

// Helper: Check if action applies to context
func (r *Registry) actionMatchesContext(action *Action, ctx Context) bool {
	for _, actCtx := range action.Contexts {
		if actCtx == ContextGlobal || actCtx == ctx {
			return true
		}
	}
	return false
}

// Helper: Remove action from slice
func (r *Registry) removeFromSlice(slice []*Action, action *Action) []*Action {
	for i, a := range slice {
		if a.ID == action.ID {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
