# Phase 7: Jira Library Integration - Orchestration Plan

**Date:** 2025-10-21
**Branch:** `feature/jira-domain-redesign`
**Director:** Claude (Main Agent)
**Status:** ✅ EXTERNAL VALIDATION COMPLETE → READY FOR IMPLEMENTATION

---

## Executive Summary

**Mission:** Complete production-ready integration of `andygrunwald/go-jira` v1.17.0 library into Ticketr's Jira adapter following hexagonal architecture principles.

**External Validation:** ✅ APPROVED by both Gemini and Codex AI architects
**Implementation Status:** Code complete but NOT integrated
**Production Blockers:** 5 critical gaps identified by external AIs

---

## Orchestration Strategy

This plan follows **strict Director methodology** with specialized agent delegation:

- **Scribe:** All documentation (ADR, Architecture.md, README updates)
- **Builder:** All code integration (main.go wiring, logging, tests)
- **Verifier:** All testing validation (CI, benchmarks, feature flag verification)
- **Director (me):** Orchestration, handover, final approval

**Parallelization:** Steps 1-2 can run in parallel, Steps 3-5 sequential.

---

## Implementation Recipe

### **STEP 1: Governance & Documentation (Scribe)**

**Agent:** Scribe
**Priority:** CRITICAL (Gemini: "non-negotiable")
**Duration:** ~30 minutes
**Can Run in Parallel:** Yes (with Step 2)

#### Task 1.1: Create Architecture Decision Record (ADR)

**File:** `docs/adr/001-adopt-go-jira-library.md`

**Required Sections (per Gemini):**
```markdown
# ADR-001: Adopt andygrunwald/go-jira Library for Jira Integration

## Status
APPROVED

## Context
- Custom HTTP client (jira_adapter.go) is 1,136 lines
- High maintenance burden (manual JSON, pagination, error handling)
- Need reliable, battle-tested Jira REST API client
- Single user, no migration concerns

## Decision
Adopt `github.com/andygrunwald/go-jira` v1.17.0 as underlying client for Jira adapter

## Rationale
- **Maturity:** 9 years production use (2015-2025)
- **Adoption:** 868 importers, 1,600 stars
- **Security:** 0 CVEs, clean govulncheck
- **Code Reduction:** 1,136 → 757 lines (33% reduction)
- **External Validation:** Approved by Gemini + Codex AI architects
- **Architecture Fit:** Respects ports & adapters, isolated to adapter layer

## Consequences

### Positive
- Reduced code complexity (-379 lines)
- Lower maintenance burden
- Improved reliability (battle-tested)
- Faster feature development
- Better error handling

### Negative
- New third-party dependency (+12 transitive deps)
- Library maintenance risk (1 commit in 2024)
- Potential behavioral divergence from V1

### Mitigation
- Feature flag system (TICKETR_JIRA_ADAPTER_VERSION)
- Hexagonal architecture allows easy adapter swap
- Integration tests validate both V1/V2 behavior
- Contingency fork plan documented

## Alternatives Considered
- **ctreminiom/go-atlassian:** Overkill (Confluence/JSM support unnecessary)
- **Keep custom HTTP client:** Technical debt, maintenance burden grows
- **Hybrid approach:** Inconsistent, difficult to debug

## References
- External validation: research/EXTERNAL-VALIDATION-REPORT.md
- Gemini consultation: /tmp/gemini-consultation-results.txt
- Codex consultation: /tmp/codex-consultation-results.txt
```

#### Task 1.2: Update docs/ARCHITECTURE.md

**Changes Required:**
1. Update "Jira Adapter" section to mention library:
   ```markdown
   ### Jira Adapter (internal/adapters/jira/)

   Implements JiraPort interface using `andygrunwald/go-jira` v1.17.0 library.

   **Implementation Versions:**
   - V1: Custom HTTP client (deprecated, removal planned v3.3.0)
   - V2: Library-based implementation (default)

   **Feature Flag:** `TICKETR_JIRA_ADAPTER_VERSION` (values: v1, v2)
   ```

2. Update dependencies section:
   ```markdown
   ### External Dependencies
   - `github.com/andygrunwald/go-jira` v1.17.0 - Jira REST API client
   - `github.com/zalando/go-keyring` v0.2.6 - OS keyring credential storage
   ```

3. Update credential flow diagram to show V2 integration

#### Task 1.3: Update internal/adapters/jira/README.md

