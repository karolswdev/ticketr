# Ticketr Template Examples

This directory contains example YAML templates for creating ticket hierarchies with Ticketr.

## Available Templates

### feature.yaml
Creates a full feature hierarchy with epic, frontend/backend/documentation stories, and tasks.

**Variables**:
- `Name`: Feature name
- `Actor`: User role (e.g., "user", "admin")
- `Goal`: What the user wants to accomplish
- `Benefit`: Why this is valuable

**Example**:
```bash
ticketr template apply examples/templates/feature.yaml
# When prompted:
# Name: Authentication
# Actor: user
# Goal: to log in securely
# Benefit: I can access my account
```

### bug-investigation.yaml
Creates a bug investigation ticket with detailed reproduction steps and investigation tasks.

**Variables**:
- `Title`: Bug title
- `Reporter`: Who reported the bug
- `Severity`: Bug severity (Critical, High, Medium, Low)
- `Environment`: Where the bug occurs
- `Steps`: Steps to reproduce
- `Expected`: Expected behavior
- `Actual`: Actual behavior

**Example**:
```bash
ticketr template apply examples/templates/bug-investigation.yaml
```

### spike.yaml
Creates a research spike epic with investigation tasks.

**Variables**:
- `Topic`: Research topic
- `Goal`: What to learn or prove
- `Timebox`: Time limit for spike

**Example**:
```bash
ticketr template apply examples/templates/spike.yaml
# When prompted:
# Topic: GraphQL Performance
# Goal: Determine if GraphQL improves API performance
# Timebox: 3 days
```

## Creating Custom Templates

Templates are YAML files with the following structure:

```yaml
name: template-name
structure:
  epic:                    # Optional
    title: "Epic Title"
    description: "Description with {{.Variable}}"
    tasks:                 # Optional
      - "Task 1"
      - "Task 2"
  stories:                 # At least one required if no epic
    - title: "Story Title with {{.Variable}}"
      description: "Story description"
      tasks:               # Optional
        - "Task 1"
        - "Task 2"
```

### Variable Substitution

Use `{{.VariableName}}` syntax in any text field. Variables are:
- Case-sensitive
- Alphanumeric only (no spaces or special characters)
- Prompted interactively when applying template

### Template Location

Templates can be stored:
- **Project-specific**: `examples/templates/` (this directory)
- **Global**: `~/.local/share/ticketr/templates/` (Linux/macOS)
- **Global**: `%LOCALAPPDATA%\ticketr\templates\` (Windows)

Use `ticketr template list` to see all available templates.

## Tips

- Use `--dry-run` flag to preview without creating tickets
- Use `ticketr template validate <file>` to check template syntax
- Keep templates in version control for team consistency
- Use descriptive variable names for clarity
