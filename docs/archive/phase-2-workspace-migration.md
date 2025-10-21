# Phase 2: Workspace Migration Guide

**Version:** 1.0
**Last Updated:** October 17, 2025
**Scope:** Migrating from Ticketr v2.x to v3.0 workspace model

## Purpose

This guide provides step-by-step instructions for migrating from Ticketr v2.x (single-project, environment variable configuration) to Ticketr v3.0 (multi-workspace model with secure credential storage).

---

## Overview

### What's Changing

**Ticketr v2.x (Environment Variable Model):**
- Single Jira project per installation
- Credentials in `.env` files or shell environment
- State tracked in local `.ticketr.state` files
- Must navigate to project directory to run commands

**Ticketr v3.0 (Workspace Model):**
- Multiple Jira projects from one installation
- Credentials stored securely in OS keychain
- Centralized SQLite database for all workspaces
- Run commands from anywhere in the filesystem

### Migration Goals

1. Preserve all existing ticket data and sync state
2. Move credentials from environment variables to OS keychain
3. Create workspaces for existing v2.x projects
4. Enable seamless switching between projects
5. Maintain backward compatibility during transition

---

## Prerequisites

### Before You Begin

**Required:**
- [ ] Ticketr v3.0 or later installed
- [ ] Backup of all `.ticketr.state` files
- [ ] List of all projects you currently manage with Ticketr
- [ ] Jira credentials for each project (URL, username, API token)
- [ ] OS keychain access (macOS Keychain, Windows Credential Manager, Linux Secret Service)

**Recommended:**
- [ ] Review current `.env` files for all projects
- [ ] Test credentials are still valid
- [ ] Export ticket data for backup: `ticketr pull --project PROJ --output backup.md`

### Version Check

Verify you have Ticketr v3.0 or later:

```bash
ticketr --version
```

Expected output:

```
ticketr version 3.0.0 (or later)
```

If you have an older version, upgrade first:

```bash
go install github.com/karolswdev/ticketr/cmd/ticketr@v3
```

---

## Migration Scenarios

### Scenario 1: Single Project User

You manage one Jira project with Ticketr v2.x.

**Current setup:**

```bash
# .env file
JIRA_URL=https://company.atlassian.net
JIRA_EMAIL=dev@company.com
JIRA_API_KEY=ATATxxxxxxxxxxxxxxxx
JIRA_PROJECT_KEY=PROJ
```

**Migration steps:**

1. **Navigate to your project directory:**

```bash
cd ~/projects/main-project
```

2. **Load existing environment variables:**

```bash
source .env
```

3. **Create workspace from environment:**

```bash
ticketr workspace create main \
  --url $JIRA_URL \
  --project $JIRA_PROJECT_KEY \
  --username $JIRA_EMAIL \
  --token $JIRA_API_KEY
```

4. **Verify workspace creation:**

```bash
ticketr workspace list
```

Expected output:

```
WORKSPACE    PROJECT    JIRA URL                           DEFAULT    LAST USED
main         PROJ       https://company.atlassian.net      yes        2025-10-17 14:30
```

5. **Test workspace:**

```bash
ticketr workspace current
ticketr pull --output test.md
```

6. **Optional: Remove .env file (credentials now in keychain):**

```bash
# Backup first
cp .env .env.backup

# Remove from version control if accidentally committed
echo ".env" >> .gitignore
git rm --cached .env
```

**Result:** You can now run `ticketr` commands from any directory.

---

### Scenario 2: Multiple Project User

You manage multiple Jira projects, each in its own directory with separate `.env` files.

**Current setup:**

```
~/projects/
├── backend/
│   ├── .env (JIRA_PROJECT_KEY=BACK)
│   └── tickets.md
├── frontend/
│   ├── .env (JIRA_PROJECT_KEY=FRONT)
│   └── tickets.md
└── mobile/
    ├── .env (JIRA_PROJECT_KEY=MOB)
    └── tickets.md
```

**Migration steps:**

1. **Create a migration script:**

Save this as `migrate-workspaces.sh`:

