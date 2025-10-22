# Phase 5 Execution Checklist

**Version:** 1.0 | **Date:** October 19, 2025

*Weekly execution tracking for Ticketr v3.0 Phase 5: Advanced Features (Weeks 18-20)*

---

## Phase 5 Overview

**Goal**: Add power-user features leveraging v3.0 architecture
**Duration**: 3 weeks (Weeks 18-20)
**Success Criteria**: Bulk operations, templates, smart sync, JQL aliases all functional with zero regressions

---

## Pre-Phase 5 Readiness

**Prerequisites** (verify before starting Week 18):

- [ ] Phase 4 TUI complete and deployed
- [ ] Milestone 18 (credential profiles) complete
- [ ] All critical bugs fixed (P0 bugs from PHASE4-WEEK16-COMPLETE.md)
- [ ] Working directory clean: `git status`
- [ ] All tests passing: `go test ./...`
- [ ] Current baseline coverage: ___% (record for comparison)

---

## Week 18: Bulk Operations

**Goal**: Multi-ticket selection and operations in CLI and TUI

### Day 1: Domain Model & Service Interface

**Builder Tasks**:
- [ ] Create `internal/core/domain/bulk_operation.go`
- [ ] Define BulkOperation struct (Action, TicketIDs, Changes)
- [ ] Add validation logic (max 100 tickets per operation)
- [ ] Create `internal/core/ports/bulk_operation_service.go` interface
- [ ] Initial tests: domain model validation

**Acceptance**:
- [ ] Domain model compiles
- [ ] Validation tests pass (>60% coverage)
- [ ] Interface defined

**Estimated**: 2 hours | **Actual**: _____

### Day 2: Service Implementation

**Builder Tasks**:
- [ ] Create `internal/core/services/bulk_operation_service.go`
- [ ] Implement BatchUpdate() method
- [ ] Add transaction safety (rollback on partial failure)
- [ ] Implement progress callback pattern
- [ ] Create comprehensive test suite

**Verifier Tasks**:
- [ ] Run full test suite: `go test ./...`
- [ ] Verify service coverage >80%
- [ ] Validate transaction rollback behavior
- [ ] Check error handling for partial failures

**Acceptance**:
- [ ] Service compiles and all tests pass
- [ ] Coverage ≥80% for BulkOperationService
- [ ] Handles partial failures gracefully
- [ ] Verifier recommendation: APPROVE

**Estimated**: 4 hours | **Actual**: _____

### Day 3: CLI Integration

**Builder Tasks**:
- [ ] Create `cmd/ticketr/bulk_commands.go`
- [ ] Implement `ticketr bulk update --ids X,Y --set status=Done`
- [ ] Implement `ticketr bulk move --ids X,Y --parent Z`
- [ ] Implement `ticketr bulk delete --ids X,Y --confirm`
- [ ] Add confirmation prompts for destructive ops
- [ ] Wire progress indicators

**Verifier Tasks**:
- [ ] Test all CLI commands manually
- [ ] Verify confirmation prompts work
- [ ] Check progress indicators display
- [ ] Run CLI command tests

**Acceptance**:
- [ ] All commands functional
- [ ] Confirmation prompts prevent accidental deletion
- [ ] Progress feedback clear
- [ ] CLI tests pass

**Estimated**: 3 hours | **Actual**: _____

### Day 4-5: TUI Integration

**Builder Tasks**:
- [ ] Modify `internal/adapters/tui/views/ticket_tree.go`
  - [ ] Add multi-select checkboxes (Space bar to toggle)
  - [ ] Add 'a' keybinding for select all
  - [ ] Visual indicators for selected tickets
- [ ] Create bulk operations modal
  - [ ] Bulk update modal (change fields)
  - [ ] Bulk move modal (select new parent)
  - [ ] Bulk delete confirmation
- [ ] Update `internal/adapters/tui/app.go` routing
- [ ] Update help view with new keybindings

