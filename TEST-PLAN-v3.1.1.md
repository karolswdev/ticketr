# Ticketr v3.1.1 - Human Validation Test Plan

**Build:** `./ticketr` (current directory)
**Focus:** Phase 6 Week 2 TUI UX Improvements
**Estimated Time:** 30-45 minutes

---

## Prerequisites

Before testing, ensure you have:
- [ ] Built binary: `./ticketr` exists in current directory
- [ ] Valid Jira credentials configured
- [ ] At least one workspace configured (or ready to create one)
- [ ] Access to a Jira instance with some tickets

**Quick Setup Verification:**
```bash
# Verify binary exists
ls -lh ./ticketr

# Test basic command (should show help)
./ticketr --help

# Check if you have workspaces configured
./ticketr workspace list
```

---

## Test Suite: Phase 6 Week 2 Features

### TEST 1: Launch TUI (Baseline)
**Feature:** TUI Menu Structure (Day 8-9)
**Priority:** CRITICAL

**Steps:**
1. Launch TUI: `./ticketr tui`
2. Observe the interface loads without errors
3. Verify you see the main ticket list view
4. Look for the **bottom action bar** (should show context-aware actions)

**Expected Results:**
- ✅ TUI launches without crashes
- ✅ Bottom action bar visible with keybindings
- ✅ Interface responsive to keyboard input

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 2: Action Bar - Main View Context
**Feature:** Context-Aware Action Bar (Day 8-9)
**Priority:** HIGH

**Steps:**
1. In TUI main view, look at bottom action bar
2. Verify you see keys like: `F1`, `F2`, `F5`, `F10`, `ESC`
3. Observe the action labels (e.g., "F1:Help", "F2:Workspaces", "F5:Sync", "F10:Quit")

**Expected Results:**
- ✅ Action bar shows 5+ actions
- ✅ Actions are context-appropriate for main view
- ✅ Keybindings clearly labeled

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 3: F-Key Shortcuts - F1 (Help)
**Feature:** F-Key Shortcuts (Day 8-9)
**Priority:** HIGH

**Steps:**
1. Press `F1` to open help/keybindings
2. Observe a modal or view showing keybindings
3. Press `ESC` to close

**Expected Results:**
- ✅ F1 opens help view
- ✅ Help content is readable and accurate
- ✅ ESC closes help view

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 4: F-Key Shortcuts - F2 (Workspaces)
**Feature:** Workspace Modal (Day 8-9)
**Priority:** HIGH

**Steps:**
1. Press `F2` to open workspace selector
2. Observe workspace list (or prompt to create one)
3. Use arrow keys to navigate workspaces
4. Press `ESC` to close without selecting

**Expected Results:**
- ✅ F2 opens workspace modal
- ✅ Workspaces listed (if any exist)
- ✅ Navigation works
- ✅ ESC closes modal

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 5: Command Palette - Invoke
**Feature:** Command Palette with Fuzzy Search (Day 8-9)
**Priority:** HIGH

**Steps:**
1. Press `Ctrl+P` or `:` to open command palette
2. Observe text input field appears
3. Type a few characters (e.g., "sync" or "help")
4. Observe fuzzy matching filters commands
5. Use arrow keys to select a command
6. Press `ESC` to cancel

**Expected Results:**
- ✅ Command palette opens
- ✅ Fuzzy search filters commands as you type
- ✅ Commands are selectable
- ✅ ESC cancels without executing

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 6: Async Pull Operation (No Progress - Baseline)
**Feature:** Async Job Queue (Day 6-7) - Baseline test
**Priority:** CRITICAL

**Steps:**
1. Ensure you have a workspace configured with Jira tickets
2. Press `F5` (Sync) or select "Sync All" command
3. Observe that the TUI does NOT freeze
4. Try pressing keys (arrow up/down) during sync
5. Wait for sync to complete

**Expected Results:**
- ✅ TUI remains responsive during sync (no freeze)
- ✅ You can navigate tickets while sync runs in background
- ✅ Sync completes successfully

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 7: Progress Indicators - Visual Display
**Feature:** Real-Time Progress Indicators (Day 10-11)
**Priority:** CRITICAL

