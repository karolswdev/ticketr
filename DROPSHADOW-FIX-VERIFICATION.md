# Drop Shadow Fix Verification Plan

**Fix Applied:** 2025-10-21
**Fix Location:** `/home/karol/dev/private/ticktr/internal/adapters/tui/views/workspace_modal.go` line 241

---

## 1. The Fix

### Change Applied

**File:** `internal/adapters/tui/views/workspace_modal.go`
**Line:** 241

```diff
- w.app.SetRoot(w.form, true)
+ w.app.SetRoot(w.Primitive(), true)
```

### Why This Fixes It

**Before:**
- `Show()` method called `app.SetRoot(w.form, true)`
- This displayed the **inner form** directly
- The `ShadowForm` wrapper was created but never rendered
- Result: No shadow visible

**After:**
- `Show()` method calls `app.SetRoot(w.Primitive(), true)`
- `Primitive()` returns `w.shadowForm` when shadows are enabled (line 519)
- tview renders the `ShadowForm` wrapper
- `ShadowForm.Draw()` renders shadow characters, then inner form
- Result: Shadow visible ✅

---

## 2. Build Verification

### Status: ✅ PASSED

```bash
$ go build ./cmd/ticketr
# (no output = success)
```

**Verdict:** Fix compiles successfully. No syntax errors, no type errors.

---

## 3. Manual Testing Checklist

### Test 1: Shadows Enabled (Default)

**Setup:**
```bash
# Default configuration has shadows enabled
# No environment variables needed
```

**Steps:**
1. Run `./ticketr tui`
2. Press 'c' to create workspace OR navigate to workspace list and press Enter on "Create New Workspace"
3. Observe the modal appearance

**Expected Result:**
```
╔════════════════════════╗▒
║  Create Workspace      ║▒
║                        ║▒
║  [Form fields...]      ║▒
║                        ║▒
║  [OK]  [Cancel]        ║▒
╚════════════════════════╝▒
 ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
```

**Verification Points:**
- [ ] Modal appears centered on screen
- [ ] `▒` characters visible on **right edge** of modal
- [ ] `▒` characters visible on **bottom edge** of modal
- [ ] Shadow characters are **gray/dim** color
- [ ] Shadow appears **behind** the modal (not overlapping content)
- [ ] Modal border is intact (not obscured by shadow)

**Acceptance:** ALL checkboxes must be checked.

---

### Test 2: Shadows Disabled via Environment Variable

**Setup:**
```bash
export TICKETR_EFFECTS_SHADOWS=false
```

**Steps:**
1. Run `./ticketr tui`
2. Press 'c' to create workspace
3. Observe the modal appearance

**Expected Result:**
```
╔════════════════════════╗
║  Create Workspace      ║
║                        ║
║  [Form fields...]      ║
║                        ║
║  [OK]  [Cancel]        ║
╚════════════════════════╝
```

**Verification Points:**
- [ ] Modal appears centered on screen
- [ ] NO `▒` characters visible
- [ ] Modal looks clean without shadows
- [ ] No visual artifacts or rendering glitches

**Acceptance:** ALL checkboxes must be checked.

---

### Test 3: Shadows Re-enabled

**Setup:**
```bash
export TICKETR_EFFECTS_SHADOWS=true
```

**Steps:**
1. Run `./ticketr tui`
2. Press 'c' to create workspace
3. Observe the modal appearance

**Expected Result:**
- Shadows visible again (same as Test 1)

**Verification Points:**
- [ ] Shadows reappear correctly
- [ ] No visual artifacts from toggling

**Acceptance:** ALL checkboxes must be checked.

---

### Test 4: Shadow Visibility Across Terminal Themes

Test shadows on different terminal background colors to ensure they're visible:

**Dark Terminal (Black/Dark Gray Background):**
- [ ] Shadows visible with sufficient contrast
- [ ] Shadow color is dim gray (tcell.ColorGray with Dim(true))

**Light Terminal (White/Light Background):**
- [ ] Shadows visible with sufficient contrast
- [ ] Shadow color adapts appropriately

**Mid-tone Terminal (Blue/Purple Background):**
- [ ] Shadows visible with sufficient contrast

**Acceptance:** Shadows must be visible on at least 2 of 3 terminal themes.

---

### Test 5: Modal Interaction with Shadows

**Steps:**
1. Open workspace modal with shadows enabled
2. Tab through form fields
3. Fill in fields
4. Submit form
5. Trigger error (e.g., empty required field)
6. Observe error state

