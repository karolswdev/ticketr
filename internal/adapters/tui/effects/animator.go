package effects

import (
	"context"
	"sync"
	"time"

	"github.com/rivo/tview"
)

// Animator is the core animation engine that manages all visual effects.
// It provides a centralized, CPU-efficient way to run animations with proper cancellation.
type Animator struct {
	app        *tview.Application
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	animations map[string]*Animation
	mu         sync.RWMutex
	enabled    bool // Global motion kill switch
}

// Animation represents a single running animation.
type Animation struct {
	Name     string
	Interval time.Duration
	Handler  func() bool // Returns false to stop animation
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewAnimator creates a new animation engine.
func NewAnimator(app *tview.Application) *Animator {
	ctx, cancel := context.WithCancel(context.Background())
	return &Animator{
		app:        app,
		ctx:        ctx,
		cancel:     cancel,
		animations: make(map[string]*Animation),
		enabled:    true, // Default: animations enabled (but individual effects default to OFF)
	}
}

// SetEnabled sets the global motion kill switch.
// When disabled, all animations are paused.
func (a *Animator) SetEnabled(enabled bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.enabled = enabled
}

// IsEnabled returns whether animations are globally enabled.
func (a *Animator) IsEnabled() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.enabled
}

// Start begins an animation with the given name and interval.
// The handler function is called periodically and should return false to stop.
// If an animation with the same name exists, it is replaced.
func (a *Animator) Start(name string, interval time.Duration, handler func() bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Cancel existing animation with same name
	if existing, ok := a.animations[name]; ok {
		existing.cancel()
	}

	// Create animation context
	ctx, cancel := context.WithCancel(a.ctx)
	animation := &Animation{
		Name:     name,
		Interval: interval,
		Handler:  handler,
		ctx:      ctx,
		cancel:   cancel,
	}

	a.animations[name] = animation
	a.wg.Add(1)

	// Start animation goroutine
	go a.runAnimation(animation)
}

// runAnimation executes an animation loop with proper timing and cancellation.
func (a *Animator) runAnimation(anim *Animation) {
	defer a.wg.Done()
	defer func() {
		a.mu.Lock()
		delete(a.animations, anim.Name)
		a.mu.Unlock()
	}()

	ticker := time.NewTicker(anim.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-anim.ctx.Done():
			return
		case <-ticker.C:
			// Skip if globally disabled
			if !a.IsEnabled() {
				continue
			}

			// Run handler and check if it wants to continue
			shouldContinue := true
			a.app.QueueUpdateDraw(func() {
				shouldContinue = anim.Handler()
			})

			if !shouldContinue {
				return
			}
		}
	}
}

// Stop stops a specific animation by name.
func (a *Animator) Stop(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if anim, ok := a.animations[name]; ok {
		anim.cancel()
		delete(a.animations, name)
	}
}

// StopAll stops all running animations.
func (a *Animator) StopAll() {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, anim := range a.animations {
		anim.cancel()
	}
	a.animations = make(map[string]*Animation)
}

// Shutdown gracefully stops all animations and waits for them to complete.
func (a *Animator) Shutdown() {
	a.cancel()
	a.wg.Wait()
}

// Pulse creates a pulsing animation effect for borders.
// Returns a function that cycles through brightness levels.
func Pulse(currentBrightness *int, maxBrightness int) func() bool {
	direction := 1
	return func() bool {
		*currentBrightness += direction
		if *currentBrightness >= maxBrightness {
			direction = -1
		} else if *currentBrightness <= 0 {
			direction = 1
		}
		return true
	}
}

// Spinner returns a function that cycles through spinner frames.
// Frames: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏ (Braille spinner)
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func Spinner(currentFrame *int) func() bool {
	return func() bool {
		*currentFrame = (*currentFrame + 1) % len(SpinnerFrames)
		return true
	}
}

// Sparkle represents a success sparkle particle.
type Sparkle struct {
	X, Y       int
	Char       string
	TTL        int // Time to live in frames
	MaxTTL     int
	Brightness float64
}

