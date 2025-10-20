# Phase 6: The Enchantment Release

**Version**: v3.1.1 Massive Re-Release
**Tagline**: "Not just functional. Beautiful."
**Goal**: Ship production-ready Ticketr with clean foundation and enchanting UX
**Status**: Planning
**Start Date**: TBD
**Target Completion**: TBD

---

## Executive Summary

Phase 6 represents the final transformation of Ticketr into a clean, production-ready product. After extensive development through Phases 1-5, we've accumulated migration code, feature flags, and transitional documentation that served their purpose but now obscure the clean architecture we've built.

**This is the massive re-release.** Version 3.1.1 ships as a mature, stable product with no feature gates, no migration commands, and no conditional behaviors. Users install Ticketr and it works—period.

### Why Phase 6?

1. **Remove Technical Debt**: Migration code from v2.x compatibility is no longer needed
2. **Improve User Experience**: Async operations, better menus, real-time progress
3. **Clean Documentation**: Single source of truth for requirements and features
4. **Production Ready**: No beta/rc labels, no feature flags, no caveats

### Success Criteria

- Zero feature flag code in codebase
- Zero migration commands or compatibility layers
- Single consolidated REQUIREMENTS.md
- Async TUI operations with cancellation
- Professional progress indicators and menus
- Comprehensive, clean documentation
- All tests passing (target: 95%+ coverage)
- Steward approval for production release

---

## Phase Objectives

### 1. Consolidate Requirements (FINAL)
Create a single, authoritative REQUIREMENTS.md that replaces all previous requirement documents. Remove obsolete migration requirements, add TUI UX requirements, and establish this as the permanent source of truth.

### 2. Remove All Migration/Feature Flag Code
Delete v3_migrate.go, feature flag logic, conditional behaviors, and v2.x compatibility code. Ship clean: Ticketr works as-is, no flags needed.

### 3. Improve TUI UX
Implement async job queue for non-blocking operations, context-aware menus, real-time progress indicators, and professional user experience.

### 4. Clean Release Preparation
Finalize documentation, complete Steward review, and execute production release of v3.1.1.

### 5. TUI Visual & Experiential Polish (The Director's Cut) ✨
Transform the TUI from functional to enchanting through subtle motion, light and shadow, atmospheric effects, and small charms of quality. Make the interface feel alive.

---

## The Four Principles of TUI Excellence

### 1. Subtle Motion is Life
Active spinners, focus pulse, modal fade-in create a living interface.

### 2. Light, Shadow, and Focus
Border styles, drop shadows, title gradients create depth and hierarchy.

### 3. Atmosphere and Ambient Effects
Themeable background effects (hyperspace, snow) add character (default OFF).

### 4. Small Charms of Quality
Success sparkles, animated toggles, polished progress bars show craftsmanship.

These principles guide all TUI work in Phase 6.

---

## Week 1: Foundation Cleanup (Days 1-5)

### Day 1: Requirements Consolidation

**Status**: ⬜ Not Started
**Assigned**: Builder + Verifier + Scribe
**Estimated**: 4 hours | **Actual**: _____

#### Builder Tasks
- [ ] Audit existing requirement documents
  - [ ] Review REQUIREMENTS-v2.md
  - [ ] Review ROADMAP.md active items
  - [ ] Check docs/ for scattered requirements
  - [ ] Identify obsolete migration requirements
- [ ] Create consolidated REQUIREMENTS.md structure
  - [ ] Functional requirements (core features)
  - [ ] Non-functional requirements (performance, security)
  - [ ] TUI/UX requirements (async, menus, progress)
  - [ ] Integration requirements (Jira, filesystem)
  - [ ] Requirement IDs (PROD-xxx format)
- [ ] Write REQUIREMENTS.md content
  - [ ] Import current active requirements
  - [ ] Add TUI UX requirements for Phase 6
  - [ ] Remove migration/compatibility requirements
  - [ ] Ensure traceability to roadmap milestones
  - [ ] Add acceptance criteria for each requirement
- [ ] Cross-reference with codebase
  - [ ] Verify all implemented features documented
  - [ ] Flag undocumented features for review
  - [ ] Remove requirements for deleted features

#### Verifier Tasks
- [ ] Validate all requirements are testable
  - [ ] Each requirement has measurable acceptance criteria
  - [ ] Can map requirement to test cases
  - [ ] No ambiguous language ("should be fast" → "respond in <2s")
- [ ] Check for contradictions
  - [ ] No conflicting requirements
  - [ ] Consistent terminology throughout
  - [ ] Priority conflicts resolved
- [ ] Ensure traceability
  - [ ] Requirements map to roadmap milestones
  - [ ] Can trace requirement → feature → test
  - [ ] Orphaned requirements identified
- [ ] Review completeness
  - [ ] All Phase 1-5 features captured
  - [ ] Phase 6 TUI improvements defined
  - [ ] Security requirements present
  - [ ] Performance requirements quantified

#### Scribe Tasks
- [ ] Update README.md with final feature list
  - [ ] Replace feature section with consolidated list
  - [ ] Remove "beta" and "coming soon" language
  - [ ] Add link to REQUIREMENTS.md
- [ ] Remove migration documentation
  - [ ] Delete v2-to-v3 migration guides
  - [ ] Remove feature flag documentation
  - [ ] Clean up installation docs (no --enable-beta)
- [ ] Update ARCHITECTURE.md
  - [ ] Reflect final architecture (no migration layers)
  - [ ] Document TUI async job queue design
  - [ ] Update component diagrams if needed
- [ ] Archive obsolete docs
  - [ ] Move REQUIREMENTS-v2.md to docs/archive/
  - [ ] Keep for historical reference only
  - [ ] Add README in archive explaining context

#### Acceptance Criteria
- [ ] Single REQUIREMENTS.md file exists at project root
- [ ] No obsolete requirements remain (migration, compatibility)
- [ ] All current features (Phase 1-5) documented with IDs
- [ ] Phase 6 TUI requirements clearly defined
- [ ] Roadmap aligns with requirements (no orphans)
- [ ] README.md reflects final feature set
- [ ] No migration documentation in user-facing docs
- [ ] Verifier sign-off on requirement quality

**Deliverables**:
- `/home/karol/dev/private/ticktr/REQUIREMENTS.md` (new, final)
- `/home/karol/dev/private/ticktr/docs/archive/REQUIREMENTS-v2.md` (archived)
- Updated `/home/karol/dev/private/ticktr/README.md`
- Updated `/home/karol/dev/private/ticktr/docs/ARCHITECTURE.md`

---

### Day 2-3: Remove Migration Code

**Status**: ⬜ Not Started
**Assigned**: Builder + Verifier + Scribe
**Estimated**: 8 hours | **Actual**: _____

#### Builder Tasks: Code Deletion
- [ ] Delete migration command file
  - [ ] Remove `cmd/ticketr/v3_migrate.go` entirely
  - [ ] Remove from build system/imports
- [ ] Remove `v3 enable` command logic
  - [ ] Audit `cmd/ticketr/main.go` for v3 enable command
  - [ ] Remove command registration
  - [ ] Remove subcommand handlers (enable beta/rc/alpha/stable)
  - [ ] Delete associated helper functions
- [ ] Remove feature flag checks in TUI
  - [ ] Search `cmd/ticketr/tui_command.go` for feature flag conditionals
  - [ ] Remove `if featureEnabled("v3")` blocks
  - [ ] Remove fallback behaviors
  - [ ] Simplify control flow
- [ ] Remove feature flag checks in services
  - [ ] Search `internal/core/services/` for feature flags
  - [ ] Remove conditional service initialization
  - [ ] Remove v2/v3 behavior switches
