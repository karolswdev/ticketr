# ROOT CAUSE ANALYSIS: Workspace Panel/Modal State Corruption Bug

## Bug Reproduction Flow

### Working Flow (before ESC corruption):
1. Press `W` (uppercase) → `globalKeyHandler` line 465 → `toggleWorkspacePanel()` → Panel shows ✅
2. Press `W` → `toggleWorkspacePanel()` → Panel closes ✅
3. Press `W` → `toggleWorkspacePanel()` → Panel shows ✅
4. Press `w` (lowercase) → Panel is visible, so `WorkspaceListView.handleInput()` line 132-137 → calls `onCreateWorkspace()` → Modal appears ✅
5. Press `ESC` → Modal ESC handler line 73-79 → calls `handleCancel()` → Modal closes ✅
6. **Press `W` → Panel BROKEN ❌** (this is where it fails)
7. **Press `w` → Modal BROKEN ❌**

## The State Corruption

### Event Routing Analysis

**Key Insight:** There are TWO different input capture mechanisms:

1. **Global Handler** (`app.go` line 368-507): `app.SetInputCapture(t.globalKeyHandler)`
   - Handles `W` (uppercase) at line 465-468
   - Handles `w` (lowercase) ONLY when `!t.inModal` (line 398)

2. **Panel-Local Handler** (`workspace_list.go` line 39): `list.SetInputCapture(view.handleInput)`
   - Handles `w` (lowercase) at line 132-137 when panel has focus
   - Returns `nil` to consume the event

### The Corruption Sequence

**After ESC closes the modal:**

1. Line 441 `workspace_modal.go`: ESC triggers `handleCancel()`
2. Line 442: Calls `w.onClose()` callback
3. Line 312 `app.go`: Sets `t.inModal = false` ✅ (correct)
4. Line 314: Calls `t.pages.HidePage("workspace-modal")` ✅ (correct)
5. Line 316: Calls `t.updateFocus()`

**The CRITICAL ISSUE in `updateFocus()`:**

Looking at `app.go` lines 538-568:

```go
func (t *TUIApp) updateFocus() {
	// Update border colors for all views
	t.ticketTreeView.SetFocused(t.currentFocus == "ticket_tree")
	t.ticketDetailView.SetFocused(t.currentFocus == "ticket_detail")

	// ... action bar update ...

	// Set application focus
	switch t.currentFocus {
	case "ticket_tree":
		t.app.SetFocus(t.ticketTreeView.Primitive())
	case "ticket_detail":
		t.app.SetFocus(t.ticketDetailView.Primitive())
	}
}
```

**THE BUG:** When the modal closes, `t.currentFocus` is still set to whatever it was before the modal opened. The modal ESC handler does NOT track or restore focus properly.

