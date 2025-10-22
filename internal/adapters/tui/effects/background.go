package effects

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// BackgroundEffect represents the type of background ambient effect.
type BackgroundEffect string

const (
	BackgroundNone       BackgroundEffect = "none"
	BackgroundHyperspace BackgroundEffect = "hyperspace"
	BackgroundSnow       BackgroundEffect = "snow"
)

// BackgroundConfig configures the background animator.
type BackgroundConfig struct {
	Effect    BackgroundEffect
	Density   float64 // Particles per 100 cells (0.01 = 1%, 0.02 = 2%)
	Speed     int     // Milliseconds per frame (higher = slower)
	Enabled   bool    // Master enable switch
	MaxFPS    int     // Frame rate limit (default: 15)
	AutoPause bool    // Auto-pause when UI is busy (default: true)
}

// DefaultBackgroundConfig returns a conservative default configuration.
func DefaultBackgroundConfig() BackgroundConfig {
	return BackgroundConfig{
		Effect:    BackgroundNone,
		Density:   0.02, // 2%
		Speed:     100,  // 100ms = 10 FPS
		Enabled:   false,
		MaxFPS:    15,
		AutoPause: true,
	}
}

// Particle represents a single background particle.
type Particle struct {
	X, Y   int
	Char   rune
	Color  tcell.Color
	VX, VY int // Velocity
	Age    int // Frames alive
	MaxAge int // Lifespan in frames
	Trail  []rune
	Fading bool
}

// BackgroundAnimator manages ambient background effects.
type BackgroundAnimator struct {
	config    BackgroundConfig
	particles []*Particle
	width     int
	height    int
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	app       *tview.Application
	overlay   *BackgroundOverlay
	paused    bool
	rng       *rand.Rand
}

// BackgroundOverlay is a custom tview primitive that renders background effects.
type BackgroundOverlay struct {
	*tview.Box
	animator *BackgroundAnimator
}

