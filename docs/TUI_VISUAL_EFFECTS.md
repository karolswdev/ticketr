# TUI Visual Effects System

## Overview

The TUI Visual Effects system transforms Ticketr's terminal interface from functional to enchanting. Built on the **Four Principles of TUI Excellence**, it provides subtle motion, depth, atmosphere, and quality craftsmanship while maintaining strict performance budgets and accessibility.

## The Four Principles

### 1. Subtle Motion is Life
Static interfaces feel dead. Animation must be subtle, never distracting. Every motion serves a purpose: feedback, guidance, or delight.

**Examples:**
- Active spinners (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏) provide real-time feedback
- Focus pulse draws attention to active panels
- Modal fade-in creates smooth transitions

### 2. Light, Shadow, and Focus
Create depth in 2D space. Guide user attention through visual hierarchy using borders, shadows, and gradients intentionally.

**Examples:**
- Double-line borders (╔═╗) for focused panels
- Drop shadows on modals (▒ offset by 1 row, 2 cols)
- Title gradients that shift color across focused panel headers

### 3. Atmosphere and Ambient Effects
Add character and personality when users opt-in. Must be themeable, default OFF, and never sacrifice performance or usability.

**Examples:**
- Hyperspace starfield (dark theme) - stars streak left to right
- Snow effect (arctic theme) - gentle snowfall from top to bottom
- Fully configurable density and speed

### 4. Small Charms of Quality
Details show craftsmanship. Tiny celebrations for user actions. Responsive feedback to every interaction.

**Examples:**
- Success sparkles (✦✧⋆∗·) - 500ms particle burst on completion
- Animated toggles ([ ]→[•]→[x]) - 3-frame transition
- Progress bar shimmer - subtle shine sweep across completed portion

## Architecture

### Package Structure

```
internal/adapters/tui/
├── effects/
│   ├── animator.go        # Core animation engine
│   ├── background.go      # Ambient background effects
│   ├── borders.go         # Border style utilities
│   ├── shadowbox.go       # Drop shadow primitives
│   └── shimmer.go         # Shimmer and gradient effects
└── theme/
    └── theme.go           # Enhanced theme with visual effects config
```

### Core Components

#### 1. Animator (`effects/animator.go`)

The centralized animation engine managing all visual effects with proper cancellation and CPU efficiency.

**Key Features:**
- Context-aware goroutines (all animations cancellable)
- Global motion kill switch (`animator.SetEnabled(false)`)
- Non-blocking event loop integration
- Automatic cleanup on shutdown

**Usage:**
```go
animator := effects.NewAnimator(app)

// Start a spinner
frame := 0
animator.Start("spinner", 80*time.Millisecond, effects.Spinner(&frame))

// Stop specific animation
animator.Stop("spinner")

// Shutdown all animations
defer animator.Shutdown()
```

#### 2. Background Animator (`effects/background.go`)

Renders ambient particle effects in the background.

**Effects:**
- **Hyperspace** - Stars streaking left to right (dark theme)
- **Snow** - Gentle snowfall from top to bottom (arctic theme)
- **None** - Default (effects OFF)

**Configuration:**
```go
config := effects.BackgroundConfig{
    Effect:     effects.BackgroundHyperspace,
    Density:    0.02,   // 2% of screen cells
    Speed:      100,    // 100ms per frame
    Enabled:    false,  // Default: OFF (user opt-in)
    MaxFPS:     15,     // Rate limiting
    AutoPause:  true,   // Pause when UI busy
}

ba := effects.NewBackgroundAnimator(app, config)
ba.Start()
defer ba.Shutdown()
```

**Performance:**
- CPU-efficient: ≤ 3% CPU on typical dev laptop
- Frame rate limited: 12-20 FPS (coalesced timers)
- Auto-pause when UI is busy
- All animations interruptible

#### 3. ShadowBox (`effects/shadowbox.go`)

Extended primitives with automatic drop shadow rendering.

**Components:**
- `ShadowBox` - Extends `tview.Box`
- `ShadowFlex` - Flex container with shadow
- `ShadowForm` - Form with shadow

**Usage:**
```go
// Create a modal with shadow
shadowForm := effects.NewShadowForm()
form := shadowForm.GetForm()

form.SetBorder(true).SetTitle(" Modal ")
form.AddInputField("Name", "", 20, nil, nil)

// Shadow is rendered automatically on Draw()
```

**Shadow Characteristics:**
- Character: ▒ (medium shade)
- Offset: 2 columns right, 1 row down
- Color: Dim gray
- Configurable via setters

#### 4. Shimmer Effects (`effects/shimmer.go`)

Animated polish for progress bars and titles.