**Add section:**
```markdown
## Library Integration

**Current:** V2 adapter using `andygrunwald/go-jira` v1.17.0

**Feature Flag:**
```bash
# Use V2 (default)
export TICKETR_JIRA_ADAPTER_VERSION=v2
ticketr pull

# Rollback to V1
export TICKETR_JIRA_ADAPTER_VERSION=v1
ticketr pull
```

**ADR:** See docs/adr/001-adopt-go-jira-library.md
```

**Deliverable:** 3 documentation files updated, ADR created

---

### **STEP 2: Code Integration (Builder)**

**Agent:** Builder
**Priority:** CRITICAL
**Duration:** ~20 minutes
**Can Run in Parallel:** Yes (with Step 1)

#### Task 2.1: Wire Factory into cmd/ticketr/main.go

**File:** `cmd/ticketr/main.go`

**Change Location:** Line 142

**Current Code:**
```go
return jira.NewJiraAdapterFromConfig(config, fieldMappings)
```

**New Code:**
```go
return jira.NewJiraAdapterFromConfigWithVersion(config, fieldMappings)
```

**Verification:**
- Factory function exists in `internal/adapters/jira/factory.go:68`
- Feature flag defaults to `v2`
- Both V1 and V2 constructors still compile

#### Task 2.2: Add Enhanced Logging with Version Tags

**Files:**
- `internal/adapters/jira/jira_adapter.go` (V1)
- `internal/adapters/jira/jira_adapter_v2.go` (V2)

**Changes Required (per Gemini):**

**In V1 Adapter:**
```go
// Wrap errors with adapter version tag
if err != nil {
    return nil, fmt.Errorf("[jira-v1] failed to search tickets: %w", err)
}
```

**In V2 Adapter:**
```go
// Wrap errors with adapter version tag
if err != nil {
    return nil, fmt.Errorf("[jira-v2] failed to search tickets: %w", err)
}
```

**Locations:**
- `SearchTickets()` - Add `[jira-v1]` / `[jira-v2]` prefix to all errors
- `CreateTicket()` - Add version prefix
- `UpdateTicket()` - Add version prefix
- `CreateTask()` - Add version prefix
- `UpdateTask()` - Add version prefix

**Purpose:** Monitor for increased errors post-deployment, identify adapter version issues

#### Task 2.3: Verify Compilation

**Commands:**
```bash
go build ./cmd/ticketr
go test ./internal/adapters/jira/...
```

**Expected:**
- ✅ Build succeeds
- ✅ All existing tests pass
- ✅ No new warnings

**Deliverable:** Factory wired, version logging added, builds successfully

---

### **STEP 3: Feature Flag Validation (Verifier)**

**Agent:** Verifier
**Priority:** HIGH (Gemini: "feature flag itself is untested")
**Duration:** ~30 minutes
**Depends On:** Step 2 complete

#### Task 3.1: Create Feature Flag Test Script

**File:** `scripts/test-adapter-versions.sh`

```bash
#!/bin/bash
set -e

echo "Testing Jira Adapter Feature Flag System"
echo "========================================="

# Test V1
echo ""
echo "Testing V1 Adapter..."
export TICKETR_JIRA_ADAPTER_VERSION=v1
go test ./internal/adapters/jira/... -v -run TestJiraAdapter
echo "✅ V1 tests passed"

# Test V2
echo ""
echo "Testing V2 Adapter..."
export TICKETR_JIRA_ADAPTER_VERSION=v2
go test ./internal/adapters/jira/... -v -run TestJiraAdapterV2
echo "✅ V2 tests passed"

# Test invalid version
echo ""
echo "Testing invalid version handling..."
export TICKETR_JIRA_ADAPTER_VERSION=invalid
if go run ./cmd/ticketr workspace current 2>&1 | grep -q "unknown adapter version"; then
    echo "✅ Invalid version correctly rejected"
else
    echo "❌ Invalid version not properly handled"
    exit 1
fi

echo ""
echo "✅ All feature flag tests passed"
```

**Make executable:**
```bash
chmod +x scripts/test-adapter-versions.sh
```

#### Task 3.2: Add CI Validation

**File:** `.github/workflows/test.yml` (or equivalent CI config)

**Add step:**
```yaml
- name: Test Jira Adapter Versions
  run: |
    # Test V1
    TICKETR_JIRA_ADAPTER_VERSION=v1 go test ./internal/adapters/jira/...

    # Test V2
    TICKETR_JIRA_ADAPTER_VERSION=v2 go test ./internal/adapters/jira/...
```