```bash
#!/bin/bash

# Project directories
PROJECTS=(
  "backend:BACK:~/projects/backend"
  "frontend:FRONT:~/projects/frontend"
  "mobile:MOB:~/projects/mobile"
)

for project in "${PROJECTS[@]}"; do
  IFS=':' read -r name key dir <<< "$project"

  echo "Migrating $name ($key)..."

  # Navigate to project directory
  cd "$dir" || continue

  # Load environment variables
  source .env

  # Create workspace
  ticketr workspace create "$name" \
    --url "$JIRA_URL" \
    --project "$key" \
    --username "$JIRA_EMAIL" \
    --token "$JIRA_API_KEY"

  echo "✓ $name workspace created"
done

# Set default workspace
ticketr workspace set-default backend

echo "Migration complete!"
ticketr workspace list
```

2. **Make script executable and run:**

```bash
chmod +x migrate-workspaces.sh
./migrate-workspaces.sh
```

3. **Verify all workspaces created:**

```bash
ticketr workspace list
```

Expected output:

```
WORKSPACE    PROJECT    JIRA URL                           DEFAULT    LAST USED
backend      BACK       https://company.atlassian.net      yes        2025-10-17 14:30
frontend     FRONT      https://company.atlassian.net      no         Never
mobile       MOB        https://mobile.atlassian.net       no         Never
```

4. **Test each workspace:**

```bash
ticketr workspace switch backend
ticketr pull --output backend-test.md

ticketr workspace switch frontend
ticketr pull --output frontend-test.md

ticketr workspace switch mobile
ticketr pull --output mobile-test.md
```

**Result:** All projects accessible from anywhere with `ticketr workspace switch <name>`.

---

### Scenario 3: Same Jira Instance, Multiple Projects

You manage multiple projects on the same Jira instance with identical credentials.

**Current setup:**

```
JIRA_URL=https://company.atlassian.net (same for all)
JIRA_EMAIL=dev@company.com (same for all)
JIRA_API_KEY=ATATxxxxxxxxxxxxxxxx (same for all)

Projects: BACK, FRONT, MOB (different project keys)
```

**Migration steps:**

1. **Create all workspaces with same credentials:**

```bash
# Set credentials once
export JIRA_URL="https://company.atlassian.net"
export JIRA_EMAIL="dev@company.com"
export JIRA_API_KEY="ATATxxxxxxxxxxxxxxxx"

# Create workspaces for each project
ticketr workspace create backend \
  --url $JIRA_URL \
  --project BACK \
  --username $JIRA_EMAIL \
  --token $JIRA_API_KEY

ticketr workspace create frontend \
  --url $JIRA_URL \
  --project FRONT \
  --username $JIRA_EMAIL \
  --token $JIRA_API_KEY

ticketr workspace create mobile \
  --url $JIRA_URL \
  --project MOB \
  --username $JIRA_EMAIL \
  --token $JIRA_API_KEY
```

2. **Set default workspace:**

```bash
ticketr workspace set-default backend
```

3. **Verify and test:**

```bash
ticketr workspace list
ticketr workspace switch frontend
ticketr pull --output frontend-tickets.md
```

**Note:** Each workspace stores its own copy of credentials in the keychain. This allows you to rotate tokens per workspace later if needed.

---

### Scenario 4: Different Jira Instances

You work with multiple clients, each with their own Jira instance and credentials.

**Current setup:**

```
Client A: https://clienta.atlassian.net
Client B: https://clientb.atlassian.net
Client C: https://clientc.atlassian.net

Different credentials for each
```

**Migration steps:**

1. **Create workspace for each client:**

```bash
# Client A
ticketr workspace create client-a \
  --url https://clienta.atlassian.net \
  --project CA \
  --username dev@clienta.com \
  --token clienta-token

# Client B
ticketr workspace create client-b \
  --url https://clientb.atlassian.net \
  --project CB \
  --username dev@clientb.com \
  --token clientb-token

# Client C
ticketr workspace create client-c \
  --url https://clientc.atlassian.net \
  --project CC \
  --username consulting@yourcompany.com \
  --token clientc-token
```

2. **Test each workspace:**

```bash
ticketr workspace switch client-a
ticketr schema  # Verify credentials work

ticketr workspace switch client-b
ticketr schema

ticketr workspace switch client-c
ticketr schema
```

**Result:** Seamless switching between client projects with isolated credentials.

---

## Default Workspace Creation

### Automatic Default Selection

When you create your first workspace, Ticketr automatically sets it as the default:

```bash
ticketr workspace create main --url ... --project PROJ
# Automatically becomes default
```

Subsequent commands use this workspace unless you explicitly switch:

```bash
ticketr push tickets.md  # Uses 'main' workspace
ticketr pull --output tickets.md  # Uses 'main' workspace
```

### Changing the Default

Set a different workspace as default:

```bash
ticketr workspace set-default frontend
```

### Checking the Default

```bash
ticketr workspace list
# Look for 'yes' in the DEFAULT column
```

Or:

```bash
ticketr workspace current
# Shows default workspace if no explicit switch
```

---

## Credential Migration Process

### From Environment Variables to Keychain

**v2.x Credential Storage:**

```bash
# .env file (plain text)
JIRA_URL=https://company.atlassian.net
JIRA_EMAIL=dev@company.com
JIRA_API_KEY=ATATxxxxxxxxxxxxxxxx
JIRA_PROJECT_KEY=PROJ
```

**v3.0 Credential Storage:**

```
OS Keychain (encrypted)
├── Service: ticketr
├── Account: <workspace-id>
└── Data: {"username": "dev@company.com", "apiToken": "ATATxxxxxxxxxxxxxxxx"}

SQLite Database (reference only)
└── workspaces.credential_keychain_id: <workspace-id>
```

### Migration Flow

1. **Read credentials from environment or .env file**
2. **Pass credentials to `ticketr workspace create`**
3. **Ticketr stores credentials in OS keychain**
4. **Ticketr stores only a reference in SQLite**
5. **Original .env file can be safely deleted**

### Verifying Credential Migration

```bash
# Check workspace shows valid credentials
ticketr workspace show main

# Output includes:
# Credential Status: Valid
```

If credential status shows "Invalid" or "Missing":

```bash
# Re-add credentials
ticketr workspace update main \
  --username dev@company.com \
  --token new-token
```

---

## State File Migration

### v2.x State Files

**Location:** `.ticketr.state` in each project directory

**Format:**

```json
{
  "PROJ-123": {
    "local_hash": "abc123...",
    "remote_hash": "def456..."
  }
}
```

### v3.0 State Storage

**Location:** SQLite database at `~/.local/share/ticketr/ticketr.db`

**Schema:**

```sql
CREATE TABLE tickets (
    id TEXT PRIMARY KEY,
    workspace_id TEXT,
    jira_id TEXT,
    local_hash TEXT,
    remote_hash TEXT,
    -- ... other fields
);
```

### State Migration Strategy

**Phase 2 (Current):** Workspaces exist, but tickets still use file-based state

- `.ticketr.state` files remain functional
- Workspace credentials stored in keychain
- Ticket sync state still tracked per-directory

**Phase 3 (Future):** Full migration to SQLite

- State imported from `.ticketr.state` files
- Centralized state in SQLite database
- Automatic migration tool: `ticketr migrate state`

**Action Required for Phase 2:** None. Your existing `.ticketr.state` files continue to work.

---

## Validation and Testing

### Post-Migration Checklist

After creating workspaces, verify everything works:

- [ ] **List workspaces:** `ticketr workspace list`
- [ ] **Check default:** `ticketr workspace current`
- [ ] **Test pull:** `ticketr pull --output test.md`
- [ ] **Test push:** `ticketr push test.md`
- [ ] **Test workspace switch:** `ticketr workspace switch <name>`
- [ ] **Verify credentials:** `ticketr schema`
- [ ] **Test from different directory:** `cd /tmp && ticketr workspace list`

### Common Validation Issues

**Issue:** Workspace created but credentials invalid

```bash
ticketr workspace show main
# Shows: Credential Status: Invalid
```

**Solution:**

```bash
ticketr workspace update main \
  --username correct-email@company.com \
  --token valid-token
```

---

**Issue:** Cannot access OS keychain

**Error:** `failed to store credentials: keychain access denied`

**Solution (macOS):**

```bash
# Unlock keychain
security unlock-keychain

# Grant access to Ticketr
# (Keychain Access app will prompt for permission)
```

**Solution (Linux):**

```bash
# Ensure GNOME Keyring or KWallet is running
ps aux | grep -i keyring

# If not running, start it
gnome-keyring-daemon --start
```

**Solution (Windows):**

