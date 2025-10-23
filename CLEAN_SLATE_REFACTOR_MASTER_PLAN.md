# Ticketr Bubbletea TUI: Clean-Slate Refactor Master Plan

**Status:** DEFINITIVE PLAN
**Type:** CLEAN SLATE - No backward compatibility required
**Timeline:** 12 weeks aggressive (vs 18 with migration overhead)
**Philosophy:** REFACTOR TO DESTINY (DENSITY)

---

## Executive Summary

### The Decision: CLEAN SLATE

**CRITICAL CONTEXT CHANGES:**
- NO parallel Tview development - We're going ALL IN on Bubbletea
- NO migration layer - We CAN break things (feature branch isolation)
- NO incremental cutover - Big bang merge when ready
- YES to aggressive timeline - No migration overhead saves 6 weeks
- YES to breaking changes - Complete architectural refresh

**Why Clean Slate?**

1. **Simpler**: No compatibility shims, no dual codebases, no feature flags
2. **Faster**: 12 weeks vs 18 weeks (33% time savings)
3. **Better Quality**: Design from scratch with lessons learned
4. **Lower Risk**: Feature branch until perfect, atomic switch
5. **Team Alignment**: Full team focus on new architecture

**Rollback Strategy:**
- It's a git branch. That's it.
- If we need to roll back: `git checkout main`
- No complex migration state to unwind

**When to Merge:**
When ALL quality gates pass (see Quality Gates section). No pressure to ship incrementally.

---

## 1. New Architecture Vision

### Philosophy

**Midnight Commander Meets Modern TUI**
- Dual-panel layout (ticket tree + detail view)
- F-key actions (F1=Help, F2=Sync, F3=Workspace, etc.)
- High information density (60 FPS rendering, no flicker)
- Space-themed cosmic effects (optional, can disable)
- Vim-like keybindings (hjkl navigation)
- Context-aware action bar

**Design Principles**
1. **Declarative over Imperative** - Describe UI, don't manage updates
2. **Messages over Callbacks** - Everything is message-driven
3. **Composition over Inheritance** - Small, composable components
4. **Performance by Default** - 60 FPS, <16ms render budget
5. **Accessibility First** - Screen reader support, high contrast mode

### Directory Structure