- [ ] Clean up configuration
  - [ ] Check `internal/config/` for feature flag storage
  - [ ] Remove flag fields from config structs
  - [ ] Remove flag parsing logic
  - [ ] Update config validation
- [ ] Remove rollback/migration strategies
  - [ ] Delete state migration logic (if any)
  - [ ] Remove schema version checks
  - [ ] Remove backward compatibility transforms
- [ ] Remove v2.x compatibility code
  - [ ] Search for "v2" comments/tags
  - [ ] Remove legacy API adapters
  - [ ] Remove deprecated function wrappers
  - [ ] Clean up compatibility shims

#### Builder Tasks: Code Simplification
- [ ] Update `cmd/ticketr/main.go`
  - [ ] Remove v3 command from command tree
  - [ ] Simplify initialization (no feature checks)
  - [ ] Update version string to 3.1.1
- [ ] Simplify service initialization
  - [ ] Remove conditional service creation
  - [ ] Direct initialization only
- [ ] Update tests
  - [ ] Remove feature flag test cases
  - [ ] Remove migration test cases
  - [ ] Update test fixtures if needed
  - [ ] Fix broken imports

#### Verifier Tasks
- [ ] Build verification
  - [ ] Run `go build ./...` - succeeds cleanly
  - [ ] No compilation errors
  - [ ] No unused imports
- [ ] Test suite verification
  - [ ] Run `go test ./...` - all pass
  - [ ] Check test coverage (target: maintain current level)
  - [ ] Verify no skipped tests due to missing flags
- [ ] Dead code detection
  - [ ] Run `deadcode ./...` or equivalent
  - [ ] Verify no orphaned functions
  - [ ] Check for unused variables/constants
- [ ] Code search validation
  - [ ] `grep -r "enable beta" .` returns empty (exclude docs/archive)
  - [ ] `grep -r "enable rc" .` returns empty
  - [ ] `grep -r "enable alpha" .` returns empty
  - [ ] `grep -r "enable stable" .` returns empty
  - [ ] `grep -r "v3_migrate" .` returns empty
  - [ ] `grep -r "feature.*flag" .` returns only historical references
  - [ ] `grep -r "v2.*compat" .` returns empty
- [ ] Integration test
  - [ ] Build binary: `go build -o ticketr ./cmd/ticketr`
  - [ ] Run `./ticketr --help` - no v3 command
  - [ ] Run `./ticketr version` - shows 3.1.1
  - [ ] Run `./ticketr tui` - works without flags

#### Scribe Tasks
- [ ] Update README.md
  - [ ] Remove "Migration from v2" section
  - [ ] Remove "Enabling v3 Features" section
  - [ ] Update installation to simple: `go install` + `ticketr tui`
  - [ ] Remove feature flag references
- [ ] Update CHANGELOG.md
  - [ ] Add v3.1.1 entry
  - [ ] Note: "Clean release - removed migration code"
  - [ ] List breaking changes (if any)
  - [ ] Document that v3 is now default behavior
- [ ] Remove migration guides
  - [ ] Delete `docs/v2-to-v3-migration.md` (if exists)
  - [ ] Delete `docs/feature-flags.md` (if exists)
  - [ ] Archive if historical value, otherwise delete
- [ ] Update command documentation
  - [ ] Review `docs/COMMANDS.md` for v3 enable references
  - [ ] Update command tree diagram
  - [ ] Remove deprecated command documentation

#### Acceptance Criteria
- [ ] Zero feature flag code in codebase (verified by grep)
- [ ] Zero migration commands (no v3_migrate.go)
- [ ] `go build ./...` succeeds with no warnings
- [ ] `go test ./...` passes (100% of remaining tests)
- [ ] Binary size reasonable (no bloat from dead code)
- [ ] `ticketr --help` clean (no migration commands)
- [ ] Documentation clean (no migration instructions)
- [ ] Verifier sign-off on code cleanliness

**Deliverables**:
- Deleted: `cmd/ticketr/v3_migrate.go`
- Modified: `cmd/ticketr/main.go` (simplified)
- Modified: `cmd/ticketr/tui_command.go` (no flags)
- Modified: Configuration files (flags removed)
- Updated: README.md, CHANGELOG.md
- Test results: Full suite passing

---

### Day 4-5: Agent Definition Review

**Status**: ⬜ Not Started
**Assigned**: Builder + Verifier + Scribe
**Estimated**: 4 hours | **Actual**: _____

#### Builder Tasks
- [ ] Review Builder agent definition
  - [ ] Read `.agents/builder.agent.md`
  - [ ] Check clarity of responsibilities
  - [ ] Verify alignment with Director methodology
  - [ ] Check workflow steps are actionable
  - [ ] Validate guardrails are specific
  - [ ] Confirm success checklist format
  - [ ] Update if improvements needed
- [ ] Review Verifier agent definition
  - [ ] Read `.agents/verifier.agent.md`
  - [ ] Check test strategy guidelines
  - [ ] Verify coverage requirements clear
  - [ ] Check integration with Builder
  - [ ] Update if improvements needed
- [ ] Review Scribe agent definition
  - [ ] Read `.agents/scribe.agent.md`
  - [ ] Check documentation standards
  - [ ] Verify deliverables clear
  - [ ] Check style guidelines
  - [ ] Update if improvements needed
- [ ] Review Steward agent definition
  - [ ] Read `.agents/steward.agent.md`
  - [ ] Check architecture governance role
  - [ ] Verify approval criteria
  - [ ] Check escalation process
  - [ ] Update if improvements needed
- [ ] Review Director agent definition
  - [ ] Read `.agents/director.agent.md`
  - [ ] Check orchestration responsibilities
  - [ ] Verify phase planning methodology
  - [ ] Check coordination process
  - [ ] Update if improvements needed
- [ ] Review TUIUX agent definition
  - [ ] Read `.agents/tuiux.agent.md`
  - [ ] Check TUI/UX specialist responsibilities
  - [ ] Verify Four Principles of TUI Excellence documented
  - [ ] Check performance budgets and accessibility requirements
  - [ ] Verify deliverables pattern (design, implementation, tests, docs, demo)
  - [ ] Update if improvements needed
- [ ] Check internal Claude agent definitions
  - [ ] Review any `.claudeagent` files
  - [ ] Check for consistency with .agents/
  - [ ] Align terminology and responsibilities

#### Builder Tasks: Updates
- [ ] Standardize format across all agent definitions
  - [ ] Consistent section structure
  - [ ] Unified voice and tone
  - [ ] Standard deliverable format
- [ ] Add cross-references
  - [ ] Link agents that interact
  - [ ] Reference Director methodology
  - [ ] Point to example workflows
- [ ] Clarify handoff points
  - [ ] Builder → Verifier criteria
  - [ ] TUIUX → Verifier criteria (visual polish validation)
  - [ ] Verifier → Scribe triggers
  - [ ] Escalation to Steward process
  - [ ] Director checkpoints

#### Verifier Tasks
- [ ] Validate agent definition consistency
  - [ ] No conflicting responsibilities
  - [ ] No gaps in workflow coverage
  - [ ] Clear ownership boundaries
- [ ] Check Director methodology alignment
  - [ ] Agent roles match methodology
  - [ ] Phase structure consistent
  - [ ] Deliverable formats aligned
- [ ] Verify no contradictions
  - [ ] Builder vs Verifier scope clear
  - [ ] Scribe not duplicating Builder docs
  - [ ] Steward authority well-defined
