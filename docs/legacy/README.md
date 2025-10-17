# Legacy Documentation

## Overview

This directory contains deprecated documentation from earlier phases of the ticketr project. These files are preserved for historical reference only.

## Contents

### REQUIREMENTS-v1.md

This is the original Software Requirements Specification drafted early in the project. It has been superseded by the consolidated requirements in `docs/development/REQUIREMENTS.md`.

**Status:** Deprecated
**Superseded by:** [docs/development/REQUIREMENTS.md](/docs/development/REQUIREMENTS.md)

## Migration Rationale

The early requirements document was based on the original "Jira Story Creator" concept, which used a STORY-centric schema. Before the 1.0 release we broadened the format so that:

1. **Generic Ticket Schema:** Jira supports multiple issue types (Story, Task, Bug, Epic, etc.), not just Stories. The story-only schema was too narrow.
2. **Flexibility:** The unified `# TICKET:` format lets ticketr represent any Jira issue type while keeping a consistent structure.
3. **Consistency:** The updated requirements align with the current codebase architecture (Hexagonal/Ports & Adapters) and operational model.

## Current Documentation

For current project requirements, architecture, and workflows, see:

- [docs/development/REQUIREMENTS.md](/docs/development/REQUIREMENTS.md) - Current requirements specification
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
