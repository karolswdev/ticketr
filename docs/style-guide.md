# Documentation Style Guide

**Version:** 1.0
**Last Updated:** October 16, 2025
**Scope:** All documentation in the ticketr project

## Purpose

This guide establishes standards for writing, formatting, and maintaining documentation in the ticketr project. Consistent documentation improves readability, maintainability, and onboarding experience.

---

## Markdown Formatting Standards

### File Naming

**Convention:** Use lowercase kebab-case with `.md` extension

**Examples:**
- ✅ `migration-guide.md`
- ✅ `integration-testing-guide.md`
- ✅ `state-management.md`
- ❌ `MigrationGuide.md`
- ❌ `integration_testing_guide.md`
- ❌ `StateManagement.MD`

**Rationale:** Consistent naming improves discoverability and cross-platform compatibility.

---

### Heading Hierarchy

**Rules:**
1. Use exactly one `#` (H1) heading per file for the title
2. Use `##` (H2) for main sections
3. Use `###` (H3) for subsections
4. Use `####` (H4) for sub-subsections
5. Avoid headings deeper than H4 (readability)
6. Do not skip heading levels (e.g., H1 → H3)

**Example Structure:**
```markdown
# Document Title

## Section 1

### Subsection 1.1

#### Detail 1.1.1

### Subsection 1.2

## Section 2
```

**Rationale:** Proper hierarchy enables automatic table of contents generation and improves screen reader accessibility.

---

### Code Blocks

**Requirement:** Always specify language hint for syntax highlighting

**Examples:**
```markdown
```bash
ticketr push tickets.md
\```

```go
func calculateHash(ticket Ticket) string {
    // implementation
}
\```

```json
{
  "field": "value"
}
\```

```yaml
field_mappings:
  "Sprint": "customfield_10020"
\```
```

**Supported Languages:**
- `bash` - Shell commands
- `go` - Go code
- `json` - JSON data
- `yaml` - YAML configuration
- `markdown` - Markdown examples
- `text` - Plain text output

**Rationale:** Syntax highlighting improves readability and helps users distinguish code from prose.

---

### Lists

#### Unordered Lists

**Use:** For items without inherent order

**Formatting:**
- Use `-` (hyphen) as bullet character
- Indent nested items with 2 spaces
- Add blank line before/after list

**Example:**
```markdown
Available commands:

- push - Upload tickets to Jira
- pull - Download tickets from Jira
  - Supports JQL queries
  - Fetches subtasks automatically
- schema - Discover field mappings
```

#### Ordered Lists

**Use:** For sequential steps or ranked items

**Formatting:**
- Use `1.`, `2.`, `3.` (numbers with periods)
- Indent nested items with 3 spaces
- Add blank line before/after list

**Example:**
```markdown
Migration steps:

1. Preview changes with dry-run mode
2. Backup your files
3. Apply migration with --write flag
4. Verify results
```

#### Task Lists

**Use:** For checklists or tracking items

**Formatting:**
- Use `- [ ]` for incomplete items
- Use `- [x]` for completed items

**Example:**
```markdown
Pre-release checklist:

- [x] All tests passing
- [x] Documentation updated
- [ ] Release notes written
- [ ] Version bumped
```

**Rationale:** Consistent list formatting improves scannability and GitHub rendering.

---

### Links

#### Internal Links (Preferred)

**Use:** Relative paths for files within the repository

**Format:**
```markdown
See [Migration Guide](docs/migration-guide.md) for details.
See [development/REQUIREMENTS.md](development/REQUIREMENTS.md) for specifications.
See [Field Inheritance](README.md#field-inheritance) for examples.
```

**Best Practices:**
- Use relative paths from current file location
- Include anchor links for specific sections (`#heading-name`)
- Use descriptive link text (not "click here")

#### External Links

**Use:** Absolute URLs for external resources

**Format:**
```markdown
Built with [Go](https://golang.org/) programming language.
See [Jira REST API](https://developer.atlassian.com/cloud/jira/platform/rest/v2/) documentation.
```