## Root Cause Identified

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go`
**Lines:** 310-317 (WorkspaceModal OnClose handler)

**The Problem:**

```go
t.workspaceModal.SetOnClose(func() {
	// FIX #1: Hide modal page instead of SetRoot
	t.inModal = false
	if t.pages != nil {
		t.pages.HidePage("workspace-modal")
	}
	t.updateFocus()  // ❌ THIS CALLS SetFocus ON MAIN VIEW, NOT PANEL
})
```

When the modal was opened from the workspace panel (pressing `w` while panel is visible):
- The panel still has the `"workspace-overlay"` page showing
- But `updateFocus()` calls `SetFocus()` on `ticketTreeView` or `ticketDetailView`
- This leaves the workspace panel VISIBLE but WITHOUT FOCUS

**After this corruption:**

1. Press `W` → `globalKeyHandler` line 465 → calls `toggleWorkspacePanel()`
2. Line 1042: Checks `t.workspaceSlideOut.IsVisible()`
3. **Returns TRUE** because panel page is still showing!
4. Line 1044-1047: Tries to HIDE the panel
5. Line 1045: `t.pages.HidePage("workspace-overlay")`
6. Line 1047: `t.updateFocus()` → focuses main view

Now panel is HIDDEN. Pressing `W` again:

1. Line 1042: Checks `t.workspaceSlideOut.IsVisible()`
2. **Returns FALSE** because we just hid it
3. Line 1050-1056: Tries to SHOW the panel
4. Line 1051: `t.workspaceSlideOut.Show()` ← **This is a widget method**
5. Line 1052: `t.pages.ShowPage("workspace-overlay")` ✅
6. Line 1054: `t.workspaceListView.SetFocused(true)` ✅
7. Line 1055: `t.app.SetFocus(t.workspaceListView.Primitive())` ✅

**But wait...** Let me check what `workspaceSlideOut.Show()` does...

## Critical Missing Information

I need to check the SlideOut widget implementation to understand if there's state inside it.

## The ACTUAL Root Cause - CONFIRMED

After analyzing SlideOut widget (`widgets/slideout.go`), I found the issue:

**State Desynchronization Between:**

1. **SlideOut.isVisible** (line 15 in `slideout.go`) - Widget's internal state
2. **Pages visibility** ("workspace-overlay" page in `app.go` line 370)

**The Bug Flow:**

1. User presses `W` → Panel shows → `workspaceSlideOut.isVisible = true` ✅
2. User presses `w` (lowercase) → Modal opens
3. Modal opens BUT does NOT close the panel → Panel still visible with `isVisible = true`
4. User presses `ESC` in modal → Modal closes → Calls `updateFocus()` on main view
5. **Panel is still showing** but now `app.SetFocus()` points to `ticketTreeView`
6. User presses `W`:
   - Line 1042 checks `t.workspaceSlideOut.IsVisible()` → Returns `true` (panel thinks it's visible)
   - Line 1044-1047: Tries to HIDE the panel
   - Line 1045: `t.pages.HidePage("workspace-overlay")` → Panel disappears
   - Panel state now: `isVisible = true` (widget state) but page is hidden (pages state)
7. User presses `W` again:
   - Line 1042 checks `t.workspaceSlideOut.IsVisible()` → Returns `true` (still!)
   - Tries to HIDE again, but page already hidden
   - **Toggle is now reversed** - W hides when it should show, shows when it should hide

**The Root Problem:** The SlideOut widget's `isVisible` flag gets out of sync with the pages overlay state when the modal opens without closing the panel first.

## State Variables Involved

1. `t.inModal` (bool) - Tracks if modal is active ✅ Works correctly
2. `t.currentFocus` (string) - Tracks which main view has focus ✅ Works correctly
3. **SlideOut.isVisible** (bool) - Widget's visibility state ❌ **GETS CORRUPTED**
4. **Pages overlay state** ("workspace-overlay" page visibility) ❌ **GETS CORRUPTED**

## The Fix - Option A (RECOMMENDED)

**Auto-close panel before modal opens** to maintain consistent state.

Modify `app.go` line 1010-1015 in `showWorkspaceModal()`:

```go
func (t *TUIApp) showWorkspaceModal() {
	// BUGFIX: Auto-close workspace panel if open to prevent state corruption
	if t.workspaceSlideOut != nil && t.workspaceSlideOut.IsVisible() {
		t.toggleWorkspacePanel() // This properly syncs both isVisible and pages state
	}
	t.inModal = true
	t.workspaceModal.OnShow()
	t.workspaceModal.Show()
}
```

This ensures:
- Panel is properly closed (both `isVisible` and pages state synced)
- Modal opens cleanly
- ESC closes modal and returns to main view
- Next `W` works correctly (panel state is clean)

## Verification

After fix, this sequence should work:
1. Press `W` → Panel shows ✅
2. Press `w` → Modal opens (panel auto-closes) ✅
3. Press `ESC` → Modal closes, focus to main view ✅
4. Press `W` → Panel shows ✅
5. Press `ESC` → Panel closes ✅
6. Press `w` → Modal opens ✅
7. Repeat...

## Implementation Plan

**Option A: Auto-close panel before modal (RECOMMENDED)**

Modify `app.go` line 1010-1015:

```go
func (t *TUIApp) showWorkspaceModal() {
	// Auto-close workspace panel if open
	if t.workspaceSlideOut != nil && t.workspaceSlideOut.IsVisible() {
		t.toggleWorkspacePanel()
	}
	t.inModal = true
	t.workspaceModal.OnShow()
	t.workspaceModal.Show()
}
```

This ensures clean state when modal opens/closes.

**Option B: Track modal parent and restore (COMPLEX)**

Add field to TUIApp:
```go
modalOpenedFrom string // "", "panel", "main"
```

Track when opening, restore when closing. More complex, harder to maintain.

## Decision: Option A

Implement the simple fix - auto-close panel before modal opens.
