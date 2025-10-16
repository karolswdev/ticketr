# Ticketr 🎫

A powerful command-line tool that bridges the gap between local Markdown files and Jira, enabling seamless story and task management with bidirectional synchronization.

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](Dockerfile)

> **Breaking Change in v2.0:** The `# STORY:` schema has been deprecated in favor of the generic `# TICKET:` schema. See REQUIREMENTS-v2.md (PROD-201) for migration guidance.

## ✨ Features

- **📝 Markdown-First Workflow**: Define stories and tasks in simple Markdown files
- **🔄 Bidirectional Sync**: Create new items or update existing ones in Jira
- **🎯 Smart Updates**: Automatically detects and updates only changed items
- **🚀 CI/CD Ready**: Built for automation with non-interactive modes
- **🐳 Docker Support**: Lightweight container (~15MB) for consistent execution
- **🔒 Secure**: Environment-based configuration keeps credentials safe

## 🚀 Quick Start

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
# TICKET: User Authentication System

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
# TICKET: [PROJ-123] User Authentication System

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

4. **Continue on validation errors** (optional):

If you encounter validation errors but want to proceed anyway:

```bash
ticketr push stories.md --force-partial-upload
```

Note: Valid tickets will be created; invalid ones will be skipped with error messages.

## 📖 Advanced Usage

### Migrating from v1.x

If you have existing `.md` files using the legacy `# STORY:` format, you'll need to update them to the canonical `# TICKET:` format.

**Quick Migration:**
```bash
# Preview changes
ticketr migrate your-file.md

# Apply changes
ticketr migrate your-file.md --write

# Batch migration
ticketr migrate tickets/*.md --write
```

**For detailed migration instructions, see:**
- [Migration Guide](docs/migration-guide.md) - Complete step-by-step guide
- [REQUIREMENTS-v2.md](REQUIREMENTS-v2.md) - PROD-201 schema specification
- CLI help: `ticketr help migrate`

**Why migrate?** The generic `# TICKET:` schema supports all Jira issue types (stories, bugs, tasks, epics) while maintaining hierarchical validation. The parser will reject legacy `# STORY:` format with helpful error messages.

### Updating Existing Items

Simply edit your file and run the tool again - it intelligently handles updates:

```markdown
# TICKET: [PROJ-123] User Authentication System (Updated)

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

# Verbose output for debugging
ticketr push stories.md --verbose

# Continue on errors (CI/CD mode) - accepts partial success
ticketr push stories.md --force-partial-upload

# Discover JIRA schema and generate configuration
ticketr schema > .ticketr.yaml

# Legacy mode (backward compatibility)
ticketr -f stories.md -v --force-partial-upload
```

### Push Command

**Note**: Ticketr validates your file for correctness before sending any data to Jira. Validation includes:
- Hierarchical rules (e.g., Sub-tasks cannot be children of Epics)
- Required fields validation
- Format validation (only `# TICKET:` format is supported)

