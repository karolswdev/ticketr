// Package theme provides comprehensive Lipgloss style catalog for all UI elements
package theme

import "github.com/charmbracelet/lipgloss"

// GetPaletteForTheme returns the color palette for a given theme
func GetPaletteForTheme(t *Theme) Palette {
	return GetPalette(t.Name)
}

// Border Styles

// BorderStyleFocused returns the style for focused panels
func BorderStyleFocused(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(p.BorderFocused))
}

// BorderStyleBlurred returns the style for blurred/unfocused panels
func BorderStyleBlurred(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(p.BorderBlur))
}

// BorderStyleRounded returns a rounded border style
func BorderStyleRounded(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(p.Border))
}

// Text Styles

// TitleStyle returns the style for panel titles
func TitleStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Primary)).
		Bold(true).
		Padding(0, 1)
}

// SubtitleStyle returns the style for subtitles and section headers
func SubtitleStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Secondary)).
		Bold(true)
}

// BodyTextStyle returns the default text style
func BodyTextStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Foreground))
}

// HelpTextStyle returns the style for help/hint text
func HelpTextStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Muted)).
		Italic(true)
}

// Component Styles

// ActionBarStyle returns the style for the action bar
func ActionBarStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Background(lipgloss.Color(p.Background)).
		Foreground(lipgloss.Color(p.Foreground)).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color(p.Border))
}

// StatusBarStyle returns the style for the status bar
func StatusBarStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Background(lipgloss.Color(p.Background)).
		Foreground(lipgloss.Color(p.Foreground)).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color(p.Border))
}

// ModalBackdropStyle returns the backdrop style for modals
func ModalBackdropStyle(t *Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Foreground(lipgloss.Color("#FFFFFF"))
}

// ModalStyle returns the style for modal dialogs
func ModalStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(p.Primary)).
		Background(lipgloss.Color(p.Background)).
		Foreground(lipgloss.Color(p.Foreground)).
		Padding(1, 2)
}

// State Styles

// SuccessStyle returns the style for success messages
func SuccessStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Success)).
		Bold(true)
}

// ErrorStyle returns the style for error messages
func ErrorStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Error)).
		Bold(true)
}

// WarningStyle returns the style for warning messages
func WarningStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Warning)).
		Bold(true)
}

// InfoStyle returns the style for info messages
func InfoStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Info))
}

// Tree Styles

// TreeItemStyle returns the default style for tree items
func TreeItemStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Foreground))
}

// TreeItemSelectedStyle returns the style for selected tree items
func TreeItemSelectedStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Primary)).
		Background(lipgloss.Color(p.Selection)).
		Bold(true)
}

// TreeItemMatchedStyle returns the style for search-matched tree items
func TreeItemMatchedStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Accent)).
		Bold(true)
}

// Panel Styles

// PanelStyle returns the base panel style
func PanelStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Background(lipgloss.Color(p.Background)).
		Foreground(lipgloss.Color(p.Foreground)).
		Padding(1)
}

// PanelTitleStyle returns the style for panel titles
func PanelTitleStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Primary)).
		Background(lipgloss.Color(p.Background)).
		Bold(true).
		Padding(0, 1)
}

// List Styles

// ListItemStyle returns the default list item style
func ListItemStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Foreground)).
		Padding(0, 1)
}

// ListItemSelectedStyle returns the style for selected list items
func ListItemSelectedStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Primary)).
		Background(lipgloss.Color(p.Selection)).
		Bold(true).
		Padding(0, 1)
}

// Input Styles

// InputStyle returns the style for text inputs
func InputStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Foreground)).
		Background(lipgloss.Color(p.Background)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(p.Border)).
		Padding(0, 1)
}

// InputFocusedStyle returns the style for focused text inputs
func InputFocusedStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Foreground)).
		Background(lipgloss.Color(p.Background)).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(p.Primary)).
		Padding(0, 1)
}

// Button Styles

// ButtonStyle returns the default button style
func ButtonStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Foreground)).
		Background(lipgloss.Color(p.Secondary)).
		Padding(0, 2).
		MarginRight(1)
}

// ButtonFocusedStyle returns the style for focused buttons
func ButtonFocusedStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Background)).
		Background(lipgloss.Color(p.Primary)).
		Bold(true).
		Padding(0, 2).
		MarginRight(1)
}

// Badge Styles

// BadgeStyle returns the style for status badges
func BadgeStyle(color string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(color)).
		Padding(0, 1).
		Bold(true)
}

// StatusBadgeStyle returns status-specific badge styles
func StatusBadgeStyle(t *Theme, status string) lipgloss.Style {
	p := GetPaletteForTheme(t)
	var color string
	switch status {
	case "synced":
		color = p.SyncedStatus
	case "local":
		color = p.LocalStatus
	case "dirty":
		color = p.DirtyStatus
	case "error":
		color = p.Error
	default:
		color = p.Muted
	}
	return BadgeStyle(color)
}

// Header Styles

// HeaderStyle returns the style for the main header
func HeaderStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Background(lipgloss.Color(p.Background)).
		Foreground(lipgloss.Color(p.Primary)).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color(p.Border))
}

// Footer Styles

// FooterStyle returns the style for the footer/action bar
func FooterStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Background(lipgloss.Color(p.Background)).
		Foreground(lipgloss.Color(p.Foreground)).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color(p.Border))
}

// KeybindingStyle returns the style for keybinding hints
func KeybindingStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Accent)).
		Bold(true)
}

// KeybindingLabelStyle returns the style for keybinding labels
func KeybindingLabelStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Muted))
}

// Progress Styles

// ProgressBarStyle returns the style for progress bars
func ProgressBarStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Success))
}

// ProgressBarEmptyStyle returns the style for empty progress bar sections
func ProgressBarEmptyStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Muted))
}

// Table Styles

// TableHeaderStyle returns the style for table headers
func TableHeaderStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Primary)).
		Background(lipgloss.Color(p.Background)).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color(p.Border))
}

// TableRowStyle returns the style for table rows
func TableRowStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Foreground)).
		Padding(0, 1)
}

// TableRowSelectedStyle returns the style for selected table rows
func TableRowSelectedStyle(t *Theme) lipgloss.Style {
	p := GetPaletteForTheme(t)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Primary)).
		Background(lipgloss.Color(p.Selection)).
		Bold(true).
		Padding(0, 1)
}
