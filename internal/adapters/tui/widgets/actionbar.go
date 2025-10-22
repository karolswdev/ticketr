package widgets

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ActionBarContext represents the current context/view state.
type ActionBarContext string

const (
	ContextWorkspaceList ActionBarContext = "workspace_list"
	ContextTicketTree    ActionBarContext = "ticket_tree"
	ContextTicketDetail  ActionBarContext = "ticket_detail"
	ContextModal         ActionBarContext = "modal"
	ContextSyncing       ActionBarContext = "syncing"
)

// KeyBinding represents a single keybinding to display.
type KeyBinding struct {
	Key         string
	Description string
}

// ActionBar is a bottom status bar widget that displays context-aware keybindings.
type ActionBar struct {
	*tview.TextView
	context  ActionBarContext
	bindings map[ActionBarContext][]KeyBinding

	// Marquee for overflow text
	marquee   *Marquee
	app       *tview.Application
	lastWidth int // FIX #3: Track terminal width for resize detection
}

// NewActionBar creates a new action bar widget.
func NewActionBar() *ActionBar {
	ab := &ActionBar{
		TextView: tview.NewTextView(),
		context:  ContextWorkspaceList,
	}

	ab.SetDynamicColors(true)
	ab.SetBorder(true)
	ab.SetBorderColor(tcell.ColorDimGray)
	ab.SetTitle(" Keybindings ")

	// Initialize default keybindings for each context
	ab.bindings = map[ActionBarContext][]KeyBinding{
		ContextWorkspaceList: {
			{Key: "Enter", Description: "Select Workspace"},
			{Key: "Tab", Description: "Next Panel"},
			{Key: "n", Description: "New Workspace"},
			{Key: "Esc/W/F3", Description: "Close Panel"},
			{Key: "?", Description: "Help"},
		},
		ContextTicketTree: {
			{Key: "Enter", Description: "Open Ticket"},
			{Key: "Space", Description: "Select/Deselect"},
			{Key: "W/F3", Description: "Workspaces"},
			{Key: "Tab", Description: "Next Panel"},
			{Key: "j/k", Description: "Navigate"},
			{Key: "h/l", Description: "Collapse/Expand"},
			{Key: "b", Description: "Bulk Ops"},
			{Key: "/", Description: "Search"},
			{Key: ":", Description: "Commands"},
			{Key: "?", Description: "Help"},
		},
		ContextTicketDetail: {
			{Key: "Esc", Description: "Back"},
			{Key: "Tab", Description: "Next Panel"},
			{Key: "W/F3", Description: "Workspaces"},
			{Key: "e", Description: "Edit"},
			{Key: ":", Description: "Commands"},
			{Key: "?", Description: "Help"},
		},
		ContextModal: {
			{Key: "Esc", Description: "Close"},
			{Key: "Enter", Description: "Confirm"},
		},
		ContextSyncing: {
			{Key: "Esc", Description: "Cancel Operation"},
			{Key: "Ctrl+C", Description: "Quit"},
		},
	}

	ab.update()
	return ab
}

// SetContext updates the action bar to show keybindings for the given context.
func (ab *ActionBar) SetContext(ctx ActionBarContext) {
	if ab.context != ctx {
		ab.context = ctx
		ab.update()
	}
}

// GetContext returns the current context.
func (ab *ActionBar) GetContext() ActionBarContext {
	return ab.context
}

// AddBinding adds a custom keybinding for a specific context.
func (ab *ActionBar) AddBinding(ctx ActionBarContext, key, description string) {
	if ab.bindings[ctx] == nil {
		ab.bindings[ctx] = []KeyBinding{}
	}
	ab.bindings[ctx] = append(ab.bindings[ctx], KeyBinding{
		Key:         key,
		Description: description,
	})
}

// SetBindings replaces all keybindings for a specific context.
func (ab *ActionBar) SetBindings(ctx ActionBarContext, bindings []KeyBinding) {
	ab.bindings[ctx] = bindings
	if ab.context == ctx {
		ab.update()
	}
}

