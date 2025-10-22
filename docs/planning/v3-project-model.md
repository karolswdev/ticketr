# Ticketr v3.0: Project-Based Architecture

**Version:** 1.0
**Date:** January 2025
**Status:** Architecture Revision

---

## Core Concept: Projects as Namespaces

Ticketr "Projects" are **isolated work contexts** - like Kubernetes namespaces or Git repositories. Each project:
- Maps to a specific JIRA instance and project key
- Has its own credentials and configuration
- Maintains its own markdown files locally
- Operates independently from other projects

```
~/.config/ticketr/
├── config.yaml                 # Global settings
├── ticketr.db                  # Central database
└── projects/                   # Project configurations
    ├── backend.yaml
    ├── frontend.yaml
    └── mobile.yaml

~/work/                         # Your actual work directories
├── backend-project/
│   ├── epics/
│   │   ├── auth-system.md     # Epic with stories and tasks
│   │   └── payment-flow.md
│   └── .ticketr/              # Project-specific cache
├── frontend-app/
│   ├── tickets.md             # All tickets in one file
│   └── .ticketr/
└── mobile/
    ├── sprint-23.md
    └── .ticketr/
```

---

## Markdown Hierarchy: Epic → Story → Task

The markdown structure defines the **logical hierarchy**, independent of JIRA's issue types:

```markdown
# EPIC: [PROJ-100] Authentication System

Epic-level description and goals.

## Fields
- Epic Name: Authentication System
- Target Release: v2.0
- Epic Owner: @john.doe

## STORY: [PROJ-101] User Registration Flow

As a new user, I want to register an account...

### Acceptance Criteria
- Email validation
- Password strength requirements
- Confirmation email sent

### Fields
- Story Points: 8
- Sprint: Sprint 23

### TASK: [PROJ-201] Create registration API endpoint

Technical implementation details...

#### Fields
- Assignee: @dev.user
- Priority: High

### TASK: [PROJ-202] Design registration form UI

UI/UX implementation...

#### Fields
- Assignee: @ui.designer
- Component: Frontend

## STORY: [PROJ-102] Login and Session Management

Another story under the same epic...

### TASK: [PROJ-203] Implement JWT tokens

...
```

**Key Points:**
- `#` = EPIC level
- `##` = STORY level
- `###` = TASK level
- Each level can have its own Fields section
- Hierarchy is **defined in markdown**, mapped to JIRA

---

## Project Model

### Project Configuration

```yaml
# ~/.config/ticketr/projects/backend.yaml
name: backend
description: Backend API Services
jira:
  url: https://company.atlassian.net
  project_key: BACK
  credential_ref: backend_creds  # Keychain reference

hierarchy:
  epic_type: Epic
  story_types: [Story, Bug]
  task_type: Sub-task

paths:
  default_location: ~/work/backend-project
  epic_pattern: "epics/*.md"
  story_pattern: "stories/*.md"

defaults:
  assignee: currentUser()
  sprint: currentSprint()
  component: Backend
```

### Project Commands

```bash
# Project management
ticketr project create backend --jira-url https://... --project-key BACK
ticketr project list
ticketr project switch backend
ticketr project current
ticketr project config backend --set defaults.sprint "Sprint 24"

# Working within a project context
cd ~/work/backend-project
ticketr push epics/auth-system.md    # Uses 'backend' project config
ticketr pull --epic BACK-100         # Pulls into local markdown structure

# Or specify project explicitly from anywhere
ticketr push ~/documents/tickets.md --project frontend
ticketr pull --project mobile --jql "sprint = 'Sprint 23'"

# TUI with project context
ticketr tui --project backend        # Opens TUI for backend project
ticketr tui                          # Opens TUI with project selector
```

---

## Database Schema (Revised)

```sql
-- Projects table (like namespaces)
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    jira_url TEXT NOT NULL,
    project_key TEXT NOT NULL,
    credential_ref TEXT NOT NULL,
    config JSON,  -- Hierarchy mappings, defaults, etc.
    is_current BOOLEAN DEFAULT FALSE,
    last_accessed TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tickets with project isolation
CREATE TABLE tickets (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    jira_id TEXT,
    parent_id TEXT REFERENCES tickets(id),
    level INTEGER CHECK(level IN (0, 1, 2)), -- 0=Epic, 1=Story, 2=Task
    issue_type TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    fields JSON,
    acceptance_criteria JSON,
    markdown_path TEXT,  -- Source file path
    markdown_line INTEGER, -- Line number in source
    local_hash TEXT,
    remote_hash TEXT,
    sync_status TEXT DEFAULT 'new',
    last_synced TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, jira_id)
);

-- Project-specific state
CREATE TABLE project_state (
    project_id TEXT PRIMARY KEY REFERENCES projects(id),
    last_push TIMESTAMP,
    last_pull TIMESTAMP,
    cached_schema JSON,
    statistics JSON  -- Ticket counts, velocity, etc.
);

-- Multi-project sync history
CREATE TABLE sync_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id TEXT REFERENCES projects(id),
    operation TEXT,
    file_path TEXT,
    tickets_affected INTEGER,
    status TEXT,
    details JSON,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ticket_project ON tickets(project_id);
CREATE INDEX idx_ticket_hierarchy ON tickets(project_id, parent_id, level);
```

---

## How Projects Work

### 1. Project Initialization

```bash
$ ticketr project create backend \
    --jira-url https://company.atlassian.net \
    --project-key BACK \
    --path ~/work/backend

Creating project 'backend'...
✓ Validated JIRA connection
✓ Stored credentials in system keychain
✓ Discovered issue types: Epic, Story, Bug, Task, Sub-task
✓ Created project configuration
✓ Set as current project

$ cd ~/work/backend
$ ticketr pull --epic BACK-100

Pulling from project 'backend' (BACK)...
✓ Found Epic: BACK-100 "Authentication System"
✓ Found 3 Stories: BACK-101, BACK-102, BACK-103
✓ Found 8 Sub-tasks
✓ Created: epics/authentication-system.md
```

