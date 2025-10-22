package effects

import (
	"testing"
	"time"

	"github.com/rivo/tview"
)

func TestAnimator(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping animator tests in short mode (they use timers)")
	}

	app := tview.NewApplication()
	animator := NewAnimator(app)

	t.Run("Create animator", func(t *testing.T) {
		if animator == nil {
			t.Fatal("Expected non-nil animator")
		}

		if !animator.IsEnabled() {
			t.Error("Expected animator to be enabled by default")
		}
	})

	t.Run("Global motion kill switch", func(t *testing.T) {
		animator.SetEnabled(false)
		if animator.IsEnabled() {
			t.Error("Expected animator to be disabled")
		}

		animator.SetEnabled(true)
		if !animator.IsEnabled() {
			t.Error("Expected animator to be enabled")
		}
	})

	t.Run("Stop nonexistent animation", func(t *testing.T) {
		// Should not panic
		animator.Stop("nonexistent")
	})

	animator.Shutdown()
}

func TestSpinner(t *testing.T) {
	frame := 0
	handler := Spinner(&frame)

	for i := 0; i < 20; i++ {
		handler()
		if frame != (i+1)%len(SpinnerFrames) {
			t.Errorf("Expected frame %d, got %d", (i+1)%len(SpinnerFrames), frame)
		}
	}
}

func TestSparkle(t *testing.T) {
	sparkle := NewSparkle(10, 5)

	if sparkle.X != 10 || sparkle.Y != 5 {
		t.Errorf("Expected position (10, 5), got (%d, %d)", sparkle.X, sparkle.Y)
	}

	if sparkle.TTL != sparkle.MaxTTL {
		t.Errorf("Expected TTL to equal MaxTTL initially")
	}

	// Update until death
	alive := true
	updates := 0
	for alive {
		alive = sparkle.Update()
		updates++
		if updates > 20 {
			t.Fatal("Sparkle should have died by now")
		}
	}

	if sparkle.TTL != 0 {
		t.Errorf("Expected TTL to be 0 when dead, got %d", sparkle.TTL)
	}
}

func TestToggleAnimation(t *testing.T) {
	completed := false
	toggle := NewToggleAnimation(func() {
		completed = true
	})

	expectedFrames := []string{"[ ]", "[•]", "[x]"}

	for i := 0; i < len(expectedFrames); i++ {
		frame := toggle.GetFrame()
		if frame != expectedFrames[i] {
			t.Errorf("Expected frame %q, got %q", expectedFrames[i], frame)
		}
		toggle.Update()
	}

	if !completed {
		t.Error("Expected completion callback to be called")
	}
}

func TestFadeAnimation(t *testing.T) {
	fade := NewFadeAnimation(100 * time.Millisecond)

	// Phase 1: Light
	if fade.GetOpacity() != "░" {
		t.Errorf("Expected initial opacity ░, got %s", fade.GetOpacity())
	}

	// Advance to middle
	fade.Update(40 * time.Millisecond)
	if fade.GetOpacity() != "▒" {
		t.Errorf("Expected middle opacity ▒, got %s", fade.GetOpacity())
	}

	// Advance to heavy
	fade.Update(30 * time.Millisecond)
	if fade.GetOpacity() != "█" {
		t.Errorf("Expected heavy opacity █, got %s", fade.GetOpacity())
	}

	// Advance to done
	alive := fade.Update(40 * time.Millisecond)
	if alive {
		t.Error("Expected fade to be done")
	}

	if !fade.IsDone() {
		t.Error("Expected IsDone to return true")
	}
}

func TestPulse(t *testing.T) {
	brightness := 0
	maxBrightness := 10
	handler := Pulse(&brightness, maxBrightness)

	// Should increase to max
	for i := 0; i < 15; i++ {
		handler()
	}

	if brightness > maxBrightness {
		t.Errorf("Brightness should not exceed max, got %d", brightness)
	}

	// Should decrease back
	for i := 0; i < 25; i++ {
		handler()
	}

	if brightness < 0 {
		t.Errorf("Brightness should not go below 0, got %d", brightness)
	}
}

// Benchmarks

func BenchmarkAnimator(b *testing.B) {
	app := tview.NewApplication()
	animator := NewAnimator(app)
	defer animator.Shutdown()

	handler := func() bool { return true }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		animator.Start("bench", 100*time.Millisecond, handler)
		animator.Stop("bench")
	}
}

func BenchmarkSpinner(b *testing.B) {
	frame := 0
	handler := Spinner(&frame)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler()
	}
}

func BenchmarkSparkleUpdate(b *testing.B) {
	sparkle := NewSparkle(10, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sparkle.Update()
		if !sparkle.Update() {
			sparkle = NewSparkle(10, 5)
		}
	}
}

func BenchmarkFadeAnimation(b *testing.B) {
	fade := NewFadeAnimation(100 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fade.Update(10 * time.Millisecond)
		if fade.IsDone() {
			fade = NewFadeAnimation(100 * time.Millisecond)
		}
	}
}