```
internal/tui-bubbletea/           # NEW - All Bubbletea code
├── app.go                        # Root Bubbletea program + setup
├── model.go                      # Root model (global state)
├── update.go                     # Root update function (message router)
├── view.go                       # Root view function (layout)
│
├── components/                   # Reusable UI components
│   ├── tree/                     # Hierarchical ticket tree
│   │   ├── model.go              # Tree component model
│   │   ├── update.go             # Tree update logic
│   │   ├── view.go               # Tree rendering
│   │   ├── delegate.go           # List item delegate for tree nodes
│   │   └── flatten.go            # Tree flattening + virtualization
│   │
│   ├── actionbar/                # Context-aware action bar
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── marquee.go            # Scrolling action hints
│   │
│   ├── modal/                    # Modal overlay system
│   │   ├── model.go              # Generic modal container
│   │   ├── confirm.go            # Confirmation dialogs
│   │   ├── input.go              # Input prompts
│   │   └── overlay.go            # Overlay positioning
│   │
│   ├── statusbar/                # Header status bar
│   │   ├── model.go
│   │   └── view.go
│   │
│   └── panel/                    # Reusable panel component
│       ├── model.go
│       └── view.go
│
├── views/                        # Screen-level views
│   ├── workspace/                # Workspace selector (slide-out)
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   │
│   ├── tickets/                  # Main dual-panel view
│   │   ├── model.go              # Coordinates tree + detail
│   │   ├── update.go
│   │   └── view.go
│   │
│   ├── detail/                   # Ticket detail panel
│   │   ├── model.go              # Display mode state
│   │   ├── edit.go               # Edit mode state
│   │   ├── update.go
│   │   └── view.go
│   │
│   ├── search/                   # Fuzzy search modal
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   │
│   └── help/                     # Help screen
│       ├── model.go
│       └── view.go
│
├── actions/                      # Extensible action system
│   ├── action.go                 # Action type definitions
│   ├── context.go                # Context manager
│   ├── registry.go               # Action registry
│   ├── resolver.go               # Keybinding resolver
│   ├── executor.go               # Execution pipeline
│   │
│   ├── predicates/               # Action predicates
│   │   └── predicates.go
│   │
│   └── builtin/                  # Built-in actions
│       ├── tickets.go            # Ticket actions (open, edit, delete)
│       ├── workspaces.go         # Workspace actions
│       ├── navigation.go         # Navigation actions
│       ├── sync.go               # Sync actions
│       └── system.go             # System actions (quit, help)
│
├── theme/                        # Theme system (Lipgloss styles)
│   ├── theme.go                  # Theme loader
│   ├── colors.go                 # Color palette
│   ├── styles.go                 # Style definitions
│   └── default.go                # Default theme
│
├── effects/                      # Visual effects (optional)
│   ├── spinner.go                # Loading spinners
│   ├── progress.go               # Progress bars
│   ├── shimmer.go                # Shimmer effect
│   └── particles.go              # Particle effects (stretch goal)
│
├── messages/                     # Message type definitions
│   ├── tickets.go                # Ticket-related messages
│   ├── workspace.go              # Workspace messages
│   ├── sync.go                   # Sync messages
│   └── ui.go                     # UI state messages
│
└── layout/                       # Layout utilities
    ├── flexbox.go                # Flexbox layout (via stickers)
    ├── responsive.go             # Responsive sizing
    └── measure.go                # Content measurement
```

### Component Library Strategy

Based on ecosystem research:

**Use These Libraries:**
1. **Bubbletea** - Core framework (TEA pattern)
2. **Lipgloss** - ALL styling (no inline styles)
3. **Bubbles** - Official components:
   - `list` - Base for tree view (custom delegate)
   - `viewport` - Detail view scrolling
   - `textinput` - Search, filters
   - `textarea` - Edit mode
   - `help` - Context-sensitive help
   - `spinner` - Loading states
   - `progress` - Sync progress
4. **Stickers** - Layout management (FlexBox for dual panels)
5. **Huh** - Forms (workspace creation, bulk operations)

**Build Custom:**
1. **Tree Component** - No good library exists for 1000+ hierarchical items
   - Base: `bubbles/list` with custom delegate
   - Virtualization: Render only visible items (50-100)
   - Flattening: Dynamic tree → flat list on expand/collapse
2. **Action System** - Extensible, plugin-ready (from design doc)
3. **Modal System** - Overlay positioning (simple, use Lipgloss PlaceOverlay)

**Skip/Decide Later:**
- **Particle effects** - Cool but optional (3-week effort, can add later)
- **Teacup** - Inactive library, don't use

---

## 2. The 12-Week Plan

### Phase 1: Foundation (Weeks 1-2)

**Goal:** Scaffold architecture, prove out patterns

**Week 1: Structure & Core Types**
- [ ] Create `internal/tui-bubbletea/` directory
- [ ] Define root Model struct (global state)
- [ ] Define message types (tickets, workspace, sync, ui)
- [ ] Implement theme system (Lipgloss styles)
  - Color palette (adaptive colors for light/dark)
  - Component styles (panels, borders, text)
  - Dynamic sizing helpers
- [ ] Setup Layout system (Stickers FlexBox)
  - Dual-panel layout (50/50 split)
  - Header bar (3 rows)
  - Action bar (3 rows)
  - Responsive sizing on WindowSizeMsg
- [ ] Write initial tests (model, themes, layout)

