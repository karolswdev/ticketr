package effects

import (
	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
	"github.com/rivo/tview"
)

// BorderChars contains the characters used for drawing borders.
type BorderChars struct {
	Horizontal     rune
	Vertical       rune
	TopLeft        rune
	TopRight       rune
	BottomLeft     rune
	BottomRight    rune
	VerticalLeft   rune
	VerticalRight  rune
	HorizontalDown rune
	HorizontalUp   rune
	Cross          rune
}

var (
	// SingleBorder uses single-line characters: ┌─┐
	SingleBorder = BorderChars{
		Horizontal:     '─',
		Vertical:       '│',
		TopLeft:        '┌',
		TopRight:       '┐',
		BottomLeft:     '└',
		BottomRight:    '┘',
		VerticalLeft:   '├',
		VerticalRight:  '┤',
		HorizontalDown: '┬',
		HorizontalUp:   '┴',
		Cross:          '┼',
	}

	// DoubleBorder uses double-line characters: ╔═╗
	DoubleBorder = BorderChars{
		Horizontal:     '═',
		Vertical:       '║',
		TopLeft:        '╔',
		TopRight:       '╗',
		BottomLeft:     '╚',
		BottomRight:    '╝',
		VerticalLeft:   '╠',
		VerticalRight:  '╣',
		HorizontalDown: '╦',
		HorizontalUp:   '╩',
		Cross:          '╬',
	}

	// RoundedBorder uses rounded corner characters: ╭─╮
	RoundedBorder = BorderChars{
		Horizontal:     '─',
		Vertical:       '│',
		TopLeft:        '╭',
		TopRight:       '╮',
		BottomLeft:     '╰',
		BottomRight:    '╯',
		VerticalLeft:   '├',
		VerticalRight:  '┤',
		HorizontalDown: '┬',
		HorizontalUp:   '┴',
		Cross:          '┼',
	}
)

// GetBorderChars returns the border characters for a given style.
func GetBorderChars(style theme.BorderStyle) BorderChars {
	switch style {
	case theme.BorderStyleDouble:
		return DoubleBorder
	case theme.BorderStyleRounded:
		return RoundedBorder
	case theme.BorderStyleSingle:
		fallthrough
	default:
		return SingleBorder
	}
}

// ApplyBorderStyle applies the appropriate border style to a Box based on focus state.
func ApplyBorderStyle(box *tview.Box, focused bool) {
	style := theme.GetBorderStyle(focused)
	_ = GetBorderChars(style) // Reserved for future use when tview supports custom borders

	box.SetBorder(true)
	box.SetBorderAttributes(0) // Reset attributes

	// Note: tview doesn't support custom border characters directly,
	// but we can influence the appearance through other means.
	// The actual border characters are hardcoded in tview.
	// This function is a placeholder for potential future enhancement.

	// For now, we just set the border color based on focus
	if focused {
		box.SetBorderColor(theme.GetPrimaryColor())
	} else {
		box.SetBorderColor(theme.GetSecondaryColor())
	}
}

// GradientTitle creates a gradient title for focused panels.
func GradientTitle(title string, focused bool) string {
	if !focused {
		return title
	}

	effects := theme.GetEffects()
	if !effects.GradientTitles {
		return title
	}

	// Get theme colors (for future use when we have proper color mapping)
	_ = theme.GetPrimaryColor()
	_ = theme.GetInfoColor()

	// Convert tcell colors to tview color names
	// This is a simplified approach - in a real implementation,
	// we'd need a proper color name mapping
	startColor := "blue"
	endColor := "cyan"

	// Determine colors based on theme
	current := theme.Current()
	switch current.Name {
	case "default":
		startColor = "green"
		endColor = "lime"
	case "dark":
		startColor = "blue"
		endColor = "cyan"
	case "arctic":
		startColor = "darkcyan"
		endColor = "cyan"
	}

	// Simple gradient: just apply different colors to first and second half
	if len(title) <= 2 {
		return "[" + startColor + "]" + title + "[-]"
	}

	mid := len(title) / 2
	return "[" + startColor + "]" + title[:mid] + "[-][" + endColor + "]" + title[mid:] + "[-]"
}

// PulseColor returns a color that pulses based on a frame counter.
// Used for focus pulse animation.
func PulseColor(baseColor string, frame int, maxFrames int) string {
	// Pulse between base color and a brighter variant
	// For now, we just return the base color
	// Full implementation would interpolate between colors
	return baseColor
}

// DrawFocusIndicator draws a subtle focus indicator (vertical bar on the left).
func DrawFocusIndicator() string {
	return "▌"
}

// TitleWithIcon adds an icon to a title based on content type.
func TitleWithIcon(title string, icon string) string {
	if icon == "" {
		return title
	}
	return icon + " " + title
}

// CommonIcons provides commonly used icons for TUI elements.
var CommonIcons = struct {
	Workspace string
	Ticket    string
	Detail    string
	Search    string
	Command   string
	Sync      string
	Settings  string
	Help      string
	Success   string
	Error     string
	Warning   string
	Info      string
	Loading   string
}{
	Workspace: "📁",
	Ticket:    "🎫",
	Detail:    "📄",
	Search:    "🔍",
	Command:   "⚡",
	Sync:      "🔄",
	Settings:  "⚙",
	Help:      "❓",
	Success:   "✓",
	Error:     "✗",
	Warning:   "⚠",
	Info:      "ℹ",
	Loading:   "⏳",
}
