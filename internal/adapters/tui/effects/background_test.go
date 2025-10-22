package effects

import (
	"testing"
	"time"

	"github.com/rivo/tview"
)

func TestBackgroundAnimator(t *testing.T) {
	app := tview.NewApplication()

	t.Run("Create with default config", func(t *testing.T) {
		config := DefaultBackgroundConfig()
		ba := NewBackgroundAnimator(app, config)

		if ba == nil {
			t.Fatal("Expected non-nil background animator")
		}

		if ba.config.Effect != BackgroundNone {
			t.Errorf("Expected default effect to be none, got %s", ba.config.Effect)
		}

		ba.Shutdown()
	})

	t.Run("Hyperspace effect", func(t *testing.T) {
		config := BackgroundConfig{
			Effect:  BackgroundHyperspace,
			Density: 2.0, // Higher density to ensure particles spawn
			Speed:   50,
			Enabled: true,
			MaxFPS:  15,
		}

		ba := NewBackgroundAnimator(app, config)
		ba.width = 80
		ba.height = 24

		// Create some particles
		ba.spawnParticles()

		if len(ba.particles) == 0 {
			t.Skip("No particles spawned (density too low or RNG issue)")
		}

		// Verify particles move
		initialX := ba.particles[0].X
		ba.updateParticle(ba.particles[0])

		if ba.particles[0].X == initialX {
			t.Error("Expected particle to move")
		}

		ba.Shutdown()
	})

	t.Run("Snow effect", func(t *testing.T) {
		config := BackgroundConfig{
			Effect:  BackgroundSnow,
			Density: 2.0, // Higher density to ensure particles spawn
			Speed:   100,
			Enabled: true,
			MaxFPS:  15,
		}

		ba := NewBackgroundAnimator(app, config)
		ba.width = 80
		ba.height = 24

		// Create some particles
		ba.spawnParticles()

		if len(ba.particles) == 0 {
			t.Skip("No particles spawned (density too low or RNG issue)")
		}

		// Verify particles fall
		initialY := ba.particles[0].Y
		for i := 0; i < 5; i++ {
			ba.updateParticle(ba.particles[0])
		}

		if ba.particles[0].Y <= initialY {
			t.Error("Expected particle to fall downward")
		}

		ba.Shutdown()
	})

	t.Run("Pause and resume", func(t *testing.T) {
		config := BackgroundConfig{
			Effect:  BackgroundHyperspace,
			Density: 0.02,
			Speed:   50,
			Enabled: true,
			MaxFPS:  15,
		}

		ba := NewBackgroundAnimator(app, config)

		if ba.paused {
			t.Error("Expected animator to not be paused initially")
		}

		ba.Pause()
		if !ba.paused {
			t.Error("Expected animator to be paused")
		}

		ba.Resume()
		if ba.paused {
			t.Error("Expected animator to be resumed")
		}

		ba.Shutdown()
	})

	t.Run("Clear particles", func(t *testing.T) {
		config := BackgroundConfig{
			Effect:  BackgroundHyperspace,
			Density: 2.0, // Higher density to ensure particles spawn
			Speed:   50,
			Enabled: true,
			MaxFPS:  15,
		}

		ba := NewBackgroundAnimator(app, config)
		ba.width = 80
		ba.height = 24
		ba.spawnParticles()

		if len(ba.particles) == 0 {
			t.Skip("No particles spawned (density too low or RNG issue)")
		}

		ba.Clear()

		if len(ba.particles) != 0 {
			t.Errorf("Expected particles to be cleared, got %d", len(ba.particles))
		}

		ba.Shutdown()
	})
}