**Week 2: Simple View Proof-of-Concept**
- [ ] Implement HelpView (simplest view)
  - Static content rendering
  - Viewport scrolling
  - Keyboard navigation (q to quit, j/k scroll)
- [ ] Implement StatusBar component
  - Sync status indicator (animated spinner)
  - Workspace name
  - Ticket count
- [ ] Wire up root app.go
  - Tea.Program initialization
  - WindowSizeMsg handling
  - Basic message routing
- [ ] Test in terminal (manual QA)
  - Resize handling
  - Keyboard input
  - Rendering performance

**Deliverable:** Running TUI that shows help screen with proper layout

---

### Phase 2: Core Views (Weeks 3-6)

**Goal:** Implement primary user flows (workspace → tree → detail)

**Week 3: Workspace Selector**
- [ ] Build WorkspaceListView
  - Use `bubbles/list` component
  - Slide-out animation (150ms from left)
  - Backdrop dimming effect
  - Search/filter input
- [ ] Integrate with workspace query layer
  - Load workspaces from DB
  - Display metadata (ticket count, last sync)
- [ ] Implement selection & switch
  - Message: `workspaceSelectedMsg`
  - Update root model state
  - Trigger ticket tree load

**Week 4-5: Ticket Tree (HARD)**
- [ ] Design custom tree component
  - Base: `bubbles/list` with custom delegate
  - TreeItem type (ID, level, hasChildren, expanded)
  - Flattening algorithm (only expanded nodes visible)
  - Virtualization (render max 100 visible items)
- [ ] Implement tree delegate
  - Indentation rendering (strings.Repeat("  ", level))
  - Expand/collapse icons (▶ collapsed, ▼ expanded)
  - Status indicators (● synced, ○ local, □ task)
  - Selection checkboxes ([x] selected, [ ] not)
- [ ] Implement keybindings
  - j/k - Navigate up/down
  - h/l - Collapse/expand
  - Enter - Open detail view
  - Space - Toggle selection
  - gg/G - Jump to top/bottom
- [ ] Integrate with ticket query layer
  - Load tickets by workspace
  - Build tree structure (subtask hierarchy)
  - Track expansion state (map[string]bool)
- [ ] Test with 1000+ tickets
  - Measure render time (<16ms target)
  - Profile memory usage
  - Ensure smooth scrolling

**Week 6: Ticket Detail View**
- [ ] Build TicketDetailView (display mode)
  - Use `bubbles/viewport` for scrolling
  - Render ticket fields (ID, summary, status, priority, etc.)
  - Render description (Markdown formatting via Glamour)
  - Render custom fields
  - Render acceptance criteria checkboxes
- [ ] Build edit mode (future: Week 7)
  - Switch between display/edit on 'e' key
  - Use `huh` forms for field editing
  - Dirty state tracking
  - Save confirmation
- [ ] Implement focus management
  - Tree focused: Double border, green
  - Detail focused: Double border, green
  - Tab to switch focus
- [ ] Wire up ticket selection
  - Message: `ticketSelectedMsg`
  - Load ticket from service
  - Display in detail panel

**Deliverable:** Functional workspace → tree → detail flow

---

### Phase 3: Modals & Search (Weeks 7-9)

**Goal:** Secondary UX features for productivity

**Week 7: Search Modal**
- [ ] Build SearchView modal
  - Center overlay (60% width, 40% height)
  - Backdrop dim + blur effect
  - Search input (`bubbles/textinput`)
  - Results list (`bubbles/list`)
- [ ] Implement fuzzy search
  - Integrate existing fuzzy match logic
  - Real-time filtering as user types
  - Score-based sorting
  - Match highlighting (yellow)
- [ ] Implement filter syntax
  - `@user` - assignee filter
  - `!priority` - priority filter
  - `#ID` - ticket ID
  - `~sprint` - sprint filter
- [ ] Open from tree view
  - Press `/` to open
  - Select result → close modal, open ticket
  - Esc to cancel

