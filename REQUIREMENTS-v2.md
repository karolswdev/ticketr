# Ticketr - Software Requirements Specification

**Version:** 2.0
**Status:** Modernization Baseline

## Introduction

This document outlines the software requirements for **Ticketr v2.0**. It serves as the single source of truth for what the system must do, the constraints under which it must operate, and the rules governing its development and deployment.

Each requirement has a **unique, stable ID** (e.g., `PROD-001`). These IDs **MUST** be used to link implementation stories and test cases back to these foundational requirements, ensuring complete traceability.

The requirement keywords (`MUST`, `MUST NOT`, `SHOULD`, `SHOULD NOT`, `MAY`) are used as defined in [RFC 2119](https://www.ietf.org/rfc/rfc2119.txt).

---

## 1. Product & Functional Requirements

*Defines what the system does; its core features and capabilities.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="PROD-001"></a>**PROD-001** | **Core:** File-Based State | The system **MUST** operate on an input Markdown file and produce a new Markdown file as its primary output. This output file **MUST** be a valid input file for subsequent runs, with Jira Issue Keys injected for all successfully processed items. Files **MUST** use the `# TICKET:` schema format. | This creates a self-contained, idempotent workflow that supports easy recovery from failures and iterative updates. The file itself becomes the record of the system's state. |
| <a name="PROD-002"></a>**PROD-002** | Hierarchical Validation | The system **MUST** validate input hierarchically. If a ticket's definition is malformed, the ticket and all its child tasks **MUST** be rejected as a single unit. | To prevent the creation of orphaned tasks or incomplete tickets in Jira, ensuring data integrity. |
| <a name="PROD-003"></a>**PROD-003** | Rich Create & Update Logic | The system **MUST** support both "create" and "update" operations for tickets and tasks, triggered by the presence of a Jira Issue Key (e.g., `[PROJ-456]`). When updating a ticket, the system **MUST** update its fields in Jira and **MUST** also process its child tasks (creating or updating them as required). | To provide a comprehensive, "single source of truth" workflow where the Markdown file's state can be fully synchronized with Jira. |
| <a name="PROD-004"></a>**PROD-004** | Human-Readable Logging | The system **MUST** write a summary of its execution to a local, human-readable flat text log file, in addition to the console. | To provide a persistent and easily reviewable record of the application's execution for auditing and debugging purposes. |
| <a name="PROD-005"></a>**PROD-005** | Rich Task Definitions | The system **MUST** be able to parse and sync a `Description` and `Acceptance Criteria` for individual tasks, in addition to their titles. Tasks **MUST** support custom `## Fields` sections for field overrides. | To allow for the creation of detailed, well-defined tasks that do not require immediate follow-up in the Jira UI. |
| <a name="PROD-009"></a>**PROD-009** | Hierarchical Field Inheritance | The system **MUST** implement field inheritance where child tasks inherit custom fields from their parent ticket, with task-specific fields overriding inherited values. | To enable consistent field management while allowing task-specific customizations. |
| <a name="PROD-010"></a>**PROD-010** | Query-Based Pull Synchronization | The system **MUST** support pulling tickets from Jira based on project, epic, or custom JQL queries, converting them to the canonical Markdown format. | To enable bidirectional synchronization and maintain Markdown as the source of truth. |
| <a name="PROD-201"></a>**PROD-201** | Generic `TICKET` Markdown Schema | The system **MUST** recognize and parse `# TICKET:` blocks with structured `## Description`, `## Fields`, `## Acceptance Criteria`, and `## Tasks` sections. | To support any Jira issue type and custom field configuration. |
| <a name="PROD-202"></a>**PROD-202** | Hierarchical Field Inheritance Logic | The system **MUST** calculate final fields for tasks by merging task-specific fields over parent ticket fields. | To ensure consistent field inheritance while allowing task-level overrides. |
| <a name="PROD-203"></a>**PROD-203** | Dynamic Field Mapping | The system **MUST** support configurable field mappings between human-readable names and JIRA field IDs, with automatic type conversion for number and array fields. | To support different JIRA configurations and custom fields across instances. |
| <a name="PROD-204"></a>**PROD-204** | State-Based Change Detection | The system **MUST** track content hashes of tickets to skip unchanged items during push operations. | To minimize API calls and improve performance for large ticket sets. |
| <a name="PROD-205"></a>**PROD-205** | Query-Based Ticket Pulling | The system **MUST** construct JQL queries combining project filters with user-provided JQL for flexible ticket retrieval. | To enable targeted synchronization of specific ticket subsets. |
| <a name="PROD-206"></a>**PROD-206** | Markdown Rendering | The system **MUST** convert JIRA tickets to well-formed Markdown documents preserving all field mappings and hierarchy. | To maintain bidirectional format consistency. |

---

## 2. User Interaction Requirements

*Defines how a user interacts with the system. Focuses on usability and user-facing workflows.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="USER-001"></a>**USER-001** | Non-Interactive Error Handling | The system **MUST** process all tickets and provide a comprehensive summary report, exiting with an error code if any failures occurred. The system **MUST NOT** stop processing on the first error, but rather continue attempting all tickets and aggregate the results. | To ensure all valid tickets are processed even if some fail, providing complete visibility into what succeeded and what failed. |
| <a name="USER-002"></a>**USER-002** | Detailed Execution Report | The final summary report **MUST** include: 1. A list of all successfully created/updated items with direct links to them in the Jira UI. 2. A list of all failed items with a clear reason for failure and their source line number. | To provide a clear, actionable summary of the outcome, enabling users to verify results and debug failures efficiently. |
| <a name="USER-003"></a>**USER-003** | Verbose Output Mode | The system **SHOULD** provide a command-line flag (e.g., `--verbose`) that prints a detailed, real-time log of all operations to the console. | To allow users to monitor the application's progress in detail during execution for debugging or observation. |
| <a name="USER-004"></a>**USER-004** | Jira Schema Discovery | The system **MUST** provide a `schema` command that discovers and generates field mappings from the connected Jira instance. | To automate configuration and support different Jira instances with varying custom fields. |
| <a name="USER-005"></a>**USER-005** | Configurable Pull Verbosity | The system **MUST** allow users to configure which fields are pulled from Jira through the configuration file. | To reduce noise and focus on relevant fields during synchronization. |
| <a name="USER-201"></a>**USER-201** | Centralized YAML Configuration & CLI | The system **MUST** read configuration from `.ticketr.yaml` files using Viper, supporting environment variable overrides. Commands **MUST** be structured as `ticketr push`, `ticketr pull`, and `ticketr schema`. | To provide flexible, maintainable configuration and intuitive command structure. |

---

## 3. Architectural Requirements

*Defines high-level, non-negotiable design principles and structural constraints.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="ARCH-001"></a>**ARCH-001** | Ports & Adapters | The system architecture **MUST** adhere to the Ports and Adapters (Hexagonal) pattern for its core components. | To decouple the core application logic from external concerns (like the Jira API or the command-line interface), improving testability, maintainability, and flexibility. |

---

## 4. Non-Functional Requirements (NFRs)

*Defines the quality attributes and operational characteristics of the system. The "-ilities".*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="NFR-001"></a>**NFR-001** | **Security:** Flexible Credentials | The system **MUST** support Jira credentials from environment variables or `.ticketr.yaml` configuration, with environment variables taking precedence. | To provide flexible secret management suitable for both local development and containerized deployments. |
| <a name="NFR-002"></a>**NFR-002** | **Reliability:** Graceful API Error Handling | The system **MUST** gracefully handle and report API errors from Jira related to user permissions and project-specific validation rules. | To provide clear, actionable feedback to the user when an API call fails, enabling them to diagnose and resolve the underlying issue in Jira. |
| <a name="NFR-201"></a>**NFR-201** | **Final SRS Conformance** | The system **MUST** implement hierarchical validation, file-based logging, and enhanced reporting as specified in the modernization plan section 4. | To complete all v2.0 requirements including validation services, comprehensive logging, and detailed execution reports. |

---

## 5. Technology & Platform Requirements

*Defines the specific technologies, frameworks, and platforms that are mandated for use.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="TECH-P-001"></a>**TECH-P-001** | **Primary Language:** Go | The application's backend **MUST** be implemented using Go. | To build a performant, single-binary executable that is well-suited for containerized, command-line applications. |

---

## 6. Operational & DevOps Requirements

*Defines the rules and constraints governing the development workflow, build process, and execution environment.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="DEV-001"></a>**DEV-001** | **Containerized Execution** | The application **MUST** be distributed and executed as a Docker container. | To ensure a consistent, portable, and isolated execution environment, abstracting away the host machine's configuration. |