**Verifier Tasks**:
- [ ] Manual TUI testing: multi-select workflow
- [ ] Verify keybindings work as documented
- [ ] Test bulk operations complete successfully
- [ ] Check visual feedback (checkboxes, status)

**Acceptance**:
- [ ] Multi-select works smoothly
- [ ] Bulk operations accessible via keys
- [ ] Status feedback after operations
- [ ] TUI tests pass (where applicable)

**Estimated**: 6 hours | **Actual**: _____

### Day 5 (afternoon): Documentation

**Scribe Tasks**:
- [ ] Update `README.md`: Add bulk operations to Features section
- [ ] Create `docs/bulk-operations-guide.md`:
  - [ ] Introduction: What are bulk operations?
  - [ ] CLI command reference with examples
  - [ ] TUI workflow with keybindings
  - [ ] Safety considerations
  - [ ] Troubleshooting
- [ ] Update `CHANGELOG.md`: Week 18 bulk operations entry
- [ ] Update `docs/v3-implementation-roadmap.md`: Mark Week 18 tasks complete
- [ ] Test all examples manually

**Acceptance**:
- [ ] All docs updated
- [ ] Examples tested and accurate
- [ ] Cross-references valid
- [ ] Markdown renders correctly

**Estimated**: 2 hours | **Actual**: _____

### Week 18 Final Checklist

**Pre-Commit**:
- [ ] All tests passing: `go test ./...`
- [ ] Build successful: `go build ./...`
- [ ] Coverage maintained or improved: ___%
- [ ] Documentation complete
- [ ] Manual testing complete (CLI + TUI)

**Git Commits**:
- [ ] Commit 1: `feat(bulk): Implement bulk operations domain and service`
- [ ] Commit 2: `feat(bulk): Add CLI commands for bulk operations`
- [ ] Commit 3: `feat(bulk): Add TUI multi-select and bulk operations`
- [ ] Commit 4: `docs(bulk): Add comprehensive bulk operations guide`

**Week 18 Complete**: [ ] YES / [ ] NO (if NO, why: _____________)

---

## Week 19: Templates + Smart Sync

**Goal**: Template system and intelligent conflict resolution

### Day 1-2: Template System

**Builder Tasks**:
- [ ] Create `internal/templates/parser.go`
- [ ] Implement YAML template parsing
- [ ] Create variable substitution engine ({{.Name}}, {{.Sprint}})
- [ ] Add template validation (schema checks)
- [ ] Create `cmd/ticketr/template_commands.go`
- [ ] Implement `ticketr template apply feature.yaml`
- [ ] Add TUI template selector modal
- [ ] Create comprehensive tests

**Verifier Tasks**:
- [ ] Run template parser tests
- [ ] Verify variable substitution works correctly
- [ ] Test template validation catches errors
- [ ] Manual test: apply template creates tickets
- [ ] Check coverage ≥70% for template package

**Acceptance**:
- [ ] Templates parse YAML correctly
- [ ] Variable substitution works (all supported types)
- [ ] Validation prevents invalid templates
- [ ] CLI and TUI integration functional
- [ ] Tests comprehensive

**Estimated**: 12 hours | **Actual**: _____

### Day 3-5: Smart Sync ✅ COMPLETE

**Status**: ✅ Delivered 2025-10-20
**Commit**: 09b0053 feat(templates): Implement template system with YAML parser (Slice 1)

**Builder Tasks**:
- ✅ Define `internal/core/ports/sync_strategy.go` interface
- ✅ Implement LocalWinsStrategy
- ✅ Implement RemoteWinsStrategy
- ✅ Implement ThreeWayMergeStrategy
- ✅ Update PullService to use strategies
- ⏸️ Add configuration option for default strategy (deferred to v3.1)
- ⏸️ Create TUI conflict resolution modal (deferred to future slice)
  - ⏸️ Show local vs remote changes
  - ⏸️ Allow user to choose strategy
  - ⏸️ Preview merge result
- ✅ Comprehensive tests (all strategies, all conflict types)

