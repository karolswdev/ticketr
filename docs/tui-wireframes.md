# Ticketr TUI Wireframes & Design System

**Version:** 1.0
**Date:** January 2025
**Status:** Design Specification

---

## Design Philosophy

**Inspiration Sources:**
- **k9s** - Kubernetes management (navigation, resource views)
- **lazygit** - Git operations (panel layout, diff views)
- **htop** - System monitoring (real-time updates, color coding)
- **vim** - Modal editing (keybindings, command mode)
- **ranger** - File management (miller columns, preview pane)

**Core Principles:**
1. **Information Density** - Maximum data, minimum chrome
2. **Keyboard-First** - Every action accessible without mouse
3. **Progressive Disclosure** - Details on demand
4. **Visual Hierarchy** - Color and spacing convey importance
5. **Contextual Actions** - Show only relevant commands

---

## 1. Main Dashboard (Project View)

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Ticketr v3.0.0  Project: backend [BACK]  ✓ Connected  ↻ 2s ago  karol@jira  ║
╠════════════════╤═════════════════════════════════════════════════════════════╣
║                │ EPIC → STORY → TASK                            [23 items]   ║
║ PROJECTS    4  │ ┌─────────────────────────────────────────────────────────┐ ║
║ ═══════════    │ │ ▼ 🎯 BACK-100 Authentication System           ●●●●○ 80% │ ║
║ ● backend      │ │   ▼ 📘 BACK-101 User Registration            ●●●●● Done │ ║
║ ○ frontend     │ │     ✓ BACK-201 Create API endpoint          ✅ Done    │ ║
║ ○ mobile       │ │     ✓ BACK-202 Design registration form     ✅ Done    │ ║
║ ○ personal     │ │     ✓ BACK-203 Add validation logic         ✅ Done    │ ║
║                │ │   ▶ 📘 BACK-102 Login Flow                  ●●●○○ 60%  │ ║
║ VIEWS       8  │ │   ▼ 📘 BACK-103 Password Reset              ●●○○○ 40%  │ ║
║ ═══════════    │ │     ⚡ BACK-206 Email template              🔄 In Prog │ ║
║ □ My Tasks     │ │     ○ BACK-207 Reset endpoint               ⏸  To Do   │ ║
║ □ Sprint 23    │ │     ○ BACK-208 Expiry logic                 ⏸  To Do   │ ║
║ □ Blocked      │ │ ▶ 🎯 BACK-110 Payment Integration            ○○○○○ 0%   │ ║
║ □ In Review    │ │ ▼ 🎯 BACK-120 Performance Optimization       ●●○○○ 40%  │ ║
║ □ Recent       │ │   ▶ 📘 BACK-121 Database indexing           ●○○○○ 20%  │ ║
║                │ │   ▼ 📘 BACK-122 Caching layer               ●●●○○ 60%  │ ║
║ FILTERS     +  │ │     ✓ BACK-221 Redis setup                  ✅ Done    │ ║
║ ═══════════    │ │     ⚡ BACK-222 Cache invalidation          🔄 In Prog │ ║
║ Type: All  ▼   │ └─────────────────────────────────────────────────────────┘ ║
║ Status: Open   │ ┌─────────────────────────────────────────────────────────┐ ║
║ Sprint: 23     │ │ 📊 Statistics         │ 🔄 Recent Activity              │ ║
║                │ │ Total: 23 tickets     │ 2m ago: BACK-222 updated       │ ║
║ ACTIONS        │ │ Done: 8 (35%)         │ 15m ago: BACK-206 started      │ ║
║ ═══════════    │ │ In Progress: 5        │ 1h ago: BACK-201-203 completed │ ║
║ [c] Create     │ │ Blocked: 0            │ 3h ago: Sprint 23 started      │ ║
║ [e] Edit       │ └─────────────────────────────────────────────────────────┘ ║
║ [p] Push       │                                                             ║
║ [P] Pull       │ Legend: 🎯Epic 📘Story ⚡Task ✅Done 🔄Progress ⏸Todo 🔴Block║
╠════════════════╧═════════════════════════════════════════════════════════════╣
║ [/]Search [:]Cmd [:w]Push [:q]Quit [Tab]Panel [?]Help [1-9]Projects │ 1:32PM ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

**Key Design Elements:**
- **Tree hierarchy** with collapsible nodes (▼/▶)
- **Progress bars** for epics/stories (●●●○○)
- **Status indicators** with emoji for quick scanning
- **Color coding**: Green=Done, Yellow=Progress, Gray=Todo, Red=Blocked
- **Statistics panel** for project health
- **Activity feed** for team awareness

