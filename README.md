# Ticketr 🎫

Manage JIRA tickets using Markdown files with bidirectional sync. Version control your backlog, automate workflows, and work offline-first.

[![CI](https://github.com/karolswdev/ticktr/workflows/CI/badge.svg)](https://github.com/karolswdev/ticktr/actions)
[![Coverage](https://img.shields.io/badge/coverage-52.5%25-brightgreen)](https://github.com/karolswdev/ticktr)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](Dockerfile)

> **Note:** Ticketr requires `# TICKET:` headings. Files using `# STORY:` must be updated before use.

## Features

- 📝 **Markdown-First**: Define tickets in simple Markdown
- 🔄 **Bidirectional Sync**: Push to and pull from JIRA
- 🎯 **Smart Updates**: Only syncs changed tickets
- 🚀 **CI/CD Ready**: Non-interactive modes for automation
- 🐳 **Docker Support**: Lightweight 15MB container
- 🔒 **Secure**: Environment-based credential management

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

```

## Key Concepts

### State Management

Ticketr tracks changes via `.ticketr.state` (gitignored). Only modified tickets sync to JIRA.

```bash
# Force full re-push
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

### Conflict Detection

Pull command detects simultaneous local/remote changes. Use `--force` to accept remote, or manually merge.

### Logging

All operations logged to `.ticketr/logs/` with sensitive data redacted. Logs auto-rotate (keeps last 10).

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

- [WORKFLOW.md](docs/WORKFLOW.md) - End-to-end usage guide
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - Technical architecture
- [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) - Common issues
- [state-management.md](docs/state-management.md) - Change detection
- [release-process.md](docs/release-process.md) - Release management
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guide

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
