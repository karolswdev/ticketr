# Integration Testing Guide - Milestone 7 Field Inheritance

## Purpose

This guide prepares you to test Milestone 7 (Field Inheritance Compliance - PROD-009/202) against a real JIRA instance. After completing unit testing with 67 passing tests, we now need to validate that field inheritance works correctly with actual JIRA API interactions.

---

## Prerequisites

### 1. JIRA Instance Access

You need:
- **JIRA URL**: Your Atlassian instance (e.g., `https://yourcompany.atlassian.net`)
- **JIRA Email**: Your account email
- **JIRA API Token**: Generated from Account Settings → Security → API tokens
- **JIRA Project Key**: A test project where you can create tickets and subtasks

### 2. Environment Configuration

**File:** `.env` (create from `.env.example` if needed)

```bash
# Copy example file
cp .env.example .env

# Edit with your actual values
JIRA_URL=https://yourcompany.atlassian.net
JIRA_EMAIL=your.email@company.com
JIRA_API_KEY=your-api-token-here
JIRA_PROJECT_KEY=TEST
```

**Security Note:** Never commit `.env` to version control. It's already in `.gitignore`.

### 3. JIRA Project Setup

Your test project should have:
- **Story/Epic issue type** configured
- **Subtask issue type** enabled
- **Custom fields** available for testing (optional but recommended):
  - Priority (usually built-in)
  - Sprint (if using Scrum)
  - Story Points (if using Scrum)
  - Team/Component (optional)
  - Labels (optional)

**Recommended:** Use a dedicated test project or sandbox to avoid polluting production data.

---

## Test Scenarios

### Scenario 1: Complete Field Inheritance (TC-INT-701.1)

**Goal:** Verify tasks inherit all parent custom fields when task has no custom fields.

**Test Markdown:** Create `test-field-inheritance-1.md`

```markdown
# TICKET: E-Commerce Platform Upgrade
**Status:** To Do
**Project:** TEST
**Description:** Modernize the legacy e-commerce platform to improve performance and user experience.

**Custom Fields:**
- Priority: High
- Sprint: Sprint 24
- Story Points: 13
- Team: Platform Engineering

## TASK: Database Migration Planning
**Status:** To Do
**Description:** Analyze current database schema and plan migration strategy.

**Acceptance Criteria:**
- [ ] Schema analysis complete
- [ ] Migration risks identified
- [ ] Rollback plan documented

## TASK: API Modernization
**Status:** To Do
**Description:** Refactor legacy REST endpoints to use modern API standards.

**Acceptance Criteria:**
- [ ] API endpoints documented
- [ ] OpenAPI specification created
```

**Test Steps:**

1. **Push to JIRA:**
   ```bash
   ticketr push test-field-inheritance-1.md
   ```

2. **Verify in JIRA UI:**
   - Open the created Epic in JIRA
   - Check parent ticket has:
     - Priority: High
     - Sprint: Sprint 24
     - Story Points: 13
     - Team: Platform Engineering
   - Open each subtask and verify they **inherited all parent fields**:
     - Priority: High (inherited)
     - Sprint: Sprint 24 (inherited)
     - Story Points: 13 (inherited)
     - Team: Platform Engineering (inherited)

3. **Expected Behavior:**
   - Both subtasks should show **identical custom field values** as parent
   - No manual field copying required

**Success Criteria:** ✅ All subtask fields match parent fields exactly

---

### Scenario 2: Partial Field Override (TC-INT-701.2)

**Goal:** Verify task-specific custom fields override parent values while inheriting non-overridden fields.

**Test Markdown:** Create `test-field-inheritance-2.md`

```markdown
# TICKET: Authentication System Hardening
**Status:** In Progress
**Project:** TEST
**Description:** Strengthen authentication mechanisms to meet SOC2 compliance requirements.

**Custom Fields:**
- Priority: High
- Sprint: Sprint 24
- Story Points: 8
- Team: Security Team

## TASK: Implement MFA
**Status:** In Progress
**Description:** Add multi-factor authentication support for all user accounts.

**Custom Fields:**
- Priority: Critical
- Story Points: 5

**Acceptance Criteria:**
- [ ] TOTP support implemented
- [ ] SMS fallback configured
- [ ] User enrollment flow complete

## TASK: Audit Log Enhancement
**Status:** To Do
**Description:** Enhance audit logging to capture authentication events.

**Custom Fields:**
- Story Points: 3

**Acceptance Criteria:**
- [ ] Login attempts logged
- [ ] Failed authentication tracked
- [ ] Compliance report generated
```

**Test Steps:**

1. **Push to JIRA:**
   ```bash
   ticketr push test-field-inheritance-2.md
   ```

