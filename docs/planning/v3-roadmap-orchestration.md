# Ticketr v3.0: Self-Orchestrated Roadmap

**Version:** 2.0 - Agent-Enabled
**Date:** January 2025
**Status:** Execution Ready

---

## Self-Orchestration Model

This roadmap is designed to be **self-executing** through the existing agent system in `.claude/agents/`:

```
Director → Orchestrates phases and delegates tasks
Builder → Implements code changes
Verifier → Tests and validates
Scribe → Documents changes
Steward → Reviews architecture and approves gates
```

---

## Phase 1: Foundation Layer (Weeks 1-4)

### Agent Assignments

| Week | Agent | Task | Deliverables |
|------|-------|------|-------------|
| 1 | Builder | Create SQLite adapter | `internal/adapters/database/sqlite_adapter.go` |
| 1 | Builder | Implement Repository interface | Unit tests passing |
| 2 | Builder | Migration tool (.state → SQLite) | `cmd/ticketr/migrate.go` |
| 2 | Verifier | Test migration with real data | Migration test suite |
| 3 | Builder | Feature flag system | `internal/config/features.go` |
| 3 | Scribe | Document migration process | `docs/v3-migration-guide.md` |
| 4 | Verifier | Performance benchmarks | Benchmark results < 100ms |
| 4 | Steward | Architecture review & gate approval | Phase 1 sign-off |

### Clear Deliverables

```yaml
phase_1_deliverables:
  code:
    - internal/adapters/database/sqlite_adapter.go
    - internal/adapters/database/migrations/001_initial.sql
    - cmd/ticketr/commands/migrate.go
    - internal/config/features.go

  tests:
    - internal/adapters/database/sqlite_adapter_test.go (>90% coverage)
    - tests/integration/migration_test.go
    - tests/benchmark/sqlite_benchmark_test.go

  documentation:
    - docs/v3-migration-guide.md
    - docs/api/database-adapter.md
    - CHANGELOG.md (Phase 1 section)

  acceptance:
    - All Repository interface methods implemented
    - Zero data loss in migration
    - Performance: 1000 tickets < 100ms query
    - Feature flag enables rollback
```

### Self-Orchestration Script

```bash
# .claude/orchestration/phase1.sh
#!/bin/bash

echo "Phase 1: Foundation Layer - Starting"

# Week 1
claude-agent run builder --task "Create SQLite adapter implementing Repository interface"
claude-agent run verifier --task "Unit test SQLite adapter"

# Week 2
claude-agent run builder --task "Create migration tool from .ticketr.state to SQLite"
claude-agent run verifier --task "Test migration with 1000+ ticket dataset"

# Week 3
claude-agent run builder --task "Implement feature flag system"
claude-agent run scribe --task "Document migration process and feature flags"

# Week 4
claude-agent run verifier --task "Run performance benchmarks"
claude-agent run steward --task "Review architecture and approve Phase 1 gate"

echo "Phase 1: Complete - Awaiting gate approval"
```

---

## Phase 2: Project Model (Weeks 5-8)

### Agent Assignments

| Week | Agent | Task | Deliverables |
|------|-------|------|-------------|
| 5 | Builder | Project domain model | `internal/core/domain/project.go` |
| 5 | Builder | Project service layer | `internal/core/services/project_service.go` |
| 6 | Builder | CLI commands for projects | `cmd/ticketr/commands/project.go` |
| 6 | Verifier | Test project isolation | Project isolation test suite |
| 7 | Builder | Credential management | `internal/security/keychain.go` |
| 7 | Verifier | Security audit | Security test report |
| 8 | Builder | Update commands for context | All commands project-aware |
| 8 | Scribe | Project documentation | `docs/project-model.md` |

### Clear Deliverables

```yaml
phase_2_deliverables:
  code:
    - internal/core/domain/project.go
    - internal/core/services/project_service.go
    - internal/core/ports/project_repository.go
    - internal/adapters/database/project_repository.go
    - cmd/ticketr/commands/project.go
    - internal/security/keychain.go

  database:
    - migrations/002_projects.sql
    - migrations/003_epic_hierarchy.sql

  tests:
    - internal/core/services/project_service_test.go
    - tests/integration/multi_project_test.go
    - tests/security/credential_test.go

  documentation:
    - docs/project-model.md
    - docs/epic-story-task-hierarchy.md
    - docs/multi-jira-setup.md
```

---

## Phase 3: Global Installation (Weeks 9-10)

### Agent Assignments

| Week | Agent | Task | Deliverables |
|------|-------|------|-------------|
| 9 | Builder | XDG path compliance | `internal/core/services/path_resolver.go` |
| 9 | Builder | Package manifests | `brew/`, `snap/`, `scoop/` formulas |
| 10 | Scribe | Installation guides | Platform-specific install docs |
| 10 | Verifier | Cross-platform testing | Test results for all platforms |

### Clear Deliverables