**Features:**
- **Progress Bar Shimmer** - Sweeping brightness wave
- **Gradient Titles** - Two-color horizontal gradients
- **Rainbow Text** - Multi-color cycling effect
- **Pulse Intensity** - Sinusoidal brightness variation

**Usage:**
```go
// Progress bar shimmer
shimmer := effects.NewProgressBarShimmer(20, true)
shimmer.Update() // Call on each frame

// Apply to bar string
result := shimmer.Apply("[██████░░░░]", 6)

// Gradient title
title := effects.GradientTitle(" Focused Panel ", true)
// Returns: "[blue]Focused[-][cyan] Panel[-]"
```

## Theme Integration

### Visual Effects Configuration

Themes now include visual effects settings:

```go
type VisualEffects struct {
    // Motion effects
    Motion        bool   // Global motion kill switch
    Spinner       bool   // Active spinners
    FocusPulse    bool   // Focus pulse animation
    ModalFadeIn   bool   // Modal fade-in effect

    // Visual polish
    DropShadows    bool   // Drop shadows on modals
    GradientTitles bool   // Gradient titles on focused panels

    // Border styles
    FocusedBorder   BorderStyle  // double, single, rounded
    UnfocusedBorder BorderStyle

    // Ambient effects (default OFF)
    AmbientEnabled bool
    AmbientMode    string  // "hyperspace", "snow", "off"
    AmbientDensity float64 // 0.01 - 0.10 (1% - 10%)
    AmbientSpeed   int     // Milliseconds per frame
}
```

### Built-in Themes

#### Default Theme
Conservative, accessible, minimal effects:
```go
Effects: VisualEffects{
    Motion:         true,
    Spinner:        true,   // Essential feedback only
    FocusPulse:     false,
    ModalFadeIn:    false,
    DropShadows:    false,
    GradientTitles: false,
    FocusedBorder:  BorderStyleDouble,
    UnfocusedBorder: BorderStyleSingle,
    AmbientEnabled: false,
    AmbientMode:    "off",
}
```

#### Dark Theme
Blue accents with hyperspace option:
```go
Effects: VisualEffects{
    Motion:         true,
    Spinner:        true,
    FocusPulse:     false,  // OFF by default
    ModalFadeIn:    false,
    DropShadows:    false,
    GradientTitles: false,
    FocusedBorder:  BorderStyleDouble,
    UnfocusedBorder: BorderStyleSingle,
    AmbientEnabled: false,  // User must enable
    AmbientMode:    "hyperspace",
    AmbientDensity: 0.02,
    AmbientSpeed:   100,
}
```

#### Arctic Theme
Cyan tones with snow effect:
```go
Effects: VisualEffects{
    Motion:         true,
    Spinner:        true,
    FocusPulse:     false,
    ModalFadeIn:    false,
    DropShadows:    false,
    GradientTitles: false,
    FocusedBorder:  BorderStyleDouble,
    UnfocusedBorder: BorderStyleRounded,
    AmbientEnabled: false,
    AmbientMode:    "snow",
    AmbientDensity: 0.015,
    AmbientSpeed:   120,
}
```

### Accessing Effects

```go
// Get current theme's effects
effects := theme.GetEffects()

if effects.DropShadows {
    // Use shadow primitives
}

// Update effects
newEffects := effects
newEffects.AmbientEnabled = true
theme.SetEffects(newEffects)
```

## Performance Budgets

**NON-NEGOTIABLE LIMITS:**

| Component | Budget | Measurement |
|-----------|--------|-------------|
| Animations | ≤ 3% CPU | On typical dev laptop |
| Background effects | 12-20 FPS | Coalesced timers |
| Spinner | ≤ 0.5% CPU | While spinning |
| Event loop | Never block | Always use `QueueUpdateDraw()` |

**Enforcement:**
- Benchmark tests with CPU assertions
- Timer coalescing (use `time.Ticker`, not busy loops)
- Context-aware cancellation (all animations)
- Frame rate limiting (background effects)

## Accessibility

**MANDATORY FEATURES:**

1. **Global Motion Kill Switch**
   ```go
   animator.SetEnabled(false) // Disable all motion
   ```

2. **Graceful Degradation**
   - Work on 256-color terminals
   - Work without Unicode support (fallback characters)
   - Never impair legibility

3. **Reduced-Motion Support**
   - Honor `ui.motion: false` config
   - Individual effect toggles

4. **Low-Contrast Ambient**
   - Background effects use dim, low-contrast colors
   - Never interfere with text readability

## Configuration

### YAML Configuration

```yaml
ui:
  motion: true              # Global kill switch
  spinner: true
  focusPulse: false
  modalFadeIn: false
  dropShadows: false
  gradientTitles: false
  ambient:
    enabled: false
    mode: "hyperspace"      # or "snow", "off"
    density: 0.02           # 2% of cells
    speed: 100              # milliseconds
```

