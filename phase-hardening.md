### **PRIME DIRECTIVE FOR THE EXECUTING AI AGENT**

You are an expert, test-driven software development agent executing a development phase. You **MUST** adhere to the following methodology without deviation:

1.  **Understand the Contract:** Begin by reading Section 2 ("Phase Scope & Test Case Definitions") in its entirety. This is your reference library for **what** to test and **how** to prove success.
2.  **Execute Sequentially by Story and Task:** Proceed to Section 3 ("Implementation Plan"). Address each **Story** in order. Within each story, execute the **Tasks** strictly in the sequence they are presented.
3.  **Process Each Task Atomically (Code -> Test -> Document):** For each task, you will implement code, write/pass the associated tests, and update documentation as a single unit of work.
4.  **Escalate Testing (Story & Phase Regression):**
    a.  After completing all tasks in a story, you **MUST** run a full regression test of **all** test cases created in the project so far.
    b.  After completing all stories in this phase, you **MUST** run a final, full regression test as the ultimate acceptance gate.
5.  **Commit Work:** You **MUST** create a Git commit at the completion of each story. This is a non-negotiable step.
6.  **Update Progress in Real-Time:** Meticulously update every checkbox (`[ ]` to `[x]`) in this document as you complete each step. Your progress tracking must be flawless.

## [ ] PHASE-3: The Hardening

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-3 | The Hardening |

> **As a** Lead Systems Engineer, **I want** to finalize the Ticketr v2.0 feature set by integrating robust validation, conflict management, and comprehensive reporting, **so that** the tool is reliable, safe, and ready for enterprise-wide adoption as the definitive "Tickets-as-Code" engine.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

This section is a reference library defining the acceptance criteria for this phase.

