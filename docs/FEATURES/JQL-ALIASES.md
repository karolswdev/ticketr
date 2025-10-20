# JQL Aliases Guide

**Version:** 1.0
**Last Updated:** 2025-10-20
**Scope:** JQL alias management for Ticketr v3.0

## Purpose

JQL Aliases provide a powerful way to create reusable, named JQL queries for faster and more consistent ticket filtering. Instead of remembering complex JQL syntax or repeatedly typing the same queries, you can create memorable aliases like `mine`, `sprint`, or `urgent-bugs` and use them across all your ticket operations.

## Table of Contents

- [Overview](#overview)
- [Predefined Aliases](#predefined-aliases)
- [Creating Custom Aliases](#creating-custom-aliases)
- [Workspace-Specific vs Global Aliases](#workspace-specific-vs-global-aliases)
- [Recursive Alias References](#recursive-alias-references)
- [CLI Command Reference](#cli-command-reference)
- [Using Aliases with Pull](#using-aliases-with-pull)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## Overview

### What are JQL Aliases?

JQL (Jira Query Language) aliases are named shortcuts for complex JQL queries. They allow you to:

- **Save time**: Create once, use everywhere
- **Ensure consistency**: Team members can share the same queries
- **Simplify workflows**: No need to remember complex JQL syntax
- **Build composable queries**: Reference aliases within other aliases

### Key Features

- **Predefined system aliases**: `mine`, `sprint`, `blocked` available out-of-the-box
- **Custom user aliases**: Create your own workspace-specific or global aliases
- **Recursive expansion**: Build complex queries by referencing other aliases
- **Circular reference detection**: Prevents infinite loops in alias chains
- **Workspace isolation**: Team-specific aliases don't conflict across projects

---

## Predefined Aliases

Ticketr includes three predefined system aliases available in all workspaces:

### mine

**JQL**: `assignee = currentUser() AND resolution = Unresolved`

**Description**: Tickets assigned to the current user that are not yet resolved.

**Use Case**: Daily standup prep, checking your work queue

**Example**:
```bash
ticketr pull --alias mine --output my-work.md
```

### sprint

**JQL**: `sprint in openSprints()`

**Description**: Tickets in currently active sprints.

**Use Case**: Sprint planning, checking sprint progress

**Example**:
```bash
ticketr pull --alias sprint --output current-sprint.md
```

### blocked

**JQL**: `status = Blocked OR labels in (blocked)`

**Description**: Tickets that are blocked or have the blocked label.

**Use Case**: Identifying and resolving blockers

**Example**:
```bash
ticketr pull --alias blocked --output blockers.md
```

**Note**: Predefined aliases cannot be modified or deleted. They are always available in all workspaces.

---

## Creating Custom Aliases

### Basic Alias Creation

Create a workspace-specific alias:

```bash
ticketr alias create <name> "<jql-query>"
```

**Example**:
```bash
ticketr alias create my-bugs "assignee = currentUser() AND type = Bug"
```

### Alias with Description

Add a description to document the alias purpose:

```bash
ticketr alias create urgent "priority = Highest AND status != Done" \
  --description "Urgent open tickets"
```

### Global Alias

Create an alias available across all workspaces:

```bash
ticketr alias create critical "priority = Critical" --global
```

### Alias Naming Rules

Alias names must follow these rules:

- **Alphanumeric characters**: a-z, A-Z, 0-9
- **Special characters**: Hyphens (-) and underscores (_) only
- **Length**: 1-64 characters
- **Reserved names**: Cannot use predefined alias names (mine, sprint, blocked)

**Valid Names**:
- `my-bugs`
- `sprint_tickets`
- `urgent_work_2024`

**Invalid Names**:
- `my bugs` (contains space)
- `my@bugs` (invalid character)
- `mine` (reserved predefined alias)

---

## Workspace-Specific vs Global Aliases

### Workspace-Specific Aliases (Default)

By default, aliases are created for the current workspace only.

**Benefits**:
- Isolated from other projects
- Can use project-specific field values
- Team collaboration within a workspace

**Example**:
```bash
# Create for current workspace
ticketr alias create backend-bugs "project = BACK AND type = Bug"
```

**Scope**: Only available when the workspace is active.

### Global Aliases

Global aliases are available across all workspaces.

**Benefits**:
- Consistent across all projects
- Useful for generic queries
- Single definition, universal access

**Example**:
```bash
# Create global alias
ticketr alias create high-priority "priority = High" --global
```

**Scope**: Available in all workspaces.

### Resolution Priority

When an alias name exists in both workspace and global scope:

1. **Workspace-specific alias** takes precedence
2. **Global alias** used as fallback
3. **Predefined alias** used if no custom aliases exist

This allows workspaces to override global aliases with project-specific versions.

---

## Recursive Alias References

### What are Recursive Aliases?

Recursive aliases allow you to reference other aliases within a new alias using the `@` syntax. This enables building complex queries from simpler building blocks.

### Basic Recursive Reference

**Syntax**: `@alias_name`

**Example**:
```bash
# Create base alias
ticketr alias create my-work "assignee = currentUser() AND resolution = Unresolved"

# Reference it in another alias
ticketr alias create urgent-work "@my-work AND priority = High"
```

When `urgent-work` is expanded:
```
(assignee = currentUser() AND resolution = Unresolved) AND priority = High
```

### Chaining Multiple References

You can reference multiple aliases in a single query:

```bash
# Base aliases
ticketr alias create my-work "assignee = currentUser()"
ticketr alias create in-progress "status IN ('In Progress', 'In Review')"

# Composite alias
ticketr alias create my-active "@my-work AND @in-progress"
```

Expands to:
```
(assignee = currentUser()) AND (status IN ('In Progress', 'In Review'))
```

### Multi-Level Nesting

Aliases can reference aliases that themselves reference other aliases:

```bash
# Level 1
ticketr alias create unresolved "resolution = Unresolved"

# Level 2 (references level 1)
ticketr alias create my-work "@unresolved AND assignee = currentUser()"

# Level 3 (references level 2)
ticketr alias create urgent-mine "@my-work AND priority = High"
```

### Circular Reference Detection

Ticketr automatically detects and prevents circular references:

```bash
# This will fail
ticketr alias create a "@b AND status = Open"
ticketr alias create b "@a AND type = Bug"
```

**Error**: Circular reference detected when expanding alias.

**Prevention**: The service tracks visited aliases during expansion and returns an error if a cycle is detected.

---

## CLI Command Reference

### List Aliases

Display all available aliases for the current workspace:

```bash
ticketr alias list
```

**Output**:
```
NAME        TYPE      DESCRIPTION              JQL
----        ----      -----------              ---
mine        system    -                        assignee = currentUser() AND resolution = Un...
sprint      system    -                        sprint in openSprints()
blocked     system    -                        status = Blocked OR labels in (blocked)
my-bugs     user      -                        assignee = currentUser() AND type = Bug
critical    global    High priority tickets    priority = Critical

Use 'ticketr alias show <name>' to see the full JQL query
Use 'ticketr pull --alias <name>' to pull tickets using an alias
```

### Create Alias

Create a new JQL alias:

```bash
ticketr alias create <name> "<jql>" [flags]
```

**Flags**:
- `--description <text>`: Add a description for the alias
- `--global`: Create a global alias (default: workspace-specific)

**Examples**:
```bash
# Basic alias
ticketr alias create my-tasks "assignee = currentUser() AND type = Task"

# With description
ticketr alias create urgent "priority = Highest" --description "Highest priority items"

# Global alias
ticketr alias create critical "priority = Critical" --global
```

### Show Alias Details

Display full details of a specific alias:

```bash
ticketr alias show <name>
```

**Output**:
```
Alias: my-bugs
Type: workspace
Description: My assigned bugs

JQL Query:
  assignee = currentUser() AND type = Bug

Created: 2025-10-20 14:30:00
```

For recursive aliases, the expanded JQL is also shown:

```bash
ticketr alias show urgent-work
```

**Output**:
```
Alias: urgent-work
Type: workspace

JQL Query:
  @my-work AND priority = High

Expanded JQL:
  (assignee = currentUser() AND resolution = Unresolved) AND priority = High

Created: 2025-10-20 14:35:00
```

### Update Alias

Update an existing alias:

```bash
ticketr alias update <name> "<new-jql>" [flags]
```

**Flags**:
- `--description <text>`: Update the description

**Examples**:
```bash
# Update JQL only
ticketr alias update my-bugs "assignee = currentUser() AND type = Bug AND status != Done"

# Update with new description
ticketr alias update urgent "priority = Critical" --description "Critical priority only"
```

**Note**: Predefined system aliases (mine, sprint, blocked) cannot be updated.

### Delete Alias

Delete a user-defined alias:

```bash
ticketr alias delete <name>
```

**Example**:
```bash
ticketr alias delete my-bugs
```

**Output**:
```
Alias 'my-bugs' deleted successfully
```

**Note**: Predefined system aliases cannot be deleted.

---

## Using Aliases with Pull

### Basic Pull with Alias

Use an alias to filter tickets when pulling from Jira:

```bash
ticketr pull --alias <name> --output <file>
```

**Example**:
```bash
# Pull your assigned tickets
ticketr pull --alias mine --output my-work.md

# Pull sprint tickets
ticketr pull --alias sprint --output sprint-tickets.md
```

### Combining Alias with Other Filters

Aliases can be combined with other pull filters:

```bash
# Alias with epic filter
ticketr pull --alias mine --epic PROJ-100 --output my-epic-work.md

# Note: Cannot use both --alias and --jql together
# This will fail:
ticketr pull --alias mine --jql "type = Bug"  # Error!
```

**Restriction**: The `--alias` and `--jql` flags are mutually exclusive. Use one or the other.

### Verifying Expanded Query

Use verbose mode to see how the alias expands:

```bash
ticketr pull --alias urgent-work --output urgent.md --verbose
```

**Output**:
```
Expanded alias 'urgent-work' to JQL: (assignee = currentUser() AND resolution = Unresolved) AND priority = High
Pulling tickets from project: PROJ
Using JQL filter: (assignee = currentUser() AND resolution = Unresolved) AND priority = High
...
```

---

## Best Practices

### 1. Use Descriptive Names

Choose names that clearly communicate the alias purpose:

**Good**:
- `my-open-bugs`
- `sprint-ready`
- `blocked-critical`

**Avoid**:
- `alias1`
- `temp`
- `x`

### 2. Document with Descriptions

Always add descriptions for team-shared aliases:

```bash
ticketr alias create sprint-ready "status = 'Ready for Dev' AND sprint in openSprints()" \
  --description "Sprint backlog items ready for development"
```

### 3. Start Simple, Build Complex

Create basic building blocks, then compose them:

```bash
# Building blocks
ticketr alias create my-work "assignee = currentUser()"
ticketr alias create open-items "resolution = Unresolved"
ticketr alias create high-priority "priority IN (High, Highest)"

# Composite aliases
ticketr alias create my-urgent "@my-work AND @open-items AND @high-priority"
```

### 4. Use Global Aliases for Common Patterns

Create global aliases for queries used across all projects:

```bash
ticketr alias create recent-updates "updated >= -7d" --global
ticketr alias create critical "priority = Critical" --global
```

### 5. Namespace Team Aliases

For large teams, use prefixes to organize aliases:

```bash
# Backend team
ticketr alias create be-bugs "component = Backend AND type = Bug"
ticketr alias create be-sprint "component = Backend AND sprint in openSprints()"

# Frontend team
ticketr alias create fe-bugs "component = Frontend AND type = Bug"
ticketr alias create fe-sprint "component = Frontend AND sprint in openSprints()"
```

### 6. Test Aliases Before Sharing

Always verify your alias expands correctly:

```bash
# Create the alias
ticketr alias create test-alias "@my-work AND @sprint"

# Check expansion
ticketr alias show test-alias

# Test with pull (dry-run with verbose)
ticketr pull --alias test-alias --output test.md --verbose
```

### 7. Keep Aliases Current

Review and update aliases as project needs evolve:

```bash
# Update outdated aliases
ticketr alias update old-sprint "sprint = 'Sprint 25'" \
  --description "Updated for current sprint"
```

### 8. Clean Up Unused Aliases

Delete aliases that are no longer relevant:

```bash
# Remove obsolete aliases
ticketr alias delete old-project-filter
```

---

## Troubleshooting

### Alias Not Found

**Error**: `alias 'my-alias' not found`

**Possible Causes**:
1. Alias doesn't exist in current workspace
2. Typo in alias name
3. Alias is in a different workspace

**Solutions**:
```bash
# List all available aliases
ticketr alias list

# Check if you're in the correct workspace
ticketr workspace current

# Create the alias if it doesn't exist
ticketr alias create my-alias "your JQL here"
```

### Circular Reference Detected

**Error**: `circular reference detected when expanding alias`

**Cause**: Alias references itself directly or indirectly through a chain.

**Example Problem**:
```bash
ticketr alias create a "@b"
ticketr alias create b "@a"
```

**Solution**: Review your alias chain and break the cycle:
```bash
# Fix by removing circular reference
ticketr alias update a "assignee = currentUser()"  # No longer references @b
```

### Cannot Modify Predefined Alias

**Error**: `cannot update predefined alias 'mine'`

**Cause**: Attempting to modify system aliases (mine, sprint, blocked).

**Solution**: Create a new alias with a different name:
```bash
# Instead of updating 'mine', create a custom version
ticketr alias create my-mine "assignee = currentUser() AND resolution = Unresolved AND status = 'In Progress'"
```

### Invalid Alias Name

**Error**: `alias name must contain only alphanumeric characters, hyphens, and underscores`

**Cause**: Alias name contains invalid characters.

**Examples of Invalid Names**:
- `my bugs` (space)
- `my@bugs` (@ symbol)
- `my.bugs` (dot)

**Solution**: Use only letters, numbers, hyphens, and underscores:
```bash
ticketr alias create my-bugs "..."      # Valid
ticketr alias create my_bugs "..."      # Valid
ticketr alias create myBugs "..."       # Valid
```

### Cannot Use Both --alias and --jql

**Error**: `Cannot use both --jql and --alias flags. Please use one or the other.`

**Cause**: Attempting to use both `--alias` and `--jql` in a pull command.

**Solution**: Choose one filtering method:
```bash
# Use alias
ticketr pull --alias mine --output my-work.md

# OR use JQL directly
ticketr pull --jql "assignee = currentUser()" --output my-work.md

# To combine alias with custom JQL, create a new alias:
ticketr alias create mine-plus "@mine AND type = Bug"
ticketr pull --alias mine-plus --output output.md
```

### Alias Expansion Failed

**Error**: `failed to expand alias reference '@unknown-alias'`

**Cause**: Referenced alias doesn't exist.

**Solution**: Create the missing referenced alias or fix the reference:
```bash
# Show the problematic alias
ticketr alias show my-alias

# If it references @unknown-alias, either:
# 1. Create the missing alias
ticketr alias create unknown-alias "your JQL here"

# 2. Or update the alias to remove the reference
ticketr alias update my-alias "corrected JQL without @reference"
```

### Workspace Issues

**Error**: `failed to get current workspace`

**Cause**: No workspace is configured or current.

**Solution**:
```bash
# Check current workspace
ticketr workspace current

# Set a default workspace
ticketr workspace set-default <workspace-name>

# Or create aliases as global
ticketr alias create my-alias "JQL here" --global
```

### JQL Query Too Long

**Error**: `JQL query must be 2000 characters or less`

**Cause**: JQL query exceeds 2000 character limit.

**Solution**: Break complex queries into smaller reusable aliases:
```bash
# Instead of one huge query
ticketr alias create huge-query "very long JQL AND more JQL AND even more..."

# Break it down
ticketr alias create filter1 "first part of query"
ticketr alias create filter2 "second part of query"
ticketr alias create combined "@filter1 AND @filter2"
```

---

## Examples

### Example 1: Daily Standup Workflow

```bash
# Create alias for your in-progress work
ticketr alias create my-active "@mine AND status = 'In Progress'"

# Pull daily
ticketr pull --alias my-active --output standup.md
```

### Example 2: Sprint Planning

```bash
# Create sprint-ready alias
ticketr alias create sprint-ready "status = 'Ready for Dev' AND sprint in openSprints()"

# Create sprint backlog alias
ticketr alias create sprint-backlog "@sprint AND status != Done"

# Pull for planning meeting
ticketr pull --alias sprint-backlog --output sprint-plan.md
```

### Example 3: Bug Triage

```bash
# Create base bug filter
ticketr alias create all-bugs "type = Bug"

# Create priority-based aliases
ticketr alias create critical-bugs "@all-bugs AND priority = Critical"
ticketr alias create unassigned-bugs "@all-bugs AND assignee is EMPTY"

# Pull for triage
ticketr pull --alias critical-bugs --output triage-critical.md
ticketr pull --alias unassigned-bugs --output triage-unassigned.md
```

### Example 4: Team Coordination

```bash
# Backend team aliases
ticketr alias create be-work "component = Backend"
ticketr alias create be-blocked "@be-work AND (status = Blocked OR labels in (blocked))"

# Frontend team aliases
ticketr alias create fe-work "component = Frontend"
ticketr alias create fe-blocked "@fe-work AND (status = Blocked OR labels in (blocked))"

# Pull team-specific blockers
ticketr pull --alias be-blocked --output backend-blockers.md
ticketr pull --alias fe-blocked --output frontend-blockers.md
```

### Example 5: Release Management

```bash
# Create version-based aliases
ticketr alias create v2-0 "fixVersion = '2.0.0'"
ticketr alias create v2-0-open "@v2-0 AND resolution = Unresolved"
ticketr alias create v2-0-critical "@v2-0-open AND priority = Critical"

# Track release progress
ticketr pull --alias v2-0-open --output release-remaining.md
ticketr pull --alias v2-0-critical --output release-blockers.md
```

---

## Related Documentation

- [README.md](../../README.md) - Main user guide with JQL aliases quick start
- [WORKFLOW.md](../WORKFLOW.md) - End-to-end workflow examples
- [workspace-management-guide.md](../workspace-management-guide.md) - Multi-workspace configuration
- [TROUBLESHOOTING.md](../TROUBLESHOOTING.md) - General troubleshooting guide

---

## Technical Details

### Storage

Aliases are stored in the SQLite database at:
- **Linux/Unix**: `~/.local/share/ticketr/ticketr.db`
- **macOS**: `~/Library/Application Support/ticketr/ticketr.db`
- **Windows**: `%LOCALAPPDATA%\ticketr\ticketr.db`

### Database Schema

```sql
CREATE TABLE jql_aliases (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    jql TEXT NOT NULL,
    description TEXT,
    is_predefined INTEGER NOT NULL DEFAULT 0,
    workspace_id TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(name, workspace_id)
);
```

### Expansion Algorithm

1. Retrieve alias from repository (user-defined or predefined)
2. Parse JQL string for `@alias_name` patterns
3. Recursively expand each reference
4. Track visited aliases to detect circular references
5. Replace `@alias_name` with `(expanded_jql)`
6. Return fully expanded JQL string

### Performance Considerations

- **Alias expansion**: O(n) where n is the number of nested references
- **Circular detection**: O(n) with memoization
- **Database lookups**: Indexed by `(name, workspace_id)` for fast retrieval
- **Max nesting depth**: No hard limit, but circular reference detection prevents infinite loops

---

**Version:** 1.0
**Last Updated:** 2025-10-20
**Feedback**: Report issues at [GitHub Issues](https://github.com/karolswdev/ticketr/issues)