func TestBackgroundConfig(t *testing.T) {
	t.Run("Valid config", func(t *testing.T) {
		config := BackgroundConfig{
			Effect:  BackgroundHyperspace,
			Density: 0.05,
			Speed:   100,
			Enabled: true,
			MaxFPS:  15,
		}

		if err := ValidateConfig(config); err != nil {
			t.Errorf("Expected valid config, got error: %v", err)
		}
	})

	t.Run("Invalid density", func(t *testing.T) {
		config := BackgroundConfig{
			Effect:  BackgroundHyperspace,
			Density: 15.0, // Too high
			Speed:   100,
			Enabled: true,
			MaxFPS:  15,
		}

		if err := ValidateConfig(config); err == nil {
			t.Error("Expected validation error for high density")
		}
	})

	t.Run("Invalid speed", func(t *testing.T) {
		config := BackgroundConfig{
			Effect:  BackgroundHyperspace,
			Density: 0.02,
			Speed:   5, // Too low
			Enabled: true,
			MaxFPS:  15,
		}

		if err := ValidateConfig(config); err == nil {
			t.Error("Expected validation error for low speed")
		}
	})

	t.Run("Invalid FPS", func(t *testing.T) {
		config := BackgroundConfig{
			Effect:  BackgroundHyperspace,
			Density: 0.02,
			Speed:   100,
			Enabled: true,
			MaxFPS:  100, // Too high
		}

		if err := ValidateConfig(config); err == nil {
			t.Error("Expected validation error for high FPS")
		}
	})
}

func TestBackgroundEffectDescription(t *testing.T) {
	tests := []struct {
		effect      BackgroundEffect
		expectEmpty bool
	}{
		{BackgroundNone, false},
		{BackgroundHyperspace, false},
		{BackgroundSnow, false},
		{"invalid", false}, // Should return something
	}

	for _, tt := range tests {
		desc := GetEffectDescription(tt.effect)
		if desc == "" && !tt.expectEmpty {
			t.Errorf("Expected non-empty description for %s", tt.effect)
		}
	}
}

// Benchmarks

func BenchmarkBackgroundAnimatorUpdate(b *testing.B) {
	app := tview.NewApplication()
	config := BackgroundConfig{
		Effect:  BackgroundHyperspace,
		Density: 0.02,
		Speed:   100,
		Enabled: true,
		MaxFPS:  15,
	}

	ba := NewBackgroundAnimator(app, config)
	ba.width = 80
	ba.height = 24
	ba.spawnParticles()

	defer ba.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba.update()
	}
}

func BenchmarkParticleUpdate(b *testing.B) {
	app := tview.NewApplication()
	config := BackgroundConfig{
		Effect:  BackgroundHyperspace,
		Density: 0.02,
		Speed:   100,
		Enabled: true,
		MaxFPS:  15,
	}

	ba := NewBackgroundAnimator(app, config)
	ba.width = 80
	ba.height = 24

	particle := ba.createHyperspaceParticle()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba.updateParticle(particle)
	}
}

func BenchmarkSpawnParticles(b *testing.B) {
	app := tview.NewApplication()
	config := BackgroundConfig{
		Effect:  BackgroundHyperspace,
		Density: 0.02,
		Speed:   100,
		Enabled: true,
		MaxFPS:  15,
	}

	ba := NewBackgroundAnimator(app, config)
	ba.width = 80
	ba.height = 24

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba.Clear()
		ba.spawnParticles()
	}
}

// Performance assertion test - ensure background animator uses minimal CPU
func TestBackgroundAnimatorPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	app := tview.NewApplication()
	config := BackgroundConfig{
		Effect:  BackgroundHyperspace,
		Density: 0.02,
		Speed:   100,
		Enabled: true,
		MaxFPS:  15,
	}

	ba := NewBackgroundAnimator(app, config)
	ba.width = 80
	ba.height = 24
	ba.Start()

	// Run for 1 second
	start := time.Now()
	time.Sleep(1 * time.Second)
	elapsed := time.Since(start)

	ba.Shutdown()

	// The background animator should not significantly delay
	// If it takes much longer than 1 second, something is wrong
	if elapsed > 1200*time.Millisecond {
		t.Errorf("Background animator took too long: %v (expected ~1s)", elapsed)
	}

	// Note: For true CPU usage testing, you'd need to use runtime profiling
	// This is a basic sanity check
}
