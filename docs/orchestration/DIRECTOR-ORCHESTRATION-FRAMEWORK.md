# TICKETR V3.0 DIRECTOR'S ORCHESTRATION FRAMEWORK
## The Definitive Methodology for Phase 5+ Execution

**Document Version:** 2.0
**Date:** October 19, 2025
**Status:** Production Framework
**Purpose:** Comprehensive orchestration methodology for Ticketr v3.0 Phase 5 and beyond

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [PART 1: METHODOLOGY](#part-1-methodology)
3. [PART 2: ORCHESTRATION LOOP CYCLE](#part-2-orchestration-loop-cycle)
4. [PART 3: AGENT DELEGATION MATRIX](#part-3-agent-delegation-matrix)
5. [PART 4: PHASE 5 EXECUTION ROADMAP](#part-4-phase-5-execution-roadmap)
6. [PART 5: ONBOARDING & CONTEXT](#part-5-onboarding--context)
7. [PART 6: QUALITY GATES & STANDARDS](#part-6-quality-gates--standards)
8. [PART 7: COMMUNICATION & ESCALATION](#part-7-communication--escalation)
9. [APPENDICES](#appendices)

---

## Executive Summary

This framework provides the **complete orchestration methodology** for executing Ticketr v3.0 Phase 5 (Advanced Features) and all subsequent development. It synthesizes proven patterns from Phase 2 (5,791 LOC delivered with zero regressions) and Phase 4 (4,400 LOC TUI implementation) into a systematic, repeatable process.

### Framework Goals

1. **Zero-Ambiguity Execution**: Every decision point has clear criteria and documented precedent
2. **Predictable Timelines**: Granular work breakdown enables accurate estimation
3. **Quality Assurance**: Multi-layered validation prevents regressions and technical debt
4. **Seamless Handoffs**: Structured agent communication ensures continuity
5. **Complete Traceability**: Every deliverable maps to requirements and roadmap milestones
6. **Easy Onboarding**: New Directors can execute effectively within hours, not days

### Who Should Use This Framework

- **Director Agents**: Primary orchestrators managing milestone execution
- **Human Operators**: Stakeholders reviewing progress and providing guidance
- **Future Contributors**: Developers onboarding to Ticketr methodology
- **Quality Auditors**: Reviewers validating process adherence

### How to Use This Document

1. **Read Sequentially First**: Understand the complete methodology before executing
2. **Reference During Execution**: Use checklists and templates as working tools
3. **Adapt Thoughtfully**: Framework is prescriptive but allows justified variation
4. **Update Continuously**: Document learnings and improve the framework

---

## PART 1: METHODOLOGY

### 1.1 Core Philosophy

Ticketr v3.0 development follows a **specialized agent-based methodology** where:

- **Director**: Orchestrates strategic planning and milestone execution
- **Builder**: Implements code changes and writes initial tests
- **Verifier**: Validates quality, extends test coverage, runs full suite
- **Scribe**: Creates and maintains all documentation
- **Steward**: Provides architectural oversight and final approval

**Key Principle**: **Sequential Delegation with Quality Gates**

Each milestone flows: Plan → Build → Verify → Document → (Approve) → Commit → Complete

### 1.2 Agent Roles and Responsibilities

#### Director (You)

**Mission**: Strategic orchestration and quality gatekeeping

**Responsibilities**:
- Break roadmap milestones into atomic, delegatable tasks
- Assign work to appropriate specialized agents
- Verify deliverables meet quality standards
- Enforce test coverage and documentation requirements
- Create logical git commits with proper attribution
- Track progress using TodoWrite tool
- Escalate blockers to human operator

**Authority**:
- ✅ Delegate work to any agent
- ✅ Approve or reject agent deliverables
- ✅ Request rework when quality insufficient
- ✅ Create git commits for completed work
- ✅ Update roadmap status

**Limitations**:
- ❌ Do NOT write production code (delegate to Builder)
- ❌ Do NOT write tests (delegate to Verifier)
- ❌ Do NOT write documentation (delegate to Scribe)
- ❌ Do NOT make solo architectural decisions (consult Steward)

#### Builder

**Mission**: Implement code changes with initial test coverage

**Inputs**: Requirements, file locations, acceptance criteria, patterns to follow

**Outputs**: Modified code, passing tests, implementation notes

**Quality Bar**:
- All modified code compiles successfully
- Initial tests pass (go test ./...)
- Follows hexagonal architecture patterns
- Proper error handling with context wrapping
- Code comments on public functions

**Typical Work**: Feature implementation, bug fixes, refactoring, new integrations

#### Verifier

**Mission**: Validate implementation quality and extend test coverage

**Inputs**: Builder's changes, coverage targets, requirements to validate

**Outputs**: Test results, coverage reports, regression checks, approval recommendation

**Quality Bar**:
- Full test suite passes (go test ./...)
- Coverage targets met (80%+ for critical paths)
- Zero regressions detected
- Requirements validated against acceptance criteria
- Clear APPROVE or REQUEST FIXES recommendation

**Typical Work**: Test suite validation, coverage extension, regression testing, QA

#### Scribe

**Mission**: Create and maintain comprehensive documentation

**Inputs**: Implementation summary, feature descriptions, technical details

**Outputs**: Updated docs, examples, guides, changelog entries

**Quality Bar**:
- All user-facing changes documented
- Examples tested and accurate
- Cross-references valid
- Consistent markdown formatting
- Technical accuracy verified

**Typical Work**: README updates, guide creation, API documentation, changelog maintenance

#### Steward

**Mission**: Architectural oversight and phase gate approval

**Inputs**: Completed milestone deliverables, test evidence, documentation

**Outputs**: Architectural review, security assessment, approval decision

**Quality Bar**:
- Hexagonal architecture compliance verified
- Security implications assessed
- Requirements satisfaction validated
- Technical debt evaluation
- Go/no-go decision with rationale

**Typical Work**: Phase gate reviews, major architectural decisions, production releases

### 1.3 Workflow Sequencing Rules

**Standard Milestone Workflow**:

```
1. DIRECTOR: Analyze milestone requirements
   ↓ Create TodoList with granular tasks

2. BUILDER: Implement code + initial tests
   ↓ Deliver: code, test results, notes

3. DIRECTOR: Review Builder output
   ↓ Decision: Accept (proceed) or Reject (rework)

4. VERIFIER: Validate quality + extend tests
   ↓ Deliver: full suite results, coverage, recommendation

5. DIRECTOR: Review Verifier output
   ↓ Decision: Approved (proceed) or Request Fixes (loop to Builder)

6. SCRIBE: Document features + update guides
   ↓ Deliver: doc diffs, examples, changelog

7. DIRECTOR: Review Scribe output
   ↓ Decision: Accept (proceed) or Revise (rework)

8. [OPTIONAL] STEWARD: Architectural review
   ↓ Deliver: approval, conditions, or rejection

9. DIRECTOR: Create git commit(s)
   ↓ Conventional commits with co-authorship

10. DIRECTOR: Mark milestone complete
    ↓ Update roadmap, close TodoList
```

**Critical Rules**:

1. **Never Skip Verifier**: Even if Builder reports passing tests, Verifier MUST validate
2. **Never Skip Scribe**: Documentation is mandatory for all user-facing changes
3. **Sequential Execution**: Complete current agent task before starting next
4. **Quality Gate Enforcement**: Do not proceed if quality standards unmet
5. **One Task In Progress**: TodoList must show exactly ONE task "in_progress" at a time

### 1.4 Decision Trees

#### When to Invoke Steward

**INVOKE Steward if**:
- ✅ Phase gate approval required (end of Phase 2, 3, 4, 5)
- ✅ Major architectural change (new adapter, domain model change)
- ✅ Production release preparation
- ✅ Uncertainty about design approach
- ✅ Security-sensitive changes (credential handling, state management)
- ✅ Breaking changes to public APIs

**DO NOT invoke Steward if**:
- ❌ Routine feature implementation within established patterns
- ❌ Documentation-only changes
- ❌ Test coverage improvements
- ❌ Bug fixes without architectural impact

#### When to Escalate to Human

**ESCALATE if**:
- ❗ Blocked for >30 minutes on same issue
- ❗ Conflicting requirements discovered
- ❗ Missing credentials or external access
- ❗ Agent repeatedly fails same task (3+ iterations)
- ❗ Architectural uncertainty with no clear precedent
- ❗ Security concerns identified
- ❗ Timeline at risk due to unexpected complexity

**DO NOT escalate if**:
- ✅ Standard agent rework (Builder fix, Verifier retest)
- ✅ Documentation clarifications
- ✅ Routine quality issues solvable through delegation

#### When to Create Multiple Commits

**SINGLE COMMIT if**:
- Small feature (<200 lines)
- Self-contained change
- Implementation + tests + docs all related

**MULTIPLE COMMITS if**:
- Large feature (>500 lines)
- Logical separation possible:
  - Commit 1: feat(scope): Core implementation
  - Commit 2: test(scope): Extended test coverage
  - Commit 3: docs(scope): Documentation updates
- Refactoring + feature addition (separate commits)
- Database migration + code change (separate commits)

---

## PART 2: ORCHESTRATION LOOP CYCLE

### 2.1 Pre-Execution Planning Phase

**Duration**: 10-20 minutes
**Outcome**: Complete TodoList and execution plan

#### Step 1: Context Gathering (5 min)

**Read Key Documents**:
```bash
# Core roadmap
Read: docs/v3-implementation-roadmap.md

# Requirements (if applicable)
Read: REQUIREMENTS-v2.md (search for relevant PROD-XXX IDs)

# Current state
Bash: git status
Bash: git log -5 --oneline
```

**Extract Information**:
- Next incomplete milestone number and title
- Acceptance criteria (checklist items)
- Dependencies on previous milestones
- Estimated complexity (lines of code, files affected)

#### Step 2: Codebase Analysis (5 min)

**Understand Current Implementation**:
```bash
# Find relevant files
Glob: internal/core/services/*.go
Glob: internal/adapters/**/*.go
Glob: cmd/ticketr/*.go

# Review existing patterns
Read: [files identified by glob]

# Search for related functionality
Grep: pattern="[relevant term]" output_mode="files_with_matches"
```

**Assess Impact**:
- Which services need modification?
- Which adapters affected?
- Are new domain models needed?
- Database migration required?

#### Step 3: Work Breakdown (5 min)

**Create Atomic Tasks**:

Break milestone into Builder-ready slices:

```markdown
Example: Milestone 19 "Bulk Operations"

Builder Tasks:
1. Implement BulkOperation domain model
2. Extend TicketService with BulkUpdate method
3. Create CLI command: ticketr bulk update
4. Add TUI bulk selection view

Verifier Tasks:
1. Validate BulkOperation service layer (target: 80% coverage)
2. Run full test suite with bulk operations
3. Regression test: ensure single-ticket operations unaffected

Scribe Tasks:
1. Update README with bulk operations section
2. Create docs/bulk-operations-guide.md
3. Add CHANGELOG entry
4. Update TUI help with bulk keybindings
```

**Size Validation**:
- Each Builder task: 1-3 hours of work
- Each Verifier task: 30-60 minutes
- Each Scribe task: 20-40 minutes

If tasks larger, break down further.

#### Step 4: TodoList Creation (2 min)

**Create Complete TodoList**:

```python
TodoWrite(todos=[
    {
        "content": "Analyze Milestone X requirements and dependencies",
        "activeForm": "Analyzing Milestone X requirements",
        "status": "in_progress"
    },
    {
        "content": "Delegate BulkOperation domain model implementation to Builder",
        "activeForm": "Delegating domain model implementation",
        "status": "pending"
    },
    {
        "content": "Review Builder deliverable and approve",
        "activeForm": "Reviewing Builder deliverable",
        "status": "pending"
    },
    {
        "content": "Delegate service layer validation to Verifier",
        "activeForm": "Delegating validation to Verifier",
        "status": "pending"
    },
    {
        "content": "Review Verifier results and approve",
        "activeForm": "Reviewing Verifier results",
        "status": "pending"
    },
    {
        "content": "Delegate documentation updates to Scribe",
        "activeForm": "Delegating documentation to Scribe",
        "status": "pending"
    },
    {
        "content": "Review Scribe deliverables and approve",
        "activeForm": "Reviewing documentation",
        "status": "pending"
    },
    {
        "content": "Create git commits with proper attribution",
        "activeForm": "Creating git commits",
        "status": "pending"
    },
    {
        "content": "Update roadmap and mark milestone complete",
        "activeForm": "Updating roadmap",
        "status": "pending"
    }
])
```

#### Step 5: Communicate Plan (3 min)

**Report to User**:

```markdown
# Milestone X Execution Plan

## Scope
[2-3 sentence description of what will be delivered]

## Work Breakdown
- Builder: [N tasks, estimated M hours]
- Verifier: [N tasks, estimated M minutes]
- Scribe: [N tasks, estimated M minutes]
- Steward: [If applicable]

## Timeline Estimate
Estimated completion: [X hours/days]

## Dependencies
[Any blockers or prerequisites]

## Risk Assessment
[Known risks or unknowns]

Proceeding with execution...
```

### 2.2 Execution Phase with Checkpoints

**Duration**: 2-8 hours (milestone dependent)
**Outcome**: Implemented, tested, documented feature

#### Checkpoint 1: Builder Delegation

**Template**:

```python
Task(
    subagent_type="general-purpose",
    description="Implement [feature name]",
    prompt="""You are the Builder agent for Ticketr v3.0 Milestone X.

## Mission
Implement [detailed feature description].

## Context
Current state:
- [Relevant existing code locations]
- [Patterns to follow]
- [Integration points]

Roadmap reference: docs/v3-implementation-roadmap.md (Milestone X)

## Requirements
From roadmap acceptance criteria:
1. [Criterion 1]
2. [Criterion 2]
3. [Criterion 3]

## Implementation Tasks
1. Create domain model in internal/core/domain/[model].go
2. Extend service layer in internal/core/services/[service].go
3. Add CLI command in cmd/ticketr/[command].go
4. Add initial tests in [corresponding _test.go files]

## File Locations
Files to create:
- [list new files]

Files to modify:
- [list existing files with expected changes]

## Quality Standards
- All code must compile: `go build ./...`
- All tests must pass: `go test ./...`
- Follow hexagonal architecture (ports/adapters pattern)
- Use proper error handling (fmt.Errorf with %w for wrapping)
- Add godoc comments for public functions
- Code coverage target: >70% for new code

## Expected Deliverables
Provide the following in your response:

1. **Files Modified Summary**:
   - File path
   - Lines added
   - Lines modified
   - Purpose of changes

2. **Test Results**:
   ```bash
   $ go test ./[package]
   [Exact command output]
   ```

3. **Build Verification**:
   ```bash
   $ go build ./...
   [Confirmation of success]
   ```

4. **Implementation Notes**:
   - Design decisions made
   - Patterns followed
   - Any deviations from plan (with rationale)

5. **Notes for Verifier**:
   - Areas needing additional test coverage
   - Edge cases to validate
   - Integration test suggestions

Begin implementation now.
"""
)
```

**Post-Builder Review Checklist**:

- [ ] All expected files created/modified?
- [ ] Build successful (go build ./...)?
- [ ] Tests passing (go test ./...)?
- [ ] Implementation notes clear and complete?
- [ ] No obvious code quality issues?
- [ ] Ready for Verifier?

**If NO to any**: Delegate back to Builder with specific fix requests.

**If ALL YES**: Update TodoList, proceed to Verifier.

```python
TodoWrite(todos=[
    {"content": "Analyze requirements", "activeForm": "...", "status": "completed"},
    {"content": "Delegate to Builder", "activeForm": "...", "status": "completed"},
    {"content": "Review Builder deliverable", "activeForm": "...", "status": "completed"},
    {"content": "Delegate to Verifier", "activeForm": "Delegating...", "status": "in_progress"},
    # ... rest pending
])
```

#### Checkpoint 2: Verifier Delegation

**Template**:

```python
Task(
    subagent_type="general-purpose",
    description="Verify [feature] implementation",
    prompt="""You are the Verifier agent for Ticketr v3.0 Milestone X.

## Mission
Validate the quality and completeness of [feature] implementation.

## Context
Builder has completed:
[Paste Builder's implementation summary]

Files modified:
- [list from Builder's deliverable]

## Validation Tasks

### 1. Full Test Suite Execution
Run complete test suite and report results:
```bash
go test ./... -v
```

Provide:
- Total tests run
- Pass count
- Fail count (with details if any)
- Skip count (with reason)

### 2. Coverage Analysis
Check coverage for new code:
```bash
go test ./internal/core/services/[service] -coverprofile=coverage.out
go tool cover -func=coverage.out | grep [NewMethods]
```

Report coverage percentage for:
- [Method 1]
- [Method 2]
- Overall package

### 3. Test Coverage Extension
If coverage <80% for critical methods, add tests:
- Edge cases
- Error paths
- Boundary conditions

### 4. Regression Testing
Verify no existing functionality broken:
- Run full suite: `go test ./...`
- Check for new failures in unrelated tests
- Validate backward compatibility

### 5. Requirements Validation
Map implementation to acceptance criteria:

| Criterion | Status | Evidence |
|-----------|--------|----------|
| [Criterion 1] | ✅/❌ | [How verified] |
| [Criterion 2] | ✅/❌ | [How verified] |

## Coverage Targets
- New domain models: >60%
- New service methods: >80%
- New CLI commands: >70%
- Overall package: maintain or improve from baseline

## Expected Deliverables

1. **Test Suite Results**:
   ```
   [Complete output from go test ./... -v]
   ```

2. **Coverage Report**:
   ```
   [Output from coverage analysis]
   ```

3. **Regression Check**:
   - Previous test count: [X]
   - Current test count: [Y]
   - New tests added: [Z]
   - Regressions detected: [None or details]

4. **Requirements Compliance Matrix**:
   [Table mapping criteria to evidence]

5. **Recommendation**:
   - **APPROVE**: All quality standards met, ready for Scribe
   - **REQUEST FIXES**: [Specific issues that must be addressed]

If REQUEST FIXES, provide:
- Exact failures or gaps
- Suggested fixes
- Priority (blocking vs. nice-to-have)

Begin verification now.
"""
)
```

**Post-Verifier Review Decision Tree**:

```
Verifier Recommendation = APPROVE?
├─ YES → Update TodoList, proceed to Scribe
└─ NO → Verifier provided REQUEST FIXES
    ├─ Analyze issues listed
    ├─ Delegate back to Builder with fixes
    ├─ Loop: Builder → Verifier until APPROVED
    └─ Do NOT proceed to Scribe until APPROVED
```

**Critical Rule**: **NEVER proceed to Scribe with failing tests or unmet coverage targets.**

#### Checkpoint 3: Scribe Delegation

**Template**:

```python
Task(
    subagent_type="general-purpose",
    description="Document [feature]",
    prompt="""You are the Scribe agent for Ticketr v3.0 Milestone X.

## Mission
Create comprehensive documentation for [feature].

## Context
Implementation complete and verified:

Builder Summary:
[Paste Builder's summary]

Verifier Results:
- Test count: [X]
- Coverage: [Y%]
- Status: APPROVED

## Documentation Tasks

### 1. README.md Updates
Update the main README with:
- Features section: Add [feature] with brief description
- Quick Start: Include example command if user-facing
- Common Commands: Add to reference table

### 2. User Guide Creation/Updates
File: docs/[feature]-guide.md

Sections to include:
- Introduction: What is [feature]?
- Getting Started: Step-by-step first use
- Command Reference: All commands with examples
- Use Cases: Common workflows
- Troubleshooting: Anticipated issues
- Best Practices: Team usage patterns

### 3. Technical Documentation
Update docs/ARCHITECTURE.md if:
- New domain models added
- New adapters created
- Ports/interfaces modified

### 4. Roadmap Updates
File: docs/v3-implementation-roadmap.md
- [x] Mark Milestone X acceptance criteria complete
- Update status to "✅ Complete"

### 5. Changelog Entry
File: CHANGELOG.md

Add to [Unreleased] section:
```markdown
### Added
- [Feature]: [Description] (Milestone X)
  - CLI command: `ticketr [command]`
  - [Key capability 1]
  - [Key capability 2]

### Changed
- [If anything changed]

### Technical
- [Implementation details for developers]
- Test coverage: [X]% for new code
```

## Quality Standards
- All code examples tested and accurate
- Cross-references valid (no broken links)
- Consistent markdown formatting:
  - Code blocks with language hints (```bash, ```go, etc.)
  - Tables properly formatted
  - Headers hierarchical (##, ###, ####)
- Command examples use realistic values
- Screenshots/diagrams if helpful (TUI features)

## Expected Deliverables

1. **Files Modified Summary**:
   | File | Lines Added | Lines Modified | Purpose |
   |------|-------------|----------------|---------|
   | README.md | [X] | [Y] | [Description] |
   | docs/[guide].md | [X] | [Y] | [Description] |
   | ... | ... | ... | ... |

2. **Cross-Reference Validation**:
   - Internal links tested: [list]
   - External references checked: [list]
   - All valid: ✅

3. **Example Verification**:
   - Commands tested: [list with results]
   - All examples working: ✅

4. **Documentation Diff Summary**:
   Brief description of what changed in each file.

Begin documentation now.
"""
)
```

**Post-Scribe Review Checklist**:

- [ ] README.md updated with feature?
- [ ] User guide comprehensive and accurate?
- [ ] Examples tested and working?
- [ ] Roadmap milestone marked complete?
- [ ] Changelog entry prepared?
- [ ] Cross-references valid?
- [ ] Markdown formatting consistent?

**If NO to any**: Provide Scribe with specific revision requests.

**If ALL YES**: Update TodoList, proceed to Git Commit.

#### Checkpoint 4: Git Commit Creation

**Commit Strategy**:

For most milestones, create **2-3 logical commits**:

1. **Implementation Commit**:
   ```bash
   git add [implementation files]
   git commit -m "$(cat <<'EOF'
   feat(scope): Brief description of feature

   Detailed explanation of what was implemented.

   Implementation:
   - Key component 1
   - Key component 2
   - Integration with [existing system]

   Testing:
   - [X] tests passing
   - [Y]% coverage for new code

   Files modified:
   - [file 1] (+N lines)
   - [file 2] (+N lines)

   🤖 Generated with [Claude Code](https://claude.com/claude-code)
   via [Happy](https://happy.engineering)

   Co-Authored-By: Claude <noreply@anthropic.com>
   Co-Authored-By: Happy <yesreply@happy.engineering>
   EOF
   )"
   ```

2. **Documentation Commit**:
   ```bash
   git add [documentation files]
   git commit -m "$(cat <<'EOF'
   docs(scope): Add comprehensive [feature] documentation

   Complete user-facing documentation for [feature].

   Documentation added:
   - README.md: Feature highlight
   - docs/[guide].md: Complete user guide
   - CHANGELOG.md: Release notes prepared

   Cross-references validated, examples tested.

   🤖 Generated with [Claude Code](https://claude.com/claude-code)
   via [Happy](https://happy.engineering)

   Co-Authored-By: Claude <noreply@anthropic.com>
   Co-Authored-By: Happy <yesreply@happy.engineering>
   EOF
   )"
   ```

**Conventional Commit Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `test`: Test additions/changes
- `refactor`: Code refactoring
- `perf`: Performance improvement
- `chore`: Maintenance

**Scope Examples**: `workspace`, `tui`, `cli`, `database`, `keychain`, `sync`

### 2.3 Post-Execution Validation

**Duration**: 5-10 minutes
**Outcome**: Verified completion and clean state

#### Final Validation Checklist

```bash
# 1. All tests passing
go test ./... -v

# 2. Build successful
go build ./...

# 3. Working directory clean
git status

# 4. Commits created
git log -3 --oneline

# 5. Roadmap updated
grep "Milestone X" docs/v3-implementation-roadmap.md
```

#### Milestone Completion Report

**Create Summary**:

```markdown
# Milestone X Complete

## Deliverables
- ✅ [Feature 1]: [Files modified, lines added]
- ✅ [Feature 2]: [Files modified, lines added]
- ✅ Tests: [X tests passing, Y% coverage]
- ✅ Documentation: [Z files updated]

## Test Evidence
[Paste test results]

## Commits Created
- [commit hash]: feat(scope): [description]
- [commit hash]: docs(scope): [description]

## Next Steps
- Milestone X+1: [Brief description]
```

#### Update TodoList

```python
# Mark all tasks complete
TodoWrite(todos=[
    {"content": "Analyze requirements", "activeForm": "...", "status": "completed"},
    {"content": "Delegate to Builder", "activeForm": "...", "status": "completed"},
    {"content": "Review Builder", "activeForm": "...", "status": "completed"},
    {"content": "Delegate to Verifier", "activeForm": "...", "status": "completed"},
    {"content": "Review Verifier", "activeForm": "...", "status": "completed"},
    {"content": "Delegate to Scribe", "activeForm": "...", "status": "completed"},
    {"content": "Review Scribe", "activeForm": "...", "status": "completed"},
    {"content": "Create commits", "activeForm": "...", "status": "completed"},
    {"content": "Update roadmap", "activeForm": "...", "status": "completed"}
])
```

### 2.4 Feedback Incorporation

**After Each Milestone**:

1. **Capture Learnings**: What went well? What was challenging?
2. **Update Framework**: Document new patterns or anti-patterns
3. **Refine Estimates**: Adjust future task sizing based on actuals
4. **Communicate**: Share insights with human operator

**Continuous Improvement**:
- Update this framework document with proven patterns
- Add to appendices: examples of good agent responses
- Document new decision criteria
- Refine quality standards based on production experience

---

## PART 3: AGENT DELEGATION MATRIX

### 3.1 Task → Agent Mapping

| Task Type | Agent | Rationale | Example |
|-----------|-------|-----------|---------|
| New feature implementation | Builder | Code creation expertise | Implement BulkOperation service |
| Bug fixes | Builder | Code modification | Fix workspace switching persistence |
| Refactoring | Builder | Code restructuring | Extract common validation logic |
| Test creation | Verifier | Quality assurance focus | Add edge case tests for profiles |
| Coverage improvement | Verifier | Testing expertise | Increase service layer to 80% |
| Regression testing | Verifier | Validation specialization | Ensure Phase 4 features still work |
| User documentation | Scribe | Writing expertise | Create bulk operations guide |
| API documentation | Scribe | Technical writing | Document CredentialStore interface |
| Changelog maintenance | Scribe | Release tracking | Prepare v3.1.0 release notes |
| Architecture review | Steward | Design expertise | Review hexagonal compliance |
| Security assessment | Steward | Security focus | Validate keychain integration |
| Phase gate approval | Steward | Strategic oversight | Approve Phase 5 completion |

### 3.2 Decision Matrix: Builder vs. General-Purpose

**Use Builder (via Task tool) when**:
- ✅ Writing new code
- ✅ Modifying existing code
- ✅ Implementing tests
- ✅ Following established patterns
- ✅ Work scope clear and bounded

**Use General-Purpose agent when**:
- ⚠ Never use general-purpose for Ticketr development
- ⚠ All work should go through specialized agents
- ⚠ Director is the only general-purpose agent in this workflow

**Rationale**: Specialized agents have context-specific prompts that ensure quality and consistency.

### 3.3 Verifier Involvement Triggers

**ALWAYS invoke Verifier if**:
- ✅ Builder completed code changes
- ✅ Tests need validation or extension
- ✅ Requirements need compliance check
- ✅ Phase gate approaching

**OPTIONAL Verifier for**:
- Documentation-only changes (no code impact)
- Roadmap updates (no implementation)

**Verifier provides**:
- Full test suite execution results
- Coverage metrics
- Regression detection
- Requirements validation
- APPROVE/REQUEST FIXES recommendation

### 3.4 Scribe Involvement Requirements

**MANDATORY Scribe invocation if**:
- ✅ User-facing feature added or changed
- ✅ CLI command added or modified
- ✅ TUI functionality changed
- ✅ API/interface changed
- ✅ Milestone completed

**MANDATORY Scribe deliverables**:
- ✅ README.md updates (if applicable)
- ✅ User guide updates (if applicable)
- ✅ Roadmap milestone checkboxes
- ✅ CHANGELOG.md entry

**OPTIONAL Scribe updates**:
- Internal code comments (Builder handles)
- Test documentation (Verifier handles)

### 3.5 Steward Escalation Criteria

**Phase Gate Reviews** (REQUIRED):
- Phase 2 completion: Workspace model
- Phase 3 completion: Global installation
- Phase 4 completion: TUI implementation
- Phase 5 completion: Advanced features

**Major Architectural Changes** (REQUIRED):
- New adapter creation
- Domain model significant changes
- Database schema evolution
- Security model changes

**Uncertainty Resolution** (OPTIONAL):
- Multiple valid design approaches
- Trade-off decisions needed
- Performance vs. maintainability

**Production Releases** (REQUIRED):
- v3.0.0 GA release
- Major version bumps

### 3.6 Research Delegation

**When to delegate research**:
- Investigating new libraries or frameworks
- Analyzing best practices from other projects
- Understanding complex technical concepts
- Evaluating alternative approaches

**Research Agent Template**:

```python
Task(
    subagent_type="general-purpose",
    description="Research [topic]",
    prompt="""Research [topic] for Ticketr v3.0.

## Research Questions
1. [Question 1]
2. [Question 2]
3. [Question 3]

## Deliverables
- Summary of findings
- Comparison of approaches
- Recommendation with rationale
- Implementation implications

Provide comprehensive research report.
"""
)
```

**Use Research sparingly**: Most work follows established patterns.

---

## PART 4: PHASE 5 EXECUTION ROADMAP

### 4.1 Phase 5 Overview

**Goal**: Add power-user features leveraging v3.0 architecture

**Duration**: Weeks 18-20 (3 weeks)

**Prerequisites**:
- ✅ Phase 4 TUI complete (Week 16)
- ✅ Milestone 18 complete (Week 17)
- ✅ All critical bugs fixed

**Success Criteria**:
- Bulk operations functional in CLI and TUI
- Template system reduces ticket creation time by 50%
- Smart sync prevents data loss during conflicts
- JQL aliases work seamlessly
- Zero regressions in existing functionality

### 4.2 Week 18: Bulk Operations

**Milestone**: Bulk Operations in CLI and TUI

**Features**:
1. Select multiple tickets (TUI checkboxes, CLI ID list)
2. Bulk update (change status, assignee, custom fields)
3. Bulk move (change project, parent)
4. Bulk delete (with confirmation)

#### Week 18 Slice 1: Domain Model (Day 1)

**Builder Task**:
```markdown
Create BulkOperation domain model and service interface.

Files to create:
- internal/core/domain/bulk_operation.go
- internal/core/ports/bulk_operation_service.go

Requirements:
- BulkOperation struct with Action, TicketIDs, Changes fields
- Validation: max 100 tickets per operation
- Support operations: update, move, delete
```

**Acceptance Criteria**:
- [ ] Domain model compiles
- [ ] Validation logic tested
- [ ] Interface defined

**Estimated Time**: 2 hours

#### Week 18 Slice 2: Service Implementation (Day 2)

**Builder Task**:
```markdown
Implement BulkOperationService with Jira API integration.

Files to create:
- internal/core/services/bulk_operation_service.go
- internal/core/services/bulk_operation_service_test.go

Requirements:
- BatchUpdate() method: update multiple tickets
- Transaction safety: rollback on partial failure
- Progress callback for CLI/TUI
```

**Acceptance Criteria**:
- [ ] Service compiles and tests pass
- [ ] >80% test coverage
- [ ] Handles partial failures gracefully

**Estimated Time**: 4 hours

#### Week 18 Slice 3: CLI Integration (Day 3)

**Builder Task**:
```markdown
Add CLI commands for bulk operations.

Files to create:
- cmd/ticketr/bulk_commands.go
- cmd/ticketr/bulk_commands_test.go

Commands:
- ticketr bulk update --ids PROJ-1,PROJ-2 --set status=Done
- ticketr bulk move --ids PROJ-1,PROJ-2 --parent PROJ-100
- ticketr bulk delete --ids PROJ-1,PROJ-2 --confirm
```

**Acceptance Criteria**:
- [ ] All commands functional
- [ ] Confirmation prompts for destructive operations
- [ ] Progress indicators during execution

**Estimated Time**: 3 hours

#### Week 18 Slice 4: TUI Integration (Day 4-5)

**Builder Task**:
```markdown
Add bulk selection and operations to TUI.

Files to modify:
- internal/adapters/tui/views/ticket_tree.go
- internal/adapters/tui/app.go

Features:
- Space bar: toggle ticket selection
- 'a': select all visible
- 'u': bulk update modal
- Visual indicators (checkboxes) for selected tickets
```

**Acceptance Criteria**:
- [ ] Multi-select works smoothly
- [ ] Bulk operations accessible via keybindings
- [ ] Status feedback after operations

**Estimated Time**: 6 hours

#### Week 18 Documentation (Day 5)

**Scribe Task**:
```markdown
Document bulk operations.

Files to update:
- README.md: Add bulk operations to features
- docs/bulk-operations-guide.md (new): Complete guide
- CHANGELOG.md: Prepare Week 18 entry
```

**Acceptance Criteria**:
- [ ] Examples tested and accurate
- [ ] TUI keybindings documented
- [ ] Safety warnings included

**Estimated Time**: 2 hours

### 4.3 Week 19: Templates + Smart Sync

**Milestone**: Template System and Conflict Resolution

#### Week 19 Slice 1: Template System (Day 1-2)

**Features**:
1. Template definitions (YAML format)
2. Variable substitution ({{.Name}}, {{.Sprint}})
3. Hierarchical templates (Epic → Stories → Tasks)
4. CLI command: ticketr template apply feature.yaml

**Builder Tasks**:
```markdown
1. Create template parser (internal/templates/parser.go)
2. Implement variable substitution engine
3. Add template validation
4. Create CLI command: ticketr template apply
5. Add TUI template selector
```

**Acceptance Criteria**:
- [ ] Templates parse YAML correctly
- [ ] Variable substitution works
- [ ] Validation prevents invalid templates
- [ ] CLI and TUI integration functional

**Estimated Time**: 12 hours

#### Week 19 Slice 2: Smart Sync (Day 3-5)

**Features**:
1. Sync strategies (SyncStrategy interface)
2. Conflict resolution policies (local-wins, remote-wins, manual)
3. Three-way merge for compatible changes
4. Conflict UI in TUI

**Builder Tasks**:
```markdown
1. Define SyncStrategy interface
2. Implement LocalWinsStrategy, RemoteWinsStrategy
3. Implement ThreeWayMergeStrategy
4. Update PullService to use strategies
5. Add TUI conflict resolution modal
```

**Acceptance Criteria**:
- [ ] Strategies implemented and tested
- [ ] No data loss in any conflict scenario
- [ ] User can choose strategy via config or flag
- [ ] TUI shows clear conflict information

**Estimated Time**: 15 hours

#### Week 19 Documentation (Day 5)

**Scribe Task**:
```markdown
Document templates and smart sync.

Files to create/update:
- docs/template-guide.md (new)
- docs/sync-strategies-guide.md (new)
- README.md: Add template and sync features
- CHANGELOG.md: Week 19 entry
```

**Estimated Time**: 3 hours

### 4.4 Week 20: JQL Aliases + Polish

**Milestone**: JQL Aliases and Production Polish

#### Week 20 Slice 1: JQL Aliases (Day 1-2)

**Features**:
1. Alias definitions in .ticketr.yaml
2. Alias expansion in pull commands
3. Predefined aliases (mine, sprint, blocked)
4. CLI: ticketr alias list, ticketr pull --alias mine

**Builder Tasks**:
```markdown
1. Add alias section to config schema
2. Implement alias expansion in JiraAdapter
3. Create CLI alias management commands
4. Add TUI quick filter using aliases
```

**Acceptance Criteria**:
- [ ] Aliases expand correctly
- [ ] Predefined aliases work out of box
- [ ] Users can define custom aliases
- [ ] TUI integrates aliases seamlessly

**Estimated Time**: 8 hours

#### Week 20 Slice 2: Performance Optimization (Day 3)

**Builder Tasks**:
```markdown
Optimize hot paths identified in profiling.

Focus areas:
- Ticket tree rendering (TUI)
- Bulk operations batching
- Database query optimization
```

**Acceptance Criteria**:
- [ ] TUI renders 1000+ tickets <100ms
- [ ] Bulk operations batch efficiently
- [ ] No performance regressions

**Estimated Time**: 4 hours

#### Week 20 Slice 3: Final Polish (Day 4-5)

**Builder Tasks**:
```markdown
Address minor UX issues and polish.

Tasks:
- Improve error messages
- Add tooltips/hints in TUI
- Optimize keybindings
- Fix any outstanding bugs
```

**Verifier Tasks**:
```markdown
Full regression suite on all Phase 5 features.

Validate:
- Bulk operations
- Templates
- Smart sync
- JQL aliases
- All existing features still work
```

**Scribe Tasks**:
```markdown
Final documentation pass.

Tasks:
- Review all Phase 5 docs for completeness
- Update CHANGELOG for v3.1.0 release
- Create Phase 5 completion report
- Update roadmap with Phase 5 complete
```

**Estimated Time**: 12 hours

### 4.5 Phase 5 Completion Checklist

**Before Steward Review**:

- [ ] All Week 18 features implemented and tested
- [ ] All Week 19 features implemented and tested
- [ ] All Week 20 features implemented and tested
- [ ] Full test suite passing (0 failures)
- [ ] Test coverage >75% overall
- [ ] All documentation complete
- [ ] CHANGELOG prepared for v3.1.0
- [ ] Roadmap updated with Phase 5 complete
- [ ] Zero P0 bugs outstanding
- [ ] Performance benchmarks met

**Steward Review Focus**:
- Architecture compliance
- Security assessment
- Production readiness
- Go/no-go decision for v3.1.0 release

---

## PART 5: ONBOARDING & CONTEXT

### 5.1 Project Overview

**Ticketr**: Bidirectional sync tool between Markdown files and Jira

**Vision**: Transform from directory-bound tool to global work platform

**Current State** (October 2025):
- Version: v3.0 (in development)
- Phase 4: TUI complete (4,400 LOC)
- Milestone 18: Credential profiles complete
- Ready for: Phase 5 advanced features

**Key Differentiators**:
- Markdown-first workflow (familiar to developers)
- State-aware sync (skip unchanged tickets)
- Hexagonal architecture (testable, maintainable)
- TUI interface (terminal-native UX)

### 5.2 Architecture Primer

**Pattern**: Hexagonal (Ports & Adapters)

**Core Layers**:

```
CLI/TUI (cmd/ticketr, internal/adapters/tui)
    ↓
Adapters (internal/adapters)
    ├─ Jira API
    ├─ SQLite Database
    ├─ OS Keychain
    └─ Filesystem
    ↓
Ports (internal/core/ports)
    └─ Interface definitions
    ↓
Services (internal/core/services)
    └─ Business logic
    ↓
Domain (internal/core/domain)
    └─ Core models
```

**Key Principles**:
1. **Domain Independence**: Core logic has no external dependencies
2. **Dependency Inversion**: Adapters depend on ports, not vice versa
3. **Testability**: Mock adapters via interfaces
4. **Flexibility**: Swap adapters without touching core

**Example Flow** (Push Command):
```
User: ticketr push tickets.md
  → CLI parses command
  → FilesystemAdapter reads tickets.md
  → Parser converts Markdown → Domain models
  → PushService (core) orchestrates
  → JiraAdapter (via JiraPort interface) pushes to Jira
  → StateManager updates .ticketr.state
```

### 5.3 Current Codebase State

**Statistics** (as of Milestone 18):
- Total LOC: ~15,000
- Test Coverage: 74.8% overall
- Files: ~150 Go files
- Tests: 450+ passing

**Key Directories**:

```
ticketr/
├── cmd/ticketr/                 # CLI entry points
│   ├── main.go                  # Root command
│   ├── workspace_commands.go    # Workspace management
│   ├── credentials_commands.go  # Credential profiles
│   ├── push.go, pull.go         # Sync operations
│   └── tui_command.go           # TUI launcher
│
├── internal/
│   ├── core/
│   │   ├── domain/              # Ticket, Workspace, CredentialProfile models
│   │   ├── ports/               # Interface definitions
│   │   ├── services/            # Business logic
│   │   └── validation/          # Validation rules
│   │
│   ├── adapters/
│   │   ├── jira/                # Jira REST API integration
│   │   ├── database/            # SQLite repository
│   │   ├── keychain/            # OS credential storage
│   │   ├── filesystem/          # Markdown file I/O
│   │   └── tui/                 # Terminal UI (tview)
│   │
│   ├── parser/                  # Markdown parsing
│   ├── renderer/                # Markdown generation
│   └── state/                   # State tracking (.ticketr.state)
│
└── docs/                        # Documentation
    ├── v3-implementation-roadmap.md
    ├── workspace-management-guide.md
    └── ARCHITECTURE.md
```

**Database Schema** (v3, SQLite):
```sql
workspaces (id, name, jira_url, project_key, ...)
credential_profiles (id, name, jira_url, username, keychain_ref, ...)
tickets (id, workspace_id, jira_id, content, ...)
```

**State File** (.ticketr.state):
```json
{
  "PROJ-123": {
    "local_hash": "abc123...",
    "remote_hash": "def456..."
  }
}
```

### 5.4 Key Files and Their Purpose

| File | Purpose | Importance |
|------|---------|------------|
| `docs/v3-implementation-roadmap.md` | Phase/milestone planning | ⭐⭐⭐ Director's primary reference |
| `docs/DIRECTOR-HANDBOOK.md` | Agent methodology guide | ⭐⭐⭐ Workflow instructions |
| `REQUIREMENTS-v2.md` | Product requirements | ⭐⭐ Feature validation |
| `docs/ARCHITECTURE.md` | System design | ⭐⭐ Architectural decisions |
| `internal/core/services/workspace_service.go` | Workspace orchestration | ⭐⭐ Core business logic |
| `internal/adapters/tui/app.go` | TUI main application | ⭐⭐ UI integration |
| `cmd/ticketr/main.go` | CLI entry point | ⭐ Command wiring |

### 5.5 Common Patterns

#### Error Handling

```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("create workspace: %w", err)
}

// Log errors before returning
if err := service.DoSomething(); err != nil {
    logger.Error("Failed to do something: %v", err)
    return err
}
```

#### Dependency Injection

```go
// Constructor pattern for services
func NewWorkspaceService(
    repo ports.WorkspaceRepository,
    credStore ports.CredentialStore,
) *WorkspaceService {
    return &WorkspaceService{
        repo: repo,
        credStore: credStore,
    }
}

// Wire in main.go or test setup
service := services.NewWorkspaceService(
    sqliteRepo,
    keychainStore,
)
```

#### Repository Pattern

```go
// Interface in ports package
type WorkspaceRepository interface {
    Create(workspace *domain.Workspace) error
    Get(id string) (*domain.Workspace, error)
    List() ([]*domain.Workspace, error)
}

// Implementation in adapters/database
type SQLiteAdapter struct {
    db *sql.DB
}

func (s *SQLiteAdapter) Create(workspace *domain.Workspace) error {
    // SQL implementation
}
```

#### Testing Pattern

```go
// Table-driven tests
func TestWorkspaceService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   WorkspaceConfig
        wantErr bool
    }{
        {"valid workspace", validConfig, false},
        {"invalid name", invalidConfig, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}

// Mock adapters
type mockRepository struct {
    createFunc func(*domain.Workspace) error
}

func (m *mockRepository) Create(w *domain.Workspace) error {
    return m.createFunc(w)
}
```

### 5.6 Anti-Patterns to Avoid

**DON'T**:
- ❌ Put business logic in adapters
- ❌ Import adapters in domain/services
- ❌ Use hardcoded file paths (use PathResolver)
- ❌ Store credentials in database directly
- ❌ Skip error context wrapping
- ❌ Commit without tests
- ❌ Commit without documentation updates

**DO**:
- ✅ Keep domain pure (no external deps)
- ✅ Use interfaces for all external calls
- ✅ Test with mocks
- ✅ Wrap all errors
- ✅ Document public APIs
- ✅ Follow conventional commit format

### 5.7 Testing Expectations

**Unit Tests**:
- Co-located with implementation (`*_test.go`)
- Mock external dependencies
- Table-driven for multiple scenarios
- Coverage target: >70% for new code

**Integration Tests**:
- Located in `tests/integration/`
- Use real adapters (SQLite, filesystem)
- Validate end-to-end workflows
- Run in CI/CD pipeline

**Manual Testing**:
- TUI workflows (visual validation)
- CLI commands (example verification)
- Cross-platform (macOS, Linux, Windows)

**Performance Benchmarks**:
- Startup time: <100ms
- 1000 ticket query: <100ms
- TUI rendering: <16ms (60fps)

### 5.8 Documentation Standards

**Markdown Formatting**:
```markdown
## Second-level heading

Brief introduction paragraph.

### Third-level heading

- Bullet points for lists
- Clear, concise language
- Examples where helpful

#### Fourth-level heading (rarely used)

Code blocks with language hints:
```bash
ticketr workspace create backend
```

Tables for structured data:
| Column 1 | Column 2 |
|----------|----------|
| Value    | Value    |
```

**Code Examples**:
- Always test examples before committing
- Use realistic values (not foo/bar)
- Include expected output when helpful
- Keep examples focused (one concept per example)

**Cross-References**:
- Relative links: `[Architecture](ARCHITECTURE.md)`
- Section links: `[Overview](#overview)`
- Validate all links before committing

---

## PART 6: QUALITY GATES & STANDARDS

### 6.1 Test Coverage Requirements

**Component-Level Targets**:

| Component | Minimum | Target | Current (M18) |
|-----------|---------|--------|---------------|
| Repository Layer | 80% | 90% | 80.6% ✅ |
| Service Layer | 70% | 85% | 75.2% ✅ |
| Domain Models | 60% | 80% | 68.4% ✅ |
| Adapters | 70% | 85% | 73.1% ✅ |
| Overall | 70% | 80% | 74.8% ✅ |

**Critical Path Requirements**:
- 100% coverage for critical business logic
- All error paths tested
- Edge cases validated

**Coverage Validation**:

```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out

# View by package
go tool cover -func=coverage.out

# HTML visualization
go tool cover -html=coverage.out -o coverage.html

# Fail if below threshold
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total | awk '{if ($3+0 < 70.0) exit 1}'
```

**Verifier Checklist**:
- [ ] Overall coverage ≥70%
- [ ] New service methods ≥80%
- [ ] Critical paths 100%
- [ ] No coverage regressions

### 6.2 Performance Benchmarks

**Response Time Targets**:

| Operation | Target | Measured (M18) | Status |
|-----------|--------|----------------|--------|
| Startup time | <100ms | ~80ms | ✅ |
| Workspace switch | <50ms | ~30ms | ✅ |
| 1000 ticket query | <100ms | ~85ms | ✅ |
| TUI rendering | <16ms | ~12ms | ✅ |
| Workspace creation | <200ms | ~150ms | ✅ |

**Benchmark Tests**:

```go
func BenchmarkWorkspaceService_Switch(b *testing.B) {
    service := setupService()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = service.Switch("test-workspace")
    }
}

// Run benchmarks
// go test -bench=. -benchmem ./internal/core/services
```

**Performance Regression Detection**:
- Compare benchmarks before/after changes
- Flag regressions >10%
- Investigate and optimize or justify

### 6.3 Code Quality Standards

**Static Analysis**:

```bash
# Format check
gofmt -l . | grep -v vendor

# Vet
go vet ./...

# Staticcheck (if available)
staticcheck ./...
```

**Verifier Checklist**:
- [ ] `go vet` clean
- [ ] All code formatted (`gofmt`)
- [ ] No TODO/FIXME without tracking issue
- [ ] Public functions documented

**Code Review Standards**:
- Clear variable names (no single-letter except loops)
- Functions <50 lines (prefer smaller)
- Cyclomatic complexity <10
- No duplicate code (DRY principle)

### 6.4 Documentation Completeness

**Mandatory Documentation**:

For every user-facing change:
- [ ] README.md updated
- [ ] User guide created or updated
- [ ] CHANGELOG.md entry prepared
- [ ] Examples tested and accurate

For every architectural change:
- [ ] ARCHITECTURE.md updated
- [ ] Design decisions documented
- [ ] Migration path described

For every API change:
- [ ] Godoc comments added
- [ ] Interface contract documented
- [ ] Breaking changes highlighted

**Scribe Checklist**:
- [ ] All files listed in deliverables
- [ ] Cross-references validated
- [ ] Examples tested manually
- [ ] Markdown rendered correctly

### 6.5 Security Standards

**Credential Handling**:
- [ ] Never store credentials in database
- [ ] Always use OS keychain
- [ ] Redact credentials from logs
- [ ] Validate credentials before storing

**Input Validation**:
- [ ] Validate user input (CLI flags, TUI forms)
- [ ] Sanitize data before SQL queries (use parameterized queries)
- [ ] Validate Jira responses

**Security Review Checklist** (Steward):
- [ ] Credentials stored securely
- [ ] No secrets in logs
- [ ] No SQL injection vectors
- [ ] No XSS in TUI (if rendering user input)
- [ ] Proper error handling (no stack traces to user)

### 6.6 Rollback Procedures

**Git Rollback**:

```bash
# Revert last commit
git revert HEAD

# Revert specific commit
git revert <commit-hash>

# Hard reset (use cautiously)
git reset --hard HEAD~1
```

**Database Rollback**:

```sql
-- Each migration includes down script
-- Example: 003_credential_profiles_down.sql
DROP TABLE IF EXISTS credential_profiles;
ALTER TABLE workspaces DROP COLUMN credential_profile_id;
```

**Feature Flag Rollback**:

```yaml
# .ticketr.yaml
features:
  enable_bulk_operations: false  # Disable feature
  enable_templates: false
```

**Rollback Decision Criteria**:
- P0 bug discovered in production
- Performance regression >20%
- Data loss risk identified
- Security vulnerability found

---

## PART 7: COMMUNICATION & ESCALATION

### 7.1 Director → User Communication

**Progress Updates**:

Provide concise updates at key milestones:

```markdown
# Milestone X Progress Update

## Status: [In Progress / Blocked / Complete]

## Work Completed
- ✅ Builder implementation complete (500 lines added)
- ✅ Verifier approval received (450 tests passing)
- 🔄 Scribe documentation in progress

## Next Steps
- Complete documentation (ETA: 30 minutes)
- Create git commits
- Mark milestone complete

## Blockers
[None / List specific blockers]

## ETA
Expected completion: [Time estimate]
```

**Frequency**:
- Start of milestone: Execution plan
- After Builder: Implementation complete
- After Verifier: Quality validated
- After Scribe: Documentation complete
- End of milestone: Completion report

**Escalation Communication**:

```markdown
# ESCALATION: [Issue Summary]

## Context
[What you were trying to accomplish]

## Problem
[What went wrong]

## Attempts Made
1. [Approach 1 - Result]
2. [Approach 2 - Result]
3. [Approach 3 - Result]

## Current Blocker
[Specific blocker preventing progress]

## Request
[Specific help needed from human]

## Impact
Milestone X delayed by [time estimate] until resolved.
```

### 7.2 Builder → Director Communication

**Expected Deliverable Format**:

```markdown
# Builder Deliverable: [Feature Name]

## Implementation Summary
[2-3 sentence summary of what was built]

## Files Modified
| File | Lines Added | Lines Modified | Purpose |
|------|-------------|----------------|---------|
| internal/core/domain/feature.go | 120 | 0 | Domain model |
| internal/core/services/feature_service.go | 250 | 15 | Business logic |
| cmd/ticketr/feature_commands.go | 180 | 0 | CLI commands |

Total: 550 lines added, 15 modified

## Test Results
```bash
$ go test ./internal/core/services
=== RUN   TestFeatureService_Create
--- PASS: TestFeatureService_Create (0.01s)
=== RUN   TestFeatureService_Update
--- PASS: TestFeatureService_Update (0.01s)
PASS
ok      github.com/karolswdev/ticketr/internal/core/services  0.156s
```

## Build Verification
```bash
$ go build ./...
[No output = success]
```

## Design Decisions
- Chose approach X over Y because [rationale]
- Followed pattern from [existing component]
- Deferred optimization until profiling confirms bottleneck

## Notes for Verifier
- Service layer coverage: estimated 75%
- Edge cases to validate: [list]
- Integration test suggestions: [list]

## Outstanding Issues
[None / List any known issues]
```

**Quality Indicators**:
- ✅ Specific file paths and line counts
- ✅ Actual test output (not paraphrased)
- ✅ Build verification shown
- ✅ Design rationale provided
- ✅ Handoff notes for Verifier

### 7.3 Verifier → Director Communication

**Expected Deliverable Format**:

```markdown
# Verifier Report: [Feature Name]

## Test Suite Results

Full suite execution:
```bash
$ go test ./... -v
[Complete output]

Summary:
- Total tests: 450
- Passed: 450
- Failed: 0
- Skipped: 3 (Jira integration tests - expected)
```

## Coverage Analysis
```bash
$ go test ./internal/core/services -coverprofile=coverage.out
$ go tool cover -func=coverage.out | grep Feature

FeatureService.Create    85.7%
FeatureService.Update    90.0%
FeatureService.Delete    75.0%
Package total           82.3%
```

## Regression Check
- Previous test count: 445
- Current test count: 450 (+5 new tests)
- Regressions detected: None
- New tests cover: edge cases, error paths

## Requirements Validation

| Requirement | Status | Evidence |
|-------------|--------|----------|
| Feature must create entities | ✅ | TestFeatureService_Create passes |
| Feature must validate input | ✅ | TestFeatureService_ValidationErrors passes |
| Feature must handle conflicts | ✅ | TestFeatureService_ConflictResolution passes |

## Additional Tests Added
- Error path: invalid input → error returned
- Edge case: empty list → no-op
- Boundary: max 100 items → validation error

## Performance Validation
- Benchmark: BenchmarkFeatureService_Create: 1.2ms average
- Within target (<5ms)

## Recommendation
**APPROVE** - All quality standards met. Ready for Scribe.

[or]

**REQUEST FIXES**:
1. FeatureService.Delete coverage 75% < 80% target → add error path tests
2. Integration test missing for CLI → feature interaction
3. TUI integration not validated → manual test needed
```

**Quality Indicators**:
- ✅ Exact test counts and pass/fail
- ✅ Actual coverage percentages
- ✅ Requirements mapped to test evidence
- ✅ Clear recommendation with rationale

### 7.4 Scribe → Director Communication

**Expected Deliverable Format**:

```markdown
# Scribe Deliverable: [Feature Name] Documentation

## Files Modified Summary

| File | Lines Added | Lines Modified | Purpose |
|------|-------------|----------------|---------|
| README.md | 25 | 3 | Feature highlight in Features section |
| docs/feature-guide.md | 350 | 0 | Complete user guide |
| CHANGELOG.md | 45 | 0 | Release notes prepared |
| docs/v3-implementation-roadmap.md | 0 | 5 | Milestone checkboxes marked |

Total: 420 lines added, 8 modified

## Cross-Reference Validation
✅ README links to docs/feature-guide.md
✅ Guide links to ARCHITECTURE.md
✅ All internal references tested
✅ No broken links

## Examples Verified
✅ Command: `ticketr feature create --name example`
   - Tested: Success
   - Output matches documentation

✅ TUI workflow: Press 'f' → feature modal
   - Tested: Modal opens correctly
   - Keybindings accurate

## Documentation Highlights

### README.md Changes
Added feature to main Features list with brief description.

### docs/feature-guide.md (NEW)
Sections:
- Introduction: What is the feature?
- Getting Started: Step-by-step first use
- Command Reference: All commands with examples
- TUI Integration: Keybindings and workflows
- Troubleshooting: Common issues
- Best Practices

Length: 350 lines, comprehensive coverage.

### CHANGELOG.md
Prepared entry in [Unreleased] section:
- Feature description
- CLI commands
- TUI enhancements
- Technical details (coverage, LOC)

## Quality Checks
✅ Markdown renders correctly (verified locally)
✅ Code blocks use language hints (```bash, ```go)
✅ Tables properly formatted
✅ Consistent voice (imperative, present tense)
✅ No typos (spell-checked)

## Roadmap Update
✅ Milestone X marked complete with all checkboxes
```

**Quality Indicators**:
- ✅ File-by-file breakdown with line counts
- ✅ Examples actually tested
- ✅ Cross-references validated
- ✅ Quality checks performed

### 7.5 Steward → Director Communication

**Expected Deliverable Format**:

```markdown
# Steward Architectural Review: Phase X / Milestone Y

## Review Scope
- [Files reviewed]
- [Architecture areas assessed]
- [Security considerations]

## Architecture Compliance

### Hexagonal Architecture
✅ Domain models remain pure (no external dependencies)
✅ Service layer uses ports, not concrete adapters
✅ Adapters implement interfaces correctly
✅ CLI/TUI layers thin (presentation only)

### Specific Findings
- [Finding 1]: Compliant
- [Finding 2]: Minor deviation - acceptable because [rationale]
- [Finding 3]: Concern - needs addressing

## Security Assessment

### Credential Management
✅ Credentials stored in OS keychain only
✅ No credentials in database or logs
✅ Proper redaction implemented

### Data Validation
✅ Input validation at service boundary
✅ SQL injection prevention (parameterized queries)
✅ Error messages don't leak sensitive data

### Findings
- [Security finding 1]
- [Security finding 2]

## Requirements Validation

| Requirement ID | Status | Evidence |
|----------------|--------|----------|
| PROD-XXX | ✅ | [How validated] |
| PROD-YYY | ✅ | [How validated] |

## Technical Debt Assessment
- [Debt item 1]: Acceptable for now, track as issue
- [Debt item 2]: Should be addressed before next phase
- Overall: Low/Medium/High

## Risk Analysis
- **Performance Risk**: Low - benchmarks within targets
- **Security Risk**: Low - proper patterns followed
- **Maintainability Risk**: Low - well-documented code
- **Scalability Risk**: Medium - [specific concern and mitigation]

## Recommendation

**APPROVED** - Phase X ready for production deployment.

[or]

**APPROVED WITH CONDITIONS**:
1. Address [issue 1] before release
2. Create follow-up issue for [issue 2]
3. Document [consideration 3] in ARCHITECTURE.md

[or]

**REJECTED** - Critical issues must be resolved:
1. [Blocking issue 1 with details]
2. [Blocking issue 2 with details]

Remediation plan required before re-submission.

## Next Steps
[Recommended actions based on decision]
```

**Quality Indicators**:
- ✅ Comprehensive architecture review
- ✅ Specific evidence cited
- ✅ Clear go/no-go decision
- ✅ Actionable recommendations

### 7.6 Escalation Paths

**Level 1: Agent Rework** (Director handles)

Issue: Builder's code doesn't compile, Verifier finds low coverage, Scribe's docs inaccurate

Action:
1. Identify specific issues
2. Delegate back to agent with clear fix requests
3. Review revised deliverable
4. Iterate until quality met

Timeline: 10-30 minutes per iteration

**Level 2: Human Guidance** (Director escalates)

Issue: Architectural uncertainty, conflicting requirements, missing context

Action:
1. Formulate specific question
2. Provide context and attempts made
3. Request guidance or decision
4. Resume work with clarity

Timeline: Wait for human response (minutes to hours)

**Level 3: Project Blocker** (Critical escalation)

Issue: Missing credentials, external dependency failure, critical bug in production

Action:
1. Document blocker completely
2. Assess impact on timeline
3. Escalate urgently to human
4. Pause affected work, pivot to unblocked tasks if possible

Timeline: Immediate escalation, resolution time varies

---

## APPENDICES

### Appendix A: TodoList Management Best Practices

**Task Granularity Guidelines**:

✅ **Good Tasks** (atomic, clear ownership):
```python
{"content": "Delegate BulkOperation domain model to Builder", "activeForm": "Delegating domain model", "status": "pending"}
{"content": "Review Builder deliverable and approve or request fixes", "activeForm": "Reviewing Builder deliverable", "status": "pending"}
{"content": "Run full test suite and verify 0 failures", "activeForm": "Running test suite", "status": "pending"}
```

❌ **Bad Tasks** (too vague, not actionable):
```python
{"content": "Work on milestone", "activeForm": "Working", "status": "pending"}
{"content": "Make sure everything works", "activeForm": "Checking", "status": "pending"}
{"content": "Do Phase 5", "activeForm": "Doing Phase 5", "status": "pending"}
```

**Status Discipline**:

```python
# CORRECT: Exactly ONE task in_progress
TodoWrite(todos=[
    {"content": "Task 1", "activeForm": "...", "status": "completed"},
    {"content": "Task 2", "activeForm": "...", "status": "completed"},
    {"content": "Task 3", "activeForm": "...", "status": "in_progress"},  # ONE
    {"content": "Task 4", "activeForm": "...", "status": "pending"},
    {"content": "Task 5", "activeForm": "...", "status": "pending"}
])

# WRONG: Multiple in_progress
TodoWrite(todos=[
    {"content": "Task 1", "activeForm": "...", "status": "in_progress"},  # BAD
    {"content": "Task 2", "activeForm": "...", "status": "in_progress"},  # BAD
])

# WRONG: Nothing in_progress (if work ongoing)
TodoWrite(todos=[
    {"content": "Task 1", "activeForm": "...", "status": "completed"},
    {"content": "Task 2", "activeForm": "...", "status": "pending"},  # Should be in_progress
])
```

**Update Frequency**:
- After completing ANY task
- Before starting new delegation
- When changing context (e.g., pausing for human input)

### Appendix B: Git Commit Message Templates

**Feature Implementation**:

```bash
feat(workspace): Add credential profile reuse functionality

Enable multiple workspaces to share a single credential profile,
reducing redundant credential entry and improving team workflows.

Implementation:
- Added CredentialProfile domain model
- Extended WorkspaceService with profile methods
- Implemented SQLite repository for profiles
- Added CLI commands: credentials profile create/list
- Integrated profile selection in TUI workspace modal

Testing:
- 25 new tests added
- Service layer coverage: 95.2%
- Integration tests validate end-to-end workflow
- All 450 tests passing

Technical Details:
- Database migration 003: credential_profiles table
- Foreign key relationship: workspaces.credential_profile_id
- Keychain storage unchanged (profiles reference existing credentials)

Files Modified:
- internal/core/domain/credential_profile.go (+150 lines)
- internal/core/services/workspace_service.go (+200 lines)
- internal/adapters/database/credential_profile_repository.go (+280 lines)
- cmd/ticketr/credentials_commands.go (+320 lines)
- internal/adapters/tui/views/workspace_modal.go (+450 lines)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
```

**Bug Fix**:

```bash
fix(workspace): Persist workspace switching across command invocations

Fixed critical bug where workspace switching didn't persist between
CLI commands. Users had to re-switch on every command invocation.

Root Cause:
- Current() relied on in-memory cache lost between commands
- Fallback to GetDefault() instead of checking last_used timestamp

Solution:
- Updated Current() to return most recently used workspace
- Changed logic to use List()[0] (sorted by last_used DESC)
- Added integration test for switch persistence

Testing:
- All 50+ workspace service tests passing
- New test: TestWorkspaceSwitchPersistence validates fix
- Manual verification: switch → exit → current → shows switched

Impact:
- Users can now switch workspaces reliably
- TUI workspace list shows correct active workspace indicator

Files Modified:
- internal/core/services/workspace_service.go (+20, -8 lines)
- internal/core/services/workspace_service_test.go (+100 lines)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
```

**Documentation Update**:

```bash
docs(workspace): Add comprehensive credential profile documentation

Created complete user guide for credential profile functionality,
including CLI and TUI workflows, security model, and troubleshooting.

Documentation Added:
- README.md: Credential profile feature highlight (+25 lines)
- docs/workspace-management-guide.md: Complete section (+400 lines)
- CHANGELOG.md: Release notes for Milestone 18 (+61 lines)
- docs/v3-implementation-roadmap.md: Milestone checkboxes (updated)

Content Highlights:
- Step-by-step CLI workflow
- TUI keyboard shortcuts and modal usage
- Security model explanation (keychain storage)
- Migration guide from direct credentials
- Troubleshooting common issues
- Best practices for team adoption

Quality Assurance:
✅ All examples tested manually
✅ Cross-references validated
✅ Markdown renders correctly
✅ Consistent formatting throughout

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
```

### Appendix C: Quality Checklist Templates

**Builder Quality Checklist**:

```markdown
## Builder Pre-Submission Checklist

### Code Quality
- [ ] All code compiles: `go build ./...`
- [ ] All tests pass: `go test ./...`
- [ ] No `go vet` warnings
- [ ] Code formatted with `gofmt`
- [ ] No TODOs without tracking issues

### Architecture Compliance
- [ ] Domain models have no external dependencies
- [ ] Services use ports (interfaces), not concrete adapters
- [ ] Adapters implement ports correctly
- [ ] No business logic in CLI/TUI layer

### Testing
- [ ] New functions have tests
- [ ] Error paths tested
- [ ] Edge cases validated
- [ ] Coverage target met (>70% for new code)

### Documentation
- [ ] Public functions have godoc comments
- [ ] Complex logic has inline comments
- [ ] README updated if user-facing

### Deliverable Completeness
- [ ] Files modified list with line counts
- [ ] Test results included (actual output)
- [ ] Build verification shown
- [ ] Design decisions documented
- [ ] Notes for Verifier provided
```

**Verifier Quality Checklist**:

```markdown
## Verifier Pre-Submission Checklist

### Test Execution
- [ ] Full suite run: `go test ./... -v`
- [ ] All tests pass (or skips documented)
- [ ] Zero new test failures (no regressions)
- [ ] Test count increased appropriately

### Coverage Validation
- [ ] Overall coverage checked
- [ ] New code coverage meets target (>70%)
- [ ] Critical paths have >80% coverage
- [ ] Coverage report generated

### Regression Detection
- [ ] Compared test count before/after
- [ ] Validated no unrelated test failures
- [ ] Checked for performance regressions (if applicable)

### Requirements Validation
- [ ] Each acceptance criterion mapped to test
- [ ] Evidence provided for each criterion
- [ ] All requirements satisfied

### Recommendation
- [ ] Clear APPROVE or REQUEST FIXES
- [ ] If REQUEST FIXES: specific issues listed
- [ ] Priority indicated for each issue
```

**Scribe Quality Checklist**:

```markdown
## Scribe Pre-Submission Checklist

### Content Completeness
- [ ] README.md updated (if user-facing change)
- [ ] User guide created or updated
- [ ] CHANGELOG.md entry prepared
- [ ] Roadmap checkboxes marked
- [ ] Technical docs updated (if architectural change)

### Accuracy
- [ ] All examples tested manually
- [ ] Command syntax verified
- [ ] TUI keybindings accurate
- [ ] Cross-references validated

### Quality
- [ ] Markdown renders correctly
- [ ] Code blocks have language hints
- [ ] Tables properly formatted
- [ ] Consistent voice (imperative, present tense)
- [ ] No typos (spell-checked)

### Deliverable Completeness
- [ ] Files modified list with line counts
- [ ] Cross-reference validation report
- [ ] Example verification results
- [ ] Documentation diff summary
```

### Appendix D: Decision Trees

**Decision Tree: Agent Selection**:

```
What type of work needs to be done?
├─ Code implementation → Builder
├─ Test validation/extension → Verifier
├─ Documentation creation/update → Scribe
├─ Architectural review → Steward
├─ Research investigation → General-purpose (rare)
└─ Planning/orchestration → Director (you)
```

**Decision Tree: Quality Gate Pass/Fail**:

```
Builder delivered code. Does it pass quality gate?
├─ All tests pass?
│  ├─ NO → Delegate back to Builder with specific failures
│  └─ YES → Continue
├─ Code compiles?
│  ├─ NO → Delegate back to Builder
│  └─ YES → Continue
├─ Architecture compliant?
│  ├─ NO → Delegate back to Builder with patterns to follow
│  └─ YES → Continue
└─ All YES → Proceed to Verifier

Verifier delivered results. Does it pass quality gate?
├─ All tests pass?
│  ├─ NO → Delegate back to Builder (via Director)
│  └─ YES → Continue
├─ Coverage targets met?
│  ├─ NO → Verifier should add tests OR request Builder fixes
│  └─ YES → Continue
├─ Zero regressions?
│  ├─ NO → Delegate back to Builder to fix
│  └─ YES → Continue
├─ Requirements validated?
│  ├─ NO → Investigate gaps, delegate to Builder
│  └─ YES → Continue
└─ All YES → Proceed to Scribe

Scribe delivered documentation. Does it pass quality gate?
├─ All required files updated?
│  ├─ NO → Delegate back to Scribe with missing items
│  └─ YES → Continue
├─ Examples tested and accurate?
│  ├─ NO → Delegate back to Scribe to fix
│  └─ YES → Continue
├─ Cross-references valid?
│  ├─ NO → Delegate back to Scribe to fix links
│  └─ YES → Continue
└─ All YES → Proceed to Git Commit
```

### Appendix E: Phase Gate Approval Template

**Phase Gate Submission to Steward**:

```python
Task(
    subagent_type="general-purpose",
    description="Steward approval for Phase X completion",
    prompt="""You are the Steward agent for Ticketr v3.0 Phase X Gate Approval.

## Mission
Provide comprehensive architectural review and go/no-go decision for Phase X completion.

## Phase X Context

### Phase Goals
[List phase goals from roadmap]

### Deliverables Summary
**Milestones Completed:**
- Milestone A: [Brief description]
- Milestone B: [Brief description]
- Milestone C: [Brief description]

**Code Statistics:**
- Total LOC added: [X]
- Files created: [Y]
- Files modified: [Z]
- Test coverage: [W%]

**Test Evidence:**
[Paste latest full test suite results]

**Documentation:**
- README.md updated
- User guides: [list]
- CHANGELOG.md prepared
- Roadmap marked complete

### Gate Requirements
From docs/phase-X-gate-approval.md:

**P0 Requirements:**
1. [Requirement 1]
2. [Requirement 2]
3. [Requirement 3]

**P1 Requirements:**
1. [Requirement 1]
2. [Requirement 2]

**P2 Requirements:**
1. [Requirement 1]

## Review Tasks

### 1. Architecture Compliance Assessment
Review the following for hexagonal architecture compliance:
- Domain models: [file list]
- Services: [file list]
- Adapters: [file list]
- Ports: [file list]

Validate:
- No domain dependencies on external packages
- Services use ports, not concrete adapters
- Adapters properly implement interfaces
- CLI/TUI layers remain thin

### 2. Security Architecture Review
Assess security implications:
- Credential storage patterns
- Data validation at boundaries
- Error handling (no sensitive leaks)
- Input sanitization

### 3. Requirements Validation
Map each gate requirement to implementation evidence.

### 4. Phase Readiness Assessment
Evaluate:
- Production deployment readiness
- Risk assessment
- Technical debt introduced
- Performance benchmarks met
- Breaking changes (if any)

## Expected Deliverables

### 1. Architecture Compliance Report
For each component, provide compliance assessment.

### 2. Security Assessment
List security considerations and how they're addressed.

### 3. Requirements Compliance Matrix
| Requirement ID | Priority | Status | Evidence |
|----------------|----------|--------|----------|
| [ID] | P0 | ✅/❌ | [How validated] |

### 4. Technical Debt Evaluation
- Debt introduced (if any)
- Mitigation plan
- Impact assessment

### 5. Risk Analysis
- Performance risks
- Security risks
- Maintainability risks
- Mitigation strategies

### 6. Final Recommendation

**Option 1: APPROVED**
Phase X ready for production deployment.

**Option 2: APPROVED WITH CONDITIONS**
Approved subject to addressing:
1. [Condition 1]
2. [Condition 2]

**Option 3: REJECTED**
Critical issues must be resolved:
1. [Blocking issue 1]
2. [Blocking issue 2]

Provide detailed justification for decision.

### 7. Next Steps Recommendation
Based on decision, recommend next actions.

Begin comprehensive review now.
"""
)
```

### Appendix F: Troubleshooting Guide

**Common Issue: Builder Returns Failing Tests**

Symptoms:
- Builder reports test failures
- Implementation appears complete

Diagnosis:
1. Read exact error messages
2. Identify failure type:
   - Logic bug in implementation
   - Test environment issue
   - Existing regression
   - Test itself incorrect

Resolution:
```python
Task(
    subagent_type="general-purpose",
    description="Fix test failures",
    prompt="""Builder agent: Fix the following test failures:

[Paste exact error output]

Analysis:
[Your diagnosis of root cause]

Required Fix:
[Specific fix needed]

Verify fix with: `go test ./... -v`
"""
)
```

**Common Issue: Verifier Finds Regressions**

Symptoms:
- New code works
- Existing tests now fail

Diagnosis:
1. Identify which tests regressed
2. Determine if true regression or legitimate change
3. Assess scope of impact

Resolution:
- If true regression: Delegate to Builder to fix implementation
- If legitimate change: Delegate to Builder to update tests
- If test environment issue: Fix environment, re-run Verifier

**Common Issue: Steward Rejects Phase**

Symptoms:
- Steward returns REJECTED or APPROVED WITH CONDITIONS
- Must address issues before proceeding

Diagnosis:
1. Read full Steward report
2. Categorize issues by severity
3. Determine which agent addresses each

Resolution:
1. Create remediation TodoList
2. Delegate fixes to appropriate agents
3. Re-verify all work
4. Update documentation
5. Re-submit to Steward with evidence

**Common Issue: Agent Doesn't Understand Task**

Symptoms:
- Agent returns confused response
- Deliverables don't match request

Diagnosis:
- Prompt too vague
- Missing context
- Task too complex

Resolution:
1. Improve prompt specificity
2. Provide more context (file locations, patterns)
3. Break task into smaller pieces
4. Include examples of expected output

---

## Conclusion

This Director's Orchestration Framework provides a comprehensive, battle-tested methodology for executing Ticketr v3.0 Phase 5 and beyond. It synthesizes learnings from:

- **Phase 2**: 5,791 LOC delivered with zero regressions
- **Phase 4**: 4,400 LOC TUI implementation
- **Milestone 18**: Credential profiles with 95%+ service coverage

**Key Success Factors**:

1. **Disciplined Agent Delegation**: Right agent for the right task
2. **Quality Gate Enforcement**: Never compromise on standards
3. **Sequential Execution**: Complete current task before next
4. **Clear Communication**: Structured formats ensure clarity
5. **Continuous Improvement**: Learn from each milestone

**For New Directors**:

1. Read this framework sequentially first time
2. Reference specific sections during execution
3. Follow templates and checklists precisely
4. Adapt thoughtfully with justification
5. Document learnings for future improvements

**Framework Maintenance**:

This document is a **living framework**. Update it when:
- New patterns proven effective
- Anti-patterns identified
- Process improvements discovered
- Technology changes (Go version, libraries)

**Next Actions**:

1. Familiarize yourself with Phase 5 roadmap
2. Execute Week 18 (Bulk Operations) following this methodology
3. Document any framework improvements discovered
4. Iterate and improve

**The framework works. Trust the process. Deliver with confidence.**

---

**Document Version:** 2.0
**Status:** Production Framework
**Maintenance:** Update with each phase completion
**Contact:** Director agents for methodology questions

**Related Documents**:
- `docs/v3-implementation-roadmap.md` - Phase/milestone details
- `docs/DIRECTOR-HANDBOOK.md` - Director role guide
- `.agents/director.agent.md` - Director agent prompt
- `REQUIREMENTS-v2.md` - Product requirements

---

*🤖 Generated with [Claude Code](https://claude.com/claude-code) via [Happy](https://happy.engineering)*

*Co-Authored-By: Claude <noreply@anthropic.com>*
*Co-Authored-By: Happy <yesreply@happy.engineering>*
