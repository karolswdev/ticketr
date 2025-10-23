// Package theme provides the comprehensive theme system for the Ticketr Bubbletea TUI.
// It defines three beautiful themes: Default (Green/Terminal), Dark (Blue/Modern), and Arctic (Cyan/Cool).
package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a complete visual theme for the TUI
type Theme struct {
	Name        string
	Background  lipgloss.AdaptiveColor
	Foreground  lipgloss.AdaptiveColor
	Primary     lipgloss.AdaptiveColor
	Secondary   lipgloss.AdaptiveColor
	Accent      lipgloss.AdaptiveColor
	Success     lipgloss.AdaptiveColor
	Warning     lipgloss.AdaptiveColor
	Error       lipgloss.AdaptiveColor
	Muted       lipgloss.AdaptiveColor
	Border      lipgloss.AdaptiveColor
	BorderFocus lipgloss.AdaptiveColor

	// Border styles
	BorderStyle      lipgloss.Border
	BorderFocusStyle lipgloss.Border
}

var (
	// DefaultTheme - Green/Terminal (Midnight Commander inspired)
	DefaultTheme = Theme{
		Name:       "Default",
		Background: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"},
		Foreground: lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"},
		Primary:    lipgloss.AdaptiveColor{Light: "#00AA00", Dark: "#00FF00"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#006600", Dark: "#00AA00"},
		Accent:     lipgloss.AdaptiveColor{Light: "#CC9900", Dark: "#FFFF00"},
		Success:    lipgloss.AdaptiveColor{Light: "#00AA00", Dark: "#00FF00"},
		Warning:    lipgloss.AdaptiveColor{Light: "#CC6600", Dark: "#FFA500"},
		Error:      lipgloss.AdaptiveColor{Light: "#CC0000", Dark: "#FF0000"},
		Muted:      lipgloss.AdaptiveColor{Light: "#808080", Dark: "#666666"},
		Border:     lipgloss.AdaptiveColor{Light: "#006600", Dark: "#006600"},
		BorderFocus: lipgloss.AdaptiveColor{Light: "#00AA00", Dark: "#00FF00"},

		BorderStyle:      lipgloss.NormalBorder(),
		BorderFocusStyle: lipgloss.DoubleBorder(),
	}

	// DarkTheme - Blue/Modern (sleek and professional)
	DarkTheme = Theme{
		Name:       "Dark",
		Background: lipgloss.AdaptiveColor{Light: "#F5F5F5", Dark: "#1E1E1E"},
		Foreground: lipgloss.AdaptiveColor{Light: "#2D2D2D", Dark: "#ABB2BF"},
		Primary:    lipgloss.AdaptiveColor{Light: "#0066CC", Dark: "#61AFEF"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#0052A3", Dark: "#528BFF"},
		Accent:     lipgloss.AdaptiveColor{Light: "#8B3DAB", Dark: "#C678DD"},
		Success:    lipgloss.AdaptiveColor{Light: "#00AA00", Dark: "#98C379"},
		Warning:    lipgloss.AdaptiveColor{Light: "#CC8800", Dark: "#E5C07B"},
		Error:      lipgloss.AdaptiveColor{Light: "#CC0000", Dark: "#E06C75"},
		Muted:      lipgloss.AdaptiveColor{Light: "#888888", Dark: "#5C6370"},
		Border:     lipgloss.AdaptiveColor{Light: "#DDDDDD", Dark: "#3E4451"},
		BorderFocus: lipgloss.AdaptiveColor{Light: "#0066CC", Dark: "#61AFEF"},

		BorderStyle:      lipgloss.NormalBorder(),
		BorderFocusStyle: lipgloss.DoubleBorder(),
	}

	// ArcticTheme - Cyan/Cool (crisp and refreshing)
	ArcticTheme = Theme{
		Name:       "Arctic",
		Background: lipgloss.AdaptiveColor{Light: "#F0F9FF", Dark: "#0A1628"},
		Foreground: lipgloss.AdaptiveColor{Light: "#0F1419", Dark: "#E0F2FE"},
		Primary:    lipgloss.AdaptiveColor{Light: "#0891B2", Dark: "#00FFFF"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#0E7490", Dark: "#00AAAA"},
		Accent:     lipgloss.AdaptiveColor{Light: "#06B6D4", Dark: "#A5F3FC"},
		Success:    lipgloss.AdaptiveColor{Light: "#10B981", Dark: "#34D399"},
		Warning:    lipgloss.AdaptiveColor{Light: "#F59E0B", Dark: "#FBBF24"},
		Error:      lipgloss.AdaptiveColor{Light: "#EF4444", Dark: "#F87171"},
		Muted:      lipgloss.AdaptiveColor{Light: "#64748B", Dark: "#475569"},
		Border:     lipgloss.AdaptiveColor{Light: "#CBD5E1", Dark: "#164E63"},
		BorderFocus: lipgloss.AdaptiveColor{Light: "#0891B2", Dark: "#00FFFF"},

		BorderStyle:      lipgloss.RoundedBorder(),
		BorderFocusStyle: lipgloss.DoubleBorder(),
	}

	// All available themes
	AllThemes = []*Theme{&DefaultTheme, &DarkTheme, &ArcticTheme}
)

// GetByName returns a theme by name (Default, Dark, Arctic)
// Returns DefaultTheme if theme is not found
func GetByName(name string) *Theme {
	for _, t := range AllThemes {
		if t.Name == name {
			return t
		}
	}
	return &DefaultTheme
}

// Next returns the next theme in the cycle
// Takes the current theme and returns the next one
func Next(current *Theme) *Theme {
	for i, t := range AllThemes {
		if t.Name == current.Name {
			return AllThemes[(i+1)%len(AllThemes)]
		}
	}
	return &DefaultTheme
}

// GetAllThemeNames returns a slice of all available theme names
func GetAllThemeNames() []string {
	names := make([]string, len(AllThemes))
	for i, t := range AllThemes {
		names[i] = t.Name
	}
	return names
}
