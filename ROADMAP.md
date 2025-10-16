# Ticketr Operational Readiness Roadmap

## Purpose

This roadmap captures the opinionated engineering plan to close the usability gaps identified during the Ticketr operational-readiness review. It is written so that a future contributor (human or LLM) can pick it up, understand the desired end-state, and execute the work without re-discovering context.

## How to Use This Document

- Work through milestones sequentially; unblocked parallelisation is called out explicitly.
- Each checkbox is actionable. Check items off only after code, tests, and documentation updates are merged.
- Keep this file updated—add notes, links to PRs, or test evidence beneath each task as you progress.
- All file paths are repository-relative.

## Key References & Background

- Product requirements: `REQUIREMENTS-v2.md`
- CLI entrypoint: `cmd/ticketr/main.go`
- Parser and writer: `internal/parser/parser.go`, `internal/adapters/filesystem/file_repository.go`
- Push/Pull services: `internal/core/services/*`
- Jira adapter: `internal/adapters/jira/jira_adapter.go`
- State management: `internal/state/manager.go`
- User docs & templates: `README.md`, `examples/`

## Execution Preparation

1. Read `REQUIREMENTS-v2.md` to confirm acceptance targets (e.g., PROD-004 logging, PROD-009 inheritance).
2. Run `go test ./...` to capture current baseline (expect failures if environment prerequisites are missing).
3. Note any uncommitted local changes so you do not overwrite them.
4. Set `JIRA_*` environment variables or prepare mocks before touching Jira integrations.

---

## Milestone 0 – Repository Recon & Standards Alignment ✅

Goal: Ensure contributors understand canonical formats before modifying behaviour.

**Status:** COMPLETE (Commit: ecc24d4)
**Completed:** 2025-10-16

- [x] Confirm `# TICKET:` is the canonical Markdown heading and that legacy `# STORY:` content will be rejected with a clear error. (`internal/parser/parser.go`)
- [x] Document the deprecation path and rationale in `README.md` and `REQUIREMENTS-v2.md`, including migration guidance.
- [x] Capture sample legacy files under `testdata/legacy_story/*.md` for regression tests.
- [x] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [x] Update documentation: refresh `README.md` overview to declare the canonical `# TICKET:` format and point migration notes to `REQUIREMENTS-v2.md`; ensure `REQUIREMENTS-v2.md` references legacy deprecation; add a short note to `docs/README.md` (create if missing) describing the legacy samples in `testdata/legacy_story/`.

**Deliverables:**
- README.md: Updated with TICKET format and migration guidance
- REQUIREMENTS-v2.md: Enhanced PROD-201 with deprecation subsection
- docs/README.md: Created documenting test fixture strategy
- All user-facing examples use canonical # TICKET: format

---

## Milestone 1 – Canonical Markdown Schema & Tooling ✅

Goal: Align documentation, templates, and tooling with the canonical schema; provide migration support.

**Status:** COMPLETE (Commits: 3a49e78, 943fbfa)
**Completed:** 2025-10-16
**Test Results:** 36 passed, 0 failed, 3 skipped (JIRA integration)

- [x] Update all examples, templates, and README snippets to use the `# TICKET:` schema. (`README.md`, `examples/*.md`)
- [x] Add a parser error explaining how to migrate when a legacy heading is detected. (`internal/parser/parser.go`)
- [x] Introduce a `ticketr migrate <file>` helper (optional but recommended) that rewrites legacy `# STORY:` blocks to `# TICKET:` while preserving content. (`cmd/ticketr/main.go`, new helper file)
- [x] Extend parser unit tests to cover canonical parsing and rejection/migration paths. (`internal/parser/parser_test.go`)
- [x] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [x] Update documentation: rewrite README quick-start and advanced sections to the new schema; add `docs/migration-guide.md` covering legacy `# STORY` conversion; update `examples/README.md` (create) to describe templates.

**Deliverables:**
- Parser rejection: Legacy # STORY: format rejected with actionable error (line number + migration command)
- Migration command: ticketr migrate with dry-run (default) and --write modes
- Test coverage: 11 new tests (3 parser + 7 migration + 1 service)
- Documentation: docs/migration-guide.md (175 lines), examples/README.md (155 lines)
- Examples: All 3 templates migrated to canonical format
- README: Enhanced migration section with quick commands and cross-references

**Known Issues:**
- Minor: Migration guide references non-existent `ticketr validate` command (non-blocking)

Dependencies: Milestone 0.

---

## Milestone 2 – Pull Conflict Resolution Flag ✅

