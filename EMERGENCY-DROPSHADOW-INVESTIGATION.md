# Emergency Drop Shadow Investigation Report

**Investigation Date:** 2025-10-21
**Investigator:** TUIUX Agent
**Issue:** Drop shadows are NOT visible on modals despite claims they work
**User Visual Quality Rating:** 4/10 - "Still looks back in the 80s"

---

## Executive Summary

**ROOT CAUSE IDENTIFIED:** Drop shadows are completely bypassed due to a critical integration gap.

The workspace modal uses `ShadowForm` wrapper (which implements shadows), but **the wrapper is NEVER rendered to the screen**. Instead, only the inner form is displayed via `app.SetRoot(w.form, true)` at line 241 of workspace_modal.go.

**Status:** ❌ Shadows are implemented but not integrated
**Severity:** HIGH - Complete feature failure despite code existing
**Fix Complexity:** LOW - One-line change to use correct primitive

---

## 1. Root Cause Analysis

### Investigation Trail

#### Evidence 1: Configuration is Correct ✅
**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/theme/theme.go`
**Line 51:** `DropShadows: true` is explicitly enabled in DefaultTheme

```go
func DefaultVisualEffects() VisualEffects {
	return VisualEffects{
		// ...
		DropShadows:     true,  // ENABLED for professional appearance
		// ...
	}
}
```

**Verdict:** Configuration is correct and loaded properly.

---

#### Evidence 2: ShadowForm Wrapper is Created ✅
**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/views/workspace_modal.go`
**Lines 58-66:** ShadowForm is properly created when shadows are enabled

```go
func (w *WorkspaceModal) setupForm() {
	// Check if shadows are enabled
	effectsConfig := theme.GetEffects()
	if effectsConfig.DropShadows {
		// Use shadow form for modal with drop shadow
		w.shadowForm = effects.NewShadowForm()
		w.form = w.shadowForm.GetForm()
	} else {
		// Use regular form without shadow
		w.form = tview.NewForm()
	}
	// ...
}
```

**Verdict:** Wrapper is created correctly and stores reference in `w.shadowForm`.

---

#### Evidence 3: ShadowForm.Draw() Implementation is Correct ✅
**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/effects/shadowbox.go`
**Lines 282-320:** Draw method properly renders shadows using `▒` characters

```go
func (sf *ShadowForm) Draw(screen tcell.Screen) {
	// Draw shadow if enabled
	if sf.shadowEnabled {
		x, y, width, height := sf.form.GetRect()

		shadowStyle := tcell.StyleDefault.
			Foreground(sf.shadowColor).
			Background(tcell.ColorDefault).
			Dim(true)

		// Draw bottom shadow (horizontal line)
		shadowY := y + height + sf.shadowOffsetY - 1
		if shadowY >= 0 {
			for sx := x + sf.shadowOffsetX; sx < x+width+sf.shadowOffsetX; sx++ {
				if sx >= 0 {
					screen.SetContent(sx, shadowY, sf.shadowChar, nil, shadowStyle)
				}
			}
		}

		// Draw right shadow (vertical line)
		shadowX := x + width + sf.shadowOffsetX - 1
		if shadowX >= 0 {
			for sy := y + sf.shadowOffsetY; sy < y+height+sf.shadowOffsetY; sy++ {
				if sy >= 0 {
					screen.SetContent(shadowX, sy, sf.shadowChar, nil, shadowStyle)
				}
			}
		}

		// Draw corner shadow
		if shadowX >= 0 && shadowY >= 0 {
			screen.SetContent(shadowX, shadowY, sf.shadowChar, nil, shadowStyle)
		}
	}

	// Draw the form
	sf.form.Draw(screen)
}
```

**Verdict:** Draw implementation is correct and will render shadows if called.

---

#### Evidence 4: ❌ THE SMOKING GUN - Wrong Primitive is Displayed

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/views/workspace_modal.go`
**Line 241:** `app.SetRoot(w.form, true)` displays the INNER form, not the ShadowForm wrapper

