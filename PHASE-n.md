
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

## [ ] PHASE-[ID]: [Phase Title]

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-[ID] | [Phase Title] |

> **As a** [User Persona], **I want** [The Goal of this specific Phase], **so that** [The Value or milestone achieved by completing this Phase].

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

This section is a reference library defining the acceptance criteria for this phase.

*   **Requirement:** **[REQ-ID]** - [Requirement Description] ([Link to file](./REQUIREMENTS.md#[requirement-anchor]))
    *   **Test Case ID:** `TC-[ID.1]`
        *   **Test Method Signature:** `[e.g., public void ClassName_MethodName_Scenario_ExpectedBehavior()]`
        *   **Test Logic:** [Describe the Arrange-Act-Assert steps for this specific test case.]
        *   **Required Proof of Passing:** [Define the exact evidence required for this test case.]
    *   **Test Case ID:** `TC-[ID.2]`
        *   **Test Method Signature:** `[...]`
        *   **Test Logic:** `[...]`
        *   **Required Proof of Passing:** `[...]`

---

### **3. Implementation Plan (The Execution)**

#### [ ] STORY-[ID.1]: [Story Title, e.g., Implement Core Parsing Logic]

1.  **Task:** [Short description of the first atomic task].
    *   **Instruction:** `[Imperative command for the LLM to execute.]`
    *   **Fulfills:** This task contributes to requirement **[REQ-ID]**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-[ID.1]`:**
            *   [ ] **Test Method Created:** Checked after the test method is written. **Evidence:** Provide the complete code for the test method.
            *   [ ] **Test Method Passed:** Checked after the test passes. **Evidence:** Provide the console output from the test runner proving the specific test passed.
        *   **Test Case `TC-[ID.2]`:**
            *   [ ] **Test Method Created:** **Evidence:** `[...]`
            *   [ ] **Test Method Passed:** **Evidence:** `[...]`
    *   **Documentation:**
        *   [ ] **Documentation Updated:** Checked after the relevant documentation is updated. **Instruction:** `Update the [DOCUMENT_NAME.md] file...` **Evidence:** Provide a diff or quote of the updated section.

---
> ### **Story Completion: STORY-[ID.1]**
>
> You may only proceed once all checkboxes for all tasks within this story are marked `[x]`. Then, you **MUST** complete the following steps in order:
>
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** Checked after running all tests created in the project up to this point.
>     *   **Instruction:** `Execute the master test command for the entire solution (e.g., 'dotnet test').`
>     *   **Evidence:** Provide the full summary output from the test runner, showing the total number of tests executed and confirming all have passed.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** Checked after creating the Git commit.
>     *   **Instruction:** `Execute 'git add .' followed by 'git commit -m "feat(story): Complete STORY-[ID.1] - [Story Title]"'.`
>     *   **Evidence:** Provide the full commit hash returned by the Git command.
> 3.  **Finalize Story:**
>     *   **Instruction:** Once the two checkboxes above are complete, you **MUST** update this story's main checkbox from `[ ]` to `[x]`.

---

#### [ ] STORY-[ID.2]: [Another Story Title]
... (Tasks for Story 2 would follow the same structure) ...

> ### **Story Completion: STORY-[ID.2]**
> ... (Completion block for Story 2 would follow the same structure) ...

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-[ID]` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   **Instruction:** You are at the final gate for this phase. Before marking the entire phase as done, you must perform one last, full regression test to ensure nothing was broken by the final commits.
*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute the master test command for the entire solution (e.g., 'dotnet test') one last time.`
    *   **Evidence:** Provide the full, final summary output from the test runner, showing the grand total of tests for this phase and confirming that 100% have passed.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox above is marked `[x]`, your final action for this phase is to modify the main title of this document, changing `[ ] PHASE-[ID]` to `[x] PHASE-[ID]`. This concludes your work on this phase file.