2. **Verify in JIRA UI:**
   - **Parent Ticket:**
     - Priority: High
     - Sprint: Sprint 24
     - Story Points: 8
     - Team: Security Team

   - **Subtask 1 (Implement MFA):**
     - Priority: **Critical** (OVERRIDDEN from High)
     - Sprint: Sprint 24 (INHERITED)
     - Story Points: **5** (OVERRIDDEN from 8)
     - Team: Security Team (INHERITED)

   - **Subtask 2 (Audit Log Enhancement):**
     - Priority: High (INHERITED)
     - Sprint: Sprint 24 (INHERITED)
     - Story Points: **3** (OVERRIDDEN from 8)
     - Team: Security Team (INHERITED)

3. **Expected Behavior:**
   - Task-specific fields take precedence
   - Non-specified fields inherit from parent
   - Each task can have different overrides

**Success Criteria:** ✅ Override semantics work correctly (task fields win, parent fields fill gaps)

---

### Scenario 3: Update Existing Ticket with Field Changes (TC-INT-701.3)

**Goal:** Verify field inheritance applies when updating existing tickets.

**Test Markdown:** Reuse `test-field-inheritance-1.md` and modify

```markdown
# TICKET: E-Commerce Platform Upgrade
**Jira ID:** TEST-123  # Use actual Jira ID from Scenario 1
**Status:** To Do
**Project:** TEST
**Description:** Modernize the legacy e-commerce platform to improve performance and user experience.

**Custom Fields:**
- Priority: Critical  # CHANGED from High
- Sprint: Sprint 25   # CHANGED from Sprint 24
- Story Points: 13
- Team: Platform Engineering

## TASK: Database Migration Planning
**Jira ID:** TEST-124  # Use actual Jira ID from Scenario 1
**Status:** To Do
**Description:** Analyze current database schema and plan migration strategy.

**Custom Fields:**
- Priority: Medium  # NEW: Override parent's Critical

**Acceptance Criteria:**
- [ ] Schema analysis complete
- [ ] Migration risks identified
- [ ] Rollback plan documented

## TASK: API Modernization
**Jira ID:** TEST-125  # Use actual Jira ID from Scenario 1
**Status:** To Do
**Description:** Refactor legacy REST endpoints to use modern API standards.

**Acceptance Criteria:**
- [ ] API endpoints documented
- [ ] OpenAPI specification created
```

**Test Steps:**

1. **Push Update:**
   ```bash
   ticketr push test-field-inheritance-1.md
   ```

2. **Verify in JIRA UI:**
   - **Parent Ticket (TEST-123):**
     - Priority: **Critical** (updated)
     - Sprint: **Sprint 25** (updated)
     - Story Points: 13 (unchanged)
     - Team: Platform Engineering (unchanged)

   - **Subtask 1 (TEST-124):**
     - Priority: **Medium** (overridden, not inherited)
     - Sprint: **Sprint 25** (inherited updated value)
     - Story Points: **13** (inherited updated value)
     - Team: Platform Engineering (inherited)

   - **Subtask 2 (TEST-125):**
     - Priority: **Critical** (inherited updated value)
     - Sprint: **Sprint 25** (inherited updated value)
     - Story Points: **13** (inherited updated value)
     - Team: Platform Engineering (inherited)

**Success Criteria:** ✅ Updates propagate correctly through inheritance

---

### Scenario 4: Multiple Custom Field Types (TC-INT-701.4)

**Goal:** Test inheritance with various JIRA field types (text, select, multi-select, number).

**Test Markdown:** Create `test-field-inheritance-3.md`

```markdown
# TICKET: Payment Gateway Integration
**Status:** To Do
**Project:** TEST
**Description:** Integrate Stripe payment gateway for subscription billing.

**Custom Fields:**
- Priority: High
- Component: Backend
- Labels: payment, integration, api
- Story Points: 8
- Due Date: 2025-11-01

## TASK: API Integration
**Status:** To Do
**Description:** Implement Stripe API client.

**Custom Fields:**
- Story Points: 3
- Component: Payment Service

**Acceptance Criteria:**
- [ ] Stripe SDK integrated
- [ ] Payment endpoints created

## TASK: Webhook Handler
**Status:** To Do
**Description:** Build webhook receiver for payment events.

**Custom Fields:**
- Priority: Critical
- Story Points: 5

**Acceptance Criteria:**
- [ ] Webhook validation implemented
- [ ] Event processing logic complete
```

**Test Steps:**

1. **Push to JIRA:**
   ```bash
   ticketr push test-field-inheritance-3.md
   ```

2. **Verify Field Type Handling:**
   - **Text fields** (Component): Inherit and override correctly
   - **Select fields** (Priority): Inherit and override correctly
   - **Multi-select fields** (Labels): Inherit correctly
   - **Number fields** (Story Points): Inherit and override correctly
   - **Date fields** (Due Date): Inherit correctly

**Success Criteria:** ✅ All field types handled correctly

---

## Validation Checklist

After running all scenarios, verify:

