# Migration Guide: Legacy STORY Format to TICKET Format

## Overview

Ticketr v2.0 introduces the generic `# TICKET:` schema, replacing the legacy `# STORY:` format. This guide explains how to migrate your existing Markdown files.

## Why the Change?

The generic `# TICKET:` schema supports multiple Jira issue types (stories, bugs, tasks, epics) while maintaining backward compatibility with hierarchical validation. The legacy `# STORY:` format was overly specific and didn't align with Jira's flexible issue type system.

**Reference:** REQUIREMENTS-v2.md (PROD-201)

## Migration Methods

### Option 1: Automated Migration (Recommended)

Use the `ticketr migrate` command to automatically convert files:

#### Dry-Run Mode (Preview Changes)
```bash
# Preview changes for a single file
ticketr migrate path/to/your-story.md

# Preview changes for multiple files
ticketr migrate examples/*.md
```

#### Write Mode (Apply Changes)
```bash
# Apply changes to a single file
ticketr migrate path/to/your-story.md --write

# Apply changes to multiple files
ticketr migrate examples/*.md --write
```

**How it works:**
- Scans file for `# STORY:` patterns
- Replaces with `# TICKET:`
- Preserves all formatting, whitespace, and indentation
- Reports number of changes made

### Option 2: Manual Migration

For simple cases, you can use find-replace in your text editor:

1. Open the Markdown file
2. Find: `# STORY:`
3. Replace with: `# TICKET:`
4. Save the file
5. Verify with: `ticketr validate path/to/file.md`

## Migration Checklist

- [ ] Identify all Markdown files using legacy `# STORY:` format
- [ ] Run `ticketr migrate <files>` in dry-run mode to preview changes
- [ ] Review the preview to ensure changes are correct
- [ ] Run `ticketr migrate <files> --write` to apply changes
- [ ] Run `ticketr validate <files>` to verify migrated files
- [ ] Test your workflow with migrated files
- [ ] Commit the changes to version control

## Common Scenarios

### Scenario 1: Single File Migration
```bash
ticketr migrate tickets/PROJ-123.md --write
```

### Scenario 2: Batch Migration (Entire Directory)
```bash
ticketr migrate tickets/*.md --write
```

### Scenario 3: Selective Migration
```bash
# Preview all files first
ticketr migrate tickets/*.md

# Migrate only specific files that need changes
ticketr migrate tickets/old-story-1.md tickets/old-story-2.md --write
```

## Troubleshooting

### Error: "Legacy '# STORY:' format detected"

If you see this error when running `ticketr push`, it means your file hasn't been migrated yet:

```
Error: Legacy '# STORY:' format detected at line 1.
Please migrate to '# TICKET:' format.
See REQUIREMENTS-v2.md (PROD-201) or use 'ticketr migrate <file>' command.
```

**Solution:** Run the migration command on the affected file:
```bash
ticketr migrate path/to/file.md --write
```

### Files Already Using # TICKET:

If you run the migrate command on files already using `# TICKET:` format, the tool will report:
```
No changes needed for: path/to/file.md
```

### Mixed Formats in One File

If a file contains both `# STORY:` and `# TICKET:` headers, the migration tool will convert all `# STORY:` instances to `# TICKET:`.

## Rollback Strategy

If you need to revert changes after migration:

1. **Using Git:**
   ```bash
   git checkout -- path/to/file.md
   ```

2. **Using Backup:**
   Before migrating, create backups:
   ```bash
   cp -r tickets/ tickets.backup/
   ```

## Schema Differences

The migration is straightforward - only the header keyword changes:

**Before (Legacy):**
```markdown
# STORY: PROJ-123 Implement user authentication
```

**After (Canonical):**
```markdown
# TICKET: PROJ-123 Implement user authentication
```

All other fields remain identical:
- Description
- Acceptance Criteria
- Tasks
- Issue Type
- Status
- Assignee
- etc.

## Validation

After migration, validate your files:

```bash
# Validate single file
ticketr validate path/to/file.md

# Validate multiple files
ticketr validate tickets/*.md
```

## Additional Resources

- **Schema Specification:** REQUIREMENTS-v2.md (PROD-201)
- **Parser Behavior:** The v2.0+ parser rejects `# STORY:` format with helpful error messages
- **Test Fixtures:** Legacy samples in `testdata/legacy_story/` demonstrate rejected formats
- **Examples:** See `examples/` directory for canonical format templates

## Support

If you encounter issues during migration:
1. Check the error message for specific guidance
2. Review REQUIREMENTS-v2.md for schema details
3. Run `ticketr help migrate` for command usage
4. Create a GitHub issue with details if problems persist
