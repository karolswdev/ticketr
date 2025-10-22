package effects

import (
	"strings"
)

// ShimmerEffect creates a shimmer animation across a progress bar.
type ShimmerEffect struct {
	offset    int
	width     int
	direction int // 1 for left-to-right, -1 for right-to-left
}

// NewShimmerEffect creates a new shimmer effect.
func NewShimmerEffect(width int) *ShimmerEffect {
	return &ShimmerEffect{
		offset:    0,
		width:     width,
		direction: 1,
	}
}

// Update advances the shimmer animation.
func (s *ShimmerEffect) Update() {
	s.offset += s.direction

	// Bounce at edges
	if s.offset >= s.width {
		s.direction = -1
	} else if s.offset <= 0 {
		s.direction = 1
	}
}

// Apply applies the shimmer effect to a progress bar string.
// The shimmer brightens a 3-character window that sweeps across the filled portion.
func (s *ShimmerEffect) Apply(barString string, filledWidth int) string {
	if filledWidth <= 0 || s.width <= 0 {
		return barString
	}

	// Find the filled portion (inside the brackets)
	startIdx := strings.Index(barString, "[")
	if startIdx == -1 {
		return barString
	}

	endIdx := strings.Index(barString, "]")
	if endIdx == -1 || endIdx <= startIdx {
		return barString
	}

	// Extract parts
	before := barString[:startIdx+1]
	bar := barString[startIdx+1 : endIdx]
	after := barString[endIdx:]

	// Apply shimmer to the filled portion
	shimmerWidth := 3
	shimmerStart := s.offset
	shimmerEnd := s.offset + shimmerWidth

	runes := []rune(bar)
	for i := 0; i < len(runes); i++ {
		// Only shimmer filled characters
		if i < filledWidth && i >= shimmerStart && i < shimmerEnd {
			// Replace with brighter version
			if runes[i] == '█' {
				runes[i] = '▓' // Slightly dimmed for shimmer effect
			}
		}
	}

	return before + string(runes) + after
}

// Reset resets the shimmer effect to the beginning.
func (s *ShimmerEffect) Reset() {
	s.offset = 0
	s.direction = 1
}

// ProgressBarShimmer wraps progress bar rendering with shimmer effect.
type ProgressBarShimmer struct {
	shimmer *ShimmerEffect
	enabled bool
}

// NewProgressBarShimmer creates a new progress bar shimmer wrapper.
func NewProgressBarShimmer(width int, enabled bool) *ProgressBarShimmer {
	return &ProgressBarShimmer{
		shimmer: NewShimmerEffect(width),
		enabled: enabled,
	}
}

// SetEnabled enables or disables the shimmer effect.
func (pbs *ProgressBarShimmer) SetEnabled(enabled bool) {
	pbs.enabled = enabled
}

// Update advances the shimmer animation.
func (pbs *ProgressBarShimmer) Update() {
	if pbs.enabled {
		pbs.shimmer.Update()
	}
}

// Apply applies the shimmer effect to a progress bar string.
func (pbs *ProgressBarShimmer) Apply(barString string, filledWidth int) string {
	if !pbs.enabled {
		return barString
	}
	return pbs.shimmer.Apply(barString, filledWidth)
}

// GradientText creates a horizontal gradient effect for text.
// Used for title gradients in focused panels.
func GradientText(text string, startColor, endColor string) string {
	if len(text) == 0 {
		return text
	}

	// Simple two-color gradient using tview color tags
	// For a smoother gradient, we'd need to interpolate colors per character
	// For now, we use a simple transition at midpoint
	midpoint := len(text) / 2

	var result strings.Builder
	result.WriteString("[" + startColor + "]")

	for i, ch := range text {
		if i == midpoint {
			result.WriteString("[-][" + endColor + "]")
		}
		result.WriteRune(ch)
	}

	result.WriteString("[-]")
	return result.String()
}

// PulseIntensity calculates the intensity for a pulse effect.
// Returns a value between 0.0 and 1.0 based on the current frame.
func PulseIntensity(frame, maxFrames int) float64 {
	if maxFrames <= 0 {
		return 1.0
	}

	// Sinusoidal pulse
	progress := float64(frame%maxFrames) / float64(maxFrames)
	// Map to 0.5 - 1.0 range (never fully dim)
	return 0.5 + 0.5*progress
}

// Rainbow creates a rainbow color effect for text.
// Cycles through colors: red, yellow, green, cyan, blue, magenta.
func Rainbow(text string) string {
	colors := []string{"red", "yellow", "green", "cyan", "blue", "magenta"}
	if len(text) == 0 {
		return text
	}

	var result strings.Builder
	colorIdx := 0

	for _, ch := range text {
		result.WriteString("[" + colors[colorIdx%len(colors)] + "]")
		result.WriteRune(ch)
		if ch != ' ' {
			colorIdx++
		}
	}

	result.WriteString("[-]")
	return result.String()
}
