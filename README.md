# Ticketr 🎫

Ticketr keeps your Jira backlog in Git. Author issues in Markdown, review changes like code, and synchronize with Jira whenever you're ready.

[![CI](https://github.com/karolswdev/ticktr/workflows/CI/badge.svg)](https://github.com/karolswdev/ticktr/actions)
[![Coverage](https://img.shields.io/badge/coverage-52.5%25-brightgreen)](https://github.com/karolswdev/ticktr)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](Dockerfile)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> Ticketr understands **only** the `# TICKET:` schema. If you still have `# STORY:` headings, rename them before running the CLI.

## Why Ticketr?

- **Tickets as code** – store Jira issues alongside source, protected by version control, reviews, and history.
- **Bidirectional sync** – push Markdown to Jira, pull Jira updates back to Markdown, resolve conflicts explicitly.
- **Safe automation** – deterministic runs, machine-friendly exit codes, and redacted logs for CI/CD.
- **Human-readable state** – `.ticketr.state` tracks hashes so only changed tickets are touched.
- **Zero lock-in** – plain Markdown files and YAML config keep your backlog portable.

## Quick Start

### 1. Install

```bash
# From source
git clone https://github.com/karolswdev/ticketr.git
cd ticketr

go build -o ticketr ./cmd/ticketr

# or
# go install github.com/karolswdev/ticketr/cmd/ticketr@latest
```

### 2. Configure credentials

Ticketr uses the Atlassian Cloud REST API.

```bash
export JIRA_URL="https://yourcompany.atlassian.net"
export JIRA_EMAIL="you@yourcompany.com"
export JIRA_API_KEY="<api-token>"    # Create at https://id.atlassian.com/manage-profile/security/api-tokens
export JIRA_PROJECT_KEY="PROJ"
```

> Tip: keep these in an `.env` file and `source .env` locally. In CI, store them as secrets.

### 3. Draft your first ticket

```markdown
# TICKET: New Authentication Flow

## Description
Ship the new login, registration, and session handling experience.

## Acceptance Criteria
- Users can sign up with email + password
- Passwords stored with bcrypt
- Sessions expire after 24 hours of inactivity

## Tasks
- Implement login and registration endpoints
- Add password reset flow
- Instrument audit logging
```

Save this as `tickets/auth.md` (or any filename you prefer).

### 4. Push to Jira

```bash
ticketr push tickets/auth.md
```

Ticketr validates your Markdown, creates missing issues, updates existing ones, and injects Jira keys back into the file.

### 5. Pull updates from Jira

```bash
ticketr pull --project PROJ --output tickets/pulled.md
```

Pull merges Jira changes into a Markdown file, highlights conflicts, and respects `.ticketr.state` so only edited tickets are touched.

## Everyday Workflow

```bash
# Author or edit Markdown in git branches
vim tickets/sprint-24.md

# Preview issues without touching Jira
ticketr push tickets/sprint-24.md --force-partial-upload

# Synchronize (create/update) in Jira
ticketr push tickets/sprint-24.md

# Bring Jira edits back to Markdown
ticketr pull --project PROJ --jql "Sprint = 'Sprint 24'" --output tickets/sprint-24.md

# Resolve conflicts if both sides changed
vim tickets/sprint-24.md
```

Automate the same flow from CI:

```yaml
# .github/workflows/jira-sync.yml
jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: go build -o ticketr ./cmd/ticketr
      - run: ./ticketr push backlog.md --force-partial-upload
        env:
          JIRA_URL: ${{ secrets.JIRA_URL }}
          JIRA_EMAIL: ${{ secrets.JIRA_EMAIL }}
          JIRA_API_KEY: ${{ secrets.JIRA_API_KEY }}
          JIRA_PROJECT_KEY: PROJ
```

## Core Concepts

### Markdown schema

Every ticket starts with `# TICKET:` followed by sections (Description, Acceptance Criteria, Tasks, custom `## Fields`, etc.). Tasks are Markdown list items that can hold their own detail blocks.

### Field inheritance

Tasks inherit any custom fields defined on their parent ticket, unless you override them explicitly.

```markdown
# TICKET: [PROJ-100] Payment Gateway

## Fields
- Priority: High
- Sprint: Sprint 24
- Component: API

## Tasks
- ### Build adapter
  #### Fields
  - Priority: Critical      # Overrides parent priority

- ### Document retry policy
  #### Fields
  - Component: Docs         # Overrides parent component
```

Ticketr will treat the second task as `Priority=High, Sprint=Sprint 24, Component=Docs` because only `Component` is overridden. See [docs/WORKFLOW.md](docs/WORKFLOW.md) for a complete breakdown.

### State tracking

Ticketr keeps `.ticketr.state` (ignored by git) with hashes of the last successful push/pull. If you delete the file, the next run treats everything as changed.

### Conflict detection

`ticketr pull` compares the state file, your Markdown, and Jira. When all three diverge, the pull fails with a conflict. Fix the Markdown manually or accept remote changes with `--force`.

### Logging

Each run writes a timestamped log in `.ticketr/logs/` with credentials redacted. The last 10 logs are retained automatically.

## CLI essentials

```bash
# Push one or more files
ticketr push backlog.md

# Merge Jira changes back into Markdown
ticketr pull --project PROJ --output backlog.md

# Force remote version when resolving conflicts
ticketr pull --project PROJ --force

# Continue despite validation errors (records partial successes)
ticketr push backlog.md --force-partial-upload

# Discover Jira fields and generate .ticketr.yaml
ticketr schema > .ticketr.yaml
```

Run `ticketr --help` or `ticketr <command> --help` for full flag descriptions.

## Templates & documentation

- [examples/](examples/) – ready-to-use Markdown templates for epics, bugs, sprints, inheritance patterns
- [docs/WORKFLOW.md](docs/WORKFLOW.md) – end-to-end walkthroughs
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) – ports & adapters design
- [docs/state-management.md](docs/state-management.md) – hash algorithm and conflict model
- [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) – common failures and fixes
- [docs/release-process.md](docs/release-process.md) – release playbook

## Troubleshooting highlights

| Symptom | Quick check |
|---------|-------------|
| `401 Unauthorized` | Ensure `JIRA_URL` includes `https://` and the API token is fresh |
| Missing custom fields | Run `ticketr schema > .ticketr.yaml` and commit the config |
| Nothing pushes | Inspect `.ticketr.state`; delete it to force a full sync |
| Pull conflicts every time | Someone or automation is editing the Markdown + Jira simultaneously – reconcile, then push |

More detail lives in [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md).

## Development

```bash
go test ./...
bash scripts/quality.sh
bash tests/smoke/smoke_test.sh
```

CI runs the same checks on every pull request. See [docs/ci.md](docs/ci.md) for workflow internals.

## Support & contributing

- 📖 Docs & wiki: https://github.com/karolswdev/ticketr/wiki
- 🐞 Issues: https://github.com/karolswdev/ticketr/issues
- 💬 Discussions: https://github.com/karolswdev/ticketr/discussions
- 🔒 Security: [SECURITY.md](SECURITY.md)
- 🤝 Contributions: [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT – see [LICENSE](LICENSE).

---

Built with ❤️ in Go, backed by the Jira Cloud REST API, and ready for teams who treat planning like code.
