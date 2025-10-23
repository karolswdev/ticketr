# Bubbletea Migration - Decision Matrix

**Quick Reference Guide for Implementation Decisions**

---

## Component Selection Matrix

| UI Need | Approved Library | Rationale | Complexity |
|---------|-----------------|-----------|------------|
| **Workspace List** | `bubbles/list` | Pagination, filtering built-in | Low |
| **Ticket Tree** | **CUSTOM** | No library exists, need full control | **High** |
| **Text Input** | `bubbles/textinput` | Standard, robust | Low |
| **Multi-line Text** | `bubbles/textarea` | For descriptions | Low |
| **Scrollable Detail** | `bubbles/viewport` | Smooth scrolling | Low |
| **Forms** | `huh/v2` | Best form library | Low |
| **Progress Bar** | `bubbles/progress` | With custom shimmer | Medium |
| **Spinner** | `bubbles/spinner` | Loading indicator | Low |
| **Help Screen** | `bubbles/help` | Auto-generated | Low |
| **Table** | `bubbles/table` | Structured data | Low |
| **Modals** | Custom (`lipgloss.Place`) | Full control | Medium |
| **Layout** | `lipgloss` (native) | Simple, no deps | Low |

---

## Architecture Pattern Decisions

### State Management
✅ **APPROVED:** Elm Architecture (Model-View-Update)
❌ **REJECTED:** Callbacks, imperative state

**Why:** Single source of truth, eliminates race conditions, testable

### Message Passing
✅ **APPROVED:** Custom message types for domain events
❌ **REJECTED:** Direct function calls, channels

**Why:** Pure functions, no side effects, predictable state transitions

### Component Communication
✅ **APPROVED:** Messages via Update() dispatcher
❌ **REJECTED:** Callbacks, shared state, event emitters

**Why:** Decoupled, testable, no hidden dependencies

### Async Operations
✅ **APPROVED:** `tea.Cmd` for all I/O
❌ **REJECTED:** Goroutines, channels

**Why:** Framework-native, prevents race conditions

### Layout System
✅ **APPROVED:** Lipgloss (JoinHorizontal/Vertical)
❌ **REJECTED:** bubblelayout, stickers, bubbleboxer

**Why:** Simple, native, no external dependencies

---

## Implementation Strategy Decisions

### Development Approach
✅ **APPROVED:** Parallel development (keep Tview working)
❌ **REJECTED:** Direct replacement, hybrid approach

**Why:** Lower risk, allows rollback, gradual testing

### Feature Flag
✅ **APPROVED:** `TICKETR_USE_BUBBLETEA=true`
❌ **REJECTED:** Hard cutover, compile-time flag

**Why:** Runtime switching, easy rollback, gradual adoption

### Testing Strategy
✅ **APPROVED:** `teatest` + golden files + unit tests
❌ **REJECTED:** Manual testing only

**Why:** Automated regression detection, fast feedback

### Rollout Timeline
✅ **APPROVED:** 16-18 weeks with beta period
❌ **REJECTED:** 12 weeks (too aggressive), 24+ weeks (too slow)

**Why:** Realistic complexity estimate, includes polish and testing

---

## Feature Inclusion Decisions

### Phase 1 (MVP)
✅ **INCLUDE:**
- Workspace list with sync indicators
- Ticket tree (custom component)
- Ticket detail view (display + edit)
- Search modal with filters
- Command palette
- Bulk operations
- Sync operations (pull, push, full)
- Action bar with context-aware keybindings
- Status bar with workspace info
- Progress bar with shimmer
- Basic animations (fade, slide)
- 3 themes (Default, Dark, Arctic)

❌ **DEFER to Phase 6.5:**
- Background particle effects (cosmic, hyperspace, snow, matrix)
- Lua plugin system
- Tri-panel mode
- Picture-in-picture
- Git integration UI
- Advanced dashboard widgets

**Why:** Particle effects are very high complexity (ANSI compositing), not required for feature parity. Ship MVP first, add later.

---

## Technical Debt Decisions

