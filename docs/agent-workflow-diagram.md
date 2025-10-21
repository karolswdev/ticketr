# Agent Workflow Diagram

**Version:** 1.0
**Last Updated:** Phase 6, Week 1 Day 4-5
**Purpose:** Visual representation of the 6-agent development methodology

---

## Overview

Ticketr development follows a **6-Agent Methodology** where specialized agents handle different aspects of the development lifecycle. This diagram illustrates the workflow, agent interactions, and handoff points.

---

## The 6 Specialized Agents

```
┌─────────────────────────────────────────────────────────────────────┐
│                         6-AGENT METHODOLOGY                         │
│                                                                       │
│  Builder  →  Verifier  →  Scribe  →  (Steward)  →  Director         │
│                                              ↓                        │
│                                          TUIUX                        │
└─────────────────────────────────────────────────────────────────────┘
```

| Agent | Role | Invocation |
|-------|------|------------|
| **Director** | Control Flow Orchestrator | `Task(subagent_type="general-purpose")` |
| **Builder** | Feature Developer | `Task(subagent_type="general-purpose")` |
| **Verifier** | Quality & Test Engineer | `Task(subagent_type="general-purpose")` |
| **Scribe** | Documentation Specialist | `Task(subagent_type="general-purpose")` |
| **Steward** | Architect & Final Approver | `Task(subagent_type="general-purpose")` |
| **TUIUX** | TUI/UX Expert | `Task(subagent_type="tuiux")` |

---

## Standard Sequential Workflow

```
┌──────────────────────────────────────────────────────────────────────┐
│  MILESTONE START                                                      │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────────┐
│  DIRECTOR: Analyze & Plan                                             │
│  ─────────────────────────                                            │
│  • Read ROADMAP.md milestone                                          │
│  • Read REQUIREMENTS.md (PROD-xxx, USER-xxx, NFR-xxx)                 │
│  • Create TodoList (one task in_progress at a time)                   │
│  • Determine if Steward approval needed                               │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────────┐
│  BUILDER: Implement                                                    │
│  ───────────────────                                                   │
│  • Write production-quality Go code                                    │
│  • Follow hexagonal architecture (ports & adapters)                    │
│  • Write initial tests (>80% coverage for critical paths)             │
│  • Run: go test ./..., gofmt, go vet                                   │
│                                                                        │
│  Deliverables:                                                         │
│  → Files modified (+X lines)                                           │
│  → Behaviors added (summary)                                           │
│  → Test results (passing)                                              │
│  → Implementation notes                                                │
│  → Notes for Verifier (areas needing validation)                      │
│  → Notes for Scribe (documentation updates)                           │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
        ┌─────────────────────────────────┐
        │  DIRECTOR: Review Builder       │
        │  ✅ Files modified as expected  │
        │  ✅ Tests passing              │
        │  ✅ No obvious issues          │
        └─────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────────┐
│  VERIFIER: Validate                                                    │
│  ─────────────────                                                     │
│  • Extend test coverage (add missing tests)                            │
│  • Run full test suite: go test ./... -v                               │
│  • Run race detector: go test -race ./...                              │
│  • Measure coverage: go test -cover ./...                              │
│  • Validate requirements compliance                                    │
│  • Check for regressions (baseline comparison)                         │
│                                                                        │
│  Deliverables:                                                         │
│  → Test execution summary (X/X passing)                                │
│  → Coverage metrics (>80% critical, >50% overall)                      │
│  → Requirements validation matrix                                      │
│  → Regression analysis (zero regressions ✅)                           │
│  → Recommendation: APPROVE or REQUEST FIXES                            │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
        ┌─────────────────────────────────┐
        │  DIRECTOR: Review Verifier      │
        │  ❓ All tests passing?          │
        │  ❓ Coverage targets met?       │
        │  ❓ Zero regressions?           │
        │  ❓ Verifier recommends APPROVE?│
        └─────────────────────────────────┘
                  ↓                  ↓
            [APPROVE]           [REQUEST FIXES]
                  ↓                  ↓
                  │          Re-delegate to Builder
                  │          with Verifier's findings,
                  │          then re-run Verifier
                  ↓
┌──────────────────────────────────────────────────────────────────────┐
│  SCRIBE: Document                                                      │
│  ────────────────                                                      │
│  • Update README.md (features, commands, examples)                     │
│  • Update REQUIREMENTS.md (traceability, status)                       │
│  • Update ROADMAP.md (milestone checkboxes)                            │
│  • Create/update feature guides in docs/                               │
│  • Update examples in examples/                                        │
│  • Update CHANGELOG.md                                                 │
│  • Validate cross-references (no broken links)                         │
│  • Spell-check all modified files                                      │
│                                                                        │
│  Deliverables:                                                         │
│  → Files updated (list with summaries)                                 │
│  → Files created (new guides)                                          │
│  → Examples added/updated                                              │
│  → Cross-references validated                                          │
│  → Quality checks passed                                               │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
        ┌─────────────────────────────────┐
        │  DIRECTOR: Review Scribe        │
        │  ✅ All docs updated            │
        │  ✅ Examples accurate           │
        │  ✅ Cross-references valid      │
        │  ✅ Quality checks passed       │
        └─────────────────────────────────┘
                              ↓
                  ┌───────────────────┐
                  │ Major Change or   │
                  │ Phase Gate?       │
                  └───────────────────┘
                  ↓                  ↓
              [YES]               [NO]
                  ↓                  ↓
                  │          Skip Steward
                  │          (routine change)
                  ↓
┌──────────────────────────────────────────────────────────────────────┐
│  STEWARD: Approve (Optional)                                           │
│  ────────────────────────                                              │
│  • Architecture compliance review (hexagonal boundaries)               │
│  • Security assessment (no secrets, credential management)             │
│  • Requirements validation (traceability chain)                        │
│  • Quality assessment (test coverage, test quality)                    │
│  • Documentation completeness review                                   │
│  • Make GO/NO-GO decision                                              │
│                                                                        │
│  Deliverables:                                                         │
│  → Architecture compliance report                                      │
│  → Security assessment report                                          │
│  → Requirements compliance matrix                                      │
│  → Quality assessment report                                           │
│  → Documentation completeness report                                   │
│  → Final decision: APPROVE / APPROVE WITH CONDITIONS / REJECT          │
│  → Rationale and remediation plan (if needed)                          │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
        ┌─────────────────────────────────┐
        │  DIRECTOR: Review Steward       │
        │  ❓ GO or NO-GO decision?       │
        └─────────────────────────────────┘
              ↓              ↓             ↓
        [APPROVE]  [WITH CONDITIONS]  [REJECT]
              ↓              ↓             ↓
              │              │       Re-delegate
              │              │       to fix issues,
              │              │       re-submit to
              │              │       Steward
              ↓              ↓
     Proceed (document      Proceed (document
     conditions for         conditions)
     future work)
                  ↓
┌──────────────────────────────────────────────────────────────────────┐
│  DIRECTOR: Create Git Commit(s)                                        │
│  ───────────────────────────────                                       │
│  • Create logical commits (not one giant commit)                       │
│  • Use conventional commit format:                                     │
│    type(scope): description                                            │
│  • Include comprehensive commit messages                               │
│  • Add Happy/Claude co-authorship                                      │
│                                                                        │
│  Typical structure (2-4 commits):                                      │
│  1. feat(scope): Implement [feature] with [key capability]             │
│  2. test(scope): Extend test coverage for [component]                  │
│  3. docs(scope): Add [feature] documentation and examples              │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────────┐
│  DIRECTOR: Report Milestone Completion                                 │
│  ──────────────────────────────────────                                │
│  • Milestone completion report                                         │
│  • Git commit hashes                                                   │
│  • Quality gates passed                                                │
│  • Next steps identified                                               │
│  • Mark TodoList complete (empty list)                                 │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────────┐
│  MILESTONE COMPLETE ✅                                                │
└──────────────────────────────────────────────────────────────────────┘
```