### Programmatic Configuration

```go
// Via theme
effects := theme.GetEffects()
effects.AmbientEnabled = true
effects.AmbientMode = "hyperspace"
theme.SetEffects(effects)

// Direct config
config := effects.BackgroundConfig{
    Effect:     effects.BackgroundHyperspace,
    Density:    0.02,
    Speed:      100,
    Enabled:    true,
    MaxFPS:     15,
}
```

## Animation Reference

### Spinner Frames
Braille spinner (80ms per frame):
```
⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
```

### Sparkle Frames
Success celebration (50ms per frame, 500ms total):
```
✦ → ✧ → ⋆ → ∗ → ·
```

### Toggle Animation
Checkbox transition (100ms per frame):
```
[ ] → [•] → [x]
```

### Fade Animation
Modal fade-in (33ms per phase, 100ms total):
```
░ (light) → ▒ (medium) → █ (solid)
```

### Border Styles

**Single Line (unfocused):**
```
┌─────┐
│     │
└─────┘
```

**Double Line (focused):**
```
╔═════╗
║     ║
╚═════╝
```

**Rounded (arctic theme):**
```
╭─────╮
│     │
╰─────╯
```

## Testing

### Unit Tests

```bash
# Run all effects tests
go test ./internal/adapters/tui/effects/... -v

# Run with short mode (skip slow tests)
go test ./internal/adapters/tui/effects/... -short

# Run benchmarks
go test ./internal/adapters/tui/effects/... -bench=.
```

### Performance Benchmarks

```bash
# Benchmark specific component
go test -bench=BenchmarkAnimator -benchmem ./internal/adapters/tui/effects/

# CPU profiling
go test -bench=BenchmarkBackgroundAnimator -cpuprofile=cpu.prof ./internal/adapters/tui/effects/
go tool pprof cpu.prof
```

### Cross-Terminal Compatibility

Test on multiple terminal emulators:
- **Linux:** xterm, gnome-terminal, kitty, alacritty
- **macOS:** Terminal.app, iTerm2
- **Windows:** Windows Terminal, ConEmu

Check for:
- Unicode character support
- Color rendering (256-color vs true-color)
- Performance on slow links (SSH)

## Troubleshooting

### Effects Not Visible

1. Check global motion switch:
   ```go
   if !animator.IsEnabled() {
       animator.SetEnabled(true)
   }
   ```

2. Verify theme effects configuration:
   ```go
   effects := theme.GetEffects()
   // Check individual flags
   ```

3. Ensure animations are started:
   ```go
   animator.Start("name", interval, handler)
   ```

### High CPU Usage

1. Check frame rate limits:
   ```go
   config.MaxFPS = 15 // Lower if needed
   ```

2. Reduce particle density:
   ```go
   config.Density = 0.01 // 1% instead of 2%
   ```

3. Disable ambient effects:
   ```go
   effects.AmbientEnabled = false
   ```

### Rendering Artifacts

1. Use `app.QueueUpdateDraw()` for all UI updates from goroutines
2. Ensure proper cleanup on animation stop
3. Check shadow rendering bounds (may clip on small terminals)

## Best Practices

1. **Default to OFF** - Users opt into enchantment
2. **Respect the kill switch** - Always check `animator.IsEnabled()`
3. **Use context for cancellation** - All animations must be stoppable
4. **Measure performance** - Run benchmarks before merging
5. **Test on multiple terminals** - Unicode/color support varies
6. **Document user-facing config** - Clear examples in config files
7. **Never block the event loop** - Use `QueueUpdateDraw()` from goroutines

## Future Enhancements

Potential additions (not in scope for Phase 6):

1. **Custom border characters** - When tview adds support
2. **Easing functions** - Smooth acceleration/deceleration
3. **Particle physics** - Gravity, collision for ambient effects
4. **Sound effects** - Terminal bell integration
5. **Transition animations** - Smooth view changes
6. **Loading skeletons** - Placeholder animations for async content

## References

- **Implementation:** `/home/karol/dev/private/ticktr/internal/adapters/tui/effects/`
- **Theme System:** `/home/karol/dev/private/ticktr/internal/adapters/tui/theme/theme.go`
- **Tests:** `/home/karol/dev/private/ticktr/internal/adapters/tui/effects/*_test.go`
- **Phase 6 Spec:** `/home/karol/dev/private/ticktr/docs/PHASE6-CLEAN-RELEASE.md`

---

**Version:** 1.0 (Phase 6, Day 12.5)
**Author:** TUIUX Agent
**Last Updated:** 2025-10-20
