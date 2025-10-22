package views

import (
	"fmt"
	"time"

	"github.com/karolswdev/ticktr/internal/adapters/tui/sync"
	"github.com/karolswdev/ticktr/internal/adapters/tui/widgets"
	"github.com/karolswdev/ticktr/internal/tui/jobs"
	"github.com/rivo/tview"
)

// SyncStatusView displays the current sync operation status.
type SyncStatusView struct {
	textView          *tview.TextView
	status            sync.SyncStatus
	showWorkspaceInfo bool   // For compact mode
	workspaceName     string // Current workspace name
	ticketCount       int    // Number of tickets

	// Progress tracking (Phase 6, Day 10-11)
	progressBar     *widgets.ProgressBar
	currentProgress *jobs.JobProgress
	jobStartTime    time.Time
	showProgress    bool

	// Animation support (Phase 6.5, Day 1)
	app             *tview.Application
	animationTicker *time.Ticker
	animationStop   chan struct{}
}

// NewSyncStatusView creates a new sync status view.
// app parameter is optional - if provided, enables spinner animations
func NewSyncStatusView(app ...*tview.Application) *SyncStatusView {
	textView := tview.NewTextView()
	textView.SetBorder(true).SetTitle(" Sync Status ")
	textView.SetDynamicColors(true).SetTextAlign(tview.AlignLeft)

	view := &SyncStatusView{
		textView:    textView,
		status:      sync.NewIdleStatus(),
		progressBar: widgets.NewProgressBar(40), // Default 40 chars for progress bar
	}

	// Store app reference if provided (for animations)
	if len(app) > 0 && app[0] != nil {
		view.app = app[0]
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

	// Start animation if syncing and app is available
	if status.State == sync.StateSyncing && v.app != nil {
		v.startAnimation()
	} else {
		v.stopAnimation()
	}

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

// UpdateProgress updates the progress bar with new job progress data (Phase 6, Day 10-11).
func (v *SyncStatusView) UpdateProgress(progress jobs.JobProgress) {
	v.currentProgress = &progress
	v.showProgress = true

	// Start progress bar timer if not already started
	if v.jobStartTime.IsZero() {
		v.jobStartTime = time.Now()
		v.progressBar.Start()
	}

	v.updateDisplay()
}

// ClearProgress hides the progress bar and clears progress data.
func (v *SyncStatusView) ClearProgress() {
	v.currentProgress = nil
	v.showProgress = false
	v.jobStartTime = time.Time{} // Reset to zero
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

		// Show progress bar if we have progress data (Phase 6, Day 10-11)
		if v.showProgress && v.currentProgress != nil {
			progressData := widgets.FromJobProgress(*v.currentProgress, v.jobStartTime)
			progressStr := v.progressBar.RenderCompact(progressData)
			text = fmt.Sprintf("%s[yellow]%s:[-] [white]%s",
				workspacePrefix,
				v.status.Operation,
				progressStr)
		} else {
			text = fmt.Sprintf("%s[yellow]%s: [white]%s",
				workspacePrefix,
				v.status.Operation,
				v.status.Progress)
		}

	case sync.StateSuccess:
		titleColor = "green"
		text = fmt.Sprintf("%s[green]%s completed: [white]%s",
			workspacePrefix,
			v.status.Operation,
			v.status.Progress)
		// Clear progress when transitioning to success
		if v.showProgress {
			v.ClearProgress()
		}

	case sync.StateError:
		titleColor = "red"
		text = fmt.Sprintf("%s[red]%s failed: [white]%s",
			workspacePrefix,
			v.status.Operation,
			v.status.Error)
		// Clear progress when transitioning to error
		if v.showProgress {
			v.ClearProgress()
		}

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

// startAnimation starts the spinner animation ticker for smooth progress updates.
func (v *SyncStatusView) startAnimation() {
	// Don't start if already running
	if v.animationTicker != nil {
		return
	}

	// Create ticker at 80ms intervals (12.5 FPS for smooth spinner)
	v.animationTicker = time.NewTicker(80 * time.Millisecond)
	v.animationStop = make(chan struct{})

	// Start animation goroutine
	go func() {
		for {
			select {
			case <-v.animationStop:
				return
			case <-v.animationTicker.C:
				// Trigger redraw via QueueUpdateDraw (thread-safe)
				if v.app != nil {
					v.app.QueueUpdateDraw(func() {
						v.updateDisplay()
					})
				}
			}
		}
	}()
}

// stopAnimation stops the spinner animation ticker.
func (v *SyncStatusView) stopAnimation() {
	if v.animationTicker == nil {
		return
	}

	// Stop ticker
	v.animationTicker.Stop()
	v.animationTicker = nil

	// Signal goroutine to stop
	if v.animationStop != nil {
		close(v.animationStop)
		v.animationStop = nil
	}
}
