package widgets

import (
	"strings"
	"testing"
	"time"
)

func TestMarquee_NeedsScrolling(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected bool
	}{
		{
			name:     "text fits exactly",
			text:     "hello",
			width:    5,
			expected: false,
		},
		{
			name:     "text is shorter",
			text:     "hi",
			width:    10,
			expected: false,
		},
		{
			name:     "text overflows",
			text:     "this is a very long text that will overflow",
			width:    20,
			expected: true,
		},
		{
			name:     "text with color codes fits",
			text:     "[yellow]hello[white]",
			width:    5,
			expected: false,
		},
		{
			name:     "text with color codes overflows",
			text:     "[yellow]this is a very long text[white]",
			width:    10,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMarquee(tt.text, tt.width)
			defer m.Shutdown()

			got := m.NeedsScrolling()
			if got != tt.expected {
				t.Errorf("NeedsScrolling() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMarquee_GetDisplayText(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
		want  string
	}{
		{
			name:  "text fits - return as-is",
			text:  "hello",
			width: 10,
			want:  "hello",
		},
		{
			name:  "text overflows - truncate",
			text:  "hello world this is long",
			width: 10,
			want:  "hello worl", // Truncated to width
		},
		{
			name:  "empty text",
			text:  "",
			width: 10,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMarquee(tt.text, tt.width)
			defer m.Shutdown()

			got := m.GetDisplayText()
			if got != tt.want {
				t.Errorf("GetDisplayText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMarquee_SetText(t *testing.T) {
	m := NewMarquee("initial", 10)
	defer m.Shutdown()

	// Verify initial state
	if got := m.GetDisplayText(); got != "initial" {
		t.Errorf("Initial text = %q, want %q", got, "initial")
	}

	// Update text
	m.SetText("updated text")
	display := m.GetDisplayText()
	// Should either show full text or truncated version
	if !strings.Contains(display, "updated") {
		t.Errorf("After SetText, display should contain 'updated', got %q", display)
	}
}

func TestMarquee_SetWidth(t *testing.T) {
	m := NewMarquee("hello world", 5)
	defer m.Shutdown()

	// Initial width = 5
	if got := m.GetDisplayText(); len(stripColorCodes(got)) > 5 {
		t.Errorf("Initial display length should be <= 5, got %d", len(stripColorCodes(got)))
	}

	// Increase width
	m.SetWidth(11)
	if got := m.GetDisplayText(); got != "hello world" {
		t.Errorf("After SetWidth(11) = %q, want %q", got, "hello world")
	}

	// Decrease width
	m.SetWidth(3)
	if got := m.GetDisplayText(); len(stripColorCodes(got)) > 3 {
		t.Errorf("After SetWidth(3), display length should be <= 3, got %d", len(stripColorCodes(got)))
	}
}

func TestMarquee_StartStop(t *testing.T) {
	config := DefaultMarqueeConfig()
	config.Items = []string{"this is a long text that will animate"}
	config.Width = 10
	config.SlideInDuration = 200 * time.Millisecond
	config.CenterDuration = 300 * time.Millisecond
	config.SlideOutDuration = 200 * time.Millisecond

	m := NewMarqueeWithConfig(config)
	defer m.Shutdown()

	// Initially not animating
	m.mu.RLock()
	if m.isAnimating {
		t.Error("Marquee should not be animating initially")
	}
	m.mu.RUnlock()

	// Start animation
	m.Start()

	// Wait for animation to begin
	time.Sleep(100 * time.Millisecond)

	m.mu.RLock()
	isAnimating := m.isAnimating
	phase := m.phase
	m.mu.RUnlock()

	if !isAnimating {
		t.Error("Marquee should be animating after Start()")
	}

	if phase == PhaseIdle {
		t.Error("Marquee should not be idle after Start()")
	}

	// Stop animation
	m.Stop()

	m.mu.RLock()
	if m.isAnimating {
		t.Error("Marquee should not be animating after Stop()")
	}
	if m.phase != PhaseIdle {
		t.Errorf("Phase should be Idle after Stop(), got %v", m.phase)
	}
	m.mu.RUnlock()
}

func TestMarquee_NoScrollWhenTextFits(t *testing.T) {
	m := NewMarquee("short", 20)
	defer m.Shutdown()

	m.Start()

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	m.mu.RLock()
	isAnimating := m.isAnimating
	m.mu.RUnlock()

	if isAnimating {
		t.Error("Marquee should not animate when text fits within width")
	}
}

func TestMarquee_VisualLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "plain text",
			input:    "hello world",
			expected: 11,
		},
		{
			name:     "text with single color",
			input:    "[yellow]hello[white]",
			expected: 5,
		},
		{
			name:     "text with multiple colors",
			input:    "[yellow]hello[white] [red]world[white]",
			expected: 11,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "only color codes",
			input:    "[yellow][white]",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := visualLength(tt.input)
			if got != tt.expected {
				t.Errorf("visualLength(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMarquee_StripColorCodes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "single color",
			input:    "[yellow]hello[white]",
			expected: "hello",
		},
		{
			name:     "multiple colors",
			input:    "[yellow]hello[white] [red]world[white]",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripColorCodes(tt.input)
			if got != tt.expected {
				t.Errorf("stripColorCodes(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMarquee_AnimationPhases(t *testing.T) {
	config := DefaultMarqueeConfig()
	config.Items = []string{"this is a very long test item that will trigger the animation"}
	config.Width = 30 // Smaller width to trigger animation
	config.SlideInDuration = 100 * time.Millisecond
	config.CenterDuration = 200 * time.Millisecond
	config.SlideOutDuration = 100 * time.Millisecond

	m := NewMarqueeWithConfig(config)
	defer m.Shutdown()

	m.Start()

	// Phase 1: Slide-In
	time.Sleep(50 * time.Millisecond)
	m.mu.RLock()
	if m.phase != PhaseSlideIn {
		t.Errorf("Phase 1: Expected PhaseSlideIn, got %v", m.phase)
	}
	m.mu.RUnlock()

	// Phase 2: Center (after slide-in completes)
	time.Sleep(100 * time.Millisecond)
	m.mu.RLock()
	if m.phase != PhaseCenter {
		t.Errorf("Phase 2: Expected PhaseCenter, got %v", m.phase)
	}
	m.mu.RUnlock()

	// Phase 3: Slide-Out (after center duration)
	time.Sleep(250 * time.Millisecond)
	m.mu.RLock()
	phase := m.phase
	m.mu.RUnlock()

	// Should be in slide-out or back to slide-in (next item)
	if phase != PhaseSlideOut && phase != PhaseSlideIn {
		t.Errorf("Phase 3: Expected PhaseSlideOut or PhaseSlideIn, got %v", phase)
	}
}

func TestMarquee_BlinkDuringCenter(t *testing.T) {
	config := DefaultMarqueeConfig()
	config.Items = []string{"this is a blink test with long text"}
	config.Width = 20 // Smaller width to trigger animation
	config.SlideInDuration = 50 * time.Millisecond
	config.CenterDuration = 500 * time.Millisecond
	config.BlinkInterval = 100 * time.Millisecond

	m := NewMarqueeWithConfig(config)
	defer m.Shutdown()

	m.Start()

	// Wait for center phase
	time.Sleep(100 * time.Millisecond)

	// Check blink state changes
	m.mu.RLock()
	if m.phase != PhaseCenter {
		t.Errorf("Should be in center phase, got %v", m.phase)
	}
	blink1 := m.blinkVisible
	m.mu.RUnlock()

	// Wait for blink interval
	time.Sleep(120 * time.Millisecond)

	m.mu.RLock()
	blink2 := m.blinkVisible
	m.mu.RUnlock()

	// Blink state should have toggled
	if blink1 == blink2 {
		t.Error("Blink state should toggle during center phase")
	}
}

func TestMarquee_MultipleItems(t *testing.T) {
	// Test with action bar format
	actionBarText := "[yellow][[white]Enter[yellow] [white]Open Ticket[yellow]] [yellow][[white]Space[yellow] [white]Select[yellow]]"

	config := DefaultMarqueeConfig()
	config.Items = []string{actionBarText}
	config.Width = 40

	m := NewMarqueeWithConfig(config)
	defer m.Shutdown()

	// Should parse into multiple items
	if len(m.items) < 2 {
		t.Errorf("Expected multiple items from action bar format, got %d", len(m.items))
	}

	// Should need scrolling (multiple items always marquee)
	if !m.NeedsScrolling() {
		t.Error("Multiple items should always need marquee animation")
	}
}

func TestMarquee_CheckResize(t *testing.T) {
	longText := "this is a very long text that will need animation"

	config := DefaultMarqueeConfig()
	config.Items = []string{longText}
	config.Width = 20

	m := NewMarqueeWithConfig(config)
	defer m.Shutdown()

	// Should need animation at 20 cols
	if !m.NeedsScrolling() {
		t.Error("Text should need animation at 20 cols")
	}

	m.Start()
	time.Sleep(50 * time.Millisecond)

	m.mu.RLock()
	wasAnimating := m.isAnimating
	m.mu.RUnlock()

	if !wasAnimating {
		t.Error("Should be animating")
	}

	// Resize to fit text
	m.CheckResize(100)
	time.Sleep(50 * time.Millisecond)

	m.mu.RLock()
	isAnimating := m.isAnimating
	m.mu.RUnlock()

	if isAnimating {
		t.Error("Should stop animating when text fits")
	}

	// Resize back to narrow
	m.CheckResize(20)
	time.Sleep(100 * time.Millisecond)

	m.mu.RLock()
	isAnimating = m.isAnimating
	m.mu.RUnlock()

	if !isAnimating {
		t.Error("Should resume animating when text overflows")
	}
}

// BenchmarkMarquee_CPU measures CPU usage of marquee animation.
func BenchmarkMarquee_CPU(b *testing.B) {
	config := DefaultMarqueeConfig()
	config.Items = []string{"This is a long text that will continuously animate for benchmarking"}
	config.Width = 20

	m := NewMarqueeWithConfig(config)
	defer m.Shutdown()

	m.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.GetDisplayText()
	}
}

func TestMarquee_Italicize(t *testing.T) {
	text := "[yellow]Test[white]"
	result := italicize(text)

	// Should have angle brackets
	if !strings.Contains(result, "⟨") || !strings.Contains(result, "⟩") {
		t.Errorf("Italicized text should have angle brackets, got %q", result)
	}

	// Should contain original text
	if !strings.Contains(result, "Test") {
		t.Errorf("Italicized text should contain original text, got %q", result)
	}
}