```powershell
# Run PowerShell as Administrator
# Credential Manager should be accessible by default
```

---

**Issue:** Multiple workspaces, wrong one being used

**Solution:**

```bash
# Check current workspace
ticketr workspace current

# Switch to correct workspace
ticketr workspace switch desired-workspace

# Or set correct default
ticketr workspace set-default desired-workspace
```

---

**Issue:** Old .env file still being read

**Solution:**

```bash
# Unset environment variables in current shell
unset JIRA_URL
unset JIRA_EMAIL
unset JIRA_API_KEY
unset JIRA_PROJECT_KEY

# Remove from .bashrc or .zshrc if present
grep -v "JIRA_" ~/.bashrc > ~/.bashrc.new
mv ~/.bashrc.new ~/.bashrc

# Restart terminal or source file
source ~/.bashrc
```

---

## Rollback Procedure

If you encounter issues and need to revert to v2.x:

### Step 1: Reinstall v2.x

```bash
go install github.com/karolswdev/ticketr/cmd/ticketr@v2.0.0
```

### Step 2: Restore .env files

```bash
# If you backed up .env files
cp .env.backup .env

# Or manually recreate
cat > .env << EOF
JIRA_URL=https://company.atlassian.net
JIRA_EMAIL=dev@company.com
JIRA_API_KEY=ATATxxxxxxxxxxxxxxxx
JIRA_PROJECT_KEY=PROJ
EOF
```

### Step 3: Verify v2.x operation

```bash
cd ~/projects/main-project
source .env
ticketr push tickets.md
```

### Step 4: Remove v3.0 database (optional)

```bash
# Backup first
cp ~/.local/share/ticketr/ticketr.db ~/ticketr-v3-backup.db

# Remove database
rm ~/.local/share/ticketr/ticketr.db
```

**Note:** v2.x and v3.0 can coexist if you use different binaries (e.g., `ticketr2` and `ticketr3`).

---

## Best Practices

### Naming Conventions

Use consistent workspace names:

- **Good:** `backend`, `frontend-web`, `mobile-ios`
- **Avoid:** `proj1`, `temp`, `test123`

Workspace names should be:
- Descriptive of the project
- Easy to remember
- Short enough to type quickly

### Credential Management

**Security tips:**

1. **Generate separate API tokens per workspace**
   - Easier to revoke if compromised
   - Clear audit trail per project

2. **Set token expiration dates**
   - Many organizations require token rotation
   - Ticketr will prompt you to update expired tokens

3. **Never commit .env files**
   - Add to `.gitignore` immediately
   - Use `.env.example` for team onboarding

4. **Use keychain password managers**
   - Store backup of tokens in 1Password, LastPass, etc.
   - In case you need to recreate workspace

### Workspace Organization

**For personal projects:**
```bash
ticketr workspace create personal-backend --project PERS
ticketr workspace create personal-mobile --project PMOB
```

**For work projects:**
```bash
ticketr workspace create work-platform --project PLAT
ticketr workspace create work-infra --project INFRA
```

**For client work:**
```bash
ticketr workspace create client-acme-backend --project ACME
ticketr workspace create client-globex-api --project GLOB
```

### Default Workspace Selection

Set your most frequently used workspace as default:

```bash
# If you primarily work on backend
ticketr workspace set-default backend

# All commands default to backend unless overridden
ticketr pull --output sprint.md  # Uses backend
```

---

## Troubleshooting

### Migration Issues

**Issue:** Workspace creation fails with "workspace already exists"

**Cause:** You've already created a workspace with that name

**Solution:**

```bash
# List existing workspaces
ticketr workspace list

# Use a different name
ticketr workspace create backend-v2 --url ...

# Or delete existing workspace first
ticketr workspace delete backend
ticketr workspace create backend --url ...
```

---

**Issue:** Credentials not working after migration

**Cause:** Token expired, incorrect username, or wrong Jira URL

**Solution:**

```bash
# Test credentials manually
curl -u your.email@company.com:ATATxxxxxxxxxxxxxxxx \
  https://company.atlassian.net/rest/api/2/myself

# If 401 Unauthorized, regenerate token
# Update workspace with new token
ticketr workspace update main --token new-token
```

---

**Issue:** Cannot find .ticketr.state files

**Cause:** Files may be in subdirectories or gitignored