**Week 8: Command Palette**
- [ ] Build CommandPaletteView
  - Center overlay (50% width, 60% height)
  - Command registry integration
  - Fuzzy search on command names
  - Category headers (Sync, File, View, etc.)
  - Keyboard shortcut hints
- [ ] Register all commands
  - sync:pull, sync:push, sync:full
  - file:open, file:save
  - workspace:switch, workspace:new
  - view:theme, view:layout
  - help:keys, help:about
- [ ] Implement execution
  - Press `:` or `Ctrl+P` to open
  - Enter to execute
  - Dispatch appropriate message

**Week 9: Action Bar + Modals**
- [ ] Implement ActionBar component
  - Context-aware keybinding display
  - Dynamic based on focus (tree vs detail)
  - Integration with action registry
  - Marquee scrolling if too wide (optional)
- [ ] Build modal system
  - Generic modal container
  - Confirmation dialogs (delete, discard changes)
  - Input prompts (quick ticket create)
  - Bulk operations modal (simplified for MVP)

**Deliverable:** Full modal + search experience

---

### Phase 4: Async & Job Integration (Weeks 10-11)

**Goal:** Jira sync, background jobs, progress tracking

**Week 10: Job Queue Integration**
- [ ] Adapt job queue to tea.Cmd pattern
  - Convert job.Start() to return tea.Cmd
  - Job progress as messages (jobProgressMsg)
  - Job completion as messages (jobCompleteMsg, jobErrorMsg)
- [ ] Implement progress UI
  - Progress bar in action bar
  - Spinner in sync status indicator
  - ETA calculation
  - Cancel support (Ctrl+C during sync)
- [ ] Wire up pull operation
  - Action: sync:pull
  - Show progress modal
  - Update tree on completion
  - Error handling + notification

**Week 11: Push & Full Sync**
- [ ] Implement push operation
  - Collect dirty tickets
  - Show confirmation modal (X tickets to push)
  - Progress tracking
  - Conflict resolution UI (if needed)
- [ ] Implement full sync
  - Combines pull + push
  - Two-phase progress (pull → push)
  - Optimistic UI updates
- [ ] Background sync coordinator
  - Auto-sync every N minutes (config)
  - Notifications on completion
  - Silent mode (no modals)

**Deliverable:** Full Jira sync integration

---

### Phase 5: Polish & Effects (Weeks 12-14)

**Goal:** Visual polish, performance tuning, nice-to-haves

**Week 12: Basic Effects**
- [ ] Spinners for loading states
  - Use `bubbles/spinner`
  - Sync in progress
  - Loading ticket details
  - Background operations
- [ ] Progress bars
  - Use `bubbles/progress`
  - Sync progress
  - Bulk operation progress
- [ ] Shimmer effect on progress bars
  - Character substitution animation
  - Configurable via theme
- [ ] Smooth transitions
  - Panel focus transitions (color fade)
  - Modal fade in/out
  - Slide-out panel animation

**Week 13-14: Advanced Effects (OPTIONAL)**
- [ ] **DECISION POINT:** Particle effects or skip?
  - If YES: 2 weeks to implement ANSI layer-based particle system
  - If NO: Use time for extra polish & testing
- [ ] Gradient effects
  - Header gradient
  - Status indicator gradients
- [ ] Themes
  - Default theme (green accent)
  - Additional themes (blue, purple, midnight)
  - Theme switcher in command palette

**Deliverable:** Polished, visually appealing TUI

---

### Phase 6: Testing & Hardening (Weeks 15-16) - EXTENDED

**Note:** Original plan was 12 weeks. Adding 4 weeks buffer for testing/hardening brings us to 16 weeks total, still under 18-week migration estimate.

**Week 15: Comprehensive Testing**
- [ ] Unit tests for all components
  - Model Update() logic
  - View() rendering (golden files via teatest)
  - Predicate functions
  - Message handling