**Steps:**
1. Start a sync operation (F5 or pull command)
2. Look at the **status bar** (likely at top or bottom)
3. Observe the progress bar during sync
4. Check for these elements:
   - Visual bar: `[█████░░░░░]` or similar
   - Percentage: `50%` or similar
   - Counts: `(25/50)` or similar
   - Time: `Elapsed: 5s` or similar
   - ETA: `ETA: 5s` or similar

**Expected Results:**
- ✅ Progress bar appears during sync
- ✅ Percentage displayed and updates
- ✅ Ticket counts shown (current/total)
- ✅ Elapsed time shown
- ✅ ETA calculated and shown
- ✅ Progress bar updates in real-time (not stuck)

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 8: Progress Indicators - Accuracy
**Feature:** Progress Calculation (Day 10-11)
**Priority:** HIGH

**Steps:**
1. During a sync, observe the progress percentage
2. Verify percentage matches count (e.g., 25/100 = 25%)
3. Watch the ETA - does it seem reasonable?
4. Verify progress reaches 100% when complete

**Expected Results:**
- ✅ Percentage calculation accurate
- ✅ ETA reasonable (not wildly off)
- ✅ Progress reaches 100% on completion
- ✅ Progress bar clears after operation

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 9: Async Cancellation
**Feature:** Graceful Job Cancellation (Day 6-7)
**Priority:** HIGH

**Steps:**
1. Start a sync operation (F5)
2. While sync is running, press `ESC`
3. Observe the sync operation cancels
4. Verify TUI returns to normal state
5. Check that partial results preserved (if applicable)

**Expected Results:**
- ✅ ESC cancels the running operation
- ✅ TUI returns to normal state
- ✅ No errors or crashes
- ✅ Partial progress preserved (tickets pulled before cancel)

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 10: Visual Effects - Default (OFF)
**Feature:** Visual Effects System (Day 12.5)
**Priority:** MEDIUM

**Steps:**
1. Launch TUI normally: `./ticketr tui`
2. Observe the visual presentation
3. Look for effects (you should NOT see them by default):
   - Drop shadows on modals
   - Shimmer/animation on progress bars
   - Background particle effects (stars, snow)

**Expected Results:**
- ✅ TUI looks clean and professional
- ✅ NO visual effects enabled by default
- ✅ No animations or drop shadows
- ✅ Fast and responsive

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 11: Visual Effects - Motion Enabled
**Feature:** Visual Effects with Motion (Day 12.5)
**Priority:** LOW

**Steps:**
1. Exit TUI if running
2. Enable motion effects:
   ```bash
   export TICKETR_EFFECTS_MOTION=true
   export TICKETR_EFFECTS_SHADOWS=true
   export TICKETR_EFFECTS_SHIMMER=true
   ./ticketr tui
   ```
3. Open workspace modal (F2)
4. Observe if modal has drop shadow (▒ characters around edges)
5. Start a sync operation
6. Observe if progress bar has shimmer animation

**Expected Results:**
- ✅ Workspace modal shows drop shadow
- ✅ Progress bar has shimmer effect during sync
- ✅ No performance degradation
- ✅ TUI still responsive

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 12: Visual Effects - Dark Theme
**Feature:** Themed Visual Effects (Day 12.5)
**Priority:** LOW

**Steps:**
1. Exit TUI if running
2. Enable dark theme:
   ```bash
   export TICKETR_THEME=dark
   export TICKETR_EFFECTS_MOTION=true
   ./ticketr tui
   ```
3. Observe the color scheme changes
4. Look for themed elements (darker background, different borders)

**Expected Results:**
- ✅ Theme changes to dark palette
- ✅ Interface remains readable
- ✅ No rendering glitches

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 13: Regression - Basic Operations Still Work
**Feature:** No regressions from Phase 6 changes
**Priority:** CRITICAL

**Steps:**
1. Launch TUI: `./ticketr tui`
2. Navigate ticket list with arrow keys
3. Select a ticket and view details (Enter)
4. Edit a ticket (if supported)
5. Exit cleanly (F10 or Ctrl+C)