**Verifier Tasks**:
- ✅ Test all sync strategies
- ✅ Verify no data loss in any conflict scenario
- ✅ Test three-way merge with compatible/incompatible changes
- ⏸️ Manual test: pull with conflicts, resolve in TUI (TUI modal deferred)
- ✅ Check coverage ≥80% for sync strategies (93.95% achieved)

**Acceptance**:
- ✅ All strategies implemented and tested
- ✅ No data loss in any conflict scenario
- ⏸️ User can choose strategy via config or flag (v3.1 feature)
- ⏸️ TUI conflict modal clear and functional (future slice)
- ✅ Tests comprehensive

**Deliverables**:
- ✅ SyncStrategy interface (`internal/core/ports/sync_strategy.go`)
- ✅ LocalWinsStrategy implementation
- ✅ RemoteWinsStrategy implementation (default)
- ✅ ThreeWayMergeStrategy implementation
- ✅ PullService integration with `NewPullServiceWithStrategy()`
- ✅ 64 new tests (55 unit + 9 integration, 93.95% coverage)
- ✅ Documentation (`docs/sync-strategies-guide.md`, 600+ lines)

**Notes**:
- CLI flag `--strategy` and config file support deferred to v3.1 (future enhancement)
- TUI conflict modal deferred to future slice (not blocking for Slice 2 completion)
- Default behavior (RemoteWins) preserves backward compatibility with v2.x
- All core functionality complete and tested
- Verifier status: APPROVED

**Estimated**: 15 hours | **Actual**: ~12 hours (ahead of schedule)

### Day 5 (afternoon): Documentation ✅ COMPLETE

**Status**: ✅ Delivered 2025-10-20

**Scribe Tasks**:
- ⏸️ Create `docs/template-guide.md` (deferred - template system Slice 1)
  - ⏸️ Template syntax reference
  - ⏸️ Variable types and usage
  - ⏸️ Example templates (epic, feature, sprint)
  - ⏸️ CLI and TUI workflows
- ✅ Create `docs/sync-strategies-guide.md`:
  - ✅ Sync strategy overview
  - ✅ When to use each strategy
  - ✅ Conflict resolution workflows
  - ✅ Three-way merge explanation
  - ✅ Field-level merging details
  - ✅ Troubleshooting and best practices
  - ✅ Examples for all three strategies
- ✅ Update `README.md`: Add smart sync features
- ✅ Update `CHANGELOG.md`: Week 19 Slice 2 entry
- ✅ Update roadmap (PHASE5-EXECUTION-CHECKLIST.md)

**Deliverables**:
- ✅ `docs/sync-strategies-guide.md` (684 lines)
  - Comprehensive guide covering all three strategies
  - Decision matrix for strategy selection
  - Field-level merge explanation with examples
  - Compatible vs incompatible change examples
  - Troubleshooting section (5 common issues)
  - Best practices (7 recommendations)
  - Technical details (coverage, performance, hash algorithm)
- ✅ README.md updates (+38 lines)
  - Added Smart Sync Strategies to Features section
  - Added detailed "Conflict Detection & Smart Sync Strategies" section
  - Added sync-strategies-guide.md to documentation references
- ✅ CHANGELOG.md updates (+46 lines)
  - Week 19 Slice 2 complete entry
  - API changes documented
  - Testing metrics documented
  - Notes on backward compatibility

**Acceptance**:
- ✅ All docs complete and accurate
- ✅ Examples tested (all strategies verified)
- ✅ Cross-references valid (all links checked)

**Estimated**: 3 hours | **Actual**: ~2.5 hours (efficient delivery)

### Week 19 Final Checklist

**Pre-Commit**:
- ✅ All tests passing (64 new tests, 98 services tests passing)
- ✅ Build successful (`go build ./...` - clean)
- ✅ Coverage maintained or improved: 93.95% (sync_strategies.go)
- ✅ Documentation complete (sync-strategies-guide.md, CHANGELOG.md, README.md)
- ⏸️ Manual testing complete (deferred - TUI modal not implemented in Slice 2)

