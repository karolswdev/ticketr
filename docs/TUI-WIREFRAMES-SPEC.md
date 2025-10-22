# Ticketr TUI - Visual Design Specification & Wireframes

**Version:** v3.1.1 Target
**Purpose:** Define the visual character, animations, and UX expectations
**Created:** 2025-10-20 (Phase 6.5)
**Status:** AUTHORITATIVE - This is what we're building toward

---

## Design Philosophy

### The Character We're Building

**Not This:** Raw tview widgets, static text, "back in the 80s"
**Instead:** Smooth, modern TUI with subtle animations, visual depth, professional polish

**Inspiration:**
- lazygit (smooth, responsive, modern)
- k9s (polished, animated, professional)
- lazydocker (clean, themed, delightful)

**Target FPS:** 30-60 FPS for animations (smooth, not choppy)
**Performance Budget:** ≤5% CPU when idle, ≤15% during animations

---

## Visual Hierarchy & Layout

### Main View (Tri-Panel Layout)

```
┌─ Ticketr ─────────────────────────────────────────────────────────────┐
│                                                                        │
│ ┌─ Workspaces ────┐ ┌─ Tickets ─────────┐ ┌─ Details ──────────────┐ │
│ │                 │ │                    │ │                         │ │
│ │ ► tbct          │ │ ◉ EPM-123         │ │ # EPM-123: Fix Bug     │ │
│ │   project-alpha │ │   EPM-456         │ │                         │ │
│ │   staging       │ │   EPM-789         │ │ Status: In Progress    │ │
│ │                 │ │                    │ │ Assignee: John Doe     │ │
│ │                 │ │ Loading...         │ │                         │ │
│ │                 │ │ ⠋ Fetching        │ │ Description:           │ │
│ │                 │ │                    │ │ Lorem ipsum...         │ │
│ └─────────────────┘ └────────────────────┘ └─────────────────────────┘ │
│                                                                        │
│ ┌─ Status ──────────────────────────────────────────────────────────┐ │
│ │ ✓ Synced 42 tickets • Last pull: 2m ago                           │ │
│ │ [████████████████████░░░░] 80% (40/50) • 5s elapsed • ETA: 2s     │ │
│ └────────────────────────────────────────────────────────────────────┘ │
│                                                                        │
│ ┌─ Actions ──────────────────────────────────────────────────────────┐│
│ │ F1:Help  F2:Workspaces  F5:Sync  P:Pull  p:Push  ?:Keys  F10:Quit ││
│ └────────────────────────────────────────────────────────────────────┘│
└────────────────────────────────────────────────────────────────────────┘
```

**Key Visual Elements:**
- **Borders:** Double-line for focused (╔═╗), single-line for unfocused (┌─┐)
- **Selection:** ► or ◉ indicator with subtle highlight
- **Spinner:** Rotating Braille characters (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏) at 80ms intervals
- **Progress:** Block characters [████░░░░] updating smoothly
- **Status:** Icons (✓ ✗ ⚠) with color coding

---

## Animation Specifications

### 1. Spinner Animation (CRITICAL)

**Location:** Status bar, loading states, sync operations
**Character Set:** ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏ (Braille spinner)
**Rotation Speed:** 80ms per frame (12.5 FPS)
**Behavior:**
- MUST animate continuously when visible
- Stops when operation complete
- Never static

**Visual Example:**
```
Frame 1: ⠋ Fetching tickets...
Frame 2: ⠙ Fetching tickets...
Frame 3: ⠹ Fetching tickets...
Frame 4: ⠸ Fetching tickets...
... (cycles continuously)
```

**Implementation:**
- Animator goroutine updates every 80ms
- Uses QueueUpdateDraw() to refresh UI
- Cancellable via context

---

### 2. Progress Bar Animation

**Location:** Status bar during sync operations
**Visual:**
```
[████████████████████░░░░░░░░░░] 65% (65/100) | 12s | ETA: 8s
 └──────────────────┘ └────────┘   └──┘ └───┘    └─┘    └──┘
 Filled (animated)    Empty      Percent Count  Time   ETA
```

**Animations:**
1. **Fill Progress:** Smooth left-to-right fill (not jumpy)
2. **Shimmer Effect:** Subtle sweep across filled portion every 2s
3. **Percentage Update:** Real-time calculation
4. **Time Updates:** Every 500ms

**Character Set:**
- Filled: █ (full block)
- Empty: ░ (light shade)
- Shimmer: Brief flash of ▓ (medium shade) sweeping right

