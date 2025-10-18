# Ticketr Vision: From Directory Tool to Global Work Platform

**Version:** 3.0 Vision Document
**Date:** January 2025
**Status:** Strategic Direction

## The Paradigm Shift

Ticketr must evolve from a **directory-bound tool** to a **global work management platform** that lives in your system, not your project.

### Current Reality (v2.x)
```bash
cd ~/projects/backend
ticketr push tickets.md  # Creates .ticketr.state HERE
cd ~/projects/frontend
ticketr push tickets.md  # Different .ticketr.state HERE
# No connection between projects!
```

### Future Vision (v3.x)
```bash
# From ANYWHERE in your system:
ticketr workspace list              # See all your workspaces
ticketr tui                         # Launch TUI from anywhere
ticketr push ~/any/path/tickets.md  # Works globally
ticketr switch backend              # Context switching
```

## Core Architecture Evolution

### 1. Global Installation & State

```
$HOME/
├── .config/ticketr/           # Or %APPDATA%\ticketr on Windows
│   ├── config.yaml            # Global configuration
│   ├── ticketr.db             # SQLite database (ALL state)
│   ├── plugins/               # Plugin directory
│   └── themes/                # TUI themes
```

**No more `.ticketr.state` files scattered everywhere!**

### 2. Centralized SQLite Database

```sql
-- Core schema
CREATE TABLE workspaces (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE,
    path TEXT,  -- Optional: can be NULL for pure-remote workspaces
    jira_url TEXT,
    project_key TEXT,
    last_accessed TIMESTAMP
);

CREATE TABLE tickets (
    id TEXT PRIMARY KEY,
    workspace_id TEXT REFERENCES workspaces(id),
    jira_id TEXT,
    local_hash TEXT,
    remote_hash TEXT,
    content JSON,
    last_synced TIMESTAMP,
    INDEX idx_workspace (workspace_id),
    INDEX idx_jira_id (jira_id)
);

CREATE TABLE sync_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workspace_id TEXT,
    operation TEXT, -- 'push', 'pull', 'conflict'
    ticket_id TEXT,
    timestamp TIMESTAMP,
    details JSON
);

CREATE TABLE user_preferences (
    key TEXT PRIMARY KEY,
    value JSON
);
```

### 3. Workspace Model

```go
type WorkspaceManager struct {
    db *sql.DB
    current *Workspace
}

type Workspace struct {
    ID         string
    Name       string
    Path       string    // Optional - can work purely remote
    JiraURL    string
    ProjectKey string
    Filters    []Filter  // Saved JQL queries
    Templates  []Template
}

// Global context switching
func (wm *WorkspaceManager) Switch(name string) error {
    wm.current = wm.GetWorkspace(name)
    return nil
}

// Work from anywhere
func (wm *WorkspaceManager) Push(filepath string) error {
    // If filepath is absolute, use it
    // If relative, check current directory
    // If no file specified, look for *.md in current dir
    // All state goes to central DB, not local files
}
```

## The TUI Revolution

### TUI as First-Class Citizen

The TUI isn't just a feature—it's the **primary interface** for complex operations, with CLI as the automation/scripting layer.

```go
// internal/adapters/tui/
type TUIAdapter struct {
    app        *tview.Application
    workspace  *WorkspaceManager
    jira       ports.JiraPort
}

// Launch from anywhere
$ ticketr tui
```

### TUI Layout Vision

```
┌─────────────────────────────────────────────────────────────┐
│ Ticketr v3.0 | Workspace: backend | Connected: JIRA         │
├──────────────┬──────────────────────────────────────────────┤
│              │                                              │
│ ▼ Workspaces │  EPIC-100: Authentication System            │
│   ● backend  │  ├── STORY-101: User Registration   [Done]  │
│   ○ frontend │  ├── STORY-102: Login Flow         [In Pr] │
│   ○ mobile   │  │   ├── TASK-201: API Endpoint           │
│              │  │   ├── TASK-202: Frontend Form           │
│ ▼ Quick JQL  │  │   └── TASK-203: Validation             │
│   My Open    │  └── STORY-103: Password Reset     [Todo]  │
│   Sprint 23  │                                              │
│   Blocked    │  ┌─────────────────────────────────────┐    │
│              │  │ Title: Add OAuth support             │    │
│ ▼ Templates  │  │ Description: [....................]  │    │
│   Feature    │  │ Story Points: [5]                   │    │
│   Bug Fix    │  │ Sprint: [Sprint 24]                 │    │
│   Refactor   │  │ Labels: [auth, oauth]               │    │
│              │  └─────────────────────────────────────┘    │
├──────────────┴──────────────────────────────────────────────┤
│ [F1]Help [F2]Search [F3]Filter [F4]Create [F5]Sync [Q]uit   │
└──────────────────────────────────────────────────────────────┘
```

### TUI Features

1. **Real-time Sync Indicators**
   - Visual diff between local/remote
   - Conflict highlighting
   - Push/pull status

2. **Interactive Editing**
   - Edit tickets directly in TUI
   - Markdown preview pane
   - Field autocomplete from Jira schema

3. **Bulk Operations**
   - Multi-select with Space
   - Batch status updates
   - Drag-and-drop hierarchy reorganization

4. **Smart Search**
   ```
   / to search
   : for commands
   @ for people
   # for ticket IDs
   ```

