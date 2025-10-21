# Workspace Management Guide

**Version:** 1.0
**Last Updated:** October 17, 2025
**Scope:** Ticketr v3.0 workspace functionality

## Purpose

This guide explains how to use Ticketr's workspace feature to manage multiple Jira projects from a single installation. Workspaces enable you to switch between different Jira instances or projects without reconfiguring credentials or managing multiple installations.

---

## What Are Workspaces?

A **workspace** in Ticketr represents a connection to a specific Jira project. Each workspace maintains:

- Independent project configuration (Jira URL, project key)
- Separate credentials stored securely in your OS keychain
- Isolated ticket data and sync state
- Individual sync history and logs

**Benefits:**
- Manage multiple Jira projects from one Ticketr installation
- Switch contexts quickly without reconfiguring environment variables
- Keep project data isolated and organized
- Secure credential management per workspace

---

## Getting Started

### Prerequisites

- Ticketr v3.0 or later installed
- Jira API credentials for each project you want to manage
- OS keychain access (macOS Keychain, Windows Credential Manager, Linux Secret Service)

### Your First Workspace

When you first use Ticketr v3.0, you'll need to create a workspace:

```bash
ticketr workspace create backend \
  --url https://company.atlassian.net \
  --project BACK \
  --username your.email@company.com \
  --token your-api-token
```

**What happens:**
1. Workspace "backend" is created
2. Credentials are stored securely in your OS keychain
3. The workspace is automatically set as your default
4. You can now start working with tickets

**Note:** The first workspace you create automatically becomes the default workspace.

---

## Creating Workspaces

### Basic Creation

```bash
ticketr workspace create <name> \
  --url <jira-url> \
  --project <project-key> \
  --username <email> \
  --token <api-token>
```

**Example:**

```bash
ticketr workspace create frontend \
  --url https://company.atlassian.net \
  --project FRONT \
  --username dev@company.com \
  --token ATATxxxxxxxxxxxxxxxx
```

### Workspace Naming Rules

Workspace names must:
- Contain only alphanumeric characters, hyphens, and underscores
- Be 64 characters or less
- Be unique across all workspaces

**Valid names:**
- `backend`
- `frontend-mobile`
- `project_123`
- `qa-automation`

**Invalid names:**
- `my project` (contains space)
- `prod@env` (contains @)
- `very-long-name-that-exceeds-the-sixty-four-character-limit-for-workspace-names` (too long)

### Interactive Creation

Ticketr can prompt for credentials if you prefer not to pass them as flags:

```bash
ticketr workspace create mobile
# You will be prompted for URL, project key, username, and token
```

**Tip:** Use the interactive mode to avoid storing credentials in shell history.

---

## Switching Between Workspaces

### Switch to a Workspace

```bash
ticketr workspace switch <name>
```

**Example:**

```bash
ticketr workspace switch frontend
```

The current workspace remains active for your entire session until you switch to another workspace or restart your terminal.

### Check Current Workspace

```bash
ticketr workspace current
```

**Output:**

```
Current workspace: frontend
Jira URL: https://company.atlassian.net
Project: FRONT
Last used: 2025-10-17 14:30:00
```

### Using the Default Workspace

If you don't explicitly switch workspaces, Ticketr uses your default workspace automatically.

To set a workspace as default:

```bash
ticketr workspace set-default backend
```

---

## Listing Workspaces

View all configured workspaces:

```bash
ticketr workspace list
```

**Output:**

```
WORKSPACE    PROJECT    JIRA URL                           DEFAULT    LAST USED
backend      BACK       https://company.atlassian.net      yes        2025-10-17 15:45
frontend     FRONT      https://company.atlassian.net      no         2025-10-16 10:20
mobile       MOB        https://mobile.atlassian.net       no         2025-10-15 08:00
```

Workspaces are listed by most recently used first.

---

## Managing Workspace Credentials

### How Credentials Are Stored

Ticketr stores credentials securely using your operating system's credential manager:

- **macOS:** Keychain Access
- **Windows:** Windows Credential Manager
- **Linux:** Secret Service (GNOME Keyring, KWallet)

Credentials are **never** stored in:
- The SQLite database
- Configuration files
- Log files
- Environment variables (after initial setup)

### Updating Credentials

If your API token expires or you need to change credentials:

```bash
ticketr workspace update backend \
  --username new.email@company.com \
  --token new-token
```

You can also update the Jira URL or project key:

```bash
ticketr workspace update backend \
  --url https://newinstance.atlassian.net \
  --project NEWPROJ
```

### Viewing Workspace Configuration

To see workspace configuration without credentials:

```bash
ticketr workspace show backend
```

**Output:**

```
Name: backend
Jira URL: https://company.atlassian.net
Project Key: BACK
Default: yes
Created: 2025-10-01 09:00:00
Last Used: 2025-10-17 15:45:00
Last Updated: 2025-10-17 14:30:00
Credential Status: Valid
```

