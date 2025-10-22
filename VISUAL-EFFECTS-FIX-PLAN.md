# Visual Effects Fix Plan - Phase 6.5 Day 1 Afternoon

**Priority:** CRITICAL - BLOCKER #4
**Timeline:** 8-10 hours total
**Owner:** TUIUX Agent (Afternoon Session)
**Reviewed By:** Director

---

## Quick Reference

**Root Cause:** Animator created but never used - no periodic redraws
**Framework:** tview is fully capable - no migration needed
**Fix Scope:** Wire existing effects infrastructure to UI components

---

## Implementation Plan

### Phase 1: Spinner Animation (2 hours) - CRITICAL

**Objective:** Make spinners animate during operations

**Files to Modify:**
1. `/home/karol/dev/private/ticktr/internal/adapters/tui/views/sync_status.go`

**Changes:**

```go
// Add spinner animation tracking
type SyncStatusView struct {
    // ... existing fields ...
    spinnerFrame int
    animator     *effects.Animator  // Add animator reference
}

// Modify NewSyncStatusView to accept animator
func NewSyncStatusView(animator *effects.Animator) *SyncStatusView {
    view := &SyncStatusView{
        textView:     tview.NewTextView(),
        animator:     animator,
        spinnerFrame: 0,
        // ...
    }
    return view
}

// Add method to start spinner animation
func (v *SyncStatusView) startSpinnerAnimation() {
    if v.animator != nil {
        v.animator.Start("status-spinner", 80*time.Millisecond, func() bool {
            v.spinnerFrame = (v.spinnerFrame + 1) % len(effects.SpinnerFrames)
            v.updateDisplay()
            return v.status.State == sync.StateSyncing
        })
    }
}

// Add method to stop spinner
func (v *SyncStatusView) stopSpinnerAnimation() {
    if v.animator != nil {
        v.animator.Stop("status-spinner")
    }
}

// Modify SetStatus to start/stop animation
func (v *SyncStatusView) SetStatus(status sync.SyncStatus) {
    v.status = status
    
    if status.State == sync.StateSyncing {
        v.startSpinnerAnimation()
    } else {
        v.stopSpinnerAnimation()
    }
    
    v.updateDisplay()
}

// Modify updateDisplay to use animated frame
func (v *SyncStatusView) updateDisplay() {
    // ... existing code ...
    
    if v.status.State == sync.StateSyncing {
        spinner := effects.SpinnerFrames[v.spinnerFrame]
        // Use spinner in display
    }
}
```

2. `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go`

**Changes:**

```go
// Pass animator to SyncStatusView
func (t *TUIApp) setupApp() error {
    // ... existing code ...
    
    // Create animator FIRST (before views)
    effectsConfig := theme.GetEffects()
    if effectsConfig.Motion {
        t.animator = effects.NewAnimator(t.app)
    }
    
    // Create sync status view WITH animator
    t.syncStatusView = views.NewSyncStatusView(t.animator)
    
    // ... rest of setup ...
}
```

**Testing:**
```bash
# Test spinner animation
go run cmd/ticketr/main.go tui
# Press 'P' to trigger pull
# Verify spinner animates: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
```

**Acceptance:**
- [ ] Spinner cycles through all 10 frames
- [ ] Animation stops when operation completes
- [ ] CPU usage <1% during animation
- [ ] No flickering

---

### Phase 2: Enable Drop Shadows (1 hour) - HIGH IMPACT

**Objective:** Show shadows on modals by default

**Files to Modify:**
1. `/home/karol/dev/private/ticktr/internal/adapters/tui/theme/theme.go`

**Changes:**

```go
func DefaultVisualEffects() VisualEffects {
    return VisualEffects{
        Motion:          true,
        Spinner:         true,
        FocusPulse:      false,
        ModalFadeIn:     false,
        DropShadows:     true,  // ← Change from false to true
        GradientTitles:  false,
        FocusedBorder:   BorderStyleDouble,
        UnfocusedBorder: BorderStyleSingle,
        AmbientEnabled:  false,
        AmbientMode:     "off",
        AmbientDensity:  0.02,
        AmbientSpeed:    100,
    }
}
```

**Testing:**
```bash
go run cmd/ticketr/main.go tui
# Press F2 to open workspace modal
# Verify shadow visible: ▒ characters on right and bottom edge
```

**Acceptance:**
- [ ] Modal has visible shadow (▒)
- [ ] Shadow offset: 2 cols right, 1 row down
- [ ] Shadow doesn't overlap modal content
- [ ] Works on all modals

---

### Phase 3: Progress Bar Animation (2 hours) - MEDIUM

**Objective:** Smooth progress bar updates during operations

**Files to Modify:**
1. `/home/karol/dev/private/ticktr/internal/adapters/tui/views/sync_status.go`

**Changes:**

