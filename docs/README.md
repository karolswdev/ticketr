# Ticketr Documentation

## Test Fixtures and Legacy Samples

### Legacy Story Schema (`testdata/legacy_story/`)

This directory contains markdown files using the deprecated `# STORY:` schema format. These samples serve a specific purpose:

**Purpose:** Regression testing to ensure the parser correctly rejects legacy format with helpful error messages.

**Important Notes:**
- These files are **NOT** valid for current Ticketr operations
- They exist solely for backward compatibility testing
- The canonical schema is `# TICKET:` (see REQUIREMENTS-v2.md PROD-201)
- Users encountering `# STORY:` format should migrate to `# TICKET:` immediately

**Related Files:**
- Migration guidance: README.md "Migrating from v1.x" section
- Schema specification: REQUIREMENTS-v2.md (PROD-201)
- Parser rejection tests: `internal/parser/parser_test.go`

---

*This documentation structure will be expanded in Milestone 10.*
