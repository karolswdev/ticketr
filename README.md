# Ticketr 🎫

Bridge Markdown and Jira. Create, update, and sync tickets from simple `.md` files — fast, reliable, and CI/CD‑friendly.

For developers who prefer editors and pull requests over tab‑heavy UIs: keep your backlog close to your code, reviewable, and scriptable. Tickets as code — not clicks. ✨

[![CI](https://github.com/karolswdev/ticketr/actions/workflows/ci.yml/badge.svg)](https://github.com/karolswdev/ticketr/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.24%2B-00ADD8?logo=go)](https://go.dev)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/karolswdev/ticketr)](https://pkg.go.dev/github.com/karolswdev/ticketr)
[![Go Report Card](https://goreportcard.com/badge/github.com/karolswdev/ticketr?refresh=1)](https://goreportcard.com/report/github.com/karolswdev/ticketr)
[![Coverage](https://codecov.io/gh/karolswdev/ticketr/branch/main/graph/badge.svg)](https://app.codecov.io/gh/karolswdev/ticketr)
[![Release](https://img.shields.io/github/v/release/karolswdev/ticketr?label=release&logo=github)](https://github.com/karolswdev/ticketr/releases)
[![Downloads](https://img.shields.io/github/downloads/karolswdev/ticketr/total.svg?label=downloads)](https://github.com/karolswdev/ticketr/releases)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Why Ticketr

- **📝 Markdown‑first**: Keep tickets next to code, review with PRs.
- **🔄 Bidirectional sync**: Push to Jira and pull back safely.
- **🧠 State‑aware**: Skips unchanged, detects conflicts, supports strategies.
- **⚙️ CI/CD friendly**: Dry‑run, partial‑upload, clear exit codes.
- **🔔 Portable**: Docker image + webhook server for real‑time sync.

## Install

- Using Go: `go install github.com/karolswdev/ticketr/cmd/ticketr@latest`
- From source: `go build -o ticketr cmd/ticketr/main.go`
- Dockerfile provided (see Getting Started for usage).

## Configure

Set environment variables (or use `.env` — see `.env.example`):

```bash
export JIRA_URL="https://yourcompany.atlassian.net"
export JIRA_EMAIL="your.email@company.com"
export JIRA_API_KEY="your-api-token"
export JIRA_PROJECT_KEY="PROJ"
```

Optional:
- `JIRA_STORY_TYPE` (default: `Task`)
- `JIRA_SUBTASK_TYPE` (default: `Sub-task`)

## 30‑Second Example 🚀

Create `tickets.md`:

```markdown
# TICKET: Implement User Login

## Description
As a user, I want to log into the app so I can access my dashboard.

## Acceptance Criteria
- Invalid credentials show an error
- Successful login redirects to dashboard

## Tasks
- Design login form UI
- Implement authentication API
```

Custom fields example (optional `## Fields` section):

```markdown
## Fields
Story Points: 5
Priority: High
Labels: auth, backend
```

Push to Jira and update the file with IDs:

```bash
ticketr push tickets.md
```

Pull updates from Jira (safe merge, conflict strategies):

```bash
ticketr pull --project PROJ -o pulled.md
# or: ticketr pull --jql "assignee=currentUser()" --strategy=local-wins
```

Quick stats 📊:

```bash
ticketr stats tickets.md
```

## Usage 📌

- `push <file>`: Validate, create/update in Jira, write IDs back.
- `pull [--project|--jql|--epic] [-o file] [--strategy=local-wins|remote-wins]`.
- `schema`: Discover fields; generate `.ticketr.yaml` template.
- `listen [--port 8080] [--path tickets.md] [--secret ...]`: Webhook server.
- `stats <file>`: Simple analytics.

Format: Use `# TICKET:` headings in Markdown files.

## Documentation 📚

- Getting Started: [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md)
- Architecture: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- Webhooks: [docs/WEBHOOKS.md](docs/WEBHOOKS.md)
- Development: [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)
- CI & Coverage: [docs/CI.md](docs/CI.md)
- CLI Reference: [docs/CLI.md](docs/CLI.md)
- Examples: [examples/](examples/)
  - Rich custom fields example: [examples/tickets_with_custom_fields.md](examples/tickets_with_custom_fields.md)
  - Sample state file: [examples/ticketr_state.example.json](examples/ticketr_state.example.json)
- Roadmap: [ROADMAP.md](ROADMAP.md)
- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- Security: [SECURITY.md](SECURITY.md)
- Support: [SUPPORT.md](SUPPORT.md)

## License

MIT — see [LICENSE](LICENSE).

## Community

- Issues: https://github.com/karolswdev/ticketr/issues
- Discussions: https://github.com/karolswdev/ticketr/discussions
