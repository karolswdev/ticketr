# Ticketr Bubbletea TUI - Wireframe Specifications

**Vision**: Midnight Commander inspired, modern dual-panel layout with space-themed effects
**Framework**: Bubbletea + Lipgloss
**Design Philosophy**: High information density, clean aesthetics, accessible

---

## Table of Contents

1. [Main View (Dual Panel)](#1-main-view-dual-panel)
2. [Workspace Selector (Slide-out)](#2-workspace-selector-slide-out)
3. [Search Modal](#3-search-modal)
4. [Command Palette](#4-command-palette)
5. [Bulk Operations Modal](#5-bulk-operations-modal)
6. [Settings/Preferences](#6-settingspreferences)
7. [Help Screen](#7-help-screen)
8. [Color Palette](#color-palette)
9. [Responsive Behavior](#responsive-behavior)
10. [Animation & Effects](#animation--effects)
11. [Accessibility Notes](#accessibility-notes)

---

## 1. Main View (Dual Panel)

**Layout**: 50/50 split (adjustable) | Min terminal: 120x30 | Optimal: 160x40

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃ ╭─ TICKETR v3.1 ────────────────────────────────────────────────────────────────────────╮ [○] Idle  ◐ Workspace: my-project  ⚡ 42 tickets ┃
┃ │ 🚀 Welcome back! Last sync: 2 min ago | Press W for workspaces, / for search, ? for help │                                                    ┃
┃ ╰───────────────────────────────────────────────────────────────────────────────────────────╯                                                    ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ ╔═══════════════ 📋 Ticket Tree ════════════════╗ │ ╔═══════════════ 📄 Ticket Detail ══════════════╗                       ┃
┃ ║                                                ║ │ ║                                                ║                       ┃
┃ ║ 🔍 Filter: [____________________________]  ⌫  ║ │ ║ PROJ-1234: Implement authentication           ║                       ┃
┃ ║                                                ║ │ ║                                                ║                       ┃
┃ ║ ▼ 🏢 Epic: User Management        [3 tickets] ║ │ ║ Status: In Progress     Priority: High         ║                       ┃
┃ ║   [ ] PROJ-1234: Implement auth...    ● synced║ │ ║ Assignee: @john.doe                            ║                       ┃
┃ ║   [x] PROJ-1235: Add OAuth support    ● synced║ │ ║ Sprint: Sprint 42       Labels: backend, auth  ║                       ┃
┃ ║   [ ] PROJ-1236: Setup 2FA            ○ local ║ │ ║                                                ║                       ┃
┃ ║     ├─ TASK: Frontend integration     □ local ║ │ ║ ─────────────── Description ────────────────── ║                       ┃
┃ ║     └─ TASK: Backend API              □ local ║ │ ║                                                ║                       ┃
┃ ║                                                ║ │ ║ Implement a comprehensive authentication       ║                       ┃
┃ ║ ▼ 🐛 Bug Fixes                    [2 tickets] ║ │ ║ system with JWT tokens, refresh tokens, and    ║                       ┃
┃ ║   [ ] PROJ-1240: Fix login bug        ● synced║ │ ║ secure session management.                     ║                       ┃
┃ ║   [ ] PROJ-1241: Error handling       ○ local ║ │ ║                                                ║                       ┃
┃ ║                                                ║ │ ║ The system should support:                     ║                       ┃
┃ ║ ▼ ✨ Features                     [5 tickets] ║ │ ║ - Email/password login                         ║                       ┃
┃ ║   [x] PROJ-1250: Dark mode            ● synced║ │ ║ - OAuth providers (Google, GitHub)             ║                       ┃
┃ ║   [ ] PROJ-1251: Notifications        ○ local ║ │ ║ - Two-factor authentication                    ║                       ┃
┃ ║   [ ] PROJ-1252: Export reports       ○ local ║ │ ║ - Password reset flow                          ║                       ┃
┃ ║                                                ║ │ ║                                                ║                       ┃
┃ ║                                                ║ │ ║ ────────── Acceptance Criteria ─────────────── ║                       ┃
┃ ║                                                ║ │ ║                                                ║                       ┃
┃ ║                                                ║ │ ║ ✓ User can log in with email/password          ║                       ┃
┃ ║                                                ║ │ ║ ✓ JWT tokens expire after 1 hour               ║                       ┃
┃ ║                                                ║ │ ║ □ Refresh tokens work correctly                ║                       ┃
┃ ║ 3 selected                                     ║ │ ║ □ Failed login attempts are logged             ║                       ┃
┃ ║                                                ║ │ ║ □ Password reset emails are sent               ║                       ┃
┃ ║                                                ║ │ ║                                                ║                       ┃
┃ ╚════════════════════════════════════════════════╝ │ ║ ──────────── Custom Fields ──────────────────  ║                       ┃
┃                   ↑ FOCUSED (double border)        │ ║                                                ║                       ┃
┃                                                    │ ║ Story Points: 8                                ║                       ┃
┃                                                    │ ║ Original Estimate: 3 days                      ║                       ┃
┃                                                    │ ║ Time Spent: 1.5 days                           ║                       ┃
┃                                                    │ ║ Components: Auth Service, API Gateway          ║                       ┃
┃                                                    │ ║                                                ║                       ┃
┃                                                    │ ╚════════════════════════════════════════════════╝                       ┃
┃                                                    │                  (single border - not focused)                          ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ ╭──────────────────────────────────── 🎮 Actions ─────────────────────────────────────────────────────────────────────────────────────────────────╮ ┃
┃ │ [Enter] Open  [Space] Select  [W] Workspaces  [Tab] Next Panel  [j/k] Navigate  [h/l] Collapse/Expand  [b] Bulk  [/] Search  [:] Cmd  [?] Help │ ┃
┃ ╰──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯ ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Dimensions & Layout
- **Header Bar**: Fixed height (3 rows) - Status indicators and welcome message
- **Main Split**: Flexible (typically 50/50)
  - Left Panel: 40-60% width (min 50 cols)
  - Right Panel: 40-60% width (min 50 cols)
- **Action Bar**: Fixed height (3 rows) - Context-aware keybindings
- **Total Height**: Header(3) + Content(H-6) + Actions(3)

### Visual Elements
- **Borders**:
  - Focused: `╔═══╗` (double line, themed color)
  - Unfocused: `╭───╮` (single line, muted)
- **Icons**: Nerd font glyphs + emoji fallbacks
  - 📋 Tree, 📄 Detail, 🏢 Epic, 🐛 Bug, ✨ Feature, 🚀 Title
  - ● Synced (filled), ○ Local (empty), □ Task unchecked, ■ Task checked
- **Selection**: `[x]` checkbox with themed background highlight
- **Sync Status**:
  - `[○]` Idle (white)
  - `[◐]` Syncing (animated rotation: ◐◓◑◒)
  - `[●]` Success (green, 2s fade)
  - `[✗]` Error (red, persistent)

### Color Scheme (Default Theme)
```
Primary (focused):   #00FF00 (green)
Secondary:           #AAAAAA (gray)
Accent:              #00FFFF (cyan)
Success:             #00FF00 (green)
Warning:             #FFFF00 (yellow)
Error:               #FF0000 (red)
Background:          #1A1A1A (dark)
Text:                #FFFFFF (white)
Muted:               #666666 (dim gray)
```

### Focus Indicators
1. **Border Style Change**: Single → Double
2. **Color Pulse**: Subtle glow effect (0.5s cycle)
3. **Title Highlight**: Gradient on focused panel title
4. **Cursor**: Block cursor in active input fields

### State Transitions
```
Normal State:
  ├─ Press [W] → Workspace Selector (overlay)
  ├─ Press [/] → Search Modal (overlay)
  ├─ Press [:] → Command Palette (overlay)
  ├─ Press [b] → Bulk Operations (if selections exist)
  ├─ Press [?] → Help Screen (full overlay)
  └─ Press [Tab] → Cycle panel focus

Selection Mode (when tickets selected):
  ├─ Border color: Teal (#00FFFF)
  ├─ Title suffix: "(3 selected)"
  └─ Action bar highlights [b] Bulk operation
```

---

## 2. Workspace Selector (Slide-out)

**Layout**: Overlay from left, 35 cols wide | Animation: 150ms slide-in

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                                     │                                                                                                            ┃
┃ ╔══ 🌍 Workspaces ══════════════╗   │  [Main content dimmed/blurred]                                                                             ┃
┃ ║                                ║   │                                                                                                            ┃
┃ ║ 🔍 [_____________________]  ⌫  ║   │  ╔═══════════ Ticket Tree ═══════════╗    ╔═══════════ Detail ═══════════╗                                ┃
┃ ║                                ║   │  ║ (content dimmed 40% opacity)       ║    ║ (content dimmed 40% opacity) ║                                ┃
┃ ║ ────────── Recent ─────────── ║   │  ║                                    ║    ║                              ║                                ┃
┃ ║                                ║   │  ║  [ ] PROJ-123...                   ║    ║  PROJ-123: Auth...           ║                                ┃
┃ ║ ▶ my-project          ★ ●     ║   │  ║  [ ] PROJ-124...                   ║    ║                              ║                                ┃
┃ ║   jira.company.com             ║   │  ║                                    ║    ║  Description...              ║                                ┃
┃ ║   42 tickets | Synced 2m ago   ║   │  ║                                    ║    ║                              ║                                ┃
┃ ║                                ║   │  ╚════════════════════════════════════╝    ╚══════════════════════════════╝                                ┃
┃ ║   staging-env            ○     ║   │                                                                                                            ┃
┃ ║   jira-staging.company.com     ║   │                                                                                                            ┃
┃ ║   15 tickets | Never synced    ║   │                                                                                                            ┃
┃ ║                                ║   │                                                                                                            ┃
┃ ║ ──────── All Workspaces ────── ║   │                                                                                                            ┃
┃ ║                                ║   │                                                                                                            ┃
┃ ║   backend-team           ●     ║   │                                                                                                            ┃
┃ ║   jira.company.com             ║   │                                                                                                            ┃
┃ ║   128 tickets | Synced 10m ago ║   │                                                                                                            ┃
┃ ║                                ║   │                                                                                                            ┃
┃ ║   frontend-team          ●     ║   │                                                                                                            ┃
┃ ║   jira.company.com             ║   │                                                                                                            ┃
┃ ║   67 tickets | Synced 5m ago   ║   │                                                                                                            ┃
┃ ║                                ║   │                                                                                                            ┃
┃ ║   mobile-app             ○     ║   │                                                                                                            ┃
┃ ║   jira.company.com             ║   │                                                                                                            ┃
┃ ║   23 tickets | Synced 1h ago   ║   │                                                                                                            ┃
┃ ║                                ║   │                                                                                                            ┃
┃ ║ ───────────────────────────── ║   │                                                                                                            ┃
┃ ║                                ║   │                                                                                                            ┃
┃ ║ [n] New Workspace              ║   │                                                                                                            ┃
┃ ║                                ║   │                                                                                                            ┃
┃ ╚════════════════════════════════╝   │                                                                                                            ┃
┃    ↑ Shadow effect (▒▒)              │                                                                                                            ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ [Enter] Select  [j/k] Navigate  [/] Filter  [n] New  [Esc/W/F3] Close                                                                             ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Dimensions
- **Width**: 35 cols (fixed)
- **Height**: Full height minus header/footer
- **Position**: Left-aligned overlay
- **Animation**: Slide from left (-35 → 0) over 150ms with ease-out

### Visual Effects
- **Backdrop**: Main content dimmed to 40% opacity + Gaussian blur (if supported)
- **Shadow**: 2-col `▒▒` shadow on right edge (dim gray)
- **Highlight**: Current workspace has `▶` marker + bold text
- **Icons**:
  - `★` Default workspace
  - `●` Synced (green)
  - `○` Not synced (gray)
  - `◐` Syncing in progress (animated)

### Interaction States
1. **Closed**: Not visible
2. **Opening**: Slide-in animation (150ms)
3. **Open**: Fully visible, main content dimmed
4. **Filtering**: Search box active, filtered list updates
5. **Closing**: Slide-out animation (150ms)

### Accessibility
- **Screen Reader**: "Workspace selector overlay. 5 workspaces available. Current: my-project."
- **High Contrast**: Border increases to 2px width
- **Keyboard Nav**: Full keyboard control, no mouse required

---

## 3. Search Modal

**Layout**: Center overlay, 60% width x 40% height | Min: 80x20

```
                         Background dimmed 60% ↓
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                                                                                                                                                    ┃
┃         (Main content dimmed)                                                                                                                      ┃
┃                                                                                                                                                    ┃
┃              ╔══════════════════════════════ 🔍 Search Tickets ═══════════════════════════════╗                                                   ┃
┃              ║                                                                                 ║                                                   ┃
┃              ║  Query: [auth @john !high ~Sprint-42_____________________________]           ⌫  ║                                                   ┃
┃              ║                                                                                 ║                                                   ┃
┃              ║  ─────────────────────────── Filters ────────────────────────────────          ║                                                   ┃
┃              ║  @user   Assignee filter          !priority  Priority filter                   ║                                                   ┃
┃              ║  #ID     Ticket ID filter         ~sprint    Sprint filter                     ║                                                   ┃
┃              ║  /regex/ Regular expression       status:    Status exact match                ║                                                   ┃
┃              ║                                                                                 ║                                                   ┃
┃              ║  ─────────────────────── Results (8 matches) ─────────────────────             ║                                                   ┃
┃              ║                                                                                 ║                                                   ┃
┃              ║  ▶ [95%] PROJ-1234: Implement authentication                                   ║                                                   ┃
┃              ║    Matched: title, description, assignee                                       ║                                                   ┃
┃              ║    @john.doe | !high | ~Sprint-42                                              ║                                                   ┃
┃              ║                                                                                 ║                                                   ┃
┃              ║    [87%] PROJ-1236: Setup 2FA authentication                                   ║                                                   ┃
┃              ║    Matched: title, labels                                                      ║                                                   ┃
┃              ║    @john.doe | !medium | ~Sprint-42                                            ║                                                   ┃
┃              ║                                                                                 ║                                                   ┃
┃              ║    [75%] PROJ-1240: Fix login authentication bug                               ║                                                   ┃
┃              ║    Matched: title, description                                                 ║                                                   ┃
┃              ║    @sarah.smith | !high | ~Sprint-41                                           ║                                                   ┃
┃              ║                                                                                 ║                                                   ┃
┃              ║    [68%] PROJ-1250: Auth token expiration                                      ║                                                   ┃
┃              ║    Matched: title                                                              ║                                                   ┃
┃              ║    @john.doe | !low | ~Sprint-43                                               ║                                                   ┃
┃              ║                                                                                 ║                                                   ┃
┃              ║  ... 4 more results                                                            ║                                                   ┃
┃              ║                                                                                 ║                                                   ┃
┃              ╚═════════════════════════════════════════════════════════════════════════════════╝                                                   ┃
┃              ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒    ← Shadow                                       ┃
┃                                                                                                                                                    ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ [Enter] Open  [↑↓] Navigate  [Ctrl+F/B] Page  [Esc] Close                                                                                         ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Dimensions
- **Width**: 60% of terminal (min 80 cols)
- **Height**: 40% of terminal (min 20 rows)
- **Position**: Centered overlay

### Visual Effects
- **Backdrop**: 60% dim + blur
- **Shadow**: 2-col offset shadow (bottom-right)
- **Match Highlight**: Matched terms highlighted in yellow
- **Score Indicator**: `[95%]` relevance score badge
- **Animation**:
  - Fade-in (200ms) with slight scale-up (0.95 → 1.0)
  - Results appear with stagger effect (50ms delay each)

### Search Features
- **Real-time filtering**: Results update as you type
- **Fuzzy matching**: Typo-tolerant search
- **Filter syntax**:
  - `@john` - assignee filter (autocomplete)
  - `!high` - priority filter
  - `#PROJ-123` - ID filter (exact or prefix)
  - `~Sprint-42` - sprint filter
  - `/regex/` - regex pattern
  - `status:Done` - exact field match
- **Score Display**: Percentage match score with color gradient
  - 90-100%: Green
  - 70-89%: Yellow
  - <70%: White

### Interaction Flow
```
Idle → Press [/]
  ↓
Modal opens (fade-in + scale)
  ↓
User types query (real-time search)
  ↓
Results appear (stagger animation)
  ↓
User navigates (↑↓ or j/k)
  ↓
Press [Enter] → Open ticket detail, modal closes
Press [Esc] → Cancel, modal closes (fade-out)
```

---

## 4. Command Palette

**Layout**: Center overlay, 50% width x 60% height | Min: 70x25

```
                         Background dimmed 60% ↓
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                                                                                                                                                    ┃
┃                                                                                                                                                    ┃
┃                  ╔═════════════════════════ ⚡ Command Palette ════════════════════════╗                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ║  : [sync pull____________________________________]               ⌫  ║                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ║  ──────────────────── Commands (3 matches) ─────────────────────   ║                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ║  Sync & Operations                                                 ║                                                          ┃
┃                  ║  ▶ 🔄 sync:pull          Pull tickets from Jira          [Ctrl+P]  ║                                                          ┃
┃                  ║    🔄 sync:push          Push tickets to Jira            [Ctrl+U]  ║                                                          ┃
┃                  ║    🔄 sync:full          Full bidirectional sync         [Ctrl+S]  ║                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ║  ────────────────────── All Commands ───────────────────────────   ║                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ║  File Operations                                                   ║                                                          ┃
┃                  ║    📂 file:open          Open ticket file in editor       [Ctrl+O]  ║                                                          ┃
┃                  ║    💾 file:save          Save current changes             [Ctrl+S]  ║                                                          ┃
┃                  ║    📋 file:copy-path     Copy file path to clipboard                ║                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ║  Workspace                                                         ║                                                          ┃
┃                  ║    🌍 workspace:switch   Switch workspace                 [Ctrl+W]  ║                                                          ┃
┃                  ║    ➕ workspace:new      Create new workspace                       ║                                                          ┃
┃                  ║    ⚙️  workspace:settings Configure workspace                       ║                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ║  View                                                              ║                                                          ┃
┃                  ║    🎨 view:theme         Change color theme                         ║                                                          ┃
┃                  ║    📊 view:layout        Toggle layout (50/50, 60/40, etc)          ║                                                          ┃
┃                  ║    🌟 view:effects       Toggle visual effects                      ║                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ║  Bulk Operations                                                   ║                                                          ┃
┃                  ║    📦 bulk:update        Bulk update tickets                 [b]    ║                                                          ┃
┃                  ║    🔀 bulk:move          Bulk move tickets                          ║                                                          ┃
┃                  ║    🗑️  bulk:delete        Bulk delete tickets                        ║                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ║  Help & Info                                                       ║                                                          ┃
┃                  ║    ❓ help:keys          Show keyboard shortcuts           [?]      ║                                                          ┃
┃                  ║    📖 help:docs          Open documentation                         ║                                                          ┃
┃                  ║    ℹ️  help:about         About Ticketr                              ║                                                          ┃
┃                  ║                                                                     ║                                                          ┃
┃                  ╚═════════════════════════════════════════════════════════════════════╝                                                          ┃
┃                  ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒     ← Shadow                                                 ┃
┃                                                                                                                                                    ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ [Enter] Execute  [↑↓] Navigate  [Ctrl+K] Clear  [Esc] Close                                                                                       ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Dimensions
- **Width**: 50% of terminal (min 70 cols)
- **Height**: 60% of terminal (min 25 rows)
- **Position**: Centered overlay

### Visual Effects
- **Backdrop**: 60% dim + blur
- **Shadow**: 2-col offset shadow
- **Fuzzy Match Highlight**: Matched characters in command names highlighted
- **Category Headers**: Dimmed separators
- **Icons**: Emoji + nerd font glyphs for visual hierarchy
- **Shortcut Badges**: Right-aligned keyboard shortcut hints

### Command Categories
1. **Sync & Operations**: Pull, push, full sync
2. **File Operations**: Open, save, copy path
3. **Workspace**: Switch, new, settings
4. **View**: Theme, layout, effects
5. **Bulk Operations**: Update, move, delete
6. **Help & Info**: Keys, docs, about

### Interaction Flow
```
Idle → Press [:] or [Ctrl+P] or [F1]
  ↓
Palette opens (fade-in)
  ↓
User types command name (fuzzy search)
  ↓
Filtered list updates instantly
  ↓
Press [Enter] → Execute command, palette closes
Press [Esc] → Cancel, palette closes
```

### Fuzzy Search Algorithm
- Match characters in order (e.g., "sp" matches "sync:pull")
- Highlight matched characters in yellow
- Sort by relevance score (prefix match > mid-word match)
- Show keyboard shortcuts for quick access

---

## 5. Bulk Operations Modal

**Layout**: Center overlay, 70% width x 70% height | Multi-step wizard

### Step 1: Operation Selection

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                                                                                                                                                    ┃
┃                                                                                                                                                    ┃
┃            ╔════════════════════════════ 📦 Bulk Operations ═══════════════════════════╗                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ║  Selected: 3 tickets                                                      ║                                                         ┃
┃            ║  • PROJ-1235: Add OAuth support                                           ║                                                         ┃
┃            ║  • PROJ-1236: Setup 2FA                                                   ║                                                         ┃
┃            ║  • PROJ-1240: Fix login bug                                               ║                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ║  ──────────────────── Choose Operation ─────────────────────              ║                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ║  ▶ 🔄 Update Fields                                                       ║                                                         ┃
┃            ║    Change status, priority, assignee, or custom fields                    ║                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ║    🔀 Move Tickets                                                        ║                                                         ┃
┃            ║    Move selected tickets under a new parent                               ║                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ║    🗑️  Delete Tickets                                                      ║                                                         ┃
┃            ║    Permanently delete selected tickets (⚠️ irreversible)                  ║                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ║  ──────────────────────────────────────────────────────────               ║                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ║                         [Continue]  [Cancel]                               ║                                                         ┃
┃            ║                                                                            ║                                                         ┃
┃            ╚════════════════════════════════════════════════════════════════════════════╝                                                         ┃
┃            ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒                                                         ┃
┃                                                                                                                                                    ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ [Enter] Continue  [↑↓] Navigate  [Esc] Cancel                                                                                                     ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Step 2: Update Fields Form

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                                                                                                                                                    ┃
┃            ╔════════════════════════════ 🔄 Bulk Update Fields ═══════════════════════╗                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  Updating 3 tickets                                        [Step 2 of 3]  ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  ──────────────────── Field Values ─────────────────────                  ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  Status: [In Progress_____________________] ▼                             ║                                                          ┃
┃            ║          ⓘ Leave empty to keep existing values                            ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  Priority: [High__________________________] ▼                             ║                                                          ┃
┃            ║            Options: High, Medium, Low, Blocker                            ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  Assignee: [@john.doe____________________] 👤                             ║                                                          ┃
┃            ║            Search users: @john → john.doe, johnathan.smith                ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  Sprint: [Sprint-43_____________________] 🏃                               ║                                                          ┃
┃            ║          Active sprints: Sprint-42, Sprint-43, Sprint-44                  ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  Labels: [backend,auth,security_________] 🏷️                              ║                                                          ┃
┃            ║          Comma-separated. Press [Tab] for autocomplete.                   ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  ─────────────── Custom Fields (Optional) ───────────────                 ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  Story Points: [8_________________________]                               ║                                                          ┃
┃            ║  Components:   [Auth Service,API Gateway__]                               ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║  ──────────────────────────────────────────────────────                   ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ║                      [← Back]  [Preview]  [Apply]                          ║                                                          ┃
┃            ║                                                                            ║                                                          ┃
┃            ╚════════════════════════════════════════════════════════════════════════════╝                                                          ┃
┃            ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒                                                          ┃
┃                                                                                                                                                    ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ [Tab] Next Field  [Shift+Tab] Previous  [Enter] Apply  [Esc] Cancel                                                                               ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Step 3: Progress Tracking

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                                                                                                                                                    ┃
┃            ╔═══════════════════════ ⏳ Processing Bulk Operation ═══════════════════════╗                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ║  Progress: 2 of 3 tickets completed                       [Step 3 of 3]    ║                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ║  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░  67%        ║                                                        ┃
┃            ║                        ↑ Shimmer effect moves →                            ║                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ║  ──────────────────────── Status ─────────────────────────                 ║                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ║  ✅ PROJ-1235: Add OAuth support                                           ║                                                        ┃
┃            ║     Updated: Status → In Progress, Priority → High, Assignee → john.doe    ║                                                        ┃
┃            ║     Duration: 1.2s                                                         ║                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ║  ✅ PROJ-1236: Setup 2FA                                                   ║                                                        ┃
┃            ║     Updated: Status → In Progress, Priority → High, Assignee → john.doe    ║                                                        ┃
┃            ║     Duration: 0.9s                                                         ║                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ║  ⏳ PROJ-1240: Fix login bug                                               ║                                                        ┃
┃            ║     Processing... ⠋ (spinner animation)                                    ║                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ║  ──────────────────────────────────────────────────────                    ║                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ║  Estimated time remaining: ~1s                                             ║                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ║                              [Cancel]                                       ║                                                        ┃
┃            ║                                                                             ║                                                        ┃
┃            ╚═════════════════════════════════════════════════════════════════════════════╝                                                        ┃
┃            ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒                                                        ┃
┃                                                                                                                                                    ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ [Esc] Cancel Operation (will rollback)                                                                                                            ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Dimensions
- **Width**: 70% of terminal (min 80 cols)
- **Height**: 70% of terminal (min 30 rows)
- **Position**: Centered overlay

### Visual Effects
- **Progress Bar**: Animated shimmer effect (sweep left to right)
- **Spinner**: Rotating braille characters (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
- **Success/Error Icons**:
  - ✅ Success (green)
  - ❌ Error (red)
  - ⏳ Processing (animated)
- **Step Indicator**: "[Step 2 of 3]" in header

### Error Handling

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                                                                                                                                                    ┃
┃            ╔═══════════════════════ ⚠️  Operation Failed ══════════════════════╗                                                                 ┃
┃            ║                                                                    ║                                                                 ┃
┃            ║  Bulk update partially completed with errors                      ║                                                                 ┃
┃            ║                                                                    ║                                                                 ┃
┃            ║  ✅ Successful: 2 tickets                                         ║                                                                 ┃
┃            ║  ❌ Failed: 1 ticket                                              ║                                                                 ┃
┃            ║                                                                    ║                                                                 ┃
┃            ║  ──────────────────────────────────────────────────               ║                                                                 ┃
┃            ║                                                                    ║                                                                 ┃
┃            ║  ✅ PROJ-1235: Add OAuth support                                  ║                                                                 ┃
┃            ║  ✅ PROJ-1236: Setup 2FA                                          ║                                                                 ┃
┃            ║                                                                    ║                                                                 ┃
┃            ║  ❌ PROJ-1240: Fix login bug                                      ║                                                                 ┃
┃            ║     Error: Permission denied. You don't have edit access.         ║                                                                 ┃
┃            ║     Code: JIRA_403_FORBIDDEN                                      ║                                                                 ┃
┃            ║                                                                    ║                                                                 ┃
┃            ║  ──────────────────────────────────────────────────               ║                                                                 ┃
┃            ║                                                                    ║                                                                 ┃
┃            ║  Would you like to rollback successful changes?                   ║                                                                 ┃
┃            ║                                                                    ║                                                                 ┃
┃            ║                  [Rollback]  [Keep Changes]  [Close]               ║                                                                 ┃
┃            ║                                                                    ║                                                                 ┃
┃            ╚════════════════════════════════════════════════════════════════════╝                                                                 ┃
┃            ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒                                                                 ┃
┃                                                                                                                                                    ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

---

## 6. Settings/Preferences

**Layout**: Center overlay, 80% width x 80% height | Tabbed interface

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                                                                                                                                                    ┃
┃       ╔══════════════════════════════════════ ⚙️  Settings ═══════════════════════════════════════╗                                              ┃
┃       ║                                                                                             ║                                              ┃
┃       ║  ┌─────────┬─────────┬──────────┬────────────┬─────────┐                                  ║                                              ┃
┃       ║  │ General │ Theming │ Behavior │ Sync & Jira│ Advanced│                                  ║                                              ┃
┃       ║  ╞═════════╧═════════╧══════════╧════════════╧═════════╧═══════════════════════════════╡  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │  ────────────────────── General Settings ──────────────────────                     │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │  Default Workspace:  [my-project__________________] ▼                               │  ║                                              ┃
┃       ║  │  Editor:             [vim_____________________________]                             │  ║                                              ┃
┃       ║  │                      ⓘ External editor for ticket files                             │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │  Auto-sync on startup:      [✓] Enabled                                             │  ║                                              ┃
┃       ║  │  Sync interval (minutes):   [5____]  (0 = manual only)                              │  ║                                              ┃
┃       ║  │  Confirmation on delete:    [✓] Enabled                                             │  ║                                              ┃
┃       ║  │  Remember window size:      [✓] Enabled                                             │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │  ────────────────────── Layout Preferences ─────────────────────                    │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │  Panel split ratio:         ( ) 40/60  (●) 50/50  ( ) 60/40                         │  ║                                              ┃
┃       ║  │  Workspace panel:           ( ) Always visible  (●) Slide-out  ( ) Hidden           │  ║                                              ┃
┃       ║  │  Action bar position:       (●) Bottom  ( ) Top  ( ) Hidden                         │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │  ───────────────────── Accessibility ─────────────────────                          │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │  High contrast mode:        [ ] Enabled                                             │  ║                                              ┃
┃       ║  │  Screen reader support:     [✓] Enabled                                             │  ║                                              ┃
┃       ║  │  Reduce motion:             [ ] Enabled  (disables animations)                      │  ║                                              ┃
┃       ║  │  Font size multiplier:      [1.0__] (0.8 - 2.0)                                     │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  │                                                                                      │  ║                                              ┃
┃       ║  └──────────────────────────────────────────────────────────────────────────────────────┘  ║                                              ┃
┃       ║                                                                                             ║                                              ┃
┃       ║                           [Apply]  [Reset to Defaults]  [Cancel]                            ║                                              ┃
┃       ║                                                                                             ║                                              ┃
┃       ╚═════════════════════════════════════════════════════════════════════════════════════════════╝                                              ┃
┃       ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒                                              ┃
┃                                                                                                                                                    ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ [Tab] Next Tab  [Shift+Tab] Previous Tab  [Enter] Apply  [Esc] Cancel                                                                             ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Theming Tab

```
┃       ║  │  ────────────────────── Color Theme ──────────────────────                          │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Theme:  (●) Default (Green)  ( ) Dark (Blue)  ( ) Arctic (Cyan)  ( ) Custom        │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  ┌─────────── Theme Preview ────────────┐                                           │  ║
┃       ║  │  │ ╔═══ Focused Panel ════╗             │                                           │  ║
┃       ║  │  │ ║ Sample content here   ║             │                                           │  ║
┃       ║  │  │ ╚═══════════════════════╝             │                                           │  ║
┃       ║  │  │ ╭─── Unfocused Panel ───╮             │                                           │  ║
┃       ║  │  │ │ Sample content here    │             │                                           │  ║
┃       ║  │  │ ╰────────────────────────╯             │                                           │  ║
┃       ║  │  │ ✅ Success  ⚠️  Warning  ❌ Error      │                                           │  ║
┃       ║  │  └────────────────────────────────────────┘                                           │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  ─────────────────── Visual Effects ────────────────────                            │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Animations:            [✓] Enabled  (fade, slide, pulse)                           │  ║
┃       ║  │  Drop shadows:          [✓] Enabled  (modals and overlays)                          │  ║
┃       ║  │  Focus pulse:           [ ] Enabled  (subtle border glow)                           │  ║
┃       ║  │  Gradient titles:       [ ] Enabled  (focused panel titles)                         │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  ────────────────── 🚀 Space Effects ──────────────────                             │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Ambient mode:          (●) Off  ( ) Hyperspace  ( ) Snow  ( ) Matrix               │  ║
┃       ║  │                         ⓘ Easter egg: animated background particles                 │  ║
┃       ║  │  Particle density:      [▓▓▓░░░░░░░]  2%  (higher = more particles)                 │  ║
┃       ║  │  Animation speed:       [▓▓▓▓▓░░░░░]  100ms per frame                               │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Easter egg trigger:    [ Konami code: ↑↑↓↓←→←→ba ]                                 │  ║
```

### Behavior Tab

```
┃       ║  │  ────────────────── Keyboard Shortcuts ──────────────────────                       │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Preset:  (●) Vim-style  ( ) Emacs  ( ) Custom                                      │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Navigation:                                                                         │  ║
┃       ║  │    Up:           [k____]    Down:         [j____]                                   │  ║
┃       ║  │    Left:         [h____]    Right:        [l____]                                   │  ║
┃       ║  │    Page Up:      [Ctrl+B]   Page Down:    [Ctrl+F]                                  │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Actions:                                                                            │  ║
┃       ║  │    Search:       [/____]    Command:      [:____]                                   │  ║
┃       ║  │    Workspace:    [W____]    Help:         [?____]                                   │  ║
┃       ║  │    Bulk ops:     [b____]    Quit:         [q____]                                   │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  ─────────────────── Tree Behavior ─────────────────────                            │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Auto-expand on load:       [ ] Enabled                                             │  ║
┃       ║  │  Remember expansion state:  [✓] Enabled                                             │  ║
┃       ║  │  Group by:                  (●) Type  ( ) Status  ( ) Assignee  ( ) None            │  ║
┃       ║  │  Sort by:                   ( ) ID  (●) Priority  ( ) Updated  ( ) Created          │  ║
┃       ║  │  Show task count:           [✓] Enabled  (e.g., "[3 tickets]")                      │  ║
```

### Sync & Jira Tab

```
┃       ║  │  ─────────────────── Sync Settings ───────────────────────                          │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Auto-sync:                 [✓] Enabled                                             │  ║
┃       ║  │  Sync interval (minutes):   [5____]  (0 = manual only)                              │  ║
┃       ║  │  Sync on startup:           [✓] Enabled                                             │  ║
┃       ║  │  Background sync:           [✓] Enabled  (non-blocking)                             │  ║
┃       ║  │  Retry failed syncs:        [✓] Enabled  (max 3 retries)                            │  ║
┃       ║  │  Notify on sync errors:     [✓] Enabled                                             │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  ──────────────────── Jira Integration ─────────────────────                        │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  API request timeout:       [30___] seconds                                         │  ║
┃       ║  │  Batch size (bulk ops):     [10___] tickets per batch                               │  ║
┃       ║  │  Rate limiting:             [✓] Enabled  (respect Jira limits)                      │  ║
┃       ║  │  Cache duration:            [60___] minutes                                         │  ║
┃       ║  │                                                                                      │  ║
┃       ║  │  Custom field mapping:      [Configure Mappings...]                                 │  ║
┃       ║  │                             ⓘ Map Jira custom fields to ticket fields               │  ║
```

---

## 7. Help Screen

**Layout**: Full screen overlay | Scrollable content

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃ ╔═══════════════════════════════════════════════ 🚀 Ticketr Help ════════════════════════════════════════════════╗                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ───────────────────────────────────── Quick Start ─────────────────────────────────────                       ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Welcome to Ticketr - Jira-Markdown Synchronization Tool                                                       ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Getting Started:                                                                                               ║                               ┃
┃ ║  1. Press [W] to open workspace selector                                                                       ║                               ┃
┃ ║  2. Select or create a workspace                                                                               ║                               ┃
┃ ║  3. Press [P] to pull tickets from Jira                                                                        ║                               ┃
┃ ║  4. Navigate with [j/k] or arrow keys                                                                          ║                               ┃
┃ ║  5. Press [Enter] to view ticket details                                                                       ║                               ┃
┃ ║  6. Press [?] anytime to return to this help                                                                   ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ──────────────────────────── Global Navigation ────────────────────────────                                   ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Tab / Shift+Tab    Cycle between panels (tree ↔ detail)                                                       ║                               ┃
┃ ║  W / F3             Toggle workspace selector                                                                  ║                               ┃
┃ ║  /                  Open fuzzy search                                                                          ║                               ┃
┃ ║  : or Ctrl+P        Open command palette                                                                       ║                               ┃
┃ ║  ?                  Toggle this help screen                                                                    ║                               ┃
┃ ║  q / Ctrl+C         Quit application                                                                           ║                               ┃
┃ ║  Esc                Go back / Close modal                                                                      ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ─────────────────────────── Ticket Tree Panel ─────────────────────────                                       ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  j / k              Navigate down / up                                                                         ║                               ┃
┃ ║  h / l              Collapse / expand node                                                                     ║                               ┃
┃ ║  g / G              Go to first / last ticket                                                                  ║                               ┃
┃ ║  Enter              Open ticket detail                                                                         ║                               ┃
┃ ║  Space              Toggle ticket selection (for bulk ops)                                                     ║                               ┃
┃ ║  a / A              Select all / deselect all                                                                  ║                               ┃
┃ ║  b                  Open bulk operations (requires selection)                                                  ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ──────────────────────── Ticket Detail Panel ──────────────────────────                                       ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  j / k              Scroll down / up (one line)                                                                ║                               ┃
┃ ║  Ctrl+F / Ctrl+B    Page down / up                                                                             ║                               ┃
┃ ║  Ctrl+D / Ctrl+U    Half-page down / up                                                                        ║                               ┃
┃ ║  g / G              Go to top / bottom                                                                         ║                               ┃
┃ ║  e                  Edit mode (if supported)                                                                   ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ───────────────────────── Sync Operations ─────────────────────────                                           ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  p                  Push local changes to Jira                                                                 ║                               ┃
┃ ║  P                  Pull tickets from Jira                                                                     ║                               ┃
┃ ║  s                  Full sync (pull + push)                                                                    ║                               ┃
┃ ║  r                  Refresh current workspace                                                                  ║                               ┃
┃ ║  Esc (during sync)  Cancel active sync operation                                                               ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ─────────────────────────── Search Modal ──────────────────────────                                           ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  /                  Open search modal                                                                          ║                               ┃
┃ ║  @user              Filter by assignee (e.g., @john)                                                           ║                               ┃
┃ ║  #ID                Filter by ticket ID (e.g., #PROJ-123)                                                      ║                               ┃
┃ ║  !priority          Filter by priority (e.g., !high)                                                           ║                               ┃
┃ ║  ~sprint            Filter by sprint (e.g., ~Sprint-42)                                                        ║                               ┃
┃ ║  /regex/            Regex pattern search                                                                       ║                               ┃
┃ ║  status:value       Exact field match (e.g., status:Done)                                                      ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Multiple filters can be combined:                                                                             ║                               ┃
┃ ║  Example: "auth @john !high ~Sprint-42" finds auth-related tickets assigned to John                            ║                               ┃
┃ ║           with high priority in Sprint 42                                                                      ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ────────────────────────── Bulk Operations ────────────────────────                                           ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  1. Select tickets with [Space]                                                                                ║                               ┃
┃ ║  2. Press [b] to open bulk operations menu                                                                     ║                               ┃
┃ ║  3. Choose operation:                                                                                          ║                               ┃
┃ ║     • Update Fields - Change status, priority, assignee, etc.                                                  ║                               ┃
┃ ║     • Move Tickets - Move under new parent                                                                     ║                               ┃
┃ ║     • Delete Tickets - Permanent deletion (with confirmation)                                                  ║                               ┃
┃ ║  4. Fill in form and apply                                                                                     ║                               ┃
┃ ║  5. Monitor progress (real-time updates)                                                                       ║                               ┃
┃ ║  6. Review results and optionally rollback on errors                                                           ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ───────────────────────── Visual Indicators ────────────────────────                                          ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Borders:                                                                                                       ║                               ┃
┃ ║    ╔═══╗  Focused panel (double line, themed color)                                                            ║                               ┃
┃ ║    ╭───╮  Unfocused panel (single line, muted)                                                                 ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Sync Status:                                                                                                   ║                               ┃
┃ ║    ● Synced with Jira (green)                                                                                  ║                               ┃
┃ ║    ○ Local only (white/gray)                                                                                   ║                               ┃
┃ ║    ◐ Syncing in progress (animated rotation)                                                                   ║                               ┃
┃ ║    ✗ Sync error (red)                                                                                          ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Selection:                                                                                                     ║                               ┃
┃ ║    [ ] Unselected ticket                                                                                       ║                               ┃
┃ ║    [x] Selected ticket                                                                                         ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ──────────────────────────── Color Themes ──────────────────────────                                          ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Default (Green)  - Classic green/white theme                                                                  ║                               ┃
┃ ║  Dark (Blue)      - Dark theme with blue accents + hyperspace effect                                           ║                               ┃
┃ ║  Arctic (Cyan)    - Light cyan theme with snow effect                                                          ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Change theme: Open settings (: → settings) or set TICKETR_THEME env var                                       ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ─────────────────────────── Easter Eggs 🥚 ──────────────────────────                                         ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  • Konami Code (↑↑↓↓←→←→ba) - Activate hyperspace mode!                                                        ║                               ┃
┃ ║  • Type "retro" in command palette - Enable 80s retro theme                                                    ║                               ┃
┃ ║  • Type "matrix" - Enter the Matrix (green rain effect)                                                        ║                               ┃
┃ ║  • Press Ctrl+Alt+S three times fast - Space invaders mini-game (WIP)                                          ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  ──────────────────────────── About Ticketr ─────────────────────────                                          ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ║  Version: 3.1.1                                                                                                 ║                               ┃
┃ ║  Architecture: Hexagonal (Ports & Adapters)                                                                    ║                               ┃
┃ ║  Framework: Bubbletea + Lipgloss                                                                               ║                               ┃
┃ ║  License: MIT                                                                                                   ║                               ┃
┃ ║  Repository: github.com/karolswdev/ticktr                                                                      ║                               ┃
┃ ║                                                                                                                 ║                               ┃
┃ ╚═════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝                               ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ [Esc] or [?] Close  [j/k] Scroll  [Ctrl+F/B] Page  [g/G] Top/Bottom                                                                               ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Dimensions
- **Layout**: Full screen overlay
- **Scrollable**: Yes (content extends beyond viewport)
- **Sections**: Collapsible with headers

---

## Color Palette

### Default Theme (Green/Classic)
```
┌─────────────────────────────────────────────┐
│ Primary (Focused):      #00FF00  ████████  │
│ Secondary (Unfocused):  #AAAAAA  ████████  │
│ Accent:                 #00FFFF  ████████  │
│ Success:                #00FF00  ████████  │
│ Warning:                #FFFF00  ████████  │
│ Error:                  #FF0000  ████████  │
│ Background:             #1A1A1A  ████████  │
│ Text Primary:           #FFFFFF  ████████  │
│ Text Muted:             #666666  ████████  │
│ Shadow:                 #000000  ████████  │
│                            (40% opacity)    │
└─────────────────────────────────────────────┘
```

### Dark Theme (Blue/Space)
```
┌─────────────────────────────────────────────┐
│ Primary:                #0080FF  ████████  │
│ Secondary:              #555555  ████████  │
│ Accent:                 #00D0FF  ████████  │
│ Success:                #00FF80  ████████  │
│ Warning:                #FFD700  ████████  │
│ Error:                  #FF4040  ████████  │
│ Background:             #0A0A0A  ████████  │
│ Text:                   #E0E0E0  ████████  │
│ Ambient (Stars):        #FFFFFF  ████████  │
│                            (20% opacity)    │
└─────────────────────────────────────────────┘
```

### Arctic Theme (Cyan/Snow)
```
┌─────────────────────────────────────────────┐
│ Primary:                #00CED1  ████████  │
│ Secondary:              #87CEEB  ████████  │
│ Accent:                 #40E0D0  ████████  │
│ Success:                #00FF00  ████████  │
│ Warning:                #FFA500  ████████  │
│ Error:                  #DC143C  ████████  │
│ Background:             #F0F8FF  ████████  │
│ Text:                   #2F4F4F  ████████  │
│ Snow Particles:         #FFFFFF  ████████  │
│                            (60% opacity)    │
└─────────────────────────────────────────────┘
```

### Color-Blind Friendly Modes
- **Deuteranopia**: Replace green with blue/yellow
- **Protanopia**: Replace red with blue/purple
- **Tritanopia**: Replace blue with red/pink

---

## Responsive Behavior

### Minimum Terminal Sizes
```
Absolute Minimum:  80 cols × 24 rows
Recommended:      120 cols × 30 rows
Optimal:          160 cols × 40 rows
Ultra-wide:       200 cols × 50 rows
```

### Layout Adaptation

#### 80×24 (Minimum)
- Single panel mode (toggle between tree/detail with Tab)
- Workspace panel: Full overlay (not slide-out)
- Action bar: Compact, scrolling text
- Help: Reduced content, essential only

#### 120×30 (Recommended)
- Dual panel: 50/50 split
- Workspace panel: Slide-out (35 cols)
- Action bar: Full keybindings visible
- Modals: 60% width

#### 160×40 (Optimal)
- Dual panel: Adjustable split (40/60, 50/50, 60/40)
- All features fully visible
- Modals: 70% width with comfortable padding
- Help: Multi-column layout

#### 200×50+ (Ultra-wide)
- Tri-panel mode option (workspace | tree | detail)
- Side-by-side modals
- Enhanced information density
- Picture-in-picture preview mode

### Breakpoint Behavior
```
Width < 80:   Error message "Terminal too narrow"
80-119:       Single panel, compact mode
120-159:      Dual panel, standard mode
160+:         Dual panel, enhanced mode

Height < 24:  Error message "Terminal too short"
24-29:        Compact vertical spacing
30-39:        Standard vertical spacing
40+:          Generous vertical spacing
```

---

## Animation & Effects

### Entry Animations
```
Modal Open:
  Duration: 200ms
  Effect: Fade-in (0% → 100% opacity) + Scale (0.95 → 1.0)
  Easing: ease-out

Slide-out Panel:
  Duration: 150ms
  Effect: Slide from left (-35 cols → 0)
  Easing: ease-in-out

Search Results:
  Duration: 50ms per item
  Effect: Stagger appearance with fade-in
  Max stagger delay: 500ms total
```

### Exit Animations
```
Modal Close:
  Duration: 150ms
  Effect: Fade-out (100% → 0% opacity) + Scale (1.0 → 0.95)
  Easing: ease-in

Slide-out Close:
  Duration: 150ms
  Effect: Slide to left (0 → -35 cols)
  Easing: ease-in-out
```

### Micro-Interactions
```
Focus Change:
  Duration: 100ms
  Effect: Border style swap + color transition
  Border: Single → Double (focused)
  Color: Secondary → Primary

Button Hover/Focus:
  Duration: 100ms
  Effect: Background color change + slight scale
  Scale: 1.0 → 1.05

Selection Toggle:
  Duration: 50ms
  Effect: Checkbox fill animation
  [ ] → [x] with flash effect
```

### Progress Indicators
```
Progress Bar Shimmer:
  Duration: 2s loop
  Effect: Light band sweeps left to right across filled portion
  Width: 3 characters
  Direction: Bounce (reverses at edges)

Spinner (Braille):
  Frames: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
  Duration: 80ms per frame
  Loop: Infinite while processing

Sync Status Pulse:
  Duration: 1s loop
  Effect: Opacity pulse (60% → 100% → 60%)
  Only when syncing
```

### Space Effects (Easter Eggs)

#### Hyperspace Mode
```
Activation: Konami code (↑↑↓↓←→←→ba) or Dark theme + ambient enabled
Effect: Animated stars moving from center to edges
Density: 2% of screen (configurable)
Speed: 100ms per frame
Appearance: White dots (●) with motion blur trails
```

#### Snow Mode
```
Activation: Arctic theme + ambient enabled
Effect: Falling snowflakes
Density: 1.5% of screen
Speed: 120ms per frame
Characters: ❄ ❅ ❆ ∗ ·
Physics: Gentle side-to-side drift
```

#### Matrix Mode
```
Activation: Type "matrix" in command palette
Effect: Falling green characters (0-9, A-Z)
Density: 5% of screen columns
Speed: 50ms per frame
Appearance: Fades as it falls (100% → 20% opacity)
```

---

## Accessibility Notes

### Screen Reader Support
```
ARIA Labels (conceptual for TUI):
- "Main ticket tree. 42 tickets. Press j/k to navigate."
- "Ticket detail: PROJ-1234. Status In Progress. Priority High."
- "Workspace selector overlay. 5 workspaces. Current: my-project."
- "Bulk operation in progress. 2 of 3 tickets completed. 67%."

Focus announcements:
- Panel focus: "Now in ticket tree panel"
- Modal open: "Search modal opened. Type to filter."
- Operation complete: "Bulk update completed. 3 tickets updated successfully."
```

### Keyboard-Only Navigation
- **All features** accessible without mouse
- **Logical tab order** through panels and modals
- **Escape** always goes back or closes
- **Enter** confirms actions
- **No dead ends** - always a way to navigate out

### High Contrast Mode
```
Enabled when: TICKETR_HIGH_CONTRAST=1 or via settings

Changes:
- Border width: 1px → 2px
- Color contrast: WCAG AAA compliant (7:1 minimum)
- Remove subtle gradients
- Increase icon size by 20%
- Bolder fonts for text
- No transparency/blur effects
```

### Reduce Motion Mode
```
Enabled when: TICKETR_REDUCE_MOTION=1 or via settings

Changes:
- Disable all fade/slide animations
- Instant modals (no transition)
- No pulse/shimmer effects
- Static progress bars (no animation)
- Snap focus changes (no smooth transition)
- Disable ambient effects completely
```

### Color-Blind Modes

#### Deuteranopia (Red-Green)
```
Changes:
- Success: Green → Blue
- Error: Red → Orange
- Warning: Yellow → Purple
- Primary: Green → Blue
```

#### Protanopia (Red-Blind)
```
Changes:
- Error: Red → Blue
- Warning: Yellow → Cyan
- Success: Green → Purple
```

#### Tritanopia (Blue-Yellow)
```
Changes:
- Primary: Blue → Red
- Info: Cyan → Magenta
- Warning: Yellow → Red
```

### Font Size Scaling
```
Multiplier: 0.8 - 2.0 (default 1.0)

Affects:
- All text rendering
- Icon sizes (proportional)
- Padding/margins (proportional)
- Modal sizes (proportional)

Terminal must support: Unicode box drawing at various scales
```

---

## Implementation Notes

### Bubbletea Model Structure
```go
type MainModel struct {
    // Layout
    width, height  int
    panels         []Panel  // [tree, detail, workspace]
    focusedPanel   int
    splitRatio     float64  // 0.5 = 50/50

    // State
    workspace      *Workspace
    tickets        []Ticket
    selectedTicket *Ticket
    selectedIDs    []string  // For bulk ops

    // Modals
    activeModal    Modal  // search, command, bulk, settings, help
    modalStack     []Modal  // For nested modals

    // Effects
    theme          Theme
    animations     []Animation
    ambientEffect  AmbientEffect

    // Sync
    syncStatus     SyncStatus  // idle, syncing, success, error
    syncProgress   float64     // 0.0 - 1.0
}
```

### Lipgloss Style Templates
```go
var (
    FocusedBorderStyle = lipgloss.NewStyle().
        Border(lipgloss.DoubleBorder()).
        BorderForeground(theme.Primary).
        Padding(1)

    UnfocusedBorderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(theme.Secondary).
        Padding(1)

    ModalStyle = lipgloss.NewStyle().
        Border(lipgloss.DoubleBorder()).
        BorderForeground(theme.Primary).
        Padding(1, 2).
        Width(80).
        Align(lipgloss.Center)

    ShadowStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("0")).
        Background(lipgloss.Color("0")).
        Opacity(0.4)
)
```

### Key Binding Map
```go
var DefaultKeyMap = KeyMap{
    // Global
    Quit:          key.NewBinding(key.WithKeys("q", "ctrl+c")),
    Help:          key.NewBinding(key.WithKeys("?")),
    TabNext:       key.NewBinding(key.WithKeys("tab")),
    TabPrev:       key.NewBinding(key.WithKeys("shift+tab")),

    // Navigation
    Up:            key.NewBinding(key.WithKeys("k", "up")),
    Down:          key.NewBinding(key.WithKeys("j", "down")),
    Left:          key.NewBinding(key.WithKeys("h", "left")),
    Right:         key.NewBinding(key.WithKeys("l", "right")),
    PageDown:      key.NewBinding(key.WithKeys("ctrl+f")),
    PageUp:        key.NewBinding(key.WithKeys("ctrl+b")),

    // Actions
    Select:        key.NewBinding(key.WithKeys("enter")),
    Back:          key.NewBinding(key.WithKeys("esc")),
    ToggleSelect:  key.NewBinding(key.WithKeys(" ")),

    // Features
    Search:        key.NewBinding(key.WithKeys("/")),
    Command:       key.NewBinding(key.WithKeys(":", "ctrl+p", "f1")),
    Workspace:     key.NewBinding(key.WithKeys("W", "f3")),
    BulkOps:       key.NewBinding(key.WithKeys("b")),

    // Sync
    Pull:          key.NewBinding(key.WithKeys("P")),
    Push:          key.NewBinding(key.WithKeys("p")),
    Sync:          key.NewBinding(key.WithKeys("s")),
    Refresh:       key.NewBinding(key.WithKeys("r")),
}
```

---

## Performance Considerations

### Virtual Scrolling
- Only render visible rows + buffer (10 rows above/below)
- Lazy load ticket details on demand
- Cache rendered strings for static content

### Debouncing
- Search input: 150ms debounce
- Window resize: 100ms debounce
- Filter changes: 200ms debounce

### Memory Management
- Limit ticket cache: 1000 tickets in memory
- Paginate large result sets
- Unload hidden modal content
- Dispose animation frames when not visible

### Rendering Optimization
- Diff-based rendering (only update changed regions)
- Batch DOM updates in single frame
- Skip animations when terminal too small
- Reduce particle count on slower terminals

---

## Testing Scenarios

### Manual Test Cases
1. **Responsive Resize**: Drag terminal from 80×24 → 200×50, verify layout adapts
2. **Keyboard Navigation**: Navigate entire app without mouse
3. **Modal Stacking**: Open search → help → settings, close in order
4. **Bulk Ops Error**: Trigger partial failure, verify rollback prompt
5. **Sync Cancel**: Start sync, press Esc mid-operation
6. **Theme Switch**: Change themes, verify all colors update
7. **Easter Egg**: Konami code activation
8. **High Contrast**: Enable accessibility mode, verify contrast
9. **Long Content**: 1000+ tickets, verify smooth scrolling
10. **Network Timeout**: Simulate slow Jira API, verify timeout handling

---

## Future Enhancements

### v3.2.0 Features
- [ ] Tri-panel mode for ultra-wide terminals
- [ ] Ticket timeline view (activity feed)
- [ ] Inline ticket editing (no modal)
- [ ] Drag-and-drop ticket reordering
- [ ] Custom dashboard widgets
- [ ] Graph view (ticket relationships)
- [ ] AI-powered search suggestions
- [ ] Voice command support (experimental)

### Easter Egg Ideas
- Space Invaders mini-game during sync wait
- Ticket tetris (organize tickets by dragging)
- Conway's Game of Life in background
- ASCII art gallery (procedural generation)
- Music visualizer sync with keypress sounds

---

**Document Version**: 1.0
**Last Updated**: 2025-10-22
**Author**: Claude (Anthropic)
**Status**: Ready for Implementation
