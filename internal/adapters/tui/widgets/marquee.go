package widgets

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// AnimationPhase represents the current animation state for a marquee item.
type AnimationPhase int

const (
	PhaseIdle     AnimationPhase = iota // Not animating
	PhaseSlideIn                        // Sliding in from right
	PhaseCenter                         // Holding in center (blinking)
	PhaseSlideOut                       // Sliding out to left
)

// MarqueeItem represents a single item in the marquee queue.
type MarqueeItem struct {
	Text      string // Text with color codes
	PlainText string // Text without color codes (for length calculation)
}

// Marquee provides theatrical per-item animations for overflow text.
// Each item slides in from right (italicized), holds center (blinking),
// then slides out to left (italicized). One item at a time.
type Marquee struct {
	// Item queue
	items       []MarqueeItem
	currentIdx  int
	singleItem  bool // True if only one item (special case: no animation)

	// Display state
	width       int

	// Animation state
	phase           AnimationPhase
	phaseStartTime  time.Time
	slideInDuration time.Duration // Default: 600ms
	centerDuration  time.Duration // Default: 1.5s
	slideOutDuration time.Duration // Default: 600ms
	blinkInterval   time.Duration // Default: 150ms

	// Position tracking (for slide animations)
	slideProgress float64 // 0.0 to 1.0
	blinkVisible  bool    // Toggle for blink effect

	// Lifecycle management
	ctx        context.Context
	cancel     context.CancelFunc
	ticker     *time.Ticker
	isAnimating bool
	mu         sync.RWMutex
}

// MarqueeConfig holds configuration for creating a marquee.
type MarqueeConfig struct {
	Items            []string      // Items to display (with color codes)
	Width            int           // Display width
	SlideInDuration  time.Duration // Time to slide in (default: 600ms)
	CenterDuration   time.Duration // Time to hold center (default: 1.5s)
	SlideOutDuration time.Duration // Time to slide out (default: 600ms)
	BlinkInterval    time.Duration // Blink interval during center (default: 150ms)
	FrameRate        time.Duration // Animation frame rate (default: 33ms = ~30 FPS)
}

// DefaultMarqueeConfig returns sensible defaults for marquee configuration.
func DefaultMarqueeConfig() MarqueeConfig {
	return MarqueeConfig{
		SlideInDuration:  600 * time.Millisecond,
		CenterDuration:   1500 * time.Millisecond,
		SlideOutDuration: 600 * time.Millisecond,
		BlinkInterval:    150 * time.Millisecond,
		FrameRate:        33 * time.Millisecond, // ~30 FPS
	}
}

// NewMarquee creates a new marquee with a single text item (legacy compatibility).
func NewMarquee(text string, width int) *Marquee {
	config := DefaultMarqueeConfig()
	config.Items = []string{text}
	config.Width = width
	return NewMarqueeWithConfig(config)
}

// NewMarqueeWithConfig creates a new marquee with custom configuration.
func NewMarqueeWithConfig(config MarqueeConfig) *Marquee {
	ctx, cancel := context.WithCancel(context.Background())

	// Parse items (split by pattern if needed, or use as-is)
	items := parseMarqueeItems(config.Items)

	return &Marquee{
		items:            items,
		currentIdx:       0,
		singleItem:       len(items) == 1,
		width:            config.Width,
		phase:            PhaseIdle,
		slideInDuration:  config.SlideInDuration,
		centerDuration:   config.CenterDuration,
		slideOutDuration: config.SlideOutDuration,
		blinkInterval:    config.BlinkInterval,
		ctx:              ctx,
		cancel:           cancel,
		blinkVisible:     true,
	}
}

