package views

import "github.com/rivo/tview"

// View represents a screen/panel in the TUI application.
// All views must implement this interface to be managed by the router.
type View interface {
	// Name returns the unique identifier for this view
	Name() string

	// Primitive returns the tview primitive to be displayed
	Primitive() tview.Primitive

	// OnShow is called when the view becomes active
	OnShow()

	// OnHide is called when the view is hidden
	OnHide()
}
