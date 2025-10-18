# Ticketr v3.0 Briefing for External AI Consultation

## Executive Summary

Ticketr is a command-line tool that bridges Markdown files and JIRA, enabling bidirectional synchronization of tickets. Currently at v2.0, we're planning a major architectural evolution to v3.0 that transforms it from a directory-bound tool into a global work management platform.

## Current State (v2.0)

### What Ticketr Does Today
- **Markdown → JIRA**: Write tickets in Markdown, push to JIRA
- **JIRA → Markdown**: Pull tickets from JIRA into Markdown files
- **Bidirectional Sync**: Detect changes, handle conflicts
- **State Management**: Track changes via `.ticketr.state` files
- **Field Inheritance**: Tasks inherit parent ticket fields

### Current Architecture
```
Hexagonal (Ports & Adapters) Pattern:
- Core Domain: Ticket, Task models
- Ports: Repository, JiraClient interfaces
- Adapters: Filesystem, JIRA API
- CLI: Cobra-based commands
```

### Current Limitations
1. **No Hierarchy**: Only 2 levels (Ticket → Task), missing Epic level
2. **Single Project**: Hardcoded `JIRA_PROJECT_KEY` environment variable
3. **Directory-Bound**: Must be in project directory to use
4. **File-Based State**: JSON files, no database, doesn't scale
5. **No Pagination**: Limited to 100 tickets per pull
6. **CLI Only**: No interactive interface

## Proposed Evolution (v3.0)

### Vision
Transform Ticketr from a "Markdown-to-JIRA sync tool" into a "Distributed Work Management Platform" that:
- Works globally from any directory
- Manages multiple JIRA instances/projects
- Provides a rich TUI (Terminal User Interface)
- Supports full Epic → Story → Task hierarchy
- Scales to thousands of tickets

### Key Architectural Changes

#### 1. Project-Based Model
```yaml
Projects (like Kubernetes namespaces):
- Each project = separate JIRA instance/credentials
- Isolated contexts with own configurations
- Markdown files stay in project directories
- Central SQLite tracks everything
```

#### 2. Global Installation
```
~/.config/ticketr/       # Configuration
├── ticketr.db           # Central SQLite database
├── projects/            # Project configurations
└── templates/           # Reusable templates

/usr/local/bin/ticketr   # Single binary in PATH
```

#### 3. Full Hierarchy Support
```markdown
# EPIC: Authentication System
## STORY: User Registration
### TASK: Create API endpoint
### TASK: Design UI form
## STORY: Login Flow
### TASK: Implement JWT
```

#### 4. Terminal User Interface (TUI)
- Interactive tree view of tickets
- Real-time sync status
- Conflict resolution UI
- Vim-style keybindings
- Inspired by k9s, lazygit, htop

### Implementation Roadmap

**5 Phases over 20 weeks:**

1. **Foundation (Weeks 1-4)**: SQLite database layer
2. **Projects (Weeks 5-8)**: Multi-project/multi-JIRA support
3. **Global Tool (Weeks 9-10)**: System-wide installation
4. **TUI (Weeks 11-16)**: Terminal user interface
5. **Advanced (Weeks 17-20)**: Templates, bulk ops, smart sync

### Self-Orchestration via Agents

We have specialized AI agents that will execute the implementation:
- **Director**: Orchestrates the roadmap
- **Builder**: Implements code changes
- **Verifier**: Tests and validates
- **Scribe**: Documents everything
- **Steward**: Architecture oversight

Each phase has:
- Clear deliverables
- Specific file outputs
- Test requirements
- Documentation needs
- Gate criteria

## Questions for External Perspective

### Strategic Questions

1. **Architecture**: Is the project-based model (similar to kubectl contexts) the right abstraction? Should we consider alternatives like workspaces or profiles?

2. **Database Choice**: SQLite for local state - correct choice? Should we consider:
   - Embedded key-value stores (BoltDB, BadgerDB)?
   - Just staying with files but using better indexing?
   - Starting with PostgreSQL for future multi-user support?

3. **TUI vs Web UI**: We chose TUI for CLI-first philosophy. Should we:
   - Also build a local web UI (localhost:8080)?
   - Focus solely on TUI?
   - Build both in parallel?

4. **Hierarchy Depth**: Currently planning 3 levels (Epic→Story→Task). Should we:
   - Support arbitrary depth?
   - Stay with fixed 3 levels?
   - Make it configurable per project?

