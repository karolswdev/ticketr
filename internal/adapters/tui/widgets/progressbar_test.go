package widgets

import (
	"strings"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/tui/jobs"
)

func TestNewProgressBar(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		expectedWidth int
	}{
		{
			name:          "normal width",
			width:         50,
			expectedWidth: 50,
		},
		{
			name:          "minimum width enforced",
			width:         5,
			expectedWidth: 10,
		},
		{
			name:          "zero width gets minimum",
			width:         0,
			expectedWidth: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := NewProgressBar(tt.width)
			if pb.width != tt.expectedWidth {
				t.Errorf("expected width %d, got %d", tt.expectedWidth, pb.width)
			}
		})
	}
}

func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		expected float64
	}{
		{
			name:     "50% complete",
			current:  50,
			total:    100,
			expected: 50.0,
		},
		{
			name:     "100% complete",
			current:  100,
			total:    100,
			expected: 100.0,
		},
		{
			name:     "0% complete",
			current:  0,
			total:    100,
			expected: 0.0,
		},
		{
			name:     "over 100% capped",
			current:  150,
			total:    100,
			expected: 100.0,
		},
		{
			name:     "zero total returns 0",
			current:  50,
			total:    0,
			expected: 0.0,
		},
		{
			name:     "fractional percentage",
			current:  1,
			total:    3,
			expected: 33.33333333333333,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePercentage(tt.current, tt.total)
			if result != tt.expected {
				t.Errorf("expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestRenderBar(t *testing.T) {
	pb := NewProgressBar(20)

	tests := []struct {
		name       string
		percentage float64
		wantFilled int
		wantEmpty  int
	}{
		{
			name:       "0% progress",
			percentage: 0,
			wantFilled: 0,
			wantEmpty:  18, // 20 - 2 (brackets)
		},
		{
			name:       "50% progress",
			percentage: 50,
			wantFilled: 9,
			wantEmpty:  9,
		},
		{
			name:       "100% progress",
			percentage: 100,
			wantFilled: 18,
			wantEmpty:  0,
		},
		{
			name:       "negative percentage treated as 0",
			percentage: -10,
			wantFilled: 0,
			wantEmpty:  18,
		},
		{
			name:       "over 100% capped at 100",
			percentage: 150,
			wantFilled: 18,
			wantEmpty:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := pb.renderBar(tt.percentage)

			// Check that bar has brackets
			if !strings.HasPrefix(bar, "[") || !strings.HasSuffix(bar, "]") {
				t.Errorf("bar should have brackets: %s", bar)
			}

			// Extract inner bar content
			inner := strings.TrimPrefix(strings.TrimSuffix(bar, "]"), "[")

			// Count filled and empty characters
			filled := strings.Count(inner, "█")
			empty := strings.Count(inner, "░")

			if filled != tt.wantFilled {
				t.Errorf("expected %d filled chars, got %d (bar: %s)", tt.wantFilled, filled, bar)
			}
			if empty != tt.wantEmpty {
				t.Errorf("expected %d empty chars, got %d (bar: %s)", tt.wantEmpty, empty, bar)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	pb := NewProgressBar(50)

	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "0 seconds",
			duration: 0,
			expected: "0s",
		},
		{
			name:     "5 seconds",
			duration: 5 * time.Second,
			expected: "5s",
		},
		{
			name:     "59 seconds",
			duration: 59 * time.Second,
			expected: "59s",
		},
		{
			name:     "1 minute",
			duration: 60 * time.Second,
			expected: "1m",
		},
		{
			name:     "1 minute 30 seconds",
			duration: 90 * time.Second,
			expected: "1m 30s",
		},
		{
			name:     "1 hour",
			duration: time.Hour,
			expected: "1h",
		},
		{
			name:     "1 hour 15 minutes",
			duration: time.Hour + 15*time.Minute,
			expected: "1h 15m",
		},
		{
			name:     "2 hours 5 minutes",
			duration: 2*time.Hour + 5*time.Minute,
			expected: "2h 5m",
		},
		{
			name:     "negative duration treated as 0",
			duration: -5 * time.Second,
			expected: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pb.formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCalculateETA(t *testing.T) {
	pb := NewProgressBar(50)

	tests := []struct {
		name        string
		current     int
		total       int
		elapsed     time.Duration
		expectedETA time.Duration
	}{
		{
			name:        "50% complete, linear extrapolation",
			current:     50,
			total:       100,
			elapsed:     10 * time.Second,
			expectedETA: 10 * time.Second, // Same time remaining as elapsed
		},
		{
			name:        "25% complete",
			current:     25,
			total:       100,
			elapsed:     10 * time.Second,
			expectedETA: 30 * time.Second, // 3x elapsed for remaining 75%
		},
		{
			name:        "90% complete, almost done",
			current:     90,
			total:       100,
			elapsed:     90 * time.Second,
			expectedETA: 10 * time.Second, // 1s per item, 10 remaining
		},
		{
			name:        "100% complete",
			current:     100,
			total:       100,
			elapsed:     100 * time.Second,
			expectedETA: 0, // Already done
		},
		{
			name:        "zero current returns 0",
			current:     0,
			total:       100,
			elapsed:     10 * time.Second,
			expectedETA: 0,
		},
		{
			name:        "zero total returns 0",
			current:     50,
			total:       0,
			elapsed:     10 * time.Second,
			expectedETA: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eta := pb.calculateETA(tt.current, tt.total, tt.elapsed)

			// Allow small rounding differences (within 1 second)
			diff := eta - tt.expectedETA
			if diff < 0 {
				diff = -diff
			}
			if diff > time.Second {
				t.Errorf("expected ETA %v, got %v (diff: %v)", tt.expectedETA, eta, diff)
			}
		})
	}
}

func TestRender(t *testing.T) {
	pb := NewProgressBar(20)
	pb.Start()

	tests := []struct {
		name           string
		data           ProgressBarData
		shouldContain  []string
		shouldNotExist []string
	}{
		{
			name: "normal progress",
			data: ProgressBarData{
				Current:    50,
				Total:      100,
				Percentage: 50.0,
				StartTime:  time.Now().Add(-10 * time.Second),
			},
			shouldContain: []string{
				"[", "]", // Bar brackets
				"50%",      // Percentage
				"(50/100)", // Count
				"Elapsed:",
				"ETA:",
			},
		},
		{
			name: "0% progress",
			data: ProgressBarData{
				Current:    0,
				Total:      100,
				Percentage: 0.0,
				StartTime:  time.Now(),
			},
			shouldContain: []string{
				"0%",
				"(0/100)",
			},
		},
		{
			name: "100% progress",
			data: ProgressBarData{
				Current:    100,
				Total:      100,
				Percentage: 100.0,
				StartTime:  time.Now().Add(-60 * time.Second),
			},
			shouldContain: []string{
				"100%",
				"(100/100)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pb.Render(tt.data)

			for _, expected := range tt.shouldContain {
				if !strings.Contains(result, expected) {
					t.Errorf("expected result to contain %q, got: %s", expected, result)
				}
			}

			for _, notExpected := range tt.shouldNotExist {
				if strings.Contains(result, notExpected) {
					t.Errorf("expected result to NOT contain %q, got: %s", notExpected, result)
				}
			}
		})
	}
}

func TestRenderCompact(t *testing.T) {
	pb := NewProgressBar(20)
	pb.Start()

	data := ProgressBarData{
		Current:    50,
		Total:      100,
		Percentage: 50.0,
		StartTime:  time.Now(),
	}

	result := pb.RenderCompact(data)

	// Compact should have bar, percentage, count
	expected := []string{"[", "]", "50%", "(50/100)"}
	for _, exp := range expected {
		if !strings.Contains(result, exp) {
			t.Errorf("compact render should contain %q, got: %s", exp, result)
		}
	}

	// Compact should NOT have time info
	notExpected := []string{"Elapsed:", "ETA:"}
	for _, notExp := range notExpected {
		if strings.Contains(result, notExp) {
			t.Errorf("compact render should NOT contain %q, got: %s", notExp, result)
		}
	}
}

func TestRenderIndeterminate(t *testing.T) {
	pb := NewProgressBar(20)
	pb.Start()

	data := ProgressBarData{
		Current:   0,
		Total:     0, // Unknown total
		StartTime: time.Now(),
		Message:   "Fetching tickets",
	}

	result := pb.renderIndeterminate(data)

	// Should show spinner and message
	if !strings.Contains(result, "Fetching tickets") {
		t.Errorf("indeterminate should show message, got: %s", result)
	}
	if !strings.Contains(result, "Elapsed:") {
		t.Errorf("indeterminate should show elapsed time, got: %s", result)
	}
}

func TestFromJobProgress(t *testing.T) {
	startTime := time.Now().Add(-30 * time.Second)

	jobProgress := jobs.JobProgress{
		JobID:      "test-job",
		Current:    45,
		Total:      120,
		Percentage: 37.5,
		Message:    "Processing tickets",
	}

	data := FromJobProgress(jobProgress, startTime)

	if data.Current != 45 {
		t.Errorf("expected current 45, got %d", data.Current)
	}
	if data.Total != 120 {
		t.Errorf("expected total 120, got %d", data.Total)
	}
	if data.Percentage != 37.5 {
		t.Errorf("expected percentage 37.5, got %.1f", data.Percentage)
	}
	if data.Message != "Processing tickets" {
		t.Errorf("expected message 'Processing tickets', got %q", data.Message)
	}
	if data.StartTime != startTime {
		t.Errorf("expected start time to match")
	}
}

func TestGetSpinner(t *testing.T) {
	pb := NewProgressBar(50)
	pb.Start()

	// Get spinner at different times
	spinner1 := pb.getSpinner()
	if spinner1 == "" {
		t.Error("spinner should not be empty")
	}

	// Spinner should be one of the expected characters
	validSpinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	found := false
	for _, valid := range validSpinners {
		if spinner1 == valid {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("spinner %q not in valid set", spinner1)
	}
}

func TestProgressBarWidth(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{"narrow", 15},
		{"medium", 30},
		{"wide", 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := NewProgressBar(tt.width)
			data := ProgressBarData{
				Current:    50,
				Total:      100,
				Percentage: 50.0,
				StartTime:  time.Now(),
			}

			result := pb.RenderCompact(data)

			// Extract bar portion
			startIdx := strings.Index(result, "[")
			endIdx := strings.Index(result, "]")
			if startIdx == -1 || endIdx == -1 {
				t.Fatalf("bar not found in result: %s", result)
			}

			// Count runes (characters) not bytes
			barPortion := result[startIdx : endIdx+1]
			barLength := len([]rune(barPortion))
			if barLength != tt.width {
				t.Errorf("expected bar length %d, got %d (result: %s)", tt.width, barLength, result)
			}
		})
	}
}