// NewBackgroundAnimator creates a new background animator.
func NewBackgroundAnimator(app *tview.Application, config BackgroundConfig) *BackgroundAnimator {
	ctx, cancel := context.WithCancel(context.Background())

	ba := &BackgroundAnimator{
		config:    config,
		particles: make([]*Particle, 0),
		width:     80, // Default, will be updated
		height:    24, // Default, will be updated
		ctx:       ctx,
		cancel:    cancel,
		app:       app,
		paused:    false,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Create overlay primitive
	box := tview.NewBox()
	box.SetBackgroundTransparent(true) // CRITICAL: Must be transparent to see particles
	ba.overlay = &BackgroundOverlay{
		Box:      box,
		animator: ba,
	}

	return ba
}

// SetConfig updates the background configuration.
func (ba *BackgroundAnimator) SetConfig(config BackgroundConfig) {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.config = config
}

// GetOverlay returns the background overlay primitive.
func (ba *BackgroundAnimator) GetOverlay() tview.Primitive {
	return ba.overlay
}

// Start begins the background animation loop.
func (ba *BackgroundAnimator) Start() {
	fmt.Fprintf(os.Stderr, "[DEBUG BackgroundAnimator] Start called, enabled: %v\n", ba.config.Enabled)
	if !ba.config.Enabled {
		fmt.Fprintf(os.Stderr, "[DEBUG BackgroundAnimator] Not starting because config.Enabled is false\n")
		return
	}

	fmt.Fprintf(os.Stderr, "[DEBUG BackgroundAnimator] Starting animation loop goroutine\n")
	ba.wg.Add(1)
	go ba.run()
}

// run is the main animation loop.
func (ba *BackgroundAnimator) run() {
	defer ba.wg.Done()

	// Calculate frame interval based on max FPS
	frameInterval := time.Duration(ba.config.Speed) * time.Millisecond
	if ba.config.MaxFPS > 0 {
		minInterval := time.Second / time.Duration(ba.config.MaxFPS)
		if frameInterval < minInterval {
			frameInterval = minInterval
		}
	}

	fmt.Fprintf(os.Stderr, "[DEBUG BackgroundAnimator] Animation loop started, frame interval: %v\n", frameInterval)

	ticker := time.NewTicker(frameInterval)
	defer ticker.Stop()

	frameCount := 0
	for {
		select {
		case <-ba.ctx.Done():
			fmt.Fprintf(os.Stderr, "[DEBUG BackgroundAnimator] Animation loop stopped (context cancelled)\n")
			return
		case <-ticker.C:
			if ba.paused || !ba.config.Enabled {
				continue
			}

			ba.update()
			ba.app.QueueUpdateDraw(func() {
				// Draw is handled by the overlay's Draw method
			})

			frameCount++
			if frameCount%30 == 0 { // Log every 30 frames
				ba.mu.RLock()
				particleCount := len(ba.particles)
				ba.mu.RUnlock()
				fmt.Fprintf(os.Stderr, "[DEBUG BackgroundAnimator] Frame %d, particles: %d\n", frameCount, particleCount)
			}
		}
	}
}

// update updates all particles and spawns new ones.
func (ba *BackgroundAnimator) update() {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	// Update terminal dimensions
	ba.updateDimensions()

	// Update existing particles
	alive := make([]*Particle, 0, len(ba.particles))
	for _, p := range ba.particles {
		if ba.updateParticle(p) {
			alive = append(alive, p)
		}
	}
	ba.particles = alive

	// Spawn new particles based on density
	ba.spawnParticles()
}

// updateDimensions updates the width and height based on overlay size.
func (ba *BackgroundAnimator) updateDimensions() {
	if ba.overlay != nil {
		_, _, w, h := ba.overlay.GetInnerRect()
		if w > 0 && h > 0 {
			ba.width = w
			ba.height = h
		}
	}
}

// updateParticle updates a single particle.
// Returns false if particle should be removed.
func (ba *BackgroundAnimator) updateParticle(p *Particle) bool {
	p.Age++

	// Check lifespan
	if p.MaxAge > 0 && p.Age >= p.MaxAge {
		return false
	}

	// Update position
	p.X += p.VX
	p.Y += p.VY

	// Effect-specific behavior
	switch ba.config.Effect {
	case BackgroundHyperspace:
		// Accelerate over time
		if p.VX > 0 {
			p.VX++
		} else if p.VX < 0 {
			p.VX--
		}

		// Wrap around or remove if out of bounds
		if p.X < 0 || p.X >= ba.width || p.Y < 0 || p.Y >= ba.height {
			return false
		}

	case BackgroundSnow:
		// Gentle drift
		if p.Age%5 == 0 {
			p.X += ba.rng.Intn(3) - 1 // Random drift -1, 0, or 1
		}

		// Remove if out of bounds
		if p.Y >= ba.height || p.X < 0 || p.X >= ba.width {
			return false
		}
	}

	return true
}

// spawnParticles creates new particles based on density configuration.
func (ba *BackgroundAnimator) spawnParticles() {
	if ba.config.Effect == BackgroundNone {
		return
	}

	totalCells := ba.width * ba.height
	targetParticles := int(float64(totalCells) * ba.config.Density / 100.0)

	// Spawn particles to maintain density
	for len(ba.particles) < targetParticles {
		p := ba.createParticle()
		if p != nil {
			ba.particles = append(ba.particles, p)
		} else {
			break // Can't create more particles
		}
	}
}

// createParticle creates a new particle based on current effect.
func (ba *BackgroundAnimator) createParticle() *Particle {
	switch ba.config.Effect {
	case BackgroundHyperspace:
		return ba.createHyperspaceParticle()
	case BackgroundSnow:
		return ba.createSnowParticle()
	default:
		return nil
	}
}

// createHyperspaceParticle creates a hyperspace star.
func (ba *BackgroundAnimator) createHyperspaceParticle() *Particle {
	chars := []rune{'.', '*', '·', '∗', '⋆'}
	colors := []tcell.Color{
		tcell.ColorWhite,
		tcell.ColorLightBlue,
		tcell.ColorLightCyan,
	}

	return &Particle{
		X:      ba.rng.Intn(ba.width / 2), // Spawn on left half
		Y:      ba.rng.Intn(ba.height),
		Char:   chars[ba.rng.Intn(len(chars))],
		Color:  colors[ba.rng.Intn(len(colors))],
		VX:     1 + ba.rng.Intn(2), // Speed 1-2
		VY:     0,
		Age:    0,
		MaxAge: ba.width * 2, // Live long enough to cross screen
	}
}

// createSnowParticle creates a snowflake.
func (ba *BackgroundAnimator) createSnowParticle() *Particle {
	chars := []rune{'*', '❄', '❅', '❆', '·', '∗'}

	return &Particle{
		X:      ba.rng.Intn(ba.width),
		Y:      0, // Spawn at top
		Char:   chars[ba.rng.Intn(len(chars))],
		Color:  tcell.ColorWhite,
		VX:     0,
		VY:     1, // Fall down
		Age:    0,
		MaxAge: ba.height * 2,
	}
}

// Pause pauses the background animation.
func (ba *BackgroundAnimator) Pause() {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.paused = true
}

// Resume resumes the background animation.
func (ba *BackgroundAnimator) Resume() {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.paused = false
}

// Clear removes all particles.
func (ba *BackgroundAnimator) Clear() {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.particles = make([]*Particle, 0)
}

// Shutdown stops the background animator.
func (ba *BackgroundAnimator) Shutdown() {
	ba.cancel()
	ba.wg.Wait()
}

// Draw renders the background overlay.
func (bo *BackgroundOverlay) Draw(screen tcell.Screen) {
	// Draw the base box (usually transparent)
	bo.Box.Draw(screen)

	// Draw particles
	bo.animator.mu.RLock()
	particleCount := len(bo.animator.particles)
	bo.animator.mu.RUnlock()

	// DEBUG: Log draw calls periodically
	if bo.animator.rng.Intn(30) == 0 { // Log ~1 in 30 frames
		fmt.Fprintf(os.Stderr, "[DEBUG BackgroundOverlay] Draw called, particle count: %d\n", particleCount)
	}

	bo.animator.mu.RLock()
	defer bo.animator.mu.RUnlock()

	x, y, width, height := bo.GetInnerRect()

	for _, p := range bo.animator.particles {
		// Calculate screen position
		px := x + p.X
		py := y + p.Y

		// Bounds check
		if px < x || px >= x+width || py < y || py >= y+height {
			continue
		}

		// Draw particle with low intensity to not interfere with content
		style := tcell.StyleDefault.Foreground(p.Color).Dim(true)
		screen.SetContent(px, py, p.Char, nil, style)
	}
}

// GetEffectDescription returns a human-readable description of an effect.
func GetEffectDescription(effect BackgroundEffect) string {
	switch effect {
	case BackgroundHyperspace:
		return "Hyperspace stars moving from left to right"
	case BackgroundSnow:
		return "Gentle snowfall from top to bottom"
	case BackgroundNone:
		return "No background effect"
	default:
		return "Unknown effect"
	}
}

// ValidateConfig validates a background configuration.
func ValidateConfig(config BackgroundConfig) error {
	if config.Density < 0 || config.Density > 10 {
		return fmt.Errorf("density must be between 0 and 10")
	}
	if config.Speed < 10 || config.Speed > 1000 {
		return fmt.Errorf("speed must be between 10 and 1000 ms")
	}
	if config.MaxFPS < 1 || config.MaxFPS > 60 {
		return fmt.Errorf("maxFPS must be between 1 and 60")
	}
	return nil
}
