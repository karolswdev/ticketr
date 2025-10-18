-- Migration 001: Initial schema for v3.0
-- This maintains backward compatibility while adding SQLite storage

-- Version tracking
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Workspaces (initially single default workspace for backward compatibility)
CREATE TABLE IF NOT EXISTS workspaces (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    jira_url TEXT,
    project_key TEXT,
    credential_ref TEXT,
    config JSON,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create default workspace for backward compatibility
INSERT OR IGNORE INTO workspaces (id, name, is_default)
VALUES ('default', 'default', TRUE);

-- Create unique index for default workspace
CREATE UNIQUE INDEX IF NOT EXISTS idx_default_workspace
ON workspaces(is_default) WHERE is_default = TRUE;

-- Tickets table with full content storage
CREATE TABLE IF NOT EXISTS tickets (
    id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL DEFAULT 'default' REFERENCES workspaces(id) ON DELETE CASCADE,
    jira_id TEXT,
    title TEXT NOT NULL,
    description TEXT,
    custom_fields JSON,
    acceptance_criteria JSON,
    tasks JSON,
    local_hash TEXT,
    remote_hash TEXT,
    sync_status TEXT DEFAULT 'new' CHECK(sync_status IN ('new', 'synced', 'modified', 'conflict')),
    source_line INTEGER,
    last_synced TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(workspace_id, jira_id)
);

-- Indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_ticket_workspace ON tickets(workspace_id);
CREATE INDEX IF NOT EXISTS idx_ticket_jira ON tickets(jira_id);
CREATE INDEX IF NOT EXISTS idx_ticket_status ON tickets(sync_status);
CREATE INDEX IF NOT EXISTS idx_ticket_modified ON tickets(updated_at);

-- State tracking (replaces .ticketr.state file)
CREATE TABLE IF NOT EXISTS ticket_state (
    ticket_id TEXT PRIMARY KEY REFERENCES tickets(id) ON DELETE CASCADE,
    local_hash TEXT NOT NULL,
    remote_hash TEXT NOT NULL,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sync operations log for audit trail
CREATE TABLE IF NOT EXISTS sync_operations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workspace_id TEXT REFERENCES workspaces(id),
    operation TEXT NOT NULL CHECK(operation IN ('push', 'pull', 'migrate', 'conflict_resolve')),
    file_path TEXT,
    ticket_count INTEGER,
    success_count INTEGER,
    failure_count INTEGER,
    conflict_count INTEGER,
    duration_ms INTEGER,
    error_details JSON,
    started_at TIMESTAMP,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sync_workspace ON sync_operations(workspace_id);
CREATE INDEX IF NOT EXISTS idx_sync_timestamp ON sync_operations(started_at);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_ticket_timestamp
AFTER UPDATE ON tickets
BEGIN
    UPDATE tickets SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_workspace_timestamp
AFTER UPDATE ON workspaces
BEGIN
    UPDATE workspaces SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Mark this migration as applied
INSERT INTO schema_migrations (version, name)
VALUES (1, '001_initial_schema');