---

## 2. Ticket Detail View (Split Screen)

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Ticketr  BACK-102: Login Flow  [Story]  Modified*  ↵ Save  ESC Cancel       ║
╠═══════════════════════════════╤══════════════════════════════════════════════╣
║ HIERARCHY                     │ DETAILS                                      ║
║ ─────────────────────         │ ──────────────────────────────────────────── ║
║ 🎯 BACK-100 Auth System       │ Title: Login Flow Implementation             ║
║ └─ 📘 BACK-102 Login Flow ←   │                                              ║
║    ├─ ✓ BACK-204 JWT tokens   │ Description:                                 ║
║    ├─ ⚡ BACK-205 Session mgmt│ ┌──────────────────────────────────────────┐ ║
║    └─ ○ BACK-206 Remember me  │ │As a user, I want to securely log into   │ ║
║                               │ │the application using my credentials so   │ ║
║ METADATA                      │ │that I can access protected resources.    │ ║
║ ─────────────────────         │ │                                           │ ║
║ Type:        Story            │ │Implementation should include:            │ ║
║ Status:      In Progress      │ │- JWT token generation                    │ ║
║ Priority:    High             │ │- Secure session management               │ ║
║ Sprint:      Sprint 23        │ │- Remember me functionality               │ ║
║ Points:      8                │ │- Rate limiting on login attempts         │ ║
║ Assignee:    @john.doe        │ └──────────────────────────────────────────┘ ║
║ Reporter:    @jane.smith      │                                              ║
║ Created:     2025-01-15       │ Acceptance Criteria:                         ║
║ Updated:     2 minutes ago    │ ☑ Users can login with email/password       ║
║                               │ ☑ Tokens expire after 24 hours              ║
║ LABELS                        │ ☐ Session persists across browser restart   ║
║ ─────────────────────         │ ☐ Failed login shows appropriate error      ║
║ [auth] [security] [p0]        │ ☐ Rate limiting prevents brute force        ║
║                               │                                              ║
║ LINKS                         │ Custom Fields:                               ║
║ ─────────────────────         │ ┌──────────────────────────────────────────┐ ║
║ Blocks:      FRONT-200        │ │Component:    Backend API                 │ ║
║ Blocked by:  None             │ │Environment:  Production                  │ ║
║ Related:     BACK-150         │ │Risk Level:   Medium                      │ ║
║                               │ └──────────────────────────────────────────┘ ║
║ ATTACHMENTS                   │                                              ║
║ ─────────────────────         │ Comments (3):                                ║
║ 📎 login-flow.png (232 KB)    │ ┌──────────────────────────────────────────┐ ║
║ 📎 api-spec.yaml (18 KB)      │ │@john.doe - 2h ago                        │ ║
║                               │ │Started implementation, JWT part done     │ ║
║                               │ ├──────────────────────────────────────────┤ ║
║                               │ │@jane.smith - 1h ago                      │ ║
║                               │ │Please ensure refresh tokens are included │ ║
║                               │ └──────────────────────────────────────────┘ ║
╠═══════════════════════════════╧══════════════════════════════════════════════╣
║ [Tab]Next Field [S-Tab]Prev [Ctrl-S]Save [Ctrl-C]Cancel [F1]Help │ EDIT MODE ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

**Design Features:**
- **Split pane layout** - Context on left, details on right
- **Inline editing** - Direct field modification
- **Visual hierarchy path** - See where ticket fits
- **Checkbox criteria** - Track completion
- **Comment thread** - Collaboration history

---