```yaml
phase_3_deliverables:
  code:
    - internal/core/services/path_resolver.go
    - internal/config/xdg.go
    - cmd/ticketr/setup/first_run.go

  packages:
    - dist/homebrew/ticketr.rb
    - dist/snap/snapcraft.yaml
    - dist/scoop/ticketr.json
    - dist/debian/ticketr.deb

  ci:
    - .github/workflows/release.yml
    - .github/workflows/package-test.yml

  documentation:
    - docs/installation/macos.md
    - docs/installation/linux.md
    - docs/installation/windows.md
```

---

## Phase 4: TUI Implementation (Weeks 11-16)

### Agent Assignments

| Week | Agent | Task | Deliverables |
|------|-------|------|-------------|
| 11 | Builder | TUI adapter skeleton | `internal/adapters/tui/app.go` |
| 12 | Builder | Main dashboard (wireframe 1) | `internal/adapters/tui/views/dashboard.go` |
| 13 | Builder | Ticket detail view (wireframe 2) | `internal/adapters/tui/views/detail.go` |
| 14 | Builder | Search & command (wireframe 3) | `internal/adapters/tui/views/search.go` |
| 15 | Builder | Conflict resolution (wireframe 4) | `internal/adapters/tui/views/conflict.go` |
| 15 | Builder | Sync status (wireframe 5) | `internal/adapters/tui/views/sync.go` |
| 16 | Verifier | TUI integration tests | Complete TUI test suite |
| 16 | Scribe | TUI user guide | `docs/tui-guide.md` |

### Clear Deliverables

```yaml
phase_4_deliverables:
  code:
    - internal/adapters/tui/app.go
    - internal/adapters/tui/router.go
    - internal/adapters/tui/keybindings.go
    - internal/adapters/tui/themes/
    - internal/adapters/tui/views/
      - dashboard.go      # Implements wireframe 1
      - detail.go         # Implements wireframe 2
      - search.go         # Implements wireframe 3
      - conflict.go       # Implements wireframe 4
      - sync_status.go    # Implements wireframe 5
      - project_switch.go # Implements wireframe 6
      - create_wizard.go  # Implements wireframe 7

  tests:
    - internal/adapters/tui/app_test.go
    - tests/integration/tui_navigation_test.go
    - tests/e2e/tui_workflow_test.go

  assets:
    - internal/adapters/tui/assets/help.txt
    - internal/adapters/tui/assets/shortcuts.yaml

  documentation:
    - docs/tui-guide.md
    - docs/tui-keybindings.md
    - docs/tui-customization.md
```

### TUI Wireframe Implementation Map

```go
// Each wireframe from docs/tui-wireframes.md maps to a view:

// Wireframe 1: Main Dashboard
type DashboardView struct {
    projects  *ProjectList    // Left panel
    hierarchy *TicketTree     // Center panel
    stats     *Statistics     // Bottom panel
}

// Wireframe 2: Ticket Detail
type DetailView struct {
    hierarchy *MiniTree       // Left context
    form      *TicketForm     // Right editor
    metadata  *MetadataPanel  // Left metadata
}

// Wireframe 3: Search & Command
type SearchView struct {
    input    *SearchInput
    results  *ResultList
    preview  *PreviewPane
}

// Wireframe 4: Conflict Resolution
type ConflictView struct {
    local   *DiffPanel
    remote  *DiffPanel
    merged  *EditablePanel
    actions *ResolutionActions
}

// Wireframe 5: Sync Status
type SyncView struct {
    progress *ProgressBar
    log      *OperationLog
    stats    *SyncStats
    errors   *ErrorList
}
```

---

## Phase 5: Advanced Features (Weeks 17-20)

### Agent Assignments

| Week | Agent | Task | Deliverables |
|------|-------|------|-------------|
| 17 | Builder | Bulk operations | `internal/core/services/bulk_service.go` |
| 18 | Builder | Template system | `internal/templates/engine.go` |
| 19 | Builder | Smart sync strategies | `internal/sync/strategies/` |
| 20 | Verifier | Full system test | Complete test report |
| 20 | Steward | Final architecture review | v3.0 approval |

---

## Agent Coordination Protocol

### Director's Master Script

```python
# .claude/orchestration/master.py

phases = [
    {"id": 1, "name": "Foundation", "weeks": 4},
    {"id": 2, "name": "Projects", "weeks": 4},
    {"id": 3, "name": "Global", "weeks": 2},
    {"id": 4, "name": "TUI", "weeks": 6},
    {"id": 5, "name": "Advanced", "weeks": 4}
]

for phase in phases:
    print(f"Starting Phase {phase['id']}: {phase['name']}")

    # Delegate to Builder
    builder_tasks = get_builder_tasks(phase['id'])
    for task in builder_tasks:
        execute_agent('builder', task)

    # Delegate to Verifier
    verifier_tasks = get_verifier_tasks(phase['id'])
    for task in verifier_tasks:
        execute_agent('verifier', task)

    # Delegate to Scribe
    scribe_tasks = get_scribe_tasks(phase['id'])
    for task in scribe_tasks:
        execute_agent('scribe', task)

    # Gate check with Steward
    gate_result = execute_agent('steward', f"Review Phase {phase['id']}")

    if not gate_result['approved']:
        print(f"Phase {phase['id']} failed gate check")
        break
```

