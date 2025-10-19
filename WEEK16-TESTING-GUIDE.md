# Week 16 Feature Testing Guide

## Quick Testing Checklist

This guide helps verify all Week 16 features are working correctly.

---

## Prerequisites

```bash
# Build the TUI
go build ./...

# Run the TUI (if you have a workspace configured)
./ticketr tui

# Or create a test workspace first
./ticketr workspace create test-ws --url https://your-jira.atlassian.net --token YOUR_TOKEN --project PROJ
```

---

## 1. Page Navigation Testing

### Ticket Detail View

1. **Open TUI and select a ticket**
   - Press `Tab` until workspace list is focused
   - Press `Enter` on a workspace
   - Press `Tab` to move to ticket tree
   - Press `Enter` on a ticket to view details

2. **Test page navigation**
   - Press `j` multiple times to scroll down (should see scrolling)
   - Press `Ctrl+F` → Should jump down by full page
   - Press `Ctrl+B` → Should jump up by full page
   - Press `Ctrl+D` → Should jump down by half page
   - Press `Ctrl+U` → Should jump up by half page
   - Press `g` → Should jump to top
   - Press `G` (Shift+g) → Should jump to bottom

3. **Verify status bar**
   - At bottom of detail view, should see: "Keys: e edit | j/k scroll | Ctrl+F/B page | Ctrl+D/U half | g/G top/bottom | Esc back"

### Help View

1. **Open help**
   - Press `?` from main view

2. **Test navigation**
   - Press `j`/`k` for single-line scrolling
   - Press `Ctrl+F` for page down
   - Press `Ctrl+B` for page up
   - Press `Ctrl+D` for half-page down
   - Press `Ctrl+U` for half-page up
   - Press `g` to jump to top
   - Press `G` to jump to bottom

3. **Verify documentation**
   - Scroll to "Page Navigation (Week 16 - NEW!)" section
   - Should see all Ctrl+F/B/D/U keybindings documented
   - "About" section should say "Phase 4 Week 16"

### Search View

1. **Open search**
   - Press `/` from main view
   - Type some search text to get results

2. **Navigate to results list**
   - Press Down arrow or `j` to move to results

3. **Test page navigation**
   - Press `Ctrl+F` → Should skip down 10 items
   - Press `Ctrl+B` → Should skip up 10 items
   - Press `Ctrl+D` → Should skip down 5 items
   - Press `Ctrl+U` → Should skip up 5 items

---

## 2. Theme System Testing

### Verify Theme Colors

1. **Check default theme**
   - Focused panels should have **green** borders
   - Unfocused panels should have **white** borders

2. **Verify consistent theming**
   - Tab through all panels (workspace, tree, detail)
   - Green border should follow focus
   - Previous panel should turn white

3. **Check status colors**
   - Sync status bar uses colors:
     - Idle: white
     - Syncing: yellow
     - Success: green
     - Error: red

### Test Theme Switching (Programmatic)

**Note:** Theme switching UI not yet implemented, but theme system is ready.

To test programmatically, modify `app.go` before building:
```go
// In setupApp() after theme.Apply(t.app), add:
theme.SetByName("dark")
theme.Apply(t.app)
```

Build and run - should see blue borders instead of green.

---

## 3. Error Modal Testing

**Note:** Error modal is created but not yet integrated into all error paths.

To test, you can trigger errors:
- Try to sync without credentials configured
- Try to switch to non-existent workspace

Error display should be clear and actionable.

---

## 4. Enhanced Status Bar Testing

### Check Status Bar Content

1. **View sync status**
   - Bottom right panel shows "Sync Status"
   - Should show current state with icon:
     - ○ Idle
     - ◌ Syncing
     - ● Success
     - ✗ Error