- [ ] Integration tests
  - Full user flows (teatest)
  - Workspace → Tree → Detail
  - Search → Select → Open
  - Sync operation end-to-end
- [ ] Performance testing
  - Profile with 1000+ tickets
  - Measure render times (<16ms)
  - Memory profiling
  - Stress testing (10K tickets)

**Week 16: Bug Fixing & Documentation**
- [ ] Fix all critical bugs
- [ ] Fix all high-priority bugs
- [ ] Document keybindings (update README)
- [ ] Write migration guide (for users)
- [ ] Update contribution guidelines
- [ ] Create demo video/GIF
- [ ] Final QA pass

**Deliverable:** Production-ready TUI

---

### Phase 7: Merge & Deploy (Week 17)

**Goal:** Merge to main, deploy, celebrate

**Pre-Merge Checklist:**
- [ ] ALL quality gates passed (see section below)
- [ ] No critical or high bugs
- [ ] Performance benchmarks met
- [ ] Code review complete
- [ ] Documentation updated
- [ ] Demo prepared

**Merge Process:**
1. Final rebase from main
2. Resolve any conflicts
3. Run full test suite
4. Create PR (feature/bubbletea-refactor → main)
5. Team review
6. Merge (big bang)
7. Tag release (v0.5.0-bubbletea or v1.0.0)
8. Deploy to users

**Rollback Plan:**
- If critical bugs discovered: Revert merge commit
- Fix on feature branch
- Re-merge when fixed

---

## 3. Quality Gates for Merge

These are HARD REQUIREMENTS. We do NOT merge until all are satisfied.

### Performance Gates

- [ ] **60 FPS rendering**: No dropped frames during normal operation
- [ ] **Render budget**: All View() functions complete in <16ms (60 FPS target)
- [ ] **Tree rendering**: 1000 tickets render in <100ms
- [ ] **Search latency**: Results appear in <200ms
- [ ] **Memory usage**: <50MB RSS with 1000 tickets loaded
- [ ] **Smooth animations**: No jank in panel transitions or modal fades

### Functional Gates

- [ ] **100% feature parity** with essential current TUI features:
  - Workspace selection
  - Ticket tree (expand/collapse, navigate, select)
  - Ticket detail view (display mode)
  - Search (fuzzy, filters)
  - Sync (pull, push, full)
  - Help screen
  - Command palette
  - Keybindings
- [ ] **Vim-like navigation** fully functional (hjkl, gg, G, etc.)
- [ ] **Context-aware actions** working (action bar updates by focus)
- [ ] **No data loss**: All ticket edits save correctly
- [ ] **No sync corruption**: Sync operations maintain data integrity

### Quality Gates

- [ ] **Test coverage**: 80%+ for all new code
- [ ] **All tests passing**: No failing unit or integration tests
- [ ] **No critical bugs**: Zero P0 bugs
- [ ] **<5 high bugs**: Less than 5 P1 bugs (with mitigation plan)
- [ ] **Code review**: All code reviewed by at least one other developer
- [ ] **Documentation**: README updated, keybindings documented

### UX Gates

- [ ] **Keyboard-only navigation**: No mouse required for any operation
- [ ] **Accessibility**: Screen reader friendly (tested)
- [ ] **Responsive**: Works on 120x30 terminals (min size)
- [ ] **No flicker**: Rendering is clean (no screen tearing)
- [ ] **Error handling**: All errors have user-friendly messages
- [ ] **Help available**: '?' key shows context-sensitive help

### Integration Gates

- [ ] **Jira sync working**: Pull and push operations succeed
- [ ] **Job queue integration**: Background jobs complete successfully
- [ ] **Database migrations**: No schema breakage
- [ ] **Config compatibility**: User config files still work (or migration provided)

---

## 4. Risk Assessment & Mitigation

