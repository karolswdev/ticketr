package theme

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Theme represents a color scheme for the TUI.
type Theme struct {
	Name       string
	Primary    tcell.Color // Focused borders
	Secondary  tcell.Color // Unfocused borders
	Success    tcell.Color // Success messages
	Error      tcell.Color // Error messages
	Warning    tcell.Color // Warning messages
	Info       tcell.Color // Info messages
	Background tcell.Color // Background
	Text       tcell.Color // Primary text
}

var (
	// DefaultTheme is the standard green/white theme (current).
	DefaultTheme = Theme{
		Name:       "default",
		Primary:    tcell.ColorGreen,
		Secondary:  tcell.ColorWhite,
		Success:    tcell.ColorGreen,
		Error:      tcell.ColorRed,
		Warning:    tcell.ColorYellow,
		Info:       tcell.ColorTeal,
		Background: tcell.ColorDefault,
		Text:       tcell.ColorDefault,
	}

	// DarkTheme is a darker color scheme with blue accents.
	DarkTheme = Theme{
		Name:       "dark",
		Primary:    tcell.ColorBlue,
		Secondary:  tcell.ColorGray,
		Success:    tcell.ColorLime,
		Error:      tcell.ColorOrangeRed,
		Warning:    tcell.ColorGold,
		Info:       tcell.ColorAqua,
		Background: tcell.ColorDefault,
		Text:       tcell.ColorDefault,
	}

	// LightTheme is a lighter color scheme with purple accents.
	LightTheme = Theme{
		Name:       "light",
		Primary:    tcell.ColorPurple,
		Secondary:  tcell.ColorSilver,
		Success:    tcell.ColorGreen,
		Error:      tcell.ColorMaroon,
		Warning:    tcell.ColorOlive,
		Info:       tcell.ColorNavy,
		Background: tcell.ColorDefault,
		Text:       tcell.ColorDefault,
	}
)

// currentTheme is the active theme (default).
var currentTheme = DefaultTheme

// Current returns the currently active theme.
func Current() Theme {
	return currentTheme
}

// Set sets the active theme.
func Set(t Theme) {
	currentTheme = t
}

// SetByName sets the active theme by name.
// Returns true if theme was found and set, false otherwise.
func SetByName(name string) bool {
	switch name {
	case "default":
		currentTheme = DefaultTheme
		return true
	case "dark":
		currentTheme = DarkTheme
		return true
	case "light":
		currentTheme = LightTheme
		return true
	default:
		return false
	}
}

// Apply applies the current theme to global tview settings.
func Apply(app *tview.Application) {
	// Set global color scheme for tview primitives
	tview.Styles.PrimitiveBackgroundColor = currentTheme.Background
	tview.Styles.ContrastBackgroundColor = currentTheme.Background
	tview.Styles.MoreContrastBackgroundColor = currentTheme.Background
	tview.Styles.PrimaryTextColor = currentTheme.Text
	tview.Styles.SecondaryTextColor = currentTheme.Secondary
	tview.Styles.TertiaryTextColor = currentTheme.Info
	tview.Styles.InverseTextColor = currentTheme.Background
	tview.Styles.BorderColor = currentTheme.Secondary
	tview.Styles.TitleColor = currentTheme.Primary
	tview.Styles.GraphicsColor = currentTheme.Secondary

	// Update application if provided
	if app != nil {
		app.Draw()
	}
}

// GetPrimaryColor returns the primary color for focused elements.
func GetPrimaryColor() tcell.Color {
	return currentTheme.Primary
}

// GetSecondaryColor returns the secondary color for unfocused elements.
func GetSecondaryColor() tcell.Color {
	return currentTheme.Secondary
}

// GetSuccessColor returns the color for success messages.
func GetSuccessColor() tcell.Color {
	return currentTheme.Success
}

// GetErrorColor returns the color for error messages.
func GetErrorColor() tcell.Color {
	return currentTheme.Error
}

// GetWarningColor returns the color for warning messages.
func GetWarningColor() tcell.Color {
	return currentTheme.Warning
}

// GetInfoColor returns the color for info messages.
func GetInfoColor() tcell.Color {
	return currentTheme.Info
}