```go
func (w *WorkspaceModal) Show() {
	// Load credential profiles
	if err := w.loadProfiles(); err != nil {
		// Show error and close
		w.showError(fmt.Sprintf("Failed to load credential profiles: %v", err))
		return
	}

	// Setup profile dropdown
	w.setupProfileDropdown()

	// Rebuild form to ensure proper layout
	w.buildForm()

	// Focus the name field
	w.app.SetFocus(w.nameField)

	// Show form as root
	w.app.SetRoot(w.form, true)  // ❌ BUG: Shows inner form, not shadowForm wrapper
}
```

**Compare with the Primitive() method at line 517:**

```go
func (w *WorkspaceModal) Primitive() tview.Primitive {
	// Return shadow form if shadows are enabled, otherwise return regular form
	if w.shadowForm != nil {
		return w.shadowForm  // ✅ CORRECT: Returns wrapper
	}
	return w.form
}
```

**The Problem:**
- `Show()` method: Uses `w.form` directly → bypasses shadows
- `Primitive()` method: Returns `w.shadowForm` → would show shadows if used
- **The Primitive() method is NEVER called when showing the modal**

---

#### Evidence 5: Integration Point in app.go
**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go`
**Lines 976-981:** showWorkspaceModal calls Show() method directly

```go
func (t *TUIApp) showWorkspaceModal() {
	t.inModal = true
	t.workspaceModal.OnShow()
	t.workspaceModal.Show()  // Calls Show() which bypasses Primitive()
}
```

---

### Root Cause Summary

**The exact failure chain:**

1. ✅ Configuration enables shadows: `DropShadows: true`
2. ✅ ShadowForm wrapper is created: `w.shadowForm = effects.NewShadowForm()`
3. ✅ Inner form stored: `w.form = w.shadowForm.GetForm()`
4. ✅ Primitive() method correctly returns wrapper
5. ❌ **FAILURE:** Show() method uses `app.SetRoot(w.form, true)` instead of `app.SetRoot(w.shadowForm, true)`
6. ❌ Result: Only inner form renders, shadow wrapper's Draw() never called
7. ❌ User sees modal without shadows

**Hypothesis Confirmed:** #6 - Complete Integration Gap

We wrote working shadow code but never integrated it properly. The Show() method hardcoded display of the inner form, completely bypassing the wrapper pattern we created.

---

## 2. Integration Status

| Component | Status | Evidence |
|-----------|--------|----------|
| Is `DropShadows` config loaded? | ✅ YES | theme.go line 51, enabled by default |
| Is modal_wrapper actually used? | ❌ NO | modal_wrapper.go exists but workspace_modal doesn't use it |
| Is ShadowForm created? | ✅ YES | workspace_modal.go line 61 |
| Is ShadowBox.Draw() implemented? | ✅ YES | shadowbox.go line 282-320 |
| Is ShadowForm.Draw() called? | ❌ NO | Show() bypasses wrapper, displays inner form directly |
| Are shadows visible to user? | ❌ NO | Wrapper never renders |

**Integration Gap:** The Show() method hardcodes `app.SetRoot(w.form, true)` instead of using the Primitive() method or directly referencing shadowForm wrapper.

---

## 3. The Fix

### Specific Solution

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/views/workspace_modal.go`
**Line:** 241
**Change:**

```go
// BEFORE (broken):
w.app.SetRoot(w.form, true)

// AFTER (fixed):
w.app.SetRoot(w.Primitive(), true)
```

**Alternative (more explicit):**

```go
// BEFORE (broken):
w.app.SetRoot(w.form, true)

// AFTER (fixed):
if w.shadowForm != nil {
	w.app.SetRoot(w.shadowForm, true)
} else {
	w.app.SetRoot(w.form, true)
}
```

### Verification

After fix, the render chain will be:

1. `app.SetRoot(w.Primitive(), true)` called
2. `Primitive()` returns `w.shadowForm` (when shadows enabled)
3. tview calls `shadowForm.Draw(screen)`
4. ShadowForm.Draw() renders shadow characters using `screen.SetContent()`
5. ShadowForm.Draw() calls `w.form.Draw(screen)` to render form on top
6. **Result:** User sees modal with drop shadows

### Testing Approach

**Manual Test:**
1. Apply one-line fix
2. Run `go build ./cmd/ticketr`
3. Run `./ticketr tui`
4. Press 'c' to create workspace (or trigger workspace modal)
5. **Expected:** Modal appears with visible `▒` shadow characters on right and bottom edges

