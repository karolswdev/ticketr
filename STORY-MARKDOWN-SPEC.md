# Ticktr Markdown Specification (TMS)

**Version:** 1.0
**Status:** Baseline

## 1. Introduction

This document provides the official Interface Control Document (ICD) for the Ticktr Markdown Syntax. All input files provided to the `ticktr` application **MUST** adhere to this specification.

## 2. File Structure

A Ticktr Markdown file consists of one or more `Story` blocks. Each `Story` block is separated by a thematic break (`---`).

## 3. Story Block

A `Story` block defines a single Jira Story and its associated tasks.

### 3.1. Story Title & ID

*   **Syntax:** The story **MUST** begin with a Level 1 ATX Heading (`#`).
*   **Content:** The heading text **MUST** start with the prefix `STORY:`. The text following the prefix is the story's title.
*   **Update Mode:** To specify an existing Jira story for an update operation, the Jira Issue Key **MUST** be included in brackets at the end of the heading line (e.g., `[PROJ-123]`).

**Example:**
```markdown
# STORY: Implement the login page [TICK-42]
```

### 3.2. Story Description

*   **Syntax:** All text following the Story Title heading and before the first Level 2 heading (`##`) is considered the story's description.

### 3.3. Sections

A story can contain the following sections, denoted by Level 2 ATX Headings (`##`).

#### 3.3.1. `## Acceptance Criteria`

*   **Syntax:** This section is optional. All content following this heading until the next heading is considered the acceptance criteria.
*   **Content:** The content is typically a bulleted list.

#### 3.3.2. `## Tasks`

*   **Syntax:** This section is required if you wish to define tasks for the story.

## 4. Task Block

A `Task` is defined within the `## Tasks` section of a `Story` block.

### 4.1. Task Title & ID

*   **Syntax:** Each task **MUST** be a bulleted list item (using `-` or `*`).
*   **Update Mode:** To specify an existing Jira task for an update, the Jira Issue Key **MUST** be the first element on the line, enclosed in brackets (e.g., `[PROJ-124]`).

**Example:**
```markdown
- [TICK-43] Design the login form fields
- Create the submit button component
```

### 4.2. Task Details

*   **Syntax:** Task details (`Description` and `Acceptance Criteria`) **MUST** be defined in an indented block directly under the parent task item.
*   **Keywords:** The keywords `**Description:**` and `**Acceptance Criteria:**` are used to denote the different detail types.

**Example:**
```markdown
- [TICK-43] Design the login form fields
  - **Description:** This task covers the visual design and layout of the email and password fields.
  - **Acceptance Criteria:**
    - Field for email is present.
    - Field for password is present and masks input.
```

## 5. Complete Example

```markdown
# STORY: Implement the login page [TICK-42]
As a user, I want to be able to log in to the application so that I can access my account.

## Acceptance Criteria
- I can enter my email and password.
- I can click a "Login" button.
- If my credentials are valid, I am taken to the dashboard.
- If my credentials are invalid, I see an error message.

## Tasks
- [TICK-43] Design the login form fields
  - **Description:** This task covers the visual design and layout of the email and password fields.
  - **Acceptance Criteria:**
    - Field for email is present.
    - Field for password is present and masks input.
- Create the submit button component
- Implement the login API call

---

# STORY: Set up the user database
This is a new story to be created.

## Tasks
- Choose a database technology
- Define the user table schema
```