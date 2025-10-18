# Ticketr v3.0 Implementation Roadmap

**Document Version:** 1.0
**Created:** January 2025
**Status:** Planning Phase
**Goal:** Transform Ticketr from directory-bound tool to global work platform

---

## Overview

This document defines the gated implementation path for Ticketr v3.0, with clear milestones, acceptance criteria, and rollback points. Each phase builds on the previous, with feature flags enabling progressive rollout.

## Guiding Principles

1. **Backward Compatibility**: v2.x workflows must continue working
2. **Progressive Enhancement**: Each phase delivers value independently
3. **Feature Flags**: New features behind flags until stable
4. **Data Safety**: No data loss during migration
5. **Incremental Delivery**: 2-week sprints, monthly releases

---

## Phase 1: Foundation Layer (Weeks 1-4)

### Milestone 14: Centralized State Management

**Goal**: Replace file-based state with SQLite while maintaining backward compatibility.

#### Technical Requirements

```go
// internal/adapters/database/sqlite_adapter.go
type SQLiteAdapter struct {
    db *sql.DB
    path string // ~/.config/ticketr/ticketr.db
}

// Implements existing Repository interface
func (s *SQLiteAdapter) GetTickets(filter string) ([]domain.Ticket, error)
func (s *SQLiteAdapter) SaveTickets(tickets []domain.Ticket) error
func (s *SQLiteAdapter) GetState(ticketID string) (*StateRecord, error)
func (s *SQLiteAdapter) UpdateState(ticketID string, state StateRecord) error
```

#### Database Schema v1

```sql
-- Initial schema focused on state management
CREATE TABLE IF NOT EXISTS schema_version (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tickets (
    id TEXT PRIMARY KEY,
    jira_id TEXT UNIQUE,
    title TEXT NOT NULL,
    description TEXT,
    content JSON, -- Full ticket data
    local_hash TEXT,
    remote_hash TEXT,
    last_modified TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_jira_id ON tickets(jira_id);
CREATE INDEX idx_modified ON tickets(last_modified);

CREATE TABLE IF NOT EXISTS sync_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    operation TEXT NOT NULL, -- 'push', 'pull', 'conflict'
    ticket_id TEXT,
    status TEXT, -- 'success', 'failed', 'partial'
    details JSON,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Implementation Steps

1. **Week 1**: Create SQLite adapter implementing Repository interface
2. **Week 2**: Add migration tool for .ticketr.state → SQLite
3. **Week 3**: Implement feature flag system
4. **Week 4**: Integration testing & performance benchmarks

#### Feature Flags

```yaml
# .ticketr.yaml
features:
  use_sqlite: false  # Default to file-based for compatibility
  sqlite_path: ~/.config/ticketr/ticketr.db
```

#### Acceptance Criteria

- [ ] SQLite adapter passes all existing Repository tests
- [ ] Migration tool converts .ticketr.state without data loss
- [ ] Performance: < 100ms for 1000 ticket queries
- [ ] Backward compatible with v2.x file-based state
- [ ] Zero changes required to existing CLI commands

#### Rollback Plan

- Feature flag `use_sqlite: false` reverts to file-based state
- Export command: `ticketr export --to-file`
- State files remain untouched during SQLite operations

---

## Phase 2: Workspace Model (Weeks 5-8)

### Milestone 15: Multi-Workspace Support

**Goal**: Enable managing multiple Jira projects from a single installation.

#### Technical Requirements

```go
// internal/core/domain/workspace.go
type Workspace struct {
    ID          string
    Name        string
    JiraURL     string
    ProjectKey  string
    Credentials CredentialRef // Reference to secure storage
    IsDefault   bool
    LastUsed    time.Time
}

// internal/core/services/workspace_service.go
type WorkspaceService struct {
    repo WorkspaceRepository
    current *Workspace
}