// parseMarqueeItems converts raw text strings into MarqueeItem objects.
// Splits action bar format into individual bindings.
func parseMarqueeItems(items []string) []MarqueeItem {
	if len(items) == 0 {
		return []MarqueeItem{}
	}

	var result []MarqueeItem

	for _, item := range items {
		// Check if this looks like action bar format with multiple bindings
		// Format: [yellow][[white]Key[yellow] [white]Desc[yellow]] ...
		if strings.Contains(item, "[yellow][[white]") {
			// Split by binding pattern
			parts := strings.Split(item, "[yellow]]")
			for _, part := range parts {
				trimmed := strings.TrimSpace(part)
				if trimmed == "" {
					continue
				}
				// Reconstruct complete binding
				binding := trimmed
				if !strings.HasSuffix(binding, "[yellow]]") && trimmed != "" {
					binding = trimmed + "[yellow]]"
				}
				if !strings.HasPrefix(binding, "[yellow][[white]") && trimmed != "" {
					binding = "[yellow]" + binding
				}

				if binding != "" && visualLength(binding) > 0 {
					result = append(result, MarqueeItem{
						Text:      binding,
						PlainText: stripColorCodes(binding),
					})
				}
			}
		} else {
			// Single item, use as-is
			result = append(result, MarqueeItem{
				Text:      item,
				PlainText: stripColorCodes(item),
			})
		}
	}

	// Fallback: if parsing failed, use original
	if len(result) == 0 && len(items) > 0 {
		for _, item := range items {
			result = append(result, MarqueeItem{
				Text:      item,
				PlainText: stripColorCodes(item),
			})
		}
	}

	return result
}

// SetText updates the marquee text and resets animation state.
// For legacy compatibility - converts single text to items.
func (m *Marquee) SetText(text string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = parseMarqueeItems([]string{text})
	m.currentIdx = 0
	m.singleItem = len(m.items) == 1
	m.phase = PhaseIdle
	m.slideProgress = 0
	m.blinkVisible = true
}

// SetWidth updates the display width.
func (m *Marquee) SetWidth(width int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.width = width
}

// NeedsScrolling returns true if animation is needed.
// For multi-item marquees, always return true.
// For single items, return true if text exceeds width.
func (m *Marquee) NeedsScrolling() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.items) == 0 {
		return false
	}

	// Multiple items always need marquee
	if len(m.items) > 1 {
		return true
	}

	// Single item: check if it fits
	if len(m.items) == 1 {
		return visualLength(m.items[0].Text) > m.width
	}

	return false
}

// Start begins the animation.
func (m *Marquee) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Already animating
	if m.isAnimating {
		return
	}

	// No items
	if len(m.items) == 0 {
		return
	}

	// Single item that fits - no animation needed
	// For multi-item marquees, always animate (theatrical slide-in/blink/slide-out)
	// For single items that overflow, use horizontal scroll animation
	if m.singleItem && visualLength(m.items[0].Text) <= m.width {
		return
	}

	// Start animation
	m.isAnimating = true
	m.phase = PhaseSlideIn
	m.phaseStartTime = time.Now()
	m.slideProgress = 0
	m.blinkVisible = true
	m.ticker = time.NewTicker(33 * time.Millisecond) // ~30 FPS

	go m.animationLoop()
}

// Stop halts the animation and resets to idle state.
func (m *Marquee) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isAnimating {
		return
	}

	m.isAnimating = false
	if m.ticker != nil {
		m.ticker.Stop()
		m.ticker = nil
	}
	m.phase = PhaseIdle
	m.currentIdx = 0
	m.slideProgress = 0
	m.blinkVisible = true
}

// GetDisplayText returns the currently visible text based on animation state.
func (m *Marquee) GetDisplayText() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.items) == 0 {
		return ""
	}

	// Single item that fits - return as-is (no animation)
	if m.singleItem && visualLength(m.items[0].Text) <= m.width {
		return m.items[0].Text
	}

	// Single item that doesn't fit - use legacy horizontal scroll
	if m.singleItem {
		return m.extractVisibleWindow()
	}

	// Multi-item marquee - render current item with animation
	return m.renderCurrentItem()
}

// renderCurrentItem renders the current item based on animation phase.
func (m *Marquee) renderCurrentItem() string {
	if m.currentIdx >= len(m.items) {
		m.currentIdx = 0
	}

	item := m.items[m.currentIdx]
	itemLen := len(item.PlainText)

	switch m.phase {
	case PhaseIdle:
		// Show centered, no italics
		return m.centerText(item.Text, itemLen)

	case PhaseSlideIn:
		// Slide in from right, italicized
		italicText := italicize(item.Text)
		position := m.calculateSlideInPosition(itemLen)
		return m.positionText(italicText, itemLen, position)

	case PhaseCenter:
		// Center, blink, no italics
		if m.blinkVisible {
			return m.centerText(item.Text, itemLen)
		}
		// Blink off - show spaces to maintain layout
		return m.centerText(strings.Repeat(" ", itemLen), itemLen)

	case PhaseSlideOut:
		// Slide out to left, italicized
		italicText := italicize(item.Text)
		position := m.calculateSlideOutPosition(itemLen)
		return m.positionText(italicText, itemLen, position)
	}

	return item.Text
}

