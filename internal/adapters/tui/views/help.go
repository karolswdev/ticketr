package views

import (
	"github.com/rivo/tview"
)

// HelpView displays keyboard shortcuts and usage information.
type HelpView struct {
	textView *tview.TextView
}

// NewHelpView creates a new help view.
func NewHelpView() *HelpView {
	textView := tview.NewTextView()
	textView.SetBorder(true).SetTitle(" Help ")
	textView.SetDynamicColors(true)
	textView.SetWordWrap(true)

	view := &HelpView{
		textView: textView,
	}

	// Set help content
	view.setContent()

	return view
}

// Name returns the view identifier.
func (v *HelpView) Name() string {
	return "help"
}

// Primitive returns the tview primitive.
func (v *HelpView) Primitive() tview.Primitive {
	return v.textView
}

// OnShow is called when the view becomes active.
func (v *HelpView) OnShow() {
	// No initialization needed
}

// OnHide is called when the view is hidden.
func (v *HelpView) OnHide() {
	// No cleanup needed
}

// setContent populates the help text with keybindings.
func (v *HelpView) setContent() {
	content := `[yellow]Ticketr TUI - Keyboard Shortcuts[-]

[::b]Global Commands[::-]
  [::b]?[::-]           Toggle this help screen
  [::b]q[::-] / [::b]Ctrl+C[::-]  Quit application
  [::b]Tab[::-]         Switch focus between panels

[::b]Workspace List (Left Panel)[::-]
  [::b]j[::-] / [::b]↓[::-]       Move down
  [::b]k[::-] / [::b]↑[::-]       Move up
  [::b]Enter[::-]       Switch to selected workspace

[::b]Ticket Tree (Right Panel)[::-]
  [::b]j[::-] / [::b]↓[::-]       Move down
  [::b]k[::-] / [::b]↑[::-]       Move up
  [::b]h[::-] / [::b]←[::-]       Collapse node
  [::b]l[::-] / [::b]→[::-]       Expand node
  [::b]Enter[::-]       View ticket details (coming in Week 13)

[::b]Panel Focus Indicators[::-]
  [green]Green border[::-]   Currently focused panel
  [white]White border[::-]   Unfocused panel

[::b]Tips[::-]
  • Run 'ticketr pull' to sync tickets from Jira
  • Switch workspaces to view different project tickets
  • Tickets are loaded from your local database
  • Use vim-style keys (j/k/h/l) or arrow keys

[::b]About[::-]
Ticketr v3 - Jira-Markdown synchronization tool
Phase 4 Week 12: Multi-panel layout with real ticket data
Architecture: Hexagonal (Ports & Adapters)

Press [::b]?[::-] again to close this help screen.
`
	v.textView.SetText(content)
}
