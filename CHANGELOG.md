# Changelog

All notable changes to Ticketr will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- SECURITY.md with responsible disclosure policy and security best practices
- Comprehensive release automation workflow
- Release process documentation

### Changed
- Improved credential management documentation

## [0.2.0] - 2025-10-16

### Added
- **Milestone 12**: Requirements consolidation and documentation governance
  - docs/ARCHITECTURE.md with comprehensive system architecture
  - docs/style-guide.md for documentation standards
  - Legacy v1 requirements archived to docs/legacy/
  - Phase playbooks archived to docs/history/
- **Milestone 11**: Quality gates and automation
  - GitHub Actions CI/CD pipeline with 5 jobs (build, test, coverage, lint, smoke-tests)
  - Smoke test suite with 7 scenarios
  - Quality automation script (scripts/quality.sh)
  - Test coverage increased to 52.5% (+27 new tests)
  - OS matrix testing (Ubuntu, macOS; Go 1.21, 1.22, 1.23)
- **Milestone 10**: Documentation and developer experience
  - docs/WORKFLOW.md with end-to-end walkthrough (379 lines)
  - CONTRIBUTING.md with testing and architecture guidelines
  - README Quick Reference section
  - Enhanced requirement traceability in docs/development/REQUIREMENTS.md

### Changed
- All documentation cross-references established
- Developer onboarding process formalized

## [0.1.0] - 2025-10-16

### Added
- **Milestone 9**: State-aware push integration (PROD-204)
  - StateManager integrated into CLI push command
  - Field inheritance in PushService
  - Unchanged tickets automatically skipped
  - TicketService deprecated in favor of PushService
- **Milestone 8**: Pulling tasks/subtasks
  - Enhanced SearchTickets() to fetch subtasks
  - parseJiraSubtask() for Jira-to-domain conversion
  - Round-trip compatibility (pull → markdown → push)
  - 4 new integration tests (TC-208.1 through TC-208.4)
- **Milestone 7**: Field inheritance compliance (PROD-009/202)
  - calculateFinalFields() method for hierarchical field merging
  - Tasks inherit parent CustomFields with local overrides
  - Integration tested with production JIRA instance
  - 4 new unit tests + 6 subtasks tested in JIRA
  - docs/field-inheritance.md and examples
- **Milestone 6**: Persistent execution log (PROD-004)
  - Human-readable timestamped logs
  - Log location: `.ticketr/logs/<timestamp>.log`
  - TICKETR_LOG_DIR environment variable support
  - Automatic redaction of sensitive data
  - Log rotation (keeps last 10 files)
  - Comprehensive logging documentation
- **Milestone 5**: Force-partial-upload semantics
  - Pre-flight validation respects --force-partial-upload flag
  - Validation errors downgraded to warnings with flag
  - Exit codes: 0 (partial success), 1 (validation), 2 (runtime errors)
  - 4 new comprehensive test cases (TC-501.1 through TC-501.4)
- **Milestone 4**: Deterministic state hashing
  - Sorted custom field keys for consistent hashing
  - Sorted task lists and nested custom fields
  - 3 new determinism tests (100-iteration stability, permutations)
  - docs/state-management.md (149 lines)
- **Milestone 3**: First-run pull safety
  - FileRepository.GetTickets maps os.ErrNotExist to ports.ErrFileNotFound
  - PullService handles missing local files gracefully
  - 3 new first-run tests (TC-303.3, TC-303.4, TC-303.5)
  - "First Pull" troubleshooting section in README
- **Milestone 2**: Pull conflict resolution flag
  - --force flag for pull command
  - Conflict detection and resolution workflow
  - TestPullService_ConflictResolvedWithForce test
- **Milestone 1**: Canonical markdown schema and tooling
  - Parser rejection of legacy # STORY: format
  - ticketr migrate command with dry-run and --write modes
  - 11 new tests (3 parser + 7 migration + 1 service)
  - docs/migration-guide.md (175 lines)
  - examples/README.md (155 lines)
- **Milestone 0**: Repository recon and standards alignment
  - Canonical # TICKET: format established
  - Legacy # STORY: format deprecated with clear errors
  - testdata/legacy_story/ fixtures for regression tests
  - docs/README.md documenting test strategy

### Changed
- Renamed from jira-story-creator to ticketr
- All examples migrated to canonical # TICKET: format
- Enhanced README with migration guidance

### Fixed
- Spurious change detection from non-deterministic map iteration
- First-run pull failures when local file doesn't exist
- Conflict resolution workflow without --force flag

## [0.0.1] - Initial Development

### Added
- Core Markdown parser for ticket definitions
- Jira adapter for creating and updating issues
- CLI interface with basic commands
- Docker support with Alpine-based image (~15MB)
- Bidirectional sync (push to Jira, pull from Jira)
- State tracking with .ticketr.state file
- Environment-based configuration
- Advanced CLI flags (--verbose, --force-partial-upload)
- Schema discovery command
- Custom field mapping support

### Features
- Markdown-first workflow for ticket management
- Smart update detection
- Hierarchical task management (parent tickets + subtasks)
- CI/CD ready with non-interactive modes

---

## Semantic Versioning Policy

Ticketr follows [Semantic Versioning 2.0.0](https://semver.org/):

Given a version number MAJOR.MINOR.PATCH, increment the:

1. **MAJOR** version when you make incompatible API changes
   - Breaking CLI argument changes
   - Incompatible Markdown schema changes
   - Removal of supported features

2. **MINOR** version when you add functionality in a backward compatible manner
   - New CLI commands
   - New features (field inheritance, logging, etc.)
   - Enhanced functionality (subtask pulling, conflict resolution)

3. **PATCH** version when you make backward compatible bug fixes
   - Bug fixes
   - Documentation updates
   - Performance improvements
   - Security patches

### Pre-1.0 Versions (0.x.x)

During pre-1.0 development:
- Breaking changes may occur in MINOR versions
- The public API is not yet stable
- Use with caution in production environments

### Version 1.0.0 Criteria

Version 1.0.0 will be released when:
- All items in docs/development/ROADMAP.md Milestone 13 are complete
- Production-ready security practices implemented
- Stable public API established
- Comprehensive test coverage (>70%)
- Complete documentation
- At least 3 months of field testing

## Release Process

See [docs/release-process.md](docs/release-process.md) for detailed release procedures.

### Quick Release Steps

1. Update CHANGELOG.md with release notes
2. Update version in relevant files
3. Create git tag: `git tag -a v0.2.0 -m "Release v0.2.0"`
4. Push tag: `git push origin v0.2.0`
5. GitHub Actions automatically builds and publishes release

---

**Note**: Dates in this changelog use YYYY-MM-DD format (ISO 8601).

[Unreleased]: https://github.com/karolswdev/ticktr/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/karolswdev/ticktr/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/karolswdev/ticktr/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/karolswdev/ticktr/releases/tag/v0.0.1
