
> ### **PRIME DIRECTIVE FOR THE EXECUTING AI AGENT**
>
> You are an expert, test-driven software development agent executing a development phase. You **MUST** adhere to the following methodology without deviation:
>
> 1.  **Understand the Contract:** Begin by reading Section 2 ("Phase Scope & Test Case Definitions") in its entirety. This is your reference library for **what** to test and **how** to prove success.
> 2.  **Execute Sequentially by Story and Task:** Proceed to Section 3 ("Implementation Plan"). Address each **Story** in order. Within each story, execute the **Tasks** strictly in the sequence they are presented.
> 3.  **Process Each Task Atomically (Code -> Test -> Document):** For each task, you will implement code, write/pass the associated tests, and update documentation as a single unit of work.
> 4.  **Escalate Testing (Story & Phase Regression):**
>     a.  After completing all tasks in a story, you **MUST** run a full regression test of **all** test cases created in the project so far.
>     b.  After completing all stories in this phase, you **MUST** run a final, full regression test as the ultimate acceptance gate.
> 5.  **Commit Work:** You **MUST** create a Git commit at the completion of each story. This is a non-negotiable step.
> 6.  **Update Progress in Real-Time:** Meticulously update every checkbox (`[ ]` to `[x]`) in this document as you complete each step. Your progress tracking must be flawless.

## [ ] PHASE-1: The Foundation - Core Logic & Parsing

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-1 | The Foundation - Core Logic & Parsing |

> **As a** developer, **I want** a robust, testable core library that can parse the custom Markdown format into domain models, **so that** I have a reliable foundation for building the Jira integration.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

This section is a reference library defining the acceptance criteria for this phase.

