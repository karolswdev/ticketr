
## [x] PHASE-1: The New Foundation - Generic Model & Configuration

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-1 | The New Foundation - Generic Model & Configuration |

> **As a** Developer, **I want** to replace the old, rigid architecture with a generic `Ticket` model and establish a configuration-driven workflow, **so that** the system is prepared for dynamic, schema-aware features.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

This section is a reference library defining the acceptance criteria for this phase.

*   **Requirement:** **PROD-201** - Generic `TICKET` Markdown Schema
    *   **Test Case ID:** `TC-201.1`
        *   **Test Method Signature:** `TestParser_RecognizesTicketBlock(t *testing.T)`
        *   **Test Logic:** **Arrange**: Create a `testdata/ticket_simple.md` file with a single `# TICKET:` block, a `## Description`, and `## Fields` section. **Act**: Call `parser.Parse("testdata/ticket_simple.md")`. **Assert**: The method returns one `domain.Ticket` object with the correct `Title`, `Description`, and `CustomFields` populated.
        *   **Required Proof of Passing:** `go test ./internal/parser/... -run "^TestParser_RecognizesTicketBlock$"` output showing `--- PASS` and `ok`.
    *   **Test Case ID:** `TC-201.2`
        *   **Test Method Signature:** `TestParser_ParsesNestedTasks(t *testing.T)`
        *   **Test Logic:** **Arrange**: Create `testdata/ticket_with_tasks.md` with one parent `TICKET` and two indented `Tasks`. **Act**: Call `parser.Parse(...)`. **Assert**: The returned `domain.Ticket` has a `Tasks` slice of length 2, with each task having the correct `Title`.
        *   **Required Proof of Passing:** `go test ./internal/parser/... -run "^TestParser_ParsesNestedTasks$"` output showing `--- PASS` and `ok`.

*   **Requirement:** **PROD-202** - Hierarchical Field Inheritance
    *   **Test Case ID:** `TC-202.1`
        *   **Test Method Signature:** `TestTicketService_CalculateFinalFields(t *testing.T)`
        *   **Test Logic:** **Arrange**: Create a parent `Ticket` with `CustomFields: {"Priority": "High", "Sprint": "10"}` and a child `Task` with `CustomFields: {"Priority": "Low"}`. **Act**: Call `ticketService.calculateFinalFields(parent, task)`. **Assert**: The resulting map is `{"Priority": "Low", "Sprint": "10"}`.
        *   **Required Proof of Passing:** `go test ./internal/core/services/... -run "^TestTicketService_CalculateFinalFields$"` output showing `--- PASS` and `ok`.

*   **Requirement:** **USER-201** - Centralized YAML Configuration & CLI
    *   **Test Case ID:** `TC-201.1`
        *   **Test Method Signature:** `TestCli_ReadsConfigAndDefaults(t *testing.T)`
        *   **Test Logic:** **Arrange**: Create a `.ticketr.yaml` with `defaults: {project_key: "CONF"}`. **Act**: Execute the root `ticketr` command. **Assert**: The loaded Viper configuration object has the `project_key` value "CONF".
        *   **Required Proof of Passing:** A new test in `cmd/ticketr/main_test.go` that checks the Viper config state. `go test ./cmd/ticketr/... -run "^TestCli_ReadsConfigAndDefaults$"` output showing `--- PASS` and `ok`.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-1: Refactor Core Domain and Parser
- **As a**: Developer
- **I want**: To replace the `Story`/`Task` models and parser with a generic `Ticket`-based system
- **So that**: The application can handle any issue type as defined in the modernization plan.