**Verification Points:**
- [ ] Shadows remain visible during form interaction
- [ ] Shadows don't interfere with input
- [ ] Shadows update correctly when modal changes size
- [ ] Error modal also shows shadows (if applicable)
- [ ] Success modal shows shadows (if applicable)

**Acceptance:** ALL checkboxes must be checked.

---

### Test 6: Shadow Performance

**Steps:**
1. Open workspace modal with shadows
2. Rapidly tab through fields
3. Type rapidly in input fields
4. Monitor CPU usage

**Verification Points:**
- [ ] No visible lag or stuttering
- [ ] CPU usage remains reasonable (<10% spike)
- [ ] No rendering artifacts during rapid interaction

**Acceptance:** ALL checkboxes must be checked.

---

## 4. Code Review Verification

### Verify Shadow Implementation Remains Intact

**Check 1: ShadowForm.Draw() Implementation**

**File:** `internal/adapters/tui/effects/shadowbox.go`
**Lines:** 282-320

```go
func (sf *ShadowForm) Draw(screen tcell.Screen) {
	// Draw shadow if enabled
	if sf.shadowEnabled {
		// ... shadow rendering code ...
	}
	// Draw the form
	sf.form.Draw(screen)
}
```

**Verification:**
- [ ] Draw() method exists
- [ ] Shadow rendering logic is intact
- [ ] shadowEnabled check is present
- [ ] Form drawn after shadow (correct z-order)

---

### Check 2: Primitive() Method Returns Correct Wrapper

**File:** `internal/adapters/tui/views/workspace_modal.go`
**Lines:** 516-522

```go
func (w *WorkspaceModal) Primitive() tview.Primitive {
	// Return shadow form if shadows are enabled, otherwise return regular form
	if w.shadowForm != nil {
		return w.shadowForm
	}
	return w.form
}
```

**Verification:**
- [ ] Primitive() checks if shadowForm exists
- [ ] Returns shadowForm when available
- [ ] Falls back to regular form if shadows disabled

---

### Check 3: Show() Method Uses Primitive()

**File:** `internal/adapters/tui/views/workspace_modal.go`
**Line:** 241

```go
w.app.SetRoot(w.Primitive(), true)
```

**Verification:**
- [ ] Show() calls Primitive() instead of accessing w.form directly
- [ ] No other SetRoot calls in workspace_modal.go use w.form directly

---

### Check 4: ShadowForm Creation in setupForm()

**File:** `internal/adapters/tui/views/workspace_modal.go`
**Lines:** 58-66

```go
effectsConfig := theme.GetEffects()
if effectsConfig.DropShadows {
	w.shadowForm = effects.NewShadowForm()
	w.form = w.shadowForm.GetForm()
} else {
	w.form = tview.NewForm()
}
```

**Verification:**
- [ ] Checks theme.GetEffects().DropShadows
- [ ] Creates ShadowForm when enabled
- [ ] Stores inner form reference in w.form
- [ ] Falls back to regular form when disabled

---

## 5. Integration Testing

### Test: Shadow Configuration from Theme

**Test Code:**
```go
// Verify default theme enables shadows
func TestDefaultThemeEnablesShadows(t *testing.T) {
	effects := theme.DefaultVisualEffects()
	if !effects.DropShadows {
		t.Error("Default theme should enable drop shadows")
	}
}

// Verify GetEffects returns correct config
func TestGetEffectsReturnsDropShadows(t *testing.T) {
	theme.Set(theme.DefaultTheme)
	effects := theme.GetEffects()
	if !effects.DropShadows {
		t.Error("GetEffects should return DropShadows: true for default theme")
	}
}
```

**Run:**
```bash
cd internal/adapters/tui/theme
go test -run TestDefaultThemeEnablesShadows
go test -run TestGetEffectsReturnsDropShadows
```

**Acceptance:**
- [ ] Default theme has DropShadows: true
- [ ] GetEffects() returns correct configuration

---

### Test: ShadowForm Wrapper Creation

**Test Code:**
```go
// Verify ShadowForm wrapper is created
func TestWorkspaceModalCreatesShadowForm(t *testing.T) {
	// Setup
	app := tview.NewApplication()
	mockService := &MockWorkspaceService{}
	modal := views.NewWorkspaceModal(app, mockService)

	// Enable shadows in theme
	theme.Set(theme.DefaultTheme)

	// Force setupForm to be called
	modal.OnShow()

	// Verify shadowForm is created
	primitive := modal.Primitive()
	if primitive == nil {
		t.Fatal("Primitive() returned nil")
	}

	// Type assertion to verify it's a ShadowForm
	_, ok := primitive.(*effects.ShadowForm)
	if !ok {
		t.Error("Primitive() should return ShadowForm when shadows enabled")
	}
}
```

