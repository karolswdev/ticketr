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
  [green]/[-]          Open search (fuzzy search with filters)
  [green]:[-]          Open command palette
  [green]?[-]          Show this help screen
  [green]q[-]          Quit application

[cyan::b]Sync Operations (Week 15)[-:-:-]
  [green]p[-]          Push tickets to Jira (async, non-blocking)
  [green]P[-]          Pull tickets from Jira (async, non-blocking)
  [green]r[-]          Refresh current workspace tickets
  [green]s[-]          Full sync (pull then push, async)

[cyan::b]Page Navigation (Week 16 - NEW!)[-:-:-]
  [green]Ctrl+F[-]     Page down (full page scroll)
  [green]Ctrl+B[-]     Page up (full page scroll)
  [green]Ctrl+D[-]     Half-page down
  [green]Ctrl+U[-]     Half-page up

  Available in: Ticket Detail, Help View, Search Results

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
  [green]j/k[-]        Scroll down/up (one line)
  [green]Ctrl+F/B[-]   Page down/up (full page)
  [green]Ctrl+D/U[-]   Half-page down/up
  [green]g[-]          Go to top
  [green]G[-]          Go to bottom (Shift+g)
  [green]Esc[-]        Return to ticket tree

[cyan::b]Ticket Detail Panel (Edit Mode)[-:-:-]
  [green]Tab[-]        Move between form fields
  [green]Save[-]       Save changes (click button or navigate to it)
  [green]Cancel[-]     Cancel editing (click button or navigate to it)
  [green]Esc[-]        Cancel editing and discard changes

[cyan::b]Search View (/)[-:-:-]
  [green]Type[-]       Filter tickets by text
  [green]@user[-]      Filter by assignee (@john)
  [green]#ID[-]        Filter by Jira ID (#BACK-123 or #BACK)
  [green]!priority[-]  Filter by priority (!high, !p1)
  [green]~sprint[-]    Filter by sprint (~Sprint-42)
  [green]/regex/[-]    Filter by regex pattern (/bug.*fix/)
  [green]↓/j[-]        Move to results list
  [green]↑/k[-]        Move back to search input
  [green]Ctrl+F/B[-]   Navigate results (10 items at a time)
  [green]Ctrl+D/U[-]   Navigate results (5 items at a time)
  [green]Enter[-]      Open selected ticket
  [green]Esc[-]        Close search and return to main view

[cyan::b]Command Palette (:)[-:-:-]
  [green]Type[-]       Filter available commands
  [green]↓/↑[-]        Navigate command list
  [green]Enter[-]      Execute selected command
  [green]Esc[-]        Close command palette

  Available Commands:
    • push   - Push tickets to Jira
    • pull   - Pull tickets from Jira
    • sync   - Full sync (pull then push)
    • refresh - Refresh ticket list
    • help   - Show this help
    • quit   - Quit application

[cyan::b]Field Validation (Edit Mode)[-:-:-]
  • Title: Required field
  • Jira ID: Must match format PROJECT-123 (uppercase, dash, numbers)
  • Custom Fields: Format as key=value (one per line)
  • Acceptance Criteria: One criterion per line

[cyan::b]Visual Indicators[-:-:-]
  [green]Green border[-]  Currently focused panel (or themed primary color)
  [white]White border[-]  Inactive panel (or themed secondary color)
  [green]●[-]           Ticket synced with Jira
  [white]○[-]           Ticket not synced
  [green]■[-]           Task synced with Jira
  [cyan]□[-]           Task not synced
  [red]*[-]            Unsaved changes in detail view

[cyan::b]Sync Status Indicators[-:-:-]
  [white]○[-]           Idle - no sync operation in progress
  [yellow]◌[-]           Syncing - operation in progress (non-blocking)
  [green]●[-]           Success - last operation completed successfully
  [red]✗[-]           Error - last operation failed

[cyan::b]Responsive Layout (Week 16 - NEW!)[-:-:-]
  • Terminal >= 100 columns: Full tri-panel layout (workspace | tree | detail)
  • Terminal 60-99 columns: Compact layout (tree | detail)
    - Workspace info shown in status bar
  • Terminal < 60 columns: Error message, minimum width required

[cyan::b]Theme System (Week 16 - NEW!)[-:-:-]
  Three built-in themes available:
  • Default: Green/white (classic)
  • Dark: Blue/gray accents
  • Light: Purple/slate accents

  Themes affect border colors, status messages, and visual indicators.

[cyan::b]Tips[-:-:-]
  • Use Tab to quickly navigate between panels
  • Press Enter on a ticket in the tree to view details
  • Use / to quickly search and filter tickets with advanced queries
  • Combine filters: "@john !high ~Sprint-42 auth bug" finds all high-priority
    tickets assigned to John in Sprint-42 with "auth bug" in the text
  • Use : to access commands (push, pull, sync, refresh, help, quit)
  • Edit mode validates fields on save attempt
  • Esc always goes back or cancels current operation
  • Sync operations (p/P/s) run asynchronously - UI remains responsive
  • Watch the sync status bar for real-time operation progress
  • Vim-style keys (j/k/h/l) work alongside arrow keys
  • Use Ctrl+F/B for fast page navigation in long content
  • Resize terminal to see responsive layout in action

[cyan::b]Performance (Week 16)[-:-:-]
  • Optimized for 1000+ tickets
  • Smooth scrolling and navigation
  • Non-blocking async operations
  • Efficient tree rendering

[cyan::b]About[-:-:-]
Ticketr v3 - Jira-Markdown synchronization tool
Phase 4 Week 16: Final polish - keybindings, themes, and responsive layout
Architecture: Hexagonal (Ports & Adapters)

Press [green]Esc[-] or [green]?[-] to close this help screen.
Press [green]Ctrl+F/B[-] to page through this help.
`
	v.textView.SetText(content)
}

func (v *HelpView) setupKeybindings() {
	v.textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			// Return to previous view (handled by global handler)
			return event
		case tcell.KeyCtrlF:
			// Page down (full page)
			_, _, _, height := v.textView.GetInnerRect()
			row, col := v.textView.GetScrollOffset()
			v.textView.ScrollTo(row+height, col)
			return nil
		case tcell.KeyCtrlB:
			// Page up (full page)
			_, _, _, height := v.textView.GetInnerRect()
			row, col := v.textView.GetScrollOffset()
			newRow := row - height
			if newRow < 0 {
				newRow = 0
			}
			v.textView.ScrollTo(newRow, col)
			return nil
		case tcell.KeyCtrlD:
			// Half-page down
			_, _, _, height := v.textView.GetInnerRect()
			row, col := v.textView.GetScrollOffset()
			v.textView.ScrollTo(row+height/2, col)
			return nil
		case tcell.KeyCtrlU:
			// Half-page up
			_, _, _, height := v.textView.GetInnerRect()
			row, col := v.textView.GetScrollOffset()
			newRow := row - height/2
			if newRow < 0 {
				newRow = 0
			}
			v.textView.ScrollTo(newRow, col)
			return nil
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