**Code Verification:**
```bash
# Add debug logging to verify wrapper is being used
# In ShadowForm.Draw(), add at line 283:
fmt.Fprintf(os.Stderr, "DEBUG: ShadowForm.Draw() called, enabled=%v\n", sf.shadowEnabled)
```

---

## 4. Framework Assessment

### Question: Can tview achieve drop shadows?

**Answer: ✅ YES - Proven by existing implementation**

The ShadowForm implementation demonstrates that tview's `screen.SetContent()` API supports rendering decorative overlay characters. The approach is sound:

1. **Primitive Interface:** ShadowForm implements `tview.Primitive` interface
2. **Draw Override:** Draw() method intercepts rendering
3. **Layer Control:** Shadow drawn first, then form drawn on top
4. **Character Overlay:** `screen.SetContent()` allows placing `▒` characters anywhere
5. **Style Control:** Shadow uses dim gray style for depth effect

### Framework Capability Evidence

**From shadowbox.go lines 287-290:**
```go
shadowStyle := tcell.StyleDefault.
	Foreground(sf.shadowColor).
	Background(tcell.ColorDefault).
	Dim(true)
```

**From shadowbox.go lines 307-308:**
```go
screen.SetContent(shadowX, sy, sf.shadowChar, nil, shadowStyle)
```

This proves tview/tcell can:
- ✅ Render custom characters at arbitrary positions
- ✅ Apply styling (color, dim effect)
- ✅ Layer multiple primitives (shadow then content)
- ✅ Create visual depth through character overlays

**Limitation:** Shadows must be part of the primitive's own Draw() method. Cannot be applied externally after rendering. This is why the wrapper pattern is necessary.

---

## 5. Honest Assessment

### Did we mislead about shadows working?

**Yes, we did.** Here's the timeline:

**Day 1 (BUILDER):**
- Claimed: "Enabled DropShadows: true in theme.go"
- Reality: Configuration was enabled, but feature wasn't wired up
- Evidence: No integration with workspace_modal at that point

**Day 2 (BUILDER):**
- Claimed: "Created modal_wrapper.go with shadow integration"
- Reality: Created effects infrastructure but didn't wire it to workspace_modal
- Evidence: workspace_modal.go line 241 still uses `w.form` directly

**Day 2 (VERIFIER):**
- Claimed: "Code-verified shadows enabled"
- Reality: Verified configuration exists, but didn't test actual rendering
- Missed: The critical integration gap in Show() method
- Methodology flaw: Code review without runtime verification

### Why did code verification miss this?

**Three-level failure:**

1. **Configuration Check:** ✅ Passed (DropShadows: true exists)
2. **Code Existence Check:** ✅ Passed (ShadowForm class exists)
3. **Integration Check:** ❌ Failed (Never verified Show() uses wrapper)

**Critical mistake:** Assumed that having shadow code + config meant shadows work. Never traced the actual call path from `showWorkspaceModal()` → `Show()` → `SetRoot()`.

**What should have been done:**
- Trace complete rendering path
- Verify Primitive() is actually called
- Check which primitive is passed to SetRoot()
- Better yet: Actually run the TUI and visually verify

### How do we prevent this in future?

**Process Changes:**

1. **Visual Regression Tests:**
   - Capture TUI screenshots (using tools like `tcell.SimulationScreen`)
   - Compare against golden images
   - Automate visual verification

2. **Integration Testing:**
   - Don't just verify code exists
   - Trace call paths end-to-end
   - Verify runtime behavior, not just static code

3. **Demo Program Required:**
   - TUIUX agent must create demo program showing feature
   - Builder must run demo before claiming completion
   - Verifier must run demo as part of verification

4. **Acceptance Criteria Must Include:**
   - "User can see [specific visual element]"
   - Not just "Code implements [feature]"
   - Visual proof required for visual features

---

## 6. Additional Findings