### Task Specification Format

Each agent receives tasks in this format:

```yaml
task:
  id: "P1W1-001"
  phase: 1
  week: 1
  agent: builder
  title: "Create SQLite adapter"

  requirements:
    - Implement Repository interface
    - Support transactions
    - Handle migrations

  inputs:
    - internal/core/ports/repository.go
    - docs/v3-technical-specification.md#database-design

  outputs:
    - internal/adapters/database/sqlite_adapter.go
    - internal/adapters/database/sqlite_adapter_test.go

  validation:
    - Unit tests pass
    - Coverage > 90%
    - Benchmark < 100ms for 1000 records

  dependencies:
    - none

  estimated_hours: 8
```

---

## Success Metrics Dashboard

```
Phase 1: Foundation Layer
├── Code Complete:     [████████████████████] 100%
├── Tests Passing:     [████████████████████] 100%
├── Docs Updated:      [████████████████████] 100%
├── Gate Approved:     ✅ Steward
└── Status:           COMPLETE

Phase 2: Project Model
├── Code Complete:     [████████████░░░░░░░] 65%
├── Tests Passing:     [████████░░░░░░░░░░░░] 40%
├── Docs Updated:      [████░░░░░░░░░░░░░░░░] 20%
├── Gate Approved:     ⏳ Pending
└── Status:           IN PROGRESS

Phase 3: Global Installation
├── Code Complete:     [░░░░░░░░░░░░░░░░░░░░] 0%
├── Tests Passing:     [░░░░░░░░░░░░░░░░░░░░] 0%
├── Docs Updated:      [░░░░░░░░░░░░░░░░░░░░] 0%
├── Gate Approved:     ⏸ Waiting
└── Status:           QUEUED
```

---

## Communication Protocol

### Daily Standup Format

```markdown
## Director's Daily Report

### Yesterday
- ✅ Builder completed SQLite adapter
- ✅ Verifier validated migration tool
- ⚠️ Scribe blocked on API documentation

### Today
- 🔨 Builder: Start project domain model
- 🧪 Verifier: Complete performance benchmarks
- 📝 Scribe: Finish migration guide

### Blockers
- Need Steward approval for database schema changes
- Waiting on credential storage decision

### Metrics
- Velocity: 8 story points/week
- Test Coverage: 87%
- Documentation: 65% complete
```

### Escalation Path

```
1. Task-level issues → Director
2. Architecture concerns → Steward
3. Resource conflicts → Director → Human
4. Gate failures → Steward → Human
```

---

## Automation Triggers

```yaml
# .github/workflows/v3-automation.yml
on:
  schedule:
    - cron: '0 9 * * 1-5'  # Daily at 9 AM weekdays

jobs:
  daily_orchestration:
    runs-on: ubuntu-latest
    steps:
      - name: Run Director
        run: |
          claude-agent run director \
            --task "Check phase progress" \
            --output progress-report.md

      - name: Execute pending tasks
        run: |
          claude-agent orchestrate \
            --config .claude/orchestration/config.yaml \
            --phase current

      - name: Update metrics
        run: |
          claude-agent run verifier \
            --task "Update test coverage metrics"
```

---

## Deliverable Validation Checklist

### Per-Phase Checklist

```markdown
## Phase 1 Validation
- [ ] SQLite adapter implements all Repository methods
- [ ] Migration preserves 100% of data
- [ ] Feature flags tested in both states
- [ ] Performance benchmarks documented
- [ ] User migration guide complete
- [ ] API documentation generated
- [ ] Integration tests passing
- [ ] No regression in v2 functionality
```

### Final Deliverable Structure

```
ticketr-v3.0.0/
├── bin/
│   └── ticketr              # Single binary
├── docs/
│   ├── installation/        # Platform guides
│   ├── migration/           # v2 → v3 guide
│   ├── tui/                 # TUI documentation
│   ├── api/                 # API reference
│   └── tutorials/           # User tutorials
├── examples/
│   ├── projects/            # Sample projects
│   ├── templates/           # Ticket templates
│   └── workflows/           # Common workflows
├── tests/
│   ├── unit/                # Unit tests
│   ├── integration/         # Integration tests
│   ├── e2e/                 # End-to-end tests
│   └── benchmarks/          # Performance tests
└── dist/
    ├── homebrew/            # macOS package
    ├── snap/                # Linux package
    ├── scoop/               # Windows package
    └── docker/              # Container image
```

---

## Conclusion

This roadmap is now:

1. **Self-Orchestrated**: Agents have clear assignments
2. **Measurable**: Specific deliverables per phase
3. **Testable**: Validation criteria defined
4. **Trackable**: Progress metrics included
5. **Executable**: Can be run autonomously

The Director agent can now take this document and orchestrate the entire v3.0 implementation with minimal human intervention.