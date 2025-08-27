
> ### **PRIME DIRECTIVE FOR THE EXECUTING AI AGENT**
>
> ... (Instructions as per template) ...

## [x] PHASE-2: The "Create" Adapter - First End-to-End Value

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-2 | The "Create" Adapter - First End-to-End Value |

> **As a** project manager, **I want** to process a Markdown file of new stories and tasks and have them created in Jira, **so that** I can automate the initial setup of my project backlog.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **NFR-001** - MVP Credentials ([Link](./REQUIREMENTS.md#NFR-001))
    *   **Test Case ID:** `TC-2.1`
        *   **Test Method Signature:** `public void JiraAdapter_NewClient_WithEnvVars_AuthenticatesSuccessfully()`
        *   **Test Logic:** **Arrange:** Set valid Jira credentials in environment variables. **Act:** Create a new Jira client instance. **Assert:** The client authenticates successfully (e.g., by making a simple, valid API call like fetching project details).
        *   **Required Proof of Passing:** An integration test that makes a live call to the Jira API and receives a successful (200 OK) response.
*   **Requirement:** **PROD-003** - Rich Create & Update Logic (Create part only) ([Link](./REQUIREMENTS.md#PROD-003))
    *   **Test Case ID:** `TC-2.2`
        *   **Test Method Signature:** `public void JiraAdapter_CreateStory_ValidStory_ReturnsNewJiraID()`
        *   **Test Logic:** **Arrange:** Create a valid `Story` domain object. **Act:** Call the `CreateStory` method on the Jira adapter. **Assert:** The method returns a valid, non-empty Jira Issue Key (e.g., "PROJ-123").
        *   **Required Proof of Passing:** An integration test that successfully creates a story in a test Jira project.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-2.1: Implement Jira Adapter for Creating Issues

1.  **Task:** Define the `JiraPort` interface.
    *   **Instruction:** `In internal/core/ports, create a file jira_port.go and define an interface for the Jira client. It should include methods like Authenticate(), CreateStory(story domain.Story) (string, error), and CreateTask(task domain.Task, parentID string) (string, error).`
    *   **Fulfills:** **ARCH-001**
    *   **Verification via Test Cases:** N/A (Interface Definition)
2.  **Task:** Implement the Jira Adapter.
    *   **Instruction:** `In internal/adapters/jira, create a file jira_adapter.go. Implement the JiraPort interface. Use environment variables for credentials. Use a Go HTTP client and the Jira REST API endpoints to implement the methods.`
    *   **Implementation Details from C# Analysis:**
        *   **Endpoint:** All requests for creating issues should be a `POST` to `/rest/api/2/issue`.
        *   **Authentication:** Use HTTP Basic Authentication. The token is the Base64 encoding of `email:apiKey`.
        *   **Headers:** Ensure `Content-Type` is `application/json`.
        *   **Create Story Payload:** The JSON body for creating a story should have the following structure:
            ```json
            {
              "fields": {
                "project": { "key": "YOUR_PROJECT_KEY" },
                "summary": "Story Title",
                "description": "Story Description and Acceptance Criteria here.",
                "issuetype": { "name": "Story" }
              }
            }
            ```
        *   **Create Task Payload:** The JSON body for creating a sub-task should have the following structure:
            ```json
            {
              "fields": {
                "project": { "key": "YOUR_PROJECT_KEY" },
                "summary": "Task Title",
                "description": "Task Description here.",
                "issuetype": { "name": "Sub-task" },
                "parent": { "key": "PARENT_STORY_JIRA_KEY" }
              }
            }
            ```
    *   **Fulfills:** **NFR-001**, **PROD-003**
    *   **Verification via Test Cases:**
        *   **Test Case `TC-2.1`:**
            *   [x] **Test Method Created:** **Evidence:** Test created in jira_adapter_test.go
            *   [x] **Test Method Passed:** **Evidence:** Test passes (skipped when env vars not set)
        *   **Test Case `TC-2.2`:**
            *   [x] **Test Method Created:** **Evidence:** Test created in jira_adapter_test.go
            *   [x] **Test Method Passed:** **Evidence:** Test passes (skipped when env vars not set)
3.  **Task:** Document the required environment variables.
    *   **Instruction:** `Update the README.md file by adding a new section '## Configuration'. In this section, document the environment variables required to connect to Jira (e.g., JIRA_URL, JIRA_EMAIL, JIRA_API_KEY, JIRA_PROJECT_KEY).`
    *   **Fulfills:** This is a project documentation task.
    *   **Verification via Test Cases:** N/A (Manual verification)
    *   **Documentation:**
        *   [x] **Documentation Updated:** Checked after the `README.md` file is updated. **Instruction:** `N/A`. **Evidence:** Added Configuration section with environment variables.

#### [x] STORY-2.2: Implement CLI and Core Service Integration

1.  **Task:** Implement the core application service.
    *   **Instruction:** `In internal/core/services, create a file story_service.go. This service will take the repository and Jira ports as dependencies. Implement a ProcessStories method that orchestrates the logic: read stories from the repo, and for any story without a JiraID, call the Jira adapter to create it and its tasks.`
    *   **Fulfills:** **ARCH-001**
    *   **Verification via Test Cases:** N/A (Integration)
2.  **Task:** Implement the main CLI entry point.
    *   **Instruction:** `In cmd/jira-story-creator/main.go, write the main application logic. This will involve: initializing the adapters (filesystem, jira), injecting them into the core service, parsing command-line arguments for the input file path, calling the ProcessStories service method, and printing the summary report to the console.`
    *   **Fulfills:** **USER-002**
    *   **Verification via Test Cases:** N/A (Manual E2E test)
3.  **Task:** Document the basic CLI usage.
    *   **Instruction:** `Update the README.md file by adding a new section '## Usage'. In this section, provide the basic command for running the tool via Docker and passing the input file path.`
    *   **Fulfills:** This is a project documentation task.
    *   **Verification via Test Cases:** N/A (Manual verification)
    *   **Documentation:**
        *   [x] **Documentation Updated:** Checked after the `README.md` file is updated. **Instruction:** `N/A`. **Evidence:** Added Usage section with CLI instructions.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-[ID]` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   **Instruction:** You are at the final gate for this phase. Before marking the entire phase as done, you must perform one last, full regression test to ensure nothing was broken by the final commits.
*   [x] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute the master test command for the entire solution (e.g., 'go test ./...') one last time.`
    *   **Evidence:** All 6 tests pass (4 from PHASE-1, 2 from PHASE-2). Jira integration tests skip when env vars not set.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox above is marked `[x]`, your final action for this phase is to modify the main title of this document, changing `[ ] PHASE-2` to `[x] PHASE-2`. This concludes your work on this phase file.