### modal_wrapper.go is Unused

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/effects/modal_wrapper.go`

This file exists with sophisticated modal centering and shadow support, but:
- ❌ Never imported by workspace_modal.go
- ❌ Never used by app.go
- ❌ Helper functions like `ShowModal()` and `CenteredForm()` are dead code

**Why it exists:** Day 2 agent created this as a "better" approach but never replaced the existing Show() implementation.

**Recommendation:** Either use modal_wrapper.go OR delete it to avoid confusion. Having two shadow implementations is technical debt.

---

## 7. Fix Validation Checklist

Before claiming shadows are fixed:

- [ ] Apply one-line fix to workspace_modal.go line 241
- [ ] Build the application: `go build ./cmd/ticketr`
- [ ] Run TUI: `./ticketr tui`
- [ ] Open workspace modal (press 'c' or trigger creation)
- [ ] **Visually verify** `▒` shadow characters visible on right and bottom edges
- [ ] Take screenshot as proof
- [ ] Test with `TICKETR_EFFECTS_SHADOWS=false` to verify shadows disappear
- [ ] Test with `TICKETR_EFFECTS_SHADOWS=true` to verify shadows reappear
- [ ] Verify shadow color is visible against terminal background
- [ ] Test on both dark and light terminal themes

---

## 8. Estimated Fix Time

**Time to implement fix:** 5 minutes
**Time to test and verify:** 15 minutes
**Time to update documentation:** 10 minutes
**Total:** ~30 minutes

**Risk Assessment:** LOW
- One-line change
- No API changes
- No new dependencies
- Existing code is correct, just not used
- Rollback is trivial (revert one line)

---

## 9. Next Steps

### Immediate Actions

1. **Apply the fix** (workspace_modal.go line 241)
2. **Visual testing** (run TUI, verify shadows visible)
3. **Screenshot proof** (attach to verification report)
4. **Update test evidence** (VERIFIER day 3 report)

### Follow-up Actions

1. **Decide on modal_wrapper.go fate:**
   - Option A: Delete it (shadows work with ShadowForm)
   - Option B: Refactor to use it (better centering)
   - Recommendation: Delete it (YAGNI - we don't need two implementations)

2. **Add visual regression test:**
   - Create test that captures modal rendering
   - Store golden image
   - Automate shadow verification

3. **Update testing process:**
   - Add "visual proof required" to acceptance criteria
   - Require demo program for all visual features
   - Verifier must run TUI, not just read code

---

## 10. Lessons Learned

### What Went Wrong

1. **Assumed code existence == feature works**
2. **Never traced complete call path**
3. **Code review without runtime verification**
4. **Created duplicate implementations (ShadowForm + modal_wrapper)**
5. **No visual proof required for visual features**

### What Went Right

1. **Shadow implementation is technically correct**
2. **Configuration system works**
3. **Theme system properly loads settings**
4. **Wrapper pattern is sound architecture**

### Key Insight

**Visual features require visual verification.** No amount of code review can substitute for actually running the TUI and looking at it with human eyes.

---

## Appendix: Technical Details

### Shadow Rendering Algorithm

**From shadowbox.go, the shadow is rendered in three parts:**

1. **Bottom shadow:** Horizontal line at `y + height + offsetY - 1`
   - Extends from `x + offsetX` to `x + width + offsetX`
   - Character: `▒` (medium shade)

2. **Right shadow:** Vertical line at `x + width + offsetX - 1`
   - Extends from `y + offsetY` to `y + height + offsetY`
   - Character: `▒` (medium shade)

3. **Corner shadow:** Single character at intersection
   - Position: `(x + width + offsetX - 1, y + height + offsetY - 1)`
   - Character: `▒` (medium shade)

**Visual Representation:**
```
╔════════════════╗
║ Modal Title    ║
║                ║
║ [Form Fields]  ║
║                ║
║  [OK] [Cancel] ║
╚════════════════╝▒
 ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
```

**Default offset:**
- X: 2 characters right
- Y: 1 character down

**Color:** Gray (tcell.ColorGray) with Dim(true) for subtle appearance

---

## Conclusion

**Drop shadows are implemented correctly but not integrated.**

The fix is trivial: Change one line in workspace_modal.go to use `Primitive()` method instead of directly accessing `w.form`.

**The real issue:** Our verification process failed to catch this integration gap because we relied on code review without visual testing. This has been a costly lesson in the importance of end-to-end testing for visual features.

**User is right:** At 4/10 visual quality, the TUI does look "back in the 80s" without the shadows we claimed to implement. This fix will immediately improve perceived quality.

---

**Report compiled by:** TUIUX Agent
**Verification status:** Root cause confirmed with evidence
**Fix confidence:** HIGH - Solution is clear and low-risk
**Ready for implementation:** YES
