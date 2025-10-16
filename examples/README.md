# Ticketr Examples

This directory contains template files demonstrating the canonical `# TICKET:` schema for various use cases.

## Available Templates

### Quick Start Template
**File:** `quick-story.md`
**Purpose:** Basic ticket template for simple stories or tasks
**When to use:**
- Creating standalone stories
- Simple bug reports
- Quick tasks without complex hierarchy

**Key features:**
- Minimal required fields
- Clean structure
- Easy to customize

### Epic Template
**File:** `epic-template.md`
**Purpose:** Large initiatives with multiple child stories/tasks
**When to use:**
- Planning major features spanning multiple sprints
- Organizing related work items
- Managing complex projects

**Key features:**
- Epic-level metadata
- Child ticket organization
- Hierarchical structure

### Sprint Template
**File:** `sprint-template.md`
**Purpose:** Sprint planning with multiple tickets
**When to use:**
- Planning sprint backlogs
- Organizing iteration work
- Batch ticket creation

**Key features:**
- Multiple tickets in one file
- Sprint-level organization
- Consistent formatting

## Using These Templates

### 1. Copy the Template
```bash
cp examples/quick-story.md tickets/PROJ-123.md
```

### 2. Customize the Content
Edit the file to match your specific requirements:
- Update ticket ID (e.g., `PROJ-123`)
- Add description and acceptance criteria
- Set issue type, status, assignee
- Add tasks or subtasks as needed

### 3. Validate the Ticket
```bash
ticketr validate tickets/PROJ-123.md
```

### 4. Push to Jira
```bash
ticketr push tickets/PROJ-123.md
```

## Template Customization Guidelines

### Required Fields
All templates must include:
- `# TICKET:` header with ticket ID and title
- Valid issue type (Story, Bug, Task, Epic, etc.)

### Optional Fields
Customize based on your needs:
- Description
- Acceptance Criteria
- Status
- Assignee
- Sprint
- Epic Link
- Story Points
- Labels
- Custom fields

### Formatting Rules
- Use `# TICKET:` (not `# STORY:`) for the main header
- Use `## TASK:` for child tasks
- Use `**Field Name:**` for metadata fields
- Maintain consistent indentation

## Schema Reference

For complete schema specification, see:
- **REQUIREMENTS-v2.md** - PROD-201 section
- **README.md** - Quick Start and Templates sections
- **docs/migration-guide.md** - Migration from legacy format

## Creating Custom Templates

To create your own templates:

1. Start with an existing template
2. Add/remove fields based on your workflow
3. Test with `ticketr validate`
4. Save in your project's templates directory
5. Share with your team

### Example: Custom Bug Template
```markdown
# TICKET: PROJ-456 Bug title here

**Issue Type:** Bug
**Priority:** High
**Status:** Open
**Assignee:** developer@example.com

**Description:**
[Describe the bug]

**Steps to Reproduce:**
1. [Step 1]
2. [Step 2]

**Expected Behavior:**
[What should happen]

**Actual Behavior:**
[What actually happens]

**Environment:**
- OS: [e.g., Ubuntu 22.04]
- Browser: [e.g., Chrome 120]
- Version: [e.g., v1.2.3]
```

## Legacy Format Notice

Important: These templates use the canonical `# TICKET:` schema introduced in v2.0. If you have legacy files using `# STORY:` format, see the migration guide:

```bash
ticketr help migrate
# or
cat docs/migration-guide.md
```

## Additional Resources

- **CLI Help:** `ticketr help`
- **Validation:** `ticketr validate --help`
- **Migration:** `ticketr migrate --help`
- **Push/Pull:** `ticketr push --help`, `ticketr pull --help`