**Note:** Actual credentials are not displayed for security reasons.

---

## Deleting Workspaces

Remove a workspace you no longer need:

```bash
ticketr workspace delete mobile
```

**Warning:** This action will:
1. Delete the workspace configuration
2. Remove credentials from the OS keychain
3. Delete all associated ticket data from the local database
4. **Cannot be undone**

**Safety check:** You cannot delete your only workspace. At least one workspace must exist.

If you delete the default workspace, Ticketr automatically sets another workspace as the default.

### Before Deleting

Consider backing up your data:

```bash
# Export workspace data
ticketr export --workspace mobile --output mobile-backup.json
```

---

## Working with Tickets in Workspaces

### Context-Aware Operations

Once you've switched to a workspace, all ticket operations use that workspace's context:

```bash
# Switch to workspace
ticketr workspace switch backend

# All subsequent commands use 'backend' workspace
ticketr push tickets.md
ticketr pull --output current-sprint.md
ticketr schema > .ticketr.yaml
```

### Explicit Workspace Override

You can override the current workspace for a single command:

```bash
# Use frontend workspace for this pull only
ticketr pull --workspace frontend --output frontend-tickets.md

# Current workspace remains unchanged
ticketr workspace current
# Output: backend
```

### Multi-Workspace Workflows

Managing tickets across multiple projects:

```bash
# Pull from backend
ticketr workspace switch backend
ticketr pull --output backend-sprint.md

# Pull from frontend
ticketr workspace switch frontend
ticketr pull --output frontend-sprint.md

# Compare and coordinate across projects
diff backend-sprint.md frontend-sprint.md
```

---

## Common Use Cases

### Use Case 1: Managing Multiple Client Projects

You work with multiple clients, each with their own Jira instance:

```bash
# Set up client workspaces
ticketr workspace create client-acme \
  --url https://acme.atlassian.net \
  --project ACME \
  --username your.email@consulting.com \
  --token token-acme

ticketr workspace create client-globex \
  --url https://globex.atlassian.net \
  --project GLOB \
  --username your.email@consulting.com \
  --token token-globex

# Switch between clients as needed
ticketr workspace switch client-acme
ticketr push acme-sprint-24.md

ticketr workspace switch client-globex
ticketr pull --output globex-backlog.md
```

### Use Case 2: Separate Environments

Manage development, staging, and production project boards:

```bash
ticketr workspace create dev-board --project DEV
ticketr workspace create staging-board --project STAGE
ticketr workspace create prod-board --project PROD

# Work with development board by default
ticketr workspace set-default dev-board

# Promote tickets to staging when ready
ticketr workspace switch staging-board
ticketr push promoted-tickets.md
```

### Use Case 3: Personal and Team Projects

Keep personal tasks separate from team work:

```bash
ticketr workspace create personal --project PERSONAL
ticketr workspace create team --project TEAM

# Set team as default for daily work
ticketr workspace set-default team

# Switch to personal for side projects
ticketr workspace switch personal
```

### Use Case 4: Migrating Between Jira Instances

Moving from one Jira instance to another:

```bash
# Create workspace for new instance
ticketr workspace create new-jira \
  --url https://newcompany.atlassian.net \
  --project PROJ

# Export from old workspace
ticketr workspace switch old-jira
ticketr pull --output migration-data.md

# Import to new workspace
ticketr workspace switch new-jira
ticketr push migration-data.md
```

---

## Security Best Practices

### Credential Management

1. **Never commit credentials:** Keep API tokens out of version control
2. **Use dedicated tokens:** Create separate API tokens for each workspace
3. **Rotate regularly:** Update tokens periodically for security
4. **Limit permissions:** Use tokens with minimal required Jira permissions

### Secure Token Generation

Generate Jira API tokens at:
https://id.atlassian.com/manage-profile/security/api-tokens

**Recommended settings:**
- Create descriptive token names (e.g., "Ticketr - Backend Project")
- Set expiration dates if your organization requires it
- Revoke unused tokens immediately

### Audit Trail

Ticketr logs all workspace operations:

```bash
# View workspace activity
tail -f ~/.local/share/ticketr/logs/workspace-activity.log
```

Log files include:
- Workspace creation/deletion
- Credential updates
- Switch operations
- Sync operations

**Note:** Credentials are never logged. Sensitive data is automatically redacted.

---

## Troubleshooting

### Cannot Create Workspace

**Error:** `workspace already exists`

**Solution:** Choose a different workspace name or delete the existing workspace first.

```bash
ticketr workspace list
ticketr workspace delete old-backend
ticketr workspace create backend --url ...
```

---

### Authentication Failure

**Error:** `failed to retrieve credentials`

**Cause:** OS keychain access denied or credentials not found

**Solutions:**

1. **macOS:** Grant Keychain Access permissions
   ```bash
   # Check keychain status
   security unlock-keychain
   ```

