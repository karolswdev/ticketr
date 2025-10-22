package theme

import (
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// BorderStyle represents the style of borders for panels.
type BorderStyle string

const (
	BorderStyleSingle  BorderStyle = "single"  // ┌─┐
	BorderStyleDouble  BorderStyle = "double"  // ╔═╗
	BorderStyleRounded BorderStyle = "rounded" // ╭─╮
)

// VisualEffects contains configuration for visual polish features.
type VisualEffects struct {
	// Motion effects
	Motion      bool // Global motion kill switch
	Spinner     bool // Active spinners
	FocusPulse  bool // Focus pulse animation
	ModalFadeIn bool // Modal fade-in effect

	// Visual polish
	DropShadows    bool // Drop shadows on modals
	GradientTitles bool // Gradient titles on focused panels

	// Border styles
	FocusedBorder   BorderStyle
	UnfocusedBorder BorderStyle

	// Ambient effects (default OFF)
	AmbientEnabled bool
	AmbientMode    string  // "hyperspace", "snow", "off"
	AmbientDensity float64 // 0.01 - 0.10 (1% - 10%)
	AmbientSpeed   int     // Milliseconds per frame
}

// DefaultVisualEffects returns conservative defaults (most effects OFF).
func DefaultVisualEffects() VisualEffects {
	return VisualEffects{
		Motion:          true,  // Motion enabled, but individual effects opt-in
		Spinner:         true,  // Spinners are essential feedback
		FocusPulse:      false, // OFF by default
		ModalFadeIn:     false, // OFF by default
		DropShadows:     true,  // ENABLED for professional appearance
		GradientTitles:  false, // OFF by default
		FocusedBorder:   BorderStyleDouble,
		UnfocusedBorder: BorderStyleSingle,
		AmbientEnabled:  false,
		AmbientMode:     "off",
		AmbientDensity:  0.02,
		AmbientSpeed:    100,
	}
}

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

	// Visual effects configuration
	Effects VisualEffects

	// Animation character sets
	SpinnerFrames []string
	SparkleFrames []string
}

var (
	// DefaultTheme is the standard green/white theme (current).
	// Effects default to OFF for conservative, accessible experience.
	DefaultTheme = Theme{
		Name:          "default",
		Primary:       tcell.ColorGreen,
		Secondary:     tcell.ColorWhite,
		Success:       tcell.ColorGreen,
		Error:         tcell.ColorRed,
		Warning:       tcell.ColorYellow,
		Info:          tcell.ColorTeal,
		Background:    tcell.ColorDefault,
		Text:          tcell.ColorDefault,
		Effects:       DefaultVisualEffects(),
		SpinnerFrames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		SparkleFrames: []string{"✦", "✧", "⋆", "∗", "·"},
	}

	// DarkTheme is a darker color scheme with blue accents.
	// Includes hyperspace ambient effect (OFF by default, but available).
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
		Effects: VisualEffects{
			Motion:          true,
			Spinner:         true,
			FocusPulse:      false,
			ModalFadeIn:     false,
			DropShadows:     false,
			GradientTitles:  false,
			FocusedBorder:   BorderStyleDouble,
			UnfocusedBorder: BorderStyleSingle,
			AmbientEnabled:  false,
			AmbientMode:     "hyperspace",
			AmbientDensity:  0.02,
			AmbientSpeed:    100,
		},
		SpinnerFrames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		SparkleFrames: []string{"✦", "✧", "⋆", "∗", "·"},
	}

	// ArcticTheme is a light blue theme with snow effect.
	// Includes snow ambient effect (OFF by default, but available).
	ArcticTheme = Theme{
		Name:       "arctic",
		Primary:    tcell.ColorDarkCyan,
		Secondary:  tcell.ColorLightCyan,
		Success:    tcell.ColorGreen,
		Error:      tcell.ColorRed,
		Warning:    tcell.ColorOrange,
		Info:       tcell.ColorBlue,
		Background: tcell.ColorDefault,
		Text:       tcell.ColorDefault,
		Effects: VisualEffects{
			Motion:          true,
			Spinner:         true,
			FocusPulse:      false,
			ModalFadeIn:     false,
			DropShadows:     false,
			GradientTitles:  false,
			FocusedBorder:   BorderStyleDouble,
			UnfocusedBorder: BorderStyleRounded,
			AmbientEnabled:  false,
			AmbientMode:     "snow",
			AmbientDensity:  0.015,
			AmbientSpeed:    120,
		},
		SpinnerFrames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		SparkleFrames: []string{"❄", "❅", "❆", "∗", "·"},
	}

	// LightTheme is a lighter color scheme with purple accents (legacy name kept for compatibility).
	LightTheme = ArcticTheme
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
	case "light", "arctic":
		currentTheme = ArcticTheme
		return true
	default:
		return false
	}
}

