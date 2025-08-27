# Ticketr 🎫

A powerful command-line tool that bridges the gap between local Markdown files and Jira, enabling seamless story and task management with bidirectional synchronization.

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](Dockerfile)

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
# Basic operation
ticketr -f stories.md

# Verbose output for debugging
ticketr -f stories.md --verbose

# Continue on errors (CI/CD mode)
ticketr -f stories.md --force-partial-upload

# Combine options
ticketr -f stories.md -v --force-partial-upload
```

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