// SparkleFrames are the characters used for sparkle animation.
var SparkleFrames = []string{"✦", "✧", "⋆", "∗", "·"}

// NewSparkle creates a sparkle at the given position.
func NewSparkle(x, y int) *Sparkle {
	maxTTL := 10 // 10 frames at 50ms = 500ms total
	return &Sparkle{
		X:          x,
		Y:          y,
		Char:       SparkleFrames[0],
		TTL:        maxTTL,
		MaxTTL:     maxTTL,
		Brightness: 1.0,
	}
}

// Update updates the sparkle state for the next frame.
// Returns false if the sparkle is dead.
func (s *Sparkle) Update() bool {
	s.TTL--
	if s.TTL <= 0 {
		return false
	}

	// Update character based on remaining TTL
	frameIndex := (s.MaxTTL - s.TTL) * len(SparkleFrames) / s.MaxTTL
	if frameIndex >= len(SparkleFrames) {
		frameIndex = len(SparkleFrames) - 1
	}
	s.Char = SparkleFrames[frameIndex]

	// Fade out
	s.Brightness = float64(s.TTL) / float64(s.MaxTTL)

	return true
}

// ToggleAnimation creates a 3-frame checkbox toggle animation.
// Frames: [ ] → [•] → [x]
type ToggleAnimation struct {
	frame      int
	frames     []string
	onComplete func()
}

// NewToggleAnimation creates a new toggle animation.
func NewToggleAnimation(onComplete func()) *ToggleAnimation {
	return &ToggleAnimation{
		frame:      0,
		frames:     []string{"[ ]", "[•]", "[x]"},
		onComplete: onComplete,
	}
}

// Update advances the toggle animation.
// Returns false when animation is complete.
func (t *ToggleAnimation) Update() bool {
	t.frame++
	if t.frame >= len(t.frames) {
		if t.onComplete != nil {
			t.onComplete()
		}
		return false
	}
	return true
}

// GetFrame returns the current frame string.
func (t *ToggleAnimation) GetFrame() string {
	if t.frame >= len(t.frames) {
		return t.frames[len(t.frames)-1]
	}
	return t.frames[t.frame]
}

// FadePhase represents a phase in modal fade-in animation.
type FadePhase int

const (
	FadePhaseLight  FadePhase = 0 // ░
	FadePhaseMiddle FadePhase = 1 // ▒
	FadePhaseHeavy  FadePhase = 2 // █
	FadePhaseDone   FadePhase = 3
)

// FadeAnimation manages modal fade-in animation.
type FadeAnimation struct {
	phase    FadePhase
	duration time.Duration
	elapsed  time.Duration
}

// NewFadeAnimation creates a new fade animation with total duration.
func NewFadeAnimation(duration time.Duration) *FadeAnimation {
	return &FadeAnimation{
		phase:    FadePhaseLight,
		duration: duration,
		elapsed:  0,
	}
}

// Update advances the fade animation by the given delta time.
// Returns false when animation is complete.
func (f *FadeAnimation) Update(delta time.Duration) bool {
	f.elapsed += delta

	// Calculate which phase we should be in
	progress := float64(f.elapsed) / float64(f.duration)
	if progress >= 1.0 {
		f.phase = FadePhaseDone
		return false
	}

	// Three phases: 0-0.33, 0.33-0.66, 0.66-1.0
	if progress < 0.33 {
		f.phase = FadePhaseLight
	} else if progress < 0.66 {
		f.phase = FadePhaseMiddle
	} else {
		f.phase = FadePhaseHeavy
	}

	return true
}

// GetOpacity returns the current opacity character.
func (f *FadeAnimation) GetOpacity() string {
	switch f.phase {
	case FadePhaseLight:
		return "░"
	case FadePhaseMiddle:
		return "▒"
	case FadePhaseHeavy:
		return "█"
	default:
		return ""
	}
}

// IsDone returns true if the fade animation is complete.
func (f *FadeAnimation) IsDone() bool {
	return f.phase == FadePhaseDone
}
