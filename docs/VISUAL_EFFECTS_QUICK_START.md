# TUI Visual Effects - Quick Start Guide

## 5-Minute Integration Guide

### Step 1: Create the Animator

```go
import "github.com/karolswdev/ticktr/internal/adapters/tui/effects"

// In your TUIApp struct
type TUIApp struct {
    app      *tview.Application
    animator *effects.Animator
    // ... other fields
}

// In NewTUIApp
func NewTUIApp() *TUIApp {
    app := tview.NewApplication()
    tuiApp := &TUIApp{
        app:      app,
        animator: effects.NewAnimator(app),
    }
    return tuiApp
}

// In shutdown
func (t *TUIApp) Shutdown() {
    t.animator.Shutdown()
}
```

### Step 2: Add a Spinner

```go
// When starting an async operation
frame := 0
t.animator.Start("sync-spinner", 80*time.Millisecond, effects.Spinner(&frame))

// Update your status bar to show the spinner
t.statusView.SetText(effects.SpinnerFrames[frame] + " Syncing...")

// When operation completes
t.animator.Stop("sync-spinner")
```

### Step 3: Add Drop Shadows to Modals

```go
import "github.com/karolswdev/ticktr/internal/adapters/tui/effects"

// Replace your modal Form
// OLD: form := tview.NewForm()
// NEW:
shadowForm := effects.NewShadowForm()
form := shadowForm.GetForm()

// Configure form as usual
form.SetBorder(true).SetTitle(" Create Workspace ")
form.AddInputField("Name", "", 20, nil, nil)

// Use shadowForm as your primitive
app.SetRoot(shadowForm, true)
```

### Step 4: Add Background Effect (Optional)

```go
// Create background animator
config := effects.DefaultBackgroundConfig()
config.Effect = effects.BackgroundHyperspace
config.Enabled = true  // User must opt-in

bgAnimator := effects.NewBackgroundAnimator(app, config)
bgAnimator.Start()

// Overlay on your main layout
overlay := bgAnimator.GetOverlay()
// Render overlay behind your main content (requires tview Grid/Layers)
```

## Common Patterns

### Pattern 1: Async Task with Spinner

```go
func (t *TUIApp) handlePull() {
    // Start spinner
    frame := 0
    t.animator.Start("pull", 80*time.Millisecond, effects.Spinner(&frame))

    go func() {
        // Do async work
        result := t.pullService.Pull()

        // Stop spinner when done
        t.app.QueueUpdateDraw(func() {
            t.animator.Stop("pull")
            t.showResult(result)
        })
    }()
}
```

### Pattern 2: Success Sparkle

```go
func (t *TUIApp) showSuccess(message string) {
    // Create sparkles at random positions
    sparkles := []*effects.Sparkle{
        effects.NewSparkle(10, 5),
        effects.NewSparkle(15, 5),
        effects.NewSparkle(20, 5),
    }

    // Animate sparkles
    t.animator.Start("sparkle", 50*time.Millisecond, func() bool {
        allDead := true
        for _, s := range sparkles {
            if s.Update() {
                allDead = false
            }
        }
        return !allDead  // Continue until all sparkles fade
    })

    t.statusView.SetText("✓ " + message)
}
```

### Pattern 3: Progress Bar with Shimmer

```go
// Create shimmer
shimmer := effects.NewProgressBarShimmer(30, true)

// In your progress update loop
t.animator.Start("shimmer", 100*time.Millisecond, func() bool {
    shimmer.Update()

    // Render progress bar
    bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"

    // Apply shimmer
    bar = shimmer.Apply(bar, filled)

    t.progressView.SetText(bar)
    return !completed
})
```

### Pattern 4: Focused Panel with Gradient Title

```go
import "github.com/karolswdev/ticktr/internal/adapters/tui/effects"

func (v *MyView) SetFocused(focused bool) {
    v.focused = focused

    // Update border color
    if focused {
        v.box.SetBorderColor(theme.GetPrimaryColor())
    } else {
        v.box.SetBorderColor(theme.GetSecondaryColor())
    }

    // Apply gradient title if focused
    title := " My Panel "
    if focused && theme.GetEffects().GradientTitles {
        title = effects.GradientTitle(title, true)
    }
    v.box.SetTitle(title)
}
```

## Configuration Examples

### Minimal (Default)

```yaml
ui:
  motion: true
  spinner: true
  # All other effects OFF
```

### Balanced

```yaml
ui:
  motion: true
  spinner: true
  focusPulse: false
  modalFadeIn: true
  dropShadows: true
  gradientTitles: false
  ambient:
    enabled: false
```

### Maximum Enchantment

```yaml
ui:
  motion: true
  spinner: true
  focusPulse: true
  modalFadeIn: true
  dropShadows: true
  gradientTitles: true
  ambient:
    enabled: true
    mode: "hyperspace"
    density: 0.03
    speed: 80
```

### Performance-Constrained

```yaml
ui:
  motion: true
  spinner: true
  # Disable everything else for max performance
  ambient:
    enabled: false
```

## Checklist

- [ ] Create animator in app initialization
- [ ] Shutdown animator on exit
- [ ] Use spinners for all async operations lasting > 500ms
- [ ] Apply shadows to all modals (if effects enabled)
- [ ] Test with effects disabled (`ui.motion: false`)
- [ ] Run performance benchmarks
- [ ] Test on multiple terminal emulators
- [ ] Document user-facing config options

## Troubleshooting

**Q: Animations not showing?**
A: Check `animator.IsEnabled()` and theme effects configuration.

**Q: High CPU usage?**
A: Reduce background effect density or disable ambient effects.

**Q: Flicker or artifacts?**
A: Ensure using `app.QueueUpdateDraw()` for all UI updates from goroutines.

**Q: Unicode characters not displaying?**
A: Terminal may not support Unicode. Provide ASCII fallback or use simpler characters.

## Performance Targets

- Animator: < 3% CPU with all effects enabled
- Background effects: 12-20 FPS maximum
- Spinner: < 0.5% CPU while active
- Zero impact when effects disabled

## Next Steps

1. Read full documentation: `docs/TUI_VISUAL_EFFECTS.md`
2. Review examples in test files: `internal/adapters/tui/effects/*_test.go`
3. Experiment with theme customization
4. Profile your specific use case

---

**Need help?** Check `docs/TUI_VISUAL_EFFECTS.md` for detailed information.