func (w *WorkspaceService) Create(name string, config WorkspaceConfig) error
func (w *WorkspaceService) Switch(name string) error
func (w *WorkspaceService) List() ([]Workspace, error)
func (w *WorkspaceService) Current() *Workspace
```

#### Database Schema v2 (Additions)

```sql
CREATE TABLE IF NOT EXISTS workspaces (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    jira_url TEXT NOT NULL,
    project_key TEXT NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    last_used TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Update tickets table
ALTER TABLE tickets ADD COLUMN workspace_id TEXT REFERENCES workspaces(id);
CREATE INDEX idx_workspace ON tickets(workspace_id);

-- Credentials in system keychain, not DB
CREATE TABLE IF NOT EXISTS credential_refs (
    workspace_id TEXT PRIMARY KEY REFERENCES workspaces(id),
    keychain_ref TEXT NOT NULL -- Reference to OS keychain
);
```

#### CLI Commands

```bash
# Workspace management
ticketr workspace create backend --url https://company.atlassian.net --project BACK
ticketr workspace list
ticketr workspace switch frontend
ticketr workspace current
ticketr workspace remove mobile

# Context-aware operations
ticketr push tickets.md --workspace backend
ticketr pull --workspace frontend --output frontend-tickets.md
```

#### Implementation Steps

1. **Week 5**: Workspace domain model and service layer
2. **Week 6**: CLI commands for workspace management
3. **Week 7**: Credential management via OS keychain
4. **Week 8**: Update existing commands for workspace context

#### Acceptance Criteria

- [ ] Can create/switch/delete workspaces
- [ ] Credentials stored securely in OS keychain
- [ ] Existing single-project workflows still work
- [ ] --workspace flag available on all commands
- [ ] Default workspace auto-selected
- [ ] No credentials in database or logs

---

## Phase 3: Global Installation (Weeks 9-10)

### Milestone 16: System-Wide Tool

**Goal**: Transform Ticketr into a globally-installed tool that works from any directory.

#### Technical Requirements

```go
// internal/core/services/path_resolver.go
type PathResolver struct {
    configHome string // XDG_CONFIG_HOME or %APPDATA%
    dataHome   string // XDG_DATA_HOME or %LOCALAPPDATA%
}

func (p *PathResolver) ConfigDir() string // ~/.config/ticketr
func (p *PathResolver) DataDir() string   // ~/.local/share/ticketr
func (p *PathResolver) CacheDir() string  // ~/.cache/ticketr
```

#### Directory Structure

```
$HOME/
├── .config/ticketr/
│   ├── config.yaml        # User configuration
│   └── workspaces.yaml    # Workspace definitions
├── .local/share/ticketr/
│   ├── ticketr.db        # SQLite database
│   ├── templates/        # User templates
│   └── plugins/          # Future: plugins
└── .cache/ticketr/
    ├── jira_schema.json  # Cached field mappings
    └── logs/             # Rotated logs
```

#### Installation Methods

```bash
# macOS (Homebrew)
brew install karolswdev/tap/ticketr

# Linux (snap)
snap install ticketr

# Windows (scoop)
scoop bucket add ticketr https://github.com/karolswdev/scoop-ticketr
scoop install ticketr

# Cross-platform (Go)
go install github.com/karolswdev/ticktr/cmd/ticketr@v3
```

#### Implementation Steps

1. **Week 9**: XDG/Windows path compliance
2. **Week 9**: Package manifests (brew, snap, scoop)
3. **Week 10**: Installation documentation
4. **Week 10**: First-run setup wizard

#### Acceptance Criteria

- [ ] Single binary in system PATH
- [ ] Works from any directory
- [ ] Respects XDG Base Directory spec (Linux/macOS)
- [ ] Respects Windows conventions (%APPDATA%)
- [ ] Auto-creates config directories on first run
- [ ] Migration tool for v2.x users

---

## Phase 4: TUI Implementation (Weeks 11-16)

### Milestone 17: Terminal User Interface

**Goal**: Implement TUI as a first-class adapter in the hexagonal architecture.

#### Technical Requirements

```go
// internal/adapters/tui/app.go
type TUIApp struct {
    app         *tview.Application
    core        ports.TicketService
    workspace   ports.WorkspaceService
    currentView View
}

// internal/adapters/tui/views/
├── workspace_list.go    # Workspace switcher
├── ticket_tree.go       # Hierarchical ticket view
├── ticket_detail.go     # Single ticket editor
├── search_results.go    # Search/filter results
└── sync_status.go       # Real-time sync monitor
```

#### TUI Architecture

```
┌─────────────────────────────────────────┐
│           TUI Adapter Layer             │
├─────────────────────────────────────────┤
│  ┌──────────┐ ┌──────────┐ ┌─────────┐ │
│  │  Router  │ │ Commands │ │ Themes  │ │
│  └──────────┘ └──────────┘ └─────────┘ │
├─────────────────────────────────────────┤
│            View Components              │
│  ┌──────────────────────────────────┐  │
│  │ TreeView │ DetailForm │ StatusBar │  │
│  └──────────────────────────────────┘  │
├─────────────────────────────────────────┤
│           Core Services                 │
│   (Same services used by CLI adapter)  │
└─────────────────────────────────────────┘
```

#### Key Bindings

```yaml
global:
  'q': quit
  'ctrl+c': quit
  '/': search
  ':': command
  '?': help
  'tab': next_panel
  'shift+tab': prev_panel

navigation:
  'j': down
  'k': up
  'h': collapse
  'l': expand
  'g': top
  'G': bottom
  'ctrl+f': page_down
  'ctrl+b': page_up

operations:
  'c': create
  'e': edit
  'd': delete
  'p': push
  'P': pull
  'r': refresh
  's': sync
  'space': select
  'a': select_all
```

#### Implementation Steps

1. **Week 11**: TUI adapter skeleton with tview
2. **Week 12**: Workspace and ticket list views
3. **Week 13**: Ticket detail editor with validation
4. **Week 14**: Search, filter, and command palette
5. **Week 15**: Real-time sync status and indicators
6. **Week 16**: Keybindings, themes, and polish

#### Acceptance Criteria

- [ ] Launch with `ticketr tui` from anywhere
- [ ] All CLI operations available in TUI
- [ ] Vim-style navigation
- [ ] Real-time sync status
- [ ] No blocking operations (async everything)
- [ ] Graceful degradation on small terminals
- [ ] Help system with key binding reference

---

## Phase 5: Advanced Features (Weeks 17-20)

### Milestone 18: Enhanced Capabilities

**Goal**: Add power-user features leveraging the new architecture.

#### Feature Set

1. **Bulk Operations**
   ```go
   type BulkOperation struct {
       Action  string   // "update", "move", "delete"
       Tickets []string // Ticket IDs
       Changes map[string]interface{}
   }
   ```

2. **Templates System**
   ```yaml
   # ~/.config/ticketr/templates/feature.yaml
   name: feature
   structure:
     epic:
       title: "Feature: {{.Name}}"
       description: |
         As a {{.Actor}}
         I want {{.Goal}}
         So that {{.Benefit}}
     stories:
       - title: "Frontend: {{.Name}}"
         tasks:
           - "Component implementation"
           - "Unit tests"
           - "Integration tests"
   ```

3. **Smart Sync**
   ```go
   type SyncStrategy interface {
       ShouldSync(ticket Ticket) bool
       ResolveConflict(local, remote Ticket) Ticket
   }
   ```

4. **JQL Aliases**
   ```yaml
   aliases:
     mine: "assignee = currentUser() AND status != Done"
     sprint: "sprint in openSprints()"
     blocked: "status = Blocked OR labels = blocked"
   ```

#### Implementation Steps

1. **Week 17**: Bulk operations in TUI
2. **Week 18**: Template system
3. **Week 19**: Smart sync with strategies
4. **Week 20**: JQL aliases and quick filters

#### Acceptance Criteria

- [ ] Can select multiple tickets for bulk update
- [ ] Templates reduce creation time by 50%
- [ ] Conflict resolution without data loss
- [ ] Aliases work in CLI and TUI

---

## Testing Strategy

### Unit Tests
- Each adapter tested independently
- Core services have 100% coverage
- Database migrations are reversible

### Integration Tests
```go
// tests/integration/v3_test.go
func TestV2BackwardCompatibility(t *testing.T)
func TestWorkspaceSwitching(t *testing.T)
func TestTUINavigation(t *testing.T)
func TestBulkOperations(t *testing.T)
```

### End-to-End Tests
```bash
# tests/e2e/scenarios/
├── new_user_setup.sh
├── v2_migration.sh
├── multi_workspace.sh
└── tui_workflow.sh
```

### Performance Benchmarks
- Startup time: < 100ms
- Workspace switch: < 50ms
- 1000 ticket query: < 100ms
- TUI refresh: < 16ms (60fps)

---

## Migration Strategy

### For v2.x Users

```bash
# Automatic migration on first v3 run
$ ticketr v3 migrate
Migrating from v2.x to v3.0...
✓ Found 3 projects with .ticketr.state files
✓ Creating workspace 'backend' from ~/projects/backend
✓ Creating workspace 'frontend' from ~/projects/frontend
✓ Creating workspace 'mobile' from ~/projects/mobile
✓ Migrating 847 tickets to SQLite
✓ Migration complete!

Your workspaces:
- backend (default)
- frontend
- mobile

Run 'ticketr workspace list' to see all workspaces
Run 'ticketr tui' to launch the new interface
```

### Data Preservation

1. Original files untouched
2. State files backed up to ~/.config/ticketr/backups/
3. Rollback command available: `ticketr v3 rollback`

---

## Release Strategy

### Version Tags
- v3.0-alpha.1: Phase 1 complete (SQLite)
- v3.0-beta.1: Phase 3 complete (Global tool)
- v3.0-rc.1: Phase 4 complete (TUI)
- v3.0.0: Phase 5 complete (Full feature set)

### Feature Flags Evolution

```yaml
# v3.0-alpha (opt-in)
features:
  use_sqlite: false
  enable_workspaces: false
  enable_tui: false

# v3.0-beta (mixed)
features:
  use_sqlite: true        # Now default
  enable_workspaces: true  # Now default
  enable_tui: false       # Still experimental

# v3.0.0 (all enabled)
# No more feature flags needed
```

### Communication Plan

1. **Blog Post**: "Ticketr v3: From Tool to Platform"
2. **Migration Guide**: Step-by-step for v2 users
3. **Video Demo**: TUI workflow demonstration
4. **Discord/Slack**: Community support channel

---

## Success Metrics

### Quantitative
- Installation success rate > 95%
- Migration success rate > 99%
- Performance regression < 5%
- Test coverage > 80%

### Qualitative
- "It just works" from any directory
- TUI feels as natural as k9s/lazygit
- Zero data loss reports
- Power users adopt TUI as primary interface

---

## Risk Mitigation

### Technical Risks

| Risk | Mitigation |
|------|------------|
| SQLite corruption | Write-ahead logging, regular backups |
| TUI complexity | Start minimal, progressive enhancement |
| Breaking changes | Feature flags, backward compatibility |
| Platform differences | CI testing on macOS/Linux/Windows |

### User Risks

| Risk | Mitigation |
|------|------------|
| Learning curve | Comprehensive docs, video tutorials |
| Migration failures | Automatic backups, rollback command |
| Muscle memory | v2 commands still work |

---

## Go/No-Go Criteria

### Phase Gates

Each phase must meet these criteria before proceeding:

1. **All tests passing** (unit, integration, e2e)
2. **Documentation complete** (user guide, API docs)
3. **Performance benchmarks met**
4. **No critical bugs** in previous phase
5. **User feedback incorporated** (alpha/beta testers)

### Rollback Triggers

Automatic rollback if:
- Data loss detected
- Performance regression > 20%
- Critical security vulnerability
- > 5% of users report failures

---

## Timeline Summary

```
Week 1-4:   Phase 1 - Foundation (SQLite)
Week 5-8:   Phase 2 - Workspaces
Week 9-10:  Phase 3 - Global Tool
Week 11-16: Phase 4 - TUI
Week 17-20: Phase 5 - Advanced Features

Milestones:
- Month 1: Alpha release (Phases 1-2)
- Month 2: Beta release (Phase 3)
- Month 3: RC release (Phase 4)
- Month 4: v3.0.0 GA (Phase 5)
```

---

## Appendix A: Technical Decisions

### Why SQLite?
- Embedded, no server required
- ACID compliant
- Excellent Go support
- Single file deployment
- Can migrate to PostgreSQL later if needed

### Why tview for TUI?
- Pure Go (no CGO dependencies)
- Cross-platform
- Rich widget set
- Active maintenance
- Good documentation

### Why XDG compliance?
- Standard locations users expect
- Clean home directories
- Backup-friendly
- Follows platform conventions

---

## Appendix B: Alternative Approaches Considered

1. **Web UI instead of TUI**: Rejected - adds complexity, requires browser
2. **PostgreSQL from start**: Rejected - overkill for single-user tool
3. **Keep file-based state**: Rejected - doesn't scale, can't query efficiently
4. **Electron app**: Rejected - large binary, not CLI-first

---

*This roadmap is a living document. Updates tracked in git history.*