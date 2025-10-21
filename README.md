# Ticketr 🎫

Manage JIRA tickets using Markdown files with bidirectional sync. Version control your backlog, automate workflows, and work offline-first.

[![CI](https://github.com/karolswdev/ticktr/workflows/CI/badge.svg)](https://github.com/karolswdev/ticktr/actions)
[![Coverage](https://img.shields.io/badge/coverage-52.5%25-brightgreen)](https://github.com/karolswdev/ticktr)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](Dockerfile)

**Current Version:** v3.1.1 (Production-Ready - Simplified v3 Architecture)

> **Complete Requirements Specification:** See [REQUIREMENTS.md](REQUIREMENTS.md) for the authoritative source of all 51 requirements, acceptance criteria, and traceability.

## Features

- 📝 **Markdown-First**: Define tickets in simple Markdown
- 🔄 **Bidirectional Sync**: Push to and pull from JIRA
- 🎯 **Smart Updates**: Only syncs changed tickets
- 🚀 **CI/CD Ready**: Non-interactive modes for automation
- 🐳 **Docker Support**: Lightweight 15MB container
- 🔒 **Secure**: OS keychain credential storage
- 🏢 **Multi-Workspace**: Manage multiple Jira projects seamlessly
- 👥 **Credential Profiles**: Reusable credentials across workspaces
- 🎨 **TUI Interface**: Full-featured terminal interface with workspace creation
- 📁 **XDG-Compliant**: Platform-standard file locations
- ⚡ **Bulk Operations**: Update, move, or delete multiple tickets at once with real-time progress
- 🔀 **Smart Sync Strategies**: Choose how conflicts are resolved during sync
  - LocalWins: Preserve your local changes
  - RemoteWins: Always accept Jira updates (default)
  - ThreeWayMerge: Intelligent field-level merging
- 🔖 **JQL Aliases**: Create reusable named queries for faster ticket filtering

## Quick Start

### Installation

```bash
# Using Go
go install github.com/karolswdev/ticketr/cmd/ticketr@latest

# Or build from source
git clone https://github.com/karolswdev/ticketr.git
cd ticketr && go build -o ticketr cmd/ticketr/main.go
```

### Configuration

```bash
export JIRA_URL="https://yourcompany.atlassian.net"
export JIRA_EMAIL="your.email@company.com"
export JIRA_API_KEY="your-api-token"        # Get from: id.atlassian.com/manage-profile/security/api-tokens
export JIRA_PROJECT_KEY="PROJ"
```

💡 Store in `.env` file (see `.env.example`)

## Workspace Management

Ticketr v3.0 supports managing multiple Jira projects from a single installation using workspaces.

### Credential Profiles (v3.0+)

For teams or individuals managing multiple projects, credential profiles allow you to reuse Jira credentials across workspaces:

```bash
# Create a reusable credential profile
ticketr credentials profile create company-admin \
  --url https://company.atlassian.net \
  --username admin@company.com \
  --token your-api-token

# List available profiles
ticketr credentials profile list

# Create workspaces using the profile
ticketr workspace create backend --profile company-admin --project BACK
ticketr workspace create frontend --profile company-admin --project FRONT
ticketr workspace create devops --profile company-admin --project OPS
```

### Creating a Workspace (Direct Method)

```bash
ticketr workspace create backend \
  --url https://company.atlassian.net \
  --project BACK \
  --username your.email@company.com \
  --token your-api-token
```

### Credential Management and TUI Support

Credentials are stored securely in your OS keychain:
- **macOS:** Keychain Access
- **Windows:** Credential Manager
- **Linux:** GNOME Keyring / KWallet

**TUI Workspace Creation**: Use `ticketr tui` and press `w` to create workspaces interactively with credential profile support.

### Listing Workspaces

```bash
ticketr workspace list
```

### Switching Workspaces

```bash
ticketr workspace switch frontend
```

### Managing Default Workspace

```bash
ticketr workspace set-default backend
ticketr workspace current
```

### Deleting a Workspace

```bash
ticketr workspace delete old-project
# Or force delete without confirmation
ticketr workspace delete old-project --force
```

### Security

- Credentials are never stored in the database
- All credentials use OS-level encryption
- Credentials are never logged or printed
- Each workspace has isolated credentials

See [docs/workspace-management-guide.md](docs/workspace-management-guide.md) for comprehensive workspace documentation.

## Bulk Operations

Perform operations on multiple tickets at once with real-time progress feedback:

### Bulk Update

Update multiple tickets with field changes:

```bash
# Update status for multiple tickets
ticketr bulk update --ids PROJ-1,PROJ-2,PROJ-3 --set status=Done

# Update multiple fields
ticketr bulk update --ids PROJ-1,PROJ-2 --set status="In Progress" --set assignee=john@example.com

# Update with spaces in values
ticketr bulk update --ids PROJ-1,PROJ-2 --set priority="High Priority"
```

### Bulk Move

Move multiple tickets to a new parent:

```bash
# Move tickets to a new parent
ticketr bulk move --ids PROJ-1,PROJ-2,PROJ-3 --parent PROJ-100

# Move sub-tasks to a different epic
ticketr bulk move --ids TASK-1,TASK-2 --parent EPIC-42
```

### Bulk Delete

Delete multiple tickets:

```bash
# Delete tickets with confirmation
ticketr bulk delete --ids PROJ-1,PROJ-2 --confirm
```

### Features

- **Real-time progress**: See [X/Y] progress as tickets are processed
- **JQL injection prevention**: Ticket IDs validated to prevent malicious input
- **Best-effort rollback**: Attempts to restore original state on partial failures
- **Maximum 100 tickets**: Per-operation limit for safety and performance
- **Workspace-scoped**: Requires active workspace (use `ticketr workspace current`)

See [docs/bulk-operations-guide.md](docs/bulk-operations-guide.md) for detailed usage and examples.

#### TUI Multi-Select and Bulk Operations

The TUI supports interactive multi-select and bulk operations:

**Multi-Select Tickets**:
```bash
# In TUI:
# 1. Navigate to ticket tree (Tab)
# 2. Press Space to select/deselect current ticket (shows [x] checkbox)
# 3. Press 'a' to select all visible tickets
# 4. Press 'A' (Shift+a) to deselect all
# 5. Selected count shows in title: "Tickets (3 selected)"
# 6. Border turns teal/blue when tickets are selected
```

**Execute Bulk Operations**:
```bash
# With tickets selected:
# 1. Press 'b' to open bulk operations menu
# 2. Choose operation:
#    - Update Fields: Change Status, Priority, Assignee, Custom Fields
#    - Move Tickets: Move all selected tickets under a new parent
#    - Delete Tickets: Warning - not yet supported in v3.0
# 3. Fill in form fields
# 4. Click Apply to execute
# 5. Watch real-time progress with success/failure indicators
# 6. Press Cancel during operation to stop (partial changes applied)
```

**TUI Features**:
- Real-time progress tracking during operations
- Green checkmark for successful tickets, red X for failures
- Context cancellation (Esc or Cancel button stops operation)
- Automatic rollback on partial failure (best-effort)
- Clear error messages and validation
- Selection state persists across navigation

**Help in TUI**:
Press `?` in the TUI to see all bulk operations keybindings and detailed usage instructions.

## File Locations

Ticketr follows platform-standard directory conventions for storing configuration, data, and cache files.

### Linux / Unix (XDG-Compliant)

```
~/.config/ticketr/          # Configuration files
~/.local/share/ticketr/     # Database and persistent data
  ├── ticketr.db            # SQLite workspace database
  ├── state.json            # State tracking
  └── backups/              # Migration backups
~/.cache/ticketr/           # Logs and temporary files
```

### macOS

```
~/Library/Application Support/ticketr/  # Database and data
~/Library/Preferences/ticketr/          # Configuration
~/Library/Caches/ticketr/               # Logs and cache
```

### Windows

```
%LOCALAPPDATA%\ticketr\     # Database and data
%APPDATA%\ticketr\          # Configuration
%TEMP%\ticketr\             # Logs and cache
```

### Migration from Legacy Versions

Users upgrading from v2.x can use the migration commands. See archived migration guides in `docs/archive/` for details:

**Manual migration**:
```bash
ticketr migrate-paths        # Migrate to global paths (v2.x → v3.x)
ticketr rollback-paths       # Rollback to v2.x local paths
ticketr migrate <file>       # Migrate STORY format to TICKET format
```

**Note (v3.1.1+)**: As of v3.1.1, all v3 features are enabled by default. The `ticketr v3 enable` and `ticketr v3 migrate` commands no longer exist. Simply install v3.1.1 and all features work immediately.

### Basic Usage

**1. Create a ticket file:**

```markdown
# TICKET: User Authentication System

## Description
Implement secure JWT-based authentication.

## Acceptance Criteria
- Users can register with email/password
- Passwords securely hashed with bcrypt
- JWT tokens expire after 24 hours

## Tasks
- Set up authentication database schema
- Implement password hashing service
- Create login/logout endpoints
- Add JWT validation middleware
```

**2. Push to JIRA:**

```bash
ticketr push tickets.md
```

**3. File updated with JIRA IDs:**

```markdown
# TICKET: [PROJ-123] User Authentication System
...
## Tasks
- [PROJ-124] Set up authentication database schema
- [PROJ-125] Implement password hashing service
...
```

## Common Commands

### Workspace Commands

| Command | Description |
|---------|-------------|
| `workspace create` | Create a new workspace with Jira credentials |
| `workspace create --profile` | Create a workspace using an existing credential profile |
| `workspace list` | List all configured workspaces |
| `workspace switch` | Switch to a different workspace |
| `workspace current` | Show the current active workspace |
| `workspace delete` | Delete a workspace and its credentials |
| `workspace set-default` | Set the default workspace |
| `credentials profile create` | Create a reusable credential profile |
| `credentials profile list` | List available credential profiles |
| `migrate-paths` | Migrate from v2.x local paths to v3.x global paths |
| `rollback-paths` | Rollback migration to v2.x local paths |

### Ticket Management Commands

```bash
# Push tickets to JIRA
ticketr push tickets.md

# Pull from JIRA
ticketr pull --project PROJ --output tickets.md

# Force remote changes on conflict
ticketr pull --project PROJ --force

# Continue on validation errors (CI/CD mode)
ticketr push tickets.md --force-partial-upload

# Discover JIRA schema/fields
ticketr schema > .ticketr.yaml

# Migrate legacy format
ticketr migrate old-tickets.md --write
```

## Key Concepts

### State Management

Ticketr tracks changes via global state file (gitignored). Only modified tickets sync to JIRA.

**State File Location**:
- **Linux/Unix**: `~/.local/share/ticketr/state.json`
- **macOS**: `~/Library/Application Support/ticketr/state.json`
- **Windows**: `%LOCALAPPDATA%\ticketr\state.json`

```bash
# Force full re-push (v3.0)
rm ~/.local/share/ticketr/state.json && ticketr push tickets.md

# Legacy v2.x path (if not migrated)
rm .ticketr.state && ticketr push tickets.md
```

See [docs/state-management.md](docs/state-management.md) for details.

### Field Inheritance

Tasks automatically inherit parent custom fields. Task-specific fields override.

```markdown
# TICKET: [PROJ-100] Payment Integration

## Fields
- Priority: High
- Sprint: Sprint 24

## Tasks
- ### Setup payment gateway
  #### Fields
  - Priority: Critical    # Overrides High
  # Inherits: Sprint 24
```

See [docs/WORKFLOW.md](docs/WORKFLOW.md) for comprehensive examples.

### Conflict Detection & Smart Sync Strategies

Ticketr provides three sync strategies to handle conflicts when pulling from Jira:

- **LocalWinsStrategy**: Keeps local changes, ignores remote updates
- **RemoteWinsStrategy**: Accepts remote changes, discards local edits (default)
- **ThreeWayMergeStrategy**: Merges compatible changes, errors on conflicts

**Example - Compatible Changes Auto-Merge**:
```
Local:  Title="Fix bug", Description="Updated locally"
Remote: Title="Fix bug", Status="In Progress" (updated in Jira)

Result with ThreeWayMerge: Both changes preserved (different fields)
```

**Example - Incompatible Changes Error**:
```
Local:  Title="Fix authentication bug"
Remote: Title="Auth system improvements"

Result with ThreeWayMerge: Error - Title field has conflicting changes
```

**Current Behavior**:
- Default strategy: RemoteWins (backward compatible)
- CLI flag and config file support available in v3.1.1

**Manual Conflict Resolution**:
```bash
# Accept remote changes (default)
ticketr pull --project PROJ --force

# Manually merge conflicts
vim tickets.md  # Edit to desired state
ticketr push tickets.md
```

For detailed guidance, see the [Sync Strategies Guide](docs/sync-strategies-guide.md).

### Logging

All operations logged to platform-standard cache directory with sensitive data redacted. Logs auto-rotate (keeps last 10).

**Log Location**:
- **Linux/Unix**: `~/.cache/ticketr/logs/`
- **macOS**: `~/Library/Caches/ticketr/logs/`
- **Windows**: `%TEMP%\ticketr\logs\`

## Advanced Usage

### Pull with Filters

```bash
# Pull from specific epic
ticketr pull --epic PROJ-100 -o sprint23.md

# Use JQL query
ticketr pull --jql "status IN ('In Progress', 'Done')" -o active.md

# Combine filters
ticketr pull --project PROJ --jql "assignee=currentUser()"
```

### JQL Aliases

Create reusable named queries for easier ticket filtering. Aliases can be workspace-specific or global, and support recursive references.

**Predefined Aliases** (available by default):
- `mine`: Tickets assigned to you
- `sprint`: Tickets in active sprints
- `blocked`: Blocked tickets or tickets with blocked label

**Basic Usage**:

```bash
# List all available aliases
ticketr alias list

# Create workspace-specific alias
ticketr alias create my-bugs "assignee = currentUser() AND type = Bug"

# Create global alias (available in all workspaces)
ticketr alias create critical "priority = Critical" --global

# Show alias details
ticketr alias show my-bugs

# Update an alias
ticketr alias update my-bugs "assignee = currentUser() AND type = Bug AND status != Done"

# Delete an alias
ticketr alias delete my-bugs
```

**Using Aliases with Pull**:

```bash
# Pull using predefined alias
ticketr pull --alias mine --output my-tickets.md

# Pull using custom alias
ticketr pull --alias my-bugs --output bugs.md

# Combine with other filters
ticketr pull --alias sprint --epic PROJ-100
```

**Advanced: Recursive Aliases**:

Create aliases that reference other aliases using the `@` syntax:

```bash
# Create base alias
ticketr alias create my-work "assignee = currentUser() AND resolution = Unresolved"

# Reference it in another alias
ticketr alias create urgent-work "@my-work AND priority = High"

# Chain multiple references
ticketr alias create critical-sprint "@urgent-work AND sprint in openSprints()"
```

See [JQL Aliases Guide](docs/FEATURES/JQL-ALIASES.md) for comprehensive documentation.

### Custom Fields

```bash
# Discover available fields
ticketr schema

# Generate config
ticketr schema > .ticketr.yaml
```

Then use fields in Markdown:

```markdown
## Fields
- Story Points: 8
- Sprint: Sprint 23
- Labels: backend, auth
- Component: API
```

### CI/CD Integration

```yaml
# .github/workflows/jira-sync.yml
- name: Sync to JIRA
  run: |
    ticketr push backlog.md --force-partial-upload
  env:
    JIRA_URL: ${{ secrets.JIRA_URL }}
    JIRA_EMAIL: ${{ secrets.JIRA_EMAIL }}
    JIRA_API_KEY: ${{ secrets.JIRA_API_KEY }}
    JIRA_PROJECT_KEY: PROJ
```

### Docker Usage

```bash
docker run --rm \
  --env-file .env \
  -v $(pwd):/workspace \
  ticketr push /workspace/tickets.md
```

Or use Docker Compose (see `docker-compose.yml`).

## Templates

See [examples/](examples/) directory for:
- Epic template
- Bug report template
- Sprint planning template
- Field inheritance examples

## Documentation

### Core Documentation
- [REQUIREMENTS.md](REQUIREMENTS.md) - Complete requirements specification (authoritative)
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - Technical architecture
- [WORKFLOW.md](docs/WORKFLOW.md) - End-to-end usage guide
- [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) - Common issues

### Feature Guides
- [workspace-management-guide.md](docs/workspace-management-guide.md) - Multi-workspace guide
- [bulk-operations-guide.md](docs/bulk-operations-guide.md) - Bulk operations guide
- [sync-strategies-guide.md](docs/sync-strategies-guide.md) - Smart sync strategies guide
- [JQL-ALIASES.md](docs/FEATURES/JQL-ALIASES.md) - JQL aliases and reusable queries
- [state-management.md](docs/state-management.md) - Change detection

### Development
- [release-process.md](docs/release-process.md) - Release management
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guide

### Archive
- [docs/archive/](docs/archive/) - Legacy migration guides for v1.x and v2.x users

## Architecture

Hexagonal (Ports & Adapters) pattern:

```
cmd/ticketr/          # CLI entry point
internal/
├── core/             # Business logic
│   ├── domain/       # Domain models
│   ├── ports/        # Interface definitions
│   └── services/     # Core services
└── adapters/         # External integrations
    ├── filesystem/   # File I/O
    └── jira/         # JIRA API client
```

See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for comprehensive details.

## Development

```bash
# Run tests
go test ./...

# Quality checks
bash scripts/quality.sh

# Smoke tests
bash tests/smoke/smoke_test.sh

# Coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -1
```

**CI/CD:** Automated checks run on all PRs. See [docs/ci.md](docs/ci.md).

## Troubleshooting

### Authentication Issues

```bash
# Test credentials
ticketr schema --verbose

# Common fixes:
export JIRA_URL="https://company.atlassian.net"  # Include https://
# Regenerate API token if expired
```

### Field Not Found

```bash
# Discover exact field names (case-sensitive!)
ticketr schema

# Check available fields
ticketr schema > fields.yaml && cat fields.yaml
```

### No Changes Detected

```bash
# Reset state to force push
rm .ticketr.state
ticketr push tickets.md
```

### Conflict on Pull

```bash
# Accept remote changes
ticketr pull --project PROJ --force

# Or manually merge and push local
vim tickets.md && ticketr push tickets.md
```

For comprehensive troubleshooting, see [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md).

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success (or partial success with `--force-partial-upload`) |
| 1 | Validation failure (without `--force-partial-upload`) |
| 2 | Runtime error (JIRA API, network issues) |

## Support

- 📖 [Documentation](https://github.com/karolswdev/ticketr/wiki)
- 🐛 [Issues](https://github.com/karolswdev/ticketr/issues)
- 💬 [Discussions](https://github.com/karolswdev/ticketr/discussions)
- 🔒 [Security Policy](SECURITY.md)
- 🆘 [Support Guide](SUPPORT.md)

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT License - see [LICENSE](LICENSE) file.

## Acknowledgments

Built with [Go](https://golang.org/), [JIRA REST API](https://developer.atlassian.com/cloud/jira/platform/rest/v2/), and [Alpine Linux](https://alpinelinux.org/).

---

**Happy Planning!** 🚀

For detailed documentation, visit the [Wiki](https://github.com/karolswdev/ticketr/wiki) or browse [docs/](docs/).