## 3. Search & Command Palette

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Ticketr  Project: backend  Search Mode                                       ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  ┌────────────────────────────────────────────────────────────────────┐     ║
║  │ 🔍 auth login jwt                                               │ × │     ║
║  └────────────────────────────────────────────────────────────────────┘     ║
║                                                                              ║
║  Search Results (7 matches)                          Sort: Relevance ▼       ║
║  ─────────────────────────────────────────────────────────────────────      ║
║                                                                              ║
║  ► BACK-102  📘 Login Flow Implementation              95% │ In Progress    ║
║    └─ Matches: title:"login" description:"JWT token generation"             ║
║                                                                              ║
║  ► BACK-204  ⚡ Implement JWT token authentication     88% │ Done           ║
║    └─ Matches: title:"JWT" description:"auth tokens"                       ║
║                                                                              ║
║  ► BACK-100  🎯 Authentication System                  73% │ In Progress    ║
║    └─ Matches: description:"login" children:3 matches                      ║
║                                                                              ║
║  ► BACK-205  ⚡ Session management with JWT            67% │ To Do          ║
║    └─ Matches: title:"JWT" labels:"auth"                                   ║
║                                                                              ║
║  ► BACK-315  🐛 JWT token expiry not working          45% │ Open           ║
║    └─ Matches: title:"JWT" type:bug                                        ║
║                                                                              ║
║  Quick Filters:                                                             ║
║  ──────────────────                                                         ║
║  [@] Assignee  [#] ID  [!] Priority  [~] Sprint  [%] Progress  [:] Command ║
║                                                                              ║
║  Examples:                                                                  ║
║  • @john.doe !high ~sprint23    - John's high priority in sprint 23        ║
║  • #BACK-102                     - Go directly to ticket                    ║
║  • :push                         - Execute push command                     ║
║  • /blocked/                     - Regex search in descriptions             ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [Enter]Open [Ctrl-O]Split [Ctrl-E]Edit [ESC]Close [↑↓]Navigate │ 7 matches  ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

**Search Features:**
- **Fuzzy matching** with relevance scores
- **Match highlighting** showing where terms found
- **Quick filters** with special operators
- **Command mode** with `:` prefix
- **Regex support** with `/pattern/`

---

## 4. Conflict Resolution View

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Ticketr  Conflict Resolution  BACK-102: Login Flow                          ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  ⚠️  CONFLICT DETECTED - Both local and remote have changed                 ║
║                                                                              ║
║  ┌─────────────────────┬─────────────────────┬─────────────────────┐       ║
║  │ LOCAL (Your Change) │ REMOTE (JIRA)       │ MERGED (Proposed)   │       ║
║  ├─────────────────────┼─────────────────────┼─────────────────────┤       ║
║  │ Title:              │ Title:              │ Title:              │       ║
║  │ Login Flow Impl.    │ Login Flow          │ Login Flow Impl.    │       ║
║  │                     │                     │ ✏️                   │       ║
║  ├─────────────────────┼─────────────────────┼─────────────────────┤       ║
║  │ Status:             │ Status:             │ Status:             │       ║
║  │ In Progress        │ In Review           │ In Review          │       ║
║  │                     │ ◀─── Changed        │ ✏️                   │       ║
║  ├─────────────────────┼─────────────────────┼─────────────────────┤       ║
║  │ Description:        │ Description:        │ Description:        │       ║
║  │ As a user, I want  │ As a user, I want  │ As a user, I want  │       ║
║  │ to securely log    │ to securely log    │ to securely log    │       ║
║  │ into the app using │ into the app using │ into the app using │       ║
║  │ my credentials.    │ email/password.    │ my credentials.    │       ║
║  │                     │                     │                     │       ║
║  │ Implementation:     │ Implementation:     │ Implementation:     │       ║
║  │ - JWT tokens       │ - JWT tokens       │ - JWT tokens       │       ║
║  │ - Session mgmt     │ - Session mgmt     │ - Session mgmt     │       ║
║  │ - Remember me      │ - Rate limiting    │ - Remember me      │       ║
║  │ ◀─── Added         │ ◀─── Added         │ - Rate limiting    │       ║
║  │                     │                     │ ✏️                   │       ║
║  ├─────────────────────┼─────────────────────┼─────────────────────┤       ║
║  │ Story Points: 8     │ Story Points: 13    │ Story Points: 13    │       ║
║  │                     │ ◀─── Changed        │ ✏️                   │       ║
║  └─────────────────────┴─────────────────────┴─────────────────────┘       ║
║                                                                              ║
║  Resolution Strategy:                                                        ║
║  ┌────────────────────────────────────────────────────────────────┐        ║
║  │ ( ) Keep Local - Overwrite remote with your changes            │        ║
║  │ ( ) Keep Remote - Discard your changes                         │        ║
║  │ (•) Manual Merge - Review and edit each field                  │        ║
║  │ ( ) Create Duplicate - Keep both as separate tickets           │        ║
║  └────────────────────────────────────────────────────────────────┘        ║
║                                                                              ║
║  [Apply to All Conflicts] □                          [3 more conflicts]     ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [←→]Switch Column [↑↓]Navigate [e]Edit [a]Accept [s]Skip [Q]Abort │ 1 of 4  ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

**Conflict Resolution Features:**
- **Three-way diff** visualization
- **Field-level granularity**
- **Visual indicators** for changes (◀─── markers)
- **Edit capability** in merged column
- **Batch operations** for multiple conflicts

---

## 5. Sync Status Dashboard

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Ticketr  Sync Operations  Live ●                                             ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  Current Operation: PUSH to backend (BACK)                                  ║
║  ════════════════════════════════════════════════════════════════           ║
║                                                                              ║
║  Progress: ████████████████████░░░░░░░░░░  67% (8/12 tickets)               ║
║                                                                              ║
║  ┌─────────────────────────────────────────────────────────────────────┐    ║
║  │ ✅ BACK-102  Login Flow                    Updated      223ms       │    ║
║  │ ✅ BACK-204  JWT Implementation            Created      456ms       │    ║
║  │ ✅ BACK-205  Session Management            Updated      189ms       │    ║
║  │ ✅ BACK-206  Remember Me Feature           No changes  12ms        │    ║
║  │ ✅ BACK-207  Password Reset                Updated      234ms       │    ║
║  │ ✅ BACK-208  Email Templates               Created      567ms       │    ║
║  │ ✅ BACK-209  Rate Limiting                 Updated      198ms       │    ║
║  │ ✅ BACK-210  2FA Support                   Updated      345ms       │    ║
║  │ ⚡ BACK-211  OAuth Integration             Pushing...   1.2s        │    ║
║  │ ⏳ BACK-212  SAML Support                  Queued                   │    ║
║  │ ⏳ BACK-213  Security Audit                Queued                   │    ║
║  │ ⏳ BACK-214  Performance Testing           Queued                   │    ║
║  └─────────────────────────────────────────────────────────────────────┘    ║
║                                                                              ║
║  Statistics:                           Rate Limit:                          ║
║  ─────────────                        ────────────                          ║
║  Created:     2 tickets                API Calls:   45/100                  ║
║  Updated:     6 tickets                Reset in:    42s                     ║
║  Unchanged:   1 ticket                 Rate:        2.3 req/s               ║
║  Failed:      0 tickets                                                     ║
║  Conflicts:   0 detected               Network:                             ║
║                                        ────────────                          ║
║  Time Elapsed:    19.5s                Latency:     234ms avg               ║
║  Time Remaining:  ~9s                  Bandwidth:   45 KB/s                 ║
║                                                                              ║
║  Recent Errors:                                                             ║
║  ──────────────                                                             ║
║  None                                                                       ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [P]Pause [C]Cancel [V]Verbose [L]Show Log [R]Retry Failed │ Auto-retry: ON  ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## 6. Project Switcher (Modal)

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                           Switch Project                                     ║
║  ┌────────────────────────────────────────────────────────────────────┐     ║
║  │                                                                    │     ║
║  │  Select Project:                                     [Quick: 1-9] │     ║
║  │  ──────────────────────────────────────────────────────────────── │     ║
║  │                                                                    │     ║
║  │  1 ● backend        BACK   company.atlassian.net    23 tickets   │     ║
║  │      Last sync: 2 minutes ago                       ▲ current    │     ║
║  │                                                                    │     ║
║  │  2 ○ frontend       FRONT  company.atlassian.net    45 tickets   │     ║
║  │      Last sync: 1 hour ago                                       │     ║
║  │                                                                    │     ║
║  │  3 ○ mobile         MOB    company.atlassian.net    12 tickets   │     ║
║  │      Last sync: 3 days ago                         ⚠ outdated   │     ║
║  │                                                                    │     ║
║  │  4 ○ personal       PERS   personal.atlassian.net   8 tickets    │     ║
║  │      Last sync: 1 week ago                                       │     ║
║  │                                                                    │     ║
║  │  + Create New Project...                                         │     ║
║  │                                                                    │     ║
║  └────────────────────────────────────────────────────────────────────┘     ║
║                                                                              ║
║  [Enter]Select [N]New [D]Delete [E]Edit [ESC]Cancel                        ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## 7. Create Ticket Wizard

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Create New Ticket  Step 2 of 4: Details                                     ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  Hierarchy:                                                                 ║
║  ─────────────────────────────────────────────────────────────              ║
║  Parent: BACK-100 Authentication System (Epic)                              ║
║  Type:   Story                                                              ║
║                                                                              ║
║  Basic Information:                                                         ║
║  ─────────────────────────────────────────────────────────────              ║
║                                                                              ║
║  Title: ┌──────────────────────────────────────────────────────────┐       ║
║         │Password Recovery Flow                                    │       ║
║         └──────────────────────────────────────────────────────────┘       ║
║                                                                              ║
║  Description:                                                               ║
║  ┌────────────────────────────────────────────────────────────────────┐     ║
║  │As a user who has forgotten their password                         │     ║
║  │I want to reset it via email                                       │     ║
║  │So that I can regain access to my account                          │     ║
║  │                                                                    │     ║
║  │Implementation notes:                                               │     ║
║  │- Send reset link to registered email                              │     ║
║  │- Link expires after 1 hour                                        │     ║
║  │- Require strong password on reset                                 │     ║
║  └────────────────────────────────────────────────────────────────────┘     ║
║                                                                              ║
║  Fields:                                      Templates:                    ║
║  ─────────────────────                       ──────────────                 ║
║  Priority:    [High        ▼]                [Load Template ▼]              ║
║  Sprint:      [Sprint 24   ▼]                • User Story                  ║
║  Points:      [5           ]                 • Bug Report                  ║
║  Assignee:    [@john.doe   ▼]                • Technical Task              ║
║  Component:   [Backend API ▼]                • Feature Epic                ║
║                                                                              ║
║  ☑ Create subtasks after saving                                            ║
║  ☐ Open in editor after creation                                           ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [Tab]Next [S-Tab]Prev [Ctrl-Enter]Create [F2]Template [ESC]Cancel │ Step 2/4║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## Design System & Components

### Color Palette

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Semantic Colors:                                                             ║
║ ────────────────                                                             ║
║ Success/Done:     ██ Green   (#10B981)  ✅ ✓                                ║
║ In Progress:      ██ Yellow  (#F59E0B)  🔄 ⚡                               ║
║ Todo/Pending:     ██ Gray    (#6B7280)  ⏸ ○                                ║
║ Error/Blocked:    ██ Red     (#EF4444)  🔴 ✗                                ║
║ Info/Epic:        ██ Blue    (#3B82F6)  🎯                                  ║
║ Story:            ██ Cyan    (#06B6D4)  📘                                  ║
║ Bug:              ██ Orange  (#FB923C)  🐛                                  ║
║                                                                              ║
║ UI Elements:                                                                 ║
║ ────────────                                                                 ║
║ Border:           ══ ─ │ ├ └ ┌ ┐ ┘ ┤ ┬ ┴ ┼                                ║
║ Selected:         ▶ ● ◆                                                     ║
║ Unselected:       ▷ ○ ◇                                                     ║
║ Progress:         ████████░░░░                                              ║
║ Tree:             └─ ├─ │                                                   ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Keyboard Navigation Map

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Global:           Navigation:        Actions:           Modes:              ║
║ ────────          ───────────        ────────           ──────              ║
║ ? - Help          j/↓ - Down         c - Create         / - Search          ║
║ q - Quit          k/↑ - Up           e - Edit           : - Command         ║
║ Q - Force quit    h/← - Left/Close   d - Delete         i - Insert          ║
║ : - Command       l/→ - Right/Open   p - Push           v - Visual          ║
║ / - Search        g - Top            P - Pull           ESC - Normal        ║
║ Tab - Next pane   G - Bottom         r - Refresh                            ║
║ 1-9 - Projects    f - Page down      s - Save                               ║
║                   b - Page up        y - Yank/Copy                          ║
║                   H/M/L - High/Mid   x - Cut                                ║
║                   ^ - First          . - Repeat                             ║
║                   $ - Last           u - Undo                               ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## Responsive Behavior

### Terminal Size Adaptation

**Full Size (>120 cols)**: All panels visible
**Medium (80-120 cols)**: Hide statistics, abbreviate
**Small (<80 cols)**: Single panel, toggle mode

### Progressive Enhancement

1. **Minimal Mode** - Essential info only
2. **Standard Mode** - Full features
3. **Rich Mode** - Animations, graphs (if supported)

---

## Animation & Feedback

### Loading States
```
Searching   [▉▉▉▉▉░░░░░] 50%
Syncing     ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏ (spinner)
Processing  ••••••••-- (dots)
```

### Status Indicators
```
✓ Success   ⚡ Active    ⏸ Paused
✗ Failed    🔄 Syncing   ⚠ Warning
● Modified  ○ Unchanged  ◐ Partial
```

---

## Accessibility Considerations

1. **High Contrast Mode** - Clear boundaries, no color-only info
2. **Screen Reader Support** - Semantic structure
3. **Keyboard Only** - No mouse required
4. **Status Announcements** - Audio cues for operations
5. **Customizable Keys** - Remappable bindings

---

This wireframe collection demonstrates a professional, powerful TUI that rivals the best terminal applications, making Ticketr not just functional but genuinely pleasant to use.