Goal: Make conflict resolution operable from the CLI.

**Status:** COMPLETE (Commit: 18c67a3)
**Completed:** 2025-10-16
**Test Results:** 39 passed, 0 failed, 3 skipped (JIRA integration)

- [x] Add a `--force` (or `--force-remote`) flag to `ticketr pull` and thread it through to `services.PullOptions.Force`. (`cmd/ticketr/main.go`, `internal/core/services/pull_service.go`)
- [x] Update user messaging to reflect the new flag when conflicts occur. (`cmd/ticketr/main.go`)
- [x] Extend pull-service tests to cover conflict detection with and without the force flag. (`internal/core/services/pull_service_test.go`)
- [x] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [x] Update documentation: extend the `README.md` pull command section with the new `--force` flag, and ensure `cmd/ticketr/main.go` Cobra descriptions include the flag in help text (via inline comments/check-in instructions).

**Deliverables:**
- CLI flag: --force wired to PullOptions.Force (line 103, 302 in main.go)
- Test: TestPullService_ConflictResolvedWithForce validates force behavior
- Documentation: README updated with flag description and conflict resolution examples
- Help text: Command description includes conflict detection guidance
- Service layer: Leveraged existing Force support (no service changes needed)

Parallelisable with Milestone 1.

---

## Milestone 3 – First-Run Pull Safety ✅

Goal: Allow `ticketr pull` to succeed when no local file exists.

**Status:** COMPLETE (Commits: 8ce7ae1, 385c2c9)
**Completed:** 2025-10-16
**Test Results:** 42 passed, 0 failed, 3 skipped (JIRA integration)

- [x] Map `os.ErrNotExist` to `ports.ErrFileNotFound` inside `FileRepository.GetTickets`. (`internal/adapters/filesystem/file_repository.go`)
- [x] Update `PullService` to treat missing local files as an empty ticket set. (`internal/core/services/pull_service.go`)
- [x] Add regression tests for first-run pull workflows. (`internal/core/services/pull_service_test.go`)
- [x] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [x] Update documentation: add a "First Pull" troubleshooting subsection to `README.md` and, if present, append a note to `docs/troubleshooting.md`.

**Deliverables:**
- Error mapping: FileRepository.GetTickets maps os.ErrNotExist to ports.ErrFileNotFound
- Service verification: PullService already correctly handled missing files (no changes needed)
- Test coverage: 3 new tests (TC-303.3, TC-303.4, TC-303.5) for first-run scenarios
- Documentation: README "First Pull" section with troubleshooting guidance
- Enhanced mocks: Added callback functions for dynamic test behavior

Dependencies: None (can run in parallel).

---

## Milestone 4 – Deterministic State Hashing

Goal: Eliminate spurious change detection caused by Go’s non-deterministic map iteration.

- [ ] Sort custom field keys and task lists (and nested custom fields) before writing to the hash in `StateManager.CalculateHash`. (`internal/state/manager.go`)
- [ ] Add unit tests ensuring that permutations of map insert order produce identical hashes. (`internal/state/manager_test.go`)
- [ ] Verify push skip logic still behaves correctly after changes. (`internal/core/services/push_service_test.go`)
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: expand the `.ticketr.state` description in `README.md` and add deterministic-hash details to `docs/state-management.md` (create if missing).

Dependencies: None.

---

## Milestone 5 – Force-Partial Upload Semantics

Goal: Make `--force-partial-upload` honour its contract during validation and processing.

- [ ] Update CLI pre-flight validation to downgrade errors to warnings when `--force-partial-upload` is set, still emitting precise line diagnostics. (`cmd/ticketr/main.go`)
- [ ] Teach `TicketService.ProcessTicketsWithOptions` (or its replacement; see Milestone 9) to record per-ticket validation failures and continue processing when forced. (`internal/core/services/ticket_service.go`)
- [ ] Expand tests for mixed-validity uploads covering both forced and normal modes. (`internal/core/services/push_service_comprehensive_test.go`)
- [ ] Refresh README usage docs for the flag. (`README.md`)
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: clarify `--force-partial-upload` semantics in README usage tables and record validation behaviour in `docs/cli-reference.md` (create if missing).

Dependencies: Milestone 0.

---

## Milestone 6 – Persistent Execution Log (PROD-004)

Goal: Fulfil PROD-004 by writing execution summaries to disk.

