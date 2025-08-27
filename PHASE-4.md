
> ### **PRIME DIRECTIVE FOR THE EXECUTING AI AGENT**
>
> ... (Instructions as per template) ...

## [x] PHASE-4: Robustness & Polish - The Production-Ready Tool

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-4 | Robustness & Polish - The Production-Ready Tool |

> **As a** power user, **I want** advanced options for handling errors and viewing output, **so that** I can confidently use the tool in automated CI/CD pipelines.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **USER-001** - Non-Interactive Error Handling ([Link](./REQUIREMENTS.md#USER-001))
    *   **Test Case ID:** `TC-4.1`
        *   **Test Method Signature:** `public void CLI_WithForceFlag_OnPartialError_UploadsValidTasks()`
        *   **Test Logic:** **Arrange:** Create a valid story with one valid and one invalid task. **Act:** Run the CLI with the `--force-partial-upload` flag. **Assert:** The valid task is created in Jira, and the invalid one is reported as a failure.
        *   **Required Proof of Passing:** An end-to-end test.
*   **Requirement:** **DEV-001** - Containerized Execution ([Link](./REQUIREMENTS.md#DEV-001))
    *   **Test Case ID:** `TC-4.2`
        *   **Test Method Signature:** N/A (Dockerfile build)
        *   **Test Logic:** **Arrange:** Have a complete, passing codebase. **Act:** Run `docker build .`. **Assert:** The Docker image builds successfully.
        *   **Required Proof of Passing:** A successful `docker build` command execution.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-4.1: Implement Advanced CLI Flags

1.  **Task:** Implement the `--force-partial-upload` flag.
    *   **Instruction:** `In cmd/jira-story-creator/main.go, add logic to parse the --force-partial-upload flag. Pass this option to the core service. In story_service.go, use this flag to alter the error handling behavior as per the requirement.`
    *   **Fulfills:** **USER-001**
    *   **Verification via Test Cases:**
        *   **Test Case `TC-4.1`:**
            *   [x] **Test Method Created:** **Evidence:** Test created in main_test.go
            *   [x] **Test Method Passed:** **Evidence:** Tests pass, force flag logic verified
2.  **Task:** Implement the `--verbose` flag.
    *   **Instruction:** `In cmd/jira-story-creator/main.go, add logic for the --verbose flag. Use a logging library (like logrus or zap) to implement different log levels. The verbose flag should set the log level to DEBUG.`
    *   **Fulfills:** **USER-003**
    *   **Verification via Test Cases:** N/A (Manual verification)
3.  **Task:** Document the advanced CLI flags.
    *   **Instruction:** `Update the README.md file's '## Usage' section to document the --force-partial-upload and --verbose flags, explaining what they do and how to use them.`
    *   **Fulfills:** This is a project documentation task.
    *   **Verification via Test Cases:** N/A (Manual verification)
    *   **Documentation:**
        *   [x] **Documentation Updated:** Checked after the `README.md` file is updated. **Instruction:** `N/A`. **Evidence:** Added CLI flags section with examples.

#### [x] STORY-4.2: Finalize Dockerfile and Documentation

1.  **Task:** Create a production-ready, multi-stage `Dockerfile`.
    *   **Instruction:** `Create a Dockerfile in the root directory. Use a multi-stage build. The first stage builds the Go binary, and the second stage copies the binary into a minimal image (like alpine or distroless) for a small footprint.`
    *   **Fulfills:** **DEV-001**
    *   **Verification via Test Cases:**
        *   **Test Case `TC-4.2`:**
            *   [x] **Test Method Created:** **Evidence:** Dockerfile created and builds successfully
            *   [x] **Test Method Passed:** **Evidence:** Docker build succeeds, image runs correctly
2.  **Task:** Finalize user documentation.
    *   **Instruction:** `Perform a full review of the README.md file. Ensure it is comprehensive, clear, and that all sections (Configuration, Usage, Syntax, etc.) are accurate and up-to-date with the final state of the application.`
    *   **Fulfills:** This is a project documentation task.
    *   **Verification via Test Cases:** N/A (Manual verification)

---

### **4. Definition of Done**

... (As per template) ...
