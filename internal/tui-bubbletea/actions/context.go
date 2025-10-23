package actions

// Context represents where the user is in the application
type Context string

const (
	ContextWorkspaceList  Context = "workspace_list"
	ContextTicketTree     Context = "ticket_tree"
	ContextTicketDetail   Context = "ticket_detail"
	ContextSearch         Context = "search"
	ContextCommandPalette Context = "command_palette"
	ContextModal          Context = "modal"
	ContextSyncing        Context = "syncing"
	ContextHelp           Context = "help"
	ContextGlobal         Context = "*" // Matches any context
)

// ContextState tracks current application context
type ContextState struct {
	Current  Context
	Previous Context
	Stack    []Context
	Metadata map[string]interface{}
}

// IsIn checks if current context matches any of the given contexts
func (cs *ContextState) IsIn(contexts ...Context) bool {
	for _, ctx := range contexts {
		if ctx == ContextGlobal || ctx == cs.Current {
			return true
		}
	}
	return false
}

// ContextManager tracks and manages application context
type ContextManager struct {
	state    ContextState
	onChange []func(old, new Context)
}

// NewContextManager creates a new context manager
func NewContextManager(initial Context) *ContextManager {
	return &ContextManager{
		state: ContextState{
			Current:  initial,
			Previous: "",
			Stack:    []Context{initial},
			Metadata: make(map[string]interface{}),
		},
		onChange: []func(Context, Context){},
	}
}

// Switch changes the current context
func (cm *ContextManager) Switch(newContext Context) {
	old := cm.state.Current
	if old == newContext {
		return
	}

	cm.state.Previous = old
	cm.state.Current = newContext
	cm.state.Stack = append(cm.state.Stack, newContext)

	// Notify observers
	for _, fn := range cm.onChange {
		fn(old, newContext)
	}
}

// Push pushes a new context onto the stack (for modals)
func (cm *ContextManager) Push(ctx Context) {
	old := cm.state.Current
	cm.state.Previous = old
	cm.state.Current = ctx
	cm.state.Stack = append(cm.state.Stack, ctx)

	// Notify observers
	for _, fn := range cm.onChange {
		fn(old, ctx)
	}
}

// Pop returns to the previous context
func (cm *ContextManager) Pop() Context {
	if len(cm.state.Stack) <= 1 {
		return cm.state.Current
	}

	old := cm.state.Current
	cm.state.Stack = cm.state.Stack[:len(cm.state.Stack)-1]
	prev := cm.state.Stack[len(cm.state.Stack)-1]
	cm.state.Previous = old
	cm.state.Current = prev

	// Notify observers
	for _, fn := range cm.onChange {
		fn(old, prev)
	}

	return prev
}

// Current returns the current context
func (cm *ContextManager) Current() Context {
	return cm.state.Current
}

// OnChange registers a context change observer
func (cm *ContextManager) OnChange(fn func(old, new Context)) {
	cm.onChange = append(cm.onChange, fn)
}

// GetMetadata retrieves context metadata
func (cm *ContextManager) GetMetadata(key string) (interface{}, bool) {
	val, ok := cm.state.Metadata[key]
	return val, ok
}

// SetMetadata stores context metadata
func (cm *ContextManager) SetMetadata(key string, value interface{}) {
	cm.state.Metadata[key] = value
}
