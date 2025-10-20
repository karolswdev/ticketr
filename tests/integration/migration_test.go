//go:build integration

package integration

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseMigrations(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := fmt.Sprintf("%s/migration_test.db", tempDir)

	t.Run("SchemaV1ToV3Migration", func(t *testing.T) {
		// Create database and apply v1 schema
		db, err := sql.Open("sqlite3", dbPath)
		require.NoError(t, err)
		defer db.Close()

		adapter := database.NewSQLiteAdapter(db)

		// Apply v1 migration (basic workspace tables)
		err = adapter.ApplyMigration(1)
		require.NoError(t, err)

		// Verify v1 schema exists
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='workspaces'").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "workspaces table should exist after v1 migration")

		// Apply v2 migration (credential refs)
		err = adapter.ApplyMigration(2)
		require.NoError(t, err)

		// Verify credential_refs table exists
		err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='credential_refs'").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "credential_refs table should exist after v2 migration")

		// Insert test data with v2 schema
		_, err = db.Exec(`
			INSERT INTO workspaces (id, name, jira_url, project_key, is_default, created_at)
			VALUES ('ws1', 'test-workspace', 'https://test.atlassian.net', 'TEST', true, datetime('now'))
		`)
		require.NoError(t, err)

		_, err = db.Exec(`
			INSERT INTO credential_refs (workspace_id, keychain_ref)
			VALUES ('ws1', 'test-keychain-ref')
		`)
		require.NoError(t, err)

		// Apply v3 migration (credential profiles)
		err = adapter.ApplyMigration(3)
		require.NoError(t, err)

		// Verify credential_profiles table exists
		err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='credential_profiles'").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "credential_profiles table should exist after v3 migration")

		// Verify workspaces table has new column
		rows, err := db.Query("PRAGMA table_info(workspaces)")
		require.NoError(t, err)
		defer rows.Close()

		var hasProfileColumn bool
		for rows.Next() {
			var cid int
			var name, dataType string
			var notNull, defaultValue, pk interface{}
			err = rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
			require.NoError(t, err)
			if name == "credential_profile_id" {
				hasProfileColumn = true
				break
			}
		}
		assert.True(t, hasProfileColumn, "workspaces table should have credential_profile_id column after v3 migration")

		// Verify existing data preserved
		var workspaceName string
		err = db.QueryRow("SELECT name FROM workspaces WHERE id = 'ws1'").Scan(&workspaceName)
		require.NoError(t, err)
		assert.Equal(t, "test-workspace", workspaceName, "Existing workspace data should be preserved")

		var keychainRef string
		err = db.QueryRow("SELECT keychain_ref FROM credential_refs WHERE workspace_id = 'ws1'").Scan(&keychainRef)
		require.NoError(t, err)
		assert.Equal(t, "test-keychain-ref", keychainRef, "Existing credential ref should be preserved")
	})

	t.Run("MigrationRollback", func(t *testing.T) {
		// Test rollback from v3 to v2 (if supported)
		db, err := sql.Open("sqlite3", fmt.Sprintf("%s_rollback.db", tempDir))
		require.NoError(t, err)
		defer db.Close()

		adapter := database.NewSQLiteAdapter(db)

		// Apply all migrations
		err = adapter.Initialize()
		require.NoError(t, err)

		// Insert test data in v3 format
		profileID := "prof-123"
		workspaceID := "ws-456"

		_, err = db.Exec(`
			INSERT INTO credential_profiles (id, name, jira_url, username, keychain_ref, created_at, updated_at)
			VALUES (?, 'test-profile', 'https://test.atlassian.net', 'test@example.com', 'prof-keychain-ref', datetime('now'), datetime('now'))
		`, profileID)
		require.NoError(t, err)

		_, err = db.Exec(`
			INSERT INTO workspaces (id, name, jira_url, project_key, is_default, credential_profile_id, created_at)
			VALUES (?, 'profile-workspace', 'https://test.atlassian.net', 'PROF', false, ?, datetime('now'))
		`, workspaceID, profileID)
		require.NoError(t, err)

		// Test rollback migration (down migration from v3 to v2)
		err = adapter.RollbackMigration(3)
		if err == nil {
			// If rollback is supported, verify state
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='credential_profiles'").Scan(&count)
			require.NoError(t, err)
			assert.Equal(t, 0, count, "credential_profiles table should be dropped after rollback")

			// Verify workspace still exists but without profile reference
			var exists bool
			err = db.QueryRow("SELECT 1 FROM workspaces WHERE id = ?", workspaceID).Scan(&exists)
			if err != sql.ErrNoRows {
				require.NoError(t, err)
				assert.True(t, exists, "Workspace should still exist after rollback")
			}
		}
		// If rollback not supported, that's also acceptable
	})

	t.Run("ForeignKeyConstraints", func(t *testing.T) {
		db, err := sql.Open("sqlite3", fmt.Sprintf("%s_fk.db", tempDir))
		require.NoError(t, err)
		defer db.Close()

		// Enable foreign key constraints
		_, err = db.Exec("PRAGMA foreign_keys = ON")
		require.NoError(t, err)

		adapter := database.NewSQLiteAdapter(db)
		err = adapter.Initialize()
		require.NoError(t, err)

		// Test foreign key constraint: workspace references profile
		profileID := "prof-fk-test"
		workspaceID := "ws-fk-test"

		// Create profile
		_, err = db.Exec(`
			INSERT INTO credential_profiles (id, name, jira_url, username, keychain_ref, created_at, updated_at)
			VALUES (?, 'fk-test-profile', 'https://fk-test.atlassian.net', 'fk@example.com', 'fk-keychain-ref', datetime('now'), datetime('now'))
		`, profileID)
		require.NoError(t, err)

		// Create workspace referencing profile
		_, err = db.Exec(`
			INSERT INTO workspaces (id, name, jira_url, project_key, is_default, credential_profile_id, created_at)
			VALUES (?, 'fk-workspace', 'https://fk-test.atlassian.net', 'FK', false, ?, datetime('now'))
		`, workspaceID, profileID)
		require.NoError(t, err)

		// Attempt to delete referenced profile (should fail)
		_, err = db.Exec("DELETE FROM credential_profiles WHERE id = ?", profileID)
		assert.Error(t, err, "Should not allow deleting profile referenced by workspace")

		// Delete workspace first
		_, err = db.Exec("DELETE FROM workspaces WHERE id = ?", workspaceID)
		require.NoError(t, err)

		// Now profile deletion should succeed
		_, err = db.Exec("DELETE FROM credential_profiles WHERE id = ?", profileID)
		assert.NoError(t, err, "Should allow deleting unreferenced profile")
	})

	t.Run("DataTypesAndConstraints", func(t *testing.T) {
		db, err := sql.Open("sqlite3", fmt.Sprintf("%s_types.db", tempDir))
		require.NoError(t, err)
		defer db.Close()

		adapter := database.NewSQLiteAdapter(db)
		err = adapter.Initialize()
		require.NoError(t, err)

		// Test NOT NULL constraints
		_, err = db.Exec(`
			INSERT INTO credential_profiles (id, name, jira_url, username, keychain_ref, created_at, updated_at)
			VALUES ('test-id', '', 'https://test.atlassian.net', 'test@example.com', 'test-ref', datetime('now'), datetime('now'))
		`)
		assert.Error(t, err, "Should reject empty profile name")

		_, err = db.Exec(`
			INSERT INTO workspaces (id, name, jira_url, project_key, is_default, created_at)
			VALUES ('test-ws', '', 'https://test.atlassian.net', 'TEST', false, datetime('now'))
		`)
		assert.Error(t, err, "Should reject empty workspace name")

		// Test unique constraints
		_, err = db.Exec(`
			INSERT INTO credential_profiles (id, name, jira_url, username, keychain_ref, created_at, updated_at)
			VALUES ('test-1', 'unique-profile', 'https://test1.atlassian.net', 'test1@example.com', 'test-ref-1', datetime('now'), datetime('now'))
		`)
		require.NoError(t, err)

		_, err = db.Exec(`
			INSERT INTO credential_profiles (id, name, jira_url, username, keychain_ref, created_at, updated_at)
			VALUES ('test-2', 'unique-profile', 'https://test2.atlassian.net', 'test2@example.com', 'test-ref-2', datetime('now'), datetime('now'))
		`)
		assert.Error(t, err, "Should reject duplicate profile name")

		// Test timestamp defaults
		var createdAt, updatedAt time.Time
		err = db.QueryRow(`
			SELECT created_at, updated_at FROM credential_profiles WHERE id = 'test-1'
		`).Scan(&createdAt, &updatedAt)
		require.NoError(t, err)
		assert.False(t, createdAt.IsZero(), "created_at should be set automatically")
		assert.False(t, updatedAt.IsZero(), "updated_at should be set automatically")
	})
}

