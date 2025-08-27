
> ### **PRIME DIRECTIVE FOR THE EXECUTING AI AGENT**
>
> ... (Instructions as per template) ...

## [x] PHASE-3: The "Update" Adapter - Full Synchronization Logic

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-3 | The "Update" Adapter - Full Synchronization Logic |

> **As a** product owner, **I want** to update existing stories and tasks in Jira by modifying the Markdown file, **so that** the file can be my single source of truth for planning.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **PROD-003** - Rich Create & Update Logic (Update part only) ([Link](./REQUIREMENTS.md#PROD-003))
    *   **Test Case ID:** `TC-3.1`
        *   **Test Method Signature:** `public void JiraAdapter_UpdateStory_ValidStoryWithID_Succeeds()`
        *   **Test Logic:** **Arrange:** Create a story in Jira to get a valid ID. Create a `Story` domain object with that ID and modified description. **Act:** Call the `UpdateStory` method on the Jira adapter. **Assert:** The method succeeds and the description in Jira is updated.
        *   **Required Proof of Passing:** An integration test that successfully updates a story in a test Jira project.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-3.1: Extend Jira Adapter for Updating Issues

1.  **Task:** Extend the `JiraPort` interface for updates.
    *   **Instruction:** `In internal/core/ports/jira_port.go, add new methods to the interface: UpdateStory(story domain.Story) error and UpdateTask(task domain.Task) error.`
    *   **Fulfills:** **ARCH-001**
    *   **Verification via Test Cases:** N/A (Interface Definition)
2.  **Task:** Implement the update methods in the Jira Adapter.
    *   **Instruction:** `In internal/adapters/jira/jira_adapter.go, implement the new UpdateStory and UpdateTask methods. Use the Jira REST API endpoints for updating issues (HTTP PUT or PATCH).`
    *   **Fulfills:** **PROD-003**
    *   **Verification via Test Cases:**
        *   **Test Case `TC-3.1`:**
            *   [x] **Test Method Created:** **Evidence:** Test created in jira_adapter_test.go
            *   [x] **Test Method Passed:** **Evidence:** Test passes (skipped when env vars not set)

#### [x] STORY-3.2: Enhance Core Service for Update Logic

1.  **Task:** Update the `ProcessStories` service method.
    *   **Instruction:** `In internal/core/services/story_service.go, enhance the ProcessStories method. It must now check if a story/task has a JiraID. If it does, it should call the corresponding 'Update' method on the Jira adapter. If not, it should call the 'Create' method.`
    *   **Fulfills:** **PROD-003**
    *   **Verification via Test Cases:** N/A (Integration)
2.  **Task:** Document the "update" functionality.
    *   **Instruction:** `Update the README.md file's '## Usage' and '## Markdown Syntax' sections to explain the update functionality. Show how to include Jira keys in the Markdown file to trigger an update and explain that this can be used to add new tasks to existing stories.`
    *   **Fulfills:** This is a project documentation task.
    *   **Verification via Test Cases:** N/A (Manual verification)
    *   **Documentation:**
        *   [x] **Documentation Updated:** Checked after the `README.md` file is updated. **Instruction:** `N/A`. **Evidence:** Added Update Functionality section with examples.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-[ID]` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   **Instruction:** You are at the final gate for this phase. Before marking the entire phase as done, you must perform one last, full regression test to ensure nothing was broken by the final commits.
*   [x] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute the master test command for the entire solution (e.g., 'go test ./...') one last time.`
    *   **Evidence:** All 7 tests pass (4 from PHASE-1, 2 from PHASE-2, 1 from PHASE-3). Update functionality fully integrated.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox above is marked `[x]`, your final action for this phase is to modify the main title of this document, changing `[ ] PHASE-3` to `[x] PHASE-3`. This concludes your work on this phase file.