**Solution:**

```bash
# Search for all .ticketr.state files
find ~ -name ".ticketr.state" 2>/dev/null

# Check gitignored files
git status --ignored | grep ticketr
```

---

## Migration Timeline

### Phase 2 (Current) - Workspace Model

**Available Now:**
- ✅ Workspace creation and management
- ✅ Credential storage in OS keychain
- ✅ Workspace switching
- ✅ Multiple Jira project support

**Still File-Based:**
- `.ticketr.state` files (per-directory)
- Ticket storage (Markdown files)

### Phase 3 (Future) - Full SQLite Migration

**Planned:**
- Centralized ticket storage in SQLite
- Automatic state file import
- Global ticket queries across workspaces
- Workspace-scoped ticket filtering

**Migration Path:**

```bash
# When Phase 3 releases
ticketr migrate state --from-files

# Imports all .ticketr.state files into SQLite
# Preserves sync state and history
```

---

## FAQ

**Q: Do I have to migrate immediately?**

A: No. Ticketr v3.0 maintains backward compatibility with v2.x workflows. You can continue using environment variables and file-based state.

---

**Q: Can I use both v2.x and v3.0 simultaneously?**

A: Yes, if you maintain separate binaries. Install v2.x as `ticketr2` and v3.0 as `ticketr3`.

---

**Q: Will my existing tickets and sync state be lost?**

A: No. Phase 2 preserves existing `.ticketr.state` files. Phase 3 will import them into SQLite.

---

**Q: What happens to my .env files after migration?**

A: You can safely delete them. Credentials are stored in the OS keychain. Keep a backup for rollback purposes.

---

**Q: Can I migrate only some projects to workspaces?**

A: Yes. Create workspaces for projects you want to manage centrally. Others can continue using v2.x workflow.

---

**Q: How do I share workspace configuration with my team?**

A: Export workspace metadata (without credentials):

```bash
ticketr workspace export > workspaces.json
```

Team members import and add their own credentials:

```bash
ticketr workspace import < workspaces.json
ticketr workspace update main --token their-token
```

---

**Q: Can I use the same workspace name across different machines?**

A: Yes. Workspace names are local to each machine. Your teammate can have a "backend" workspace independent of yours.

---

**Q: What if I forget which workspace I'm using?**

A: Check anytime with:

```bash
ticketr workspace current
```

Or add to your shell prompt (see [Advanced Configuration](#advanced-configuration)).

---

## Advanced Configuration

### Shell Prompt Integration

Add current workspace to your shell prompt:

**Bash:**

```bash
# Add to ~/.bashrc
ticketr_workspace() {
  ticketr workspace current --short 2>/dev/null
}

export PS1="\u@\h [\$(ticketr_workspace)] \w \$ "
```

**Zsh:**

```zsh
# Add to ~/.zshrc
ticketr_workspace() {
  ticketr workspace current --short 2>/dev/null
}

PROMPT='%n@%m [$(ticketr_workspace)] %~ %# '
```

**Result:**

```
user@laptop [backend] ~/projects $
```

### Workspace Aliases

Create shell aliases for quick switching:

```bash
# Add to ~/.bashrc or ~/.zshrc
alias ws-backend='ticketr workspace switch backend'
alias ws-frontend='ticketr workspace switch frontend'
alias ws-mobile='ticketr workspace switch mobile'
```

Usage:

```bash
ws-backend
ticketr pull --output sprint.md
```

---

## Related Documentation

- [Workspace Guide](workspace-guide.md) - Complete workspace usage guide
- [ARCHITECTURE.md](ARCHITECTURE.md) - Workspace technical architecture
- [v3 Implementation Roadmap](v3-implementation-roadmap.md) - v3.0 feature roadmap
- [Migration Guide](migration-guide.md) - v1.x to v2.0 migration

---

## Support

If you encounter issues during migration:

1. **Check this guide** for troubleshooting steps
2. **Review logs** at `~/.cache/ticketr/logs/`
3. **Open an issue** at https://github.com/karolswdev/ticketr/issues
4. **Join discussions** at https://github.com/karolswdev/ticketr/discussions

---

**Document Version:** 1.0
**Migration Status:** Phase 2 Complete
**Next Phase:** SQLite state migration (Phase 3)