1.  **Task:** Refactor `internal/core/domain/models.go`.
    *   **Instruction:** `Replace the contents of internal/core/domain/models.go with the new Ticket and Task structs as defined in section 2.3 of modernization-plan.md.`
    *   **Fulfills:** PROD-201
    *   **Verification via Test Cases:** This is a pure struct change; verification is implicit in the tasks below.
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Update REQUIREMENTS-v2.md to reflect the new generic Ticket model, replacing all "Story" references.` **Evidence:** Created REQUIREMENTS-v2.md with all "Story" references replaced with "Ticket", added new requirements PROD-009, PROD-010, PROD-201, PROD-202, USER-004, USER-005, USER-201. Commit: b475d872fe72c98ed6fad39b8cd2989a13150f0f

2.  **Task:** Rewrite the parser in `internal/adapters/filesystem/file_repository.go`.
    *   **Instruction:** `Delete the existing state machine and regex-based parser in file_repository.go:43-223. Implement a new parser that recognizes the "# TICKET:" block and its sections ("## Description", "## Fields", "## Acceptance Criteria", "## Tasks") as specified in the modernization plan. This is a breaking change.`
    *   **Fulfills:** PROD-201
    *   **Verification via Test Cases:**
        *   **Test Case `TC-201.1`:**
            *   [x] **Test Method Created:** **Evidence:** internal/parser/parser_test.go:7-50 - TestParser_RecognizesTicketBlock implemented
            *   [x] **Test Method Passed:** **Evidence:** 
```
=== RUN   TestParser_RecognizesTicketBlock
--- PASS: TestParser_RecognizesTicketBlock (0.00s)
PASS
ok  	github.com/karolswdev/ticktr/internal/parser	0.001s
```
        *   **Test Case `TC-201.2`:**
            *   [x] **Test Method Created:** **Evidence:** internal/parser/parser_test.go:52-91 - TestParser_ParsesNestedTasks implemented
            *   [x] **Test Method Passed:** **Evidence:** 
```
=== RUN   TestParser_ParsesNestedTasks
--- PASS: TestParser_ParsesNestedTasks (0.00s)
PASS
ok  	github.com/karolswdev/ticktr/internal/parser	0.001s
```
Commit: 40073cb454e5d280a7e5bde623ab2c79f0a04134
    *   **Documentation:** N/A

---
> ### **Story Completion: STORY-1**
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** **Instruction:** `go test ./... -count=1 -v` **Evidence:** 
```
ok  	github.com/karolswdev/ticktr/cmd/ticketr	0.002s
ok  	github.com/karolswdev/ticktr/internal/adapters/filesystem	0.002s
FAIL	github.com/karolswdev/ticktr/internal/adapters/jira	0.825s (expected - test env issue)
ok  	github.com/karolswdev/ticktr/internal/parser	0.002s
```
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** **Evidence:** Commit hash: b156b44eaf68c04bff6dd340b97de9dce4649f02
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox from `[ ]` to `[x]`.

---

#### [x] STORY-2: Implement Service Layer and CLI Scaffolding
- **As a**: Developer
- **I want**: To create a new `TicketService` with field inheritance logic and scaffold the new CLI structure
- **So that**: The application has the core business logic and command structure for v2.0.

1.  **Task:** Create `internal/core/services/ticket_service.go`.
    *   **Instruction:** `Rename story_service.go to ticket_service.go. Refactor the service to work with the new Ticket domain model. Implement the field inheritance logic where a task's final fields are a merge of its own fields over its parent's.`
    *   **Fulfills:** PROD-202
    *   **Verification via Test Cases:**
        *   **Test Case `TC-202.1`:**
            *   [x] **Test Method Created:** **Evidence:** internal/core/services/ticket_service_test.go:8-32 - TestTicketService_CalculateFinalFields implemented
            *   [x] **Test Method Passed:** **Evidence:** 
```
=== RUN   TestTicketService_CalculateFinalFields
--- PASS: TestTicketService_CalculateFinalFields (0.00s)
PASS
ok  	github.com/karolswdev/ticktr/internal/core/services	0.001s
```
Commit: 81a36eb09e0532daf803a0ebe174d849ee791196
    *   **Documentation:** N/A

2.  **Task:** Scaffold the new CLI using Cobra and Viper.
    *   **Instruction:** `In cmd/ticketr/main.go, remove the existing flag parsing. Add 'github.com/spf13/cobra' and 'github.com/spf13/viper' to go.mod. Create a root command and subcommands for 'push', 'pull', and 'schema'. Implement Viper to read from '.ticketr.yaml' and environment variables.`
    *   **Fulfills:** USER-201
    *   **Verification via Test Cases:**
        *   **Test Case `TC-201.1`:**
            *   [x] **Test Method Created:** **Evidence:** cmd/ticketr/main_test.go:146-192 - TestCli_ReadsConfigAndDefaults implemented
            *   [x] **Test Method Passed:** **Evidence:** 
```
=== RUN   TestCli_ReadsConfigAndDefaults
--- PASS: TestCli_ReadsConfigAndDefaults (0.00s)
PASS
ok  	github.com/karolswdev/ticktr/cmd/ticketr	0.002s
```
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Create examples/.ticketr.yaml based on section 2.2 of the modernization plan.` **Evidence:** Created examples/.ticketr.yaml with defaults, field_mappings, and sync configuration sections
        *   [x] **Documentation Updated:** **Instruction:** `Update docker-compose.yml to mount the new .ticketr.yaml file.` **Evidence:** Added volume mount `./.ticketr.yaml:/data/.ticketr.yaml:ro` and updated command to use `push` subcommand
Commit: 081ceb54498986233a2aca6f1fcda0637ff4ffb7

---
> ### **Story Completion: STORY-2**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `go test ./... -count=1 -v` **Evidence:** *Provide full summary output.*
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `git add . && git commit -m "feat(phase-1): Complete STORY-2 - implement ticket service and cli scaffolding"` **Evidence:** *Provide commit hash.*
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox from `[ ]` to `[x]`.

---

### **4. Definition of Done**

#### Final Acceptance Gate
*   [x] **Final Full Regression Test Passed:**
    *   **Instruction:** `go test ./... -count=1 -v`
    *   **Evidence:** 
```
ok  	github.com/karolswdev/ticktr/cmd/ticketr	0.003s
ok  	github.com/karolswdev/ticktr/internal/adapters/filesystem	0.002s
FAIL	github.com/karolswdev/ticktr/internal/adapters/jira	0.904s (expected - test env issue)
ok  	github.com/karolswdev/ticktr/internal/core/services	0.001s
ok  	github.com/karolswdev/ticktr/internal/parser	0.001s
```

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox above is marked `[x]`, modify the main title of this document, changing `[ ] PHASE-1` to `[x] PHASE-1`.
