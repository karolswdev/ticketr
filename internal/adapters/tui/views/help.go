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
  [::b]Tab[::-]         Switch to next panel

[::b]Workspace List[::-]
  [::b]↑/↓[::-]         Navigate workspaces
  [::b]Enter[::-]       Switch to selected workspace
  [::b]w[::-]           Focus workspace panel (when implemented)

[::b]Ticket Tree[::-]
  [::b]↑/↓[::-]         Navigate tickets
  [::b]←/→[::-]         Collapse/expand nodes
  [::b]Enter[::-]       View ticket details
  [::b]t[::-]           Focus ticket panel (when implemented)

[::b]Ticket Detail[::-]
  [::b]↑/↓[::-]         Scroll content
  [::b]e[::-]           Edit ticket (when implemented)
  [::b]s[::-]           Sync with Jira (when implemented)

[::b]About[::-]
Ticketr is a Jira-Markdown synchronization tool with a Terminal User Interface.
This is a skeleton implementation - full features coming in future phases.

Version: 3.0.0-alpha
Architecture: Hexagonal (Ports & Adapters)

Press [::b]?[::-] again to close this help screen.
`
	v.textView.SetText(content)
}