---

## Alternative Flow: TUI Visual Polish (TUIUX Agent)

For TUI visual polish tasks, the Director may invoke the TUIUX agent instead of or alongside the Builder:

```
┌──────────────────────────────────────────────────────────────────────┐
│  DIRECTOR: TUI Polish Task Identified                                  │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────────┐
│  TUIUX: Design & Implement Visual Polish                              │
│  ────────────────────────────────────────                              │
│  • Design specification (visual mockups, timing diagrams)              │
│  • Implementation (effects, widgets, themes)                           │
│  • Integration code (hooks, middleware)                                │
│  • Tests (unit, performance benchmarks <3% CPU)                        │
│  • Documentation (visual guide, config reference)                      │
│  • Demo program showcasing features                                    │
│                                                                        │
│  Four Principles of TUI Excellence:                                    │
│  1. Subtle Motion is Life                                              │
│  2. Light, Shadow, and Focus                                           │
│  3. Atmosphere and Ambient Effects                                     │
│  4. Small Charms of Quality                                            │
└──────────────────────────────────────────────────────────────────────┘
                              ↓
        ┌─────────────────────────────────┐
        │  DIRECTOR: Review TUIUX         │
        └─────────────────────────────────┘
                              ↓
             Continue with Verifier → Scribe → (Steward) → Commit
```

---

## Handoff Criteria

### Builder → Verifier
✅ Ready when:
- All targeted tests passing locally
- Code builds successfully
- Coverage targets met for new code
- `gofmt` and `go vet` clean
- Implementation summary documented