**Acceptance:**
- [ ] Test passes
- [ ] Primitive() returns ShadowForm instance

---

## 6. Visual Regression Test (Future)

**Goal:** Automate shadow verification using tcell.SimulationScreen

**Approach:**
1. Create SimulationScreen
2. Render workspace modal
3. Inspect screen cells at shadow positions
4. Verify `▒` characters present with correct style

**Status:** ⏳ NOT YET IMPLEMENTED
**Priority:** MEDIUM (nice to have, but manual testing is sufficient for now)

---

## 7. Documentation Updates

### Files to Update After Verification

**1. TUI-GUIDE.md**
- [ ] Add section on visual effects configuration
- [ ] Document TICKETR_EFFECTS_SHADOWS environment variable
- [ ] Include screenshot/ASCII art showing shadows

**2. VISUAL_EFFECTS_CONFIG.md**
- [ ] Update with verified shadow behavior
- [ ] Document shadow rendering details
- [ ] Add troubleshooting section (if shadows not visible)

**3. CHANGELOG.md**
- [ ] Add entry: "Fixed drop shadows not rendering on workspace modal"

---

## 8. Rollback Plan

**If fix causes issues:**

**Rollback Command:**
```bash
git checkout HEAD -- internal/adapters/tui/views/workspace_modal.go
```

**Rollback Verification:**
1. Rebuild: `go build ./cmd/ticketr`
2. Verify TUI still works (without shadows)
3. No crashes or rendering issues

**Risk:** VERY LOW - one-line change, easy to revert

---

## 9. Success Criteria

**Fix is considered successful when:**

### Functional Requirements (ALL must pass)
- [x] Code compiles without errors
- [ ] Manual Test 1 passes (shadows visible)
- [ ] Manual Test 2 passes (shadows can be disabled)
- [ ] Manual Test 3 passes (shadows can be re-enabled)
- [ ] Code Review checks ALL pass
- [ ] No regressions in modal functionality

### Visual Requirements (ALL must pass)
- [ ] Shadow characters (`▒`) visible on right edge
- [ ] Shadow characters (`▒`) visible on bottom edge
- [ ] Shadow color is dim gray (not too bright)
- [ ] Shadow doesn't obscure modal content
- [ ] Shadow appears behind modal (correct z-order)

### Performance Requirements (ALL must pass)
- [ ] No visible lag when opening modal
- [ ] CPU usage remains reasonable
- [ ] No rendering artifacts

### Configuration Requirements (ALL must pass)
- [ ] Shadows enabled by default (DropShadows: true)
- [ ] TICKETR_EFFECTS_SHADOWS=false disables shadows
- [ ] TICKETR_EFFECTS_SHADOWS=true enables shadows

---

## 10. Testing Timeline

**Estimated Time:** 30-45 minutes

**Phase 1: Build Verification (5 min)**
- [x] Compile code
- [x] Verify no errors

**Phase 2: Manual Testing (20 min)**
- [ ] Test 1: Shadows enabled (5 min)
- [ ] Test 2: Shadows disabled (3 min)
- [ ] Test 3: Toggle shadows (2 min)
- [ ] Test 4: Terminal themes (5 min)
- [ ] Test 5: Interaction (3 min)
- [ ] Test 6: Performance (2 min)

**Phase 3: Code Review (5 min)**
- [ ] Verify all 4 code checks

**Phase 4: Documentation (10 min)**
- [ ] Take screenshots
- [ ] Update verification report
- [ ] Update investigation report

---

## 11. Screenshot Evidence

**Required Screenshots:**

### Screenshot 1: Shadows Enabled (Default)
**File:** `screenshots/workspace-modal-shadows-on.png`
**Description:** Workspace modal with drop shadows visible
**Required content:**
- Full TUI screen
- Workspace modal centered
- Shadow characters visible on right and bottom
- Clear enough to see `▒` character detail

### Screenshot 2: Shadows Disabled
**File:** `screenshots/workspace-modal-shadows-off.png`
**Description:** Workspace modal with shadows disabled
**Required content:**
- Full TUI screen
- Workspace modal centered
- NO shadow characters visible
- Clean modal appearance

### Screenshot 3: Side-by-Side Comparison
**File:** `screenshots/workspace-modal-comparison.png`
**Description:** Before/after or on/off comparison
**Format:** Two terminal windows side-by-side or stacked

