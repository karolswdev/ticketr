# Ticketr CLI Reference

Concise reference for commands, flags, and common examples.

## Global

- `-v, --verbose`: Enable detailed logs
- `--config <file>`: Use an alternate `.ticketr.yaml`

Environment (typical): `JIRA_URL`, `JIRA_EMAIL`, `JIRA_API_KEY`, `JIRA_PROJECT_KEY`, optional `JIRA_STORY_TYPE`, `JIRA_SUBTASK_TYPE`.

## push

Validate Markdown, create/update in Jira, and write IDs back to the file.

Usage:
```
ticketr push <file>
```

Flags:
- `--dry-run`: Validate and print intended actions (no writes)
- `--force-partial-upload`: Continue after errors (good for CI)

Example:
```
ticketr push tickets.md -v --dry-run
```

## pull

Fetch from Jira and merge into Markdown with conflict strategies.

Usage:
```
ticketr pull [flags]
```

Flags:
- `--project <KEY>`: Project key (falls back to `JIRA_PROJECT_KEY`)
- `--jql <query>`: Extra JQL to filter
- `--epic <KEY>`: Filter by epic
- `-o, --output <file>`: Output file (default `pulled_tickets.md`)
- `--strategy <mode>`: `local-wins` | `remote-wins`

Examples:
```
ticketr pull --project PROJ -o backlog.md
ticketr pull --jql "assignee=currentUser() AND sprint in openSprints()" --strategy=local-wins
```

## stats

Analyze a Markdown file and print a simple report.

Usage:
```
ticketr stats <file>
```

Example:
```
ticketr stats tickets.md
```

## schema

Discover available fields and output a `.ticketr.yaml` template.

Usage:
```
ticketr schema > .ticketr.yaml
```

Example:
```
ticketr schema -v
```

## listen

Start a lightweight webhook server for near‑real‑time sync.

Usage:
```
ticketr listen [flags]
```

Flags:
- `--port <n>`: Port to listen on (default `8080`)
- `--path <file>`: Markdown file to update
- `--secret <value>`: Webhook signature secret (recommended)

Examples:
```
ticketr listen
ticketr listen --port 3000 --path project.md --secret "$WEBHOOK_SECRET"
```

## Markdown Format

- Use `# TICKET: <title>` (IDs auto‑inserted as `# TICKET: [PROJ-123] <title>`)
- Optional sections: `## Description`, `## Fields`, `## Acceptance Criteria`, `## Tasks`
- Tasks use `- <title>` or `- [PROJ-456] <title>` with optional nested sections

## State and Conflicts

Ticketr maintains a JSON state file (default: `.ticketr.state`) mapping Jira IDs to local/remote hashes to detect changes and conflicts.

Example: [examples/ticketr_state.example.json](../examples/ticketr_state.example.json)

Conflict strategies when pulling:

```
ticketr pull --project PROJ --strategy=local-wins    # keep local changes
ticketr pull --project PROJ --strategy=remote-wins   # take remote changes
```

See a richer Markdown example with custom fields: [examples/tickets_with_custom_fields.md](../examples/tickets_with_custom_fields.md)