// update refreshes the displayed keybindings based on current context.
// Implements responsive truncation to fit within terminal width.
// Uses marquee for overflow text scrolling.
func (ab *ActionBar) update() {
	bindings := ab.bindings[ab.context]
	if len(bindings) == 0 {
		ab.SetText("")
		// Stop marquee if running
		if ab.marquee != nil {
			ab.marquee.Stop()
		}
		return
	}

	// Get available width (subtract border and padding)
	// Note: GetInnerRect() may return a small default before layout, so use a reasonable minimum
	_, _, width, _ := ab.GetInnerRect()
	if width < 80 {
		// Use a standard 80-column terminal as default for initial render
		// This ensures tests and initial display work correctly before actual layout
		width = 80
	}

	// Reserve space for borders, padding
	const borderOverhead = 4 // Border chars + padding
	const minWidthPerBinding = 10 // Minimum space needed for one binding

	availableWidth := width - borderOverhead

	// If terminal is too narrow, show minimal message
	if availableWidth < minWidthPerBinding {
		ab.SetText("[yellow]Press ? for help")
		// Stop marquee if running
		if ab.marquee != nil {
			ab.marquee.Stop()
		}
		return
	}

	// Build full bindings string (no truncation)
	text := ab.buildFullBindings(bindings)

	// Check if text overflows
	visualLen := ab.visualLength(text)
	if visualLen > availableWidth {
		// Text overflows - use marquee
		ab.enableMarquee(text, availableWidth)
	} else {
		// Text fits - disable marquee and show directly
		ab.disableMarquee()
		ab.SetText(text)
	}
}

// enableMarquee starts marquee scrolling for overflow text.
func (ab *ActionBar) enableMarquee(text string, width int) {
	fmt.Printf("[DEBUG] enableMarquee called: width=%d, textLen=%d\n", width, len(text))
	fmt.Printf("[DEBUG] text sample: %s...\n", text[:min(100, len(text))])

	if ab.marquee == nil {
		fmt.Println("[DEBUG] Creating new marquee")
		// EMERGENCY FIX: Use legacy NewMarquee() constructor which creates single-item mode
		// This bypasses the broken parsing logic and uses simple horizontal scroll
		ab.marquee = NewMarquee(text, width)
		fmt.Println("[DEBUG] Calling marquee.Start()")
		ab.marquee.Start()

		// FIX #3: Start resize monitoring
		fmt.Println("[DEBUG] Starting resize monitoring")
		go ab.monitorTerminalSize()
		// Start update loop for marquee
		fmt.Println("[DEBUG] Starting marquee update loop")
		go ab.marqueeUpdateLoop()
	} else {
		fmt.Println("[DEBUG] Updating existing marquee")
		// Update existing marquee
		ab.marquee.SetText(text)
		ab.marquee.CheckResize(width) // FIX #3: Use CheckResize instead of SetWidth
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// disableMarquee stops marquee scrolling.
func (ab *ActionBar) disableMarquee() {
	if ab.marquee != nil {
		ab.marquee.Stop()
	}
}

// marqueeUpdateLoop periodically updates the action bar with animated text.
// Updates at ~30 FPS to match marquee animation frame rate.
func (ab *ActionBar) marqueeUpdateLoop() {
	ticker := time.NewTicker(33 * time.Millisecond) // ~30 FPS for smooth animations
	defer ticker.Stop()

	for range ticker.C {
		if ab.marquee == nil {
			return // Only exit if marquee is completely gone
		}

		// Get current display text from marquee (with animation)
		displayText := ab.marquee.GetDisplayText()

		// Update action bar if app is available
		if ab.app != nil {
			ab.app.QueueUpdateDraw(func() {
				ab.SetText(displayText)
			})
		} else {
			// No app - update directly
			ab.SetText(displayText)
		}
	}
}

// monitorTerminalSize monitors terminal size changes and updates marquee accordingly.
// FIX #3: Add terminal resize monitoring to stop scrolling when text fits.
func (ab *ActionBar) monitorTerminalSize() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if ab.marquee == nil {
			return
		}

		// Check if inner width has changed (indicates terminal resize)
		_, _, innerWidth, _ := ab.GetInnerRect()
		if innerWidth > 0 && innerWidth != ab.lastWidth {
			ab.lastWidth = innerWidth

			// Calculate available width for marquee
			if innerWidth < 80 {
				innerWidth = 80
			}
			availableWidth := innerWidth - 4 // Border overhead

			// Update marquee with new width using app queue for thread safety
			if ab.app != nil {
				ab.app.QueueUpdateDraw(func() {
					if ab.marquee != nil {
						ab.marquee.CheckResize(availableWidth)
					}
				})
			}
		}
	}
}