**Expected Results:**
- ✅ Arrow key navigation works
- ✅ Ticket selection works
- ✅ Ticket viewing works
- ✅ Clean exit (no errors)

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 14: Command Palette - Execute Command
**Feature:** Command Execution (Day 8-9)
**Priority:** HIGH

**Steps:**
1. Open command palette (Ctrl+P or `:`)
2. Type "help" and press Enter
3. Verify help command executes
4. Open palette again
5. Type "workspace" and press Enter
6. Verify workspace modal opens

**Expected Results:**
- ✅ Command palette executes selected commands
- ✅ Commands have correct behavior
- ✅ Multiple invocations work

**Pass/Fail:** ______

**Notes:**
```


```

---

### TEST 15: Performance - Large Workspace (if available)
**Feature:** Performance with many tickets
**Priority:** MEDIUM

**Steps:**
1. If you have a workspace with 50+ tickets, use it
2. Launch TUI and select that workspace
3. Scroll through the ticket list (rapid arrow key presses)
4. Start a sync operation
5. Observe responsiveness

**Expected Results:**
- ✅ List scrolling smooth (no lag)
- ✅ Sync operation doesn't freeze UI
- ✅ Memory usage reasonable

**Pass/Fail:** ______

**Notes:**
```


```

---

## Test Results Summary

**Date:** ____________
**Tester:** ____________
**Build:** `./ticketr`

**Critical Tests (MUST PASS):**
- [ ] TEST 1: Launch TUI
- [ ] TEST 6: Async Pull Operation
- [ ] TEST 7: Progress Indicators - Visual Display
- [ ] TEST 9: Async Cancellation
- [ ] TEST 13: No Regressions

**High Priority Tests:**
- [ ] TEST 2: Action Bar
- [ ] TEST 3: F1 Help
- [ ] TEST 4: F2 Workspaces
- [ ] TEST 5: Command Palette
- [ ] TEST 8: Progress Accuracy
- [ ] TEST 14: Command Execution

**Medium/Low Priority Tests:**
- [ ] TEST 10: Visual Effects Default (OFF)
- [ ] TEST 11: Visual Effects Motion
- [ ] TEST 12: Dark Theme
- [ ] TEST 15: Performance

---

## Issues Found

**Blocker Issues (Release Blockers):**
```
1.
2.
3.
```

**High Priority Issues (Should Fix):**
```
1.
2.
3.
```

**Low Priority Issues (Can Defer):**
```
1.
2.
3.
```

---

## Final Decision

**GO FOR RELEASE:** YES / NO / CONDITIONAL

**Conditions (if applicable):**
```


```

**Tester Signature:** ____________
**Date:** ____________

---

## Quick Test Commands

```bash
# Build fresh binary
go build -o ticketr ./cmd/ticketr

# Basic launch
./ticketr tui

# With visual effects (full enchantment)
TICKETR_THEME=dark TICKETR_EFFECTS_MOTION=true TICKETR_EFFECTS_SHADOWS=true TICKETR_EFFECTS_SHIMMER=true ./ticketr tui

# Reset to defaults
unset TICKETR_THEME TICKETR_EFFECTS_MOTION TICKETR_EFFECTS_SHADOWS TICKETR_EFFECTS_SHIMMER
./ticketr tui
```

---

## What to Look For

**GOOD Signs:**
- TUI launches instantly
- No freezing during sync operations
- Progress bar shows real-time updates
- Action bar changes based on context
- F-keys work as labeled
- ESC cancels operations cleanly
- Visual effects are OFF by default (fast, accessible)

**BAD Signs (Report Immediately):**
- TUI crashes on launch
- UI freezes during sync
- Progress bar stuck at 0%
- F-keys don't work
- Cancellation doesn't work (ESC has no effect)
- Visual glitches or garbled text
- High CPU usage when idle
- Memory leaks (usage growing constantly)

---

**Good luck testing! Report any issues before we proceed to Day 15 Release.**
