package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/rivo/tview"
)

// WorkspaceListView displays available workspaces and allows switching between them.
type WorkspaceListView struct {
	list              *tview.List
	workspaceService  *services.WorkspaceService
	onSwitch          func(name string) error
	onWorkspaceChange func(workspaceID string)
}

// NewWorkspaceListView creates a new workspace list view.
func NewWorkspaceListView(workspaceService *services.WorkspaceService) *WorkspaceListView {
	list := tview.NewList()
	list.SetBorder(true).SetTitle(" Workspaces ")
	list.ShowSecondaryText(false)

	view := &WorkspaceListView{
		list:             list,
		workspaceService: workspaceService,
	}

	// Configure keybindings
	list.SetInputCapture(view.handleInput)

	return view
}

// Name returns the view identifier.
func (v *WorkspaceListView) Name() string {
	return "workspace_list"
}

// Primitive returns the tview primitive.
func (v *WorkspaceListView) Primitive() tview.Primitive {
	return v.list
}

// OnShow is called when the view becomes active.
func (v *WorkspaceListView) OnShow() {
	v.refresh()
}

// OnHide is called when the view is hidden.
func (v *WorkspaceListView) OnHide() {
	// No cleanup needed
}

// SetSwitchHandler sets the callback for workspace switching.
func (v *WorkspaceListView) SetSwitchHandler(handler func(name string) error) {
	v.onSwitch = handler
}

// SetWorkspaceChangeHandler sets the callback for workspace changes (receives workspace ID).
func (v *WorkspaceListView) SetWorkspaceChangeHandler(handler func(workspaceID string)) {
	v.onWorkspaceChange = handler
}

// refresh loads workspaces from the service and updates the list.
func (v *WorkspaceListView) refresh() {
	v.list.Clear()

	workspaces, err := v.workspaceService.List()
	if err != nil {
		v.list.AddItem(fmt.Sprintf("Error: %v", err), "", 0, nil)
		return
	}

	if len(workspaces) == 0 {
		v.list.AddItem("No workspaces found", "Create one with 'ticketr workspace create'", 0, nil)
		return
	}

	current, _ := v.workspaceService.Current()
	for _, ws := range workspaces {
		name := ws.Name
		if current != nil && ws.Name == current.Name {
			name = fmt.Sprintf("[::b]* %s[::-]", ws.Name)
		}

		v.list.AddItem(name, ws.JiraURL, 0, nil)
	}
}

// handleInput processes keyboard input for the workspace list.
func (v *WorkspaceListView) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		// Get selected workspace and trigger switch
		index := v.list.GetCurrentItem()
		if index >= 0 {
			workspaces, err := v.workspaceService.List()
			if err == nil && index < len(workspaces) {
				selectedWS := workspaces[index]
				if v.onSwitch != nil {
					_ = v.onSwitch(selectedWS.Name)
					v.refresh() // Refresh to show new current workspace
				}
				// Notify workspace change callback with ID
				if v.onWorkspaceChange != nil {
					v.onWorkspaceChange(selectedWS.ID)
				}
			}
		}
		return nil
	}

	return event
}

// SetFocused updates border color when focus changes.
func (v *WorkspaceListView) SetFocused(focused bool) {
	color := tcell.ColorWhite
	if focused {
		color = tcell.ColorGreen
	}
	v.list.SetBorderColor(color)
}