---

## 12. Known Limitations

**Documented Limitations:**

1. **Shadow Character Support:**
   - Requires terminal with Unicode support
   - `▒` character may not render on all terminals
   - Graceful degradation: if character doesn't render, shadow just won't appear

2. **Terminal Background Color:**
   - Shadow visibility depends on terminal background
   - Very light backgrounds may make shadows hard to see
   - No automatic shadow color adjustment (manual theme selection needed)

3. **Shadow Offset Fixed:**
   - Offset hardcoded to (2, 1) in ShadowForm
   - Not configurable via theme/environment
   - Future enhancement: make offset configurable

4. **Other Modals:**
   - Fix only applies to workspace modal
   - Bulk operations modal uses different architecture (tview.Pages)
   - Error modals use tview.Modal directly (no shadow support yet)
   - Future: Apply shadows to all modal types

---

## 13. Next Steps After Verification

### If ALL Tests Pass:
1. Mark fix as VERIFIED in investigation report
2. Update VERIFIER day 3 report with visual proof
3. Close emergency investigation
4. Consider applying shadows to other modals (future enhancement)

### If ANY Test Fails:
1. Document failure in detail
2. Re-investigate root cause
3. Apply additional fixes
4. Re-run verification
5. DO NOT mark as complete until all tests pass

---

## Appendix: Debug Checklist

**If shadows still don't appear after fix:**

### Debug Step 1: Verify Configuration Loaded
Add debug logging to workspace_modal.go setupForm():
```go
effectsConfig := theme.GetEffects()
fmt.Fprintf(os.Stderr, "DEBUG: DropShadows=%v\n", effectsConfig.DropShadows)
```

**Expected output:** `DEBUG: DropShadows=true`

---

### Debug Step 2: Verify ShadowForm Created
Add debug logging to workspace_modal.go setupForm():
```go
if effectsConfig.DropShadows {
	w.shadowForm = effects.NewShadowForm()
	fmt.Fprintf(os.Stderr, "DEBUG: ShadowForm created: %p\n", w.shadowForm)
	w.form = w.shadowForm.GetForm()
}
```

**Expected output:** `DEBUG: ShadowForm created: 0x...` (non-nil pointer)

---

### Debug Step 3: Verify Primitive() Returns ShadowForm
Add debug logging to workspace_modal.go Primitive():
```go
func (w *WorkspaceModal) Primitive() tview.Primitive {
	if w.shadowForm != nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Primitive() returning shadowForm\n")
		return w.shadowForm
	}
	fmt.Fprintf(os.Stderr, "DEBUG: Primitive() returning regular form\n")
	return w.form
}
```

**Expected output:** `DEBUG: Primitive() returning shadowForm`

---

### Debug Step 4: Verify Draw() Called
Add debug logging to shadowbox.go ShadowForm.Draw():
```go
func (sf *ShadowForm) Draw(screen tcell.Screen) {
	fmt.Fprintf(os.Stderr, "DEBUG: ShadowForm.Draw() called, enabled=%v\n", sf.shadowEnabled)
	if sf.shadowEnabled {
		// ... shadow rendering ...
		fmt.Fprintf(os.Stderr, "DEBUG: Shadow rendered at rect=%v\n", sf.form.GetRect())
	}
	sf.form.Draw(screen)
}
```

**Expected output:**
```
DEBUG: ShadowForm.Draw() called, enabled=true
DEBUG: Shadow rendered at rect=(x, y, width, height)
```

---

### Debug Step 5: Verify Shadow Characters Written
Add debug logging to shadowbox.go shadow rendering loop:
```go
for sx := x + sf.shadowOffsetX; sx < x+width+sf.shadowOffsetX; sx++ {
	if sx >= 0 {
		screen.SetContent(sx, shadowY, sf.shadowChar, nil, shadowStyle)
		if sx == x + sf.shadowOffsetX {
			fmt.Fprintf(os.Stderr, "DEBUG: Shadow char '%c' at (%d,%d)\n", sf.shadowChar, sx, shadowY)
		}
	}
}
```

**Expected output:** `DEBUG: Shadow char '▒' at (x,y)`

---

## Conclusion

This verification plan provides comprehensive testing coverage to confirm the drop shadow fix works correctly.

**Status:** ⏳ AWAITING MANUAL TESTING
**Next Action:** Run Manual Test 1 and capture screenshot

**Estimated completion:** 30-45 minutes from start of testing.
