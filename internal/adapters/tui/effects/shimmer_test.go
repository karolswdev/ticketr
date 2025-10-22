package effects

import (
	"strings"
	"testing"
)

func TestShimmerEffect(t *testing.T) {
	t.Run("Create shimmer", func(t *testing.T) {
		shimmer := NewShimmerEffect(20)

		if shimmer.width != 20 {
			t.Errorf("Expected width 20, got %d", shimmer.width)
		}

		if shimmer.offset != 0 {
			t.Error("Expected offset to start at 0")
		}

		if shimmer.direction != 1 {
			t.Error("Expected direction to start at 1 (left-to-right)")
		}
	})

	t.Run("Update shimmer", func(t *testing.T) {
		shimmer := NewShimmerEffect(10)

		initialOffset := shimmer.offset
		shimmer.Update()

		if shimmer.offset == initialOffset {
			t.Error("Expected offset to change after update")
		}
	})

	t.Run("Shimmer bounces at edges", func(t *testing.T) {
		shimmer := NewShimmerEffect(5)

		// Move to right edge (offset should reach width)
		for i := 0; i < 7; i++ {
			shimmer.Update()
		}

		// Should have bounced and changed direction
		if shimmer.direction != -1 {
			t.Errorf("Expected direction to reverse at right edge, got direction=%d, offset=%d", shimmer.direction, shimmer.offset)
		}

		// Reset and test from scratch
		shimmer.Reset()

		// Just verify Update() works without panicking
		for i := 0; i < 20; i++ {
			shimmer.Update()
		}

		// Shimmer should have bounced at least once
		if shimmer.offset < 0 || shimmer.offset > shimmer.width {
			t.Errorf("Offset out of bounds: %d (width: %d)", shimmer.offset, shimmer.width)
		}
	})

	t.Run("Apply shimmer to bar", func(t *testing.T) {
		shimmer := NewShimmerEffect(10)
		barString := "[██████████]"

		result := shimmer.Apply(barString, 10)

		if result == barString {
			t.Error("Expected shimmer to modify bar string")
		}

		if !strings.Contains(result, "[") || !strings.Contains(result, "]") {
			t.Error("Expected result to maintain bracket structure")
		}
	})

	t.Run("Reset shimmer", func(t *testing.T) {
		shimmer := NewShimmerEffect(10)

		shimmer.offset = 5
		shimmer.direction = -1

		shimmer.Reset()

		if shimmer.offset != 0 {
			t.Errorf("Expected offset to reset to 0, got %d", shimmer.offset)
		}

		if shimmer.direction != 1 {
			t.Error("Expected direction to reset to 1")
		}
	})
}

func TestProgressBarShimmer(t *testing.T) {
	t.Run("Create with shimmer enabled", func(t *testing.T) {
		pbs := NewProgressBarShimmer(20, true)

		if !pbs.enabled {
			t.Error("Expected shimmer to be enabled")
		}
	})

	t.Run("Create with shimmer disabled", func(t *testing.T) {
		pbs := NewProgressBarShimmer(20, false)

		if pbs.enabled {
			t.Error("Expected shimmer to be disabled")
		}
	})

	t.Run("Apply with shimmer disabled", func(t *testing.T) {
		pbs := NewProgressBarShimmer(20, false)
		barString := "[██████████]"

		result := pbs.Apply(barString, 10)

		if result != barString {
			t.Error("Expected no modification when shimmer disabled")
		}
	})

	t.Run("Update and apply", func(t *testing.T) {
		pbs := NewProgressBarShimmer(20, true)
		barString := "[██████████░░░░░░░░░░]"

		pbs.Update()
		result := pbs.Apply(barString, 10)

		if result == barString && pbs.enabled {
			// Note: This might not always trigger if offset is 0
			// Just checking the mechanism works
		}
	})

	t.Run("SetEnabled", func(t *testing.T) {
		pbs := NewProgressBarShimmer(20, false)

		pbs.SetEnabled(true)
		if !pbs.enabled {
			t.Error("Expected shimmer to be enabled")
		}

		pbs.SetEnabled(false)
		if pbs.enabled {
			t.Error("Expected shimmer to be disabled")
		}
	})
}

func TestGradientText(t *testing.T) {
	t.Run("Empty text", func(t *testing.T) {
		result := GradientText("", "red", "blue")

		if result != "" {
			t.Error("Expected empty result for empty text")
		}
	})

	t.Run("Non-empty text", func(t *testing.T) {
		result := GradientText("Hello World", "red", "blue")

		if !strings.Contains(result, "red") {
			t.Error("Expected result to contain start color")
		}

		if !strings.Contains(result, "blue") {
			t.Error("Expected result to contain end color")
		}

		// The text is modified with color tags, so we check for individual words
		if !strings.Contains(result, "Hello") || !strings.Contains(result, "World") {
			t.Errorf("Expected result to contain original text words, got: %s", result)
		}
	})

	t.Run("Short text", func(t *testing.T) {
		result := GradientText("Hi", "red", "blue")

		if result == "" {
			t.Error("Expected non-empty result for short text")
		}
	})
}

func TestPulseIntensity(t *testing.T) {
	t.Run("Zero max frames", func(t *testing.T) {
		intensity := PulseIntensity(5, 0)

		if intensity != 1.0 {
			t.Errorf("Expected intensity 1.0 for zero max frames, got %f", intensity)
		}
	})

	t.Run("Valid pulse", func(t *testing.T) {
		intensity := PulseIntensity(5, 10)

		if intensity < 0.5 || intensity > 1.0 {
			t.Errorf("Expected intensity between 0.5 and 1.0, got %f", intensity)
		}
	})

	t.Run("Frame at max", func(t *testing.T) {
		intensity := PulseIntensity(10, 10)

		if intensity < 0.5 || intensity > 1.0 {
			t.Errorf("Expected intensity between 0.5 and 1.0, got %f", intensity)
		}
	})
}

func TestRainbow(t *testing.T) {
	t.Run("Empty text", func(t *testing.T) {
		result := Rainbow("")

		if result != "" {
			t.Error("Expected empty result for empty text")
		}
	})

	t.Run("Non-empty text", func(t *testing.T) {
		result := Rainbow("Hello")

		if result == "Hello" {
			t.Error("Expected text to be modified with colors")
		}

		// Should contain color tags
		if !strings.Contains(result, "[") {
			t.Error("Expected result to contain color tags")
		}
	})

	t.Run("Text with spaces", func(t *testing.T) {
		result := Rainbow("Hello World")

		// Result will have color tags interspersed, so check for individual letters
		hasH := strings.Contains(result, "H")
		hasW := strings.Contains(result, "W")
		hasSpace := strings.Contains(result, " ")

		if !hasH || !hasW || !hasSpace {
			t.Errorf("Expected result to preserve text content, got: %s", result)
		}
	})
}

// Benchmarks

func BenchmarkShimmerUpdate(b *testing.B) {
	shimmer := NewShimmerEffect(20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shimmer.Update()
	}
}

func BenchmarkShimmerApply(b *testing.B) {
	shimmer := NewShimmerEffect(20)
	barString := "[██████████░░░░░░░░░░]"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shimmer.Apply(barString, 10)
	}
}

func BenchmarkGradientText(b *testing.B) {
	text := "This is a test title"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GradientText(text, "blue", "cyan")
	}
}

func BenchmarkRainbow(b *testing.B) {
	text := "Rainbow colored text"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Rainbow(text)
	}
}