### Risk Matrix

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| **Tree component performance with 1000+ items** | Medium | High | Virtualization + profiling early (Week 5) |
| **Particle effects take 3+ weeks** | Medium | Medium | Make optional, skip if timeline tight |
| **Job queue integration breaks** | Low | High | Thorough testing, keep existing job logic intact |
| **Team velocity lower than estimated** | Medium | High | 4-week buffer built in (Weeks 15-16) |
| **Critical bug discovered post-merge** | Low | Critical | Rollback plan (revert merge commit) |
| **User resistance to new UX** | Low | Medium | Demo early, gather feedback, iterate |
| **Scope creep (new features)** | High | Medium | Strict scope control, new features deferred to v0.6 |

### Mitigation Strategies

**Performance Risk:**
- Profile early and often (Week 5 tree component)
- Set performance budgets (16ms render time)
- Use pprof and benchmarking
- Fallback: Pagination if virtualization insufficient

**Timeline Risk:**
- Daily standups to catch blockers
- Weekly demos to stakeholders
- 4-week buffer for unknowns
- Flexible on optional features (particles, advanced effects)

**Quality Risk:**
- Automated testing from Week 1
- Code review for all PRs
- QA pass before merge
- Beta testing period (if time allows)

**Integration Risk:**
- Keep domain/business logic unchanged
- Only replace TUI adapter layer
- Extensive integration testing
- Rollback plan ready

---

## 5. Development Workflow

### Branch Strategy

```
main (production)
  ↑
  └─ feature/bubbletea-refactor (clean slate)
       ↑
       ├─ feature/tree-component
       ├─ feature/action-system
       ├─ feature/search-modal
       └─ ... (short-lived feature branches)
```

**Rules:**
1. **ALL work** happens on `feature/bubbletea-refactor` or sub-branches
2. **NO merges to main** until Week 17 (full feature complete)
3. **Rebase frequently** from main to stay current
4. **Sub-branches** merge to `feature/bubbletea-refactor` via PR
5. **Delete Tview code** once Bubbletea merged (cleanup PR)

### Daily Workflow

**Morning:**
1. Pull latest from `feature/bubbletea-refactor`
2. Review open PRs
3. Standup (blockers, progress, plan)

**During Day:**
1. Work on assigned component
2. Write tests alongside code
3. Manual testing in terminal
4. Open PR when feature complete

**End of Day:**
1. Commit work (even if not done)
2. Push to remote
3. Update task tracker

### Code Review Process

**PR Requirements:**
- [ ] All tests passing
- [ ] No lint errors
- [ ] Screenshot or video of feature (if UI)
- [ ] Description of changes
- [ ] Link to relevant design docs

**Review Checklist:**
- [ ] Code quality (readability, maintainability)
- [ ] Test coverage (80%+ for new code)
- [ ] Performance (no obvious bottlenecks)
- [ ] UX (matches wireframes, good keybindings)
- [ ] Documentation (if needed)

**Merge Process:**
- 1 approval required (minimum)
- 2 approvals for high-risk changes
- Author merges after approval
- Delete branch after merge

---

## 6. Component Priority Matrix

### Must-Have (Blocker for Merge)

These MUST be implemented and working to merge:

| Component | Complexity | Owner | Status |
|-----------|-----------|-------|--------|
| Root Model & Layout | Medium | TBD | Not Started |
| Theme System | Low | TBD | Not Started |
| Workspace Selector | Medium | TBD | Not Started |
| Ticket Tree | **High** | TBD | Not Started |
| Ticket Detail (Display) | Medium | TBD | Not Started |
| Action System | High | TBD | Not Started |
| Search Modal | Medium | TBD | Not Started |
| Sync Integration (Pull/Push) | High | TBD | Not Started |
| Help Screen | Low | TBD | Not Started |
| Action Bar | Medium | TBD | Not Started |

### Should-Have (Desired for Merge)

These are highly desired but can be deferred if timeline tight:

| Component | Complexity | Owner | Status |
|-----------|-----------|-------|--------|
| Command Palette | Medium | TBD | Not Started |
| Ticket Detail (Edit Mode) | High | TBD | Not Started |
| Bulk Operations Modal | Medium | TBD | Not Started |
| Progress Bars | Low | TBD | Not Started |
| Spinners | Low | TBD | Not Started |
| Modal System (Generic) | Medium | TBD | Not Started |

### Nice-to-Have (Post-Merge)

These can ship in v0.6 or later:

| Component | Complexity | Owner | Status |
|-----------|-----------|-------|--------|
| Particle Effects | **Very High** | TBD | Deferred |
| Advanced Themes | Low | TBD | Deferred |
| Marquee Animation | Medium | TBD | Deferred |
| Gradient Effects | Low | TBD | Deferred |
| Shimmer Effect | Low | TBD | Optional |
| Custom Borders | Low | TBD | Optional |

---

## 7. Success Metrics

### Quantitative Metrics

**Performance:**
- Render time (avg): <10ms (target <16ms)
- Tree load time (1000 tickets): <100ms
- Search latency: <200ms
- Memory usage: <50MB RSS
- Frame rate: 60 FPS sustained

**Quality:**
- Test coverage: >80%
- Critical bugs: 0
- High bugs: <5
- Code review coverage: 100%

**Timeline:**
- Merge date: Week 17 or earlier
- Feature complete: Week 14
- Testing complete: Week 16

### Qualitative Metrics

**User Satisfaction:**
- "This is MUCH smoother than the old TUI" (positive feedback)
- "I can navigate faster with vim keys" (productivity improvement)
- "The search is so much better" (UX win)

**Developer Satisfaction:**
- "The code is so much cleaner" (maintainability)
- "Adding new features is easy now" (extensibility)
- "Tests give me confidence" (quality)

**Product:**
- "This feels modern and professional" (polish)
- "60 FPS is silky smooth" (performance)
- "I love the action system" (UX innovation)

---

## 8. Communication Plan

### Weekly Demo

**Every Friday:**
- Demo progress to stakeholders
- Show new features in terminal
- Gather feedback
- Adjust plan if needed

### Status Updates

**Every Monday:**
- Written status update (what shipped last week)
- Blockers identified
- Plan for current week
- Risks/concerns

### Decision Points

**Key Decisions (Require Team Consensus):**
1. **Week 5:** Tree component architecture finalized
2. **Week 9:** Action system API frozen
3. **Week 13:** Particle effects - implement or skip?
4. **Week 14:** Feature freeze for testing phase
5. **Week 16:** Merge decision (all gates passed?)

---

## 9. The Path to Destiny

### Vision Statement

> "Ticketr's Bubbletea TUI will be the DEFINITIVE terminal Jira client.
> Faster, smoother, more powerful, and more extensible than anything before.
> This is our DESTINY."

### What Makes This DESTINY-tier?

1. **Performance** - 60 FPS, no compromises
2. **UX** - Midnight Commander meets modern design
3. **Extensibility** - Action system ready for plugins
4. **Quality** - 80%+ test coverage, zero critical bugs
5. **Maintainability** - Clean architecture, no tech debt
6. **Accessibility** - Keyboard-first, screen reader friendly
7. **Beauty** - Cosmic effects, smooth animations, Lipgloss styling

### Success Looks Like

**Week 17:**
- Feature branch merges to main
- All quality gates GREEN
- Team celebrates with demo
- Users love the new experience
- Code is clean and tested
- Documentation is complete

**3 Months Later:**
- No major bugs reported
- Users rave about performance
- Contributors can add features easily
- Plugin system ready for extensions
- Ticketr is the #1 Jira TUI

### Beyond the Refactor

**Future (v0.6+):**
- Lua plugin system (action extensions)
- Particle effects (if skipped)
- Additional themes (community contributions)
- Mobile/tablet support (SSH-based)
- Cloud sync (cross-device state)
- AI-powered ticket search
- Integration with other tools (GitHub, Slack)

