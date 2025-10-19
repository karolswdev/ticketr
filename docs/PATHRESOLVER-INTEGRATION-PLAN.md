# PathResolver Integration Plan

## Status: Phase 3 Complete, Integration Pending

**Last Updated**: 2025-01-18
**Phase**: 3 (Global Installation)
**Status**: Core implementation complete, system-wide integration planned

---

## Executive Summary

The PathResolver service has been **successfully implemented and tested** as part of Phase 3 (Global Installation). The service provides cross-platform path resolution following XDG Base Directory specification on Linux/macOS and Windows conventions.

**Current Achievement**:
- ✅ PathResolver service implemented (290 lines)
- ✅ Comprehensive unit tests (450 lines, 92.9% coverage)
- ✅ Integration tests (353 lines, 11/11 passing)
- ✅ Documentation (GLOBAL-INSTALLATION.md, ARCHITECTURE.md)
- ✅ All tests passing, no race conditions

**Next Step**: Integrate PathResolver into existing services (SQLiteAdapter, ConfigLoader, etc.)

---

## Phase 3 Completion Criteria

### ✅ Acceptance Criteria Met

- [x] XDG Base Directory specification compliance (Linux/macOS)
- [x] Windows standard path conventions (APPDATA, LOCALAPPDATA)
- [x] Environment variable overrides (XDG_CONFIG_HOME, etc.)
- [x] Automatic directory creation with proper permissions (0755)
- [x] Cross-platform path resolution
- [x] Comprehensive test coverage (>80%)
- [x] Documentation complete

### 📋 Integration Criteria (Pending)

- [ ] SQLiteAdapter uses PathResolver instead of hardcoded paths
- [ ] ConfigLoader uses PathResolver for config files
- [ ] Workspace commands initialize with PathResolver
- [ ] All adapters use PathResolver for filesystem operations
- [ ] Migration from old `~/.ticketr` directory (if exists)

---

## Integration Architecture

### Current State

```
┌─────────────────────────────────────────────────┐
│              CLI Commands                       │
│         (workspace_commands.go)                 │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│          initWorkspaceService()                 │
│   adapter = NewSQLiteAdapter(hardcodedPath) ❌  │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│           SQLiteAdapter                         │
│      (accepts string path parameter)            │
└─────────────────────────────────────────────────┘
```

### Target State

```
┌─────────────────────────────────────────────────┐
│              CLI Commands                       │
│         (workspace_commands.go)                 │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│          initWorkspaceService()                 │
│   pathResolver = NewPathResolver()       ✅     │
│   adapter = NewSQLiteAdapter(pathResolver) ✅   │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│           SQLiteAdapter                         │
│    (accepts *PathResolver parameter)            │
│    dbPath = pathResolver.DatabasePath()         │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│           PathResolver                          │
│   - ConfigDir(), DataDir(), CacheDir()          │
│   - DatabasePath(), ConfigPath()                │
│   - EnsureDirectories()                         │
└─────────────────────────────────────────────────┘
```

---

## Implementation Steps

### Step 1: Update SQLiteAdapter

**File**: `internal/adapters/database/sqlite_adapter.go`

**Current Signature**:
```go
func NewSQLiteAdapter(dbPath string) (*SQLiteAdapter, error)
```

**Target Signature**:
```go
func NewSQLiteAdapter(pathResolver *services.PathResolver) (*SQLiteAdapter, error)
```

**Changes Required**:

```go
// Add PathResolver field to struct
type SQLiteAdapter struct {
    db           *sql.DB
    pathResolver *services.PathResolver  // NEW
}

// Update constructor
func NewSQLiteAdapter(pathResolver *services.PathResolver) (*SQLiteAdapter, error) {
    if pathResolver == nil {
        return nil, fmt.Errorf("pathResolver cannot be nil")
    }

    // Ensure database directory exists
    if err := pathResolver.EnsureDirectories(); err != nil {
        return nil, fmt.Errorf("failed to create database directory: %w", err)
    }

    // Get database path from PathResolver
    dbPath := pathResolver.DatabasePath()

    // Open database connection
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Rest of initialization...
    adapter := &SQLiteAdapter{
        db:           db,
        pathResolver: pathResolver,  // NEW
    }

    // Initialize schema
    if err := adapter.initSchema(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to initialize schema: %w", err)
    }

    return adapter, nil
}

// Add getter for PathResolver
func (a *SQLiteAdapter) PathResolver() *services.PathResolver {
    return a.pathResolver
}
```