// buildFullBindings constructs the complete keybinding string without truncation.
func (ab *ActionBar) buildFullBindings(bindings []KeyBinding) string {
	var builder strings.Builder

	for i, binding := range bindings {
		// Add separator if not first
		if i > 0 {
			builder.WriteString(" ")
		}

		// Format: [Key Action]
		bindingText := fmt.Sprintf("[yellow][%s[white] %s[yellow]]", binding.Key, binding.Description)
		builder.WriteString(bindingText)
	}

	return builder.String()
}

// buildResponsiveBindings constructs the keybinding string, truncating if necessary.
func (ab *ActionBar) buildResponsiveBindings(bindings []KeyBinding, maxWidth int, truncationMsg string) string {
	var builder strings.Builder
	truncationNeeded := false

	for i, binding := range bindings {
		// Format: [Key Action]
		bindingText := fmt.Sprintf("[yellow][%s[white] %s[yellow]]", binding.Key, binding.Description)

		// Calculate visual length (strip color codes for measurement)
		visualLen := ab.visualLength(bindingText)

		// Add separator if not first
		if i > 0 {
			visualLen += 1 // Space separator
		}

		// Check if adding this binding would overflow
		currentLen := ab.visualLength(builder.String())
		if currentLen+visualLen+ab.visualLength(truncationMsg) > maxWidth {
			// Would overflow - truncate here
			truncationNeeded = true
			break
		}

		// Add separator
		if i > 0 {
			builder.WriteString(" ")
		}

		// Add binding
		builder.WriteString(bindingText)
	}

	// Handle truncation
	if truncationNeeded {
		if builder.Len() > 0 {
			// We fit at least one binding - add truncation message
			builder.WriteString(truncationMsg)
		} else {
			// Terminal too narrow - couldn't fit any bindings
			// Just show help hint
			return "[yellow]Press ? for help"
		}
	}

	return builder.String()
}

// visualLength calculates the visual length of a string, stripping tview color tags.
// This is a simplified version that handles basic tview color codes.
func (ab *ActionBar) visualLength(s string) int {
	// Simple state machine to skip color tags: [color]text[color]
	length := 0
	inTag := false

	for _, r := range s {
		switch {
		case r == '[':
			inTag = true
		case r == ']' && inTag:
			inTag = false
		case !inTag:
			length++
		}
	}

	return length
}

// Primitive returns the underlying tview primitive.
func (ab *ActionBar) Primitive() tview.Primitive {
	return ab.TextView
}

// Refresh forces a refresh of the action bar content.
// Useful when terminal is resized or context changes.
func (ab *ActionBar) Refresh() {
	ab.update()
}

// SetApp sets the tview application for marquee updates.
func (ab *ActionBar) SetApp(app *tview.Application) {
	ab.app = app
}

// Shutdown stops any active marquee and cleans up resources.
func (ab *ActionBar) Shutdown() {
	if ab.marquee != nil {
		ab.marquee.Shutdown()
		ab.marquee = nil
	}
}