**Git Commits**:
- ✅ Commit 1: `09b0053` - `feat(templates): Implement template system with YAML parser` (Slice 1)
- ✅ Commit 2: `0f98cc9` - `feat(sync): Add smart sync strategies for conflict resolution` (Slice 2)
- ⏸️ Commit 3: `feat(tui): Add template selector and conflict resolution modals` (deferred to future)
- ⏸️ Commit 4: `docs(phase5): Document templates and sync strategies` (merged with Commit 2)

**Week 19 Complete**: ✅ **YES** (Slice 2 delivered, TUI components deferred to v3.1)

---

## Week 20: JQL Aliases + Polish

**Goal**: JQL aliases, performance optimization, final polish

### Day 1-2: JQL Aliases ✅ COMPLETE

**Status**: ✅ Delivered 2025-10-20
**Slice**: Week 20 Slice 1 (JQL Aliases)

**Builder Tasks**:
- ✅ Create domain model (`internal/core/domain/jql_alias.go`)
  - ✅ JQLAlias struct with validation
  - ✅ Predefined aliases: mine, sprint, blocked
  - ✅ Name validation (alphanumeric, hyphens, underscores)
- ✅ Implement alias repository (`internal/adapters/database/alias_repository.go`)
  - ✅ SQLite persistence with migration
  - ✅ CRUD operations (Create, GetByName, List, Update, Delete)
  - ✅ Workspace isolation support
- ✅ Implement alias service (`internal/core/services/alias_service.go`)
  - ✅ Business logic layer
  - ✅ Recursive alias expansion with @ syntax
  - ✅ Circular reference detection
  - ✅ Predefined alias fallback
- ✅ Create CLI commands (`cmd/ticketr/alias_commands.go`):
  - ✅ `ticketr alias list` - Display all aliases with type indicators
  - ✅ `ticketr alias create` - Create workspace/global aliases with descriptions
  - ✅ `ticketr alias show` - Show full alias details and expanded JQL
  - ✅ `ticketr alias update` - Update existing aliases
  - ✅ `ticketr alias delete` - Delete user-defined aliases
- ✅ Integrate with pull command (`cmd/ticketr/main.go`):
  - ✅ `ticketr pull --alias <name>` flag implementation
  - ✅ Alias expansion in pull workflow
  - ✅ Mutual exclusivity with --jql flag
  - ✅ Verbose output showing expanded JQL
- ✅ Comprehensive test coverage:
  - ✅ Domain model validation tests
  - ✅ Repository integration tests
  - ✅ Service unit tests (expansion, recursion, circular detection)
  - ✅ 100% test pass rate

**Verifier Tasks**:
- ✅ Test alias expansion correctness
- ✅ Verify predefined aliases work (mine, sprint, blocked)
- ✅ Test custom alias creation (workspace and global)
- ✅ Manual test: pull with alias - PASSED
- ✅ Coverage verified (domain, service, repository all tested)
- ✅ Recursive alias expansion verified
- ✅ Circular reference detection verified

**Acceptance**:
- ✅ Aliases expand correctly with recursive support
- ✅ Predefined aliases work out of box
- ✅ Users can define custom aliases (workspace/global)
- ✅ CLI integration complete and functional
- ⏸️ TUI integration deferred to future enhancement

**Deliverables**:
- ✅ Domain model: `jql_alias.go` (112 lines)
- ✅ Repository: `alias_repository.go` + migration SQL
- ✅ Service: `alias_service.go` (247 lines, recursive expansion)
- ✅ CLI commands: `alias_commands.go` (405 lines, 5 subcommands)
- ✅ Pull integration: expandAlias() in main.go
- ✅ Comprehensive test suite (unit + integration)
- ✅ Documentation: README.md updated + comprehensive feature guide

