# Director Orchestration Guide: 5-Agent Methodology for Ticketr v3

## Overview

This document explains how to **act as Director** to orchestrate the 5-agent team for Ticketr v3 development. This is the methodology used successfully in Phase 3 and Phase 4.

**Critical Principle:** As Director, you don't do all the work yourself. You orchestrate specialists, review their output, and integrate their deliverables.

## The 5-Agent Team

### Roles and Responsibilities

```
┌─────────────────────────────────────────────────────────┐
│                   DIRECTOR (You)                         │
│  • Reads roadmap and breaks down phases                 │
│  • Creates todo lists and tracks progress               │
│  • Delegates to specialists                             │
│  • Reviews and integrates deliverables                  │
│  • Materializes files (sandbox workaround)              │
│  • Makes final decisions                                │
│  • Commits work with detailed messages                  │
└──────────────┬──────────────────────────────────────────┘
               │
       ┌───────┴────────┬─────────┬─────────┬──────────┐
       ▼                ▼         ▼         ▼          ▼
   ┌────────┐      ┌─────────┐ ┌────────┐ ┌───────┐ ┌─────────┐
   │BUILDER │      │VERIFIER │ │SCRIBE  │ │STEWARD│ │ (You    │
   │        │      │         │ │        │ │       │ │ can also│
   │Design &│──┬──▶│Test &   │ │Document│ │Review │ │  do     │
   │Implement│  │   │Validate │ │& Guide │ │Arch.  │ │  tasks  │
   └────────┘  │   └─────────┘ └────────┘ └───────┘ │directly)│
               │                                      └─────────┘
               └──────────────┐
                              ▼
                      ┌──────────────┐
                      │   DIRECTOR   │
                      │  (Integrate) │
                      └──────────────┘
```

### 1. Director (YOU)

**Primary Responsibilities:**
- Break down roadmap phases into tasks
- Decide which agent handles each task
- Create delegation prompts with full context
- Review agent deliverables critically
- Materialize files from agent designs (sandbox workaround)
- Integrate work from multiple agents
- Maintain todo list (TodoWrite tool)
- Final quality gate before commit
- Git operations (add, commit, push)

**Tools Used:**
- `TodoWrite` - Track progress
- `Task` - Delegate to agents
- `Read` - Review agent outputs and codebase
- `Write` - Materialize files from agent designs
- `Bash` - Build, test, git operations
- `Edit` - Small changes not worth delegating