```go
// Add progress animation
func (v *SyncStatusView) UpdateProgress(progress jobs.JobProgress) {
    v.currentProgress = &progress
    v.showProgress = true

    if v.jobStartTime.IsZero() {
        v.jobStartTime = time.Now()
        v.progressBar.Start()
        
        // Start progress bar redraw loop
        if v.animator != nil {
            v.animator.Start("progress-update", 100*time.Millisecond, func() bool {
                v.updateDisplay()
                return v.showProgress
            })
        }
    }

    v.updateDisplay()
}

// Stop progress animation when cleared
func (v *SyncStatusView) ClearProgress() {
    v.currentProgress = nil
    v.showProgress = false
    v.jobStartTime = time.Time{}
    
    if v.animator != nil {
        v.animator.Stop("progress-update")
    }
    
    v.updateDisplay()
}
```

**Testing:**
```bash
go run cmd/ticketr/main.go tui
# Press 'P' to trigger pull with progress
# Verify progress bar updates smoothly
# Verify percentage increases in real-time
```

**Acceptance:**
- [ ] Progress bar fills smoothly
- [ ] Percentage updates in real-time
- [ ] ETA calculation visible
- [ ] CPU usage <2%

---

### Phase 4: Shimmer Effect (2 hours) - OPTIONAL

**Objective:** Add shimmer sweep to progress bars

**Files to Modify:**
1. `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/progressbar.go`

**Changes:**

```go
// Enable shimmer in NewProgressBar
func NewProgressBar(width int) *ProgressBar {
    if width < 10 {
        width = 10
    }

    effectsConfig := theme.GetEffects()
    shimmerEnabled := effectsConfig.Motion && effectsConfig.FocusPulse  // Enable shimmer

    return &ProgressBar{
        width:   width,
        shimmer: effects.NewProgressBarShimmer(width, shimmerEnabled),
    }
}
```

2. Enable in theme (optional):

```go
// internal/adapters/tui/theme/theme.go
func DefaultVisualEffects() VisualEffects {
    return VisualEffects{
        // ...
        FocusPulse:      true,  // ← Enable for shimmer
        // ...
    }
}
```

**Testing:**
```bash
export TICKETR_EFFECTS_SHIMMER=true
go run cmd/ticketr/main.go tui
# Watch for shimmer sweep across progress bar
```

**Acceptance:**
- [ ] Shimmer sweeps left to right
- [ ] Sweep completes every 2 seconds
- [ ] Subtle, not distracting
- [ ] CPU <1%

---

### Phase 5: Testing & Polish (2 hours) - MANDATORY

**Objective:** Verify all effects work correctly

**Test Checklist:**

**Visual Acceptance:**
- [ ] Record 10-second video of TUI
- [ ] Verify spinner animates
- [ ] Verify progress bar smooth
- [ ] Verify modals have shadows
- [ ] Verify no flickering

**Performance:**
- [ ] Measure CPU during animation (<5%)
- [ ] Check FPS (should be 30+ for animations)
- [ ] Verify no memory leaks

**Cross-Terminal:**
- [ ] Test in iTerm2
- [ ] Test in gnome-terminal
- [ ] Test in kitty
- [ ] Test in standard xterm

**Accessibility:**
- [ ] Test with TICKETR_EFFECTS_MOTION=false
- [ ] Verify graceful degradation
- [ ] Check 256-color terminals

**User Scenarios:**
- [ ] Pull operation shows animated spinner
- [ ] Push operation shows progress bar
- [ ] Workspace modal has shadow
- [ ] ESC cancels without artifacts

---

## Minimal Viable Fix (4 hours)

If time-constrained, implement ONLY:

1. **Phase 1:** Spinner Animation (2 hours) - MUST HAVE
2. **Phase 2:** Drop Shadows (1 hour) - MUST HAVE
3. **Phase 5:** Basic Testing (1 hour) - MUST HAVE

**Defer to v3.2.0:**
- Progress bar animation refinement
- Shimmer effects
- Focus pulse
- Ambient effects

---

## Common Issues & Solutions

**Issue:** Spinner doesn't animate
**Solution:** Check animator is passed to SyncStatusView, verify Motion=true

**Issue:** Shadow not visible
**Solution:** Check DropShadows=true in theme config

**Issue:** High CPU usage
**Solution:** Reduce animation frame rate, check for animation leaks

**Issue:** Flickering
**Solution:** Use QueueUpdateDraw() not app.Draw() for updates

**Issue:** Animation doesn't stop
**Solution:** Verify animator.Stop() is called on state change

---

## Code Quality Checklist

Before marking complete:

- [ ] All goroutines have context cancellation
- [ ] No busy loops (all use time.Ticker)
- [ ] QueueUpdateDraw() used for UI updates
- [ ] Animations stop on completion
- [ ] CPU usage within budget (<5%)
- [ ] No memory leaks
- [ ] Code follows existing patterns
- [ ] Comments explain WHY not just WHAT

---

## Handoff to Builder (if needed)

If TUIUX needs Builder assistance:

**What to hand off:**
- This plan
- Specific file changes
- Test requirements

**What Builder provides:**
- Code implementation
- Basic testing
- Integration

**What TUIUX validates:**
- Visual quality
- Animation smoothness
- UX polish

---

**Created:** 2025-10-21
**Owner:** TUIUX Agent
**Status:** Ready for Implementation