// calculateSlideInPosition calculates the X position during slide-in phase.
// Returns position where 0 = center, positive = right of center.
func (m *Marquee) calculateSlideInPosition(itemLen int) int {
	// Start position: right edge of screen
	startPos := m.width
	// End position: center
	endPos := (m.width - itemLen) / 2

	// Linear interpolation
	currentPos := startPos - int(float64(startPos-endPos)*m.slideProgress)
	return currentPos
}

// calculateSlideOutPosition calculates the X position during slide-out phase.
// Returns position where 0 = center, negative = left of center.
func (m *Marquee) calculateSlideOutPosition(itemLen int) int {
	// Start position: center
	startPos := (m.width - itemLen) / 2
	// End position: left edge (negative = off-screen)
	endPos := -itemLen

	// Linear interpolation
	currentPos := startPos + int(float64(endPos-startPos)*m.slideProgress)
	return currentPos
}

// centerText centers the text within the width.
func (m *Marquee) centerText(text string, visualLen int) string {
	if visualLen >= m.width {
		return text
	}

	padding := (m.width - visualLen) / 2
	return strings.Repeat(" ", padding) + text
}

// positionText positions the text at the given X coordinate.
func (m *Marquee) positionText(text string, visualLen int, xPos int) string {
	if xPos < 0 {
		// Off-screen left - show partial or nothing
		overflow := -xPos
		if overflow >= visualLen {
			return "" // Completely off-screen
		}
		// Show right portion
		return extractSubstring(text, overflow, visualLen-overflow)
	}

	if xPos >= m.width {
		// Off-screen right - not visible yet
		return ""
	}

	if xPos+visualLen > m.width {
		// Partially visible on right edge
		visibleLen := m.width - xPos
		leftPadding := strings.Repeat(" ", xPos)
		truncated := extractSubstring(text, 0, visibleLen)
		return leftPadding + truncated
	}

	// Fully visible - add left padding
	return strings.Repeat(" ", xPos) + text
}

// extractSubstring extracts a substring based on visual length, preserving color codes.
func extractSubstring(text string, start, length int) string {
	if length <= 0 {
		return ""
	}

	runes := []rune(text)
	visualPos := 0
	inTag := false
	startIdx := -1
	endIdx := len(runes)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '[' {
			inTag = true
		} else if r == ']' && inTag {
			inTag = false
			continue
		}

		if !inTag {
			if visualPos == start && startIdx == -1 {
				startIdx = i
				// Include any preceding tag
				for j := i - 1; j >= 0; j-- {
					if runes[j] == '[' {
						startIdx = j
						break
					}
					if runes[j] == ']' {
						break
					}
				}
			}

			if visualPos == start+length {
				endIdx = i
				// Include any following tag close
				for j := i; j < len(runes); j++ {
					if runes[j] == ']' {
						endIdx = j + 1
						break
					}
					if runes[j] == '[' {
						break
					}
				}
				break
			}

			visualPos++
		}
	}

	if startIdx == -1 || startIdx >= len(runes) {
		return ""
	}

	if endIdx > len(runes) {
		endIdx = len(runes)
	}

	return string(runes[startIdx:endIdx])
}

// italicize adds italic visual effect to text.
// Uses Unicode Mathematical Italic characters for supported letters.
func italicize(text string) string {
	// For now, use simple prefix/suffix markers
	// Unicode italic characters are complex and may not render well in all terminals
	// Alternative: Use slanted brackets or visual markers
	return "⟨" + text + "⟩" // Angle brackets suggest motion/skew
}

// animationLoop runs in a goroutine to update animation state.
func (m *Marquee) animationLoop() {
	for {
		// Safety check: ensure ticker exists
		m.mu.RLock()
		ticker := m.ticker
		m.mu.RUnlock()

		if ticker == nil {
			return
		}

		select {
		case <-m.ctx.Done():
			return

		case <-ticker.C:
			m.mu.Lock()

			if !m.isAnimating || m.ticker == nil {
				m.mu.Unlock()
				return
			}

			m.updateAnimationState()

			m.mu.Unlock()
		}
	}
}