**Notes**:
- TUI quick filter dropdown deferred to v3.1 (not blocking for Slice 1)
- Predefined aliases implemented as in-memory constants (fast, no DB overhead)
- Workspace isolation ensures team-specific aliases don't conflict
- Recursive expansion supports unlimited nesting with cycle detection
- CLI provides helpful error messages for common issues
- Verifier status: APPROVED

**Estimated**: 8 hours | **Actual**: ~10 hours (within tolerance)

### Day 3: Performance Optimization

**Builder Tasks**:
- [ ] Profile hot paths: `go test -cpuprofile=cpu.prof`
- [ ] Optimize ticket tree rendering (if needed)
- [ ] Optimize bulk operation batching
- [ ] Optimize database queries (add indexes if needed)
- [ ] Run benchmarks before/after

**Verifier Tasks**:
- [ ] Run performance benchmarks
- [ ] Verify targets met:
  - [ ] TUI renders 1000+ tickets <100ms
  - [ ] Bulk operations batch efficiently
  - [ ] No regressions in existing operations
- [ ] Compare before/after benchmarks

**Acceptance**:
- [ ] All performance targets met
- [ ] No performance regressions
- [ ] Benchmarks documented

**Estimated**: 4 hours | **Actual**: _____

### Day 4-5: Final Polish

**Builder Tasks**:
- [ ] Review all Phase 5 error messages (improve clarity)
- [ ] Add TUI tooltips/hints where helpful
- [ ] Optimize keybindings (check for conflicts)
- [ ] Fix any outstanding minor bugs
- [ ] Code cleanup (remove debug logging, comments)

**Verifier Tasks**:
- [ ] Full regression suite on all Phase 5 features:
  - [ ] Bulk operations (CLI + TUI)
  - [ ] Templates (apply, validate)
  - [ ] Smart sync (all strategies)
  - [ ] JQL aliases (expand, use)
- [ ] Verify all existing features still work:
  - [ ] Workspace management
  - [ ] Credential profiles
  - [ ] Push/pull operations
  - [ ] TUI navigation
- [ ] Final test suite run: `go test ./... -v`
- [ ] Final coverage check: ___%

**Scribe Tasks**:
- [ ] Review all Phase 5 documentation for completeness
- [ ] Update `CHANGELOG.md` for v3.1.0 release
- [ ] Create Phase 5 completion report (`PHASE5-COMPLETE.md`)
- [ ] Update `docs/v3-implementation-roadmap.md`: Mark Phase 5 complete
- [ ] Final proofreading pass

**Acceptance**:
- [ ] Error messages improved
- [ ] TUI polish complete
- [ ] All Phase 5 features tested
- [ ] Full regression clean
- [ ] Documentation comprehensive

**Estimated**: 12 hours | **Actual**: _____

### Week 20 Final Checklist

**Pre-Commit**:
- [ ] All tests passing (final count: ___)
- [ ] Build successful
- [ ] Final coverage: ___% (≥75% overall)
- [ ] Documentation complete
- [ ] Regression testing complete
- [ ] Performance benchmarks met

**Git Commits**:
- [ ] Commit 1: `feat(jql): Add JQL alias system`
- [ ] Commit 2: `perf(phase5): Optimize hot paths for Phase 5 features`
- [ ] Commit 3: `chore(phase5): Final polish and bug fixes`
- [ ] Commit 4: `docs(phase5): Phase 5 completion report and v3.1.0 release notes`

**Week 20 Complete**: [ ] YES / [ ] NO (if NO, why: _____________)

---

## Phase 5 Completion Gate

**Prerequisites for Steward Review**:

### Code Quality
- [ ] All tests passing: `go test ./...` (0 failures)
- [ ] Build successful: `go build ./...`
- [ ] Test coverage ≥75% overall
- [ ] No P0 or P1 bugs outstanding
- [ ] Code quality standards met (gofmt, go vet clean)

### Feature Completeness
- [ ] Bulk operations: CLI + TUI functional
- [ ] Templates: YAML parser and application working
- [ ] Smart sync: All strategies implemented
- [ ] JQL aliases: Expansion and usage working
- [ ] All acceptance criteria met

