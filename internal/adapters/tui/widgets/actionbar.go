package widgets

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ActionBarContext represents the current context/view state.
type ActionBarContext string

const (
	ContextWorkspaceList ActionBarContext = "workspace_list"
	ContextTicketTree    ActionBarContext = "ticket_tree"
	ContextTicketDetail  ActionBarContext = "ticket_detail"
	ContextModal         ActionBarContext = "modal"
	ContextSyncing       ActionBarContext = "syncing"
)

// KeyBinding represents a single keybinding to display.
type KeyBinding struct {
	Key         string
	Description string
}

// ActionBar is a bottom status bar widget that displays context-aware keybindings.
type ActionBar struct {
	*tview.TextView
	context  ActionBarContext
	bindings map[ActionBarContext][]KeyBinding
}

// NewActionBar creates a new action bar widget.
func NewActionBar() *ActionBar {
	ab := &ActionBar{
		TextView: tview.NewTextView(),
		context:  ContextWorkspaceList,
	}

	ab.SetDynamicColors(true)
	ab.SetBorder(true)
	ab.SetBorderColor(tcell.ColorDimGray)
	ab.SetTitle(" Keybindings ")

	// Initialize default keybindings for each context
	ab.bindings = map[ActionBarContext][]KeyBinding{
		ContextWorkspaceList: {
			{Key: "Enter", Description: "Select Workspace"},
			{Key: "Tab", Description: "Next Panel"},
			{Key: "n", Description: "New Workspace"},
			{Key: "?", Description: "Help"},
			{Key: "q/Ctrl+C", Description: "Quit"},
		},
		ContextTicketTree: {
			{Key: "Enter", Description: "Open Ticket"},
			{Key: "Space", Description: "Select/Deselect"},
			{Key: "Tab", Description: "Next Panel"},
			{Key: "Esc", Description: "Back"},
			{Key: "j/k", Description: "Navigate"},
			{Key: "h/l", Description: "Collapse/Expand"},
			{Key: "b", Description: "Bulk Ops"},
			{Key: "/", Description: "Search"},
			{Key: ":", Description: "Commands"},
			{Key: "?", Description: "Help"},
		},
		ContextTicketDetail: {
			{Key: "Esc", Description: "Back"},
			{Key: "Tab", Description: "Next Panel"},
			{Key: "e", Description: "Edit"},
			{Key: ":", Description: "Commands"},
			{Key: "?", Description: "Help"},
		},
		ContextModal: {
			{Key: "Esc", Description: "Close"},
			{Key: "Enter", Description: "Confirm"},
		},
		ContextSyncing: {
			{Key: "Esc", Description: "Cancel Operation"},
			{Key: "Ctrl+C", Description: "Quit"},
		},
	}

	ab.update()
	return ab
}

// SetContext updates the action bar to show keybindings for the given context.
func (ab *ActionBar) SetContext(ctx ActionBarContext) {
	if ab.context != ctx {
		ab.context = ctx
		ab.update()
	}
}

// GetContext returns the current context.
func (ab *ActionBar) GetContext() ActionBarContext {
	return ab.context
}

// AddBinding adds a custom keybinding for a specific context.
func (ab *ActionBar) AddBinding(ctx ActionBarContext, key, description string) {
	if ab.bindings[ctx] == nil {
		ab.bindings[ctx] = []KeyBinding{}
	}
	ab.bindings[ctx] = append(ab.bindings[ctx], KeyBinding{
		Key:         key,
		Description: description,
	})
}

// SetBindings replaces all keybindings for a specific context.
func (ab *ActionBar) SetBindings(ctx ActionBarContext, bindings []KeyBinding) {
	ab.bindings[ctx] = bindings
	if ab.context == ctx {
		ab.update()
	}
}

// update refreshes the displayed keybindings based on current context.
func (ab *ActionBar) update() {
	bindings := ab.bindings[ab.context]
	if len(bindings) == 0 {
		ab.SetText("")
		return
	}

	// Format bindings as: [Key Action] [Key Action] ...
	var text string
	for i, binding := range bindings {
		if i > 0 {
			text += " "
		}
		text += fmt.Sprintf("[yellow][%s[white] %s[yellow]]", binding.Key, binding.Description)
	}

	ab.SetText(text)
}

// Primitive returns the underlying tview primitive.
func (ab *ActionBar) Primitive() tview.Primitive {
	return ab.TextView
}
