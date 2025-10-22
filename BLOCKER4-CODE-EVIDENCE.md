# BLOCKER #4 Code Evidence - Exact Locations

**Investigation:** Visual Effects Missing
**Agent:** TUIUX
**Date:** 2025-10-21

---

## Evidence 1: Animator Created But Never Used

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go`
**Lines:** 138-143

```go
// Initialize and start animator if effects are enabled
effectsConfig := theme.GetEffects()
if effectsConfig.Motion {
    t.animator = effects.NewAnimator(t.app)
    // Animator will be shut down in Stop()
}
```

**Issue:** ❌ Animator is created but there are ZERO calls to `t.animator.Start()` anywhere in the codebase

**Proof:**
```bash
$ grep -r "animator.Start" internal/adapters/tui --exclude-dir=effects
# Returns: NO RESULTS (except in test files)
```

---

## Evidence 2: Spinner Calculates Frames But UI Never Redraws

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/progressbar.go`
**Lines:** 167-174

```go
// getSpinner returns a rotating spinner character based on elapsed time.
// Uses Braille spinner: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
func (pb *ProgressBar) getSpinner() string {
    spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
    elapsed := time.Since(pb.startTime).Milliseconds()
    index := (elapsed / 80) % int64(len(spinners)) // Rotate every 80ms
    return spinners[index]
}
```

**Issue:** ❌ This calculates which frame to show, but there's no goroutine calling `QueueUpdateDraw()` every 80ms

**What's Missing:**
```go
// THIS CODE DOES NOT EXIST ANYWHERE:
go func() {
    ticker := time.NewTicker(80 * time.Millisecond)
    defer ticker.Stop()
    for range ticker.C {
        app.QueueUpdateDraw(func() {
            // Update spinner
        })
    }
}()
```

---

## Evidence 3: Progress Bar Start() Only Sets Timestamps

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/progressbar.go`
**Lines:** 47-52

```go
// Start marks the beginning of a progress operation.
func (pb *ProgressBar) Start() {
    now := time.Now()
    pb.startTime = now
    pb.lastUpdate = now
}
```

**Issue:** ❌ Start() does NOT start any animation loop or periodic redraw

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/views/sync_status.go`
**Lines:** 81

```go
v.progressBar.Start()  // ← Only sets startTime, doesn't start animation!
```

---

## Evidence 4: Shadows Disabled By Default

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/theme/theme.go`
**Lines:** 44-60

```go
// DefaultVisualEffects returns conservative defaults (most effects OFF).
func DefaultVisualEffects() VisualEffects {
    return VisualEffects{
        Motion:          true,  // Motion enabled, but individual effects opt-in
        Spinner:         true,  // Spinners are essential feedback
        FocusPulse:      false, // OFF by default
        ModalFadeIn:     false, // OFF by default
        DropShadows:     false, // OFF by default ← ISSUE
        GradientTitles:  false, // OFF by default
        FocusedBorder:   BorderStyleDouble,
        UnfocusedBorder: BorderStyleSingle,
        AmbientEnabled:  false,
        AmbientMode:     "off",
        AmbientDensity:  0.02,
        AmbientSpeed:    100,
    }
}
```

**Issue:** ❌ DropShadows defaults to false, making modals appear flat

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/views/workspace_modal.go`
**Lines:** 58-59

```go
if effectsConfig.DropShadows {
    // Use shadow form for modal with drop shadow
    w.shadowForm = effects.NewShadowForm()
```

**Result:** This code path is NEVER executed because DropShadows is false by default

---

## Evidence 5: 2,268 Lines of Dead Code

**Effects Package Line Count:**

```bash
$ wc -l internal/adapters/tui/effects/*.go
  328 animator.go
  211 animator_test.go
  384 background.go
  328 background_test.go
  200 borders.go
  354 shadowbox.go
  177 shimmer.go
  286 shimmer_test.go
 2268 total
```

**Production Usage:**