*   **Requirement:** **PROD-001** - File-Based State ([Link](./REQUIREMENTS.md#PROD-001))
    *   **Test Case ID:** `TC-1.1`
        *   **Test Method Signature:** `public void Parser_ParseInput_ValidFile_ReturnsCorrectStoryCount()`
        *   **Test Logic:** **Arrange:** Create a string representing a valid Markdown input with two distinct stories. **Act:** Pass the string to the parser. **Assert:** The parser returns a slice containing exactly two `Story` objects.
        *   **Required Proof of Passing:** The unit test for this case must pass.
*   **Requirement:** **PROD-005** - Rich Task Definitions ([Link](./REQUIREMENTS.md#PROD-005))
    *   **Test Case ID:** `TC-1.2`
        *   **Test Method Signature:** `public void Parser_ParseInput_TaskWithDetails_CorrectlyPopulatesTaskFields()`
        *   **Test Logic:** **Arrange:** Create a Markdown string for a single story with one task that has a nested `Description` and `Acceptance Criteria`. **Act:** Parse the string. **Assert:** The resulting `Task` object has its `Description` and `AcceptanceCriteria` fields correctly populated.
        *   **Required Proof of Passing:** The unit test for this case must pass.
*   **Requirement:** **PROD-003** - Create & Update Operations ([Link](./REQUIREMENTS.md#PROD-003))
    *   **Test Case ID:** `TC-1.3`
        *   **Test Method Signature:** `public void Parser_ParseInput_WithAndWithoutJiraKeys_CorrectlyPopulatesIDs()`
        *   **Test Logic:** **Arrange:** Create a Markdown string with one story and one task that include Jira keys (e.g., `[PROJ-123]`), and one story and one task without. **Act:** Parse the string. **Assert:** The `JiraID` field is correctly populated for the items that have keys and is empty for those that do not.
        *   **Required Proof of Passing:** The unit test for this case must pass.
*   **Requirement:** **PROD-002** - Hierarchical Validation ([Link](./REQUIREMENTS.md#PROD-002))
    *   **Test Case ID:** `TC-1.4`
        *   **Test Method Signature:** `public void Parser_ParseInput_MalformedStoryHeading_ReturnsErrorAndNoStories()`
        *   **Test Logic:** **Arrange:** Create a Markdown string where a story heading is malformed (e.g., `## STORY:` instead of `# STORY:`). **Act:** Parse the string. **Assert:** The parser returns an error and an empty slice of stories.
        *   **Required Proof of Passing:** The unit test for this case must pass.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-1.1: Scaffold Project and Define Core Models

1.  **Task:** Initialize the Git repository.
    *   **Instruction:** `Execute the command 'git init'.`
    *   **Fulfills:** This is a foundational project setup task.
    *   **Verification via Test Cases:** N/A (Repo setup)
    *   **Documentation:** N/A
2.  **Task:** Create the initial directory structure for the Go application, adhering to the Ports and Adapters pattern.
    *   **Instruction:** `Create the following directories: cmd/jira-story-creator, internal/core/domain, internal/core/ports, internal/core/services, internal/adapters/cli, internal/adapters/filesystem, internal/adapters/jira.`
    *   **Fulfills:** This task contributes to requirement **ARCH-001**.
    *   **Verification via Test Cases:** N/A (Structural setup)
    *   **Documentation:** N/A
3.  **Task:** Initialize the Go module.
    *   **Instruction:** `Execute the command 'go mod init github.com/karolswdev/ticktr' in the root directory.`
    *   **Fulfills:** This task contributes to requirement **TECH-P-001**.
    *   **Verification via Test Cases:** N/A (Project setup)
    *   **Documentation:** N/A
4.  **Task:** Define the core domain models (`Story`, `Task`) in Go.
    *   **Instruction:** `Create a file named internal/core/domain/models.go and define the Go structs for Story and Task based on the requirements. Include fields for Title, Description, AcceptanceCriteria, JiraID, and a slice of Tasks for the Story struct.`
    *   **Fulfills:** This task contributes to requirements **PROD-001**, **PROD-003**, **PROD-005**.
    *   **Verification via Test Cases:** N/A (Model definition)
    *   **Documentation:**
        *   [x] **Documentation Updated:** Checked after the Go structs are defined. **Instruction:** `Add comments to the structs in models.go explaining the purpose of each field.` **Evidence:** Provide the complete code for the `models.go` file.
5.  **Task:** Create the initial project README file.
    *   **Instruction:** `Create a file named README.md in the root directory. Add a main title '# ticktr', a brief description of the project's purpose, and a section '## Project Structure' that lists and explains the directories created in this story.`
    *   **Fulfills:** This is a project documentation task.
    *   **Verification via Test Cases:** N/A (Manual verification)
    *   **Documentation:**
        *   [x] **Documentation Updated:** Checked after the `README.md` file is created. **Instruction:** `N/A`. **Evidence:** Provide the full content of the new `README.md` file.

---
> ### **Story Completion: STORY-1.1**
>
> You may only proceed once all checkboxes for all tasks within this story are marked `[x]`. Then, you **MUST** complete the following steps in order:
>
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** Checked after running all tests created in the project up to this point.
>     *   **Instruction:** `Execute the master test command for the entire solution (e.g., 'go test ./...').`
>     *   **Evidence:** Provide the full summary output from the test runner.
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** Checked after creating the Git commit.
>     *   **Instruction:** `Execute 'git add .' followed by 'git commit -m "feat(story): Complete STORY-1.1 - Scaffold Project and Define Core Models"'.`
>     *   **Evidence:** Provide the full commit hash returned by the Git command.
> 3.  **Finalize Story:**
>     *   **Instruction:** Once the two checkboxes above are complete, you **MUST** update this story's main checkbox from `[ ]` to `[x]`.

---

#### [ ] STORY-1.2: Implement Markdown Parser

1.  **Task:** Create the file-based repository port and adapter.
    *   **Instruction:** `In internal/core/ports, create a file repository.go and define a Repository interface with methods like GetStories(filepath string) ([]domain.Story, error) and SaveStories(filepath string, stories []domain.Story) error. Then, create an initial implementation of this interface in internal/adapters/filesystem/file_repository.go.`
    *   **Fulfills:** This task contributes to requirement **ARCH-001**.
    *   **Verification via Test Cases:** N/A (Interface definition)
    *   **Documentation:** N/A
2.  **Task:** Implement the parsing logic for the `GetStories` method.
    *   **Instruction:** `In the file_repository.go adapter, implement the GetStories method. It should read the file content and parse it according to our defined Markdown syntax. Use a suitable Markdown parsing library if available, or write a manual line-by-line parser.`
    *   **Implementation Details from C# Analysis:**
        *   The parser should operate like a state machine, processing the file line by line.
        *   The current "state" (e.g., parsing a story's description, parsing tasks, parsing a task's AC) determines how each line is interpreted.
        *   A line starting with `# STORY:` marks the beginning of a new story object and resets the state.
        *   A line starting with `## Acceptance Criteria` or `## Tasks` changes the state to parse the corresponding section.
        *   Within the "Tasks" state, a line starting with `-` or `*` indicates a new task.
        *   Indented lines underneath a task item should be parsed as that task's details (e.g., `  - Description:`).
        *   The parser must be able to extract Jira keys (e.g., `[PROJ-123]`) from story and task lines.
        *   Use a `---` line to separate stories, which is a more robust separator than the C# implementation.
    *   **Fulfills:** This task contributes to requirements **PROD-001**, **PROD-002**, **PROD-003**, **PROD-005**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-1.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the complete code for the test method.
            *   [ ] **Test Method Passed:** **Evidence:** Provide the console output from the test runner.
        *   **Test Case `TC-1.2`:**
            *   [ ] **Test Method Created:** **Evidence:** `[...]`
            *   [ ] **Test Method Passed:** **Evidence:** `[...]`
        *   **Test Case `TC-1.3`:**
            *   [ ] **Test Method Created:** **Evidence:** `[...]`
            *   [ ] **Test Method Passed:** **Evidence:** `[...]`
        *   **Test Case `TC-1.4`:**
            *   [ ] **Test Method Created:** **Evidence:** `[...]`
            *   [ ] **Test Method Passed:** **Evidence:** `[...]`
    *   **Documentation:** N/A
3.  **Task:** Create the Markdown Specification ICD.
    *   **Instruction:** `Create the file STORY-MARKDOWN-SPEC.md with the approved content, formally defining the custom Markdown syntax.`
    *   **Fulfills:** This is a project documentation task.
    *   **Verification via Test Cases:** N/A (Manual verification)
    *   **Documentation:**
        *   [ ] **Documentation Updated:** Checked after the `STORY-MARKDOWN-SPEC.md` file is created. **Instruction:** `N/A`. **Evidence:** Provide the full content of the new `STORY-MARKDOWN-SPEC.md` file.
4.  **Task:** Link to the Specification from the README.
    *   **Instruction:** `Update the README.md file. In the '## Markdown Syntax' section, remove the detailed syntax and replace it with a link to the formal specification document, e.g., "The full specification for the Ticktr Markdown Syntax can be found in [STORY-MARKDOWN-SPEC.md](./STORY-MARKDOWN-SPEC.md)."`
    *   **Fulfills:** This is a project documentation task.
    *   **Verification via Test Cases:** N/A (Manual verification)
    *   **Documentation:**
        *   [ ] **Documentation Updated:** Checked after the `README.md` file is updated. **Instruction:** `N/A`. **Evidence:** Provide a diff of the changes to `README.md`.

---
> ### **Story Completion: STORY-1.2**
>
> ... (Completion block for Story 1.2 would follow the same structure) ...

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-[ID]` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   **Instruction:** You are at the final gate for this phase. Before marking the entire phase as done, you must perform one last, full regression test to ensure nothing was broken by the final commits.
*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute the master test command for the entire solution (e.g., 'go test ./...') one last time.`
    *   **Evidence:** Provide the full, final summary output from the test runner, showing the grand total of tests for this phase and confirming that 100% have passed.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox above is marked `[x]`, your final action for this phase is to modify the main title of this document, changing `[ ] PHASE-1` to `[x] PHASE-1`. This concludes your work on this phase file.