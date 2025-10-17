# Pull with Subtasks - Example

This example demonstrates pulling tickets with subtasks from Jira.

## Command

```bash
ticketr pull --project DEMO --jql "sprint = 'Sprint 1'" --output sprint-tickets.md
```

## Output (sprint-tickets.md)

```markdown
# TICKET: [DEMO-100] Setup development environment

## Description
Configure local development environment with database, cache, and API server.

## Fields
Type: Story
Sprint: Sprint 1
Story Points: 13
Priority: High
Assignee: Jane Developer

## Acceptance Criteria
- Docker Compose setup with all services
- README with setup instructions
- All services start with single command

## Tasks
- [DEMO-101] Configure PostgreSQL database
  ## Description
  Setup PostgreSQL container with initial schema and migrations.

  ## Fields
  Type: Sub-task
  Story Points: 5
  Assignee: Jane Developer

  ## Acceptance Criteria
  - Database initializes with schema.sql
  - Connection pooling configured

- [DEMO-102] Setup Redis cache
  ## Description
  Add Redis container for session caching and rate limiting.

  ## Fields
  Type: Sub-task
  Story Points: 3
  Assignee: John Engineer

  ## Acceptance Criteria
  - Redis container starts on port 6379
  - Connection tested in health check

- [DEMO-103] Create Docker Compose configuration
  ## Description
  Write docker-compose.yml orchestrating all services.

  ## Fields
  Type: Sub-task
  Story Points: 5
  Assignee: Jane Developer

  ## Acceptance Criteria
  - Single `docker-compose up` starts all services
  - Volume mounts for database persistence
  - Health checks for all containers

---

# TICKET: [DEMO-110] Implement API rate limiting

## Description
Add rate limiting middleware to prevent API abuse.

## Fields
Type: Story
Sprint: Sprint 1
Story Points: 8
Priority: Medium
Assignee: John Engineer

## Acceptance Criteria
- Rate limit: 100 requests per minute per IP
- Returns 429 status when exceeded
- Rate limit headers in response

## Tasks
- [DEMO-111] Implement rate limit middleware
  ## Description
  Create Express middleware using Redis for rate counting.

  ## Fields
  Type: Sub-task
  Story Points: 5

  ## Acceptance Criteria
  - Uses Redis INCR with TTL
  - Configurable limits via environment variables

- [DEMO-112] Add rate limit tests
  ## Description
  Write integration tests for rate limiting behavior.

  ## Fields
  Type: Sub-task
  Story Points: 3

  ## Acceptance Criteria
  - Tests verify 100 req/min limit
  - Tests verify 429 response
  - Tests verify rate limit headers
```

## Round-Trip Workflow

1. **Pull tickets** with subtasks from Jira
2. **Edit locally** - modify descriptions, add acceptance criteria
3. **Push changes** back to Jira
4. **Pull again** to verify synchronization

This enables Markdown as the single source of truth for ticket definitions.

## Key Features Demonstrated

### Automatic Subtask Fetching
- When pulling parent tickets, subtasks are automatically included
- No additional configuration required
- Complete hierarchy in a single pull operation

### Rich Task Details
- Each subtask includes Description, Fields, and Acceptance Criteria
- All custom fields are pulled from Jira
- Human-readable field names (e.g., "Story Points" instead of "customfield_10020")

### Field Inheritance
- When pushing, tasks inherit parent fields by default
- Task-specific fields override inherited values
- During pull, exact Jira values are retrieved

### Round-Trip Compatibility
- Pull → Edit → Push → Pull cycle preserves all data
- Field mappings consistent in both directions
- Markdown serves as the canonical source of truth

## Common Use Cases

### Sprint Planning
Pull all tickets for a sprint, review them as Markdown, then push updates back to Jira:

```bash
# Pull sprint backlog
ticketr pull --project DEMO --jql "sprint = 'Sprint 1'" -o sprint-1.md

# Edit sprint-1.md (add acceptance criteria, update estimates)
vim sprint-1.md

# Push changes back
ticketr push sprint-1.md
```

### Epic Management
Pull all tickets in an epic with their subtasks:

```bash
ticketr pull --epic DEMO-100 -o epic-demo-100.md
```

### Team Review
Pull tickets assigned to your team for offline review:

```bash
ticketr pull --jql "assignee in (jane.developer, john.engineer)" -o team-tickets.md
```

## Technical Notes

### Subtask Query
Ticketr uses JQL queries to fetch subtasks:
```
parent = "DEMO-100"
```

### Field Mapping
Custom fields are automatically mapped using reverse field mapping:
- `customfield_10020` → `Story Points`
- `customfield_10021` → `Sprint`
- Field mappings discovered via `ticketr schema` command

### Error Handling
If subtasks cannot be fetched (e.g., permissions), parent tickets are still returned:
- Non-fatal subtask fetch errors logged
- Parent ticket processing continues
- Partial data better than no data

## See Also

- [README.md - Pull Command](../README.md#pull-command)
- [README.md - Field Inheritance](../README.md#field-inheritance)
- [Requirements - PROD-010](../docs/development/REQUIREMENTS.md#prod-010)