- [ ] Review completeness
  - [ ] All phases covered
  - [ ] All task types addressed
  - [ ] Emergency procedures defined

#### Scribe Tasks
- [ ] Document agent roles in CONTRIBUTING.md
  - [ ] Add "Agent Roles" section
  - [ ] Summarize each agent's purpose
  - [ ] Explain when each agent engages
  - [ ] Link to detailed .agents/*.md files
- [ ] Update Director Handbook (if needed)
  - [ ] Reflect any methodology refinements
  - [ ] Add Phase 6 as case study
  - [ ] Update examples with Ticketr context
- [ ] Create agent workflow diagram
  - [ ] Visual representation of agent interactions
  - [ ] Show handoff points
  - [ ] Include escalation paths
  - [ ] Add to docs/

#### Acceptance Criteria
- [ ] All six agent definitions reviewed (.agents/*.md)
- [ ] Definitions aligned with Director methodology
- [ ] Clear, non-overlapping responsibilities
- [ ] Consistent format and structure
- [ ] Cross-references added where helpful
- [ ] CONTRIBUTING.md documents 6-agent system
- [ ] Verifier confirms no contradictions
- [ ] Scribe confirms documentation complete

**Deliverables**:
- Updated: `.agents/builder.agent.md`
- Updated: `.agents/verifier.agent.md`
- Updated: `.agents/scribe.agent.md`
- Updated: `.agents/steward.agent.md`
- Updated: `.agents/director.agent.md`
- Updated: `.agents/tuiux.agent.md`
- Updated: `CONTRIBUTING.md` (agent roles section)
- New: `docs/agent-workflow-diagram.md` (or .png)

---

## Week 2: TUI UX Improvements (Days 6-12)

### Day 6-7: Async Job Queue Architecture

**Status**: ⬜ Not Started
**Assigned**: Builder + Verifier
**Estimated**: 12 hours | **Actual**: _____

#### Builder Tasks: Design
- [ ] Design job queue system architecture
  - [ ] Define Job interface (Execute, Cancel, Progress)
  - [ ] Design JobQueue manager (submit, cancel, status)
  - [ ] Plan progress reporting channel design
  - [ ] Sketch context cancellation flow
  - [ ] Document concurrency patterns
- [ ] Create design document
  - [ ] Add to `docs/tui-async-architecture.md`
  - [ ] Include sequence diagrams
  - [ ] Document error handling
  - [ ] Specify thread safety approach

#### Builder Tasks: Implementation
- [ ] Implement Job interface
  - [ ] Create `internal/tui/jobs/job.go`
  - [ ] Define Job, JobStatus, JobProgress types
  - [ ] Add JobID generation (UUID or sequential)
- [ ] Implement JobQueue
  - [ ] Create `internal/tui/jobs/queue.go`
  - [ ] Add goroutine worker pool
  - [ ] Implement Submit(job Job) JobID
  - [ ] Implement Cancel(jobID JobID) error
  - [ ] Implement Status(jobID JobID) JobStatus
  - [ ] Add progress channel: Progress() <-chan JobProgress
  - [ ] Thread-safe with mutex/channels
- [ ] Implement PullJob
  - [ ] Create `internal/tui/jobs/pull_job.go`
  - [ ] Wrap PullService.Pull() in Job interface
  - [ ] Add context cancellation support
  - [ ] Emit progress events (ticket count, percentage)
  - [ ] Handle errors gracefully
- [ ] Wire PullService to run async
  - [ ] Modify `cmd/ticketr/tui_command.go` or equivalent
  - [ ] Submit pull operations to JobQueue
  - [ ] Update UI to show job status
  - [ ] Non-blocking: user can navigate while pulling
- [ ] Add ESC/Ctrl+C cancellation
  - [ ] Capture ESC key in TUI event loop
  - [ ] Capture Ctrl+C signal (SIGINT)
  - [ ] Call JobQueue.Cancel(currentJobID)
  - [ ] Update UI to show "Cancelling..." state
  - [ ] Ensure cleanup on cancel (partial results, connections)

#### Builder Tasks: Integration
- [ ] Integrate with existing TUI
  - [ ] Add job status widget to TUI layout
  - [ ] Subscribe to JobQueue.Progress() channel
  - [ ] Update status bar with job progress
  - [ ] Show spinner or progress bar
- [ ] Update error handling
  - [ ] Surface job errors to user
  - [ ] Add retry mechanism (optional)
  - [ ] Log errors for debugging

#### Verifier Tasks
- [ ] Test async pull operations
  - [ ] Start pull, verify TUI remains responsive
  - [ ] Navigate to different views during pull
  - [ ] Verify ticket list updates after pull completes
- [ ] Verify cancellation works
  - [ ] Press ESC during pull, verify cancel
  - [ ] Press Ctrl+C during pull, verify cancel
  - [ ] Check partial results handled correctly
  - [ ] Ensure no orphaned goroutines (use pprof)
- [ ] Check for race conditions
  - [ ] Run tests with `go test -race ./internal/tui/jobs/...`
  - [ ] Check JobQueue thread safety
  - [ ] Verify progress channel doesn't deadlock
- [ ] Validate progress reporting
  - [ ] Progress events received in order
  - [ ] Percentages accurate (0-100)
  - [ ] Ticket counts match actual pull results
- [ ] Performance test
  - [ ] Pull 500+ tickets, measure responsiveness
  - [ ] Check memory usage (no leaks)
  - [ ] Verify graceful degradation under load
- [ ] Error scenario testing
  - [ ] Network failure during pull
  - [ ] Jira API error during pull
  - [ ] Cancel during network call
  - [ ] Multiple rapid cancellations

#### Acceptance Criteria
- [ ] JobQueue implemented with clean interface
- [ ] PullService operations non-blocking
- [ ] User can navigate TUI during pull
- [ ] ESC cancels active job
- [ ] Ctrl+C cancels active job and exits cleanly
- [ ] Progress updates in real-time (visible in UI)
- [ ] No race conditions (`go test -race` passes)
- [ ] No goroutine leaks (verified with pprof)
- [ ] Partial results handled on cancel
- [ ] Tests pass for all async scenarios
- [ ] Verifier sign-off on async implementation

**Deliverables**:
- New: `internal/tui/jobs/job.go`
- New: `internal/tui/jobs/queue.go`
- New: `internal/tui/jobs/pull_job.go`
- New: `internal/tui/jobs/queue_test.go`
- Modified: `cmd/ticketr/tui_command.go` (async integration)
- New: `docs/tui-async-architecture.md`
- Test results: Race detector clean, performance benchmarks

---

### Day 8-9: TUI Menu Structure

**Status**: ⬜ Not Started
**Assigned**: Builder + Verifier + Scribe
**Estimated**: 8 hours | **Actual**: _____

#### Builder Tasks: Bottom Action Bar
- [ ] Design bottom action bar widget
  - [ ] Sketch layout (full-width, 1-2 rows)
  - [ ] Decide on content: [F1 Help] [F2 Sync] [Ctrl+C Exit] etc.
  - [ ] Choose styling (colors, borders)
- [ ] Implement action bar widget
  - [ ] Create `internal/tui/widgets/actionbar.go`
  - [ ] Render method with bubbletea Model interface
  - [ ] Update method for context changes
  - [ ] View method returns formatted string
- [ ] Add context-aware keybindings
  - [ ] Show different keys based on active view
  - [ ] Example: List view shows [Enter Open] [Space Select]
  - [ ] Example: Detail view shows [Esc Back] [E Edit]
  - [ ] Update action bar dynamically
- [ ] Integrate with main TUI
  - [ ] Add action bar to TUI layout (bottom row)
  - [ ] Wire context changes to update action bar
  - [ ] Ensure consistent across all views

#### Builder Tasks: Command Palette Enhancement
- [ ] Enhance existing command palette (if exists)
  - [ ] Add descriptions to each command
  - [ ] Example: "sync - Pull latest tickets from Jira"
  - [ ] Group commands by category (View, Edit, Sync, etc.)
- [ ] Or implement new command palette
  - [ ] Create `internal/tui/widgets/palette.go`
  - [ ] Trigger with Ctrl+P or F1
  - [ ] Fuzzy search commands
  - [ ] Show descriptions and keybindings
- [ ] Add command metadata
  - [ ] Command name, description, keybinding, category
  - [ ] Centralize in `internal/tui/commands/registry.go`

#### Builder Tasks: F-Key Shortcuts
- [ ] Evaluate F-key shortcuts
  - [ ] F1: Help/Command Palette
  - [ ] F2: Sync/Pull
  - [ ] F5: Refresh view
  - [ ] F10: Exit (common convention)
  - [ ] Balance discoverability vs. complexity
- [ ] Implement chosen F-key shortcuts
  - [ ] Add to TUI event handler
  - [ ] Document in action bar
  - [ ] Add to help screen

#### Verifier Tasks
- [ ] Test all menu interactions
  - [ ] Every keybinding in action bar works
  - [ ] Command palette opens and executes commands
  - [ ] F-keys trigger correct actions
  - [ ] No keybinding conflicts
- [ ] Verify keybinding consistency
  - [ ] Same key does same thing across views (or clearly different)
  - [ ] Standard conventions followed (Esc back, Enter confirm)
  - [ ] No confusing overlaps
- [ ] Check discoverability
  - [ ] New user can find commands without docs
  - [ ] Help/palette accessible from any view
  - [ ] Action bar visible at all times
- [ ] Usability testing
  - [ ] Navigate TUI using only menu hints
  - [ ] Check for missing actions
  - [ ] Verify clarity of descriptions

#### Scribe Tasks
- [ ] Document new menu structure
  - [ ] Add to `docs/TUI-GUIDE.md` or README
  - [ ] Explain action bar concept
  - [ ] List all available keybindings by view
- [ ] Update keybinding reference
  - [ ] Create `docs/KEYBINDINGS.md` (if not exists)
  - [ ] Table format: Key | Action | Context | Description
  - [ ] Include F-keys and standard keys
- [ ] Add screenshots/descriptions
  - [ ] Screenshot of action bar in different views
  - [ ] Screenshot of command palette
  - [ ] Annotate with explanations
  - [ ] Add to docs/images/

#### Acceptance Criteria
- [ ] Bottom action bar visible in all TUI views
- [ ] Action bar shows context-aware keybindings
- [ ] Command palette functional with descriptions
- [ ] F-key shortcuts implemented (if adopted)
- [ ] All actions discoverable without external docs
- [ ] No keybinding conflicts
- [ ] Help comprehensive and accessible
- [ ] Verifier confirms usability
- [ ] Scribe confirms documentation complete

**Deliverables**:
- New: `internal/tui/widgets/actionbar.go`
- Modified: Main TUI layout (action bar integrated)
- Modified/New: `internal/tui/widgets/palette.go`
- New: `internal/tui/commands/registry.go`
- New: `docs/KEYBINDINGS.md`
- Updated: `docs/TUI-GUIDE.md`
- New: Screenshots in `docs/images/`

---

### Day 10-11: Progress Indicators

**Status**: ⬜ Not Started
**Assigned**: Builder + Verifier + Scribe
**Estimated**: 6 hours | **Actual**: _____

#### Builder Tasks: Progress Callbacks
- [ ] Wire progress callbacks to status bar
  - [ ] JobQueue already emits progress events (from Day 6-7)
  - [ ] Subscribe to progress channel in TUI
  - [ ] Update status bar widget with progress data
- [ ] Design progress data structure
  - [ ] Current count, total count
  - [ ] Percentage (calculated or provided)
  - [ ] Time elapsed, ETA (estimated)
  - [ ] Status message (optional)

#### Builder Tasks: Progress Bar Widget
- [ ] Implement ASCII progress bar
  - [ ] Create `internal/tui/widgets/progressbar.go`
  - [ ] Render method: [=====>    ] 50% (45/120)
  - [ ] Configurable width (adapt to terminal size)
  - [ ] Optional: Use Unicode box-drawing chars for polish
- [ ] Add ticket count display
  - [ ] Format: "45/120 tickets"
  - [ ] Update in real-time as progress events arrive
- [ ] Add time display
  - [ ] Track start time when job begins
  - [ ] Show elapsed: "Elapsed: 12s"
  - [ ] Calculate ETA: "ETA: 15s" (simple linear extrapolation)
  - [ ] Update every second or on progress event

#### Builder Tasks: Integration
- [ ] Integrate progress bar with status bar
  - [ ] Add progress bar to status bar widget
  - [ ] Show when job active, hide when idle
  - [ ] Smooth updates (avoid flicker)
- [ ] Handle edge cases
  - [ ] Total unknown (indeterminate progress: spinner)
  - [ ] Rapid updates (throttle to avoid UI churn)
  - [ ] Completion (briefly show 100%, then hide)

#### Verifier Tasks
- [ ] Test progress display accuracy
  - [ ] Start pull, verify count matches actual tickets fetched
  - [ ] Check percentage calculation correct
  - [ ] Verify progress bar width scales with terminal
- [ ] Verify no UI freezing
  - [ ] Progress updates don't block main thread
  - [ ] TUI responsive during long pulls
  - [ ] No dropped input events
- [ ] Check large dataset performance
  - [ ] Pull 500+ tickets, verify progress smooth
  - [ ] Check memory usage during progress reporting
  - [ ] Ensure ETA calculation reasonable
- [ ] Visual testing
  - [ ] Progress bar renders correctly (no garbled chars)
  - [ ] Time displays formatted well (no "Elapsed: 3661s", use 1h 1m 1s)
  - [ ] Colors/styling consistent with TUI theme

#### Scribe Tasks
- [ ] Document progress indicators
  - [ ] Add to `docs/TUI-GUIDE.md`
  - [ ] Explain what progress bar shows
  - [ ] Mention time estimates
- [ ] Update user guide
  - [ ] Add section: "Monitoring Long Operations"
  - [ ] Screenshot of progress bar in action
  - [ ] Explain cancellation during progress

#### Acceptance Criteria
- [ ] Real-time progress visible during async operations
- [ ] Accurate ticket counts (current/total)
- [ ] Percentage displayed and correct
- [ ] ASCII progress bar renders cleanly
- [ ] Time elapsed shown and accurate
- [ ] ETA shown (even if approximate)
- [ ] No UI freezing or performance degradation
- [ ] Progress bar adapts to terminal width
- [ ] Tests pass for progress calculation
- [ ] Verifier confirms professional UX
- [ ] Scribe confirms documentation updated

**Deliverables**:
- New: `internal/tui/widgets/progressbar.go`
- Modified: Status bar widget (progress integration)
- Modified: TUI main loop (progress subscription)
- Updated: `docs/TUI-GUIDE.md`
- Screenshot: Progress bar in action

---

### Day 12: Integration Testing

**Status**: ⬜ Not Started
**Assigned**: Verifier
**Estimated**: 6 hours | **Actual**: _____

#### Verifier Tasks: Full Regression Suite
- [ ] Run complete test suite
  - [ ] Execute: `go test ./... -v -cover`
  - [ ] Verify all tests pass (100% of active tests)
  - [ ] Check coverage report (target: 95%+)
  - [ ] Document test count and coverage percentage
- [ ] Identify and fix flaky tests
  - [ ] Re-run tests 3x to catch intermittent failures
  - [ ] Fix or document flaky tests
- [ ] Review test output
  - [ ] No skipped tests (unless justified)
  - [ ] No panics or race warnings
  - [ ] Reasonable execution time (full suite <2 min)

#### Verifier Tasks: 500+ Ticket Stress Test
- [ ] Prepare large dataset
  - [ ] Mock Jira API with 500+ tickets
  - [ ] Or use staging Jira project with bulk data
- [ ] Execute pull stress test
  - [ ] Run `ticketr tui`, trigger sync
  - [ ] Monitor progress indicators
  - [ ] Verify all 500+ tickets fetched
  - [ ] Check TUI responsive during pull
- [ ] Measure performance
  - [ ] Total time to pull 500 tickets
  - [ ] Memory usage (baseline vs. peak)
  - [ ] CPU usage during operation
- [ ] Verify results
  - [ ] All tickets appear in TUI
  - [ ] No duplicates or missing tickets
  - [ ] State file updated correctly

#### Verifier Tasks: Async Operation Stress Test
- [ ] Test concurrent operations
  - [ ] Start pull, immediately cancel, restart
  - [ ] Rapid ESC presses during pull
  - [ ] Multiple quick Ctrl+C presses
- [ ] Test long-running operations
  - [ ] Simulate slow Jira API (500ms per ticket)
  - [ ] Verify progress updates smooth
  - [ ] Verify cancel responsive even during slow calls
- [ ] Test error scenarios
  - [ ] Network timeout mid-pull
  - [ ] Jira API error (500, 403) during pull
  - [ ] Verify graceful error handling
  - [ ] Verify UI doesn't crash

#### Verifier Tasks: Memory Leak Check
- [ ] Run with memory profiler
  - [ ] Execute: `go test ./internal/tui/jobs/... -memprofile=mem.prof`
  - [ ] Analyze with `go tool pprof mem.prof`
  - [ ] Check for growing allocations (leaks)
- [ ] Long-running test
  - [ ] Run TUI for 10+ minutes with repeated operations
  - [ ] Monitor memory with `top` or `ps`
  - [ ] Verify memory stable (no continuous growth)
- [ ] Goroutine leak check
  - [ ] Capture goroutine profile: `GODEBUG=gctrace=1`
  - [ ] Verify goroutine count stable after operations complete
  - [ ] Check for orphaned workers

#### Verifier Tasks: Race Detection
- [ ] Run tests with race detector
  - [ ] Execute: `go test -race ./...`
  - [ ] Verify zero race conditions reported
  - [ ] Fix any races detected
- [ ] Run TUI with race detector
  - [ ] Build: `go build -race -o ticketr-race ./cmd/ticketr`
  - [ ] Run: `./ticketr-race tui`
  - [ ] Perform async operations (pull, cancel, navigate)
  - [ ] Check for race warnings in output

#### Acceptance Criteria
- [ ] All tests pass (go test ./...)
- [ ] Test coverage ≥95% (or document gaps)
- [ ] No flaky tests (3x re-run clean)
- [ ] 500+ ticket stress test passes
  - [ ] All tickets fetched
  - [ ] TUI responsive
  - [ ] Performance acceptable (<30s for 500 tickets)
- [ ] Async stress test passes
  - [ ] Cancel/restart stable
  - [ ] Error handling graceful
- [ ] No memory leaks detected
- [ ] No goroutine leaks detected
- [ ] Race detector clean (zero races)
- [ ] Verifier approves for release

**Deliverables**:
- Test report: Coverage percentage, test count, execution time
- Stress test report: Performance metrics (time, memory, CPU)
- Race detector report: Clean output
- Memory profile: No leaks confirmed
- Verifier sign-off document

---

### Day 12.5: TUI Visual & Experiential Polish (Director's Cut) ✨

**Status**: ⬜ Not Started
**Assigned**: TUIUX + Verifier + Scribe
**Goal**: Transform the TUI from functional to enchanting

**TUIUX Agent Tasks**:
- [ ] Create `internal/adapters/tui/effects/` package
- [ ] Implement Background Animator system
  - [ ] Goroutine-based rendering
  - [ ] Themeable parameters (characters, density, speed)
  - [ ] CPU-efficient with auto-pause when UI busy
  - [ ] Default themes: `default` (OFF), `dark` (hyperspace), `arctic` (snow)
- [ ] Create `ShadowBox` primitive
  - [ ] Extends tview.Box with automatic drop shadow
  - [ ] Uses ▒ offset characters (1 row, 2 cols)
  - [ ] Integrate into all modal views
- [ ] Enhance Theme system
  - [ ] Add border style parameters (focused/unfocused)
  - [ ] Add animation character sets (spinners, sparkles)
  - [ ] Add background effect configuration
  - [ ] Update `internal/adapters/tui/theme.go`
- [ ] Implement Animation Helpers
  - [ ] Create `internal/adapters/tui/effects/animator.go`
  - [ ] Success sparkle (500ms particle burst)
  - [ ] Checkbox toggle animation (3-frame: [ ]→[•]→[x])
  - [ ] Progress bar shimmer effect
- [ ] Integrate Motion
  - [ ] Focus pulse in main layout (`app.go`)
  - [ ] Modal fade-in (100ms dithered: ░→▒→█)
  - [ ] Active spinner in status bar (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
- [ ] Implement Border Styles
  - [ ] Focused panels: Double-line (╔═╗)
  - [ ] Unfocused panels: Single-line (┌─┐) or rounded (╭─╮)
- [ ] Add Title Gradients
  - [ ] Horizontal gradient in focused panel titles
  - [ ] Use tview dynamic color tags
- [ ] Create Polished Progress Bar
  - [ ] Block characters: [█████░░░░░]
  - [ ] Shimmer animation across completed portion

**Verifier Tasks**:
- [ ] Performance benchmarks with all effects enabled
  - [ ] Test with 1000+ ticket dataset
  - [ ] Verify CPU usage remains minimal
  - [ ] Ensure UI never lags
  - [ ] Measure memory footprint
- [ ] Visual glitch testing
  - [ ] Test on iTerm2, Alacritty, Windows Terminal
  - [ ] Aggressive window resizing
  - [ ] Color depth variations (256 color, true color)
  - [ ] Find and document rendering artifacts
- [ ] Theme switching validation
  - [ ] All themes render correctly
  - [ ] Background effects enable/disable properly
  - [ ] Border styles apply correctly
  - [ ] No visual regressions
- [ ] A/B comparison testing
  - [ ] Compare "enchanted" vs effects-disabled
  - [ ] Subjective: Enhanced version more pleasant?
  - [ ] Verify no distraction from core functionality

**Scribe Tasks**:
- [ ] Document theme customization
  - [ ] Add "Visual Effects" section to TUI guide
  - [ ] Explain how to enable/disable effects
  - [ ] Document theme parameters
  - [ ] Provide example theme configurations
- [ ] Create "Wow" GIF for marketing
  - [ ] Record high-quality GIF (1920x1080, 60fps)
  - [ ] Showcase: Animations, shadows, background effects
  - [ ] Add to README.md
  - [ ] Include in release notes
  - [ ] Upload to GitHub release assets
- [ ] Update README with aesthetic philosophy
  - [ ] Add "Experience" section
  - [ ] Explain the four principles
  - [ ] Highlight visual polish as differentiator

**Acceptance**:
- [ ] Background animator system functional
- [ ] All modals have drop shadows
- [ ] Theme system supports visual effects
- [ ] Animations smooth and non-intrusive
- [ ] Performance tests pass (no lag with effects)
- [ ] Multi-terminal compatibility verified
- [ ] Marketing GIF created and impressive
- [ ] Documentation complete

**Estimated**: 12 hours | **Actual**: _____

**Architecture Notes**:
- Effects system must be OPTIONAL (theme-controlled)
- Zero performance impact when effects disabled
- CPU-efficient goroutines with rate limiting
- All animations interruptible
- Graceful degradation on limited terminals

---

## Week 3: Release Preparation (Days 13-15)

### Day 13: Documentation Finalization

**Status**: ⬜ Not Started
**Assigned**: Scribe
**Estimated**: 6 hours | **Actual**: _____

#### Scribe Tasks: README.md Update
- [ ] Update version references to v3.1.1
- [ ] Remove beta/rc language
  - [ ] No "experimental" warnings
  - [ ] No "under development" caveats
  - [ ] Present as stable, production-ready
- [ ] Update feature list
  - [ ] Ensure all Phase 1-5 features listed
  - [ ] Add Phase 6 TUI improvements (async, menus, progress)
  - [ ] Remove deprecated features
- [ ] Simplify installation instructions
  - [ ] Single command: `go install github.com/user/ticketr/cmd/ticketr@latest`
  - [ ] No feature flags to enable
  - [ ] Quick start: `ticketr tui`
- [ ] Update screenshots (if needed)
  - [ ] Show new TUI with action bar and progress
  - [ ] Ensure visuals reflect v3.1.1 UI
- [ ] Polish sections
  - [ ] Introduction: Highlight stability and maturity
  - [ ] Usage: Clear, concise examples
  - [ ] Configuration: Current options only
  - [ ] Contributing: Link to CONTRIBUTING.md

#### Scribe Tasks: CHANGELOG.md Entry
- [ ] Add v3.1.1 section at top
  - [ ] Release date: TBD (will update on Day 15)
  - [ ] Version: 3.1.1
- [ ] Summarize changes
  - [ ] **Changed**: Massive clean release, removed migration code
  - [ ] **Added**: Async job queue for non-blocking operations
  - [ ] **Added**: Context-aware TUI menus and action bar
  - [ ] **Added**: Real-time progress indicators with ETA
  - [ ] **Removed**: Feature flag system (v3 enable commands)
  - [ ] **Removed**: Migration commands (v3_migrate.go)
  - [ ] **Fixed**: [Any bugs fixed during Phase 6]
- [ ] Note breaking changes (if any)
  - [ ] Example: "v3 enable commands removed; v3 is now default"
  - [ ] Migration path: "No action needed; v3 behavior is default"
- [ ] Link to documentation
  - [ ] Reference REQUIREMENTS.md for full feature list
  - [ ] Link to TUI-GUIDE.md for new TUI features

#### Scribe Tasks: Release Notes
- [ ] Write release notes document
  - [ ] Create `docs/RELEASE-NOTES-3.1.1.md`
  - [ ] Target audience: Users and stakeholders
- [ ] Structure
  - [ ] **Headline**: "Ticketr v3.1.1: Clean, Production-Ready Release"
  - [ ] **Overview**: Describe the massive re-release
  - [ ] **Key Features**: Async operations, improved TUI, clean codebase
  - [ ] **What's New**: Summarize Phase 6 improvements
  - [ ] **Breaking Changes**: Feature flag removal (transparent to new users)
  - [ ] **Upgrade Instructions**: "Simply update to v3.1.1, no migration needed"
  - [ ] **Known Issues**: Document any remaining issues (or "None")
  - [ ] **Roadmap**: Link to ROADMAP.md for future plans
- [ ] Polish for distribution
  - [ ] Professional tone
  - [ ] Clear value proposition
  - [ ] Thank contributors

#### Scribe Tasks: General Documentation Review
- [ ] Review all docs for v3.1.1 accuracy
  - [ ] Check `docs/*.md` files
  - [ ] Ensure no migration references (except in archive)
  - [ ] Update version numbers where mentioned
- [ ] Review ARCHITECTURE.md
  - [ ] Reflect async job queue architecture
  - [ ] Remove migration layer diagrams
  - [ ] Update component diagram if needed
- [ ] Review CONTRIBUTING.md
  - [ ] Ensure agent roles section current (from Day 4-5)
  - [ ] Update development workflow if changed
  - [ ] Check code style guidelines current
- [ ] Review TUI-GUIDE.md (or create if missing)
  - [ ] Document new action bar
  - [ ] Document progress indicators
  - [ ] Document async operations and cancellation
  - [ ] Include keybinding reference
- [ ] Final proofread
  - [ ] Spelling and grammar check all docs
  - [ ] Consistent terminology (e.g., "Jira" not "JIRA")
  - [ ] Consistent formatting (headings, lists, code blocks)
  - [ ] Check all links work (internal and external)

#### Acceptance Criteria
- [ ] README.md updated for v3.1.1 (no beta language)
- [ ] CHANGELOG.md has v3.1.1 entry with complete change list
- [ ] Release notes written (`docs/RELEASE-NOTES-3.1.1.md`)
- [ ] All documentation reviewed and updated
- [ ] No references to migration/feature flags in user-facing docs
- [ ] Screenshots current (if TUI visuals changed significantly)
- [ ] All links functional (checked with link validator)
- [ ] Spelling and grammar clean (checked with spellchecker)
- [ ] Scribe confirms documentation complete and polished

**Deliverables**:
- Updated: `/home/karol/dev/private/ticktr/README.md`
- Updated: `/home/karol/dev/private/ticktr/CHANGELOG.md`
- New: `/home/karol/dev/private/ticktr/docs/RELEASE-NOTES-3.1.1.md`
- Updated: All docs in `/home/karol/dev/private/ticktr/docs/`
- Updated screenshots (if applicable)

---

### Day 14: Steward Final Approval

**Status**: ⬜ Not Started
**Assigned**: Steward
**Estimated**: 4 hours | **Actual**: _____

#### Steward Tasks: Architecture Compliance Check
- [ ] Review architecture against ARCHITECTURE.md
  - [ ] Async job queue follows clean architecture patterns
  - [ ] TUI widgets properly separated (presentation layer)
  - [ ] No business logic in TUI (delegated to services)
- [ ] Check dependency graph
  - [ ] No circular dependencies
  - [ ] Clear layering (adapters → core → TUI)
  - [ ] Third-party dependencies justified and minimal
- [ ] Verify design patterns
  - [ ] Consistent use of ports/adapters
  - [ ] Proper use of interfaces
  - [ ] No tight coupling introduced

#### Steward Tasks: Security Review
- [ ] Credential handling
  - [ ] Jira API tokens stored securely (keyring or encrypted)
  - [ ] No credentials in logs or error messages
  - [ ] Config files have appropriate permissions
- [ ] Input validation
  - [ ] User input sanitized (TUI, CLI)
  - [ ] Jira API responses validated
  - [ ] No injection vulnerabilities (SQL, command, etc.)
- [ ] Dependency security
  - [ ] Run `go list -m all | nancy sleuth` (or similar)
  - [ ] Check for known vulnerabilities in dependencies
  - [ ] Update vulnerable dependencies if found
- [ ] Error handling
  - [ ] Sensitive data not exposed in errors
  - [ ] Stack traces controlled (not shown to end users)

#### Steward Tasks: Requirements Validation
- [ ] Review REQUIREMENTS.md completeness
  - [ ] All implemented features documented
  - [ ] All requirements have acceptance criteria
  - [ ] Phase 6 requirements fully met
- [ ] Trace requirements to code
  - [ ] Sample 10 requirements, verify implementation
  - [ ] Check requirement IDs referenced in commits/code
- [ ] Trace requirements to tests
  - [ ] Sample 10 requirements, verify test coverage
  - [ ] Ensure critical requirements have integration tests

#### Steward Tasks: Release Readiness Assessment
- [ ] Code quality
  - [ ] Code follows Go idioms and style
  - [ ] Comments adequate (not excessive, not missing)
  - [ ] No technical debt introduced
- [ ] Test quality
  - [ ] Coverage ≥95% (per Verifier report from Day 12)
  - [ ] Integration tests cover critical paths
  - [ ] No ignored or skipped tests without reason
- [ ] Documentation quality
  - [ ] User-facing docs complete (per Scribe report from Day 13)
  - [ ] Internal docs (architecture, design decisions) current
  - [ ] No outdated information
- [ ] Performance
  - [ ] Stress test results acceptable (per Verifier report)
  - [ ] No known performance regressions
  - [ ] Reasonable resource usage (memory, CPU)
- [ ] Stability
  - [ ] No critical bugs open
  - [ ] No known crashes or data corruption issues
  - [ ] Graceful degradation on errors

#### Steward Tasks: Go/No-Go Decision
- [ ] Review all Phase 6 deliverables
  - [ ] Requirements consolidated (Day 1)
  - [ ] Migration code removed (Day 2-3)
  - [ ] Agent definitions updated (Day 4-5)
  - [ ] Async architecture implemented (Day 6-7)
  - [ ] TUI menus improved (Day 8-9)
  - [ ] Progress indicators added (Day 10-11)
  - [ ] Integration tests passed (Day 12)
  - [ ] Documentation finalized (Day 13)
- [ ] Check outstanding issues
  - [ ] Review issue tracker for blockers
  - [ ] Verify all P0/P1 issues resolved
  - [ ] Document any known issues for release notes
- [ ] Final approval decision
  - [ ] **GO**: Approve release to production
  - [ ] **NO-GO**: Document blockers and required remediation
  - [ ] Sign-off document with rationale

#### Acceptance Criteria
- [ ] Architecture compliance verified
- [ ] Security review complete (no critical vulnerabilities)
- [ ] Requirements fully validated (implemented and tested)
- [ ] Release readiness criteria met:
  - [ ] Code quality acceptable
  - [ ] Test quality acceptable (≥95% coverage)
  - [ ] Documentation complete
  - [ ] Performance acceptable
  - [ ] Stability acceptable
- [ ] Go/No-Go decision documented
- [ ] Steward formal approval granted (if GO)

**Deliverables**:
- Steward Review Report:
  - Architecture compliance: PASS/FAIL
  - Security review: PASS/FAIL (with findings if any)
  - Requirements validation: PASS/FAIL
  - Release readiness: GO/NO-GO
  - Sign-off signature (name, date)
- Remediation plan (if NO-GO)

---

### Day 15: Release

**Status**: ⬜ Not Started
**Assigned**: Director
**Estimated**: 4 hours | **Actual**: _____

#### Director Tasks: Pre-Release Verification
- [ ] Verify Steward approval received
  - [ ] Check for GO decision from Day 14
  - [ ] Review any conditional approvals
- [ ] Final smoke test
  - [ ] Build release binary: `go build -o ticketr ./cmd/ticketr`
  - [ ] Run `./ticketr version` - verify 3.1.1
  - [ ] Run `./ticketr tui` - verify TUI launches
  - [ ] Execute basic workflow (pull, view, navigate)
  - [ ] Verify no crashes or obvious issues

#### Director Tasks: Tag Release
- [ ] Update version in code (if not already)
  - [ ] Check `cmd/ticketr/version.go` or equivalent
  - [ ] Set to `3.1.1`
  - [ ] Commit: `git commit -m "chore: Bump version to 3.1.1 for release"`
- [ ] Create Git tag
  - [ ] Tag: `git tag -a v3.1.1 -m "Release v3.1.1: Clean production release"`
  - [ ] Verify: `git tag -l -n1 v3.1.1`
- [ ] Push tag to remote
  - [ ] Push: `git push origin v3.1.1`
  - [ ] Verify on GitHub/GitLab: Tag visible

#### Director Tasks: Build Binaries
- [ ] Build for target platforms
  - [ ] Linux amd64: `GOOS=linux GOARCH=amd64 go build -o ticketr-linux-amd64 ./cmd/ticketr`
  - [ ] Linux arm64: `GOOS=linux GOARCH=arm64 go build -o ticketr-linux-arm64 ./cmd/ticketr`
  - [ ] macOS amd64: `GOOS=darwin GOARCH=amd64 go build -o ticketr-darwin-amd64 ./cmd/ticketr`
  - [ ] macOS arm64: `GOOS=darwin GOARCH=arm64 go build -o ticketr-darwin-arm64 ./cmd/ticketr`
  - [ ] Windows amd64: `GOOS=windows GOARCH=amd64 go build -o ticketr-windows-amd64.exe ./cmd/ticketr`
- [ ] Verify binaries
  - [ ] Check each binary runs: `./ticketr-* version`
  - [ ] Verify size reasonable (not bloated)
- [ ] Create checksums
  - [ ] Generate: `sha256sum ticketr-* > checksums.txt`
  - [ ] Or use `shasum -a 256`

#### Director Tasks: Create GitHub Release
- [ ] Draft GitHub release
  - [ ] Navigate to GitHub Releases page
  - [ ] Click "Draft a new release"
  - [ ] Tag version: v3.1.1
  - [ ] Release title: "Ticketr v3.1.1: Clean Production Release"
  - [ ] Description: Copy from `docs/RELEASE-NOTES-3.1.1.md`
- [ ] Upload binaries
  - [ ] Attach all platform binaries
  - [ ] Attach checksums.txt
- [ ] Publish release
  - [ ] Mark as latest release
  - [ ] Do NOT mark as pre-release (this is stable)
  - [ ] Publish

#### Director Tasks: Update Documentation Site (if applicable)
- [ ] Update online docs
  - [ ] If docs hosted separately (ReadTheDocs, GitHub Pages, etc.)
  - [ ] Update version to 3.1.1
  - [ ] Ensure docs reflect new features
- [ ] Update installation instructions
  - [ ] Ensure `go install` command has correct version/path
  - [ ] Update download links for binaries

#### Director Tasks: Announce Release
- [ ] Internal announcement
  - [ ] Notify team/stakeholders via Slack/email
  - [ ] Highlight key improvements (async, TUI UX, clean codebase)
- [ ] External announcement (if applicable)
  - [ ] Post to project blog/website
  - [ ] Tweet or social media post
  - [ ] Update project description on GitHub
  - [ ] Post to relevant communities (Reddit, forums, etc.)
- [ ] Update CHANGELOG
  - [ ] Set release date in CHANGELOG.md
  - [ ] Commit: `git commit -m "docs: Set release date for v3.1.1"`
  - [ ] Push to main

#### Acceptance Criteria
- [ ] Version tagged in Git (v3.1.1)
- [ ] Binaries built for all target platforms
- [ ] Checksums generated
- [ ] GitHub release published with binaries
- [ ] Documentation site updated (if applicable)
- [ ] Announcement made (internal and external)
- [ ] CHANGELOG.md has release date set
- [ ] Director confirms successful release

**Deliverables**:
- Git tag: `v3.1.1`
- Release binaries: `ticketr-{platform}-{arch}`
- Checksums: `checksums.txt`
- GitHub release: https://github.com/{user}/ticketr/releases/tag/v3.1.1
- Announcement posts/emails
- Updated CHANGELOG.md with release date

---

## Phase 6 Final Checklist

### Pre-Release Verification

**Foundation Cleanup (Week 1)**:
- [ ] Requirements finalized (REQUIREMENTS.md exists and complete)
- [ ] All migration code removed (cmd/ticketr/v3_migrate.go deleted)
- [ ] Feature flags removed (grep verification clean)
- [ ] Agent definitions reviewed and updated
- [ ] Documentation clean (no migration references)

**TUI UX Improvements (Week 2)**:
- [ ] Async job queue implemented and tested
- [ ] TUI async operations working (non-blocking pull)
- [ ] ESC/Ctrl+C cancellation functional
- [ ] TUI menu structure improved (action bar, palette)
- [ ] Progress indicators functional (count, percentage, ETA)
- [ ] Integration tests passed (500+ ticket stress test)
- [ ] Race detector clean (no race conditions)
- [ ] Memory leak check passed (no leaks)

**Release Preparation (Week 3)**:
- [ ] All tests passing (go test ./...)
  - [ ] Test count: _____
  - [ ] Coverage: _____%
- [ ] Documentation complete (README, CHANGELOG, release notes)
- [ ] Steward approval received (GO decision)
- [ ] No critical bugs open

### Release Execution

- [ ] Version: 3.1.1 (set in code)
- [ ] Git tag created: v3.1.1
- [ ] Binaries built (Linux, macOS, Windows)
- [ ] Checksums generated
- [ ] GitHub release published
- [ ] Documentation site updated
- [ ] Announcement made

### Success Metrics

**Code Quality**:
- [ ] Zero feature flags remaining in codebase
- [ ] Zero migration code remaining
- [ ] Clean build: `go build ./...` (no warnings)
- [ ] Test coverage: ≥95% (actual: ____%)
- [ ] Race detector clean: `go test -race ./...`

**UX Quality**:
- [ ] Async operations non-blocking
- [ ] Progress indicators accurate and smooth
- [ ] TUI discoverable (menus, help, keybindings)
- [ ] Professional appearance (polished widgets)

**Documentation Quality**:
- [ ] Single REQUIREMENTS.md (consolidated)
- [ ] README.md clean (no beta language)
- [ ] CHANGELOG.md complete (v3.1.1 entry)
- [ ] Release notes comprehensive
- [ ] All docs proofread and current

**Release Quality**:
- [ ] Steward approved
- [ ] Verifier approved (tests passed)
- [ ] Scribe approved (docs complete)
- [ ] Director executed release
- [ ] Binaries available for download
- [ ] Announcement published

---

## Time Estimates

### Week 1: Foundation Cleanup
- Day 1 (Requirements): 4 hours
- Day 2-3 (Migration Code): 8 hours
- Day 4-5 (Agent Definitions): 4 hours
- **Week 1 Total**: 16 hours

### Week 2: TUI UX Improvements
- Day 6-7 (Async Architecture): 12 hours
- Day 8-9 (Menu Structure): 8 hours
- Day 10-11 (Progress Indicators): 6 hours
- Day 12 (Integration Testing): 6 hours
- Day 12.5 (Visual & Experiential Polish): 12 hours
- **Week 2 Total**: 44 hours

### Week 3: Release Preparation
- Day 13 (Documentation): 6 hours
- Day 14 (Steward Approval): 4 hours
- Day 15 (Release): 4 hours
- **Week 3 Total**: 14 hours

### Phase 6 Total
**Total Estimated Hours**: 74 hours (~3.5 calendar weeks at 20-25 hours/week)

---

## Dependencies and Risks

### Dependencies
1. **Steward Approval**: Release (Day 15) depends on Steward GO decision (Day 14)
2. **Migration Code Removal**: TUI work (Week 2) depends on clean codebase (Day 2-3)
3. **Requirements**: All work depends on finalized requirements (Day 1)
4. **Integration Tests**: Release depends on passing tests (Day 12)

### Risks and Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Async implementation complex, takes longer | Week 2 delays | Medium | Allocate buffer time; start simple (queue + cancel) |
| Steward finds critical issue on Day 14 | Release delays | Low | Continuous review throughout phase; early Steward check-ins |
| Migration code deeply intertwined | Week 1 extends | Medium | Thorough code audit Day 2; flag complex areas early |
| Test coverage drops during refactor | Verifier blocks release | Low | TDD approach; write tests before removing code |
| Documentation incomplete | Release delays | Low | Scribe works in parallel with Builder throughout |

---

## Communication Plan

### Daily Standups (Async)
- **Builder**: What I implemented, what's next, blockers
- **Verifier**: What I tested, results, concerns
- **Scribe**: What I documented, gaps identified
- **Steward**: Architecture reviews, guidance
- **Director**: Progress tracking, coordination

### Weekly Milestones
- **End of Week 1**: Foundation clean, ready for TUI work
- **End of Week 2**: TUI UX complete, integration tests passed
- **End of Week 3**: Release published

### Escalation Path
- **Issue**: Builder → Director (coordination)
- **Architecture**: Any agent → Steward (governance)
- **Blocker**: Director → Steward (decision)

---

## Notes

### Philosophy
This is the **MASSIVE RE-RELEASE**. Ticketr v3.1.1 ships clean, production-ready, with excellent UX. No feature flags, no migrations, no conditional behaviors. This is the real product.

### Quality Over Speed
Phase 6 prioritizes quality and cleanliness. If a task takes longer to do right, that's acceptable. The goal is a stable, maintainable foundation for future development.

### User-Centric
All TUI improvements are driven by user experience. Async operations, clear menus, progress indicators—these make Ticketr pleasant to use. Never compromise UX for implementation convenience.

### Clean Slate
After Phase 6, Ticketr has a clean slate. No legacy code, no transitional logic. Future features build on this solid foundation. Treat this phase as the reset button for technical debt.

---

## The Director's Vision

This is not just a release. This is a statement.

We are crafting an experience that respects the user's time, rewards their attention, and delights with small touches of excellence. Every animation, every shadow, every subtle motion is intentional. We're not building a tool that "works." We're building a tool that feels *alive*.

The TUI is not just a client. It is the primary interface. It must be beautiful.

**Make it so.**

---

## Post-Phase 6 Next Steps

After successful v3.1.1 release:

1. **Monitor**: Watch for user feedback and bug reports
2. **Support**: Address critical issues with patch releases (3.1.2, etc.)
3. **Plan**: Begin Phase 7 planning based on roadmap
4. **Celebrate**: Acknowledge team effort on massive re-release

---

## Appendix: Checklist Summary

For quick reference, the critical checkboxes:

### Week 1
- [ ] REQUIREMENTS.md created and finalized
- [ ] cmd/ticketr/v3_migrate.go deleted
- [ ] All feature flag code removed
- [ ] Agent definitions updated
- [ ] Documentation cleaned

### Week 2
- [ ] Async job queue implemented
- [ ] Progress indicators working
- [ ] TUI menus improved
- [ ] All integration tests passed
- [ ] Race detector clean

### Week 3
- [ ] Documentation finalized
- [ ] Steward approval received
- [ ] v3.1.1 released
- [ ] Announcement published

---

**End of Phase 6 Execution Plan**

This document is the single source of truth for Phase 6 execution. Track progress by checking boxes. Update "Actual" hours for learning and future estimation.

**Success**: When all checkboxes are checked, Ticketr v3.1.1 is shipped, and the team celebrates a clean, production-ready release.