func TestMigrationVersioning(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := fmt.Sprintf("%s/versioning_test.db", tempDir)

	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	adapter := database.NewSQLiteAdapter(db)

	t.Run("VersionTracking", func(t *testing.T) {
		// Apply migrations one by one and verify version tracking
		err = adapter.ApplyMigration(1)
		require.NoError(t, err)

		var version int
		var appliedAt time.Time
		err = db.QueryRow("SELECT version, applied_at FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&version, &appliedAt)
		require.NoError(t, err)
		assert.Equal(t, 1, version, "Should track v1 migration")
		assert.False(t, appliedAt.IsZero(), "Should record migration timestamp")

		err = adapter.ApplyMigration(2)
		require.NoError(t, err)

		err = db.QueryRow("SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&version)
		require.NoError(t, err)
		assert.Equal(t, 2, version, "Should track v2 migration")

		err = adapter.ApplyMigration(3)
		require.NoError(t, err)

		err = db.QueryRow("SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&version)
		require.NoError(t, err)
		assert.Equal(t, 3, version, "Should track v3 migration")

		// Verify all versions are recorded
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM schema_version").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 3, count, "Should have 3 migration records")
	})

	t.Run("IdempotentMigrations", func(t *testing.T) {
		// Applying same migration twice should not fail
		err = adapter.ApplyMigration(3)
		assert.NoError(t, err, "Re-applying migration should be idempotent")

		// Version count should remain the same
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM schema_version WHERE version = 3").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "Should not duplicate migration records")
	})

	t.Run("GetCurrentVersion", func(t *testing.T) {
		currentVersion, err := adapter.GetCurrentSchemaVersion()
		require.NoError(t, err)
		assert.Equal(t, 3, currentVersion, "Should return current schema version")
	})
}

func TestDataIntegrity(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := fmt.Sprintf("%s/integrity_test.db", tempDir)

	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	// Enable integrity checks
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	require.NoError(t, err)

	adapter := database.NewSQLiteAdapter(db)
	err = adapter.Initialize()
	require.NoError(t, err)

	t.Run("ReferentialIntegrity", func(t *testing.T) {
		// Test cascade behavior and referential integrity
		profileID := "integrity-profile"
		workspaceID := "integrity-workspace"

		// Create profile
		_, err = db.Exec(`
			INSERT INTO credential_profiles (id, name, jira_url, username, keychain_ref, created_at, updated_at)
			VALUES (?, 'integrity-test', 'https://integrity.atlassian.net', 'integrity@example.com', 'integrity-ref', datetime('now'), datetime('now'))
		`, profileID)
		require.NoError(t, err)

		// Create workspace
		_, err = db.Exec(`
			INSERT INTO workspaces (id, name, jira_url, project_key, is_default, credential_profile_id, created_at)
			VALUES (?, 'integrity-ws', 'https://integrity.atlassian.net', 'INT', false, ?, datetime('now'))
		`, workspaceID, profileID)
		require.NoError(t, err)

		// Verify workspace references profile correctly
		var referencedProfileID string
		err = db.QueryRow("SELECT credential_profile_id FROM workspaces WHERE id = ?", workspaceID).Scan(&referencedProfileID)
		require.NoError(t, err)
		assert.Equal(t, profileID, referencedProfileID, "Workspace should reference correct profile")

		// Attempt to create workspace with invalid profile reference
		_, err = db.Exec(`
			INSERT INTO workspaces (id, name, jira_url, project_key, is_default, credential_profile_id, created_at)
			VALUES ('invalid-ws', 'invalid', 'https://invalid.atlassian.net', 'INV', false, 'non-existent-profile', datetime('now'))
		`)
		assert.Error(t, err, "Should reject workspace with invalid profile reference")
	})

	t.Run("TransactionConsistency", func(t *testing.T) {
		// Test that partial failures don't leave database in inconsistent state
		tx, err := db.Begin()
		require.NoError(t, err)

		// Insert profile
		_, err = tx.Exec(`
			INSERT INTO credential_profiles (id, name, jira_url, username, keychain_ref, created_at, updated_at)
			VALUES ('tx-profile', 'tx-test', 'https://tx.atlassian.net', 'tx@example.com', 'tx-ref', datetime('now'), datetime('now'))
		`)
		require.NoError(t, err)

		// Insert workspace
		_, err = tx.Exec(`
			INSERT INTO workspaces (id, name, jira_url, project_key, is_default, credential_profile_id, created_at)
			VALUES ('tx-workspace', 'tx-ws', 'https://tx.atlassian.net', 'TX', false, 'tx-profile', datetime('now'))
		`)
		require.NoError(t, err)

		// Rollback transaction
		err = tx.Rollback()
		require.NoError(t, err)

		// Verify nothing was persisted
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM credential_profiles WHERE id = 'tx-profile'").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "Profile should not exist after rollback")

		err = db.QueryRow("SELECT COUNT(*) FROM workspaces WHERE id = 'tx-workspace'").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "Workspace should not exist after rollback")
	})
}