**Testing**:
```go
func TestNewSQLiteAdapter(t *testing.T) {
    tempDir := t.TempDir()
    pr, err := services.NewPathResolverWithOptions("ticketr",
        func(key string) string { return "" },
        func() (string, error) { return tempDir, nil })
    require.NoError(t, err)

    adapter, err := NewSQLiteAdapter(pr)
    require.NoError(t, err)
    defer adapter.Close()

    assert.NotNil(t, adapter.DB())
    assert.NotNil(t, adapter.PathResolver())
}
```

---

### Step 2: Update CLI Commands

**File**: `cmd/ticketr/workspace_commands.go`

**Current Code** (lines ~160-170):
```go
func initWorkspaceService() (*services.WorkspaceService, error) {
    // Load features config
    features, err := config.LoadFeatures()
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    // Initialize database
    adapter, err := database.NewSQLiteAdapter(features.SQLitePath)  // ❌ OLD
    if err != nil {
        return nil, fmt.Errorf("failed to initialize database: %w", err)
    }
    // ...
}
```

**Target Code**:
```go
func initWorkspaceService() (*services.WorkspaceService, error) {
    // Create PathResolver
    pathResolver, err := services.NewPathResolver()
    if err != nil {
        return nil, fmt.Errorf("failed to create path resolver: %w", err)
    }

    // Initialize database with PathResolver
    adapter, err := database.NewSQLiteAdapter(pathResolver)  // ✅ NEW
    if err != nil {
        return nil, fmt.Errorf("failed to initialize database: %w", err)
    }

    // Create workspace repository
    repo := database.NewWorkspaceRepository(adapter.DB())

    // Create credential store
    credStore := keychain.NewKeychainStore()

    // Create and return workspace service
    return services.NewWorkspaceService(repo, credStore), nil
}
```

**Impact**: Removes dependency on `features.SQLitePath` configuration

---

### Step 3: Update ConfigLoader (If Exists)

**File**: `internal/config/loader.go` (if it exists)

**Changes**:
- Accept PathResolver in constructor
- Use `pathResolver.ConfigPath()` for config file location
- Use `pathResolver.WorkspacesPath()` for workspaces file

**Example**:
```go
type ConfigLoader struct {
    pathResolver *services.PathResolver
}

func NewConfigLoader(pathResolver *services.PathResolver) *ConfigLoader {
    return &ConfigLoader{pathResolver: pathResolver}
}

func (cl *ConfigLoader) LoadConfig() (*Config, error) {
    configPath := cl.pathResolver.ConfigPath()

    // Ensure config directory exists
    if err := cl.pathResolver.EnsureDirectories(); err != nil {
        return nil, fmt.Errorf("failed to create config directory: %w", err)
    }

    // Rest of loading logic...
}
```

---

### Step 4: Update All Adapter Tests

**Pattern for Test Updates**:

**Old Pattern**:
```go
func TestSomething(t *testing.T) {
    dbPath := filepath.Join(t.TempDir(), "test.db")
    adapter, err := NewSQLiteAdapter(dbPath)
    // ...
}
```

**New Pattern**:
```go
func TestSomething(t *testing.T) {
    tempDir := t.TempDir()
    pr, err := services.NewPathResolverWithOptions("ticketr",
        func(key string) string { return "" },
        func() (string, error) { return tempDir, nil })
    require.NoError(t, err)

    adapter, err := NewSQLiteAdapter(pr)
    require.NoError(t, err)
    defer adapter.Close()
    // ...
}
```

**Files to Update**:
- `internal/adapters/database/sqlite_adapter_test.go`
- `internal/adapters/database/workspace_repository_test.go`
- Any other tests that create database adapters

---

### Step 5: Migration Support

**Goal**: Detect and optionally migrate from old `~/.ticketr` directory

**Detection**:
```go
// In initWorkspaceService or first-run logic
func checkLegacyInstallation(pathResolver *services.PathResolver) (bool, string, error) {
    homeDir := os.Getenv("HOME")
    if runtime.GOOS == "windows" {
        homeDir = os.Getenv("USERPROFILE")
    }

    legacyDir := filepath.Join(homeDir, ".ticketr")

    // Check if legacy directory exists
    if _, err := os.Stat(legacyDir); os.IsNotExist(err) {
        return false, "", nil
    }

    // Check for legacy database
    legacyDB := filepath.Join(legacyDir, "ticketr.db")
    if _, err := os.Stat(legacyDB); err == nil {
        return true, legacyDB, nil
    }

    return false, "", nil
}
```

