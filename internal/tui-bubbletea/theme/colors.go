// Package theme provides color palettes for all themes
package theme

import "github.com/charmbracelet/lipgloss"

// Palette defines a complete color scheme for the TUI
type Palette struct {
	Name string

	// Base colors
	Primary       string
	Secondary     string
	Background    string
	Foreground    string

	// Border colors
	Border        string
	BorderFocused string
	BorderBlur    string

	// Accent colors
	Accent        string
	Error         string
	Success       string
	Warning       string
	Info          string
	Muted         string

	// UI element colors
	Highlight     string
	Selection     string
	Cursor        string

	// Status colors
	SyncedStatus  string
	LocalStatus   string
	DirtyStatus   string
}

// DefaultPalette - Green/Terminal theme (Midnight Commander inspired)
var DefaultPalette = Palette{
	Name:          "Default",
	Primary:       "#00FF00",
	Secondary:     "#00AA00",
	Background:    "#000000",
	Foreground:    "#FFFFFF",
	Border:        "#006600",
	BorderFocused: "#00FF00",
	BorderBlur:    "#003300",
	Accent:        "#FFFF00",
	Error:         "#FF0000",
	Success:       "#00FF00",
	Warning:       "#FFA500",
	Info:          "#00FFFF",
	Muted:         "#666666",
	Highlight:     "#00FF00",
	Selection:     "#00AA00",
	Cursor:        "#00FF00",
	SyncedStatus:  "#00FF00",
	LocalStatus:   "#AAAAAA",
	DirtyStatus:   "#FFFF00",
}

// DarkPalette - Blue/Modern theme
var DarkPalette = Palette{
	Name:          "Dark",
	Primary:       "#61AFEF",
	Secondary:     "#528BFF",
	Background:    "#1E1E1E",
	Foreground:    "#ABB2BF",
	Border:        "#3E4451",
	BorderFocused: "#61AFEF",
	BorderBlur:    "#2C313C",
	Accent:        "#C678DD",
	Error:         "#E06C75",
	Success:       "#98C379",
	Warning:       "#E5C07B",
	Info:          "#56B6C2",
	Muted:         "#5C6370",
	Highlight:     "#61AFEF",
	Selection:     "#3E4451",
	Cursor:        "#528BFF",
	SyncedStatus:  "#98C379",
	LocalStatus:   "#5C6370",
	DirtyStatus:   "#E5C07B",
}

// ArcticPalette - Cyan/Cool theme
var ArcticPalette = Palette{
	Name:          "Arctic",
	Primary:       "#00FFFF",
	Secondary:     "#00AAAA",
	Background:    "#0A1628",
	Foreground:    "#E0F2FE",
	Border:        "#164E63",
	BorderFocused: "#00FFFF",
	BorderBlur:    "#083344",
	Accent:        "#A5F3FC",
	Error:         "#F87171",
	Success:       "#34D399",
	Warning:       "#FBBF24",
	Info:          "#60A5FA",
	Muted:         "#475569",
	Highlight:     "#00FFFF",
	Selection:     "#164E63",
	Cursor:        "#00FFFF",
	SyncedStatus:  "#34D399",
	LocalStatus:   "#475569",
	DirtyStatus:   "#FBBF24",
}

// GetPalette returns the palette for the given theme name
func GetPalette(name string) Palette {
	switch name {
	case "Dark":
		return DarkPalette
	case "Arctic":
		return ArcticPalette
	default:
		return DefaultPalette
	}
}

// ToAdaptiveColor converts a hex color to a lipgloss AdaptiveColor
func (p Palette) ToAdaptiveColor(hex string) lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{Light: hex, Dark: hex}
}