2. **Trigger sync operations**
   - Press `p` to push (may fail if not configured - that's OK)
   - Watch status bar update in real-time
   - Press `r` to refresh
   - Status should change from idle → syncing → success/error

3. **Verify workspace info**
   - Status bar should show sync state
   - In future compact mode, workspace name would appear here

---

## 5. Help Documentation Testing

### Verify Help Content

1. **Open help** (`?`)

2. **Check new sections**
   - [ ] "Page Navigation (Week 16 - NEW!)" section exists
   - [ ] Ctrl+F/B/D/U documented
   - [ ] "Responsive Layout (Week 16 - NEW!)" section exists
   - [ ] "Theme System (Week 16 - NEW!)" section exists
   - [ ] "Performance (Week 16)" section exists
   - [ ] "About" says "Phase 4 Week 16"

3. **Verify accuracy**
   - [ ] All keybindings work as documented
   - [ ] No outdated information
   - [ ] Tips section includes new features

---

## 6. Responsive Layout Testing (Framework Only)

**Note:** Dynamic layout switching not yet implemented, but framework is in place.

### Verify No Crashes

1. **Test normal terminal**
   - Run TUI in terminal >= 100 columns wide
   - Should see tri-panel layout (workspace | tree | detail)

2. **Test small terminal (if possible)**
   - Resize terminal to 80 columns
   - TUI should still work (no crash)
   - Full layout still shown (compact mode not yet active)

3. **Test very small terminal**
   - Resize to 50 columns
   - TUI should still attempt to render (may look cramped but no crash)

---

## 7. Integration Testing

### End-to-End User Flow

1. **Navigate between panels**
   - Tab through all panels
   - Theme colors should update correctly

2. **View ticket details**
   - Select ticket from tree
   - Use all navigation keys (j/k/Ctrl+F/B/D/U/g/G)
   - All should work smoothly

3. **Search tickets**
   - Press `/`
   - Type search query
   - Navigate results with Ctrl+F/B
   - Select ticket, should open detail view

4. **View help**
   - Press `?`
   - Page through help with Ctrl+F/B
   - All shortcuts documented

5. **Sync operations**
   - Press `p` or `P` or `s`
   - Watch status bar update
   - UI should remain responsive (async)

---

## 8. Performance Testing

### With Large Ticket Count

1. **Load 100+ tickets**
   - Switch to workspace with many tickets
   - Tree should render quickly (< 100ms)

2. **Test search**
   - Open search (`/`)
   - Type query
   - Results should appear instantly (< 50ms)

3. **Test navigation**
   - Page through ticket detail with Ctrl+F/B
   - Should be instant (< 10ms)

4. **Memory usage**
   - Run TUI for several minutes
   - Perform various operations
   - Should not leak memory (use `top` to monitor)

---

## 9. Edge Case Testing

### Boundary Conditions

1. **Scroll to top**
   - In ticket detail, press `g` (go to top)
   - Press `Ctrl+B` (page up)
   - Should stay at top, no crash

2. **Scroll to bottom**
   - Press `G` (go to bottom)
   - Press `Ctrl+F` (page down)
   - Should stay at bottom, no crash

3. **Empty results**
   - Search for non-existent text
   - Press Ctrl+F/B in empty list
   - Should not crash

4. **Single ticket**
   - If only one ticket in workspace
   - Navigation should still work

---

## 10. Regression Testing

### Verify Old Features Still Work

1. **Basic navigation** (Week 11-12)
   - [ ] Tab cycles focus
   - [ ] Esc goes back
   - [ ] q quits

2. **Ticket editing** (Week 13)
   - [ ] Press `e` in detail view
   - [ ] Edit fields
   - [ ] Save changes
   - [ ] Validation works

3. **Search & filters** (Week 14)
   - [ ] `/` opens search
   - [ ] `@user` filter works
   - [ ] `#ID` filter works
   - [ ] `:` opens command palette

4. **Async sync** (Week 15)
   - [ ] `p` push (async)
   - [ ] `P` pull (async)
   - [ ] `s` sync (async)
   - [ ] UI remains responsive during sync

---

## Known Issues & Limitations

### Expected Behaviors

1. **Terminal size detection**
   - Currently uses full layout regardless of terminal size
   - No dynamic switching to compact mode yet
   - Framework in place for future enhancement

2. **Theme switching UI**
   - Theme system works programmatically
   - No command palette command yet
   - Can be added in Phase 5

3. **Error modal integration**
   - ErrorModal created but not yet used everywhere
   - Some errors still show in status bar
   - Can be improved incrementally

### Pre-existing Issues

1. **Keychain concurrent test**
   - One flaky test in keychain package
   - Not introduced in Week 16
   - Does not affect TUI functionality

---

## Success Criteria

All tests should pass with:
- ✅ No crashes
- ✅ Smooth navigation
- ✅ Consistent theme colors
- ✅ Accurate help documentation
- ✅ Responsive UI during async operations
- ✅ All keybindings work as expected

---

## Reporting Issues

If you find issues during testing:

1. Note the exact steps to reproduce
2. Screenshot if visual issue
3. Check if regression (old feature broken) or new bug
4. Report to Verifier with:
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - Terminal size and OS

---

## Quick Test Commands

```bash
# Build
go build ./...

# Run tests
go test ./...

# Test specific package
go test ./internal/adapters/tui/search/... -v

# Run TUI
./ticketr tui

# Check version/build info
./ticketr version
```

---

## Testing Checklist Summary

**Page Navigation:**
- [ ] Ctrl+F/B work in ticket detail
- [ ] Ctrl+D/U work in ticket detail
- [ ] Ctrl+F/B work in help view
- [ ] Ctrl+F/B work in search results
- [ ] Status bars show new keys

**Theme System:**
- [ ] Green borders on focus
- [ ] White borders unfocused
- [ ] Consistent across all views
- [ ] Status colors correct

**Help Documentation:**
- [ ] Week 16 sections present
- [ ] All shortcuts documented
- [ ] About updated
- [ ] Tips accurate

**Status Bar:**
- [ ] Shows sync state
- [ ] Updates in real-time
- [ ] Icons display correctly

**Performance:**
- [ ] Fast with 100+ tickets
- [ ] Instant navigation
- [ ] No lag

**Regression:**
- [ ] All Week 11-15 features still work
- [ ] No existing functionality broken

---

**Happy Testing! 🧪**

If all tests pass, Phase 4 Week 16 is verified and ready for Phase 5!
