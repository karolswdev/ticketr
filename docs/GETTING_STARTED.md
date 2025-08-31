# Getting Started with Ticketr

Welcome to Ticketr! This guide will walk you through setting up and using Ticketr for the first time.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Your First Ticket](#your-first-ticket)
- [Common Workflows](#common-workflows)
- [Tips and Best Practices](#tips-and-best-practices)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before you begin, ensure you have:

1. **Go 1.24+** installed (for building from source)
2. **JIRA Account** with API access enabled
3. **JIRA API Token** (not your password!)
   - Generate at: https://id.atlassian.com/manage-profile/security/api-tokens
4. **Project permissions** in JIRA to create and modify issues

## Installation

### Option 1: Install with Go

```bash
go install github.com/karolswdev/ticketr/cmd/ticketr@latest
```

### Option 2: Build from Source

```bash
git clone https://github.com/karolswdev/ticketr.git
cd ticketr
go build -o ticketr cmd/ticketr/main.go
sudo mv ticketr /usr/local/bin/
```

### Option 3: Use Docker

```bash
docker pull ghcr.io/karolswdev/ticketr:latest
```

Verify installation:
```bash
ticketr --version
```

## Configuration

### Step 1: Set Environment Variables

Create a `.env` file in your project directory:

```bash
# JIRA Configuration
JIRA_URL=https://yourcompany.atlassian.net
JIRA_EMAIL=your.email@company.com
JIRA_API_KEY=your-api-token-here
JIRA_PROJECT_KEY=PROJ
```

Load the environment:
```bash
source .env
```

💡 **Tip**: Add `.env` to your `.gitignore` to keep credentials safe!

### Step 2: Verify Connection

Test your configuration:
```bash
ticketr schema
```

This should display your JIRA project's available fields.

## Your First Ticket

### Step 1: Create a Ticket File

Create `my-first-ticket.md`:

```markdown
# TICKET: Implement User Login

## Description
As a user, I want to log into the application
so that I can access my personal dashboard.

## Acceptance Criteria
- Users can enter username and password
- Invalid credentials show error message
- Successful login redirects to dashboard
- Session persists for 24 hours

## Tasks
- Design login form UI
- Implement authentication API
- Add form validation
- Create session management
- Write unit tests
```

### Step 2: Push to JIRA

```bash
ticketr push my-first-ticket.md
```

Output:
```
Validating tickets...
✓ Validation passed

Pushing to JIRA...
✓ Created ticket: PROJ-123
✓ Created subtask: PROJ-124
✓ Created subtask: PROJ-125
✓ Created subtask: PROJ-126
✓ Created subtask: PROJ-127
✓ Created subtask: PROJ-128

File updated with JIRA IDs!
```

### Step 3: View Updated File

Your file now includes JIRA IDs:

```markdown
# TICKET: [PROJ-123] Implement User Login

## Description
As a user, I want to log into the application
so that I can access my personal dashboard.

## Acceptance Criteria
- Users can enter username and password
- Invalid credentials show error message
- Successful login redirects to dashboard
- Session persists for 24 hours

## Tasks
- [PROJ-124] Design login form UI
- [PROJ-125] Implement authentication API
- [PROJ-126] Add form validation
- [PROJ-127] Create session management
- [PROJ-128] Write unit tests
```

## Common Workflows

### Sprint Planning

1. **Create Sprint Backlog**
```bash
vim sprint-24.md
# Add all planned tickets
```

2. **Push to JIRA**
```bash
ticketr push sprint-24.md
```

3. **Track Progress**
Update task statuses with emojis:
```markdown
## Tasks
- [PROJ-124] ✅ Design complete
- [PROJ-125] 🚧 API in progress
- [PROJ-126] Validation pending
```

### Daily Standup

Pull latest updates from JIRA:
```bash
ticketr pull --jql "assignee=currentUser() AND sprint in openSprints()" -o my-tasks.md
```

### End of Sprint

Generate sprint report:
```bash
ticketr stats sprint-24.md
```

### Bulk Updates

1. **Pull tickets from epic**
```bash
ticketr pull --epic PROJ-100 -o epic-tickets.md
```

2. **Make bulk edits** in your editor

3. **Push updates**
```bash
ticketr push epic-tickets.md
```

## Tips and Best Practices

### Organizing Tickets

Use descriptive titles with type prefixes:
```markdown
# TICKET: [Feature] Dark Mode Support
# TICKET: [Bug] Login fails with special characters
# TICKET: [Epic] Q4 Platform Migration
```

### Using Custom Fields

Add custom fields in the Fields section:
```markdown
## Fields
Story Points: 5
Sprint: Sprint 24
Priority: High
Component: Authentication
```

### Tracking Progress

Use emojis for visual status:
- ✅ Completed
- 🚧 In Progress
- 🔄 In Review
- ⏸️  Blocked
- 🐛 Has bugs

### CI/CD Integration

Automate ticket sync in GitHub Actions.

Option A — Install CLI and run:
```yaml
jobs:
  ticketr:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Install ticketr
        run: go install github.com/karolswdev/ticketr/cmd/ticketr@latest
      - name: Push tickets
        env:
          JIRA_URL: ${{ secrets.JIRA_URL }}
          JIRA_EMAIL: ${{ secrets.JIRA_EMAIL }}
          JIRA_API_KEY: ${{ secrets.JIRA_API_KEY }}
          JIRA_PROJECT_KEY: PROJ
        run: ticketr push stories/backlog.md
```

Option B — Use the composite action from this repo:
```yaml
- name: Push tickets to JIRA
  uses: karolswdev/ticketr/.github/actions/ticketr-sync@main # Pin to a tag in production
  with:
    jira-url: ${{ secrets.JIRA_URL }}
    jira-email: ${{ secrets.JIRA_EMAIL }}
    jira-api-key: ${{ secrets.JIRA_API_KEY }}
    jira-project-key: PROJ
    command: push
    file-path: stories/backlog.md
```

### Real-time Sync

Start webhook server for automatic updates:
```bash
ticketr listen --port 8080 --path tickets.md
```

Configure JIRA webhook to point to your server.

## Troubleshooting

### Authentication Errors

**Problem**: "Error: Authentication failed"

**Solution**:
1. Verify API token (not password!)
2. Check email is correct
3. Ensure JIRA_URL includes https://
4. Test with: `curl -u email:token https://company.atlassian.net/rest/api/2/myself`

### Permission Errors

**Problem**: "Error: You do not have permission to create issues"

**Solution**:
1. Check project permissions in JIRA
2. Verify JIRA_PROJECT_KEY is correct
3. Ensure your account can create issues in the project

### Parsing Errors

**Problem**: "Error: Failed to parse ticket file"

**Solution**:
1. Use correct format: `# TICKET:` (not `# STORY:`)
2. Ensure proper Markdown formatting
3. Check for unclosed code blocks
4. Validate with: `ticketr push --dry-run file.md`

### Sync Conflicts

**Problem**: "Conflict detected! Tickets have both local and remote changes"

**Solution**:
```bash
# Keep local changes
ticketr pull --strategy=local-wins

# Use remote changes
ticketr pull --strategy=remote-wins
```

### Network Issues

**Problem**: "Error: Connection timeout"

**Solution**:
1. Check internet connection
2. Verify firewall allows HTTPS to Atlassian
3. Test JIRA access in browser
4. Try with verbose mode: `ticketr push file.md -v`

## Next Steps

- 📖 Read the [Architecture Guide](ARCHITECTURE.md) to understand how Ticketr works
- 🔧 Explore [Advanced Configuration](.ticketr.yaml) options
- 🔄 Set up [Webhook Integration](WEBHOOKS.md) for real-time sync
- 🤖 Configure [GitHub Actions](../.github/workflows/) for automation
- 📊 Use `ticketr stats` to track project metrics

## Getting Help

- 💬 [GitHub Discussions](https://github.com/karolswdev/ticketr/discussions)
- 🐛 [Report Issues](https://github.com/karolswdev/ticketr/issues)
- 📚 Project docs live in the `docs/` folder of this repo

---

Happy ticket management! 🎫