By default, validation errors prevent pushing to JIRA (fail-fast behavior). Use `--force-partial-upload` to override this and process valid items even when validation errors exist. See [Understanding --force-partial-upload](#understanding---force-partial-upload) for details.

### Understanding --force-partial-upload

The `--force-partial-upload` flag modifies how Ticketr handles validation and runtime errors:

**Pre-Flight Validation Behavior:**
- **Without flag**: Validation errors (e.g., hierarchy violations, missing required fields) cause immediate exit with code 1
- **With flag**: Validation errors are downgraded to warnings, and processing continues

**Runtime Error Behavior:**
- **Without flag**: If any ticket/task fails to create/update in JIRA, exit code 2 is returned
- **With flag**: Processing continues for all tickets, valid items succeed, exit code 0 (partial success accepted)

**Use Cases:**
- **CI/CD Pipelines**: Ensure automated jobs complete even if some tickets fail
- **Bulk Operations**: Process large files where some items may have transient JIRA errors
- **Development/Testing**: Continue workflow despite validation issues during rapid iteration

**Example Output:**

Without flag (validation error):
```bash
$ ticketr push stories.md
❌ Validation errors found:
  - Task 'Invalid Task': A 'Story' cannot be the child of an 'Epic'

2 validation error(s) found. Fix these issues before pushing to JIRA.
Tip: Use --force-partial-upload to continue despite validation errors.
```

With flag (validation warning):
```bash
$ ticketr push stories.md --force-partial-upload
⚠️  Validation warnings (processing will continue with --force-partial-upload):
  - Task 'Invalid Task': A 'Story' cannot be the child of an 'Epic'

2 validation warning(s) found. Some items may fail during upload.

=== Summary ===
Tickets created: 3
Tasks created: 5

=== Errors (1) ===
  - Failed to create task 'Invalid Task': Invalid issue type hierarchy

Processing complete!
```

**Exit Code Summary:**
- Exit 0: Success (or partial success with --force-partial-upload)
- Exit 1: Pre-flight validation failure (without flag)
- Exit 2: Runtime errors (without flag)

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
- `--force` - Force overwrite local changes with remote changes when conflicts are detected

**Conflict Detection:**

The pull command now features intelligent conflict detection:
- **Safe Merge**: Automatically updates tickets that have only changed remotely
- **Conflict Detection**: Identifies when both local and remote versions have changed
- **Local Preservation**: Keeps local changes when only local has been modified
- **State Tracking**: Uses `.ticketr.state` to track both local and remote changes

When conflicts are detected, you'll see:
```
⚠️  Conflict detected! The following tickets have both local and remote changes:
  - TICKET-123
  - TICKET-456

To force overwrite local changes with remote changes, use --force flag
```

**Resolving Conflicts:**

When you're ready to accept remote changes, use the `--force` flag:
```bash
# Force overwrite local changes with remote
ticketr pull --project PROJ --force

# Force with specific filters
ticketr pull --epic PROJ-100 --force -o sprint_23.md
```

**Warning**: Using `--force` will permanently overwrite your local changes with the remote version from JIRA. Make sure to backup or commit your local changes before forcing.

### First Pull

When running `ticketr pull` for the first time (no local file exists), Ticketr will:
- Create a new file with all tickets from JIRA
- Initialize the `.ticketr.state` file to track future changes
- Set up conflict detection for subsequent pulls

**Example:**
```bash
# First pull - creates new file
ticketr pull --project PROJ -o my_tickets.md

# Result: my_tickets.md created with all PROJ tickets
# .ticketr.state initialized
```

**Troubleshooting First Pull:**

If you encounter issues on first pull:

1. **"failed to load local tickets" error**: This should not occur on first run. If it does, ensure:
   - You have write permissions to the output directory
   - The output path is valid
   - Report as a bug if the issue persists

2. **No tickets pulled**: Verify your JIRA connection:
   ```bash
   # Test authentication with verbose output
   ticketr pull --project PROJ -o test.md --verbose
   ```

3. **State file conflicts**: On first pull, delete any existing `.ticketr.state`:
   ```bash
   rm .ticketr.state
   ticketr pull --project PROJ -o my_tickets.md
   ```

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

Ticketr automatically tracks changes to prevent redundant updates to JIRA:

```bash
# The .ticketr.state file is created automatically
# It stores SHA256 hashes of ticket content

# Only changed tickets are pushed to JIRA
ticketr push stories.md  # Skips unchanged tickets

# The state file contains:
# - Ticket ID to content hash mappings
# - Both local and remote hash values
# - Automatically updated after each successful push/pull
```

**Hash Calculation:**

Ticketr uses deterministic SHA256 hashing (Milestone 4) to reliably detect changes:
- Custom field keys are sorted alphabetically before hashing
- Ensures identical content always produces identical hashes
- Prevents false positives from Go's non-deterministic map iteration

**Note**: The `.ticketr.state` file should be added to `.gitignore` as it's environment-specific.

**For detailed state management documentation, see:** [docs/state-management.md](docs/state-management.md)

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

## 📋 Ticket Templates

### Epic Template
```markdown
# TICKET: [Epic] Cloud Migration Initiative

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
# TICKET: [Bug] Login fails with special characters

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
# TICKET: Dark Mode Support

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
      - uses: actions/checkout@v2
      
      - name: Sync to Jira
        run: |
          docker run --rm \
            -e JIRA_URL=${{ secrets.JIRA_URL }} \
            -e JIRA_EMAIL=${{ secrets.JIRA_EMAIL }} \
            -e JIRA_API_KEY=${{ secrets.JIRA_API_KEY }} \
            -e JIRA_PROJECT_KEY=${{ secrets.JIRA_PROJECT_KEY }} \
            -v ${{ github.workspace }}:/data \
            ticketr \
            -f /data/stories/backlog.md \
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

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

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

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- 📖 [Documentation](https://github.com/karolswdev/ticketr/wiki)
- 🐛 [Issue Tracker](https://github.com/karolswdev/ticketr/issues)
- 💬 [Discussions](https://github.com/karolswdev/ticketr/discussions)

## 🙏 Acknowledgments

Built with ❤️ using:
- [Go](https://golang.org/) - The programming language
- [Jira REST API](https://developer.atlassian.com/cloud/jira/platform/rest/v2/) - Atlassian's API
- [Alpine Linux](https://alpinelinux.org/) - Container base image

---

**Happy Planning!** 🚀