### Documentation
- [ ] README.md updated with all Phase 5 features
- [ ] User guides complete:
  - [ ] docs/bulk-operations-guide.md
  - [ ] docs/template-guide.md
  - [ ] docs/sync-strategies-guide.md
- [ ] CHANGELOG.md prepared for v3.1.0
- [ ] Roadmap updated (Phase 5 complete)
- [ ] Phase 5 completion report created

### Performance
- [ ] All benchmarks met:
  - [ ] TUI rendering <100ms (1000+ tickets)
  - [ ] Bulk operations efficient
  - [ ] No regressions in existing operations
- [ ] Profiling results documented

### Regression Testing
- [ ] Full Phase 5 feature testing complete
- [ ] Existing features verified:
  - [ ] Workspace management
  - [ ] Credential profiles
  - [ ] Push/pull sync
  - [ ] TUI navigation
- [ ] Zero regressions detected

### Git Hygiene
- [ ] All work committed with proper attribution
- [ ] Conventional commit format used
- [ ] Working directory clean: `git status`
- [ ] Ready for release tagging

**Phase 5 Ready for Steward Review**: [ ] YES / [ ] NO

---

## Steward Phase Gate Review

**Delegate to Steward** (after all above complete):

```python
Task(
    subagent_type="general-purpose",
    description="Steward approval for Phase 5",
    prompt="""Steward: Review Phase 5 completion for v3.1.0 release.

Review scope:
- Architecture compliance
- Security assessment
- Requirements validation
- Production readiness

Gate requirements: docs/phase-5-gate-approval.md

Deliverables:
- Architecture compliance report
- Security assessment
- Final recommendation: APPROVE / APPROVED WITH CONDITIONS / REJECTED
"""
)
```

**Steward Decision**: [ ] APPROVED / [ ] APPROVED WITH CONDITIONS / [ ] REJECTED

**Conditions (if any)**:
1. _______________
2. _______________
3. _______________

**Remediation Plan** (if REJECTED or CONDITIONS):
_______________
_______________

---

## Post-Phase 5 Actions

**After Steward Approval**:

- [ ] Tag release: `git tag v3.1.0`
- [ ] Push to remote: `git push origin feature/v3 --tags`
- [ ] Create GitHub release with CHANGELOG.md content
- [ ] Update project README with v3.1.0 status
- [ ] Announce release (internal/external)

**Retrospective**:

- [ ] Document learnings from Phase 5
- [ ] Update DIRECTOR-ORCHESTRATION-FRAMEWORK.md with improvements
- [ ] Archive Phase 5 completion report
- [ ] Plan Phase 6 (if applicable) or v3.2.0 features

---

## Phase 5 Statistics

**Code Metrics**:
- Total LOC added: _____
- Total LOC modified: _____
- Files created: _____
- Files modified: _____

**Test Metrics**:
- Tests added: _____
- Final test count: _____
- Final coverage: _____%
- Test execution time: _____s

**Documentation Metrics**:
- Docs created: _____
- Docs updated: _____
- Total doc lines: _____

**Time Metrics**:
- Week 18 actual time: _____ hours
- Week 19 actual time: _____ hours
- Week 20 actual time: _____ hours
- Total Phase 5 time: _____ hours
- Estimate vs. actual: _____%

**Quality Metrics**:
- P0 bugs found: _____
- P1 bugs found: _____
- Regressions introduced: _____
- Regressions fixed: _____

---

## Notes and Learnings

**What Went Well**:
-
-
-

**What Could Be Improved**:
-
-
-

**Blockers Encountered**:
-
-
-

**Framework Improvements**:
-
-
-

---

**Phase 5 Status**: [ ] IN PROGRESS / [ ] COMPLETE / [ ] BLOCKED

**Completion Date**: _______________

**Next Phase**: _______________

---

*For detailed methodology, see DIRECTOR-ORCHESTRATION-FRAMEWORK.md*
