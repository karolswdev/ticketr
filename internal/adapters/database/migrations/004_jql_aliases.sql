-- Migration 004: JQL Aliases
-- Adds support for named, reusable JQL queries

-- JQL aliases table
CREATE TABLE IF NOT EXISTS jql_aliases (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    jql TEXT NOT NULL,
    description TEXT,
    is_predefined BOOLEAN DEFAULT FALSE,
    workspace_id TEXT REFERENCES workspaces(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, workspace_id)
);

-- Indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_alias_workspace ON jql_aliases(workspace_id);
CREATE INDEX IF NOT EXISTS idx_alias_name ON jql_aliases(name);
CREATE INDEX IF NOT EXISTS idx_alias_predefined ON jql_aliases(is_predefined);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_alias_timestamp
AFTER UPDATE ON jql_aliases
BEGIN
    UPDATE jql_aliases SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Mark this migration as applied
INSERT INTO schema_migrations (version, name)
VALUES (4, '004_jql_aliases');