- [ ] Create a logging utility that writes to `.ticketr/logs/<timestamp>.log` (configurable path via env/config). (`internal/logging` package – new)
- [ ] Pipe push/pull command summaries and errors into the log while preserving console output. (`cmd/ticketr/main.go`, `internal/core/services/*`)
- [ ] Ensure sensitive values (API keys, email) are redacted. (`internal/logging/*`)
- [ ] Document log location, rotation policy, and how to disable/override logging. (`README.md`)
- [ ] Add tests verifying log file creation (can use temp directories). (`internal/core/services/tests`, new logging tests)
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: add a “Logging” section to `README.md`, create `docs/logging.md` with file locations/rotation tips, and update any troubleshooting references.

Dependencies: None, but coordinate with Milestone 9 for shared logging hooks.

---

## Milestone 7 – Field Inheritance Compliance (PROD-009/202)

Goal: Ensure tasks inherit parent fields with local overrides before hitting Jira.

- [ ] Use `calculateFinalFields` (or successor) to merge parent `CustomFields` into each task before calling the Jira adapter. (`internal/core/services/ticket_service.go`, `internal/core/services/push_service.go`)
- [ ] Update Jira adapter methods to accept merged fields and map them to correct Jira payload fields. (`internal/adapters/jira/jira_adapter.go`)
- [ ] Add unit tests verifying that tasks inherit parent values unless overridden. (`internal/core/services/ticket_service_test.go`)
- [ ] Update documentation to describe field inheritance rules. (`README.md`, `REQUIREMENTS-v2.md`)
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: detail field inheritance rules in `README.md` and align `REQUIREMENTS-v2.md` acceptance text; create `docs/field-inheritance.md` summarising examples if not already present.

Dependencies: Milestone 9 (if push service replacement happens first).

---

## Milestone 8 – Pulling Tasks/Subtasks

Goal: Round-trip task hierarchies between Markdown and Jira.

- [ ] Enhance `JiraAdapter.SearchTickets` to request subtasks, building `ticket.Tasks` with consistent schema. (`internal/adapters/jira/jira_adapter.go`)
- [ ] Ensure pulled Markdown includes tasks (writer already supports it; verify). (`internal/adapters/filesystem/file_repository.go`)
- [ ] Handle parent-child resolution for subtasks, including field inheritance. (`internal/adapters/jira/jira_adapter.go`)
- [ ] Add integration-style tests covering push → Jira mock → pull cycles. (`internal/adapters/jira/jira_adapter_dynamic_test.go`, new fixtures)
- [ ] Document the pull behaviour with examples. (`README.md`)
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: refresh pull examples in `README.md`, document push/pull round-trip in `docs/workflow.md`, and note subtask behaviour in `docs/field-inheritance.md`.

Dependencies: Milestone 7 (shared understanding of task fields).

---

## Milestone 9 – State-Aware Push Integration

Goal: Replace the legacy push flow with the stateful service that honours PROD-204.

- [ ] Swap `TicketService` for `PushService` in `cmd/ticketr/main.go`, wiring up state manager initialisation and options. (`cmd/ticketr/main.go`)
- [ ] Align `--force-partial-upload` behaviour between CLI and push service. (`internal/core/services/push_service.go`)
- [ ] Update push tests to include state-skipping scenarios. (`internal/core/services/push_service_test.go`)
- [ ] Clean up any dead code paths left in `TicketService` (or repurpose it if needed).
- [ ] Document state file semantics for users. (`README.md`)
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: update README push command narrative to reference the state-aware flow and adjust `docs/state-management.md` for the new CLI wiring.

Dependencies: Milestones 4 & 5 (hash determinism and force semantics).

---

## Milestone 10 – Documentation & Developer Experience

Goal: Capture the new behaviours and ensure contributors/users understand them.

- [ ] Refresh the README with updated command usage, logging, conflict resolution, state handling, and migration info. (`README.md`)
- [ ] Update requirement traceability to show compliance with PROD/NFR/DEV items. (`REQUIREMENTS-v2.md`)
- [ ] Provide an end-to-end walkthrough (e.g., sample Markdown → push → pull → log review). (`README.md` or new `docs/WORKFLOW.md`)
- [ ] Add CONTRIBUTING notes about testing expectations and logging artefacts.
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: finalise `README.md` and `docs/workflow.md`, add `CONTRIBUTING.md` detailing doc/test expectations, and ensure `docs/style-guide.md` links are current.

Dependencies: All prior milestones.

---

## Milestone 11 – Quality Gates & Automation

Goal: Lock down regressions via automated checks.