**FPS Target:** 30 FPS for shimmer animation

---

### 3. Modal Fade-In

**Location:** All modals (workspace creation, command palette, etc.)
**Animation:** 150ms fade from ░ → ▒ → █
**Frames:**
```
Frame 1 (0ms):   Background ░ (25% opacity equivalent)
Frame 2 (50ms):  Background ▒ (50% opacity)
Frame 3 (100ms): Background ▓ (75% opacity)
Frame 4 (150ms): Background █ (100% solid)
```

**Modal Shadow:**
```
┌─ Create Workspace ──────┐▒
│                          │▒
│ Name: [____________]     │▒
│                          │▒
│ [Cancel]  [Create]       │▒
└──────────────────────────┘▒
 ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
```
- Shadow offset: 1 row down, 2 cols right
- Shadow character: ▒ (medium shade)

---

### 4. Focus Pulse (Subtle)

**Location:** Focused panel border
**Animation:** Subtle brightness pulse every 2s
**Effect:** Border slightly brighter/dimmer to indicate focus
**FPS:** 15 FPS (subtle, not distracting)

**Visual:**
```
Frame 1: ╔═══════╗ (100% brightness)
Frame 2: ╔═══════╗ (110% brightness) ← subtle pulse
Frame 3: ╔═══════╗ (100% brightness)
... (cycles every 2s)
```

---

### 5. Background Effects (OPTIONAL - Default OFF)

**Hyperspace Theme:**
```
Terminal background with stars moving left to right
★  →  →  →  ✦  →  →  ✧  →  →
   →  ✦  →  →  ★  →  →  →  ✦
```
- Character set: ✦ ✧ ★
- Speed: Slow (20-30 chars/sec)
- Density: 1-2% of screen
- FPS: 12-15 FPS
- CPU: ≤3% target

**Snow Theme:**
```
Snowflakes falling from top to bottom
  ❄     ❅       ❄
    ❅       ❄      ❅
       ❄         ❅
```
- Character set: ❄ ❅
- Speed: Gentle (10-15 chars/sec)
- Density: 1-2% of screen
- FPS: 12-15 FPS
- CPU: ≤3% target

---

## Wireframe: Pull Operation (Animated)

### State 1: Idle (Before Pull)
```
┌─ Status ────────────────────────────────┐
│ ✓ Synced 42 tickets • Last: 2m ago      │
└──────────────────────────────────────────┘
┌─ Actions ────────────────────────────────┐
│ F1:Help  F5:Sync  P:Pull  F10:Quit      │
└──────────────────────────────────────────┘
```

### State 2: User Presses 'P' → Connecting (0-500ms)
```
┌─ Status ────────────────────────────────┐
│ ⠋ Connecting to Jira...                 │  ← SPINNER ANIMATING
└──────────────────────────────────────────┘
┌─ Actions ────────────────────────────────┐
│ ESC:Cancel  (Pull in progress...)        │  ← Context changes
└──────────────────────────────────────────┘
```

### State 3: Querying (500ms-2s)
```
┌─ Status ────────────────────────────────┐
│ ⠙ Querying project EPM...               │  ← SPINNER ANIMATING
└──────────────────────────────────────────┘
┌─ Actions ────────────────────────────────┐
│ ESC:Cancel  (Pull in progress...)        │
└──────────────────────────────────────────┘
```

### State 4: Fetching Tickets (2s-30s)
```
┌─ Status ────────────────────────────────────────────┐
│ ⠹ Fetching tickets...                               │  ← SPINNER
│ [████████░░░░░░░░░░] 40% (20/50) | 8s | ETA: 12s   │  ← PROGRESS
└──────────────────────────────────────────────────────┘
┌─ Actions ────────────────────────────────────────────┐
│ ESC:Cancel  (Pull in progress...)                    │
└──────────────────────────────────────────────────────┘
```
**Animations Active:**
- Spinner rotating every 80ms
- Progress bar filling smoothly
- Percentage updating in real-time
- Time/ETA updating every 500ms
- Shimmer sweeping across filled portion

### State 5: Complete (30s+)
```
┌─ Status ────────────────────────────────────────────┐
│ ✓ Pull complete: 20 new, 15 updated, 15 skipped    │
│ [████████████████████████] 100% (50/50) | 30s      │
└──────────────────────────────────────────────────────┘
┌─ Actions ────────────────────────────────────────────┐
│ F1:Help  F5:Sync  P:Pull  F10:Quit                  │
└──────────────────────────────────────────────────────┘
```
**Animation:**
- Progress bar briefly shows 100%
- After 3s, bar fades out
- Status message remains

