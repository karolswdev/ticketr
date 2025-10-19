package tui

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/adapters/tui/views"
	"github.com/rivo/tview"
)

// Router manages view navigation and lifecycle.
type Router struct {
	views       map[string]views.View
	currentView views.View
	pages       *tview.Pages
}

// NewRouter creates a new view router.
func NewRouter() *Router {
	return &Router{
		views: make(map[string]views.View),
		pages: tview.NewPages(),
	}
}

// Register adds a view to the router.
func (r *Router) Register(view views.View) error {
	name := view.Name()
	if _, exists := r.views[name]; exists {
		return fmt.Errorf("view %s already registered", name)
	}

	r.views[name] = view
	r.pages.AddPage(name, view.Primitive(), true, false)
	return nil
}

// Show displays the specified view.
func (r *Router) Show(name string) error {
	view, exists := r.views[name]
	if !exists {
		return fmt.Errorf("view %s not found", name)
	}

	// Hide current view
	if r.currentView != nil {
		r.currentView.OnHide()
	}

	// Show new view
	r.pages.SwitchToPage(name)
	r.currentView = view
	view.OnShow()

	return nil
}

// Current returns the currently active view.
func (r *Router) Current() views.View {
	return r.currentView
}

// ClearCurrent clears the current view state.
func (r *Router) ClearCurrent() {
	if r.currentView != nil {
		r.currentView.OnHide()
		r.currentView = nil
	}
}

// Pages returns the tview Pages primitive.
func (r *Router) Pages() *tview.Pages {
	return r.pages
}

// ToggleView switches between two views (useful for help toggle).
func (r *Router) ToggleView(name1, name2 string) error {
	if r.currentView == nil {
		return r.Show(name1)
	}

	currentName := r.currentView.Name()
	if currentName == name1 {
		return r.Show(name2)
	}
	return r.Show(name1)
}