**Migration Command** (Future):
```go
var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Migrate from v2.x to v3.x",
    RunE: func(cmd *cobra.Command, args []string) error {
        pr, _ := services.NewPathResolver()

        hasLegacy, legacyPath, _ := checkLegacyInstallation(pr)
        if !hasLegacy {
            fmt.Println("No legacy installation found")
            return nil
        }

        fmt.Printf("Found legacy database at: %s\n", legacyPath)
        fmt.Println("Migration will copy data to new location:")
        fmt.Printf("  New location: %s\n", pr.DatabasePath())

        // Prompt for confirmation
        // Perform migration
        // ...
    },
}
```

---

## Testing Strategy

### Unit Tests

**Coverage Target**: 80%+ for all modified files

**Test Scenarios**:
1. PathResolver creation and initialization
2. Database adapter with PathResolver
3. Directory creation and permissions
4. Path resolution on different platforms
5. Error handling (nil PathResolver, permission denied, etc.)

### Integration Tests

**Test Scenarios**:
1. End-to-end workspace creation with PathResolver
2. Database operations with resolved paths
3. Config loading from resolved paths
4. Cross-platform path verification
5. Migration detection and handling

### Manual Testing Checklist

- [ ] `ticketr workspace create` creates database at correct location
- [ ] Database persists between runs
- [ ] Workspace list shows created workspaces
- [ ] Switch workspace works correctly
- [ ] Delete workspace removes data
- [ ] Works on Linux (XDG paths)
- [ ] Works on macOS (Application Support)
- [ ] Works on Windows (AppData)

---

## Rollout Plan

### Phase 1: Core Integration (Week 1)
- [ ] Update SQLiteAdapter
- [ ] Update workspace_commands.go
- [ ] Update all adapter tests
- [ ] Run full test suite
- [ ] Manual testing on Linux

### Phase 2: Extended Integration (Week 2)
- [ ] Update ConfigLoader (if exists)
- [ ] Update any other file-based adapters
- [ ] Add migration detection
- [ ] Cross-platform testing (macOS, Windows)

### Phase 3: Migration Support (Week 3)
- [ ] Implement migration command
- [ ] Add backup functionality
- [ ] Update documentation
- [ ] Release notes

### Phase 4: Release (Week 4)
- [ ] Final testing
- [ ] Update CHANGELOG
- [ ] Tag release v3.0.0
- [ ] Announce breaking changes

---

## Breaking Changes

### For Users

**Old Behavior** (v2.x):
```
Database location: ~/.ticketr/ticketr.db
Config location:   ~/.ticketr/config.yaml
```

**New Behavior** (v3.x):

**Linux**:
```
Database: ~/.local/share/ticketr/ticketr.db
Config:   ~/.config/ticketr/config.yaml
Cache:    ~/.cache/ticketr/
```

**macOS**:
```
Database: ~/Library/Application Support/ticketr/ticketr.db
Config:   ~/Library/Application Support/ticketr/config.yaml
Cache:    ~/Library/Caches/ticketr/
```

**Windows**:
```
Database: %LOCALAPPDATA%\ticketr\ticketr.db
Config:   %APPDATA%\ticketr\config.yaml
Cache:    %LOCALAPPDATA%\ticketr\cache\
```

**Migration Path**:
1. Run `ticketr migrate` to check for legacy installation
2. Confirm migration
3. Data copied to new location
4. Legacy directory preserved (user can delete manually)

### For Developers

**API Changes**:
```go
// Old
adapter := database.NewSQLiteAdapter("/path/to/db")

// New
pr, _ := services.NewPathResolver()
adapter := database.NewSQLiteAdapter(pr)
```

**Impact**: Medium (affects database initialization, tests)

---

## Risk Assessment

### Technical Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Database migration failures | Medium | High | Add backup before migration |
| Path permission issues | Low | Medium | Proper error messages, fallback options |
| Cross-platform bugs | Medium | High | Comprehensive platform testing |
| Test failures | Low | Low | Comprehensive test suite already exists |

### User Impact Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Lost data during migration | Low | Critical | Automatic backup, clear warnings |
| Confusion about new paths | High | Low | Clear documentation, migration guide |
| Workflow disruption | Medium | Medium | Smooth migration path, v2 compatibility period |

---

## Success Metrics

### Technical Metrics
- [ ] 100% of tests passing
- [ ] Test coverage ≥80%
- [ ] Zero compile errors
- [ ] Zero race conditions
- [ ] All platforms tested

### User Metrics
- [ ] <1% migration failures
- [ ] Clear error messages for all failure scenarios
- [ ] Positive user feedback on new directory structure
- [ ] Zero data loss incidents

---

## Timeline