*   **Requirement:** **`PROD-201`** - **Generic `TICKET` Markdown Schema** ([Link](./REQUIREMENTS-v2.md#PROD-201))
    *   **Test Case ID:** `TC-301.1`
        *   **Test Method Signature:** `func TestTicketService_RejectsLegacyStoryFormat(t *testing.T)`
        *   **Test Logic:** (Arrange) Create a Markdown file containing the old `# STORY:` format. (Act) Pass this file to the `ticket_service`. (Assert) The service returns an error and the `ProcessResult` indicates zero tickets were processed.
        *   **Required Proof of Passing:** Console output from `go test` showing the `TestTicketService_RejectsLegacyStoryFormat` test passing.

*   **Requirement:** **`PROD-002`** - **Hierarchical Validation** ([Link](./REQUIREMENTS-v2.md#PROD-002))
    *   **Test Case ID:** `TC-302.1`
        *   **Test Method Signature:** `func TestPushCommand_FailsFastOnValidationError(t *testing.T)`
        *   **Test Logic:** (Arrange) Create a Markdown file with a known validation error (e.g., a "Sub-task" under an "Epic"). Mock the `JiraAdapter` to fail the test if any of its `Create/Update` methods are called. (Act) Execute the `push` command logic. (Assert) The command exits with a non-zero status code, prints the validation error, and the mock confirms that no API calls were made.
        *   **Required Proof of Passing:** A test that mocks the CLI execution and verifies the `JiraAdapter` was never called, along with the captured error output.

*   **Requirement:** **`N/A`** - **Conflict Detection (New Functionality)**
    *   **Test Case ID:** `TC-303.1`
        *   **Test Method Signature:** `func TestPullService_DetectsConflictState(t *testing.T)`
        *   **Test Logic:** (Arrange) Create a `pull_service` and a `StateManager`. Pre-populate the state file with `{"TICKET-1": {"local_hash": "A", "remote_hash": "B"}}`. Prepare a local Markdown file whose TICKET-1 content hashes to "C", and mock a Jira response for TICKET-1 that hashes to "D". (Act) Run the `pull` service. (Assert) The service returns a specific `ErrConflictDetected` error for TICKET-1.
        *   **Required Proof of Passing:** Console output from `go test` showing the `TestPullService_DetectsConflictState` test passing.

*   **Requirement:** **`USER-001`** - **Non-Interactive Error Handling** ([Link](./REQUIREMENTS-v2.md#USER-001))
    *   **Test Case ID:** `TC-304.1`
        *   **Test Method Signature:** `func TestPushService_ProcessesAllAndReportsFailures(t *testing.T)`
        *   **Test Logic:** (Arrange) Create a Markdown file with three tickets. Mock the `JiraAdapter` to succeed on ticket 1 and 3, but fail on ticket 2. (Act) Run the `push` service without the `--force` flag. (Assert) The `ProcessResult` contains 2 successes and 1 failure. The service itself returns an error, but the mock confirms that API calls were attempted for all three tickets.
        *   **Required Proof of Passing:** Console output from `go test` showing the test passing, along with the contents of the final `ProcessResult` struct.

---

### **3. Implementation Plan (The Execution)**

#### [ ] STORY-301: Solidify the Core Workflow

1.  **Task:** Eliminate all legacy `Story`-based models and code paths.
    *   **Instruction:** `Perform a global search-and-replace for the 'Story' domain model and its related service/repository methods. Remove the type aliases and legacy methods. Refactor all calling code, including all test files, to use the generic 'Ticket' model and its associated methods exclusively.`
    *   **Fulfills:** This task contributes to requirement **`PROD-201`**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-301.1`:**
            *   [x] **Test Method Created:** Checked after the test method is written. **Evidence:** Complete test in `/home/karol/dev/private/ticktr/internal/core/services/ticket_service_test.go` lines 12-54.
            *   [x] **Test Method Passed:** Checked after the test passes. **Evidence:** 
```
=== RUN   TestTicketService_RejectsLegacyStoryFormat
--- PASS: TestTicketService_RejectsLegacyStoryFormat (0.00s)
PASS
ok  	github.com/karolswdev/ticktr/internal/core/services	0.002s
```
    *   **Documentation:**
        *   [x] **Documentation Updated:** Checked after the relevant documentation is updated. **Instruction:** `Update the HANDOFF-BEFORE-PHASE-3.md file, removing all sections related to "Backward Compatibility" and type aliases. State clearly that v2.0 is a breaking change.` **Evidence:** Updated section 1 of Critical Implementation Notes to state "v2.0 is a breaking change from v1.0. The legacy Story model and all related code paths have been removed."

2.  **Task:** Integrate pre-flight validation directly into the `push` command.
    *   **Instruction:** `In cmd/ticketr/main.go, within the runPush function, add logic to instantiate the internal/core/validation.Validator. Before calling the push_service, execute a full validation pass on the parsed tickets. If any validation errors are found, print them to the console and os.Exit(1).`
    *   **Fulfills:** This task contributes to requirement **`PROD-002`**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-302.1`:**
            *   [x] **Test Method Created:** **Evidence:** Complete test in `/home/karol/dev/private/ticktr/cmd/ticketr/main_validation_test.go` lines 13-65.
            *   [x] **Test Method Passed:** **Evidence:** 
```
=== RUN   TestPushCommand_FailsFastOnValidationError
--- PASS: TestPushCommand_FailsFastOnValidationError (0.00s)
PASS
ok  	github.com/karolswdev/ticktr/cmd/ticketr	0.002s
```
    *   **Documentation:**
        *   [x] **Documentation Updated:** Checked after the relevant documentation is updated. **Instruction:** `Update the README.md in the "Push Command" section to include a note: "Note: Ticketr validates your file for correctness before sending any data to Jira, preventing partial failures."` **Evidence:** Added new Push Command section at line 137-142 with validation details.

---
> ### **Story Completion: STORY-301**
>
> You may only proceed once all checkboxes for all tasks within this story are marked `[x]`. Then, you **MUST** complete the following steps in order:
>
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** Checked after running all tests created in the project up to this point.
>     *   **Instruction:** `Execute 'go test ./... -v'.`
>     *   **Evidence:** Provide the full summary output from the test runner, showing the total number of tests executed and confirming all have passed.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** Checked after creating the Git commit.
>     *   **Instruction:** `Execute 'git add .' followed by 'git commit -m "feat(core): Solidify core workflow and integrate pre-flight validation"'.`
>     *   **Evidence:** Provide the full commit hash returned by the Git command.
> 3.  **Finalize Story:**
>     *   **Instruction:** Once the two checkboxes above are complete, you **MUST** update this story's main checkbox from `[ ]` to `[x]`.

---

#### [ ] STORY-302: Implement Bidirectional Conflict Management

1.  **Task:** Evolve the State Manager for bidirectional hash tracking.
    *   **Instruction:** `Modify internal/state/manager.go. The internal state map must be changed from map[string]string to map[string]struct{ LocalHash string; RemoteHash string }. Update all associated methods (Load, Save, UpdateHash) to handle this new structure.`
    *   **Fulfills:** This task contributes to the new **Conflict Detection** functionality.
    *   **Verification via Test Cases:**
        *   [ ] **Tests Updated & Passed:** **Instruction:** `Update existing StateManager tests to reflect the new data structure.` **Evidence:** `[...]`
    *   **Documentation:**
        *   [ ] **Documentation Updated:** Checked after the relevant documentation is updated. **Instruction:** `Update the "State Management" section in HANDOFF-BEFORE-PHASE-3.md to describe the new state file JSON structure.` **Evidence:** `[...]`

2.  **Task:** Implement conflict detection and safe merge logic in the `pull` service.
    *   **Instruction:** `Create internal/core/services/pull_service.go. Implement the pull logic which modifies an existing file. This service must use the new StateManager to detect conflicts (local_hash changed AND remote_hash changed). If a conflict is detected, it must return a specific error. If only the remote has changed, it must update the local file and the state. The runPull function in main.go must be updated to use this new service.`
    *   **Fulfills:** This task contributes to **Conflict Detection** and **`PROD-010`**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-303.1`:**
            *   [ ] **Test Method Created:** **Evidence:** `[...]`
            *   [ ] **Test Method Passed:** **Evidence:** `[...]`
    *   **Documentation:**
        *   [ ] **Documentation Updated:** Checked after the relevant documentation is updated. **Instruction:** `Update the "Pull Command" section in README.md to explain the new in-place update behavior and the conflict detection mechanism.` **Evidence:** `[...]`

3.  **Task:** Refactor `push` service for comprehensive reporting.
    *   **Instruction:** `Modify internal/core/services/push_service.go to remove the "fail-fast" behavior. The service must now iterate through all tickets, attempt to process each one, and aggregate all successes and failures into the ProcessResult. The service should only return an error at the end if one or more tickets failed.`
    *   **Fulfills:** This task contributes to requirement **`USER-001`**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-304.1`:**
            *   [ ] **Test Method Created:** **Evidence:** `[...]`
            *   [ ] **Test Method Passed:** **Evidence:** `[...]`
    *   **Documentation:**
        *   [ ] **Documentation Updated:** Checked after the relevant documentation is updated. **Instruction:** `Update the REQUIREMENTS-v2.md description for USER-001 to reflect that the tool processes all tickets and provides a summary report, exiting with an error code if any failures occurred.` **Evidence:** `[...]`

---
> ### **Story Completion: STORY-302**
>
> You may only proceed once all checkboxes for all tasks within this story are marked `[x]`. Then, you **MUST** complete the following steps in order:
>
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** Checked after running all tests created in the project up to this point.
>     *   **Instruction:** `Execute 'go test ./... -v'.`
>     *   **Evidence:** Provide the full summary output from the test runner, showing the total number of tests executed and confirming all have passed.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** Checked after creating the Git commit.
>     *   **Instruction:** `Execute 'git add .' followed by 'git commit -m "feat(sync): Implement bidirectional conflict management"'.`
>     *   **Evidence:** Provide the full commit hash returned by the Git command.
> 3.  **Finalize Story:**
>     *   **Instruction:** Once the two checkboxes above are complete, you **MUST** update this story's main checkbox from `[ ]` to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-[ID]` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   **Instruction:** You are at the final gate for this phase. Before marking the entire phase as done, you must perform one last, full regression test to ensure nothing was broken by the final commits.
*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./... -v' one last time.`
    *   **Evidence:** Provide the full, final summary output from the test runner, showing the grand total of tests for this phase and confirming that 100% have passed.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox above is marked `[x]`, your final action for this phase is to modify the main title of this document, changing `[ ] PHASE-3` to `[x] PHASE-3`. This concludes your work on this phase file.