❌ NOT ready if:
- Tests failing
- Build errors
- Coverage below targets
- Linting issues unresolved

### Verifier → Scribe
✅ Ready when:
- All tests passing (100% of active tests)
- Coverage targets met (>80% critical, >50% overall)
- Zero regressions detected
- All requirements validated
- No flaky tests
- Race detector clean
- Verifier recommends APPROVE

❌ REQUEST FIXES if:
- Tests failing
- Coverage below targets
- Regressions detected
- Requirements not validated
- Flaky tests present
- Race conditions found

### Scribe → Director (or Steward)
✅ Ready when:
- All assigned documentation tasks completed
- README.md updated for user-facing changes
- REQUIREMENTS.md traceability updated
- ROADMAP.md checkboxes updated
- Examples tested or marked as hypothetical
- Cross-references validated
- Spell-check passed
- Markdown renders correctly

❌ NOT ready if:
- Examples don't match actual CLI behavior
- Broken cross-references exist
- Requirements traceability incomplete
- Roadmap checkboxes not updated
- Spelling/grammar errors present

### Steward → Director
✅ Complete review when:
- All five review areas assessed (architecture, security, requirements, quality, documentation)
- Findings documented with evidence
- Clear GO/NO-GO decision made
- Rationale provided for decision
- Remediation plan included (if rejected or conditional)

---

## Quality Gates Enforced

| Gate | Criteria | Enforced By |
|------|----------|-------------|
| Code Quality | `gofmt` clean, `go vet` clean, hexagonal architecture | Builder, Steward |
| Test Coverage | >80% critical paths, >70% service layer, >50% overall | Verifier |
| Test Quality | No regressions, race detector clean, error paths tested | Verifier |
| Documentation | README updated, guides created, requirements traced | Scribe |
| Security | No secrets, credentials in keychain, `.gitignore` proper | Steward |
| Architecture | Hexagonal boundaries respected, dependency direction correct | Steward |

---

## Workflow Rules (NEVER Violate)

❌ **Never skip Verifier** (even if Builder tests pass)
❌ **Never skip Scribe** (documentation is mandatory)
❌ **Always run agents sequentially** for same milestone (no parallel)
❌ **Complete current task before starting next**
❌ **Never write code/tests/docs yourself** (delegate to specialists)

✅ **Always create TodoList** at start of milestone
✅ **Maintain exactly ONE task in_progress**
✅ **Mark tasks complete immediately**
✅ **Review all agent deliverables**
✅ **Create logical commits with attribution**

---

## Agent Interaction Matrix

| From → To | Interaction Type | Content Exchanged |
|-----------|------------------|-------------------|
| Director → Builder | Delegation | Task description, requirements, files to modify, quality standards |
| Builder → Director | Deliverable | Implementation summary, test results, notes for Verifier/Scribe |
| Director → Verifier | Delegation | Builder's work, requirements to validate, coverage targets |
| Verifier → Director | Deliverable | Test results, coverage metrics, requirements matrix, APPROVE/REJECT |
| Director → Scribe | Delegation | Implementation summary, test results, documentation tasks |
| Scribe → Director | Deliverable | Files updated/created, examples, cross-references, quality checks |
| Director → Steward | Delegation | All agent deliverables, phase gate criteria, approval request |
| Steward → Director | Deliverable | Comprehensive review, GO/NO-GO decision, remediation plan |
| Director → TUIUX | Delegation | TUI polish tasks, Four Principles, performance budgets |
| TUIUX → Director | Deliverable | Visual polish implementation, tests, docs, demo |

---

## Success Metrics

**Director's Success Metrics:**
- [ ] All milestones follow sequential workflow
- [ ] Exactly ONE task in_progress at a time (TodoList discipline)
- [ ] Zero regressions introduced
- [ ] 100% of milestones have Verifier approval
- [ ] 100% of milestones have Scribe documentation
- [ ] All commits have proper Happy/Claude attribution

**Phase 2 Demonstrated:** This methodology delivered 5,791 lines of production-quality code with zero regressions, proving the workflow's effectiveness.

---

## Related Documentation

- **Agent Definitions:** `.agents/*.agent.md` - Complete specifications for all 6 agents
- **Director's Handbook:** `docs/DIRECTOR-HANDBOOK.md` - Complete methodology guide (1,477 lines)
- **Requirements:** `REQUIREMENTS.md` - Single source of truth
- **Roadmap:** `ROADMAP.md` - Milestone tracking
- **Architecture:** `docs/ARCHITECTURE.md` - Hexagonal architecture patterns
- **Contributing:** `CONTRIBUTING.md` - Agent roles and development methodology

---

**Document Version:** 1.0
**Status:** Active
**Maintained by:** Director

*Generated during Phase 6, Week 1 Day 4-5: Agent Definition Review*