**CRITICAL:** UI must remain responsive throughout ALL states
- User can Tab between panes
- User can press ESC to cancel
- User can navigate tickets (even while pull running)

---

## Wireframe: Modal Design (Workspace Creation)

### Modern Design (Target)
```
                    ╔══════════════════════════════╗▒
                    ║  Create Workspace            ║▒
                    ╠══════════════════════════════╣▒
                    ║                              ║▒
                    ║  Workspace Name:             ║▒
                    ║  ┌────────────────────────┐  ║▒
                    ║  │ tbct                   │  ║▒  ← Input focused
                    ║  └────────────────────────┘  ║▒
                    ║                              ║▒
                    ║  Project Key (Jira):         ║▒
                    ║  ┌────────────────────────┐  ║▒
                    ║  │ EPM                    │  ║▒
                    ║  └────────────────────────┘  ║▒
                    ║                              ║▒
                    ║  Credential Profile:         ║▒
                    ║  ┌────────────────────────┐  ║▒
                    ║  │ ► default              │  ║▒
                    ║  │   production           │  ║▒
                    ║  └────────────────────────┘  ║▒
                    ║                              ║▒
                    ║     [Cancel]    [Create]     ║▒
                    ║                              ║▒
                    ╚══════════════════════════════╝▒
                     ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
```

**Visual Features:**
- Double-line border (╔═╗║╚╝)
- Drop shadow (▒ characters offset)
- Centered on screen
- Fade-in animation on open
- Themed colors applied
- Input fields with single-line boxes
- Button highlighting on focus

**vs. "Back in the 80s" (Current - BAD):**
```
+------------------------------+
| Create Workspace             |
+------------------------------+
|                              |
| Name: ____________           |
|                              |
| Project: ____________        |
|                              |
| [Cancel]  [Create]           |
+------------------------------+
```
❌ ASCII borders only
❌ No shadow
❌ No theming
❌ No animation
❌ Looks like 1985

---

## Wireframe: Keybindings Overflow Solution

### Problem (Current - BAD):
```
┌─ Actions ──────────────────────────────────────────────┐
│ F1:Help F2:Workspaces F5:Sync P:Pull p:Push s:FullSync │  ← Cuts off
```
Text overflows, user can't see all actions.

### Solution 1: Marquee (Preferred for ≤80 cols)
```
Frame 1 (0ms):
┌─ Actions ──────────────────────────────────────────────┐
│ F1:Help  F2:Workspaces  F5:Sync  P:Pull  p:Push       │

Frame 2 (3000ms):
┌─ Actions ──────────────────────────────────────────────┐
│ :Workspaces  F5:Sync  P:Pull  p:Push  s:FullSync      │  ← Scrolled

Frame 3 (6000ms):
┌─ Actions ──────────────────────────────────────────────┐
│  F5:Sync  P:Pull  p:Push  s:FullSync  b:Bulk  ?:Keys  │  ← Scrolled
```
**Marquee Animation:**
- Speed: 1 character every 150ms
- Pause: 3s at start, 2s at end
- Loop continuously
- Smooth scroll (not jumpy)
- FPS: 7-10 FPS

### Solution 2: Priority Filter (Simpler)
```
┌─ Actions ──────────────────────────────────────────────┐
│ F1:Help  F5:Sync  P:Pull  ESC:Cancel  F10:Quit  +more │  ← Shows most important
```
**Priority Order:**
1. Help (F1)
2. Context-specific primary action (F5, P, etc.)
3. Cancel (ESC) if job running
4. Quit (F10)
5. "+N more" indicator if overflow

### Solution 3: Multi-line (For wide terminals)
```
┌─ Actions ──────────────────────────────────────────────┐
│ F1:Help  F2:Workspaces  F5:Sync  P:Pull  p:Push       │
│ s:FullSync  b:Bulk  ?:Keys  /:Search  F10:Quit        │
└─────────────────────────────────────────────────────────┘
```
**Adaptive:** Use 2 lines if terminal width >120 cols

**Recommended:** Start with Solution 2 (priority filter), add marquee later if desired.

---

## Color Theme Specification

