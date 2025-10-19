# Workspace Management Guide

**Last Updated:** October 18, 2025
**Version:** Ticketr v3.0
**Status:** Production-ready

---

## Table of Contents

1. [Introduction](#introduction)
2. [What Are Workspaces?](#what-are-workspaces)
3. [When to Use Workspaces](#when-to-use-workspaces)
4. [Getting Started](#getting-started)
5. [Workspace Commands Reference](#workspace-commands-reference)
6. [Security Model](#security-model)
7. [Troubleshooting](#troubleshooting)
8. [Best Practices](#best-practices)
9. [Migrating from v2.x](#migrating-from-v2x)

---

## Introduction

Ticketr v3.0 introduces **workspace management**, enabling you to work with multiple Jira projects from a single Ticketr installation. Each workspace represents a distinct Jira instance and project configuration with isolated credentials.

### Key Benefits

- **Multi-Project Support**: Manage tickets across different Jira projects seamlessly
- **Secure Credential Storage**: OS-level keychain encryption protects your API tokens
- **Fast Context Switching**: Switch between workspaces instantly
- **Isolated Environments**: Each workspace maintains separate credentials and state

---

## What Are Workspaces?

A **workspace** is a named configuration that contains:

- **Jira Instance URL** (e.g., `https://company.atlassian.net`)
- **Project Key** (e.g., `BACKEND`, `FRONTEND`)
- **Credentials** (username and API token, stored in OS keychain)
- **Metadata** (creation date, last used timestamp)

### Workspace vs. Project

| Concept | Description | Example |
|---------|-------------|---------|
| **Workspace** | Ticketr configuration for a Jira project | `backend`, `frontend`, `mobile` |
| **Jira Project** | Actual project in Jira Cloud/Server | `BACK`, `FRONT`, `MOB` |

**Note:** One workspace = One Jira project. You can create multiple workspaces for different Jira projects or even different Jira instances.

---

## When to Use Workspaces

### Use Cases

**1. Multiple Teams/Projects**
```bash
# DevOps team
ticketr workspace create devops --url https://company.atlassian.net --project DEVOPS

# Backend team
ticketr workspace create backend --url https://company.atlassian.net --project BACK

# Frontend team
ticketr workspace create frontend --url https://company.atlassian.net --project FRONT
```

**2. Multiple Jira Instances**
```bash
# Company Jira
ticketr workspace create work --url https://company.atlassian.net --project WORK

# Personal Jira
ticketr workspace create personal --url https://personal.atlassian.net --project PERS
```

**3. Client Projects**
```bash
# Client A
ticketr workspace create client-a --url https://clienta.atlassian.net --project PROJA

# Client B
ticketr workspace create client-b --url https://clientb.atlassian.net --project PROJB
```

**4. Development Environments**
```bash
# Production Jira
ticketr workspace create production --url https://prod.company.net --project PROD

# Staging Jira
ticketr workspace create staging --url https://staging.company.net --project STAGE
```

### When NOT to Use Multiple Workspaces

- **Single Project**: If you only work with one Jira project, you don't need multiple workspaces
- **Same Project, Different Filters**: Use JQL queries instead (e.g., `ticketr pull --jql "assignee=currentUser()"`)
- **Temporary Access**: For one-time operations, environment variables may be simpler

---

## Getting Started

### Prerequisites

- Ticketr v3.0+ installed
- Jira API token (get from [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens))
- Access to your OS keychain (macOS Keychain Access, Windows Credential Manager, or Linux Secret Service)

### Step 1: Create Your First Workspace

```bash
ticketr workspace create backend \
  --url https://company.atlassian.net \
  --project BACK \
  --username your.email@company.com \
  --token your-api-token
```

**Expected Output:**
```
✓ Workspace 'backend' created successfully
✓ Credentials stored securely in OS keychain
✓ Set as default workspace
```

**What Happens:**
1. Ticketr validates your Jira credentials
2. Credentials are encrypted and stored in your OS keychain
3. Workspace metadata is saved to SQLite database
4. The first workspace is automatically set as default

### Step 2: Verify Your Workspace

```bash
ticketr workspace list
```

**Expected Output:**
```
Available Workspaces:
  * backend (default) - https://company.atlassian.net - Project: BACK
    Last used: 2025-10-18 14:30:00
```

### Step 3: Use Your Workspace

```bash
# Push tickets using the default workspace
ticketr push tickets.md

# Pull tickets from the default workspace
ticketr pull --output tickets.md
```

### Step 4: Create Additional Workspaces

```bash
ticketr workspace create frontend \
  --url https://company.atlassian.net \
  --project FRONT \
  --username your.email@company.com \
  --token your-api-token
```

### Step 5: Switch Between Workspaces

```bash
# Switch to frontend workspace
ticketr workspace switch frontend

# Verify current workspace
ticketr workspace current
```

**Output:**
```
Current workspace: frontend
Jira URL: https://company.atlassian.net
Project: FRONT
```

---

## Workspace Commands Reference

### Create a Workspace

**Syntax:**
```bash
ticketr workspace create <name> \
  --url <jira-url> \
  --project <project-key> \
  --username <email> \
  --token <api-token>
```

**Parameters:**
- `<name>`: Workspace name (alphanumeric, hyphens, underscores, max 64 chars)
- `--url`: Jira instance URL (e.g., `https://company.atlassian.net`)
- `--project`: Jira project key (e.g., `BACK`, `FRONT`)
- `--username`: Your Jira email address
- `--token`: Jira API token

**Examples:**
```bash
# Create workspace for backend project
ticketr workspace create backend \
  --url https://company.atlassian.net \
  --project BACK \
  --username dev@company.com \
  --token abc123xyz

# Create workspace for personal Jira
ticketr workspace create personal \
  --url https://me.atlassian.net \
  --project HOME \
  --username me@gmail.com \
  --token xyz789abc
```

**Validation Rules:**
- Workspace name must be unique
- Workspace name cannot contain spaces or special characters (except `-` and `_`)
- Jira URL must be valid and reachable
- Credentials must be valid (tested during creation)

---

### List Workspaces

**Syntax:**
```bash
ticketr workspace list
```

**Output Format:**
```
Available Workspaces:
  * backend (default) - https://company.atlassian.net - Project: BACK
    Last used: 2025-10-18 14:30:00

  frontend - https://company.atlassian.net - Project: FRONT
    Last used: 2025-10-17 10:15:00

  mobile - https://company.atlassian.net - Project: MOB
    Last used: 2025-10-16 09:00:00
```

**Legend:**
- `*` = Default workspace
- Sorted by last used (most recent first)

---

### Switch Workspaces

**Syntax:**
```bash
ticketr workspace switch <name>
```

**Examples:**
```bash
# Switch to frontend workspace
ticketr workspace switch frontend

# Switch to mobile workspace
ticketr workspace switch mobile
```

**Expected Output:**
```
✓ Switched to workspace 'frontend'
Jira URL: https://company.atlassian.net
Project: FRONT
```

**Effects:**
- All subsequent `push`/`pull` commands use the new workspace
- `LastUsed` timestamp is updated
- Workspace credentials are loaded from keychain

---

### Show Current Workspace

**Syntax:**
```bash
ticketr workspace current
```

**Output:**
```
Current workspace: backend
Jira URL: https://company.atlassian.net
Project: BACK
Last used: 2025-10-18 14:30:00
```

---

### Set Default Workspace

**Syntax:**
```bash
ticketr workspace set-default <name>
```

**Examples:**
```bash
# Set backend as default
ticketr workspace set-default backend
```

**Expected Output:**
```
✓ Workspace 'backend' set as default
```

**Effects:**
- The specified workspace becomes the default
- When Ticketr starts, it automatically uses the default workspace
- Only one workspace can be default at a time

---

### Delete a Workspace

**Syntax:**
```bash
ticketr workspace delete <name> [--force]
```

**Parameters:**
- `<name>`: Workspace name to delete
- `--force`: Skip confirmation prompt

**Examples:**
```bash
# Delete workspace with confirmation
ticketr workspace delete old-project

# Delete workspace without confirmation
ticketr workspace delete old-project --force
```

**Expected Output:**
```
⚠ Warning: This will delete workspace 'old-project' and all associated credentials.
Continue? (y/N): y
✓ Workspace 'old-project' deleted successfully
✓ Credentials removed from keychain
```

**Effects:**
- Workspace metadata deleted from database
- Credentials removed from OS keychain
- If deleting the default workspace, another workspace is automatically set as default
- Cannot delete the only remaining workspace

---

## Security Model

### Credential Storage Architecture

Ticketr v3.0 uses a **zero-trust credential model**:

```
User Input → CredentialStore → OS Keychain (encrypted)
                     ↓
              SQLite Database (CredentialRef only, NO actual credentials)
```

### What Is Stored Where

| Data | Location | Security |
|------|----------|----------|
| **Workspace name** | SQLite database | Plaintext (non-sensitive) |
| **Jira URL** | SQLite database | Plaintext (non-sensitive) |
| **Project key** | SQLite database | Plaintext (non-sensitive) |
| **Username** | OS keychain | Encrypted by OS |
| **API token** | OS keychain | Encrypted by OS |
| **Credential reference** | SQLite database | Reference only (keychain ID) |

### Platform-Specific Keychain Locations

**macOS:**
- **Keychain:** Keychain Access (`/Library/Keychains/`)
- **Service Name:** `ticketr`
- **Account Name:** `<workspace-id>`
- **Encryption:** 256-bit AES (managed by macOS)

**Windows:**
- **Keychain:** Windows Credential Manager
- **Service Name:** `ticketr`
- **Account Name:** `<workspace-id>`
- **Encryption:** DPAPI (Data Protection API)

**Linux:**
- **Keychain:** Secret Service (GNOME Keyring or KWallet)
- **Service Name:** `ticketr`
- **Account Name:** `<workspace-id>`
- **Encryption:** Varies by keyring implementation (typically AES-256)

### Security Guarantees

1. **No credentials in database**: SQLite contains only references to keychain entries
2. **No credentials in logs**: Automatic redaction prevents logging of sensitive data
3. **No credentials in memory longer than needed**: Cleared after use
4. **OS-level encryption**: All credentials encrypted at rest by OS keychain
5. **Per-user isolation**: Credentials accessible only to the OS user account

### Accessing Your Credentials

**macOS:**
1. Open **Keychain Access** app
2. Search for `ticketr`
3. View credential entries by workspace ID

**Windows:**
1. Open **Credential Manager** (Control Panel → User Accounts → Credential Manager)
2. Select **Windows Credentials**
3. Search for `ticketr`

**Linux:**
```bash
# Using GNOME Keyring
seahorse  # GUI

# Using secret-tool
secret-tool search service ticketr
```

### Revoking Credentials

If you need to revoke Jira access:

1. **Delete workspace** (removes from keychain automatically):
   ```bash
   ticketr workspace delete <name> --force
   ```

2. **Manually delete from keychain**:
   - macOS: Keychain Access → Search "ticketr" → Delete
   - Windows: Credential Manager → Remove
   - Linux: `secret-tool clear service ticketr account <workspace-id>`

3. **Revoke API token in Jira**:
   - Visit [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens)
   - Revoke the token

---

## Troubleshooting

### Common Issues

#### 1. Keychain Access Denied

**Symptom:**
```
Error: Failed to store credentials in keychain
Reason: Access denied
```

**Solution (macOS):**
1. Open **Keychain Access**
2. Right-click `ticketr` entry → Get Info
3. Select **Access Control** tab
4. Add `ticketr` binary to allowed applications

**Solution (Linux):**
```bash
# Unlock keyring
secret-tool store --label='test' service test account test

# Check keyring daemon is running
ps aux | grep keyring
```

**Solution (Windows):**
- Run Ticketr as the same user who created the workspace
- Check Windows Credential Manager permissions

---

#### 2. Workspace Not Found

**Symptom:**
```
Error: Workspace 'backend' not found
```

**Solution:**
```bash
# List all workspaces
ticketr workspace list

# Verify workspace name (case-sensitive)
ticketr workspace create backend --url ... --project ...
```

---

#### 3. Invalid Credentials

**Symptom:**
```
Error: Authentication failed
Status: 401 Unauthorized
```

**Solution:**
1. **Regenerate API token**:
   - Visit [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens)
   - Create new token

2. **Update workspace credentials**:
   ```bash
   # Delete old workspace
   ticketr workspace delete backend --force

   # Recreate with new token
   ticketr workspace create backend \
     --url https://company.atlassian.net \
     --project BACK \
     --username your.email@company.com \
     --token new-api-token
   ```

---

#### 4. No Default Workspace

**Symptom:**
```
Error: No default workspace configured
```

**Solution:**
```bash
# Set a workspace as default
ticketr workspace set-default backend
```

---

#### 5. Keyring Daemon Not Running (Linux)

**Symptom:**
```
Error: Cannot connect to secret service
```

**Solution:**
```bash
# Start GNOME Keyring
gnome-keyring-daemon --start

# Or use KWallet
kwalletd5
```

---

### Debug Mode

Enable verbose logging for troubleshooting:

```bash
# Enable debug mode
export TICKETR_LOG_LEVEL=debug

# Run workspace command
ticketr workspace list --verbose
```

**Log Location:** `.ticketr/logs/<timestamp>.log`

---

## Best Practices

### 1. Naming Conventions

**Good Names:**
```bash
ticketr workspace create backend-prod --url ... --project BACK
ticketr workspace create frontend-staging --url ... --project FRONT
ticketr workspace create client-acme --url ... --project ACME
```

**Avoid:**
```bash
ticketr workspace create "My Backend Project"  # Spaces not allowed
ticketr workspace create backend/prod          # Special chars not allowed
ticketr workspace create ThisIsAReallyLongWorkspaceNameThatExceeds64Characters  # Too long
```

**Recommendations:**
- Use lowercase for consistency
- Use hyphens to separate words (`backend-prod`, not `backend_prod`)
- Include environment suffix for clarity (`-prod`, `-staging`, `-dev`)
- Keep names short (< 20 chars ideal)

---

### 2. Credential Rotation

Rotate API tokens every 90 days:

```bash
# 1. Generate new token in Jira
# 2. Delete old workspace
ticketr workspace delete backend --force

# 3. Recreate with new token
ticketr workspace create backend \
  --url https://company.atlassian.net \
  --project BACK \
  --username your.email@company.com \
  --token new-token-here

# 4. Verify
ticketr workspace current
```

---

### 3. Workspace Organization

**For Teams:**
```bash
# Organize by team/project
ticketr workspace create backend --url ... --project BACK
ticketr workspace create frontend --url ... --project FRONT
ticketr workspace create mobile --url ... --project MOB
ticketr workspace create devops --url ... --project OPS
```

**For Consultants:**
```bash
# Organize by client
ticketr workspace create client-a --url ... --project PROJA
ticketr workspace create client-b --url ... --project PROJB
ticketr workspace create internal --url ... --project INT
```

**For Environments:**
```bash
# Organize by deployment stage
ticketr workspace create prod --url ... --project PROD
ticketr workspace create staging --url ... --project STAGE
ticketr workspace create dev --url ... --project DEV
```

---

### 4. Default Workspace Strategy

**Single Primary Project:**
```bash
# Set most frequently used project as default
ticketr workspace set-default backend
```

**Multiple Projects:**
```bash
# Switch workspaces as needed, no default
ticketr workspace switch backend
ticketr push tickets.md

ticketr workspace switch frontend
ticketr pull --output frontend-tickets.md
```

---

### 5. Backup and Recovery

**Backup Workspace List:**
```bash
ticketr workspace list > workspaces-backup.txt
```

**Recreate Workspaces:**
```bash
# Manually recreate from backup list
ticketr workspace create backend --url ... --project ...
ticketr workspace create frontend --url ... --project ...
```

**Note:** Credentials are NOT backed up (stored securely in OS keychain). You must regenerate API tokens when recreating workspaces.

---

## Migrating from v2.x

### Overview

Ticketr v2.x used environment variables for Jira credentials. Ticketr v3.0 introduces workspaces with OS keychain storage.

### Migration Steps

**Step 1: Verify Current Configuration**

In v2.x, you had environment variables:
```bash
export JIRA_URL="https://company.atlassian.net"
export JIRA_EMAIL="your.email@company.com"
export JIRA_API_KEY="your-api-token"
export JIRA_PROJECT_KEY="PROJ"
```

**Step 2: Create Workspace from Environment Variables**

```bash
ticketr workspace create main \
  --url "$JIRA_URL" \
  --project "$JIRA_PROJECT_KEY" \
  --username "$JIRA_EMAIL" \
  --token "$JIRA_API_KEY"
```

**Step 3: Verify Workspace**

```bash
ticketr workspace current
```

**Expected Output:**
```
Current workspace: main
Jira URL: https://company.atlassian.net
Project: PROJ
```

**Step 4: Test Operations**

```bash
# Test push (same as v2.x)
ticketr push tickets.md

# Test pull (same as v2.x)
ticketr pull --output tickets.md
```

**Step 5: Remove Environment Variables (Optional)**

```bash
# Remove from .env or shell profile
unset JIRA_URL
unset JIRA_EMAIL
unset JIRA_API_KEY
unset JIRA_PROJECT_KEY
```

**Note:** Environment variables still work in v3.0 for backward compatibility, but workspaces are recommended for multi-project workflows.

### Migration Checklist

- [ ] Backup existing `.env` file
- [ ] Create workspace from environment variables
- [ ] Verify workspace credentials work
- [ ] Test push/pull operations
- [ ] Update CI/CD pipelines (if applicable)
- [ ] Document workspace names for team
- [ ] Remove environment variables (optional)

---

## Advanced Topics

### Using Workspaces in Scripts

```bash
#!/bin/bash
# sync-all-projects.sh

# Switch to backend and sync
ticketr workspace switch backend
ticketr push backend-tickets.md

# Switch to frontend and sync
ticketr workspace switch frontend
ticketr push frontend-tickets.md

# Switch back to default
ticketr workspace switch backend
```

### CI/CD Integration

**GitHub Actions:**
```yaml
- name: Create temporary workspace
  run: |
    ticketr workspace create ci-backend \
      --url ${{ secrets.JIRA_URL }} \
      --project BACK \
      --username ${{ secrets.JIRA_EMAIL }} \
      --token ${{ secrets.JIRA_API_KEY }}

- name: Sync tickets
  run: ticketr push tickets.md

- name: Cleanup workspace
  run: ticketr workspace delete ci-backend --force
```

---

## Summary

Ticketr v3.0 workspaces provide:

- **Multi-project support** for complex workflows
- **Secure credential storage** using OS keychain
- **Fast context switching** between Jira projects
- **Clean migration path** from v2.x environment variables

For more information:
- [README.md](../README.md) - Quick start guide
- [ARCHITECTURE.md](ARCHITECTURE.md) - Technical architecture
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues

---

**Document Version:** 1.0
**Status:** Production-ready
**Feedback:** [GitHub Issues](https://github.com/karolswdev/ticketr/issues)
