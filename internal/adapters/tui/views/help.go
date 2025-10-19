package views

import (
	"github.com/gdamore/tcell/v2"
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
	textView.SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)

	view := &HelpView{
		textView: textView,
	}

	// Set help content
	view.setContent()

	// Setup keybindings
	view.setupKeybindings()

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
	v.textView.ScrollToBeginning()
}

// OnHide is called when the view is hidden.
func (v *HelpView) OnHide() {
	// No cleanup needed
}

// setContent populates the help text with keybindings.
func (v *HelpView) setContent() {
	content := `[yellow::b]Ticketr TUI - Keyboard Shortcuts[-:-:-]

[cyan::b]Global Navigation[-:-:-]
  [green]Tab[-]        Cycle focus forward (workspace → tree → detail → workspace)
  [green]Shift+Tab[-]  Cycle focus backward
  [green]Esc[-]        Go back one panel (detail → tree → workspace)
  [green]?[-]          Show this help screen
  [green]q[-]          Quit application

[cyan::b]Workspace List Panel[-:-:-]
  [green]j/k[-]        Move down/up in list
  [green]↓/↑[-]        Move down/up in list
  [green]Enter[-]      Select workspace and load tickets

[cyan::b]Ticket Tree Panel[-:-:-]
  [green]j/k[-]        Move down/up in tree
  [green]↓/↑[-]        Move down/up in tree
  [green]h/l[-]        Collapse/expand node
  [green]←/→[-]        Collapse/expand node
  [green]Enter[-]      Open ticket detail view

[cyan::b]Ticket Detail Panel (Read-Only Mode)[-:-:-]
  [green]e[-]          Enter edit mode
  [green]j/k[-]        Scroll down/up
  [green]g[-]          Go to top
  [green]G[-]          Go to bottom (Shift+g)
  [green]Esc[-]        Return to ticket tree

[cyan::b]Ticket Detail Panel (Edit Mode)[-:-:-]
  [green]Tab[-]        Move between form fields
  [green]Save[-]       Save changes (click button or navigate to it)
  [green]Cancel[-]     Cancel editing (click button or navigate to it)
  [green]Esc[-]        Cancel editing and discard changes

[cyan::b]Field Validation (Edit Mode)[-:-:-]
  • Title: Required field
  • Jira ID: Must match format PROJECT-123 (uppercase, dash, numbers)
  • Custom Fields: Format as key=value (one per line)
  • Acceptance Criteria: One criterion per line

[cyan::b]Visual Indicators[-:-:-]
  [green]Green border[-]  Currently focused panel
  [white]White border[-]  Inactive panel
  [green]●[-]           Ticket synced with Jira
  [white]○[-]           Ticket not synced
  [green]■[-]           Task synced with Jira
  [cyan]□[-]           Task not synced
  [red]*[-]            Unsaved changes in detail view

[cyan::b]Tips[-:-:-]
  • Use Tab to quickly navigate between panels
  • Press Enter on a ticket in the tree to view details
  • Edit mode validates fields on save attempt
  • Esc always goes back or cancels current operation
  • Run 'ticketr pull' to sync tickets from Jira
  • Vim-style keys (j/k/h/l) work alongside arrow keys

[cyan::b]About[-:-:-]
Ticketr v3 - Jira-Markdown synchronization tool
Phase 4 Week 13: Ticket detail editor with validation
Architecture: Hexagonal (Ports & Adapters)

Press [green]Esc[-] or [green]?[-] to close this help screen.
`
	v.textView.SetText(content)
}

func (v *HelpView) setupKeybindings() {
	v.textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			// Return to previous view (handled by global handler)
			return event
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				// Let global handler catch this
				return event
			case 'j':
				// Scroll down
				row, col := v.textView.GetScrollOffset()
				v.textView.ScrollTo(row+1, col)
				return nil
			case 'k':
				// Scroll up
				row, col := v.textView.GetScrollOffset()
				v.textView.ScrollTo(row-1, col)
				return nil
			case 'g':
				v.textView.ScrollToBeginning()
				return nil
			case 'G':
				v.textView.ScrollToEnd()
				return nil
			}
		}
		return event
	})
}