#### Task 3.3: Create Integration Test Suite

**File:** `internal/adapters/jira/integration_test.go`

**Purpose (per Gemini):** "Ensure consistent behavior across V1/V2 for edge cases"

```go
//go:build integration

package jira

import (
    "testing"
    "os"
)

// TestAdapterBehaviorParity validates V1 and V2 handle same scenarios identically
func TestAdapterBehaviorParity(t *testing.T) {
    tests := []struct {
        name     string
        testFunc func(t *testing.T, adapter JiraPort)
    }{
        {"invalid JQL syntax", testInvalidJQL},
        {"ticket not found", testTicketNotFound},
        {"permission error", testPermissionError},
        {"large result set pagination", testLargePagination},
    }

    for _, tt := range tests {
        t.Run(tt.name+" (V1)", func(t *testing.T) {
            os.Setenv("TICKETR_JIRA_ADAPTER_VERSION", "v1")
            adapter := createTestAdapter(t)
            tt.testFunc(t, adapter)
        })

        t.Run(tt.name+" (V2)", func(t *testing.T) {
            os.Setenv("TICKETR_JIRA_ADAPTER_VERSION", "v2")
            adapter := createTestAdapter(t)
            tt.testFunc(t, adapter)
        })
    }
}
```

**Deliverable:** Feature flag validated in CI, integration tests created

---

### **STEP 4: Performance Validation (Verifier)**

**Agent:** Verifier
**Priority:** RECOMMENDED (Gemini: "optional but recommended")
**Duration:** ~20 minutes
**Depends On:** Step 2 complete

#### Task 4.1: Create Benchmark Test

**File:** `internal/adapters/jira/adapter_bench_test.go`

```go
package jira

import (
    "context"
    "os"
    "testing"
)

// BenchmarkSearchTickets_V1 benchmarks V1 adapter search performance
func BenchmarkSearchTickets_V1(b *testing.B) {
    os.Setenv("TICKETR_JIRA_ADAPTER_VERSION", "v1")
    adapter := createBenchAdapter(b)
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := adapter.SearchTickets(ctx, "PROJ", "project = PROJ", nil)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// BenchmarkSearchTickets_V2 benchmarks V2 adapter search performance
func BenchmarkSearchTickets_V2(b *testing.B) {
    os.Setenv("TICKETR_JIRA_ADAPTER_VERSION", "v2")
    adapter := createBenchAdapter(b)
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := adapter.SearchTickets(ctx, "PROJ", "project = PROJ", nil)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

#### Task 4.2: Run Benchmark Comparison

**Commands:**
```bash
# Run benchmarks
go test -bench=BenchmarkSearchTickets -benchmem ./internal/adapters/jira/

# Expected output format:
# BenchmarkSearchTickets_V1-8    100    12.5 ms/op    4096 B/op    50 allocs/op
# BenchmarkSearchTickets_V2-8    100    13.2 ms/op    5120 B/op    55 allocs/op
```

**Acceptance Criteria:**
- V2 performance within 20% of V1
- No memory leaks
- Reasonable allocation counts

**Deliverable:** Benchmark tests created, performance validated

---

### **STEP 5: Final Verification & Deployment Prep (Director + Verifier)**

**Agent:** Director (me) + Verifier
**Priority:** CRITICAL
**Duration:** ~15 minutes
**Depends On:** Steps 1-4 complete

#### Task 5.1: Final Checklist Validation

**Gemini's Production Checklist:**
- [ ] Code: Factory wired into main.go
- [ ] Test: CI validates both v1 and v2
- [ ] Benchmark: Performance comparison complete
- [ ] Governance: ADR created
- [ ] Documentation: Architecture.md updated
- [ ] Logging: Version tags added to errors
- [ ] Integration Tests: V1/V2 parity validated

#### Task 5.2: Create Deployment Plan

**File:** `docs/deployment/JIRA-LIBRARY-ROLLOUT.md`

```markdown
# Jira Library Integration Rollout Plan

## Deployment Strategy

**Feature Flag:** `TICKETR_JIRA_ADAPTER_VERSION`
**Default:** `v2` (library-based)
**Rollback:** Set to `v1` (instant fallback)

## Phase 1: Initial Deployment (Day 1)
- Deploy with V2 as default
- Monitor logs for `[jira-v2]` errors
- Watch for performance degradation

## Phase 2: Monitoring (Week 1-2)
- Compare error rates V1 vs V2
- Validate custom field handling
- Check pagination performance

