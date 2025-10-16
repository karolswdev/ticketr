# Legacy Documentation

## Overview

This directory contains deprecated documentation from earlier phases of the ticketr project. These files are preserved for historical reference only.

## Contents

### REQUIREMENTS-v1.md

This is the original Software Requirements Specification (v1.0) from the project's inception. It has been superseded by `/REQUIREMENTS-v2.md` in the root directory.

**Status:** Deprecated
**Superseded by:** [/REQUIREMENTS-v2.md](/REQUIREMENTS-v2.md)

## Migration Rationale

The v1 requirements document was based on the original "Jira Story Creator" concept, which used a STORY-centric schema. As the project evolved, we recognized that:

1. **Generic Ticket Schema:** Jira supports multiple issue types (Story, Task, Bug, Epic, etc.), not just Stories. The v1 schema was too narrow.
2. **Flexibility:** The `# TICKET:` format in v2 allows users to specify any Jira issue type, making ticketr more versatile.
3. **Consistency:** The v2 requirements align with the current codebase architecture (Hexagonal/Ports & Adapters) and operational model.

## Current Documentation

For current project requirements, architecture, and workflows, see:

- [REQUIREMENTS-v2.md](/REQUIREMENTS-v2.md) - Current requirements specification
- [ARCHITECTURE.md](/ARCHITECTURE.md) - System architecture and design
- [docs/WORKFLOW.md](/docs/WORKFLOW.md) - Developer workflows
- [CONTRIBUTING.md](/CONTRIBUTING.md) - Contribution guidelines
- [ROADMAP.md](/ROADMAP.md) - Project roadmap and milestones

## Historical Context

These documents guided the project through milestones 0-7 and represent the foundation upon which the current system was built. They are kept for:

- Historical reference
- Understanding design evolution
- Traceability of architectural decisions
- Onboarding context for new contributors

**Note:** Do not use these documents for current development. Always refer to the latest documentation in the root and `/docs` directories.