### What to Migrate
✅ **MIGRATE:**
- All view logic (app.go, views/*.go)
- All widget logic (widgets/*.go)
- Core effects (animator, shimmer, shadowbox, borders)
- Marquee animation
- Theme system
- Keybinding management

❌ **KEEP AS-IS:**
- Domain logic (`internal/domain/`)
- Service layer (`internal/adapters/jira/`, etc.)
- Job queue (`internal/tui/jobs/`)
- Database adapters

⚠️ **DEFER:**
- Background particles (effects/background.go)
- Plugin system (future feature)

**Why:** Domain/services are TUI-agnostic, job queue is shared infrastructure

### Code Organization
✅ **APPROVED:** New package `internal/adapters/bubbletea/`
❌ **REJECTED:** Modify `internal/adapters/tui/` in-place

**Why:** Parallel development, easy comparison, rollback safety

---

## Performance Decisions

### Target Metrics
| Metric | Target | Measurement |
|--------|--------|-------------|
| **Frame Rate** | 60 FPS | Profiling |
| **Tree Render (1000 items)** | <16ms | Benchmarks |
| **Search Response** | <100ms | Integration tests |
| **Key Press Latency** | <50ms | User testing |
| **Startup Time** | <500ms | Benchmarks |
| **Memory Usage** | <50MB | Profiling |

### Optimization Strategy
✅ **APPROVED:**
- Lazy rendering (only visible tree nodes)
- Viewport optimization (pagination)
- String pooling for repeated renders
- Benchmark suite in CI

❌ **REJECTED:**
- Premature optimization
- Custom rendering engine
- Bypassing Bubbletea renderer

**Why:** Bubbletea renderer is fast, trust the framework

---

## Quality Gate Decisions

### Test Coverage Requirements
| Area | Minimum | Target |
|------|---------|--------|
| **Overall** | 70% | >80% |
| **Action System** | 100% | 100% |
| **Predicates** | 100% | 100% |
| **Message Routing** | 100% | 100% |
| **Components** | 70% | >80% |
| **Critical Paths** | 100% | 100% |

### Code Quality Gates
✅ **ENFORCE:**
- Cyclomatic complexity <15
- Function length <100 lines
- File length <500 lines
- No TODO comments before release
- All functions documented

❌ **IGNORE:**
- Line length (trust gofmt)
- Comment length
- Struct field order

### Visual Quality Gates
✅ **REQUIRE:**
- No flicker reported
- No screen tearing
- Responsive resize (<100ms)
- All components themed
- High contrast mode readable

❌ **NICE-TO-HAVE:**
- Pixel-perfect alignment
- Custom fonts
- Terminal-specific optimizations

---

## Risk Mitigation Decisions

### High-Risk: Tree Component
✅ **MITIGATE:**
- Allocate 2 weeks (Weeks 4-5)
- Early prototype in Phase 0
- Performance benchmarks mandatory (<16ms)
- Fallback: bubbles/list with indentation

❌ **ACCEPT RISK:**
- No fallback plan
- "We'll figure it out"

**Why:** Tree is critical path, must have contingency

### Medium-Risk: Background Particles
✅ **MITIGATE:**
- Defer to Phase 6.5 (stretch goal)
- Not required for MVP
- Revisit after v3.2.0 stable

❌ **BLOCK RELEASE:**
- Require particles for v3.2.0
- Allocate 4+ weeks to implement

**Why:** Particles are complex, not essential, high risk of delay

### Medium-Risk: Learning Curve
✅ **MITIGATE:**
- Documentation-first approach
- Pair programming encouraged
- Weekly architecture syncs
- Code examples everywhere

❌ **IGNORE:**
- "Team will figure it out"
- No training plan

**Why:** Elm Architecture is paradigm shift, needs support

---

## Rollout Strategy Decisions

### Beta Rollout
✅ **APPROVED:**
- Week 17: Internal dogfooding (3 days)
- Week 17: v3.2.0-beta.1 (limited external)
- Week 18: v3.2.0-beta.2 (wider testing)
- Week 18: v3.2.0 stable (Bubbletea default)

❌ **REJECTED:**
- No beta period
- Extended 6-month beta

**Why:** Balanced approach, enough time to find bugs, not too long

### Rollback Plan
✅ **TRIGGER ROLLBACK IF:**
- Critical crash (data loss)
- Performance regression >50%
- Security vulnerability
- Major feature broken

❌ **IGNORE ROLLBACK:**
- Minor bugs
- Visual glitches
- Edge cases

**Why:** User trust is paramount, stability over new features

### Deprecation Timeline
✅ **APPROVED:**
- v3.2.0: Bubbletea default, Tview fallback
- v3.3.0 (+6 months): Tview deprecated warning
- v4.0.0 (+12 months): Tview removed

❌ **REJECTED:**
- Immediate Tview removal
- Indefinite dual support

**Why:** Give users time to adapt, but don't maintain two codebases forever

---

## When to Escalate

### Immediate Escalation (Block Sprint)
- Tree component performance <16ms cannot be achieved
- Critical bugs in beta (data loss, crashes)
- Security vulnerabilities discovered
- Performance regression >50% from Tview

### Weekly Escalation (Discuss in Sync)
- Phase timeline slipping >1 week
- Test coverage below target
- New risks identified
- User feedback strongly negative

### No Escalation Needed
- Minor bugs (log as issues)
- Visual polish improvements
- Documentation improvements
- Feature requests for future

---

## Decision Authority

| Decision Type | Authority | Escalation |
|--------------|-----------|------------|
| **Architecture** | Director | Steward (if major change) |
| **Component choice** | Builder | Director (if diverging) |
| **Timeline** | Director | Steward (if >2 week slip) |
| **Quality gates** | Verifier | Director (if conflict) |
| **Documentation** | Scribe | Director (if incomplete) |
| **Release decision** | Steward | Team consensus |

---

## Quick Decision Flowchart

```
Need to make a decision?
│
├─ Does it affect architecture?
│  └─ Check "Architecture Pattern Decisions"
│
├─ Does it affect UI component choice?
│  └─ Check "Component Selection Matrix"
│
├─ Does it affect timeline?
│  └─ Check "Implementation Strategy Decisions"
│
├─ Does it affect what we ship?
│  └─ Check "Feature Inclusion Decisions"
│
├─ Does it affect quality?
│  └─ Check "Quality Gate Decisions"
│
├─ Does it affect risk?
│  └─ Check "Risk Mitigation Decisions"
│
└─ Still unsure?
   └─ Escalate to Director
```

---

## Version History

- **v1.0.0** (2025-10-22): Initial decision matrix
- Future updates tracked in git history

---

**Remember:** When in doubt, refer to the master plan (`TUI_MIGRATION_MASTER_PLAN.md`). This matrix is for quick lookups only.
