package views

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/adapters/tui/sync"
	"github.com/rivo/tview"
)

// SyncStatusView displays the current sync operation status.
type SyncStatusView struct {
	textView          *tview.TextView
	status            sync.SyncStatus
	showWorkspaceInfo bool   // For compact mode
	workspaceName     string // Current workspace name
	ticketCount       int    // Number of tickets
}

// NewSyncStatusView creates a new sync status view.
func NewSyncStatusView() *SyncStatusView {
	textView := tview.NewTextView()
	textView.SetBorder(true).SetTitle(" Sync Status ")
	textView.SetDynamicColors(true).SetTextAlign(tview.AlignLeft)

	view := &SyncStatusView{
		textView: textView,
		status:   sync.NewIdleStatus(),
	}

	view.updateDisplay()
	return view
}

// Primitive returns the underlying tview primitive.
func (v *SyncStatusView) Primitive() tview.Primitive {
	return v.textView
}

// SetStatus updates the displayed sync status.
func (v *SyncStatusView) SetStatus(status sync.SyncStatus) {
	v.status = status
	v.updateDisplay()
}

// GetStatus returns the current sync status.
func (v *SyncStatusView) GetStatus() sync.SyncStatus {
	return v.status
}

// SetShowWorkspaceInfo enables displaying workspace info in status bar (for compact mode).
func (v *SyncStatusView) SetShowWorkspaceInfo(show bool) {
	v.showWorkspaceInfo = show
	v.updateDisplay()
}

// SetWorkspaceInfo sets the workspace name and ticket count for display.
func (v *SyncStatusView) SetWorkspaceInfo(name string, count int) {
	v.workspaceName = name
	v.ticketCount = count
	v.updateDisplay()
}

// updateDisplay refreshes the view content based on current status.
func (v *SyncStatusView) updateDisplay() {
	// Format the status message with colors
	var text string
	var titleColor string

	// Build workspace info prefix if enabled
	var workspacePrefix string
	if v.showWorkspaceInfo && v.workspaceName != "" {
		workspacePrefix = fmt.Sprintf("[cyan]Workspace:[-] [white]%s[-] | [cyan]Tickets:[-] [white]%d[-] | ",
			v.workspaceName, v.ticketCount)
	}

	switch v.status.State {
	case sync.StateIdle:
		titleColor = "white"
		text = fmt.Sprintf("%s[white]%s", workspacePrefix, v.status.Progress)

	case sync.StateSyncing:
		titleColor = "yellow"
		text = fmt.Sprintf("%s[yellow]%s: [white]%s",
			workspacePrefix,
			v.status.Operation,
			v.status.Progress)

	case sync.StateSuccess:
		titleColor = "green"
		text = fmt.Sprintf("%s[green]%s completed: [white]%s",
			workspacePrefix,
			v.status.Operation,
			v.status.Progress)

	case sync.StateError:
		titleColor = "red"
		text = fmt.Sprintf("%s[red]%s failed: [white]%s",
			workspacePrefix,
			v.status.Operation,
			v.status.Error)

	default:
		titleColor = "white"
		text = workspacePrefix + "[white]Unknown status"
	}

	// Update border title with colored state indicator
	stateIndicator := v.getStateIndicator()
	v.textView.SetTitle(fmt.Sprintf(" Sync Status [%s]%s[-] ", titleColor, stateIndicator))

	// Set the text content
	v.textView.SetText(text)
}

// getStateIndicator returns a visual indicator for the current state.
func (v *SyncStatusView) getStateIndicator() string {
	switch v.status.State {
	case sync.StateIdle:
		return "○" // Empty circle
	case sync.StateSyncing:
		return "◌" // Dotted circle (in progress)
	case sync.StateSuccess:
		return "●" // Filled circle (success)
	case sync.StateError:
		return "✗" // X mark (error)
	default:
		return "?"
	}
}