### Default Theme (Light)
```
Background:      Terminal default
Foreground:      Terminal default
Border (focused):    Blue (#0078D4)
Border (unfocused):  Gray (#808080)
Selection:           Cyan background (#00FFFF)
Success:             Green (#00FF00) ✓
Error:               Red (#FF0000) ✗
Warning:             Yellow (#FFFF00) ⚠
Progress Bar:        Blue (#0078D4)
Spinner:             Cyan (#00FFFF)
```

### Dark Theme (Hyperspace)
```
Background:      #1E1E1E
Foreground:      #D4D4D4
Border (focused):    Cyan (#00FFFF)
Border (unfocused):  Gray (#505050)
Selection:           Blue background (#0078D4)
Success:             Green (#00FF00) ✓
Error:               Red (#FF3030) ✗
Warning:             Yellow (#FFFF00) ⚠
Progress Bar:        Cyan (#00FFFF)
Spinner:             Cyan (#00FFFF)
Background Effect:   Stars (✦ ✧ ★) moving right
```

### Arctic Theme (Light + Snow)
```
Background:      #F0F8FF (Alice Blue)
Foreground:      #2F4F4F (Dark Slate Gray)
Border (focused):    Blue (#4682B4)
Border (unfocused):  Light Gray (#C0C0C0)
Selection:           Light Blue (#ADD8E6)
Success:             Dark Green (#006400) ✓
Error:               Dark Red (#8B0000) ✗
Warning:             Orange (#FFA500) ⚠
Progress Bar:        Steel Blue (#4682B4)
Spinner:             Blue (#4682B4)
Background Effect:   Snowflakes (❄ ❅) falling down
```

---

## Performance Budgets

**Target Metrics:**
- **Idle CPU:** ≤1% (no animations)
- **Active Animations:** ≤5% CPU (spinner + progress)
- **Full Effects:** ≤15% CPU (background + all animations)
- **Memory:** ≤50MB RSS
- **Frame Rate:** 30-60 FPS for animations
- **Input Latency:** <16ms (60 FPS feels instant)

**If Exceeding Budget:**
1. Reduce animation FPS
2. Disable background effects
3. Simplify shimmer
4. Profile and optimize

---

## Acceptance Criteria (Visual)

### Must Have (v3.1.1)
- [ ] Spinner animates smoothly (not static)
- [ ] Progress bar fills smoothly (not jumpy)
- [ ] Modals have drop shadows
- [ ] Focused borders distinct from unfocused
- [ ] Themed colors applied throughout
- [ ] UI responsive during all operations

### Should Have (v3.1.1)
- [ ] Progress bar shimmer animation
- [ ] Modal fade-in animation
- [ ] Keybindings don't overflow (priority filter)
- [ ] Status icons (✓ ✗ ⚠) visible

### Nice to Have (v3.2.0)
- [ ] Focus pulse animation
- [ ] Background effects (hyperspace, snow)
- [ ] Marquee scrolling for long action bars
- [ ] Smooth theme switching

---

## Implementation Notes

### Framework: tview vs. bubbletea

**Current (tview):**
- ✓ Mature, stable
- ✓ Widget-based (less boilerplate)
- ✗ Animations harder to implement
- ✗ Less modern

**Alternative (bubbletea):**
- ✓ Built for animations
- ✓ Modern, active development
- ✓ Elm architecture (clean)
- ✗ Requires complete rewrite
- ✗ More boilerplate

**Decision for Phase 6.5:**
1. **Try with tview first** - Fix animations in current framework
2. **If impossible** - Evaluate bubbletea migration timeline
3. **v3.1.1** - Ship with tview (if fixable)
4. **v3.2.0** - Consider bubbletea if needed

---

## Testing Checklist

**Visual Acceptance Testing:**
- [ ] Record video of TUI in action (10s)
- [ ] Verify spinner animates
- [ ] Verify progress bar smooth
- [ ] Verify modals have shadows
- [ ] Verify colors match theme
- [ ] Verify no flickering
- [ ] Measure FPS (should be 30-60)
- [ ] Measure CPU (should be ≤15%)

**User Validation:**
- [ ] Human tester confirms "smooth, modern, professional"
- [ ] Human tester confirms NOT "back in the 80s"
- [ ] Human tester approves visual quality

---

**This specification is authoritative. Any implementation that doesn't match this is incomplete.**

**Created:** 2025-10-20 (Phase 6.5)
**Owner:** Director + TUIUX + Builder
**Approver:** Human (UAT)