// GetEffects returns the visual effects configuration for the current theme.
func GetEffects() VisualEffects {
	return currentTheme.Effects
}

// SetEffects updates the visual effects configuration for the current theme.
func SetEffects(effects VisualEffects) {
	currentTheme.Effects = effects
}

// GetBorderStyle returns the border style for the given focus state.
func GetBorderStyle(focused bool) BorderStyle {
	if focused {
		return currentTheme.Effects.FocusedBorder
	}
	return currentTheme.Effects.UnfocusedBorder
}

// GetSpinnerFrames returns the spinner frames for the current theme.
func GetSpinnerFrames() []string {
	if len(currentTheme.SpinnerFrames) == 0 {
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	}
	return currentTheme.SpinnerFrames
}

// GetSparkleFrames returns the sparkle frames for the current theme.
func GetSparkleFrames() []string {
	if len(currentTheme.SparkleFrames) == 0 {
		return []string{"✦", "✧", "⋆", "∗", "·"}
	}
	return currentTheme.SparkleFrames
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

// LoadThemeFromEnv loads theme configuration from environment variables.
// Supports the following environment variables:
//
//	TICKETR_THEME - Theme name (default, dark, arctic)
//	TICKETR_EFFECTS_MOTION - Enable/disable motion (true/false)
//	TICKETR_EFFECTS_SHADOWS - Enable/disable drop shadows (true/false)
//	TICKETR_EFFECTS_SHIMMER - Enable/disable shimmer effect (true/false)
//	TICKETR_EFFECTS_AMBIENT - Enable/disable ambient effects (true/false)
func LoadThemeFromEnv() {
	// Load theme by name
	themeName := os.Getenv("TICKETR_THEME")
	if themeName != "" {
		SetByName(strings.ToLower(themeName))
	}

	// Load effects configuration
	effects := currentTheme.Effects

	if motionStr := os.Getenv("TICKETR_EFFECTS_MOTION"); motionStr != "" {
		if motion, err := strconv.ParseBool(motionStr); err == nil {
			effects.Motion = motion
		}
	}

	if shadowsStr := os.Getenv("TICKETR_EFFECTS_SHADOWS"); shadowsStr != "" {
		if shadows, err := strconv.ParseBool(shadowsStr); err == nil {
			effects.DropShadows = shadows
		}
	}

	if shimmerStr := os.Getenv("TICKETR_EFFECTS_SHIMMER"); shimmerStr != "" {
		if shimmer, err := strconv.ParseBool(shimmerStr); err == nil {
			effects.FocusPulse = shimmer
		}
	}

	if ambientStr := os.Getenv("TICKETR_EFFECTS_AMBIENT"); ambientStr != "" {
		if ambient, err := strconv.ParseBool(ambientStr); err == nil {
			effects.AmbientEnabled = ambient
		}
	}

	// Apply the updated effects
	SetEffects(effects)
}

// LoadThemeFromFlags loads theme configuration from command-line flags.
// This is a helper for CLI integration.
func LoadThemeFromFlags(themeName string, enableEffects bool) {
	if themeName != "" {
		SetByName(strings.ToLower(themeName))
	}

	if !enableEffects {
		// Disable all effects
		effects := currentTheme.Effects
		effects.Motion = false
		effects.DropShadows = false
		effects.FocusPulse = false
		effects.AmbientEnabled = false
		SetEffects(effects)
	}
}