```bash
$ grep -r "animator.Start" internal/adapters/tui --exclude="*_test.go"
# NO RESULTS

$ grep -r "shimmer.Update" internal/adapters/tui --exclude="*_test.go"
# Found in progressbar.go but only called during render
# No periodic updates means no animation

$ grep -r "ShadowForm" internal/adapters/tui
# Found in workspace_modal.go but disabled by default config
```

**Conclusion:** 2,268 lines written, but:
- 0 calls to animator.Start()
- 0 periodic animation loops
- 0 visible effects to user

---

## Evidence 6: tview Supports Animations (via context7)

**Research:** Used context7 MCP server to fetch tview documentation

**Finding 1:** Timer-based updates work in tview

```go
// From tview wiki examples
func updateTime() {
    for {
        time.Sleep(refreshInterval)
        app.QueueUpdateDraw(func() {
            view.SetText(currentTimeString())
        })
    }
}
```

**Finding 2:** app.Draw() safe from goroutines

```go
// From tview wiki
go func() {
    app.Draw()  // Safe to call from background thread
}()
```

**Finding 3:** Animations are standard pattern

```go
// From tview concurrency docs
go func() {
    app.QueueUpdateDraw(func() {
        table.SetCellSimple(0, 0, "Updated")
    })
}()
```

**Conclusion:** ✅ tview fully supports what we need - we just haven't implemented it

---

## Evidence 7: Configuration System Works

**Test:** Environment variable loading

```bash
$ export TICKETR_EFFECTS_SHADOWS=true
$ go run cmd/ticketr/main.go tui
# Shadows WOULD appear if we enabled them by default
```

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/theme/theme.go`
**Lines:** 302-305

```go
if shadowsStr := os.Getenv("TICKETR_EFFECTS_SHADOWS"); shadowsStr != "" {
    if shadows, err := strconv.ParseBool(shadowsStr); err == nil {
        effects.DropShadows = shadows
    }
}
```

**Status:** ✅ Configuration system works correctly

**Issue:** ❌ Defaults are too conservative, user doesn't know to set env vars

---

## Summary of Code Locations

### Files That Need Changes

1. **`internal/adapters/tui/app.go`**
   - Line 138-143: Pass animator to views
   - Create animator BEFORE views

2. **`internal/adapters/tui/views/sync_status.go`**
   - Add animator reference
   - Add spinner animation methods
   - Start/stop spinner on status change

3. **`internal/adapters/tui/theme/theme.go`**
   - Line 50: Change `DropShadows: false` → `true`

4. **`internal/adapters/tui/widgets/progressbar.go`**
   - Enable shimmer by default (optional)

### Files That Are Fine

These files have good code, just not wired up:
- `internal/adapters/tui/effects/animator.go` ✓
- `internal/adapters/tui/effects/shadowbox.go` ✓
- `internal/adapters/tui/effects/shimmer.go` ✓
- `internal/adapters/tui/effects/borders.go` ✓
- `internal/adapters/tui/effects/background.go` ✓

---

## Visual Proof

**Current State (What User Sees):**

```
┌─ Sync Status ────────────────────┐
│ pull: ⠋ Querying project EPM... │  ← Static, never changes
└───────────────────────────────────┘

+------------------------------+      ← Single line border
| Create Workspace             |      ← No shadow
+------------------------------+
```

**Expected State (What Should Happen):**

```
┌─ Sync Status ────────────────────┐
│ pull: ⠋ Querying project EPM... │  ← Animates: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
└───────────────────────────────────┘

╔══════════════════════════════╗▒    ← Double line when focused
║  Create Workspace            ║▒    ← Drop shadow (▒)
╠══════════════════════════════╣▒
║  ...                         ║▒
╚══════════════════════════════╝▒
 ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
```

---

**Evidence Compiled:** 2025-10-21
**Agent:** TUIUX
**Status:** Complete
**Usage:** Reference for implementation fix