- [ ] Expand unit/integration test coverage to include new features. (`internal/...`)
- [ ] Add smoke tests for CLI flows (e.g., scripts under `evidence/` or new `integration/` directory).
- [ ] Integrate `go vet`, `staticcheck`, and formatting checks into CI (GitHub Actions or local script). (`.github/workflows/*.yml` or new pipeline)
- [ ] Document a verification checklist (commands to run, expected outputs) in `HANDOFF-BEFORE-PHASE-3.md` or a new QA doc.
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: create/refresh `docs/qa-checklist.md` and document CI tooling in `docs/ci.md`, referencing them from `CONTRIBUTING.md`.

Dependencies: Milestones 1–10.

---

## Milestone 12 – Requirements Consolidation & Documentation Governance

Goal: Base the product exclusively on the v2 specification and publish an authoritative documentation set.

- [ ] Deprecate `REQUIREMENTS.md` (v1) by archiving it under `docs/legacy/` and updating `REQUIREMENTS-v2.md` with any missing requirements or clarifications. (`REQUIREMENTS.md`, `REQUIREMENTS-v2.md`)
- [ ] Rewrite README examples and workflows to use the canonical `# TICKET` schema, replacing all legacy `# STORY` snippets. (`README.md:52`, `README.md:79`, `README.md:104`, `README.md:250`)
- [ ] Update `examples/*.md` to the new schema or move them under an archived legacy folder with clear warnings. (`examples/quick-story.md`, `examples/epic-template.md`, `examples/sprint-template.md`)
- [ ] Retire phase playbooks (`PHASE-1.md`, `PHASE-2.md`, `PHASE-3.md`, `phase-hardening.md`, `PHASE-n.md`) by moving them to `docs/history/` with a README explaining their legacy status.
- [ ] Refresh `HANDOFF-BEFORE-PHASE-3.md` to reflect the current architecture or replace it with an up-to-date engineering overview. (`HANDOFF-BEFORE-PHASE-3.md`)
- [ ] Add a docs style guide / contribution rules covering tone, code fences, anchors, and file naming, referenced from `CONTRIBUTING.md` (create if missing).
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: move v1 docs to `docs/legacy/` with a README explaining scope, update `REQUIREMENTS-v2.md` cross-links, formalise `docs/style-guide.md`, and expand `CONTRIBUTING.md` with documentation governance.

Dependencies: Milestones 0–10 (documentation updates can run in parallel but depend on canonical schema decisions).

---

## Milestone 13 – Repository Hygiene & Release Readiness

Goal: Ensure the repository is production-ready, free of sensitive artefacts, and publishes releasable builds.

- [ ] Remove stray build artefacts from source control (e.g., the root `ticketr` binary) and verify `.gitignore` patterns are precise (`/ticketr` already present). Confirm the repo is clean afterwards.
- [ ] Purge the committed `.env` secrets, rotate the exposed API key, and document the safe way to manage credentials locally. (`.env`, `README.md` configuration section)
- [ ] Add SECURITY.md (responsible disclosure) and CODE_OF_CONDUCT.md to set community expectations.
- [ ] Establish release management: introduce CHANGELOG.md, SemVer policy, and create a GitHub Actions workflow for tagged releases that builds and uploads binaries/containers.
- [ ] Add CI status badges (build/test, coverage) to the README and document support cadence (issues/discussions response expectations).
- [ ] Verify `check-issue-types.sh` and other scripts have homes (move to `scripts/`), documenting utility scripts in README or CONTRIBUTING.
- [ ] Update/extend automated tests affected by this milestone and run `go test ./...`.
- [ ] Update documentation: add `SECURITY.md`, `CODE_OF_CONDUCT.md`, `CHANGELOG.md`, document release steps in `docs/release-process.md`, and update README badges/support policy.

Dependencies: Milestones 4, 6, 9–12 (requires stable behaviours and documentation).

---

## Delivery Checklist

- [ ] All milestone tasks completed and checked off.
- [ ] Tests (`go test ./...`) passing locally and in CI.
- [ ] Documentation reflects final behaviour.
- [ ] `ROADMAP.md` updated with completion notes and any follow-up ideas.

Use this checklist before closing out the readiness effort.

---

## Testing Discipline (Applies to Every Milestone)

- [ ] For each code change, create or update unit/integration tests that exercise the new behaviour.
- [ ] After finishing the work for a milestone (or substantial subtask), execute `go test ./...` locally and capture failures immediately.
- [ ] Before opening a PR or promoting a release candidate, ensure CI runs the same test suite and publishes results.
- [ ] Record notable test evidence (commands, outputs, coverage deltas) in commit messages or accompanying docs so future auditors can trace verification.

These checks are cumulative—do not mark a milestone complete without meeting them.