---

## 10. Appendix: Key Files Reference

### Research Documents
- `/home/karol/dev/private/ticktr/COMPONENT_ECOSYSTEM_DEEP_DIVE.md` - Library analysis
- `/home/karol/dev/private/ticktr/BUBBLETEA_ARCHITECTURE_RESEARCH.md` - Bubbletea patterns
- `/home/karol/dev/private/ticktr/TICKETR_CURRENT_ARCHITECTURE_ANALYSIS.md` - Current TUI analysis
- `/home/karol/dev/private/ticktr/EXTENSIBLE_ACTION_SYSTEM_DESIGN.md` - Action system design
- `/home/karol/dev/private/ticktr/TUI_WIREFRAMES.md` - UI wireframes

### Current TUI Code (Reference Only)
- `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go` - Current app (DO NOT MODIFY)
- `/home/karol/dev/private/ticktr/internal/adapters/tui/views/` - Current views (REFERENCE)
- `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/` - Current widgets (REFERENCE)
- `/home/karol/dev/private/ticktr/internal/adapters/tui/effects/` - Current effects (REFERENCE)

### New Code (To Be Created)
- `/home/karol/dev/private/ticktr/internal/tui-bubbletea/` - ALL NEW CODE GOES HERE

---

## 11. Final Checklist (Pre-Merge)

Before merging `feature/bubbletea-refactor` to `main`:

### Code Quality
- [ ] All code reviewed
- [ ] No TODO comments left
- [ ] No debug print statements
- [ ] All tests passing
- [ ] 80%+ test coverage
- [ ] No lint warnings
- [ ] No known race conditions

### Functionality
- [ ] All must-have features implemented
- [ ] All should-have features implemented (or deferred with plan)
- [ ] All keybindings working
- [ ] All views rendering correctly
- [ ] Sync operations working
- [ ] Search working
- [ ] Help screen complete

### Performance
- [ ] Profiled with pprof
- [ ] 60 FPS sustained
- [ ] <16ms render time
- [ ] <50MB memory usage
- [ ] Tree loads in <100ms

### Quality
- [ ] Zero critical bugs
- [ ] <5 high bugs (with fixes planned)
- [ ] No data corruption issues
- [ ] Error handling robust
- [ ] Logging appropriate

### Documentation
- [ ] README updated
- [ ] Keybindings documented
- [ ] Migration guide written (if needed)
- [ ] Code comments for complex logic
- [ ] Changelog updated

### User Experience
- [ ] Tested on multiple terminal emulators (iTerm, Alacritty, etc.)
- [ ] Tested at min terminal size (120x30)
- [ ] Keyboard-only navigation works
- [ ] Help screen accessible
- [ ] Error messages user-friendly

### Integration
- [ ] Jira sync working
- [ ] Database queries working
- [ ] Config file loading
- [ ] Logging working
- [ ] Exit handling clean

### Pre-Release
- [ ] Demo recorded
- [ ] Team sign-off
- [ ] Stakeholder approval
- [ ] Release notes drafted
- [ ] Version bumped (v1.0.0 or v0.5.0)

---

## Conclusion

This is THE plan. Clean slate. 12 weeks aggressive timeline (16 with buffer). No compromises on quality. No backward compatibility shackles. Just excellent, modern, performant Bubbletea TUI architecture.

**Let's build something LEGENDARY.**

**Next Steps:**
1. Team review of this plan
2. Assign owners to components
3. Create `feature/bubbletea-refactor` branch
4. Start Week 1 tasks
5. Ship in 12-16 weeks

---

**Document Version:** 1.0
**Last Updated:** 2025-10-22
**Status:** APPROVED PENDING TEAM REVIEW
**Next Review:** Week 4 (Mid-Phase 2)
