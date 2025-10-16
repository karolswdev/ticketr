# Field Inheritance Examples

This file demonstrates Ticketr's hierarchical field inheritance feature, where tasks automatically inherit custom fields from their parent ticket.

## Example 1: Complete Inheritance (No Task Fields)

```markdown
# TICKET: E-Commerce Platform Upgrade

## Description
Modernize the e-commerce platform with new features and improved performance.

## Fields
- Priority: High
- Sprint: Sprint 23
- Epic Link: PROJ-500
- Story Points: 21
- Team: Platform

## Acceptance Criteria
- All tasks completed within sprint
- Zero regressions in existing features
- Performance improved by 30%

## Tasks
- Upgrade database schema
- Implement caching layer
- Optimize API endpoints
- Update frontend components
```

**Result:** All four tasks will be created in Jira with ALL parent fields:
- Priority: High
- Sprint: Sprint 23
- Epic Link: PROJ-500
- Story Points: 21
- Team: Platform

---

## Example 2: Partial Override (Task-Specific Fields)

```markdown
# TICKET: Authentication System Hardening

## Description
Improve security and resilience of the authentication system.

## Fields
- Priority: High
- Sprint: Sprint 24
- Component: Security
- Story Points: 13

## Tasks
- ### Implement rate limiting
  #### Fields
  - Priority: Critical
  - Assignee: security-team@company.com

- ### Add MFA support
  #### Fields
  - Story Points: 8
  - Labels: [security, mfa]

- ### Update session management
  # No custom fields - inherits everything
```

**Result:**

**Task 1 - "Implement rate limiting":**
- Priority: Critical (overridden)
- Sprint: Sprint 24 (inherited)
- Component: Security (inherited)
- Story Points: 13 (inherited)
- Assignee: security-team@company.com (task-specific)

**Task 2 - "Add MFA support":**
- Priority: High (inherited)
- Sprint: Sprint 24 (inherited)
- Component: Security (inherited)
- Story Points: 8 (overridden)
- Labels: [security, mfa] (task-specific)

**Task 3 - "Update session management":**
- Priority: High (inherited)
- Sprint: Sprint 24 (inherited)
- Component: Security (inherited)
- Story Points: 13 (inherited)

---

## Example 3: Multiple Field Types

```markdown
# TICKET: Payment Gateway Integration

## Description
Integrate Stripe payment gateway for subscription billing.

## Fields
- Priority: High
- Sprint: Sprint 25
- Component: Backend
- Story Points: 21
- Epic Link: PROJ-600
- Labels: [payment, integration]
- Affects Version: 2.0.0

## Tasks
- ### Set up Stripe API integration
  #### Fields
  - Priority: Critical
  - Assignee: backend-team@company.com
  - Story Points: 8

- ### Implement webhook handlers
  #### Fields
  - Component: Infrastructure
  - Labels: [webhooks, async]

- ### Add subscription management UI
  #### Fields
  - Component: Frontend
  - Assignee: frontend-team@company.com
  - Story Points: 5

- ### Write integration tests
  # Inherits all parent fields
```

**Result:**

**Task 1:**
- Priority: Critical (override)
- Sprint: Sprint 25 (inherit)
- Component: Backend (inherit)
- Story Points: 8 (override)
- Epic Link: PROJ-600 (inherit)
- Labels: [payment, integration] (inherit)
- Affects Version: 2.0.0 (inherit)
- Assignee: backend-team@company.com (task-specific)

**Task 2:**
- Priority: High (inherit)
- Sprint: Sprint 25 (inherit)
- Component: Infrastructure (override)
- Story Points: 21 (inherit)
- Epic Link: PROJ-600 (inherit)
- Labels: [webhooks, async] (override)
- Affects Version: 2.0.0 (inherit)

**Task 3:**
- Priority: High (inherit)
- Sprint: Sprint 25 (inherit)
- Component: Frontend (override)
- Story Points: 5 (override)
- Epic Link: PROJ-600 (inherit)
- Labels: [payment, integration] (inherit)
- Affects Version: 2.0.0 (inherit)
- Assignee: frontend-team@company.com (task-specific)

**Task 4:**
- (Inherits ALL parent fields exactly as defined)

---

## Example 4: Real-World Sprint Planning

```markdown
# TICKET: User Dashboard Feature

## Description
Implement comprehensive user dashboard with analytics and customization.

## Fields
- Priority: Medium
- Sprint: Sprint 26
- Component: Frontend
- Story Points: 21
- Epic Link: PROJ-700
- Team: Product

## Acceptance Criteria
- Dashboard loads in under 2 seconds
- All widgets are customizable
- Data updates in real-time

## Tasks
- ### Design dashboard wireframes
  #### Description
  Create wireframes for all dashboard layouts and widget types.

  #### Acceptance Criteria
  - Wireframes approved by product team
  - Mobile and desktop versions

  #### Fields
  - Priority: High
  - Assignee: design-team@company.com
  - Story Points: 5

- ### Implement dashboard API
  #### Description
  Build RESTful API endpoints for dashboard data.

  #### Fields
  - Component: Backend
  - Assignee: backend-dev@company.com
  - Story Points: 8

- ### Build widget framework
  #### Description
  Create reusable widget framework with drag-and-drop support.

  #### Fields
  - Assignee: frontend-dev@company.com
  - Labels: [framework, reusable]
  - Story Points: 8

- ### Add analytics integration
  #### Description
  Integrate Google Analytics and custom event tracking.

  #### Fields
  - Component: Analytics
  - Priority: Low
  - Story Points: 3
```

**Result:**

**Task 1 - Design:**
- Priority: High (override)
- Component: Frontend (inherit)
- Sprint: Sprint 26 (inherit)
- Epic Link: PROJ-700 (inherit)
- Team: Product (inherit)
- Story Points: 5 (override)
- Assignee: design-team@company.com (task-specific)

**Task 2 - API:**
- Priority: Medium (inherit)
- Component: Backend (override)
- Sprint: Sprint 26 (inherit)
- Epic Link: PROJ-700 (inherit)
- Team: Product (inherit)
- Story Points: 8 (override)
- Assignee: backend-dev@company.com (task-specific)

**Task 3 - Framework:**
- Priority: Medium (inherit)
- Component: Frontend (inherit)
- Sprint: Sprint 26 (inherit)
- Epic Link: PROJ-700 (inherit)
- Team: Product (inherit)
- Story Points: 8 (override)
- Assignee: frontend-dev@company.com (task-specific)
- Labels: [framework, reusable] (task-specific)

**Task 4 - Analytics:**
- Priority: Low (override)
- Component: Analytics (override)
- Sprint: Sprint 26 (inherit)
- Epic Link: PROJ-700 (inherit)
- Team: Product (inherit)
- Story Points: 3 (override)

---

## Key Takeaways

1. **Inheritance is Automatic**: No configuration needed - just define parent fields
2. **Flexible Overrides**: Override any field at the task level as needed
3. **Consistent Defaults**: Parent fields provide sensible defaults for all tasks
4. **Granular Control**: Mix inherited and task-specific fields freely
5. **Sprint Planning**: Perfect for maintaining consistency across sprint tasks

## Technical References

- Implementation: `internal/core/services/ticket_service.go` (calculateFinalFields)
- Test Coverage: TC-701.1, TC-701.2, TC-701.3, TC-701.4
- Requirements: PROD-009, PROD-202
- Documentation: README.md "Field Inheritance" section