**Estimated Effort**: 3-4 weeks

| Week | Tasks | Deliverables |
|------|-------|--------------|
| 1 | Core integration | SQLiteAdapter updated, tests passing |
| 2 | Extended integration | All adapters using PathResolver |
| 3 | Migration support | Migration command implemented |
| 4 | Testing & release | v3.0.0 tagged and released |

**Dependencies**:
- Phase 3 PathResolver implementation (✅ Complete)
- Comprehensive test suite (✅ Exists)
- Documentation (✅ Created)

---

## References

### Related Documentation
- [GLOBAL-INSTALLATION.md](GLOBAL-INSTALLATION.md) - Installation guide
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [v3-implementation-roadmap.md](v3-implementation-roadmap.md) - Overall v3 plan

### Code References
- PathResolver: `internal/core/services/path_resolver.go`
- PathResolver Tests: `internal/core/services/path_resolver_test.go`
- Integration Tests: `internal/core/services/path_resolver_integration_test.go`
- SQLiteAdapter: `internal/adapters/database/sqlite_adapter.go`
- Workspace Commands: `cmd/ticketr/workspace_commands.go`

### External Standards
- [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
- [macOS File System Programming Guide](https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/FileSystemProgrammingGuide/FileSystemOverview/FileSystemOverview.html)
- [Windows Known Folders](https://docs.microsoft.com/en-us/windows/win32/shell/known-folders)

---

## Appendix: Code Snippets

### Complete SQLiteAdapter Example

```go
package database

import (
    "database/sql"
    "fmt"

    "github.com/karolswdev/ticktr/internal/core/services"
    _ "github.com/mattn/go-sqlite3"
)

type SQLiteAdapter struct {
    db           *sql.DB
    pathResolver *services.PathResolver
}

func NewSQLiteAdapter(pathResolver *services.PathResolver) (*SQLiteAdapter, error) {
    if pathResolver == nil {
        return nil, fmt.Errorf("pathResolver cannot be nil")
    }

    if err := pathResolver.EnsureDirectories(); err != nil {
        return nil, fmt.Errorf("failed to create database directory: %w", err)
    }

    dbPath := pathResolver.DatabasePath()
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    if err := db.Ping(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    adapter := &SQLiteAdapter{
        db:           db,
        pathResolver: pathResolver,
    }

    if err := adapter.initSchema(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to initialize schema: %w", err)
    }

    return adapter, nil
}

func (a *SQLiteAdapter) DB() *sql.DB {
    return a.db
}

func (a *SQLiteAdapter) PathResolver() *services.PathResolver {
    return a.pathResolver
}

func (a *SQLiteAdapter) Close() error {
    if a.db != nil {
        return a.db.Close()
    }
    return nil
}

func (a *SQLiteAdapter) initSchema() error {
    schema := `
    CREATE TABLE IF NOT EXISTS workspaces (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        jira_url TEXT,
        project_key TEXT,
        credential_ref TEXT,
        is_default INTEGER NOT NULL DEFAULT 0,
        last_used INTEGER,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
    );

    CREATE INDEX IF NOT EXISTS idx_workspaces_name ON workspaces(name);
    CREATE INDEX IF NOT EXISTS idx_workspaces_default ON workspaces(is_default);
    `

    _, err := a.db.Exec(schema)
    return err
}
```

### Complete Test Example

```go
package database

import (
    "testing"

    "github.com/karolswdev/ticktr/internal/core/services"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupTestAdapter(t *testing.T) *SQLiteAdapter {
    tempDir := t.TempDir()
    pr, err := services.NewPathResolverWithOptions("ticketr",
        func(key string) string { return "" },
        func() (string, error) { return tempDir, nil })
    require.NoError(t, err)

    adapter, err := NewSQLiteAdapter(pr)
    require.NoError(t, err)

    t.Cleanup(func() {
        adapter.Close()
    })

    return adapter
}

func TestSQLiteAdapter_Integration(t *testing.T) {
    adapter := setupTestAdapter(t)

    // Verify database is accessible
    var count int
    err := adapter.DB().QueryRow("SELECT COUNT(*) FROM workspaces").Scan(&count)
    require.NoError(t, err)
    assert.Equal(t, 0, count)

    // Verify PathResolver is available
    assert.NotNil(t, adapter.PathResolver())
    assert.Contains(t, adapter.PathResolver().DatabasePath(), "ticketr.db")
}
```

---

**Document Version**: 1.0
**Status**: Integration plan documented, ready for implementation
**Next Action**: Begin Step 1 (Update SQLiteAdapter) when ready to proceed