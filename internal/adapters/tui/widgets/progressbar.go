package widgets

import (
	"fmt"
	"strings"
	"time"

	"github.com/karolswdev/ticktr/internal/adapters/tui/effects"
	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
	"github.com/karolswdev/ticktr/internal/tui/jobs"
)

// ProgressBarData contains all information needed to render a progress bar.
type ProgressBarData struct {
	Current    int
	Total      int
	Percentage float64
	StartTime  time.Time
	Message    string
}

// ProgressBar renders ASCII progress bars with counts, percentages, time, and ETA.
type ProgressBar struct {
	width      int
	startTime  time.Time
	lastUpdate time.Time
	shimmer    *effects.ProgressBarShimmer
}

// NewProgressBar creates a new progress bar with the specified width.
// Width should be the character width available for the progress bar itself.
func NewProgressBar(width int) *ProgressBar {
	if width < 10 {
		width = 10 // Minimum width for usability
	}

	// Check if shimmer effect is enabled in theme
	effectsConfig := theme.GetEffects()
	shimmerEnabled := effectsConfig.Motion && effectsConfig.FocusPulse

	return &ProgressBar{
		width:   width,
		shimmer: effects.NewProgressBarShimmer(width, shimmerEnabled),
	}
}

// Start marks the beginning of a progress operation.
func (pb *ProgressBar) Start() {
	now := time.Now()
	pb.startTime = now
	pb.lastUpdate = now
}

// Render creates a formatted progress bar string from progress data.
// Format: [=====>    ] 50% (45/120) | Elapsed: 12s | ETA: 15s
func (pb *ProgressBar) Render(data ProgressBarData) string {
	if data.Total <= 0 {
		// Indeterminate progress - show spinner
		return pb.renderIndeterminate(data)
	}

	// Calculate metrics
	elapsed := pb.getElapsed(data.StartTime)
	eta := pb.calculateETA(data.Current, data.Total, elapsed)

	// Build components
	bar := pb.renderBar(data.Percentage)

	// Apply shimmer effect if enabled
	if pb.shimmer != nil {
		pb.shimmer.Update()
		innerWidth := pb.width - 2
		if innerWidth > 0 {
			filledWidth := int((data.Percentage / 100.0) * float64(innerWidth))
			bar = pb.shimmer.Apply(bar, filledWidth)
		}
	}

	percentage := fmt.Sprintf("%.0f%%", data.Percentage)
	count := fmt.Sprintf("(%d/%d)", data.Current, data.Total)
	elapsedStr := pb.formatDuration(elapsed)
	etaStr := pb.formatDuration(eta)

	// Assemble full progress string
	// Format: [=====>    ] 50% (45/120) | Elapsed: 12s | ETA: 15s
	return fmt.Sprintf("%s %s %s | Elapsed: %s | ETA: %s",
		bar, percentage, count, elapsedStr, etaStr)
}

// RenderCompact creates a compact version suitable for narrow displays.
// Format: [=====>] 50% (45/120)
func (pb *ProgressBar) RenderCompact(data ProgressBarData) string {
	if data.Total <= 0 {
		return pb.renderIndeterminateCompact(data)
	}

	bar := pb.renderBar(data.Percentage)

	// Apply shimmer effect if enabled
	if pb.shimmer != nil {
		pb.shimmer.Update()
		innerWidth := pb.width - 2
		if innerWidth > 0 {
			filledWidth := int((data.Percentage / 100.0) * float64(innerWidth))
			bar = pb.shimmer.Apply(bar, filledWidth)
		}
	}

	percentage := fmt.Sprintf("%.0f%%", data.Percentage)
	count := fmt.Sprintf("(%d/%d)", data.Current, data.Total)

	return fmt.Sprintf("%s %s %s", bar, percentage, count)
}

// renderBar creates the ASCII progress bar visualization.
// Uses block characters for a clean look: [█████░░░░░]
func (pb *ProgressBar) renderBar(percentage float64) string {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	// Calculate how many filled characters we need
	// Reserve 2 chars for brackets []
	innerWidth := pb.width - 2
	if innerWidth < 1 {
		innerWidth = 1
	}

	filledWidth := int((percentage / 100.0) * float64(innerWidth))
	emptyWidth := innerWidth - filledWidth

	// Build bar using block characters
	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)

	return fmt.Sprintf("[%s%s]", filled, empty)
}

// renderIndeterminate shows a spinner for unknown progress.
func (pb *ProgressBar) renderIndeterminate(data ProgressBarData) string {
	spinner := pb.getSpinner()
	elapsed := pb.getElapsed(data.StartTime)
	elapsedStr := pb.formatDuration(elapsed)

	message := data.Message
	if message == "" {
		message = "Processing"
	}

	return fmt.Sprintf("%s %s | Elapsed: %s", spinner, message, elapsedStr)
}

// renderIndeterminateCompact shows a compact spinner.
func (pb *ProgressBar) renderIndeterminateCompact(data ProgressBarData) string {
	spinner := pb.getSpinner()
	message := data.Message
	if message == "" {
		message = "Processing"
	}

	return fmt.Sprintf("%s %s", spinner, message)
}

// getSpinner returns a rotating spinner character based on elapsed time.
// Uses Braille spinner: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
func (pb *ProgressBar) getSpinner() string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	elapsed := time.Since(pb.startTime).Milliseconds()
	index := (elapsed / 80) % int64(len(spinners)) // Rotate every 80ms
	return spinners[index]
}

// calculateETA estimates time remaining using linear extrapolation.
// Returns 0 if current is 0 or total is 0 to avoid division by zero.
func (pb *ProgressBar) calculateETA(current, total int, elapsed time.Duration) time.Duration {
	if current <= 0 || total <= 0 {
		return 0
	}

	if current >= total {
		return 0 // Already complete
	}

	// Linear extrapolation: (elapsed / current) * (total - current)
	avgTimePerItem := float64(elapsed) / float64(current)
	remaining := total - current
	eta := time.Duration(avgTimePerItem * float64(remaining))

	return eta
}

// getElapsed returns the time elapsed since the given start time.
// If start time is zero, returns time since progress bar start.
func (pb *ProgressBar) getElapsed(startTime time.Time) time.Duration {
	if startTime.IsZero() {
		return time.Since(pb.startTime)
	}
	return time.Since(startTime)
}

// formatDuration formats a duration into a human-readable string.
// Examples: "5s", "1m 23s", "2h 15m"
func (pb *ProgressBar) formatDuration(d time.Duration) string {
	if d < 0 {
		return "0s"
	}

	// Round to seconds
	seconds := int(d.Seconds())

	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	seconds = seconds % 60

	if minutes < 60 {
		if seconds > 0 {
			return fmt.Sprintf("%dm %ds", minutes, seconds)
		}
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	minutes = minutes % 60

	if minutes > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dh", hours)
}

// FromJobProgress converts a jobs.JobProgress to ProgressBarData.
func FromJobProgress(progress jobs.JobProgress, startTime time.Time) ProgressBarData {
	return ProgressBarData{
		Current:    progress.Current,
		Total:      progress.Total,
		Percentage: progress.Percentage,
		StartTime:  startTime,
		Message:    progress.Message,
	}
}

// CalculatePercentage computes completion percentage from current and total.
// Returns 0 if total is 0 to avoid division by zero.
func CalculatePercentage(current, total int) float64 {
	if total <= 0 {
		return 0
	}
	percentage := (float64(current) / float64(total)) * 100
	if percentage > 100 {
		return 100
	}
	return percentage
}