### 2. Project Isolation

Each project is **completely isolated**:

```bash
# In backend project directory
$ ticketr push epics/auth.md
Pushing to 'backend' project (BACK)...

# In frontend project directory
$ ticketr push components/nav.md
Pushing to 'frontend' project (FRONT)...

# Projects don't interfere with each other
$ ticketr project list
* backend   (current) - BACK - Last sync: 2 min ago
  frontend            - FRONT - Last sync: 1 hour ago
  mobile              - MOB - Last sync: 3 days ago
```

### 3. Hierarchy Preservation

The Epic → Story → Task hierarchy is **preserved in markdown**:

```bash
$ ticketr pull --epic BACK-100 --output epics/auth.md

# Creates structured markdown:
# EPIC: [BACK-100] Authentication System
## STORY: [BACK-101] User Registration
### TASK: [BACK-201] API Endpoint
### TASK: [BACK-202] Frontend Form
## STORY: [BACK-102] Login Flow
### TASK: [BACK-203] JWT Implementation
```

### 4. Smart Context Detection

Ticketr automatically detects project context:

```go
func DetectProject(path string) (*Project, error) {
    // 1. Check for .ticketr/project.yaml in current or parent dirs
    // 2. Check if path is under a known project directory
    // 3. Fall back to 'current' project
    // 4. Prompt user if ambiguous
}
```

---

## TUI with Project Context

The TUI becomes project-aware with a namespace selector:

```
┌─────────────────────────────────────────────────────────────┐
│ Ticketr v3.0 | Project: [backend ▼] | JIRA: Connected      │
├──────────────┬──────────────────────────────────────────────┤
│              │                                              │
│ Projects     │  BACK-100: Authentication System     [Epic]  │
│ ● backend    │  ├── BACK-101: User Registration    [Story] │
│ ○ frontend   │  │   ├── BACK-201: API Endpoint     [Task]  │
│ ○ mobile     │  │   └── BACK-202: Frontend Form    [Task]  │
│ + Add...     │  └── BACK-102: Login Flow           [Story] │
│              │      ├── BACK-203: JWT Tokens       [Task]  │
│ Quick Filter │      └── BACK-204: Session Mgmt     [Task]  │
│ ○ My Tasks   │                                              │
│ ○ Sprint 23  │  ┌─────────────────────────────────────┐    │
│ ○ Epics Only │  │ Title: Add OAuth support             │    │
│              │  │ Parent: BACK-102 (Login Flow)       │    │
│              │  │ Type: [Sub-task ▼]                 │    │
│              │  │ Description: [....................]  │    │
├──────────────┴──────────────────────────────────────────────┤
│ [Tab]Switch Project [/]Search [c]Create [p]Push [P]Pull    │
└──────────────────────────────────────────────────────────────┘
```

---

## Benefits of Project Model

### 1. **True Multi-Instance Support**
- Work with different JIRA instances simultaneously
- Each project has its own credentials
- No configuration bleeding

### 2. **Clear Context**
- Always know which JIRA you're pushing to
- Project key is inherent to the project
- No more `JIRA_PROJECT_KEY` environment variable confusion

### 3. **Hierarchy in Markdown**
- Epic → Story → Task structure in files
- Visual hierarchy matches logical structure
- Easy to reorganize in your editor

### 4. **Scalability**
- SQLite indexes by project
- Quick project switching
- Efficient queries within project namespace

### 5. **Team Collaboration**
- Share project configs (without credentials)
- Consistent structure across team
- Git-friendly markdown files

---

## Migration from v2

```bash
$ ticketr migrate v2-to-v3

Detected v2 installation...
Found JIRA_PROJECT_KEY=PROJ in environment

Would you like to create a project 'PROJ' with these settings? [Y/n]

Creating project 'PROJ'...
✓ Imported existing tickets from .ticketr.state
✓ Maintained ticket IDs and relationships
✓ Created project configuration

Migration complete!
- Your tickets are now in project 'PROJ'
- Run 'ticketr project list' to see all projects
- Run 'ticketr pull --project PROJ' to sync
```

---

## Example Workflows

### Workflow 1: Multiple JIRA Instances

```bash
# Personal projects on personal JIRA
ticketr project create personal \
  --jira-url https://personal.atlassian.net \
  --project-key PERSONAL

# Work projects on company JIRA
ticketr project create work \
  --jira-url https://company.atlassian.net \
  --project-key WORK

# Switch contexts easily
ticketr project switch personal
ticketr push ideas.md

ticketr project switch work
ticketr push sprint-tasks.md
```

### Workflow 2: Cross-Project Dependencies

```markdown
# In backend project epic
## STORY: [BACK-150] API Authentication

### Fields
- Blocks: FRONT-200  # Cross-project reference

### Description
This must be completed before frontend can implement...
```

### Workflow 3: Project Templates

```bash
# Create project from template
ticketr project create mobile --template ios-app

# Template includes:
# - Standard epic structure
# - Default field mappings
# - Common workflows
```

---

## Summary

The **Project Model** gives us:

1. **Namespace isolation** - Like Kubernetes namespaces
2. **Multi-instance support** - Different JIRAs, different credentials
3. **Markdown hierarchy** - Epic→Story→Task in your files
4. **Local-first** - Markdown files live in your project directories
5. **Central management** - SQLite database tracks everything
6. **Smart context** - Auto-detects which project you're in

This is **much cleaner** than workspaces - it's more like how `git` repositories work, but with central tracking and management.