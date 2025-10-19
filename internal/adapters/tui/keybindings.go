package tui

import (
	"github.com/gdamore/tcell/v2"
)

// KeyHandler processes global keyboard events.
type KeyHandler struct {
	app    *TUIApp
	router *Router
}

// NewKeyHandler creates a new key handler.
func NewKeyHandler(app *TUIApp, router *Router) *KeyHandler {
	return &KeyHandler{
		app:    app,
		router: router,
	}
}

// Handle processes a keyboard event and returns whether it was handled.
func (h *KeyHandler) Handle(event *tcell.EventKey) bool {
	// Global keybindings
	switch event.Key() {
	case tcell.KeyCtrlC:
		// Quit application
		h.app.Stop()
		return true

	case tcell.KeyTab:
		// Switch to next panel (future implementation)
		// For now, just pass through
		return false
	}

	// Character-based keybindings
	switch event.Rune() {
	case 'q':
		// Quit application
		h.app.Stop()
		return true

	case '?':
		// Toggle help view
		current := h.router.Current()
		if current != nil && current.Name() == "help" {
			// Return to previous view (for now, workspace_list)
			_ = h.router.Show("workspace_list")
		} else {
			_ = h.router.Show("help")
		}
		return true
	}

	return false
}