### Technical Questions

5. **Sync Strategy**: How should we handle conflicts?
   - Current: Manual resolution
   - Alternative: Automatic merging with rules
   - Alternative: CRDT-based eventual consistency

6. **Performance**: Target is 100ms for 1000 tickets. For 10,000+ tickets:
   - Pagination everywhere?
   - Lazy loading?
   - Background sync workers?

7. **Plugin System**: Should v3 include a plugin architecture?
   - Pros: Extensibility, community contributions
   - Cons: Complexity, security concerns

### Implementation Questions

8. **Migration Path**: How aggressive should we be?
   - Gradual: Feature flags, slow rollout
   - Big Bang: Clean break, v3 is separate tool
   - Hybrid: v3 can import v2 but not vice versa

9. **Testing Strategy**: Current plan is 80% coverage. Should we:
   - Aim for 100% on critical paths?
   - Focus more on integration tests?
   - Add property-based testing?

10. **Agent Orchestration**: We have AI agents to build this. Should they:
    - Work fully autonomously?
    - Require human approval at gates?
    - Pair program with humans?

## Technical Context

### Technology Stack
- **Language**: Go 1.22+
- **CLI Framework**: Cobra
- **TUI Framework**: tview (considering bubbletea)
- **Database**: SQLite (via database/sql)
- **Testing**: Standard library + testify
- **CI/CD**: GitHub Actions

### File Structure
```
ticketr/
├── cmd/ticketr/         # CLI entry point
├── internal/
│   ├── core/            # Business logic
│   │   ├── domain/      # Models
│   │   ├── ports/       # Interfaces
│   │   └── services/    # Use cases
│   └── adapters/        # External integrations
│       ├── filesystem/
│       ├── jira/
│       ├── database/    # NEW in v3
│       └── tui/         # NEW in v3
```

### Design Patterns
- Hexagonal Architecture (Ports & Adapters)
- Repository Pattern
- Command Pattern (CLI commands)
- Strategy Pattern (sync strategies)

## Specific Areas for Feedback

1. **Database Schema**: Review our planned schema - any obvious issues?

2. **TUI Wireframes**: We have 7 detailed wireframes - are we missing key views?

3. **Project Model**: Is our project isolation model sound?

4. **Self-Orchestration**: Can our agent-based development work?

5. **Performance Goals**: Are our targets realistic?
   - 100ms startup
   - 50ms project switch
   - 100ms for 1000 ticket query
   - 16ms TUI refresh (60fps)

## Risk Assessment

### Technical Risks
- SQLite corruption (mitigation: WAL mode)
- Cross-platform TUI issues (mitigation: extensive testing)
- JIRA API changes (mitigation: versioned adapter)

### User Adoption Risks
- Learning curve for TUI (mitigation: progressive disclosure)
- Migration friction (mitigation: automated tooling)
- Muscle memory disruption (mitigation: v2 command compatibility)

### Development Risks
- 20-week timeline aggressive (mitigation: phased delivery)
- Agent coordination complexity (mitigation: human oversight)
- Feature creep (mitigation: strict phase gates)

## Success Criteria

### Quantitative
- Installation success rate > 95%
- Migration success rate > 99%
- Performance targets met
- Test coverage > 80%
- Zero data loss incidents

### Qualitative
- "Just works" from any directory
- TUI feels as natural as k9s
- Faster than web UI for common tasks
- Community contributions increase

## The Ask

Please review this plan and provide:

1. **Strategic Validation**: Is this the right direction?
2. **Technical Critique**: What are we missing or overthinking?
3. **Risk Assessment**: What will likely go wrong?
4. **Alternative Approaches**: What would you do differently?
5. **Priority Adjustment**: Should we reorder phases?

## Appendix: Links to Key Documents

- Current Architecture: `docs/ARCHITECTURE.md`
- V3 Roadmap: `docs/v3-implementation-roadmap.md`
- Technical Spec: `docs/v3-technical-specification.md`
- Project Model: `docs/v3-project-model.md`
- TUI Wireframes: `docs/tui-wireframes.md`
- Self-Orchestration: `docs/v3-roadmap-orchestration.md`

---

*This briefing prepared for external AI consultation. Please provide candid feedback - we want to build this right.*