// updateAnimationState updates the animation state machine.
// MUST be called with lock held.
func (m *Marquee) updateAnimationState() {
	now := time.Now()
	elapsed := now.Sub(m.phaseStartTime)

	switch m.phase {
	case PhaseSlideIn:
		// Update slide progress
		m.slideProgress = float64(elapsed) / float64(m.slideInDuration)

		if m.slideProgress >= 1.0 {
			// Transition to center phase
			m.phase = PhaseCenter
			m.phaseStartTime = now
			m.slideProgress = 0
			m.blinkVisible = true
		}

	case PhaseCenter:
		// Update blink state
		blinkCycle := int(elapsed / m.blinkInterval)
		m.blinkVisible = (blinkCycle % 2) == 0

		if elapsed >= m.centerDuration {
			// Transition to slide-out phase
			m.phase = PhaseSlideOut
			m.phaseStartTime = now
			m.slideProgress = 0
		}

	case PhaseSlideOut:
		// Update slide progress
		m.slideProgress = float64(elapsed) / float64(m.slideOutDuration)

		if m.slideProgress >= 1.0 {
			// Advance to next item
			m.currentIdx++
			if m.currentIdx >= len(m.items) {
				m.currentIdx = 0
			}

			// Reset to slide-in phase
			m.phase = PhaseSlideIn
			m.phaseStartTime = now
			m.slideProgress = 0
			m.blinkVisible = true
		}

	case PhaseIdle:
		// Start animation
		m.phase = PhaseSlideIn
		m.phaseStartTime = now
		m.slideProgress = 0
		m.blinkVisible = true
	}
}

// Shutdown stops the marquee and cleans up resources.
func (m *Marquee) Shutdown() {
	m.Stop()
	m.cancel()
}

// extractVisibleWindow returns the visible portion of text for single-item horizontal scroll.
// Legacy method for single long items (not used in new multi-item marquee).
func (m *Marquee) extractVisibleWindow() string {
	if len(m.items) == 0 {
		return ""
	}

	if m.width <= 0 {
		return ""
	}

	// Simple truncation for now (could add horizontal scroll if needed)
	text := m.items[0].Text
	visualLen := visualLength(text)

	if visualLen <= m.width {
		return text
	}

	// Truncate to width
	return extractSubstring(text, 0, m.width)
}

// visualLength calculates visible character count, ignoring color codes.
func visualLength(s string) int {
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

// stripColorCodes removes tview color tags from a string.
func stripColorCodes(s string) string {
	var result []rune
	inTag := false

	for _, r := range s {
		switch {
		case r == '[':
			inTag = true
		case r == ']' && inTag:
			inTag = false
		case !inTag:
			result = append(result, r)
		}
	}

	return string(result)
}

// GetRawText returns the full text with color codes (useful for debugging).
func (m *Marquee) GetRawText() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.items) == 0 {
		return ""
	}

	// Return first item for legacy compatibility
	return m.items[0].Text
}

// GetVisualLength returns the visual length of the current item.
func (m *Marquee) GetVisualLength() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.items) == 0 {
		return 0
	}

	return len(m.items[m.currentIdx].PlainText)
}

// CheckResize checks if width has changed and adjusts behavior.
func (m *Marquee) CheckResize(newWidth int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.width == newWidth {
		return
	}

	oldWidth := m.width
	m.width = newWidth

	// For multi-item marquees, width change doesn't affect animation need
	if len(m.items) > 1 {
		return
	}

	// For single items, check if animation need changed
	if len(m.items) == 1 {
		needsScroll := visualLength(m.items[0].Text) > newWidth

		if !needsScroll && m.isAnimating {
			// Text now fits - stop animation
			m.isAnimating = false
			if m.ticker != nil {
				m.ticker.Stop()
				m.ticker = nil
			}
			m.phase = PhaseIdle
		} else if needsScroll && !m.isAnimating {
			// Text now overflows - start animation
			m.isAnimating = true
			m.phase = PhaseSlideIn
			m.phaseStartTime = time.Now()
			m.slideProgress = 0
			m.blinkVisible = true
			m.ticker = time.NewTicker(33 * time.Millisecond)
			go m.animationLoop()
		}
	}

	// Log resize for debugging
	_ = fmt.Sprintf("Marquee resize: %d → %d", oldWidth, newWidth)
}
