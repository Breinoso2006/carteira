package database

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewDatabase_Success verifies that a database can be created at a valid
// path and that the connection is usable.
func TestNewDatabase_Success(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := NewDatabase(dbPath)
	require.NoError(t, err, "NewDatabase should succeed for a valid path")
	require.NotNil(t, db)
	t.Cleanup(func() { db.Close() })

	// Verify the underlying connection is alive.
	require.NotNil(t, db.GetDB())
	assert.NoError(t, db.GetDB().Ping())
}

// TestNewDatabase_InvalidPath verifies that NewDatabase returns an error when
// the path is not writable (e.g. a non-existent root directory on Linux).
func TestNewDatabase_InvalidPath(t *testing.T) {
	// /nonexistent/path does not exist and cannot be created by a normal user.
	db, err := NewDatabase("/nonexistent/path/db.db")
	assert.Error(t, err, "NewDatabase should fail for an invalid path")
	assert.Nil(t, db)
}

// TestNewDatabase_SchemaCreation verifies that all expected tables are present
// after initialisation.
func TestNewDatabase_SchemaCreation(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "schema_test.db")

	db, err := NewDatabase(dbPath)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	expectedTables := []string{"portfolio_entries", "stock_cache", "schema_version"}
	for _, table := range expectedTables {
		var name string
		err := db.GetDB().QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?", table,
		).Scan(&name)
		require.NoError(t, err, "table %q should exist", table)
		assert.Equal(t, table, name)
	}
}

// TestNewDatabase_MigrationRecorded verifies that schema_version contains
// version 1 after a fresh initialisation.
func TestNewDatabase_MigrationRecorded(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "migration_test.db")

	db, err := NewDatabase(dbPath)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	var version int
	err = db.GetDB().QueryRow("SELECT MAX(version) FROM schema_version").Scan(&version)
	require.NoError(t, err)
	assert.Equal(t, 1, version, "schema_version should record version 1 after init")
}

// TestNewDatabase_Idempotent verifies that opening the same database file twice
// does not produce errors (the schema DDL uses CREATE TABLE IF NOT EXISTS).
func TestNewDatabase_Idempotent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "idempotent_test.db")

	db1, err := NewDatabase(dbPath)
	require.NoError(t, err, "first open should succeed")
	require.NoError(t, db1.Close())

	db2, err := NewDatabase(dbPath)
	require.NoError(t, err, "second open of the same file should succeed")
	require.NotNil(t, db2)
	t.Cleanup(func() { db2.Close() })

	// Schema version should still be 1.
	var version int
	err = db2.GetDB().QueryRow("SELECT MAX(version) FROM schema_version").Scan(&version)
	require.NoError(t, err)
	assert.Equal(t, 1, version)
}

// TestDatabase_Close verifies that Close() returns no error on a live
// connection and that a second call is also safe.
func TestDatabase_Close(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "close_test.db")

	db, err := NewDatabase(dbPath)
	require.NoError(t, err)

	assert.NoError(t, db.Close(), "first Close() should succeed")
	// A second Close() on a nil/closed db should not panic.
	assert.NoError(t, db.Close(), "second Close() should not error")
}