5. **Vim Keybindings**
   - `j/k` - Navigate
   - `o` - Open/expand
   - `dd` - Delete
   - `yy` - Yank (copy)
   - `p` - Paste

## Integration Points

### 1. System Tray Integration (Future)
```go
// Ticketr running as background service
type TicketrDaemon struct {
    tray *systray.App
    db   *sql.DB
}

// Right-click menu:
// - Quick Create
// - Recent Tickets
// - Sync Now
// - Open TUI
```

### 2. Shell Integration
```bash
# ~/.zshrc or ~/.bashrc
eval "$(ticketr shell-init zsh)"

# Enables:
# - Automatic workspace detection based on pwd
# - Completions for ticket IDs
# - Quick aliases (tw = ticketr workspace)
```

### 3. Editor Integration
```vim
" Vim plugin
:TicketrSearch PROJ-123
:TicketrCreate
:TicketrPush
```

## The Adapter Layer Architecture

```go
// internal/adapters/
├── cli/          # Traditional CLI (Cobra)
├── tui/          # TUI Interface (tview/bubbletea)
├── api/          # REST API (future)
├── grpc/         # gRPC service (future)
├── filesystem/   # File operations
├── jira/         # Jira API
└── database/     # SQLite/PostgreSQL

// All adapters implement the same ports
type UIPort interface {
    DisplayTickets([]Ticket)
    HandleInput() Command
    ShowError(error)
}

// TUI is just another adapter!
type TUIAdapter struct {
    core *core.TicketService
}

func (t *TUIAdapter) DisplayTickets(tickets []Ticket) {
    // Render in TUI grid
}

// CLI adapter
type CLIAdapter struct {
    core *core.TicketService
}

func (c *CLIAdapter) DisplayTickets(tickets []Ticket) {
    // Print to stdout
}
```

## Multi-Mode Operation

### 1. Global Mode (Default)
```bash
ticketr push ~/documents/sprint-23.md
# Uses current workspace context
# Stores state in ~/.config/ticketr/ticketr.db
```

### 2. Project Mode (Backward Compatible)
```bash
cd ~/project
ticketr push tickets.md --local
# Creates .ticketr/ in current directory
# For teams that want project-specific state
```

### 3. Remote-Only Mode
```bash
ticketr remote pull --jql "project = PROJ" --no-local
# Don't save anything locally
# Pure Jira browser mode
```

## Progressive Enhancement Path

### Phase 1: Centralized State (Q1 2025)
- Move state to `~/.config/ticketr/ticketr.db`
- Add workspace concept
- Keep CLI interface unchanged

### Phase 2: TUI Introduction (Q2 2025)
- Basic TUI with tree view
- Interactive ticket editing
- Real-time sync status

### Phase 3: Global Tool (Q2 2025)
- Install to system PATH
- Remove directory dependency
- Add workspace switching

### Phase 4: Advanced TUI (Q3 2025)
- Split panes
- Markdown preview
- Bulk operations
- Vim bindings

### Phase 5: Service Mode (Q4 2025)
- Background daemon
- System tray
- Webhooks
- Auto-sync

## The Philosophy

**"Ticketr lives with you, not in your project"**

Just like `git` is a global tool that operates on repositories, Ticketr should be a global tool that operates on workspaces. The difference:

- **Git**: Requires being inside a repository
- **Ticketr**: Works from anywhere, on any workspace

This makes Ticketr more like `kubectl` or `docker`:
```bash
kubectl get pods --context production
docker ps --all

# Future Ticketr:
ticketr list --workspace backend
ticketr tui --workspace frontend
```

## Configuration Evolution

### Global Config (`~/.config/ticketr/config.yaml`)
```yaml
default_workspace: backend
editor: vim
tui:
  theme: dracula
  vim_mode: true
  auto_refresh: 30s

workspaces:
  backend:
    url: https://company.atlassian.net
    project: BACK
    default_jql: "assignee = currentUser()"

  frontend:
    url: https://company.atlassian.net
    project: FRONT

  personal:
    url: https://personal.atlassian.net
    project: PERSONAL

aliases:
  my: "assignee = currentUser() AND status != Done"
  sprint: "sprint in openSprints()"
  blocked: "status = Blocked"
```

## Why This Matters

1. **Reduced Friction**: No need to navigate to project directories
2. **Unified View**: See all your work across all projects
3. **Better State Management**: One source of truth in SQLite
4. **Power User Features**: TUI for complex operations
5. **Automation Ready**: CLI remains for scripts
6. **Modern UX**: Like k9s for Kubernetes, we become the "t9s" for Jira

## Success Metrics

- **Installation**: One-liner that adds to PATH globally
- **First Run**: Works immediately without configuration
- **Workspace Switch**: < 100ms
- **TUI Launch**: < 500ms from anywhere
- **Database Size**: < 100MB for 10,000 tickets

## The End Goal

```bash
# The dream workflow
$ ticketr tui  # From anywhere, instant TUI
# OR
$ tk my        # Alias for common queries
# OR
$ ticketr create --template bug "Login button not working"
# OR
$ ticketr daemon start  # Run in background with system tray
```

Ticketr becomes not just a tool you use, but an **ambient workspace** that's always available, always synced, and always ready. It's not about where your files are—it's about where your work is.

---

*"The best tools disappear into your workflow. Ticketr should be everywhere and nowhere—global when you need it, invisible when you don't."*