### Functional Requirements
- [ ] **PROD-009**: Tasks inherit parent custom fields ✓
- [ ] **PROD-202**: calculateFinalFields merges fields correctly ✓
- [ ] Tasks with no custom fields inherit all parent fields
- [ ] Task-specific fields override parent values
- [ ] Non-overridden fields inherit from parent
- [ ] Updates to parent fields propagate to child tasks (unless overridden)

### Technical Validation
- [ ] No errors in console output
- [ ] Log files created in `.ticketr/logs/` with execution details
- [ ] Sensitive data redacted in logs (API keys, emails)
- [ ] State file (`.ticketr.state`) updated correctly
- [ ] Jira IDs populated in Markdown after push
- [ ] No duplicate tickets created

### JIRA API Validation
- [ ] Subtasks created with correct parent links
- [ ] Custom fields mapped correctly to JIRA field IDs
- [ ] Field values appear in JIRA UI as expected
- [ ] Updates modify existing tickets (no duplicates)
- [ ] Subtask issue type used (not Story/Epic)

---

## Troubleshooting

### Issue: "Failed to create task - parent ticket has no Jira ID"

**Cause:** Parent ticket creation failed, blocking task creation.

**Solution:**
1. Check logs in `.ticketr/logs/<timestamp>.log` for parent creation errors
2. Verify JIRA project exists and you have create permissions
3. Check Story/Epic issue type is available in project
4. Run with `--force-partial-upload` to see detailed errors

### Issue: Custom fields not appearing in JIRA

**Cause:** Field names in Markdown don't match JIRA custom field names.

**Solution:**
1. Check JIRA project settings → Fields → Custom fields
2. Ensure exact name match (case-sensitive)
3. Some fields may require JIRA admin configuration
4. Review logs for field mapping warnings

### Issue: Subtasks showing wrong field values

**Cause:** Field inheritance not applied or JIRA API caching.

**Solution:**
1. Refresh JIRA page (hard refresh with Ctrl+Shift+R)
2. Check logs for "Updated task" confirmation messages
3. Verify test scenario markdown has correct syntax
4. Run `go test ./internal/core/services -v -run TestTicketService_CalculateFinalFields` to confirm unit tests pass

### Issue: API authentication errors

**Cause:** Invalid credentials or expired API token.

**Solution:**
1. Verify `.env` file has correct values
2. Regenerate JIRA API token from Account Settings → Security
3. Check JIRA_EMAIL matches your Atlassian account email
4. Ensure JIRA_URL doesn't have trailing slash

---

## Expected Log Output

When running integration tests, you should see log entries like:

```
=== Ticketr Execution Log ===
Command: push
File: test-field-inheritance-1.md
Timestamp: 2025-10-16 15:30:00

--- Initialization ---
✓ Environment variables loaded
✓ JIRA connection verified
✓ State manager initialized

--- Validation ---
✓ 1 ticket(s) validated
✓ 2 task(s) validated
✓ No validation errors

--- Processing ---
Created ticket 'E-Commerce Platform Upgrade' with Jira ID: TEST-123
  Created task 'Database Migration Planning' with Jira ID: TEST-124
    Inherited fields: Priority=High, Sprint=Sprint 24, Story Points=13, Team=Platform Engineering
  Created task 'API Modernization' with Jira ID: TEST-125
    Inherited fields: Priority=High, Sprint=Sprint 24, Story Points=13, Team=Platform Engineering

--- Summary ---
Tickets created: 1
Tickets updated: 0
Tasks created: 2
Tasks updated: 0
Errors: 0

=== Execution Complete ===
```

---

## Next Steps After Integration Testing

1. **Document Results:**
   - Record JIRA ticket IDs for each test scenario
   - Capture screenshots of JIRA UI showing field values
   - Note any unexpected behaviors or edge cases

2. **Report Issues:**
   - If field inheritance doesn't work as expected, check unit test coverage
   - File GitHub issues with reproduction steps and JIRA screenshots
   - Include log files (with sensitive data redacted)

3. **Clean Up Test Data:**
   - Archive or delete test tickets in JIRA
   - Remove test Markdown files from repository
   - Clear `.ticketr.state` if needed for fresh testing

4. **Proceed to Milestone 8:**
   - Once field inheritance is validated, move to Pulling Tasks/Subtasks
   - Field inheritance knowledge will apply to pull operations

---

## References

- **Implementation:** `internal/core/services/ticket_service.go:39-53,109-114`
- **Tests:** `internal/core/services/ticket_service_test.go:148-280` (TC-701.1 through TC-701.4)
- **Requirements:** `REQUIREMENTS-v2.md` (PROD-009, PROD-202)
- **Documentation:** `README.md:350-447` (Field Inheritance section)
- **Examples:** `examples/field-inheritance-example.md`

---

**Ready to test?** Follow the scenarios above sequentially, validating each before moving to the next. Good luck!
