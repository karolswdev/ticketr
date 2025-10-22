# Boot to Phase 7: Jira Library Integration
## Zero-Context Director Onboarding Document

**Purpose:** This document provides complete context for a fresh Director agent (or human) to execute Phase 7 of the Ticketr project with zero prior knowledge.

**Last Updated:** 2025-10-21
**Status:** ✅ READY FOR EXECUTION
**External Validation:** ✅ APPROVED by Gemini + Codex AI architects
**Current Branch:** `feature/jira-domain-redesign`

---

## Table of Contents

1. [What is Ticketr?](#what-is-ticketr)
2. [What is Phase 7?](#what-is-phase-7)
3. [Why Are We Doing This?](#why-are-we-doing-this)
4. [What Has Already Been Done?](#what-has-already-been-done)
5. [What Must Be Done Now?](#what-must-be-done-now)
6. [How to Execute Phase 7](#how-to-execute-phase-7)
7. [Critical Context](#critical-context)
8. [Success Criteria](#success-criteria)
9. [Troubleshooting](#troubleshooting)
10. [Appendix](#appendix)

---

## What is Ticketr?

### 30-Second Overview

**Ticketr** is a production-ready Go CLI tool (v3.1.1) that synchronizes Markdown files with Jira tickets bidirectionally.

```
Markdown Files ⟷ Ticketr ⟷ Jira Cloud
```

**Key Stats:**
- **Language:** Go 1.22+
- **Code:** ~16,000 lines
- **Tests:** 147 tests, 74.8% coverage
- **Architecture:** Hexagonal (Ports & Adapters)
- **Users:** 1 (single user, no migration concerns)
- **Current Version:** v3.1.1

### Core Functionality

```bash
# Pull Jira tickets → Generate Markdown
ticketr pull --project PROJ

# Push Markdown → Create/Update Jira tickets
ticketr push tickets.md
```

### Architecture Pattern

Ticketr follows **Hexagonal Architecture (Ports & Adapters)**:

```
┌─────────────────────────────────────────────────────────┐
│                      CLI Layer                          │
│                   (cmd/ticketr/)                        │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│                   Adapters Layer                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐            │
│  │Filesystem│  │   Jira   │  │ Keychain │            │
│  │ Adapter  │  │ Adapter  │  │  Store   │            │
│  └──────────┘  └──────────┘  └──────────┘            │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│                    Ports Layer                          │
│   Interfaces: Repository, JiraPort, CredentialStore    │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│                Core Business Logic                      │
│  WorkspaceService, PullService, PushService            │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│                   Domain Layer                          │
│   Models: Ticket, Task, Workspace, WorkspaceConfig     │
└─────────────────────────────────────────────────────────┘
```

**Key Principle:** Business logic depends on **interfaces (ports)**, not concrete implementations (adapters). This allows swapping infrastructure components (like the Jira client) without touching core business logic.

---

## What is Phase 7?

### Mission Statement

**Phase 7:** Complete production-ready integration of `andygrunwald/go-jira` v1.17.0 library into Ticketr's Jira adapter, replacing the custom HTTP client.

### The Problem

**Current State (V1):** Custom HTTP client in `internal/adapters/jira/jira_adapter.go`
- **Lines of Code:** 1,136 lines
- **Implementation:** Manual `http.NewRequest()` calls, manual JSON marshaling, custom pagination
- **Maintenance Burden:** HIGH - must manually handle all Jira API changes
- **Dependencies:** 12 (all transitive, no Jira library)

### The Solution

**Proposed State (V2):** Use battle-tested library `andygrunwald/go-jira` v1.17.0
- **Lines of Code:** 757 lines (33% reduction, -379 lines)
- **Implementation:** Library handles HTTP/JSON/pagination automatically
- **Maintenance Burden:** LOW - library maintained by community (868 importers)
- **Dependencies:** 24 (+12 from library, acceptable tradeoff)

### Strategic Value

This is a **technical debt reduction** initiative with long-term sustainability benefits:
1. Less code to maintain
2. More reliable (battle-tested library)
3. Faster feature development
4. Better error handling
5. Focus on business logic, not HTTP plumbing

---

## Why Are We Doing This?

### Historical Context

**Phase 6.5:** Emergency UAT fixes revealed systemic issues with the Jira integration layer. User demanded:
> "I DEMAND A THOROUGH REVIEW PHASE OF THE CURRENT JIRA SERVICE, JIRA ADAPTER. I WANT YOU TO QUESTION EVERYTHING IN THAT PHASE, THE DESIGN, ARE WE ROBUST ENOUGH?"

**Phase 7 Launch:** Complete architectural review of Jira domain led to discovery:
- Custom HTTP client is 1,136 lines of maintenance burden
- Go ecosystem has mature, proven library: `andygrunwald/go-jira`
- Hexagonal architecture makes adapter swap trivial

### External Validation Mandate

User explicitly required world-class validation:
> "Sure, but I still want you to consult codex and gemini... Make sure you invoke them with enough context so they can speak intelligently about this stuff. Look, this needs to be world class."

**Result:** Both Gemini and Codex AI architects **unanimously approved** the library integration.

### Decision Drivers

1. **Maturity:** 9 years production use (2015-2025)
2. **Adoption:** 868 packages use it, 1,600 GitHub stars
3. **Security:** 0 CVEs, clean `govulncheck` scan
4. **Code Reduction:** 33% less code (-379 lines)
5. **Architecture Fit:** Respects ports & adapters boundaries
6. **Risk Mitigation:** Feature flag allows instant rollback to V1

---

## What Has Already Been Done?

### ✅ Completed Work

#### 1. Research & Validation (Steward Agent)

**Files Created:**
- `research/LIBRARY-INTEGRATION-RESEARCH.md` - Initial library research
- `research/EXTERNAL-VALIDATION-REPORT.md` - 15+ source validation
- `research/jira_adapter_v2_example.go` - Proof-of-concept code

**Findings:**
- Evaluated `andygrunwald/go-jira` vs `ctreminiom/go-atlassian`
- Chose `andygrunwald/go-jira` (better fit for our needs)
- Validated: 0 CVEs, 868 importers, 9 years stable

#### 2. External AI Consultation (Director)

**Consultations Completed:**
- ✅ **Gemini CLI:** "Architectural Blessing" - approved with detailed implementation checklist
- ✅ **Codex CLI:** "Pragmatic, battle-tested choice" - validated decision

**Key Recommendations:**
- Both unanimously approved library integration
- Identified 5 critical production gaps (ADR, docs, tests, logging, CI)
- Provided 8-item production readiness checklist

**Evidence:**
- Gemini full consultation output received (2,474 tokens)
- Codex full consultation output received (detailed assessment)

#### 3. V2 Implementation (Builder Agent)

**Files Created:**
- `internal/adapters/jira/jira_adapter_v2.go` - 757 lines, full JiraPort implementation
- `internal/adapters/jira/factory.go` - 68 lines, feature flag system
- `internal/adapters/jira/jira_adapter_v2_test.go` - 37 new tests (all passing)

**Implementation Quality:**
- ✅ Implements same `JiraPort` interface as V1 (drop-in replacement)
- ✅ All 37 tests passing (100% success rate)
- ✅ Uses proper credential flow via `WorkspaceConfig`
- ✅ Feature flag: `TICKETR_JIRA_ADAPTER_VERSION` (values: v1, v2, default: v2)

#### 4. Dependency Addition

**File Modified:**
- `go.mod` - Added `github.com/andygrunwald/go-jira v1.17.0`

**Dependency Impact:**
- Before: 12 dependencies
- After: 24 dependencies (+12 from library, all vetted)

#### 5. Orchestration Planning (Director)

**File Created:**
- `.agents/phase7/ORCHESTRATION-PLAN.md` - Comprehensive 5-step execution plan

### ❌ Critical Gaps (MUST BE COMPLETED)

These gaps were identified by external AI architects as **blockers to production**:

#### Gap 1: Integration Not Wired
**Problem:** `cmd/ticketr/main.go:142` still calls V1 directly
**Current Code:**
```go
return jira.NewJiraAdapterFromConfig(config, fieldMappings)
```
**Required Code:**
```go
return jira.NewJiraAdapterFromConfigWithVersion(config, fieldMappings)
```
**Status:** ❌ NOT DONE

#### Gap 2: No Architecture Decision Record (ADR)
**Problem:** No formal documentation of this architectural decision
**Required:** `docs/adr/001-adopt-go-jira-library.md`
**Status:** ❌ NOT DONE
**Gemini:** "Non-negotiable"

#### Gap 3: Architecture Docs Outdated
**Problem:** `docs/ARCHITECTURE.md` has no mention of library or V2 adapter
**Required:** Update Jira Adapter section, add dependencies section
**Status:** ❌ NOT DONE

#### Gap 4: Feature Flag Untested in CI
**Problem:** No CI validation that both V1 and V2 adapters work
**Required:** CI runs tests with `TICKETR_JIRA_ADAPTER_VERSION=v1` AND `v2`
**Status:** ❌ NOT DONE
**Gemini:** "Feature flag itself is untested"

#### Gap 5: No Version Logging
**Problem:** Can't distinguish V1 vs V2 errors in production logs
**Required:** Wrap all errors with `[jira-v1]` or `[jira-v2]` prefix
**Status:** ❌ NOT DONE
**Purpose:** Monitor for increased errors post-deployment

---

## What Must Be Done Now?

### Execution Summary

**5 Steps | 3 Specialized Agents | ~75 minutes**

```
Director (You) orchestrates:
├─ Scribe:    Documentation + ADR (Step 1)
├─ Builder:   Code integration + logging (Step 2)
├─ Verifier:  Testing + benchmarks (Steps 3-4)
└─ Director:  Final verification (Step 5)
```

### Step-by-Step Breakdown

#### **STEP 1: Governance & Documentation (Scribe Agent)**
**Duration:** 30 minutes
**Can Run in Parallel:** Yes (with Step 2)

**Tasks:**
1. Create `docs/adr/001-adopt-go-jira-library.md` (ADR)
2. Update `docs/ARCHITECTURE.md` (Jira Adapter section)
3. Update `internal/adapters/jira/README.md` (feature flag docs)

**Deliverable:** 3 files updated/created, governance complete

#### **STEP 2: Code Integration (Builder Agent)**
**Duration:** 20 minutes
**Can Run in Parallel:** Yes (with Step 1)

**Tasks:**
1. Wire factory into `cmd/ticketr/main.go:142`
2. Add `[jira-v1]` / `[jira-v2]` version tags to all errors
3. Verify compilation and existing tests pass

**Deliverable:** Factory integrated, version logging operational

#### **STEP 3: Feature Flag Validation (Verifier Agent)**
**Duration:** 30 minutes
**Depends On:** Step 2 complete

**Tasks:**
1. Create `scripts/test-adapter-versions.sh` (test both adapters)
2. Add CI validation for both V1 and V2
3. Create `internal/adapters/jira/integration_test.go` (V1/V2 parity tests)

**Deliverable:** Feature flag validated in CI

#### **STEP 4: Performance Validation (Verifier Agent)**
**Duration:** 20 minutes
**Depends On:** Step 2 complete

**Tasks:**
1. Create `internal/adapters/jira/adapter_bench_test.go` (benchmarks)
2. Run benchmarks comparing V1 vs V2 performance
3. Validate V2 within 20% of V1 performance

**Deliverable:** Performance validated, no regression

#### **STEP 5: Final Verification (Director - You)**
**Duration:** 15 minutes
**Depends On:** Steps 1-4 complete

**Tasks:**
1. Validate Gemini's 8-item production checklist complete
2. Create `docs/deployment/JIRA-LIBRARY-ROLLOUT.md` (deployment plan)
3. Create `.agents/phase7/PHASE7-COMPLETION-REPORT.md` (final summary)

**Deliverable:** Phase 7 production-ready, deployment plan approved

---

## How to Execute Phase 7

### Prerequisites Check

Before starting, verify these conditions:

```bash
# 1. On correct branch
git branch --show-current
# Expected: feature/jira-domain-redesign

# 2. Clean working directory
git status
# Expected: no uncommitted changes (or only documentation files)

# 3. V2 adapter exists
ls internal/adapters/jira/jira_adapter_v2.go
# Expected: file exists

# 4. Factory exists
ls internal/adapters/jira/factory.go
# Expected: file exists

# 5. Library installed
grep "andygrunwald/go-jira" go.mod
# Expected: github.com/andygrunwald/go-jira v1.17.0

# 6. Tests currently pass
go test ./internal/adapters/jira/...
# Expected: All tests pass
```

**If any check fails:** Review "What Has Already Been Done" section to understand what's missing.

### Execution Sequence

#### **Phase 1: Launch Parallel Agents (Steps 1-2)**

**Invoke Scribe Agent:**
```
Agent: Scribe
Mission: Complete governance and documentation for Jira library integration
Context: Read .agents/phase7/ORCHESTRATION-PLAN.md Step 1
Tasks:
1. Create docs/adr/001-adopt-go-jira-library.md following template in ORCHESTRATION-PLAN
2. Update docs/ARCHITECTURE.md Jira Adapter section with V2 library reference
3. Update internal/adapters/jira/README.md with feature flag documentation

Deliverable: Provide completion report with all 3 files updated/created
```

**Invoke Builder Agent (simultaneously):**
```
Agent: Builder
Mission: Integrate V2 factory into main.go and add version logging
Context: Read .agents/phase7/ORCHESTRATION-PLAN.md Step 2
Tasks:
1. Update cmd/ticketr/main.go line 142 to use NewJiraAdapterFromConfigWithVersion
2. Add [jira-v1] error prefix to all errors in jira_adapter.go
3. Add [jira-v2] error prefix to all errors in jira_adapter_v2.go
4. Verify: go build ./cmd/ticketr succeeds
5. Verify: go test ./internal/adapters/jira/... passes

Deliverable: Provide completion report with factory wired and version logging operational
```

**Director Action:** Wait for BOTH agents to complete and provide handover reports.

#### **Phase 2: Sequential Testing (Steps 3-4)**

**After Phase 1 Complete, Invoke Verifier Agent:**
```
Agent: Verifier
Mission: Validate feature flag system and performance
Context: Read .agents/phase7/ORCHESTRATION-PLAN.md Steps 3-4
Tasks:
1. Create scripts/test-adapter-versions.sh per template in ORCHESTRATION-PLAN
2. Update .github/workflows/test.yml to test both V1 and V2
3. Create internal/adapters/jira/integration_test.go with V1/V2 parity tests
4. Create internal/adapters/jira/adapter_bench_test.go with benchmarks
5. Run benchmarks and verify V2 performance acceptable (within 20% of V1)

Deliverable: Provide test results showing both adapters working, benchmarks complete
```

**Director Action:** Validate Verifier handover report shows all tests passing.

#### **Phase 3: Final Verification (Step 5)**

**Director Actions (You):**

1. **Validate Production Checklist:**
```bash
# Check all items complete:
# [ ] Factory wired into main.go
grep "NewJiraAdapterFromConfigWithVersion" cmd/ticketr/main.go

# [ ] ADR created
ls docs/adr/001-adopt-go-jira-library.md

# [ ] Architecture.md updated
grep "andygrunwald/go-jira" docs/ARCHITECTURE.md

# [ ] CI tests both adapters
grep "TICKETR_JIRA_ADAPTER_VERSION" .github/workflows/test.yml

# [ ] Version logging present
grep "\[jira-v" internal/adapters/jira/jira_adapter*.go

# [ ] Integration tests exist
ls internal/adapters/jira/integration_test.go

# [ ] Benchmarks exist
ls internal/adapters/jira/adapter_bench_test.go

# [ ] All tests pass
go test ./...
```

2. **Create Deployment Plan:**
```
File: docs/deployment/JIRA-LIBRARY-ROLLOUT.md
Template: Use template in ORCHESTRATION-PLAN.md Step 5.2
```

3. **Create Completion Report:**
```
File: .agents/phase7/PHASE7-COMPLETION-REPORT.md
Contents:
- External validation results (Gemini + Codex)
- Implementation changes summary
- Test coverage results
- Performance benchmark results
- Deployment plan reference
- Success metrics validated
```

4. **Final Verification:**
```bash
# Build succeeds
go build ./cmd/ticketr

# All tests pass
go test ./...

# Test feature flag manually
TICKETR_JIRA_ADAPTER_VERSION=v1 go test ./internal/adapters/jira/...
TICKETR_JIRA_ADAPTER_VERSION=v2 go test ./internal/adapters/jira/...

# Both should pass
```

#### **Phase 4: User Handover**

**Report to User:**
```
Phase 7 Complete ✅

External Validation:
- Gemini: Architectural Blessing
- Codex: Pragmatic, battle-tested choice

Implementation:
- Factory wired: cmd/ticketr/main.go:142
- Version logging: Operational ([jira-v1] / [jira-v2])
- Tests: All passing (147 + new integration tests)
- Benchmarks: V2 performance acceptable

Governance:
- ADR created: docs/adr/001-adopt-go-jira-library.md
- Architecture docs updated
- Deployment plan documented

Ready for:
- User Acceptance Testing (UAT)
- Production deployment with V2 as default
- V1 deprecation in v3.2.0 or v3.3.0

Feature Flag:
- Default: v2 (library-based)
- Rollback: export TICKETR_JIRA_ADAPTER_VERSION=v1
```

---

## Critical Context

### The JiraPort Interface

**File:** `internal/core/ports/jira_port.go`

This interface is the **contract** that both V1 and V2 adapters must implement:

```go
type JiraPort interface {
    Authenticate() error
    SearchTickets(ctx context.Context, projectKey, jql string,
                  progressCallback JiraProgressCallback) ([]domain.Ticket, error)
    CreateTicket(ticket domain.Ticket) (string, error)
    UpdateTicket(ticket domain.Ticket) error
    CreateTask(task domain.Task, parentID string) (string, error)
    UpdateTask(task domain.Task) error
    GetProjectIssueTypes() (map[string][]string, error)
    GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error)
}
```

**Why This Matters:** As long as V2 implements this interface identically to V1, the rest of the application doesn't care which implementation is used. This is the power of hexagonal architecture.

### Credential Flow Architecture

**Critical Security Consideration:** Credentials are stored in **OS keyring**, NOT in files or database.

**Flow Diagram:**
```
User Login
    ↓
WorkspaceService.SetConfig()
    ↓
CredentialStore.Store(workspaceID, credentials)
    ↓
OS Keyring (macOS Keychain / Linux Secret Service / Windows Credential Manager)
    ↓
[Encrypted at OS level]

Later...

initJiraAdapter() called
    ↓
WorkspaceService.GetConfig(workspaceID)
    ↓
CredentialStore.Retrieve(workspaceID)
    ↓
OS Keyring returns credentials
    ↓
NewJiraAdapterV2FromConfig(config) ← config contains Username + APIToken
    ↓
jira.BasicAuthTransport{Username: config.Username, Password: config.APIToken}
```

**Key Insight (From Gemini):**
> "The adapter is 'dumb.' It doesn't know WHERE the credentials came from (keyring, env var, file). It is simply configured with them. This upholds the principle of dependency inversion and keeps your security-sensitive logic centralized in the WorkspaceService."

**V2 Credential Handling:**
```go
// From jira_adapter_v2.go
func NewJiraAdapterV2FromConfig(config *domain.WorkspaceConfig, ...) (ports.JiraPort, error) {
    // V2 receives credentials as simple struct values
    tp := jira.BasicAuthTransport{
        Username: config.Username,  // From keyring (via WorkspaceService)
        Password: config.APIToken,  // From keyring (via WorkspaceService)
    }

    client, err := jira.NewClient(tp.Client(), config.JiraURL)
    // Adapter doesn't know credentials came from keyring
}
```

**Security Guarantees:**
- ✅ No credentials in database (only workspace ID reference)
- ✅ No credentials in logs (automatic redaction)
- ✅ OS-level encryption
- ✅ Per-user isolation
- ✅ Adapter layer doesn't handle credential storage

**Gemini Verdict:** "This is a model implementation."

### Feature Flag System

**How It Works:**

**Environment Variable:**
```bash
export TICKETR_JIRA_ADAPTER_VERSION=v2  # Use V2 (default)
export TICKETR_JIRA_ADAPTER_VERSION=v1  # Use V1 (rollback)
```

**Factory Logic (factory.go):**
```go
func NewJiraAdapterFromConfigWithVersion(config *domain.WorkspaceConfig, ...) (ports.JiraPort, error) {
    version := os.Getenv("TICKETR_JIRA_ADAPTER_VERSION")
    if version == "" {
        version = "v2" // Default to V2
    }

    switch version {
    case "v1":
        return NewJiraAdapterFromConfig(config, fieldMappings)  // Old implementation
    case "v2":
        return NewJiraAdapterV2FromConfig(config, fieldMappings) // New library
    default:
        return nil, fmt.Errorf("unknown adapter version: %s", version)
    }
}
```

**Integration Point (main.go:142):**
```go
// BEFORE (V1 only):
return jira.NewJiraAdapterFromConfig(config, fieldMappings)

// AFTER (V1 or V2 based on env var):
return jira.NewJiraAdapterFromConfigWithVersion(config, fieldMappings)
```

**Rollback Strategy:**
```bash
# If V2 has issues in production
export TICKETR_JIRA_ADAPTER_VERSION=v1
ticketr pull  # Instant rollback to V1

# No code changes needed, no recompilation required
```

**Deprecation Timeline:**
- **v3.1.1 (current):** V2 ships as default, V1 available via flag
- **v3.2.0 or v3.3.0:** Remove V1 code, remove feature flag, V2 only

### External AI Architect Verdicts

#### Gemini's Assessment

**Overall:** "Architectural Blessing"

**Quote:**
> "This is an **architecturally sound and highly recommended decision.** Your team's adoption of Hexagonal Architecture is paying dividends precisely as intended: you can swap a significant, complex infrastructure component (the Jira client) with minimal-to-no impact on your core business logic."

**Key Recommendations:**
1. ✅ Architecture correct - library isolated to adapter layer
2. ✅ Credential handling is "model implementation"
3. ❌ Do NOT fork library proactively (creates unnecessary burden)
4. ✅ Trust library, verify with integration tests
5. ✅ Keep V1 for 1-2 releases, then MUST remove dead code
6. ❌ ADR creation is "non-negotiable"

**Production Checklist (8 items):**
1. Wire factory into main.go
2. Add CI testing for both V1 and V2
3. Create benchmark tests
4. Create ADR
5. Update Architecture.md
6. Merge and deploy with V2 default
7. Monitor logs for V2 errors
8. Deprecate V1 in 1-2 releases

#### Codex's Assessment

**Overall:** "Pragmatic, battle-tested choice"

**Quote:**
> "go-jira remains a pragmatic, battle-tested client for Jira Cloud in 2025 if you value stability over cutting-edge coverage and are prepared to patch edge cases locally."

**Key Recommendations:**
1. ✅ Mature, widely used, workable for Jira Cloud REST v2
2. ⚠️ Expect to fork/vendor for quick fixes if needed
3. ✅ Light wrappers/patches locally for edge cases acceptable
4. ✅ Given maintenance constraint, library minimizes long-term burden
5. ⚠️ Sporadic maintenance cadence (new endpoints may lag)
6. ✅ Hybrid pattern: library for transport, domain layer for business logic

**Alternatives Comparison:**
- `andygrunwald/go-jira`: ✅ Best fit (battle-tested, right scope)
- `ctreminiom/go-atlassian`: Broader coverage but younger, less hardened
- No official Atlassian Go SDK exists

**Consensus:** Both AIs unanimously approved the library integration as the correct architectural decision with manageable risks.

### Risk Register

**From External AI Analysis:**

| Risk | Likelihood | Impact | Mitigation | Owner |
|------|-----------|--------|------------|-------|
| Library abandonment | Medium | Medium | Hexagonal architecture allows adapter swap; contingency fork plan | Builder |
| Behavioral divergence V1→V2 | Medium | High | Integration tests validate parity; feature flag rollback | Verifier |
| Performance regression | Low | Medium | Benchmarks validate acceptable performance | Verifier |
| Integration issues | Low | High | CI tests both adapters; manual UAT before V1 removal | Director |
| Missing documentation | HIGH (current) | Medium | ADR creation mandated; Scribe delegated | Scribe |

**Gemini's Risk Philosophy:**
> "The Adapter is Your Shield: If the library becomes untenable, you only need to write a JiraAdapterV3, not rewrite your application."

---

## Success Criteria

### Phase 7 Completion Checklist

**Code Integration:**
- [ ] Factory wired into `cmd/ticketr/main.go:142`
- [ ] Version logging operational (`[jira-v1]` / `[jira-v2]` tags)
- [ ] `go build ./cmd/ticketr` succeeds
- [ ] `go test ./...` passes (all 147+ tests)

**Testing:**
- [ ] Feature flag script created: `scripts/test-adapter-versions.sh`
- [ ] CI validates both V1 and V2 adapters
- [ ] Integration tests created: `internal/adapters/jira/integration_test.go`
- [ ] Benchmarks created: `internal/adapters/jira/adapter_bench_test.go`
- [ ] V2 performance within 20% of V1 (acceptable)

**Governance:**
- [ ] ADR created: `docs/adr/001-adopt-go-jira-library.md`
- [ ] Architecture.md updated with library reference
- [ ] Jira adapter README updated with feature flag docs
- [ ] Deployment plan created: `docs/deployment/JIRA-LIBRARY-ROLLOUT.md`

**Verification:**
- [ ] Gemini's 8-item checklist complete
- [ ] All agent handovers received and validated
- [ ] Completion report created: `.agents/phase7/PHASE7-COMPLETION-REPORT.md`
- [ ] User handover report prepared

**Production Readiness:**
- [ ] Feature flag defaults to V2
- [ ] Rollback procedure documented and tested
- [ ] Monitoring strategy defined (version-tagged logs)
- [ ] V1 deprecation timeline documented (v3.2.0 or v3.3.0)

### Quality Gates

**Must Pass Before User Handover:**

1. **Build Quality:**
   ```bash
   go build ./cmd/ticketr          # Must succeed
   go test ./...                   # Must pass 100%
   go vet ./...                    # No warnings
   ```

2. **Feature Flag Validation:**
   ```bash
   TICKETR_JIRA_ADAPTER_VERSION=v1 go test ./internal/adapters/jira/...  # Pass
   TICKETR_JIRA_ADAPTER_VERSION=v2 go test ./internal/adapters/jira/...  # Pass
   TICKETR_JIRA_ADAPTER_VERSION=invalid ...                              # Error correctly
   ```

3. **Documentation Complete:**
   - ADR exists and follows template
   - Architecture.md mentions library
   - Deployment plan exists

4. **External AI Approval:**
   - ✅ Gemini approved (already received)
   - ✅ Codex approved (already received)

---

## Troubleshooting

### Common Issues

#### Issue 1: Factory Not Found

**Symptom:**
```
undefined: jira.NewJiraAdapterFromConfigWithVersion
```

**Cause:** Factory function doesn't exist or not exported

**Fix:**
```bash
# Verify factory exists
grep "func NewJiraAdapterFromConfigWithVersion" internal/adapters/jira/factory.go

# If missing, Builder agent didn't create it
# Re-read: "What Has Already Been Done" → V2 Implementation
```

#### Issue 2: Tests Fail After Integration

**Symptom:**
```
--- FAIL: TestJiraAdapter (0.00s)
```

**Cause:** Factory changes broke existing tests

**Debug:**
```bash
# Run with verbose output
go test -v ./internal/adapters/jira/...

# Check environment variable
echo $TICKETR_JIRA_ADAPTER_VERSION

# Try forcing V1
TICKETR_JIRA_ADAPTER_VERSION=v1 go test ./internal/adapters/jira/...
```

**Fix:** Ensure tests don't set `TICKETR_JIRA_ADAPTER_VERSION` or handle both versions

#### Issue 3: Credentials Not Working with V2

**Symptom:**
```
[jira-v2] failed to authenticate: 401 Unauthorized
```

**Cause:** V2 adapter not receiving credentials properly from WorkspaceConfig

**Debug:**
```bash
# Check workspace config
ticketr workspace current

# Verify credentials in keyring
# (This retrieves from keyring and tests connection)
ticketr workspace validate
```

**Fix:** Ensure `NewJiraAdapterV2FromConfig` receives correct `WorkspaceConfig` struct

#### Issue 4: CI Fails on Feature Flag Test

**Symptom:**
```
CI: TICKETR_JIRA_ADAPTER_VERSION=v1 go test ... FAILED
```

**Cause:** V1 tests broken or V1 adapter removed prematurely

**Fix:**
- If V1 tests actually broken: Fix V1 adapter
- If V1 should be deprecated: Update CI to only test V2
- During transition: BOTH must pass

#### Issue 5: Benchmark Shows Performance Regression

**Symptom:**
```
BenchmarkSearchTickets_V2    50% slower than V1
```

**Cause:** Library overhead or inefficient V2 implementation

**Analysis:**
```bash
# Run detailed benchmark
go test -bench=BenchmarkSearchTickets -benchmem -cpuprofile=cpu.out ./internal/adapters/jira/

# Analyze profile
go tool pprof cpu.out
```

**Decision:**
- If <20% slower: Acceptable (Gemini's threshold)
- If >20% slower: Investigate library usage, consider optimizations
- If >50% slower: Escalate to user, may need hybrid approach

### Agent Handover Failures

#### Scribe Fails to Create ADR

**Symptom:** Scribe reports completion but ADR file missing

**Verify:**
```bash
ls docs/adr/001-adopt-go-jira-library.md
```

**Fix:**
- Re-invoke Scribe with specific instruction: "Create docs/adr/001-adopt-go-jira-library.md using template in ORCHESTRATION-PLAN.md Step 1.1"
- Verify Scribe has write access to docs/adr/ directory
- Check Scribe handover report for error messages

#### Builder Reports Success But Factory Not Wired

**Symptom:** Builder says complete but main.go:142 unchanged

**Verify:**
```bash
grep "NewJiraAdapterFromConfigWithVersion" cmd/ticketr/main.go
```

**Fix:**
- Check git diff to see what Builder actually changed
- Re-invoke Builder with explicit instruction: "Update cmd/ticketr/main.go line 142 to call jira.NewJiraAdapterFromConfigWithVersion instead of jira.NewJiraAdapterFromConfig"
- Verify Builder handover report shows file modification

#### Verifier Reports Tests Pass But CI Fails

**Symptom:** Verifier says all tests pass locally, but CI fails

**Causes:**
- Environment differences (local vs CI)
- Missing CI configuration updates
- Cached dependencies

**Fix:**
```bash
# Reproduce CI environment locally
go clean -testcache
go test ./...

# Check CI logs for specific failure
# Update .github/workflows/test.yml if needed
```

---

## Appendix

### A. File Locations Quick Reference

**Already Exist:**
```
internal/adapters/jira/jira_adapter.go        # V1 implementation (1,136 lines)
internal/adapters/jira/jira_adapter_v2.go     # V2 implementation (757 lines) ✅
internal/adapters/jira/factory.go              # Feature flag system (68 lines) ✅
internal/core/ports/jira_port.go               # JiraPort interface
cmd/ticketr/main.go                            # Application entry point
research/EXTERNAL-VALIDATION-REPORT.md         # 15+ source validation ✅
.agents/phase7/ORCHESTRATION-PLAN.md           # Execution plan ✅
.agents/phase7/BOOT-TO-PHASE-7.md              # This document ✅
```

**Must Be Created:**
```
docs/adr/001-adopt-go-jira-library.md          # ADR (Step 1) ❌
scripts/test-adapter-versions.sh               # Feature flag test (Step 3) ❌
internal/adapters/jira/integration_test.go     # V1/V2 parity tests (Step 3) ❌
internal/adapters/jira/adapter_bench_test.go   # Benchmarks (Step 4) ❌
docs/deployment/JIRA-LIBRARY-ROLLOUT.md        # Deployment plan (Step 5) ❌
.agents/phase7/PHASE7-COMPLETION-REPORT.md     # Final report (Step 5) ❌
```

**Must Be Updated:**
```
cmd/ticketr/main.go                            # Line 142 (Step 2) ❌
docs/ARCHITECTURE.md                           # Jira Adapter section (Step 1) ❌
internal/adapters/jira/README.md               # Feature flag docs (Step 1) ❌
internal/adapters/jira/jira_adapter.go         # Add [jira-v1] tags (Step 2) ❌
internal/adapters/jira/jira_adapter_v2.go      # Add [jira-v2] tags (Step 2) ❌
.github/workflows/test.yml                     # Add V1/V2 testing (Step 3) ❌
```

### B. Command Reference

**Build & Test:**
```bash
# Full build
go build ./cmd/ticketr

# All tests
go test ./...

# Specific package tests
go test ./internal/adapters/jira/...

# Verbose tests
go test -v ./internal/adapters/jira/...

# With coverage
go test -cover ./internal/adapters/jira/...

# Clear test cache
go clean -testcache
```

**Feature Flag Testing:**
```bash
# Test V1
TICKETR_JIRA_ADAPTER_VERSION=v1 go test ./internal/adapters/jira/...

# Test V2
TICKETR_JIRA_ADAPTER_VERSION=v2 go test ./internal/adapters/jira/...

# Test invalid version
TICKETR_JIRA_ADAPTER_VERSION=invalid ./ticketr workspace current
# Should error: "unknown adapter version: invalid"
```

**Benchmarks:**
```bash
# Run benchmarks
go test -bench=. ./internal/adapters/jira/

# With memory stats
go test -bench=. -benchmem ./internal/adapters/jira/

# Specific benchmark
go test -bench=BenchmarkSearchTickets ./internal/adapters/jira/
```

**Code Quality:**
```bash
# Static analysis
go vet ./...

# Format check
go fmt ./...

# Dependency check
go mod tidy
go mod verify

# Security scan
govulncheck ./...
```

**Git Operations:**
```bash
# Check current branch
git branch --show-current

# View changes
git status
git diff

# Commit work
git add .
git commit -m "feat: integrate andygrunwald/go-jira library with feature flag system"

# Push to remote
git push origin feature/jira-domain-redesign
```

### C. Agent Invocation Templates

**Scribe Agent:**
```
You are the Scribe agent. Your mission is to complete governance and documentation for Phase 7 Jira library integration.

Context:
- Read .agents/phase7/BOOT-TO-PHASE-7.md for complete background
- Read .agents/phase7/ORCHESTRATION-PLAN.md Step 1 for detailed tasks

Tasks:
1. Create docs/adr/001-adopt-go-jira-library.md following ADR template in ORCHESTRATION-PLAN
2. Update docs/ARCHITECTURE.md Jira Adapter section to mention library
3. Update internal/adapters/jira/README.md with feature flag documentation

Requirements:
- Follow templates exactly as specified in ORCHESTRATION-PLAN
- Include all sections (Context, Decision, Consequences, Alternatives)
- Reference external validation (Gemini + Codex approvals)

Deliverable:
Provide a completion report listing:
- All files created/updated
- Verification that all 3 documentation tasks complete
- Any issues encountered
```

**Builder Agent:**
```
You are the Builder agent. Your mission is to integrate the V2 Jira adapter factory into main.go and add version logging.

Context:
- Read .agents/phase7/BOOT-TO-PHASE-7.md for complete background
- Read .agents/phase7/ORCHESTRATION-PLAN.md Step 2 for detailed tasks

Tasks:
1. Update cmd/ticketr/main.go line 142:
   OLD: return jira.NewJiraAdapterFromConfig(config, fieldMappings)
   NEW: return jira.NewJiraAdapterFromConfigWithVersion(config, fieldMappings)

2. Add version logging to all errors in jira_adapter.go:
   - Wrap all errors with fmt.Errorf("[jira-v1] ... %w", err)

3. Add version logging to all errors in jira_adapter_v2.go:
   - Wrap all errors with fmt.Errorf("[jira-v2] ... %w", err)

4. Verify compilation:
   - Run: go build ./cmd/ticketr
   - Run: go test ./internal/adapters/jira/...

Requirements:
- Do NOT modify any business logic, only error wrapping
- Ensure all tests still pass after changes
- Version tags must be consistent: [jira-v1] or [jira-v2]

Deliverable:
Provide a completion report listing:
- All files modified
- Test results (must pass)
- Build verification (must succeed)
- Any issues encountered
```

**Verifier Agent:**
```
You are the Verifier agent. Your mission is to validate the feature flag system and benchmark performance.

Context:
- Read .agents/phase7/BOOT-TO-PHASE-7.md for complete background
- Read .agents/phase7/ORCHESTRATION-PLAN.md Steps 3-4 for detailed tasks

Tasks:
1. Create scripts/test-adapter-versions.sh following template in ORCHESTRATION-PLAN
2. Update .github/workflows/test.yml to test both V1 and V2
3. Create internal/adapters/jira/integration_test.go with V1/V2 parity tests
4. Create internal/adapters/jira/adapter_bench_test.go with benchmarks
5. Run benchmarks and verify V2 performance within 20% of V1

Requirements:
- All tests must pass for BOTH V1 and V2
- Invalid version must be rejected with error
- Benchmarks must show acceptable performance (<20% regression)
- Integration tests must validate behavioral parity

Deliverable:
Provide a completion report with:
- All files created
- Test results (V1 pass, V2 pass, invalid version rejected)
- Benchmark results (V1 vs V2 comparison)
- Performance assessment (acceptable/unacceptable)
- Any issues encountered
```

### D. External AI Consultation Evidence

**Gemini Consultation:**
- **Date:** 2025-10-21
- **Method:** `cat /tmp/jira-integration-consultation-context.md | gemini -p "..."`
- **Result:** "Architectural Blessing"
- **Output:** 2,474 tokens
- **Key Quote:** "This is an architecturally sound and highly recommended decision."

**Codex Consultation:**
- **Date:** 2025-10-21
- **Method:** `cat /tmp/jira-validation-prompt.txt | codex exec --full-auto -`
- **Result:** "Pragmatic, battle-tested choice"
- **Output:** Detailed technical assessment
- **Key Quote:** "go-jira remains a pragmatic, battle-tested client for Jira Cloud in 2025"

**Consensus:** UNANIMOUS APPROVAL ✅

### E. Glossary

**ADR:** Architecture Decision Record - formal documentation of architectural choices

**Adapter:** Implementation of a port interface (infrastructure layer)

**Feature Flag:** Environment variable controlling which implementation to use

**Gemini/Codex:** External AI architects consulted for validation

**Hexagonal Architecture:** Ports & Adapters pattern, business logic isolated from infrastructure

**JiraPort:** Interface defining contract for Jira integration operations

**Phase 7:** Current phase - Jira library integration

**Ports:** Interfaces defining boundaries between core and infrastructure

**Scribe/Builder/Verifier/Director:** Specialized agents in orchestration methodology

**V1:** Original custom HTTP client implementation (1,136 lines)

**V2:** New library-based implementation using andygrunwald/go-jira (757 lines)

**WorkspaceConfig:** Domain struct containing Jira credentials and configuration

---

## Document Metadata

**Created:** 2025-10-21
**Author:** Director Agent (Phase 7 orchestration)
**Version:** 1.0
**Purpose:** Zero-context onboarding for fresh Director agent
**Audience:** New Director agent OR human taking over Phase 7 execution
**Dependencies:**
- `.agents/phase7/ORCHESTRATION-PLAN.md` (execution recipe)
- `research/EXTERNAL-VALIDATION-REPORT.md` (library validation)
- `internal/adapters/jira/jira_adapter_v2.go` (V2 implementation)

**Next Steps:** Execute 5-step orchestration plan in ORCHESTRATION-PLAN.md

---

**End of Boot-to-Phase-7 Document**

🚀 **You are now fully contextualized. Execute the ORCHESTRATION-PLAN.md to complete Phase 7.**