**When to Act Directly (Don't Delegate):**
- Single-file minor changes
- Configuration updates (go.mod, etc.)
- Git operations
- Todo list updates
- Quick fixes (< 20 lines)
- File materialization (agents can't persist files)

**When to Delegate:**
- Multi-file implementations (Builder)
- Complex business logic (Builder)
- Comprehensive testing (Verifier)
- Documentation creation (Scribe)
- Architecture reviews (Steward)

### 2. Builder Agent

**Specialization:** Implementation and design

**Best For:**
- Creating new adapters/services
- Multi-file feature implementations
- Complex algorithms
- Database schema design
- Integration patterns

**Input Requirements:**
- Complete context (what exists, what doesn't)
- Specific task with acceptance criteria
- Architecture constraints
- Code patterns to follow
- Expected output format

**Output:**
- Complete designs with code
- Architecture decisions documented
- Integration points identified
- Test strategy outlined
- Implementation report

**Limitations:**
- ⚠️ Files don't persist (sandbox) - Director must materialize
- May not have access to entire codebase
- Can't run builds or tests
- Can't commit changes

### 3. Verifier Agent

**Specialization:** Testing and quality assurance

**Best For:**
- Creating comprehensive test suites
- Running integration tests
- Code coverage analysis
- Finding edge cases
- Regression testing
- Performance validation

**Input Requirements:**
- Code to test (file paths)
- Test requirements/coverage targets
- Existing test patterns
- Edge cases to consider

**Output:**
- Test files created
- Test execution results
- Coverage reports
- Issues found
- Recommendations for improvements

**Limitations:**
- Focuses on testing, not implementation
- May need Builder output first

### 4. Scribe Agent

**Specialization:** Documentation and guides

**Best For:**
- User guides
- API documentation
- Architecture documentation
- README updates
- CHANGELOG entries
- Migration guides

**Input Requirements:**
- What was implemented (from Builder)
- User-facing features
- Technical details to document
- Target audience (users vs developers)
- Existing documentation style

**Output:**
- Markdown documentation files
- Code comments (if needed)
- Diagrams (text-based)
- Examples and tutorials
- Updated README/CHANGELOG

**Limitations:**
- Needs completed implementation to document
- Can't verify technical accuracy without context

### 5. Steward Agent

**Specialization:** Architecture review and oversight

**Best For:**
- Reviewing proposed designs
- Ensuring architecture compliance
- Identifying technical debt
- Long-term maintainability review
- Pattern consistency checks
- Security considerations

**Input Requirements:**
- Proposed design or implementation
- Architecture constraints
- Project standards
- Long-term roadmap context

**Output:**
- Architecture review report
- Approval or concerns
- Alternative approaches
- Technical debt assessment
- Recommendations

**Limitations:**
- Advisory role, doesn't implement
- Needs complete design to review

## Director Orchestration Workflows

### Workflow 1: New Feature Implementation (Full 5-Agent)

**Use When:** Implementing a major roadmap phase (like Week 11 TUI skeleton)

**Steps:**

#### 1. Director: Planning Phase
```bash
# Read roadmap
Read docs/v3-implementation-roadmap.md

# Analyze current state
Glob **/services/*.go
Read internal/core/services/workspace_service.go

# Create todo list
TodoWrite [
    {"content": "Delegate design to Builder", "status": "pending", ...},
    {"content": "Review Builder output", "status": "pending", ...},
    {"content": "Materialize files", "status": "pending", ...},
    {"content": "Delegate testing to Verifier", "status": "pending", ...},
    {"content": "Delegate docs to Scribe", "status": "pending", ...},
    {"content": "Request Steward review", "status": "pending", ...},
    {"content": "Commit integrated work", "status": "pending", ...}
]
```

#### 2. Director → Builder: Delegation
```bash
Task(
    subagent_type: "builder",
    description: "Implement TUI adapter skeleton",
    prompt: """
    You are implementing Phase 4 Week 11 of Ticketr v3: TUI adapter skeleton.

    **Context:**
    - Ticketr is a Jira-Markdown sync tool
    - Module: github.com/karolswdev/ticktr
    - Architecture: Hexagonal (ports and adapters)
    - Phase 3 complete (PathResolver implemented)
    - WorkspaceService exists at internal/core/services/workspace_service.go

    **Your Task:**
    Implement TUI adapter skeleton using tview framework.

    **Requirements from roadmap:**
    1. Add tview dependency (go get github.com/rivo/tview)
    2. Create internal/adapters/tui/ structure
    3. Implement View interface pattern
    4. Create Router for navigation
    5. Implement workspace list view (real data)
    6. Create placeholder ticket views
    7. Add global keybindings (q, ?, Ctrl+C)
    8. Integrate with CLI (ticketr tui command)

    **Architecture Constraints:**
    - TUI adapter depends ONLY on services (not repositories)
    - Use constructor injection for dependencies
    - Follow existing patterns (see workspace_commands.go)
    - No breaking changes to existing code

    **Acceptance Criteria:**
    [ ] tview added to go.mod
    [ ] TUIApp struct with WorkspaceService + PathResolver
    [ ] View interface defined (Name, Primitive, OnShow, OnHide)
    [ ] Router manages view lifecycle
    [ ] Workspace list shows real data
    [ ] Basic keybindings work (q, ?, Ctrl+C)
    [ ] CLI command: ticketr tui

    **What to Return:**
    1. Complete file contents for all files
    2. Architecture decisions made
    3. Code snippets for key integration points
    4. Line counts per file
    5. Any deviations from plan with rationale
    6. Next steps for Week 12

    Begin implementation. You have full autonomy within constraints.
    """
)
```

**Key Elements of Builder Prompt:**
- ✅ **Context:** What exists, what project is
- ✅ **Specific Task:** Clear outcome
- ✅ **Requirements:** From roadmap, explicit list
- ✅ **Constraints:** Architecture rules
- ✅ **Acceptance Criteria:** Checklist
- ✅ **Expected Output:** What to return
- ✅ **Autonomy:** "Begin implementation" (don't wait for more input)

#### 3. Director: Review Builder Output

When Builder returns (check with eye icon or wait for completion):

```bash
# Builder provides detailed report with:
# - Files created (but not persisted - sandbox!)
# - Code contents
# - Design decisions
# - Line counts
# - Deviations explained

# Director reviews:
# 1. Does it meet acceptance criteria?
# 2. Are architecture constraints followed?
# 3. Is code quality acceptable?
# 4. Any red flags?

# Update todo
TodoWrite [
    {"content": "Delegate design to Builder", "status": "completed", ...},
    {"content": "Review Builder output", "status": "completed", ...},
    {"content": "Materialize files", "status": "in_progress", ...},
    ...
]
```

#### 4. Director: Materialize Files

**Critical:** Builder's files exist only in sandbox. Director must create them:

```bash
# Add dependencies first
Bash: go get github.com/rivo/tview@latest

# Create directory structure
Bash: mkdir -p internal/adapters/tui/views

# Materialize each file from Builder's design
Write(
    file_path: /home/karol/dev/code/ticketr/internal/adapters/tui/views/view.go,
    content: [exact content from Builder's report]
)

Write(
    file_path: /home/karol/dev/code/ticketr/internal/adapters/tui/app.go,
    content: [exact content from Builder's report]
)

# Repeat for all files...

# Verify files exist
Bash: ls -la internal/adapters/tui/
Bash: go build ./cmd/ticketr
```

**Why This Step Exists:**
Agents run in sandboxes - their file operations don't persist to real filesystem. Director has real filesystem access.

#### 5. Director → Verifier: Testing

```bash
Task(
    subagent_type: "verifier",
    description: "Test TUI implementation",
    prompt: """
    You are verifying the TUI adapter implementation for Ticketr v3.

    **What Was Implemented:**
    - TUI adapter at internal/adapters/tui/
    - 11 files, 694 lines of code
    - View interface, Router, KeyHandler
    - Workspace list view with live data
    - Help view, placeholder ticket views
    - CLI integration: ticketr tui

    **Your Task:**
    Create comprehensive tests and validate implementation.

    **Test Requirements:**
    1. Unit tests for TUIApp initialization
    2. Router view registration and navigation tests
    3. Workspace list view tests with mocks
    4. Keybinding handler tests
    5. Integration test: full TUI initialization
    6. Coverage target: 80%+

    **Existing Patterns:**
    - Mock pattern: See internal/core/services/*_test.go
    - Table-driven tests preferred
    - Use testify for assertions

    **What to Return:**
    1. Test files created (with contents)
    2. Test execution results
    3. Coverage report
    4. Issues found (if any)
    5. Recommendations for improvement

    Run all tests and report results.
    """
)
```

#### 6. Director: Review & Integrate Tests

```bash
# Verifier returns test results
# Director materializes test files
Write: internal/adapters/tui/app_test.go
Write: internal/adapters/tui/router_test.go
# etc...

# Run tests yourself
Bash: go test ./internal/adapters/tui/... -v -cover

# Update todo
TodoWrite: Mark "Delegate testing to Verifier" as completed
```

#### 7. Director → Scribe: Documentation

```bash
Task(
    subagent_type: "scribe",
    description: "Document TUI implementation",
    prompt: """
    You are documenting the TUI implementation for Ticketr v3.

    **What Was Implemented:**
    - Terminal User Interface using tview
    - Workspace list view (live data)
    - Navigation: q (quit), ? (help), Tab (future)
    - Architecture: TUI adapter → Services (hexagonal)
    - Files: internal/adapters/tui/ (11 files)

    **Your Task:**
    Create user-facing and developer documentation.

    **Documents Needed:**
    1. docs/TUI-USER-GUIDE.md
       - How to launch TUI
       - Keyboard shortcuts reference
       - Feature overview
       - Screenshots (ASCII art)

    2. docs/TUI-ARCHITECTURE.md
       - View interface pattern
       - Router design
       - Integration with services
       - Extension guide for Week 12+

    3. Update docs/README.md
       - Add TUI section to features
       - Add 'ticketr tui' to quick start

    4. Update CHANGELOG.md
       - Add Phase 4 Week 11 entry

    **Target Audience:**
    - User guide: End users (CLI familiarity assumed)
    - Architecture: Developers extending TUI

    **Style:**
    - Clear, concise
    - Code examples where helpful
    - Markdown formatting

    **What to Return:**
    All document contents with line counts.
    """
)
```

#### 8. Director: Review & Integrate Docs

```bash
# Scribe returns documentation
# Director materializes
Write: docs/TUI-USER-GUIDE.md
Write: docs/TUI-ARCHITECTURE.md

# Review and edit if needed
Edit: docs/README.md (add TUI section)

# Update todo
TodoWrite: Mark "Delegate docs to Scribe" as completed
```

#### 9. Director → Steward: Architecture Review

```bash
Task(
    subagent_type: "steward",
    description: "Review TUI architecture",
    prompt: """
    You are reviewing the TUI adapter architecture for Ticketr v3.

    **What Was Implemented:**
    - TUI adapter at internal/adapters/tui/
    - Depends on WorkspaceService, PathResolver (services)
    - View interface pattern for polymorphism
    - Router pattern for lifecycle management
    - Keybinding handler for global events

    **Your Task:**
    Review architecture for compliance and long-term maintainability.

    **Review Criteria:**
    1. Hexagonal architecture compliance
       - Does TUI depend only on services?
       - Are boundaries clean?
    2. Dependency injection usage
       - Constructor injection used?
       - Testability?
    3. Pattern consistency
       - Follows existing project patterns?
    4. Extensibility
       - Easy to add new views?
       - Ready for Week 12 multi-panel?
    5. Technical debt
       - Any shortcuts taken?
       - Future refactoring needs?

    **Project Constraints:**
    - Must follow hexagonal architecture
    - No direct database/external access from adapters
    - Clean dependency graph (no cycles)

    **What to Return:**
    1. Architecture approval (yes/no/conditional)
    2. Compliance assessment
    3. Issues found (if any)
    4. Recommendations
    5. Technical debt assessment

    Provide detailed review.
    """
)
```

#### 10. Director: Final Integration & Commit

```bash
# Steward returns review
# If approved or minor issues:

# Build verification
Bash: go build ./...
Bash: ./ticketr tui --help

# Final todo update
TodoWrite: Mark all items completed

# Git commit
Bash: git add internal/adapters/tui/ cmd/ticketr/tui_command.go go.mod go.sum
Bash: git commit -m "feat(tui): Implement Phase 4 Week 11 - TUI adapter skeleton

[Detailed commit message with architecture, features, technical details]

Generated with [Claude Code](https://claude.ai/code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"
```

### Workflow 2: Incremental Enhancement (Partial 5-Agent)

**Use When:** Adding feature to existing code (like Week 12 multi-panel)

**Abbreviated Steps:**

```
1. Director: Read roadmap, analyze current state
2. Director: Create todo list
3. Director → Builder: "Enhance TUI with multi-panel layout"
   (OR Director implements directly if simple)
4. Director: Materialize files
5. Director: Build and manual test
6. Director → Verifier: "Test multi-panel focus switching" (optional)
7. Director → Scribe: "Update TUI docs for multi-panel" (optional)
8. Director: Commit

Skip Steward if changes are incremental and pattern-compliant.
```

### Workflow 3: Bug Fix (Director Solo)

**Use When:** Fixing bugs, small changes

**Steps:**

```
1. Director: Identify issue
2. Director: Create fix (Edit tool)
3. Director: Test fix (Bash)
4. Director: Commit

No agent delegation needed for small fixes.
```

## Constructing Effective Delegation Prompts

### Template for Builder Agent

```markdown
You are implementing <SPECIFIC_FEATURE> for <PROJECT_NAME>.

**Context:**
- Project description: <what it does>
- Module name: <go module path>
- Architecture: <pattern used>
- Current phase: <roadmap phase>
- Existing components: <what exists>
- What doesn't exist: <what you're building>

**Your Task:**
<One clear sentence describing the outcome>

**Requirements from roadmap:**
<Copy exact requirements from roadmap>
1. Requirement 1
2. Requirement 2
...

**Architecture Constraints:**
- Constraint 1 (e.g., "Use hexagonal architecture")
- Constraint 2 (e.g., "No direct DB access")
- Constraint 3 (e.g., "Constructor injection only")

**Existing Patterns to Follow:**
- Pattern 1: <file reference>
- Pattern 2: <code example>

**Acceptance Criteria:**
[ ] Criterion 1
[ ] Criterion 2
[ ] Criterion 3
...

**What to Return:**
1. All file contents (complete code)
2. Architecture decisions made and why
3. Key code snippets for integration
4. Line counts per file
5. Any deviations from plan with justification
6. Next steps or recommendations

**Important:**
- You have full autonomy within the constraints
- Make design decisions and document them
- Prioritize clean, maintainable code
- Follow existing project style

Begin implementation now.
```

### Template for Verifier Agent

```markdown
You are verifying <WHAT_WAS_IMPLEMENTED> for <PROJECT_NAME>.

**What Was Implemented:**
- Feature description
- Files created/modified: <list>
- Line counts: <X lines>
- Key components: <list>

**Your Task:**
Create comprehensive tests and validate the implementation.

**Test Requirements:**
1. Unit tests for <component 1>
2. Integration tests for <workflow>
3. Coverage target: <percentage>
4. Edge cases: <list>
5. Performance: <if applicable>

**Existing Test Patterns:**
- See: <reference test file>
- Framework: <testing framework used>
- Mocking: <mocking approach>

**What to Return:**
1. Test files created (with complete contents)
2. Test execution results (pass/fail counts)
3. Coverage report
4. Issues found (bugs, edge cases missed)
5. Recommendations for improvement

Run all tests and report comprehensive results.
```

### Template for Scribe Agent

```markdown
You are documenting <WHAT_WAS_IMPLEMENTED> for <PROJECT_NAME>.

**What Was Implemented:**
- Feature summary
- User-facing functionality
- Technical details
- Files involved

**Your Task:**
Create documentation for <TARGET_AUDIENCE>.

**Documents Needed:**
1. <Document 1>
   - Purpose: <why>
   - Audience: <who>
   - Key sections: <what to cover>

2. <Document 2>
   - Purpose: <why>
   - Sections: <what>

**Existing Documentation Style:**
- See: <reference doc>
- Tone: <formal/casual>
- Format: <Markdown/etc>

**What to Return:**
1. All document contents (complete markdown)
2. Line counts per document
3. Any sections needing developer input
4. Recommendations for additional docs

Provide complete, ready-to-use documentation.
```

### Template for Steward Agent

```markdown
You are reviewing <WHAT_WAS_IMPLEMENTED> for <PROJECT_NAME>.

**What Was Implemented:**
- Architecture summary
- Key design decisions
- Files created/modified
- Patterns used

**Your Task:**
Review architecture for compliance and long-term maintainability.

**Review Criteria:**
1. <Architecture pattern> compliance
2. Dependency management
3. Pattern consistency with existing code
4. Extensibility for future features
5. Technical debt assessment

**Project Standards:**
- Architecture: <pattern>
- Dependency rules: <rules>
- Code standards: <standards>

**What to Return:**
1. Approval status (approved / approved with conditions / rejected)
2. Compliance assessment per criterion
3. Issues found with severity (critical/major/minor)
4. Recommendations for improvement
5. Technical debt introduced (if any)

Provide detailed, actionable review.
```

## Handling Agent Outputs

### Builder Output Processing

**Builder Returns:**
- Report with file contents
- Design decisions
- Line counts

**Director Actions:**
1. **Review:** Read entire report, check acceptance criteria
2. **Question:** If anything unclear, ask Builder to clarify
3. **Materialize:** Create files using Write tool
4. **Verify:** Build and basic functionality test
5. **Accept/Reject:** Mark todo as complete or request revisions

**Example:**
```bash
# Builder finished, review output
# Look for: "Builder Deliverable: TUI Adapter Skeleton"

# Check acceptance criteria
✅ tview added to go.mod
✅ TUIApp struct implements initialization
✅ Can launch with ticketr tui
... (all checked)

# Materialize files
Write(file_path: internal/adapters/tui/app.go, content: [from Builder])
Write(file_path: internal/adapters/tui/router.go, content: [from Builder])
# ... repeat for all files

# Verify
Bash: go build ./cmd/ticketr
Bash: ./ticketr tui --help

# Update todo
TodoWrite: Mark "Delegate design to Builder" as completed
```

### Verifier Output Processing

**Verifier Returns:**
- Test files
- Execution results
- Coverage report
- Issues found

**Director Actions:**
1. **Review Results:** Check pass/fail, coverage percentage
2. **Materialize Tests:** Create test files
3. **Re-run Tests:** Verify locally
4. **Address Issues:** Fix any bugs found
5. **Accept:** Mark testing complete

**Red Flags:**
- ⚠️ Coverage < target (investigate why)
- ⚠️ Tests failing (fix before proceeding)
- ⚠️ Critical issues found (may need Builder revision)

### Scribe Output Processing

**Scribe Returns:**
- Documentation files
- README/CHANGELOG updates

**Director Actions:**
1. **Review Docs:** Read for accuracy and clarity
2. **Fact-Check:** Ensure technical details correct
3. **Materialize:** Create documentation files
4. **Edit If Needed:** Minor corrections with Edit tool
5. **Accept:** Mark documentation complete

**Common Edits:**
- Update version numbers
- Correct file paths
- Add missing sections
- Adjust tone if needed

### Steward Output Processing

**Steward Returns:**
- Approval status
- Compliance report
- Recommendations

**Director Actions:**
1. **Check Approval:** Approved, conditional, or rejected?
2. **Review Issues:** Severity and impact
3. **Decide:**
   - Approved: Proceed to commit
   - Conditional: Address minor issues, then proceed
   - Rejected: Major rework needed, back to Builder
4. **Document:** Note any technical debt or future work

## Parallel vs Sequential Delegation

### Parallel Delegation (Faster)

**Use When:** Tasks are independent

**Example:**
```bash
# Week 11: After implementation
Task(builder: "Implement TUI skeleton")  # Wait for completion first

# Then parallel:
Task(verifier: "Test TUI components")    # Can run together
Task(scribe: "Document TUI usage")       # Independent of tests
Task(steward: "Review TUI architecture") # Independent of tests

# All three agents work simultaneously
```

**Syntax:**
```bash
# Single message with multiple Task calls
Task(subagent_type: "verifier", ...)
Task(subagent_type: "scribe", ...)
Task(subagent_type: "steward", ...)
```

### Sequential Delegation (Dependencies)

**Use When:** Later task needs earlier output

**Example:**
```bash
# Must be sequential
Task(builder: "Implement TUI skeleton")
# Wait for Builder to finish

Task(verifier: "Test TUI implementation")
# Wait for Verifier to finish

Task(scribe: "Document tested features")
# Scribe needs to know what passed tests
```

**Syntax:**
```bash
# Separate messages, wait for completion between each
Task(subagent_type: "builder", ...)
# [Wait for completion, review output]

Task(subagent_type: "verifier", ...)
# [Wait for completion, review output]

Task(subagent_type: "scribe", ...)
```

## Decision Matrix: Delegate or Do It Yourself?

### Delegate to Builder When:
- ✅ Multi-file implementation (3+ files)
- ✅ Complex business logic
- ✅ New architecture patterns
- ✅ Unfamiliar technology
- ✅ Time to review > time to implement

### Do It Yourself When:
- ✅ Single file change
- ✅ Simple refactoring
- ✅ Configuration updates
- ✅ Pattern you know well
- ✅ < 50 lines of code

### Delegate to Verifier When:
- ✅ Need comprehensive test coverage
- ✅ Complex edge cases to find
- ✅ Integration testing required
- ✅ You're unsure what to test

### Write Tests Yourself When:
- ✅ Simple unit tests
- ✅ Testing your own code immediately
- ✅ TDD approach (test-first development)

### Delegate to Scribe When:
- ✅ User-facing documentation needed
- ✅ Multiple documents to create
- ✅ Don't know documentation structure

### Write Docs Yourself When:
- ✅ Quick README update
- ✅ Code comments
- ✅ Simple CHANGELOG entry

### Delegate to Steward When:
- ✅ Major architectural changes
- ✅ Unsure if approach is sound
- ✅ Want second opinion on design
- ✅ Complex compliance requirements

### Skip Steward When:
- ✅ Following established patterns exactly
- ✅ Trivial changes
- ✅ Incremental enhancements

## The Sandbox Limitation Workaround

### The Problem

Agents run in isolated sandboxes. File operations (Write, Edit) in agent context don't persist to real filesystem.

**Symptoms:**
- Agent reports "Created file X" but `ls` shows nothing
- Agent says "All files ready" but build fails with "file not found"
- Agent provides complete code but you can't run it

### The Solution

**Director must materialize agent-created files:**

```bash
# Agent (Builder) executes:
Write(file_path: /path/to/file.go, content: "package main...")
# ⚠️ This file exists in AGENT'S sandbox only

# Director must execute:
Write(file_path: /path/to/file.go, content: "package main...")
# ✅ This file exists in REAL filesystem
```

### Materialization Workflow

**Step 1: Agent Returns Report**
```
Builder provides:
---
## Files Created

### internal/adapters/tui/app.go (116 lines)
```go
package tui

import (...)

type TUIApp struct {
    ...
}
...
```

### internal/adapters/tui/router.go (78 lines)
```go
package tui

type Router struct {
    ...
}
...
```
---
```

**Step 2: Director Materializes Each File**
```bash
# Create directories first
Bash: mkdir -p internal/adapters/tui/views

# Materialize file 1
Write(
    file_path: /home/karol/dev/code/ticketr/internal/adapters/tui/app.go,
    content: """package tui

import (...)

type TUIApp struct {
    ...
}
..."""
)

# Materialize file 2
Write(
    file_path: /home/karol/dev/code/ticketr/internal/adapters/tui/router.go,
    content: """package tui

type Router struct {
    ...
}
..."""
)

# Repeat for ALL files agent created
```

**Step 3: Verify Materialization**
```bash
# Check files exist
Bash: ls -la internal/adapters/tui/
# Should show: app.go, router.go, etc.

# Verify contents
Read: internal/adapters/tui/app.go
# Should match agent's design

# Build test
Bash: go build ./...
# Should compile successfully
```

### Optimization: Batch Materialization

For many files, create a checklist:

```markdown
Files to materialize from Builder:
- [ ] internal/adapters/tui/app.go
- [ ] internal/adapters/tui/router.go
- [ ] internal/adapters/tui/keybindings.go
- [ ] internal/adapters/tui/views/view.go
- [ ] internal/adapters/tui/views/workspace_list.go
- [ ] internal/adapters/tui/views/ticket_tree.go
- [ ] internal/adapters/tui/views/ticket_detail.go
- [ ] internal/adapters/tui/views/help.go
- [ ] cmd/ticketr/tui_command.go

Use Write tool for each, copying exact content from Builder report.
```

## Integration Checklist

After all agents complete, before commit:

```markdown
Integration Checklist:
- [ ] All Builder files materialized
- [ ] Build succeeds: `go build ./...`
- [ ] All Verifier tests materialized
- [ ] Tests pass: `go test ./...`
- [ ] Coverage meets target: `go test -cover`
- [ ] All Scribe docs materialized
- [ ] Docs reviewed for accuracy
- [ ] Steward issues addressed (if any)
- [ ] Manual functionality test passed
- [ ] Todo list all items completed
- [ ] Git status clean (no unexpected changes)
- [ ] Commit message prepared
```

## Example Director Session (Week 11)

### Actual Commands Used

```bash
# 1. Planning
Read: docs/v3-implementation-roadmap.md (lines 272-369)
Glob: **/services/*.go
Read: internal/core/services/workspace_service.go
Read: go.mod

TodoWrite: [11 tasks for Week 11]

# 2. Delegation
Task(
    subagent_type: "builder",
    description: "Implement TUI adapter skeleton",
    prompt: [comprehensive prompt with context, requirements, constraints]
)

# 3. Wait for Builder (agent runs autonomously)
# Builder completes, provides detailed report

# 4. Review Builder Output
# Read Builder's report
# Verify acceptance criteria met
# Check design decisions

# 5. Materialize Files
Bash: go get github.com/rivo/tview@latest
Bash: mkdir -p internal/adapters/tui/views

Write: internal/adapters/tui/views/view.go [content from Builder]
Write: internal/adapters/tui/views/workspace_list.go [content]
Write: internal/adapters/tui/views/ticket_tree.go [content]
Write: internal/adapters/tui/views/ticket_detail.go [content]
Write: internal/adapters/tui/views/help.go [content]
Write: internal/adapters/tui/router.go [content]
Write: internal/adapters/tui/keybindings.go [content]
Write: internal/adapters/tui/app.go [content]
Write: cmd/ticketr/tui_command.go [content]

TodoWrite: Mark "Create TUI adapter skeleton" completed

# 6. Verification
Bash: go build ./cmd/ticketr
Bash: ./ticketr tui --help

# 7. Testing (skipped Verifier, did manually)
# Could have delegated: Task(subagent_type: "verifier", ...)

# 8. Documentation (skipped Scribe, did manually)
# Could have delegated: Task(subagent_type: "scribe", ...)

# 9. Architecture Review (skipped Steward, simple pattern)

# 10. Commit
Bash: git add internal/adapters/tui/ cmd/ticketr/tui_command.go go.mod go.sum
Bash: git commit -m "feat(tui): Implement Phase 4 Week 11..."

TodoWrite: Mark all tasks completed
```

### Time Breakdown

- Planning (Director): 15 min
- Builder delegation: 5 min
- Builder execution: 20 min (autonomous, parallel work time)
- Review Builder output: 10 min
- Materialize files: 30 min
- Verification: 15 min
- Commit: 5 min

**Total:** ~2 hours
**Without Builder:** Estimated 4+ hours (implementation from scratch)

## Common Orchestration Mistakes

### Mistake 1: Not Materializing Agent Files

**Problem:**
```bash
Task(builder: "Create TUI skeleton")
# Builder completes
# Immediately try to build
Bash: go build ./...
# ERROR: files not found
```

**Solution:**
```bash
Task(builder: "Create TUI skeleton")
# Builder completes
# MATERIALIZE FILES FIRST
Write: [all files from Builder report]
# THEN build
Bash: go build ./...
```

### Mistake 2: Incomplete Delegation Prompts

**Problem:**
```bash
Task(builder: "Add TUI to Ticketr")
# Too vague, Builder confused
```

**Solution:**
```bash
Task(builder: "Implement TUI adapter skeleton",
    prompt: """
    Complete context...
    Specific requirements...
    Architecture constraints...
    Acceptance criteria...
    What to return...
    """)
```

### Mistake 3: Not Reviewing Agent Output

**Problem:**
```bash
Task(builder: "Implement feature X")
# Builder completes
# Immediately materialize without reading
Write: [blindly copy everything]
# Later find bugs or misunderstandings
```

**Solution:**
```bash
Task(builder: "Implement feature X")
# Builder completes
# READ THE REPORT THOROUGHLY
# Check acceptance criteria
# Verify design decisions
# Question anything unclear
# THEN materialize
```

### Mistake 4: Sequential When Could Be Parallel

**Problem:**
```bash
Task(verifier: "Test TUI")
# Wait for completion
Task(scribe: "Document TUI")
# Wait for completion
# Wasted time - these are independent
```

**Solution:**
```bash
# Single message with both
Task(verifier: "Test TUI")
Task(scribe: "Document TUI")
# Both run in parallel
```

### Mistake 5: Over-Delegation

**Problem:**
```bash
# Simple one-line change
Task(builder: "Add comment to file")
# Overkill, slower than doing it yourself
```

**Solution:**
```bash
# Just do it
Edit(file_path: file.go, old: "func Foo", new: "// Foo does X\nfunc Foo")
```

## Measuring Orchestration Success

### Velocity Metrics

**Good Orchestration:**
- Feature delivered in 1-2 sessions
- All acceptance criteria met
- Build succeeds first time
- Tests pass with good coverage
- Documentation complete
- Clean commit history

**Poor Orchestration:**
- Multiple sessions for simple features
- Half-finished deliverables
- Build failures after agent work
- Missing tests or docs
- Confusing commit messages

### Quality Metrics

**Good Orchestration:**
- Architecture compliant (Steward approved)
- Test coverage > 80%
- No technical debt introduced
- Patterns consistent with existing code

**Poor Orchestration:**
- Architecture violations
- Low test coverage
- Technical debt accumulating
- Inconsistent patterns

### Efficiency Metrics

**Good Orchestration:**
- Parallel agent usage when possible
- Minimal rework cycles
- Clear delegation prompts
- Quick review and integration

**Poor Orchestration:**
- Sequential when could be parallel
- Multiple rework cycles
- Vague prompts requiring clarification
- Long review/integration time

## Summary: Director's Responsibilities

As Director, you are the conductor of the orchestra:

1. **Plan:** Break down roadmap into tasks
2. **Delegate:** Assign tasks to specialist agents
3. **Provide Context:** Comprehensive delegation prompts
4. **Monitor:** Track progress with TodoWrite
5. **Review:** Critically evaluate agent outputs
6. **Materialize:** Create files from agent designs (sandbox workaround)
7. **Integrate:** Combine work from multiple agents
8. **Verify:** Build, test, and validate
9. **Commit:** Record work with detailed messages
10. **Document Process:** Maintain handover docs

**You are NOT:**
- Writing all code yourself
- A passive observer of agents
- Accepting agent output without review
- Just a file copier

**You ARE:**
- Orchestrating specialists
- Making final decisions
- Ensuring quality
- Maintaining architecture
- The integration point

Use the 5-agent methodology to **multiply your effectiveness**, not replace your judgment.

---

**Roadmap Progress:** Phase 4 Week 11 complete using this methodology
**Next:** Apply same orchestration approach to Week 12
**Reference:** This guide for all future phases
