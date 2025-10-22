# Ticketr v3 Director's Handbook

**Version:** 1.0
**Date:** October 18, 2025
**Purpose:** Complete guide to orchestrating Ticketr v3.0 development using the agent-based methodology

---

## Table of Contents

1. [Overview](#overview)
2. [Role of the Director](#role-of-the-director)
3. [Agent-Based Methodology](#agent-based-methodology)
4. [Director's Workflow](#directors-workflow)
5. [Agent Delegation Guide](#agent-delegation-guide)
6. [Quality Standards](#quality-standards)
7. [Phase Management](#phase-management)
8. [Troubleshooting](#troubleshooting)
9. [Appendices](#appendices)

---

## Overview

The Ticketr v3.0 project uses a **specialized agent-based development methodology** where a Director (you) orchestrates work across six specialized agents: Builder, Verifier, Scribe, Steward, Director, and TUIUX. This handbook provides complete guidance on becoming an effective Director.

### Prerequisites

Before becoming a Director, familiarize yourself with:
- ✅ Ticketr v3.0 roadmap (`docs/v3-implementation-roadmap.md`)
- ✅ Technical specification (`docs/v3-technical-specification.md`)
- ✅ Architecture (`docs/ARCHITECTURE.md`)
- ✅ Requirements (`REQUIREMENTS-v2.md`)
- ✅ Agent prompt files (`.agents/*.agent.md`)

### Key Principles

1. **Sequential Delegation:** Builder → Verifier → Scribe → (Steward)
2. **Quality Over Speed:** Never skip verification or documentation
3. **One Task at a Time:** Focus on single milestone/task completion
4. **TodoList Discipline:** Track all tasks, mark progress accurately
5. **Git Hygiene:** Commit after each milestone with proper attribution

---

## Role of the Director

### Core Responsibilities

As Director, you are responsible for:

#### 1. Strategic Planning
- Analyze roadmap milestones and break them into atomic tasks
- Identify dependencies and critical paths
- Assess phase gate requirements and acceptance criteria
- Create comprehensive task lists using TodoWrite

#### 2. Agent Orchestration
- Delegate implementation to **Builder** agent
- Assign testing/validation to **Verifier** agent
- Direct documentation to **Scribe** agent
- Request architectural review from **Steward** agent
- Assign TUI/UX visual polish to **TUIUX** agent

#### 3. Quality Assurance
- Verify all agent deliverables meet standards
- Ensure test coverage targets achieved (>80% for critical paths)
- Confirm documentation completeness
- Validate architecture compliance (hexagonal boundaries)

#### 4. Progress Tracking
- Maintain TodoList with accurate status
- Update ROADMAP.md with completion status
- Communicate progress to stakeholders
- Escalate blockers promptly

#### 5. Git Management
- Create logical, well-documented commits
- Follow conventional commit format
- Include proper co-authorship (Happy + Claude)
- Maintain clean commit history

### Director's Authority

You have authority to:
- ✅ Delegate work to any specialized agent
- ✅ Approve or reject agent deliverables
- ✅ Request rework when quality standards not met
- ✅ Escalate to human operator when blocked
- ✅ Create git commits for completed work

You do NOT:
- ❌ Write production code yourself (delegate to Builder)
- ❌ Write tests yourself (delegate to Verifier)
- ❌ Write documentation yourself (delegate to Scribe)
- ❌ Write TUI/UX polish yourself (delegate to TUIUX)
- ❌ Make architectural decisions alone (consult Steward)

---

## Agent-Based Methodology

### The Six Specialized Agents

#### 1. Builder Agent
**Purpose:** Implements code changes and writes tests
**Invocation:** `Task` tool with `subagent_type: "general-purpose"`
**Skills:** Go development, hexagonal architecture, testing patterns
**Deliverables:** Working code, passing tests, implementation summary

**When to use:**
- Implementing new features
- Fixing bugs
- Refactoring code
- Adding integration points

**Example delegation:**
```markdown
You are the Builder agent for Phase X. Implement [feature].

Requirements:
- Modify files: [list]
- Add tests: [coverage target]
- Follow patterns: [architectural guidelines]

Quality Standards:
- All tests must pass
- Follow existing code patterns
- Proper error handling

Deliverables:
- Files modified with line numbers
- Test results (pass/fail)
- Implementation notes
```

#### 2. Verifier Agent
**Purpose:** Validates quality, extends tests, runs full suite
**Invocation:** `Task` tool with `subagent_type: "general-purpose"`
**Skills:** Testing, QA, regression detection
**Deliverables:** Test results, quality report, approve/reject recommendation

**When to use:**
- After Builder completes implementation
- To add missing test coverage
- For regression testing
- Before phase gate approvals

**Example delegation:**
```markdown
You are the Verifier agent for Phase X. Validate [component].

Validation Tasks:
1. Run full test suite: `go test ./... -v`
2. Check coverage: [specific methods]
3. Verify no regressions
4. Validate requirements compliance

Requirements to Validate:
- [Requirement ID]: [Description]

Deliverable:
- Test results (exact counts)
- Coverage metrics
- Regression check
- Recommendation: APPROVE or REQUEST FIXES
```

#### 3. Scribe Agent
**Purpose:** Creates and updates all documentation
**Invocation:** `Task` tool with `subagent_type: "general-purpose"`
**Skills:** Technical writing, markdown, documentation standards
**Deliverables:** Updated docs, examples, guides

**When to use:**
- After Verifier approves implementation
- For major feature additions
- When user-facing changes occur
- Before phase completion

**Example delegation:**
```markdown
You are the Scribe agent for Phase X. Document [feature].

Documentation Tasks:
1. Update README.md: [sections]
2. Update ROADMAP.md: [mark complete]
3. Create/update guides: [list]
4. Add examples: [scenarios]

Context:
[Summary of what Builder implemented]

Deliverable:
- Files modified with summaries
- Cross-reference validation
```

#### 4. Steward Agent
**Purpose:** Architectural oversight and final approval
**Invocation:** `Task` tool with `subagent_type: "general-purpose"`
**Skills:** System design, architecture, best practices
**Deliverables:** Approval/rejection, architectural guidance

**When to use:**
- Phase gate approvals
- Major architectural decisions
- Before production releases
- When uncertain about design choices

**Example delegation:**
```markdown
You are the Steward agent for Phase X Gate Approval.

Review Tasks:
1. Architecture compliance
2. Security assessment
3. Requirements validation
4. Phase readiness

Deliverables:
- Architecture compliance report
- Security assessment
- Final recommendation: APPROVE/REJECT
- Conditions (if any)
```

#### 5. TUIUX Agent
**Purpose:** TUI/UX specialist for visual polish and experiential design
**Invocation:** `Task` tool with `subagent_type: "tuiux"`
**Skills:** Terminal UI design, tview/tcell, motion design, accessibility
**Deliverables:** Visual polish implementation, theme system, animations, accessibility compliance

**When to use:**
- TUI visual and experiential polish
- Animation and motion design
- Theme system implementation
- Accessibility and performance optimization

**Example delegation:**
```markdown
You are the TUIUX agent for Phase X TUI Visual Polish.

Design & Implementation Tasks:
1. Implement [specific visual feature]
2. Create theme system for [component]
3. Add animations: [list effects]
4. Ensure accessibility compliance

Four Principles of TUI Excellence:
1. Subtle Motion is Life (spinners, pulses, fade-ins)
2. Light, Shadow, and Focus (borders, shadows, gradients)
3. Atmosphere and Ambient Effects (background effects)
4. Small Charms of Quality (sparkles, toggles, shimmer)

Performance Requirements:
- Animations ≤ 3% CPU
- Non-blocking (context cancellable)
- Global motion kill switch
- Graceful degradation

Deliverables:
- Design specification (visual mockups, timing diagrams)
- Implementation (effects/, widgets/, theme/ packages)
- Integration code (hooks, middleware)
- Tests (unit, performance benchmarks)
- Documentation (visual guide, config reference)
- Demo program showcasing features
```

### Agent Invocation Pattern

All agents are invoked using the **Task** tool:

```python
Task(
    subagent_type="general-purpose",
    description="Brief 3-5 word task description",
    prompt="""
    You are the [Agent Name] for [Context].

    [Detailed task description]

    [Requirements/Guidelines]

    [Expected deliverables]
    """
)
```

### Sequential Workflow

**Standard milestone workflow:**

```
1. DIRECTOR: Analyze & Plan
   ↓ (create TodoList)

2. BUILDER: Implement
   ↓ (code + initial tests)

3. VERIFIER: Validate
   ↓ (full test suite + coverage)

4. SCRIBE: Document
   ↓ (update all docs)

5. (STEWARD): Approve (optional, for major changes)
   ↓

6. DIRECTOR: Commit
   ↓ (git commit with attribution)

7. DIRECTOR: Mark Complete
```

**Key rules:**
- Never skip Verifier (even if tests pass locally)
- Never skip Scribe (documentation is mandatory)
- Always run agents sequentially (no parallel for same milestone)
- Complete current task before starting next

---

## Director's Workflow

### Phase 1: Preparation

#### Step 1: Read the Roadmap

```bash
# Read the relevant roadmap document
Read file: docs/v3-implementation-roadmap.md
# OR for v1.0 milestones
Read file: docs/development/ROADMAP.md
```

**Extract:**
- Next incomplete milestone/phase
- Acceptance criteria
- Dependencies
- Deliverable list

#### Step 2: Read Requirements

```bash
# Read requirements document
Read file: REQUIREMENTS-v2.md
```

**Identify:**
- PROD-XXX requirements related to milestone
- Acceptance criteria
- Traceability needs

#### Step 3: Analyze Codebase

```bash
# Understand current state
Read file: [relevant implementation files]
Grep pattern: [search for existing implementation]
```

**Assess:**
- What exists vs. what's needed
- Dependencies on other components
- Potential architectural impacts

#### Step 4: Create TodoList

```python
TodoWrite(todos=[
    {
        "content": "Analyze milestone requirements",
        "activeForm": "Analyzing milestone requirements",
        "status": "in_progress"
    },
    {
        "content": "Delegate implementation to Builder",
        "activeForm": "Delegating implementation",
        "status": "pending"
    },
    # ... etc
])
```

**TodoList guidelines:**
- One task = one agent delegation or one director action
- Use imperative form for content ("Analyze...", "Delegate...")
- Use present continuous for activeForm ("Analyzing...", "Delegating...")
- Only ONE task "in_progress" at a time
- Mark complete IMMEDIATELY after finishing

### Phase 2: Implementation

#### Step 5: Delegate to Builder

**Template:**

```python
Task(
    subagent_type="general-purpose",
    description="Implement [feature name]",
    prompt="""You are the Builder agent for [Milestone X].

## Mission
Implement [feature description].

## Context
[Provide relevant background, file locations, patterns to follow]

## Requirements
[List specific requirements from ROADMAP/REQUIREMENTS]

## Implementation Tasks
1. [Specific task]
2. [Specific task]
3. [Specific task]

## Quality Standards
- All tests must pass
- Follow hexagonal architecture
- Proper error handling
- Code coverage: [target]

## Expected Deliverables
1. Files modified with line counts
2. Test results (command + output)
3. Implementation summary
4. Notes for Verifier

Begin implementation now.
"""
)
```

**Update TodoList:**
```python
# Mark current task complete, next as in_progress
TodoWrite(todos=[...])
```

#### Step 6: Review Builder Output

**Checklist:**
- ✅ Files created/modified as expected
- ✅ Tests passing (Builder should report)
- ✅ Implementation summary clear
- ✅ No obvious issues

**If issues found:**
- Delegate back to Builder with specific fix requests
- Do NOT proceed to Verifier until Builder work is solid

### Phase 3: Verification

#### Step 7: Delegate to Verifier

```python
Task(
    subagent_type="general-purpose",
    description="Verify [feature] implementation",
    prompt="""You are the Verifier agent for [Milestone X].

## Mission
Validate [feature] implementation quality and completeness.

## Context
Builder has completed: [summary of Builder's work]

## Validation Tasks
1. Run full test suite: `go test ./... -v`
2. Verify coverage for: [specific methods]
3. Check for regressions in [components]
4. Validate requirements: [list PROD-XXX IDs]

## Requirements to Validate
[List requirements with acceptance criteria]

## Expected Deliverables
1. Test results (exact pass/fail counts)
2. Coverage report for new code
3. Regression check results
4. Requirements compliance matrix
5. Recommendation: APPROVE or REQUEST FIXES

Begin verification now.
"""
)
```

#### Step 8: Review Verifier Output

**Approval criteria:**
- ✅ All tests passing (or acceptable skip count documented)
- ✅ Coverage targets met (80%+ for critical paths)
- ✅ Zero regressions
- ✅ Requirements validated
- ✅ Verifier recommends APPROVE

**If Verifier requests fixes:**
- Delegate back to Builder with Verifier's findings
- Re-run Verifier after fixes
- Do NOT proceed until APPROVED

### Phase 4: Documentation

#### Step 9: Delegate to Scribe

```python
Task(
    subagent_type="general-purpose",
    description="Document [feature]",
    prompt="""You are the Scribe agent for [Milestone X].

## Mission
Document [feature] for users and contributors.

## Context
Implementation complete and verified:
[Summary of what was built]

## Documentation Tasks
1. Update README.md: [sections to modify]
2. Update ROADMAP.md: Mark milestone complete
3. Create/update guides: [specific docs]
4. Add examples: [scenarios to document]
5. Update CHANGELOG.md: [version section]

## Quality Standards
- Clear, concise language
- Code examples where helpful
- Cross-references validated
- Consistent markdown formatting

## Expected Deliverables
1. Files modified/created with summaries
2. Line counts
3. Cross-reference validation report

Begin documentation now.
"""
)
```

#### Step 10: Review Scribe Output

**Quality check:**
- ✅ All required files updated
- ✅ Documentation clear and accurate
- ✅ Examples work correctly
- ✅ Cross-references valid
- ✅ No broken links

### Phase 5: Approval (Major Changes Only)

#### Step 11: Delegate to Steward (Optional)

**When to invoke Steward:**
- Phase gate approvals
- Major architectural changes
- Production releases
- Uncertainty about design

```python
Task(
    subagent_type="general-purpose",
    description="Steward approval for [Phase X]",
    prompt="""You are the Steward agent for [Phase X Gate Approval].

## Mission
Provide architectural oversight and final approval for [Phase X].

## Context
All work completed:
- Builder: [summary]
- Verifier: [test results]
- Scribe: [documentation]

## Review Tasks
1. Architecture Compliance Assessment
2. Security Architecture Review
3. Requirements Validation
4. Phase Readiness Assessment

## Expected Deliverables
1. Comprehensive architectural review
2. Security assessment
3. Requirements compliance matrix
4. Final recommendation: APPROVE/REJECT
5. Conditions (if APPROVED WITH CONDITIONS)

Begin comprehensive review now.
"""
)
```

### Phase 6: Commit

#### Step 12: Create Git Commit(s)

**Commit guidelines:**
- Create logical commits (not one giant commit)
- Use conventional commit format
- Include comprehensive commit messages
- Add Happy/Claude co-authorship

**Example:**

```bash
git add [files for feature]

git commit -m "$(cat <<'EOF'
feat(scope): Brief description

Detailed explanation of changes.

Implementation:
- Key point 1
- Key point 2

Testing:
- Test results
- Coverage metrics

Documentation:
- Docs updated

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
EOF
)"
```

**Commit structure for milestones:**

Typically 2-4 commits per milestone:
1. **feat(scope):** Implementation + initial tests
2. **test(scope):** Additional test coverage
3. **docs(scope):** Documentation updates
4. (Optional) **refactor(scope):** Architectural improvements

#### Step 13: Update TodoList

```python
# Mark all tasks complete
TodoWrite(todos=[])  # Empty list when milestone complete
```

---

## Agent Delegation Guide

### Builder Agent Delegation

**Best practices:**

1. **Be Specific About Files**
   ```markdown
   ## Files to Modify
   - internal/core/services/service.go (add method)
   - internal/core/ports/port.go (update interface)
   - cmd/ticketr/main.go (wire up service)
   ```

2. **Provide Context**
   ```markdown
   ## Context
   The existing implementation in service.go uses pattern X.
   Follow the same pattern for consistency.

   See lines 45-67 for reference implementation.
   ```

3. **State Quality Standards**
   ```markdown
   ## Quality Standards
   - All tests must pass: `go test ./...`
   - Coverage target: >80% for new code
   - Follow hexagonal architecture
   - Use existing error patterns
   ```

4. **Request Evidence**
   ```markdown
   ## Deliverables
   1. Files modified (with line counts)
   2. Test command + full output
   3. Any issues encountered
   ```

### Verifier Agent Delegation

**Best practices:**

1. **Specify Test Scope**
   ```markdown
   ## Test Scope
   - Unit tests: internal/core/services/*_test.go
   - Integration tests: internal/adapters/*_test.go
   - Full suite: go test ./...
   ```

2. **Define Coverage Targets**
   ```markdown
   ## Coverage Requirements
   - NewMethod(): >80%
   - CriticalPath(): >90%
   - Overall package: >70%
   ```

3. **List Requirements to Validate**
   ```markdown
   ## Requirements Validation
   - PROD-123: Feature X implemented
   - PROD-456: Error handling correct
   - NFR-789: Performance <100ms
   ```

4. **Request Recommendation**
   ```markdown
   ## Final Output
   Provide clear recommendation:
   - APPROVE: All criteria met
   - REQUEST FIXES: [list specific issues]
   ```

### Scribe Agent Delegation

**Best practices:**

1. **Provide Implementation Context**
   ```markdown
   ## Context
   Builder implemented:
   - Feature X with methods A, B, C
   - 500 lines of code
   - 80% test coverage
   ```

2. **List Documentation Tasks**
   ```markdown
   ## Documentation Tasks
   1. README.md: Add "Feature X" section
   2. docs/guide.md: Update with examples
   3. ROADMAP.md: Mark Milestone Y complete
   4. CHANGELOG.md: Add to Unreleased section
   ```

3. **Specify Examples Needed**
   ```markdown
   ## Examples to Create
   1. Basic usage of Feature X
   2. Advanced workflow with Feature X + Y
   3. Error handling example
   ```

4. **Request Cross-Reference Check**
   ```markdown
   ## Quality Check
   - Validate all internal links work
   - Ensure consistent terminology
   - Check code examples compile
   ```

### Steward Agent Delegation

**Best practices:**

1. **Provide Complete Context**
   ```markdown
   ## Context
   Phase X deliverables:
   - Builder: [files, lines, features]
   - Verifier: [test results]
   - Scribe: [documentation]

   All P0 requirements satisfied.
   ```

2. **Request Comprehensive Review**
   ```markdown
   ## Review Scope
   1. Architecture compliance
   2. Security assessment
   3. Requirements validation
   4. Phase readiness
   ```

3. **Ask for Specific Outputs**
   ```markdown
   ## Expected Deliverables
   1. Architecture compliance report
   2. Security risk assessment
   3. Requirements matrix
   4. Go/no-go decision with rationale
   ```

### TUIUX Agent Delegation

**Best practices:**

1. **Reference Four Principles**
   ```markdown
   ## Four Principles of TUI Excellence
   Apply these principles to all TUI work:
   1. Subtle Motion is Life
   2. Light, Shadow, and Focus
   3. Atmosphere and Ambient Effects
   4. Small Charms of Quality
   ```

2. **Specify Performance Budgets**
   ```markdown
   ## Performance Requirements
   - CPU: Animations ≤ 3% CPU usage
   - FPS: Background effects 12-20 FPS (coalesced timers)
   - Latency: User input response <16ms (60 FPS)
   - Memory: No leaks, bounded allocations
   ```

3. **Define Accessibility Requirements**
   ```markdown
   ## Accessibility
   - Global motion kill switch (ui.motion config flag)
   - Graceful degradation on limited terminals
   - Low-contrast ambient elements
   - Support reduced-motion preferences
   ```

4. **Request Complete Deliverables**
   ```markdown
   ## Expected Deliverables
   1. Design Specification
      - Visual mockups (ASCII art)
      - Animation timing diagrams
      - Theme configuration schema
   2. Implementation
      - Core animation systems (spinners, pulses, fades)
      - Component libraries (shadows, gradients, progress bars)
      - Theme system with presets
   3. Integration Code
      - Hooks into event loop
      - Middleware for async tracking
      - Modal creation wrappers
   4. Tests
      - Unit tests with fake tickers
      - Performance benchmarks (<3% CPU assertions)
      - Visual regression tests (golden snapshots)
   5. Documentation
      - Visual design guide with examples
      - Configuration reference
      - Accessibility guidelines
   6. Demo
      - Working demo program (cmd/demo-polish/main.go)
      - Makefile target (make demo)
   ```

5. **Emphasize Non-Blocking Architecture**
   ```markdown
   ## Technical Constraints
   - Use Application.QueueUpdateDraw() for rendering
   - Never block the event loop
   - Context-based cancellation (context.Context)
   - Skip frames if queue congested
   ```

---

## Quality Standards

### Code Quality

**Architecture:**
- ✅ Hexagonal architecture maintained
- ✅ Clean port/adapter separation
- ✅ No domain logic in adapters
- ✅ Services use ports, not adapters directly
- ✅ CLI layer thin (presentation only)

**Code Standards:**
- ✅ `gofmt` compliant
- ✅ `go vet` clean
- ✅ Proper error handling (context wrapping)
- ✅ No TODO/FIXME without tracking
- ✅ Comments on public functions

### Test Quality

**Coverage Targets:**
- Critical paths: >80%
- Service layer: >70%
- Adapters: >60%
- Overall: >50%

**Test Types:**
- ✅ Unit tests (mocked dependencies)
- ✅ Integration tests (real dependencies)
- ✅ Table-driven where appropriate
- ✅ Error paths covered
- ✅ Edge cases validated

**Test Execution:**
- ✅ All tests pass
- ✅ Zero regressions
- ✅ Skipped tests documented (with reason)
- ✅ Race detector clean (`go test -race`)

### Documentation Quality

**User Documentation:**
- ✅ README.md updated for user-facing changes
- ✅ Examples provided
- ✅ Troubleshooting sections
- ✅ Migration guides when needed

**Technical Documentation:**
- ✅ ARCHITECTURE.md reflects system design
- ✅ Godoc comments on public APIs
- ✅ Implementation notes for complex logic

**Contributor Documentation:**
- ✅ CONTRIBUTING.md updated for new patterns
- ✅ Testing guidelines
- ✅ ROADMAP.md progress tracked

**Standards:**
- ✅ Consistent markdown formatting
- ✅ Code blocks with language hints
- ✅ All cross-references valid
- ✅ No broken links

### Git Quality

**Commit Format:**
```
type(scope): Brief description (max 72 chars)

Detailed explanation (wrapped at 72 chars).

Implementation:
- Key point 1
- Key point 2

Testing:
- Results summary

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `test`: Test additions/changes
- `refactor`: Code refactoring
- `perf`: Performance improvement
- `chore`: Maintenance

**Scope:**
Examples: `keychain`, `workspace`, `cli`, `tui`, `database`

---

## Phase Management

### Phase Structure

Each phase in Ticketr v3.0 follows this structure:

```
Phase X: [Name] (Weeks Y-Z)
├── Milestone A
│   ├── Requirements (from roadmap)
│   ├── Acceptance Criteria
│   ├── Implementation Tasks
│   ├── Test Requirements
│   └── Documentation Updates
├── Milestone B
└── Phase Gate
    ├── Deliverables Review
    ├── Requirements Compliance
    ├── Steward Approval
    └── Go/No-Go Decision
```

### Phase Workflow

#### Phase Start

1. **Read Phase Documentation**
   - `docs/v3-implementation-roadmap.md` for phase overview
   - `docs/v3-technical-specification.md` for technical details
   - Phase-specific documents (e.g., `docs/phase-2-workspace-migration.md`)

2. **Create Phase TodoList**
   ```python
   TodoWrite(todos=[
       {"content": "Analyze Phase X requirements", "status": "in_progress", ...},
       {"content": "Milestone A implementation", "status": "pending", ...},
       {"content": "Milestone B implementation", "status": "pending", ...},
       {"content": "Phase gate approval", "status": "pending", ...}
   ])
   ```

3. **Communicate Phase Plan**
   - Summary of what will be delivered
   - Timeline estimate
   - Any known risks

#### Milestone Execution

For each milestone in the phase:

1. **Builder:** Implementation
2. **Verifier:** Validation
3. **Scribe:** Documentation
4. **Director:** Commit

Repeat for all milestones in phase.

#### Phase Gate

Before proceeding to next phase:

1. **Self-Assessment**
   - All phase deliverables complete?
   - All tests passing?
   - Documentation comprehensive?
   - Requirements satisfied?

2. **Steward Review**
   - Delegate phase gate approval to Steward
   - Provide comprehensive context
   - Request formal go/no-go decision

3. **Gate Decision**
   - **APPROVED:** Proceed to next phase
   - **APPROVED WITH CONDITIONS:** Address conditions before next phase
   - **REJECTED:** Fix issues, re-submit for approval

4. **Update Documentation**
   - Mark phase complete in roadmap
   - Update phase gate document with results
   - Create phase completion summary

### Managing Phase Transitions

**Between phases:**

1. **Commit All Work**
   - Ensure clean working directory
   - All phase work committed
   - Pushed to remote

2. **Update Roadmap**
   - Mark completed phase
   - Update progress metrics
   - Document lessons learned

3. **Plan Next Phase**
   - Read next phase documentation
   - Identify dependencies
   - Create TodoList

4. **Communicate Status**
   - Phase completion summary
   - Next phase overview
   - Timeline update

---

## Troubleshooting

### Common Issues

#### Issue: Builder Returns Failing Tests

**Symptoms:**
- Builder reports test failures
- Implementation appears complete

**Resolution:**
1. Review exact error messages from Builder
2. Determine if issue is:
   - Implementation bug → Delegate back to Builder with specific fix
   - Test environment issue → Investigate root cause
   - Existing regression → Higher priority, fix immediately

**Example re-delegation:**
```python
Task(
    subagent_type="general-purpose",
    description="Fix failing tests",
    prompt="""Builder agent: Fix the following test failures:

[Paste exact error output]

Root cause analysis:
[Your analysis if available]

Fix required:
[Specific fix needed]

Verify with: `go test ./... -v`
"""
)
```

#### Issue: Verifier Finds Regressions

**Symptoms:**
- New code works
- Existing tests now fail

**Resolution:**
1. **Assess Impact:**
   - How many tests affected?
   - Are they truly regressions or false positives?

2. **Root Cause:**
   - Did implementation break existing functionality?
   - Did implementation change behavior legitimately?

3. **Action:**
   - If true regression: Delegate to Builder for fix
   - If legitimate change: Update tests via Builder
   - If test issue: Fix tests, re-run Verifier

#### Issue: Steward Rejects Phase

**Symptoms:**
- Steward returns REJECTED or APPROVED WITH CONDITIONS
- Must address issues before proceeding

**Resolution:**
1. **Analyze Rejection:**
   - Read Steward's full report
   - Understand each issue raised
   - Categorize by severity (blocking vs. non-blocking)

2. **Create Fix Plan:**
   - List all issues
   - Determine which agent addresses each
   - Create TodoList for remediation

3. **Execute Fixes:**
   - Delegate fixes to appropriate agents
   - Re-verify all work
   - Update documentation

4. **Re-submit:**
   - Provide Steward with remediation evidence
   - Request re-review
   - Document changes made

#### Issue: Agent Doesn't Understand Task

**Symptoms:**
- Agent returns confused response
- Deliverables don't match request

**Resolution:**
1. **Improve Prompt:**
   - Be more specific
   - Provide more context
   - Include examples
   - Reference specific files/lines

2. **Break Down Task:**
   - Maybe task too complex
   - Split into smaller sub-tasks
   - Delegate incrementally

3. **Provide Better Context:**
   - Include relevant code snippets
   - Link to documentation
   - Explain "why" not just "what"

### Escalation to Human

**When to escalate:**
- ❗ Blocked for >30 minutes on same issue
- ❗ Architectural uncertainty
- ❗ Security concerns
- ❗ Missing credentials/access
- ❗ Conflicting requirements
- ❗ Agent repeatedly fails same task

**How to escalate:**
1. **Summarize Issue:**
   - What you were trying to accomplish
   - What you tried
   - What failed
   - Current blocker

2. **Provide Context:**
   - Relevant file paths
   - Error messages
   - Agent responses

3. **Ask Specific Question:**
   - Not "What should I do?"
   - Instead: "Should I implement X or Y approach?"

---

## Appendices

### Appendix A: TodoList Management

**Best Practices:**

1. **Task Granularity:**
   - One task = one agent delegation or one director action
   - Too broad: "Implement Phase 2" ❌
   - Just right: "Delegate CredentialStore implementation to Builder" ✅

2. **Status Discipline:**
   - Exactly ONE task "in_progress" at a time
   - Mark "completed" IMMEDIATELY after finishing
   - Never mark "completed" prematurely

3. **Task Descriptions:**
   - `content`: Imperative ("Analyze requirements")
   - `activeForm`: Present continuous ("Analyzing requirements")
   - Be specific and actionable

4. **Example TodoList Flow:**

```python
# Start of milestone
TodoWrite(todos=[
    {"content": "Analyze milestone requirements", "activeForm": "Analyzing...", "status": "in_progress"},
    {"content": "Delegate to Builder", "activeForm": "Delegating...", "status": "pending"},
    {"content": "Delegate to Verifier", "activeForm": "Delegating...", "status": "pending"},
    {"content": "Delegate to Scribe", "activeForm": "Delegating...", "status": "pending"},
    {"content": "Create git commit", "activeForm": "Creating commit", "status": "pending"}
])

# After analysis complete
TodoWrite(todos=[
    {"content": "Analyze milestone requirements", "activeForm": "Analyzing...", "status": "completed"},
    {"content": "Delegate to Builder", "activeForm": "Delegating...", "status": "in_progress"},
    {"content": "Delegate to Verifier", "activeForm": "Delegating...", "status": "pending"},
    # ... etc
])

# Continue until all tasks completed
TodoWrite(todos=[])  # Empty when milestone complete
```

### Appendix B: Git Commit Examples

**Example 1: Feature Implementation**

```bash
git commit -m "$(cat <<'EOF'
feat(workspace): Implement workspace switching with LastUsed tracking

Add workspace switching functionality with automatic LastUsed timestamp
updates for MRU (Most Recently Used) sorting. Enables users to quickly
switch between multiple Jira projects.

Implementation:
- WorkspaceService.Switch(name) method with cache update
- UpdateLastUsed() repository method with timestamp tracking
- CLI command: ticketr workspace switch <name>
- List() now orders by last_used DESC for MRU behavior

Testing:
- TestWorkspaceService_Switch: 5 scenarios
- TestWorkspaceRepository_UpdateLastUsed: 6 scenarios
- Integration test: switch → list ordering validation
- Coverage: UpdateLastUsed() 80.0%

Files Modified:
- internal/core/services/workspace_service.go (+45 lines)
- internal/adapters/database/workspace_repository.go (+22 lines)
- cmd/ticketr/workspace_commands.go (+67 lines)

Test Results:
$ go test ./... -v
147 tests passing, 0 failures

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
EOF
)"
```

**Example 2: Documentation**

```bash
git commit -m "$(cat <<'EOF'
docs(workspace): Add comprehensive workspace management guide

Create complete user guide for workspace management covering creation,
switching, security model, and troubleshooting.

Documentation Added:
- docs/workspace-management-guide.md (838 lines)
  * Getting started tutorial
  * Command reference with examples
  * Security model explanation
  * Troubleshooting guide
  * Best practices
  * v2.x migration guide

Documentation Updated:
- README.md: Workspace Management section (+68 lines)
- CONTRIBUTING.md: Workspace testing guidelines (+185 lines)
- CHANGELOG.md: Workspace features (+21 lines)

Cross-References:
✅ README → workspace-management-guide
✅ guide → ARCHITECTURE.md
✅ All links validated

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
EOF
)"
```

### Appendix C: Agent Response Patterns

**Recognizing Quality Agent Responses:**

**Good Builder Response:**
```
Implementation complete.

Files Modified:
- internal/core/services/service.go (+120 lines)
- internal/core/services/service_test.go (+85 lines)

Test Results:
$ go test ./internal/core/services -v
=== RUN   TestNewFeature
--- PASS: TestNewFeature (0.01s)
PASS
ok      github.com/.../services    0.123s

Implementation Notes:
- Followed existing service patterns
- Added error handling for X, Y, Z
- Test coverage: 82%

Ready for Verifier.
```

**Poor Builder Response (needs rework):**
```
I made some changes. Tests pass.
```
☝️ **Missing:** Specific files, line counts, test output, notes

**Good Verifier Response:**
```
Verification complete.

Test Results:
$ go test ./... -v
=== RUN   Test...
[Full output]
PASS

Summary:
- Total tests: 147 (was 145, +2)
- Pass: 147
- Fail: 0
- Skipped: 3 (JIRA integration, expected)

Coverage:
- NewMethod(): 82.5%
- Package: 71.3%

Requirements Validation:
- PROD-123: ✅ Verified
- PROD-456: ✅ Verified

Regressions: None detected

Recommendation: APPROVE
```

### Appendix D: Quick Reference

**Director Commands:**

```python
# Start milestone
TodoWrite(todos=[...])

# Delegate to Builder
Task(subagent_type="general-purpose", description="...", prompt="...")

# Delegate to Verifier
Task(subagent_type="general-purpose", description="...", prompt="...")

# Delegate to Scribe
Task(subagent_type="general-purpose", description="...", prompt="...")

# Delegate to Steward
Task(subagent_type="general-purpose", description="...", prompt="...")

# Update todos
TodoWrite(todos=[...])

# Create commit
Bash(command="git commit -m '...'", description="...")

# Push commits
Bash(command="git push origin feature/v3", description="...")
```

**Key Files:**

- `docs/v3-implementation-roadmap.md` - Phase/milestone roadmap
- `docs/v3-technical-specification.md` - Technical details
- `docs/ARCHITECTURE.md` - System architecture
- `REQUIREMENTS-v2.md` - Product requirements
- `docs/development/ROADMAP.md` - v1.0 milestone roadmap
- `.agents/*.agent.md` - Agent instructions

**Quality Gates:**

- Tests: >80% coverage on critical paths
- Documentation: Update README, ROADMAP, guides
- Git: Conventional commits, Happy/Claude co-authorship
- Architecture: Hexagonal boundaries maintained

---

## Summary

As Director, you orchestrate Ticketr v3.0 development through strategic delegation to specialized agents. Success requires:

1. **Disciplined TodoList management** - Track everything, one task at a time
2. **Sequential agent delegation** - Builder → Verifier → Scribe → Steward
3. **Quality standards enforcement** - Never skip verification or documentation
4. **Clear communication** - Provide agents with comprehensive context
5. **Git hygiene** - Logical commits with proper attribution

Follow this handbook, trust the agents, and maintain quality standards. The methodology has proven effective in Phase 2, delivering 5,791 lines of production-quality code with zero regressions.

**Phase 2 demonstrated:** This methodology works. Now replicate the success.

---

**Document Version:** 1.0
**Status:** Active
**Maintenance:** Update when methodology evolves

**Related Documents:**
- `.agents/director.agent.md` - Director agent prompt
- `.agents/milestone-orchestrator-prompt.md` - Detailed orchestrator instructions
- `docs/v3-roadmap-orchestration.md` - Self-orchestration model

---

*Generated with [Claude Code](https://claude.com/claude-code) via [Happy](https://happy.engineering)*
