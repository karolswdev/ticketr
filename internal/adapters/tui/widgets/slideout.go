package widgets

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SlideOut is an overlay panel that slides in from the left side of the screen.
// It displays content on-demand and can be toggled on/off with keyboard shortcuts.
type SlideOut struct {
	*tview.Flex

	content          tview.Primitive
	width            int
	isVisible        bool
	onClose          func()
	background       *tview.Box // Semi-transparent background overlay
	cosmicBackground tview.Primitive // Optional cosmic background effect (Phase 6.5)

	// Layered container for cosmic background (Phase 6.5)
	pages *tview.Pages
}

// NewSlideOut creates a new slide-out overlay panel.
// The content primitive is displayed inside the slide-out.
// Width specifies how many columns the slide-out should occupy.
func NewSlideOut(content tview.Primitive, width int) *SlideOut {
	// Create background overlay (dimmed background)
	background := tview.NewBox()
	// Use a dim color to show it's an overlay without being too dark
	background.SetBackgroundColor(tcell.NewRGBColor(0, 0, 0))

	// Create main flex container
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create pages container for layering (Phase 6.5)
	// This allows cosmic background to be rendered behind the content
	pages := tview.NewPages()

	so := &SlideOut{
		Flex:       flex,
		content:    content,
		width:      width,
		isVisible:  false,
		background: background,
		pages:      pages,
	}

	// Set up input capture for ESC key
	flex.SetInputCapture(so.handleInput)

	return so
}

// Show displays the slide-out panel.
func (so *SlideOut) Show() {
	if so.isVisible {
		return
	}

	so.isVisible = true
	so.updateLayout()
}

// Hide conceals the slide-out panel.
func (so *SlideOut) Hide() {
	if !so.isVisible {
		return
	}

	so.isVisible = false
	so.updateLayout()

	// Call close callback if set
	if so.onClose != nil {
		so.onClose()
	}
}

// Toggle shows or hides the slide-out panel.
func (so *SlideOut) Toggle() {
	if so.isVisible {
		so.Hide()
	} else {
		so.Show()
	}
}

// IsVisible returns whether the slide-out is currently visible.
func (so *SlideOut) IsVisible() bool {
	return so.isVisible
}

// SetOnClose sets a callback to be invoked when the slide-out is closed.
func (so *SlideOut) SetOnClose(callback func()) {
	so.onClose = callback
}

// SetWidth updates the slide-out width.
func (so *SlideOut) SetWidth(width int) {
	so.width = width
	if so.isVisible {
		so.updateLayout()
	}
}

// updateLayout rebuilds the flex layout based on visibility.
func (so *SlideOut) updateLayout() {
	fmt.Fprintf(os.Stderr, "[DEBUG SlideOut] updateLayout called, isVisible: %v, cosmicBackground is nil: %v\n", so.isVisible, so.cosmicBackground == nil)

	// Clear existing items
	so.Flex.Clear()

	if so.isVisible {
		if so.cosmicBackground != nil {
			fmt.Fprintf(os.Stderr, "[DEBUG SlideOut] Using cosmic background layered approach\n")

			// PHASE 6.5 FIX: Use layered approach with Pages for cosmic background
			// Layer 1 (bottom): Cosmic background fills entire space
			// Layer 2 (top): Content panel on left side (transparent background)

			// Create a NEW pages container each time (Pages doesn't support clearing)
			so.pages = tview.NewPages()

			// Create a flex for the content panel (left side only)
			contentFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
			contentFlex.
				AddItem(so.content, so.width, 0, true).    // Content on left
				AddItem(tview.NewBox(), 0, 1, false)       // Transparent spacer on right

			// Add layers to pages (back to front)
			so.pages.
				AddPage("cosmic-bg", so.cosmicBackground, true, true).  // Layer 1: Cosmic background
				AddPage("content", contentFlex, true, true)             // Layer 2: Content panel

			// CRITICAL FIX: Add pages to flex AFTER setting it up
			so.Flex.AddItem(so.pages, 0, 1, true)
			fmt.Fprintf(os.Stderr, "[DEBUG SlideOut] Cosmic background layers added to pages and pages added to flex\n")
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG SlideOut] Using plain background (no cosmic effect)\n")
			// Original behavior: plain background on right
			so.Flex.
				AddItem(so.content, so.width, 0, true).      // Fixed width on left
				AddItem(so.background, 0, 1, false)          // Plain background on right
		}
	} else {
		fmt.Fprintf(os.Stderr, "[DEBUG SlideOut] Slide-out is hidden\n")
		// Hidden: show nothing (handled by tview.Pages layer)
	}
}

// SetBackgroundOverlay sets a cosmic background overlay for the slide-out (Phase 6.5).
// If set, this will be used instead of the plain background when the slide-out is visible.
func (so *SlideOut) SetBackgroundOverlay(overlay tview.Primitive) {
	fmt.Fprintf(os.Stderr, "[DEBUG SlideOut] SetBackgroundOverlay called, overlay is nil: %v\n", overlay == nil)
	so.cosmicBackground = overlay
	fmt.Fprintf(os.Stderr, "[DEBUG SlideOut] isVisible: %v\n", so.isVisible)
	if so.isVisible {
		fmt.Fprintf(os.Stderr, "[DEBUG SlideOut] Refreshing layout because slide-out is visible\n")
		so.updateLayout() // Refresh layout if already visible
	}
}

// handleInput processes keyboard input for the slide-out.
func (so *SlideOut) handleInput(event *tcell.EventKey) *tcell.EventKey {
	// ESC closes the slide-out
	if event.Key() == tcell.KeyEsc {
		so.Hide()
		return nil
	}

	return event
}

// Primitive returns the underlying tview primitive.
func (so *SlideOut) Primitive() tview.Primitive {
	return so.Flex
}