2. **Linux:** Ensure Secret Service is running
   ```bash
   # Check for running keyring
   ps aux | grep -i keyring
   ```

3. **Windows:** Run as administrator if Credential Manager access fails

4. **Re-create credentials:**
   ```bash
   ticketr workspace update <name> --username ... --token ...
   ```

---

### Workspace Not Found

**Error:** `workspace 'frontend' not found`

**Solution:** Verify workspace name and list available workspaces

```bash
ticketr workspace list
ticketr workspace switch backend  # Use correct name
```

---

### Cannot Delete Last Workspace

**Error:** `cannot delete the only workspace`

**Cause:** Ticketr requires at least one workspace

**Solution:** Create a new workspace before deleting the last one

```bash
ticketr workspace create new-workspace --url ...
ticketr workspace delete old-workspace
```

---

### Credentials Expired

**Error:** `Jira authentication failed: 401 Unauthorized`

**Cause:** API token expired or revoked

**Solution:** Generate a new token and update workspace

```bash
# Generate new token at: https://id.atlassian.com/manage-profile/security/api-tokens
ticketr workspace update backend --token new-token
```

---

### Multiple Default Workspaces

**Error:** `multiple default workspaces found`

**Cause:** Database corruption or manual database modification

**Solution:** Reset default workspace

```bash
# Set a single default
ticketr workspace set-default backend

# If error persists, contact support or check database
```

---

## Advanced Configuration

### Database Location

Workspaces are stored in:

```
~/.local/share/ticketr/ticketr.db (Linux)
~/Library/Application Support/ticketr/ticketr.db (macOS)
%LOCALAPPDATA%/ticketr/ticketr.db (Windows)
```

### Backup Workspaces

Export workspace configuration (without credentials):

```bash
ticketr workspace export > workspaces-backup.json
```

Restore from backup:

```bash
ticketr workspace import < workspaces-backup.json
# Note: You must re-enter credentials manually
```

### Environment Variable Override

Temporarily override workspace for scripting:

```bash
TICKETR_WORKSPACE=frontend ticketr pull --output temp.md
```

---

## Workspace Lifecycle

### Typical Workflow

1. **Initial Setup**
   ```bash
   ticketr workspace create main --url ... --project MAIN
   ```

2. **Daily Work**
   ```bash
   ticketr workspace switch main
   ticketr pull --output today.md
   # Edit tickets
   ticketr push today.md
   ```

3. **Context Switching**
   ```bash
   ticketr workspace switch side-project
   ticketr pull --output side-tasks.md
   ```

4. **Cleanup**
   ```bash
   ticketr workspace delete old-project
   ```

### Migration from v2.x

If migrating from Ticketr v2.x (environment variable configuration):

```bash
# v2.x environment variables still in .env
source .env

# Create workspace from existing config
ticketr workspace create main \
  --url $JIRA_URL \
  --project $JIRA_PROJECT_KEY \
  --username $JIRA_EMAIL \
  --token $JIRA_API_KEY

# Set as default
ticketr workspace set-default main

# Verify migration
ticketr workspace show main
```

See [Phase 2 Migration Guide](archive/phase-2-workspace-migration.md) for detailed migration instructions (archived).

---

## Related Documentation

- [Phase 2 Migration Guide](archive/phase-2-workspace-migration.md) - Migrating from v2.x (archived)
- [ARCHITECTURE.md](ARCHITECTURE.md) - Workspace domain model and security
- [v3 Implementation Roadmap](v3-implementation-roadmap.md) - v3.0 feature roadmap
- [v3 Technical Specification](v3-technical-specification.md) - Technical details

---

## FAQ

**Q: Can I use the same Jira credentials for multiple workspaces?**

A: Yes, you can use the same username and token across workspaces. Each workspace stores its own copy securely.

---

**Q: Do I need separate API tokens for each workspace?**

A: No, but it's recommended for security. Separate tokens allow you to revoke access per workspace without affecting others.

---

**Q: Can I rename a workspace?**

A: Not directly. You must create a new workspace and delete the old one. Use export/import to preserve data.

---

**Q: What happens to tickets when I delete a workspace?**

A: Local ticket data is deleted. Jira tickets are unaffected. Export data before deleting if needed.

---

**Q: Can I share workspace configurations with my team?**

A: You can export workspace metadata, but credentials must be configured individually for security.

---

**Q: How do I see which workspace a ticket belongs to?**

A: Tickets are automatically scoped to the current workspace. Check with `ticketr workspace current`.

---

**Q: Can I have workspaces pointing to different Jira instances?**

A: Yes! Workspaces support different Jira URLs, making it easy to manage multiple Jira instances.

---

**Q: Is there a limit to the number of workspaces?**

A: No hard limit, but performance may degrade with hundreds of workspaces. Recommended: < 50 workspaces.

---

**Document Version:** 1.0
**Status:** Complete
**Feedback:** Open an issue at https://github.com/karolswdev/ticketr/issues
