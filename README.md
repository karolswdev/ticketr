# Ticketr 🎫

A powerful command-line tool that bridges the gap between local Markdown files and Jira, enabling seamless story and task management with bidirectional synchronization.

[![CI](https://github.com/karolswdev/ticketr/actions/workflows/ci.yml/badge.svg)](https://github.com/karolswdev/ticketr/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://go.dev)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/karolswdev/ticketr)](https://pkg.go.dev/github.com/karolswdev/ticketr)
[![Go Report Card](https://goreportcard.com/badge/github.com/karolswdev/ticketr?refresh=1)](https://goreportcard.com/report/github.com/karolswdev/ticketr)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](Dockerfile)

## ✨ Features

- **📝 Markdown-First Workflow**: Define stories and tasks in simple Markdown files
- **🔄 Bidirectional Sync**: Create new items or update existing ones in Jira
- **🎯 Smart Updates**: Automatically detects and updates only changed items
- **🚀 CI/CD Ready**: Built for automation with non-interactive modes
- **🐳 Docker Support**: Lightweight container (~15MB) for consistent execution
- **🔒 Secure**: Environment-based configuration keeps credentials safe
- **📊 Analytics**: Built-in reporting and progress tracking

## 🚀 Quick Start

> **New to Ticketr?** Check out our comprehensive [Getting Started Guide](docs/GETTING_STARTED.md)!

### Installation

#### Using Go
```bash
go install github.com/karolswdev/ticketr/cmd/ticketr@latest
```

#### Building from Source
```bash
git clone https://github.com/karolswdev/ticketr.git
cd ticketr
go build -o ticketr cmd/ticketr/main.go
```

### Configuration

Set up your Jira credentials as environment variables:

```bash
export JIRA_URL="https://yourcompany.atlassian.net"
export JIRA_EMAIL="your.email@company.com"
export JIRA_API_KEY="your-api-token"
export JIRA_PROJECT_KEY="PROJ"
```

💡 **Tip**: Store these in a `.env` file for convenience (see `.env.example`)

### Basic Usage

1. **Create a story file** (`stories.md`):

```markdown
# STORY: User Authentication System

## Description
As a developer, I want to implement a secure authentication system
so that users can safely access the application.

## Acceptance Criteria
- Users can register with email and password
- Passwords are securely hashed
- Session management is implemented

## Tasks
- Set up authentication database schema
- Implement password hashing service
- Create login/logout endpoints
- Add session middleware
```

2. **Sync with Jira**:

```bash
ticketr -f stories.md
```

3. **Result**: Your file is updated with Jira IDs:

```markdown
# STORY: [PROJ-123] User Authentication System

## Description
As a developer, I want to implement a secure authentication system
so that users can safely access the application.

## Acceptance Criteria
- Users can register with email and password
- Passwords are securely hashed
- Session management is implemented

## Tasks
- [PROJ-124] Set up authentication database schema
- [PROJ-125] Implement password hashing service
- [PROJ-126] Create login/logout endpoints
- [PROJ-127] Add session middleware
```

## 📖 Advanced Usage

### Updating Existing Items

Simply edit your file and run the tool again - it intelligently handles updates:

```markdown
# STORY: [PROJ-123] User Authentication System (Updated)

## Tasks
- [PROJ-124] Set up authentication database schema ✅
- [PROJ-125] Implement password hashing service
- Add JWT token generation  # New task will be created
```

### Command-Line Options

```bash
# Push tickets to JIRA (with pre-flight validation)
ticketr push stories.md

# Pull tickets from JIRA to Markdown
ticketr pull --project PROJ --jql "status=Done" -o done_tickets.md

# Pull tickets from a specific epic
ticketr pull --epic PROJ-100 -o epic_tickets.md

# Analyze tickets and display statistics
ticketr stats stories.md

# Verbose output for debugging
ticketr push stories.md --verbose

# Continue on errors (CI/CD mode)
ticketr push stories.md --force-partial-upload

# Dry run - validate and preview changes without modifying JIRA
ticketr push stories.md --dry-run

# Discover JIRA schema and generate configuration
ticketr schema > .ticketr.yaml

# Legacy mode (backward compatibility)
ticketr -f stories.md -v --force-partial-upload
```

### Push Command

The `ticketr push` command synchronizes your local Markdown tickets with JIRA:

**Options:**
- `--dry-run` - Validate tickets and show what would be done without making any changes to JIRA
- `--force-partial-upload` - Continue processing even if some tickets fail (useful for CI/CD)
- `--verbose` or `-v` - Enable detailed logging output

**Note**: Ticketr validates your file for correctness before sending any data to Jira, preventing partial failures. Validation includes:
- Hierarchical rules (e.g., Sub-tasks cannot be children of Epics)
- Required fields validation
- Format validation (only `# TICKET:` format is supported, legacy `# STORY:` format is rejected)

### Pull Command

The `ticketr pull` command fetches tickets from JIRA and intelligently merges them with your local file:

```bash
# Pull all tickets from a project
ticketr pull --project PROJ

# Pull tickets using JQL query
ticketr pull --jql "status IN ('In Progress', 'Done')"

# Pull tickets from a specific epic
ticketr pull --epic PROJ-100 --output sprint_23.md

# Combine filters
ticketr pull --project PROJ --jql "assignee=currentUser()" -o my_tickets.md
```

**Pull Command Options:**
- `--project` - JIRA project key to pull from (uses JIRA_PROJECT_KEY env var if not specified)
- `--epic` - Filter tickets by epic key
- `--jql` - Custom JQL query for filtering
- `-o, --output` - Output file path (default: pulled_tickets.md)

**Conflict Detection:**

The pull command now features intelligent conflict detection:
- **Safe Merge**: Automatically updates tickets that have only changed remotely
- **Conflict Detection**: Identifies when both local and remote versions have changed
- **Local Preservation**: Keeps local changes when only local has been modified
- **State Tracking**: Uses `.ticketr.state` to track both local and remote changes

When conflicts are detected, you can resolve them using strategies:
```
⚠️  Conflict detected! The following tickets have both local and remote changes:
  - TICKET-123
  - TICKET-456

To resolve conflicts, use --strategy flag with one of:
  --strategy=local-wins   Keep local changes
  --strategy=remote-wins  Use remote changes
```

Example usage:
```bash
# Keep local changes when conflicts occur
ticketr pull --project PROJ --strategy=local-wins

# Use remote changes when conflicts occur
ticketr pull --project PROJ --strategy=remote-wins
```

### Analytics and Reporting

The `ticketr stats` command provides detailed analytics about your tickets:

```bash
# Analyze tickets in a file
ticketr stats stories.md

# Example output:
╔══════════════════════════════════════╗
║        TICKET ANALYTICS REPORT       ║
╚══════════════════════════════════════╝

📊 Overall Statistics
────────────────────
  Total Tickets:      10
  Total Tasks:        25
  Total Story Points: 42.5
  Acceptance Criteria: 35

🔄 JIRA Synchronization
────────────────────
  Tickets Synced: 8/10 (80%)
  Tasks Synced:   20/25 (80%)

📋 Tickets by Type
────────────────────
  Story:     ████████░░░░░░░░░░░░ 5
  Bug:       ████░░░░░░░░░░░░░░░░ 2
  Feature:   ████░░░░░░░░░░░░░░░░ 2
  Epic:      ██░░░░░░░░░░░░░░░░░░ 1

📈 Tickets by Status
────────────────────
  Done:        ████████░░░░░░░░░░░░ 4
  In Progress: ██████░░░░░░░░░░░░░░ 3
  To Do:       ██████░░░░░░░░░░░░░░ 3

🎯 Progress Summary
────────────────────
  Overall Completion: 48%
  Items Completed:    12/35
```

The stats command helps you:
- Track overall project progress
- Monitor JIRA synchronization status
- Visualize work distribution by type and status
- Identify bottlenecks and areas needing attention

### Real-time Webhook Synchronization

The `ticketr listen` command starts a webhook server for automatic, real-time synchronization:

```bash
# Start webhook server with default settings
ticketr listen

# Custom port and file
ticketr listen --port 3000 --path project-tickets.md

# With webhook signature validation (recommended)
ticketr listen --secret "your-webhook-secret"
```

This enables instant updates to your local Markdown files whenever tickets change in JIRA. See the [Webhook Configuration Guide](docs/WEBHOOKS.md) for detailed setup instructions including:
- JIRA webhook configuration steps
- Security best practices
- Production deployment options
- Troubleshooting guide

### Schema Discovery

The `ticketr schema` command helps you discover available fields in your JIRA instance and generate a proper configuration file:

```bash
# Discover fields and generate configuration
ticketr schema > .ticketr.yaml

# View available fields with verbose output
ticketr schema -v

# The command will output field mappings like:
# field_mappings:
#   "Story Points":
#     id: "customfield_10010"
#     type: "number"
#   "Sprint": "customfield_10020"
#   "Epic Link": "customfield_10014"
```

This is especially useful when working with custom fields that vary between JIRA instances.

### State Management

Ticketr maintains a `.ticketr.state` file to track content hashes and support intelligent sync flows (e.g., conflict detection on pull):

```bash
# The .ticketr.state file is created/updated as you sync
# It stores SHA256 hashes to detect local vs remote changes

# Pull uses state to detect and resolve conflicts
ticketr pull --strategy=local-wins
```

<<<<<<< HEAD
Push is now state-aware by default. The CLI uses the `PushService` so unchanged tickets are skipped. In `--dry-run` mode, Ticketr shows all intended operations without writing to JIRA or the file/state.
=======
Note: State-aware skipping for `push` exists in the `PushService`, but the default CLI path currently uses `TicketService` (always processes all tickets). A future release will wire `push` to the state-aware flow. Until then, all tickets are processed during `push`.
>>>>>>> origin/main

The `.ticketr.state` file is environment-specific and ignored by default via `.gitignore`.

### Docker Usage

Build and run using Docker:

```bash
# Build the Docker image
docker build -t ticketr .

# Run with Docker
docker run --rm \
  -e JIRA_URL="$JIRA_URL" \
  -e JIRA_EMAIL="$JIRA_EMAIL" \
  -e JIRA_API_KEY="$JIRA_API_KEY" \
  -e JIRA_PROJECT_KEY="$JIRA_PROJECT_KEY" \
  -v $(pwd)/stories.md:/data/stories.md \
  ticketr -f /data/stories.md

# Or use Docker Compose (reads .env automatically)
docker-compose run --rm ticketr
```

## 📋 Story Templates

### Epic Template
```markdown
# STORY: [Epic] Cloud Migration Initiative

## Description
Migrate all services to cloud infrastructure for improved scalability and reliability.

## Acceptance Criteria
- All services running in cloud
- Zero data loss during migration
- Downtime < 1 hour

## Tasks
- Audit current infrastructure
- Design cloud architecture
- Set up cloud environments
- Migrate databases
- Migrate services
- Update DNS and routing
```

### Bug Report Template
```markdown
# STORY: [Bug] Login fails with special characters

## Description
Users cannot login when password contains special characters like & or %.

## Acceptance Criteria
- All special characters work in passwords
- Existing users can still login
- No security vulnerabilities introduced

## Tasks
- Reproduce the issue
- Fix password encoding
- Add comprehensive tests
- Update documentation
```

### Feature Template
```markdown
# STORY: Dark Mode Support

## Description
As a user, I want to switch between light and dark themes
so that I can use the app comfortably in different lighting conditions.

## Acceptance Criteria
- Toggle switch in settings
- Theme preference persisted
- All UI elements properly themed

## Tasks
- Design dark color palette
- Implement theme context
- Update all components
- Add theme toggle to settings
- Test across all pages
```

## 🔄 Workflow Examples

### Sprint Planning Workflow

1. **Create sprint backlog** in Markdown:
```bash
vim sprint-23.md  # Define all stories for the sprint
```

2. **Review with team** (stories still in Markdown)

3. **Push to Jira** when approved:
```bash
ticketr -f sprint-23.md
```

4. **Track progress** by updating the file:
```markdown
## Tasks
- [PROJ-124] Database setup ✅ DONE
- [PROJ-125] API implementation 🚧 IN PROGRESS
- [PROJ-126] Frontend integration
```

### CI/CD Integration

#### Using the Official GitHub Action (Recommended)

```yaml
# .github/workflows/jira-sync.yml
name: Sync Stories to Jira
on:
  push:
    paths:
      - 'stories/*.md'

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Sync to JIRA
        uses: ./.github/actions/ticketr-sync
        with:
          jira-url: ${{ secrets.JIRA_URL }}
          jira-email: ${{ secrets.JIRA_EMAIL }}
          jira-api-key: ${{ secrets.JIRA_API_KEY }}
          jira-project-key: ${{ secrets.JIRA_PROJECT_KEY }}
          command: 'push'
          file-path: 'stories/backlog.md'
          verbose: 'true'
      
      - name: Commit updates
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add stories/backlog.md .ticketr.state || true
          git diff --staged --quiet || git commit -m "Update JIRA IDs [skip ci]"
          git push
```

#### Using Docker (Alternative)

```yaml
jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Sync to Jira
        run: |
          docker run --rm \
            -e JIRA_URL=${{ secrets.JIRA_URL }} \
            -e JIRA_EMAIL=${{ secrets.JIRA_EMAIL }} \
            -e JIRA_API_KEY=${{ secrets.JIRA_API_KEY }} \
            -e JIRA_PROJECT_KEY=${{ secrets.JIRA_PROJECT_KEY }} \
            -v ${{ github.workspace }}:/data \
            ticketr \
            push /data/stories/backlog.md \
            --force-partial-upload
```

## 🏗️ Architecture

Ticketr follows a clean architecture pattern:

```
ticketr/
├── cmd/ticketr/               # CLI entry point
├── internal/
│   ├── core/                  # Business logic
│   │   ├── domain/            # Domain models
│   │   ├── ports/             # Interface definitions
│   │   └── services/          # Core services
│   └── adapters/              # External integrations
│       ├── cli/               # Command-line interface
│       ├── filesystem/        # File I/O operations
│       └── jira/              # Jira API client
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md).

### Development Setup

```bash
# Clone the repository
git clone https://github.com/karolswdev/ticketr.git
cd ticketr

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o ticketr cmd/ticketr/main.go
```

## 📚 Documentation

- **[Getting Started Guide](docs/GETTING_STARTED.md)** - Step-by-step tutorial for new users
- **[Architecture Guide](docs/ARCHITECTURE.md)** - System design and component overview
- **[Development Guide](docs/DEVELOPMENT.md)** - Local setup, testing, and debugging
- **[Webhook Configuration](docs/WEBHOOKS.md)** - Real-time sync with JIRA webhooks
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute to the project
- **[API Documentation](https://pkg.go.dev/github.com/karolswdev/ticketr)** - Go package documentation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- 📖 [Getting Started](docs/GETTING_STARTED.md)
- 🐛 [Issue Tracker](https://github.com/karolswdev/ticketr/issues)
- 💬 [Discussions](https://github.com/karolswdev/ticketr/discussions)
- 🔐 [Security](SECURITY.md)


## 🙏 Acknowledgments

Built with ❤️ using:
- [Go](https://golang.org/) - The programming language
- [Jira REST API](https://developer.atlassian.com/cloud/jira/platform/rest/v2/) - Atlassian's API
- [Alpine Linux](https://alpinelinux.org/) - Container base image

---

**Happy Planning!** 🚀
