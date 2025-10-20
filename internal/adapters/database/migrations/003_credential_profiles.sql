-- Migration 003: Add credential profiles for reusable Jira credentials
-- This allows multiple workspaces to share the same credential configuration

-- Create credential_profiles table
CREATE TABLE IF NOT EXISTS credential_profiles (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    jira_url TEXT NOT NULL,
    username TEXT NOT NULL,
    keychain_ref TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add credential_profile_id foreign key to workspaces table
-- This allows workspaces to optionally reference a credential profile
ALTER TABLE workspaces
    ADD COLUMN credential_profile_id TEXT REFERENCES credential_profiles(id);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_credential_profile_name ON credential_profiles(name);
CREATE INDEX IF NOT EXISTS idx_workspace_profile ON workspaces(credential_profile_id);

-- Add trigger to update updated_at timestamp for credential_profiles
CREATE TRIGGER IF NOT EXISTS update_credential_profile_timestamp
AFTER UPDATE ON credential_profiles
BEGIN
    UPDATE credential_profiles SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Mark this migration as applied
INSERT INTO schema_migrations (version, name)
VALUES (3, '003_credential_profiles');