## Phase 3: V1 Deprecation (v3.2.0 or v3.3.0)
- Remove jira_adapter.go (V1)
- Remove feature flag factory
- Remove `TICKETR_JIRA_ADAPTER_VERSION` support
- Update docs to remove V1 references

## Rollback Procedure

If critical issues found:
```bash
export TICKETR_JIRA_ADAPTER_VERSION=v1
ticketr pull  # Instant rollback
```

## Monitoring Queries

```bash
# Check for V2 errors
grep "\[jira-v2\]" logs/ | wc -l

# Compare error rates
grep "\[jira-v1\]" logs/ | wc -l
grep "\[jira-v2\]" logs/ | wc -l
```
```

#### Task 5.3: Create Migration Complete Report

**File:** `.agents/phase7/PHASE7-COMPLETION-REPORT.md`

Template for final report documenting:
- External validation results
- Implementation changes
- Test coverage
- Performance benchmarks
- Deployment plan
- Success metrics

**Deliverable:** Final checklist verified, deployment plan created

---

## Orchestration Execution Plan

### **Parallel Execution (Steps 1-2)**

**Invoke simultaneously:**
1. **Scribe Agent:** Documentation + ADR creation (Step 1)
2. **Builder Agent:** Code integration + logging (Step 2)

**Director monitors:** Both agents complete, handover reports received

### **Sequential Execution (Steps 3-5)**

**After Steps 1-2 complete:**
3. **Verifier Agent:** Feature flag validation (Step 3)
4. **Verifier Agent:** Performance benchmarks (Step 4)
5. **Director:** Final verification + deployment prep (Step 5)

---

## Success Metrics

### **Code Quality**
- ✅ All tests pass (147 existing + new integration tests)
- ✅ Build succeeds with factory wired
- ✅ No new linter warnings
- ✅ Coverage maintained ≥74.8%

### **Governance**
- ✅ ADR created and approved
- ✅ Architecture.md updated
- ✅ Deployment plan documented

### **Testing**
- ✅ Feature flag validated in CI
- ✅ Both V1 and V2 tests pass
- ✅ Integration tests cover edge cases
- ✅ Benchmarks show acceptable performance

### **Deployment Readiness**
- ✅ Feature flag defaults to V2
- ✅ Rollback procedure documented
- ✅ Monitoring strategy defined
- ✅ Version logging operational

---

## Risk Mitigation

### **Risk 1: V2 Behavioral Divergence**
- **Mitigation:** Integration tests validate V1/V2 parity
- **Fallback:** Feature flag instant rollback to V1

### **Risk 2: Performance Regression**
- **Mitigation:** Benchmarks validate acceptable performance
- **Acceptance:** V2 within 20% of V1 performance

### **Risk 3: Library Abandonment**
- **Mitigation:** Hexagonal architecture allows adapter swap
- **Contingency:** Fork plan documented (if critical bug unfixed in X days)

### **Risk 4: Integration Issues**
- **Mitigation:** All tests run against both adapters in CI
- **Verification:** Manual UAT before removing V1

---

## Timeline Estimate

| Step | Agent | Duration | Dependencies |
|------|-------|----------|--------------|
| 1. Documentation | Scribe | 30 min | None |
| 2. Integration | Builder | 20 min | None |
| 3. Feature Flag Tests | Verifier | 30 min | Step 2 |
| 4. Benchmarks | Verifier | 20 min | Step 2 |
| 5. Final Verification | Director | 15 min | Steps 1-4 |

**Total:** ~75 minutes (with parallelization)
**Serial:** ~115 minutes (without parallelization)

---

## Handover Requirements

Each agent must provide:
1. **Completion Report:** Files modified, tests added, verification results
2. **Test Evidence:** All tests passing, benchmark results
3. **Integration Verification:** Feature works end-to-end
4. **Documentation Updates:** All docs current and accurate

**Final Approval:** Director validates all handovers before marking phase complete.

---

## Appendix: External AI Recommendations

### Gemini's Critical Advice
> "You have done excellent preparatory work. Completing these final steps will ensure this change is not just a technical success, but a well-governed and sustainable improvement to your application."

### Codex's Strategic Guidance
> "Given your maintenance constraint and existing 1,136-line custom client, starting from a library plus light wrappers still minimizes long-term burden."

### Consensus Verdict
**Both AIs unanimously approved** the library integration as architecturally sound, with proper mitigation strategies for identified risks.

---

**End of Orchestration Plan**