**Best Practices:**
- Use HTTPS URLs when available
- Include protocol (https://)
- Test links periodically

#### Reference-Style Links

**Use:** For frequently referenced URLs or long URLs

**Format:**
```markdown
See the [official documentation][jira-api] for more details.

[jira-api]: https://developer.atlassian.com/cloud/jira/platform/rest/v2/
```

**Rationale:** Relative links prevent broken references when files move. Descriptive link text improves accessibility.

---

### Tables

**Formatting Rules:**
1. Use `|` to separate columns
2. Use `:---` for left-aligned columns
3. Use `:---:` for center-aligned columns
4. Use `---:` for right-aligned columns
5. Add header separator row
6. Align columns for readability (optional but recommended)

**Example:**
```markdown
| Command | Description | Example |
|:--------|:------------|:--------|
| push    | Upload tickets to Jira | `ticketr push tickets.md` |
| pull    | Download tickets from Jira | `ticketr pull --project PROJ` |
| schema  | Discover field mappings | `ticketr schema > .ticketr.yaml` |
```

**Complex Tables:**
For tables with many columns or complex content, consider:
- Breaking into multiple smaller tables
- Using definition lists instead
- Creating a separate reference page

**Rationale:** Well-formatted tables improve readability in both Markdown source and rendered output.

---

### Horizontal Rules

**Use:** To separate major sections or content blocks

**Format:** Use three hyphens with blank lines before and after

```markdown
## Section 1

Content here.

---

## Section 2

More content.
```

**When to Use:**
- Between major document sections
- Before/after large code blocks
- To separate examples from explanatory text

**When NOT to Use:**
- Between every paragraph (overuse reduces effectiveness)
- Within lists or tables
- In place of proper heading hierarchy

**Rationale:** Horizontal rules provide visual breaks without disrupting document hierarchy.

---

## Tone and Voice Guidelines

### General Principles

**Be:**
- **Technical but accessible** - Use precise terminology, but explain jargon
- **Concise** - Respect reader's time, avoid verbosity
- **Helpful** - Anticipate questions, provide context
- **Active voice** - "Run the command" not "The command should be run"
- **Direct** - "Use `--force` to override" not "You might want to consider using `--force`"

**Avoid:**
- Excessive exclamation marks (use sparingly for genuine emphasis)
- Emojis in technical documentation (acceptable in README intro only)
- Condescending language ("simply", "just", "obviously")
- Passive voice when possible
- Unnecessary qualifiers ("very", "really", "quite")

### Examples

**Good:**
```
Run the quality script before creating a pull request:
\```bash
bash scripts/quality.sh
\```
```

**Avoid:**
```
You should probably really make sure to run the quality script before you
create a pull request (it's super important!):
\```bash
bash scripts/quality.sh
\```
```

---

### Addressing the Reader

**Use:** Second person ("you") for user-facing documentation

**Examples:**
- "You can configure field mappings in `.ticketr.yaml`"
- "Run `ticketr push` to upload your tickets"
- "See the Migration Guide for step-by-step instructions"

**Exception:** Technical architecture docs may use third person or passive voice when describing system behavior.

---

### Imperative vs. Indicative

**Imperative (Instructions):**
Use for commands, procedures, and how-to guides:

```markdown
## Installation

1. Clone the repository
2. Install dependencies
3. Build the binary
```

**Indicative (Descriptions):**
Use for explanations, concepts, and reference material:

```markdown
## Architecture

Ticketr follows the Hexagonal Architecture pattern. The core business
logic is isolated from external dependencies through ports and adapters.
```

**Rationale:** Imperative mood creates clear, actionable instructions. Indicative mood provides informative descriptions.

---

## Cross-Referencing Standards

### Bidirectional Links

**Best Practice:** When referencing another document, ensure the referenced document links back if relevant.

**Example:**

In `migration-guide.md`:
```markdown
See [development/REQUIREMENTS.md](development/REQUIREMENTS.md) for schema specification (PROD-201).
```

In `development/REQUIREMENTS.md`:
```markdown
**Migration Guide:** See [docs/migration-guide.md](docs/migration-guide.md)
```

### Section Anchors

**Format:** Use lowercase with hyphens for section anchors

**Example:**
```markdown
See [Field Inheritance](#field-inheritance) section below.

...

## Field Inheritance

Content here.
```

**Auto-Generated Anchors:**
- "Field Inheritance" → `#field-inheritance`
- "Pre-Commit Checklist" → `#pre-commit-checklist`
- "v2.0 Breaking Changes" → `#v20-breaking-changes`

---

### Related Documentation Sections

**Recommendation:** Include "Related Documentation" or "See Also" section at end of document

**Example:**
```markdown
## Related Documentation

- [README.md](README.md) - User guide and quick start
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines
- [docs/ci.md](docs/ci.md) - CI/CD pipeline details
```

**Rationale:** Cross-references help users discover related content and understand the documentation structure.

---

## Example Quality Standards

### Code Examples

**Requirements:**
1. Complete and runnable
2. Include expected output when helpful
3. Add comments for complex logic
4. Use realistic variable names
5. Show both success and error cases

**Example:**

```markdown
### Pushing Tickets with Validation Errors

If validation errors are detected, the push fails by default:

\```bash
$ ticketr push tickets.md

❌ Validation errors found:
  - Line 12: Missing required field "Type"

Fix these issues before pushing to JIRA.
\```

To continue despite errors, use the --force-partial-upload flag:

\```bash
$ ticketr push tickets.md --force-partial-upload

⚠️  Validation warnings (processing will continue):
  - Line 12: Missing required field "Type"

=== Summary ===
Tickets created: 2
Errors: 1
\```
```

**Rationale:** Complete examples with expected output help users verify correct behavior.

---

### Command Examples

**Format:**
- Show command prompt (`$` or `#`)
- Include flags and arguments
- Show output when relevant
- Explain non-obvious options

**Example:**

```markdown
Pull tickets from a specific epic:

\```bash
$ ticketr pull --epic PROJ-100 --output epic-tickets.md

Successfully pulled 15 tickets with 42 subtasks
\```

The --epic flag filters tickets belonging to the specified epic.
```

**Rationale:** Clear command examples reduce user errors and improve self-service success rates.

---

## Line Length Recommendations

**Target:** ~80-120 characters per line for prose
**Hard Limit:** 120 characters (enforced for code, recommended for prose)

**Exceptions:**
- Long URLs (do not break)
- Code blocks (follow language-specific conventions)
- Tables (may exceed for readability)
- Literal command output

**Rationale:**
- Improves readability in narrow viewports
- Easier diff reviews in version control
- Better compatibility with various editors

**Tool Support:** Most editors support soft wrapping for viewing while maintaining source line length.

---

## File Organization

### Front Matter

**Recommended sections for guides:**

```markdown
# Document Title

**Version:** 1.0
**Last Updated:** 2025-10-16
**Scope:** Brief description of document scope

## Purpose

One or two paragraphs explaining why this document exists and who should read it.

---

[Main content begins here]
```

### Table of Contents

**When to Include:**
- Documents longer than 3 pages
- Complex documents with many sections
- Reference documentation

**Format:** Use Markdown list with section links

**Example:**
```markdown
## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Push Command](#push-command)
  - [Pull Command](#pull-command)
- [Troubleshooting](#troubleshooting)
```

**Note:** GitHub automatically generates TOC for README.md files. Manual TOC may be redundant.

---

### Section Ordering

**Typical Structure:**

1. **Title and metadata** (version, date, scope)
2. **Purpose/Overview** - Why this document exists
3. **Prerequisites** - What reader needs before proceeding
4. **Main Content** - Ordered logically (simple → complex, general → specific)
5. **Examples** - Practical applications
6. **Troubleshooting** - Common issues and solutions
7. **Related Documentation** - Cross-references
8. **Appendix** - Reference material (optional)

**Rationale:** Consistent structure helps readers navigate unfamiliar documents quickly.

---

## Special Formatting

### Admonitions (Notes, Warnings, Tips)

**Format:** Use blockquotes with emoji or bold prefix

**Note:**
```markdown
**Note:** This feature requires Go 1.22 or later.
```

**Warning:**
```markdown
**Warning:** Using `--force` will permanently overwrite local changes with remote data.
```

**Tip:**
```markdown
**Tip:** Use `ticketr schema` to discover available custom fields in your Jira instance.
```

**Rationale:** Visual distinction helps critical information stand out.

---

### Inline Code

**Use:** Backticks for code, commands, filenames, and technical terms

**Examples:**
- Command: `ticketr push tickets.md`
- Filename: `.ticketr.state`
- Function: `calculateFinalFields()`
- Variable: `JIRA_URL`
- Field: `customfield_10020`

**When NOT to Use:**
- Regular words that happen to be technical (Git, Jira, Markdown)
- Complete sentences
- Emphasis (use **bold** or *italic* instead)

---

### Emphasis

**Bold (`**text**`):**
- Important terms on first use
- Strong emphasis
- Section labels in lists

**Italic (`*text*`):**
- Subtle emphasis
- Technical terms from other domains
- Book/document titles

**Example:**
```markdown
The **state file** tracks ticket hashes. See *State Management* in the
architecture documentation for details.
```

**Rationale:** Consistent emphasis improves scannability without overusing visual noise.

---

## Documentation Maintenance

### Version Control

**Best Practices:**
1. Commit documentation changes with related code changes
2. Use descriptive commit messages for doc-only changes
3. Review diffs carefully (typos are easy to miss)
4. Update "Last Updated" date when making substantive changes

**Commit Message Examples:**
```
docs: update migration guide with batch processing examples
docs: fix broken link in CONTRIBUTING.md
docs: add troubleshooting section to state-management.md
```

---

### Keeping Documentation Current

**When Code Changes:**
- [ ] Update affected documentation in same commit/PR
- [ ] Check cross-references are still valid
- [ ] Update code examples if affected
- [ ] Verify command syntax is current

**Periodic Reviews:**
- Review all docs quarterly for accuracy
- Update version numbers and dates
- Check for broken links
- Identify outdated examples

**Deprecation:**
- Mark deprecated content clearly
- Provide migration path
- Archive to `docs/legacy/` or `docs/history/` when fully obsolete

---

### Documentation Checklist

Before committing documentation changes:

- [ ] Spell check completed (use editor's built-in tools)
- [ ] Links tested (internal and external)
- [ ] Code examples tested and verified
- [ ] Markdown syntax validated (use linter if available)
- [ ] Headings follow hierarchy rules
- [ ] Code blocks have language hints
- [ ] Tables formatted correctly
- [ ] Tone and voice consistent with project standards
- [ ] Cross-references updated bidirectionally
- [ ] Version/date updated if substantive changes

---

## Tools and Resources

### Recommended Editors

**Markdown-Aware Editors:**
- VS Code with Markdown extensions
- IntelliJ IDEA with Markdown plugin
- Typora (WYSIWYG)
- MacDown (macOS)

**Helpful Extensions/Plugins:**
- Markdown linters (e.g., markdownlint)
- Spell checkers
- Link validators
- Table formatters

---

### Validation Tools

**Markdown Linters:**
```bash
# markdownlint-cli (Node.js)
npm install -g markdownlint-cli
markdownlint '**/*.md'

# mdl (Ruby)
gem install mdl
mdl docs/
```

**Link Checkers:**
```bash
# markdown-link-check (Node.js)
npm install -g markdown-link-check
markdown-link-check README.md
```

**Spell Checkers:**
- Built-in editor spell check
- `aspell` or `hunspell` (CLI)
- VS Code spell checker extensions

---

### Preview Tools

**Local Preview:**
- GitHub-flavored Markdown preview in most editors
- `grip` (renders GitHub-style locally)
- Static site generators (Jekyll, Hugo)

**GitHub Preview:**
- GitHub automatically renders `.md` files
- Preview tab in PR interface
- GitHub Pages for published documentation

---

## Style Guide Governance

### Updates to This Guide

**Process:**
1. Propose changes via GitHub issue or PR
2. Discuss with maintainers
3. Update style guide
4. Announce changes to contributors
5. Apply retroactively to existing docs as time permits

### Exceptions

**When to Deviate:**
- External documentation requires different format (e.g., API docs)
- Auto-generated documentation (e.g., godoc)
- Legacy documentation in `docs/legacy/` or `docs/history/` (preserve as-is)

**How to Request Exception:**
1. Document reason in PR description
2. Get maintainer approval
3. Add note in document explaining exception

---

## Examples of Well-Formatted Documentation

**Reference Examples in This Repository:**

Excellent examples demonstrating these standards:

- **README.md** - User-facing guide with clear structure, code examples, cross-references
- **docs/WORKFLOW.md** - End-to-end workflow guide with step-by-step instructions
- **docs/ci.md** - Comprehensive technical documentation with troubleshooting
- **docs/migration-guide.md** - Procedural guide with clear steps and examples
- **ARCHITECTURE.md** - Technical reference with diagrams, tables, and cross-references

Study these documents when writing new documentation.

---

## Summary

**Key Takeaways:**

1. **Consistency is paramount** - Follow established patterns
2. **Clarity over brevity** - Be concise, but prioritize understanding
3. **Examples matter** - Show, don't just tell
4. **Maintain actively** - Keep docs current with code changes
5. **Cross-reference generously** - Help users discover related content
6. **Test everything** - Verify examples work, links resolve, commands execute

**Questions?**

- Refer to existing documentation for patterns
- Ask maintainers for clarification
- Propose style guide updates via GitHub issue

---

**Style Guide Version:** 1.0
**Adopted:** October 16, 2025
**Next Review:** January 16, 2026
