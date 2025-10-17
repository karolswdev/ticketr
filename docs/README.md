# Ticketr Documentation

## Test Fixtures and Unsupported Samples

### Unsupported Story Schema (`testdata/unsupported_story/`)

This directory contains markdown files using the old `# STORY:` heading. These samples serve a specific purpose:

**Purpose:** Regression testing to ensure the parser correctly rejects unsupported headings with helpful error messages.

**Important Notes:**
- These files are **NOT** valid for current Ticketr operations
- They exist solely for regression testing
- The canonical schema is `# TICKET:` (see development/REQUIREMENTS.md PROD-201)
- Users encountering `# STORY:` headings should rename them to `# TICKET:` manually before using Ticketr

**Related Files:**
- Schema specification: development/REQUIREMENTS.md (PROD-201)
- Parser rejection tests: `internal/parser/parser_test.go`

---

*This documentation structure will be expanded in Milestone 10.*
