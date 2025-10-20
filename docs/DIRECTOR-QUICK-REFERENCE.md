# Director's Quick Reference Guide

**Version:** 2.0 | **Date:** October 19, 2025

*Companion to DIRECTOR-ORCHESTRATION-FRAMEWORK.md*

---

## 1-Page Cheat Sheet

### Standard Milestone Workflow

```
1. PLAN    → Read roadmap, create TodoList (10-20 min)
2. BUILD   → Delegate to Builder (2-4 hours)
3. VERIFY  → Delegate to Verifier (30-60 min)
4. DOCUMENT → Delegate to Scribe (30-60 min)
5. COMMIT  → Create git commits (5-10 min)
6. CLOSE   → Update roadmap, mark complete (5 min)
```

### Quality Gates (NEVER SKIP)

```
Builder   → Code compiles, tests pass, architecture compliant
Verifier  → Coverage ≥70%, zero regressions, APPROVE recommendation
Scribe    → All docs updated, examples tested, cross-refs valid
Steward   → (Optional) Architecture approved, security assessed
```

### Agent Invocation Pattern

```python
Task(
    subagent_type="general-purpose",
    description="[5-word description]",
    prompt="""You are the [Agent Name] for [Context].

## Mission
[What to accomplish]

## Context
[Background info]

## Tasks
1. [Specific task]
2. [Specific task]

## Expected Deliverables
[Exact format]
"""
)
```

---

## Decision Trees

### Which Agent to Use?

```
Need code written? → Builder
Need tests validated? → Verifier
Need docs created? → Scribe
Need architecture review? → Steward
Need research? → General-purpose (rare)
```

### When to Escalate to Human?

```
Blocked >30 min? → YES
Missing credentials? → YES
Architectural uncertainty? → YES
Agent fails 3+ times? → YES
Routine rework? → NO
```

### How Many Git Commits?

```
Small feature (<200 lines)? → 1 commit
Large feature (>500 lines)? → 2-3 commits
  - feat: Implementation
  - test: Extended coverage
  - docs: Documentation
```

---

## Common Commands

### Pre-Execution Planning

```bash
# Read roadmap
Read: docs/v3-implementation-roadmap.md

# Check current state
Bash: git status
Bash: git log -5 --oneline

# Find relevant files
Glob: internal/core/services/*.go
Grep: pattern="Feature" output_mode="files_with_matches"

# Create TodoList
TodoWrite(todos=[...])
```

### Builder Delegation

```python
Task(
    subagent_type="general-purpose",
    description="Implement [feature]",
    prompt="""Builder: Implement [feature]

Files to create/modify: [list]
Requirements: [criteria]
Quality standards: tests pass, >70% coverage
Deliverables: code, test results, notes for Verifier
"""
)
```

### Verifier Delegation

```python
Task(
    subagent_type="general-purpose",
    description="Verify [feature]",
    prompt="""Verifier: Validate [feature]

Run: go test ./... -v
Check coverage: >70% for new code
Validate requirements: [list]
Deliverables: test results, coverage, APPROVE/REQUEST FIXES
"""
)
```

### Scribe Delegation

```python
Task(
    subagent_type="general-purpose",
    description="Document [feature]",
    prompt="""Scribe: Document [feature]

Update: README.md, docs/[guide].md, CHANGELOG.md, roadmap
Test examples: all commands must work
Deliverables: file diffs, cross-ref validation
"""
)
```

### Git Commit

```bash
git add [files]
git commit -m "$(cat <<'EOF'
feat(scope): Brief description

Detailed explanation.

Implementation:
- Key point 1
- Key point 2

Testing:
- [X] tests passing
- [Y]% coverage

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
EOF
)"
```

---

## Quality Standards

### Test Coverage Targets

```
Overall: ≥70%
Services: ≥80%
Repositories: ≥80%
Critical paths: 100%
```

### Performance Benchmarks

```
Startup: <100ms
Workspace switch: <50ms
1000 ticket query: <100ms
TUI rendering: <16ms (60fps)
```

### Documentation Requirements

```
User-facing change? → Update README + guide
CLI command added? → Add to command reference
TUI feature added? → Update help + keybindings
Milestone complete? → Update roadmap + CHANGELOG
```

---

## TodoList Management

### Good Task Format

```python
{
    "content": "Delegate BulkOperation implementation to Builder",
    "activeForm": "Delegating implementation to Builder",
    "status": "pending"
}
```

### Status Rules

```
✅ Exactly ONE task "in_progress" at a time
✅ Mark "completed" IMMEDIATELY after finishing
✅ Update TodoList after EVERY delegation/review
❌ Never multiple tasks "in_progress"
❌ Never mark "completed" prematurely
```

---

## Troubleshooting Quick Fixes

### Builder Returns Failing Tests

```
1. Read exact error messages
2. Diagnose: logic bug, test bug, or regression?
3. Delegate back to Builder with specific fix
4. Do NOT proceed to Verifier until passing
```

### Verifier Finds Regressions

```
1. Identify which tests regressed
2. Determine: true regression or legitimate change?
3. If regression: Builder fixes
4. If legitimate: Builder updates tests
5. Re-run Verifier after fixes
```

### Scribe Docs Inaccurate

```
1. Identify specific inaccuracies
2. Delegate back to Scribe with corrections
3. Request example re-testing
4. Verify cross-references
```

---

## Phase 5 Week-by-Week

### Week 18: Bulk Operations (5 days)

```
Day 1: Domain model
Day 2: Service implementation
Day 3: CLI integration
Day 4-5: TUI integration + docs
```

### Week 19: Templates + Smart Sync (5 days)

```
Day 1-2: Template system
Day 3-5: Smart sync + conflict resolution + docs
```

### Week 20: JQL Aliases + Polish (5 days)

```
Day 1-2: JQL aliases
Day 3: Performance optimization
Day 4-5: Final polish + full regression + docs
```

---

## Key File Locations

```
Roadmap: docs/v3-implementation-roadmap.md
Requirements: REQUIREMENTS-v2.md
Architecture: docs/ARCHITECTURE.md
Director Guide: docs/DIRECTOR-HANDBOOK.md
Framework: docs/DIRECTOR-ORCHESTRATION-FRAMEWORK.md

Services: internal/core/services/
Adapters: internal/adapters/
Domain: internal/core/domain/
CLI: cmd/ticketr/
TUI: internal/adapters/tui/
```

---

## Quality Checklist (Pre-Commit)

```
[ ] All tests passing: go test ./...
[ ] Build successful: go build ./...
[ ] Coverage ≥70%: go tool cover -func=coverage.out
[ ] Documentation updated
[ ] Roadmap checkboxes marked
[ ] TodoList cleared
[ ] Git commits created with co-authorship
[ ] Working directory clean: git status
```

---

## Emergency Contacts

```
Blocked? → Escalate to human with context
Security issue? → Escalate immediately
Data loss risk? → Halt, escalate
Timeline at risk? → Communicate early
```

---

## Remember

```
✅ Trust the process
✅ Follow the sequence
✅ Enforce quality gates
✅ Document everything
✅ One task at a time
```

**The methodology works. Proven across 10,000+ LOC delivered.**

---

*For detailed guidance, see DIRECTOR-ORCHESTRATION-FRAMEWORK.md*
