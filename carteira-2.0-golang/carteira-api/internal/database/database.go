package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const currentSchemaVersion = 1

// schema is the initial DDL executed on every startup (all statements are idempotent).
const schema = `
-- Table: portfolio_entries
CREATE TABLE IF NOT EXISTS portfolio_entries (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    ticker               TEXT    NOT NULL UNIQUE,
    fundamentalist_grade REAL    NOT NULL,
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: stock_cache
CREATE TABLE IF NOT EXISTS stock_cache (
    symbol        TEXT PRIMARY KEY,
    price         REAL,
    pe            REAL,
    pbv           REAL,
    psr           REAL,
    bvps          REAL,
    eps           REAL,
    dy            REAL,
    source        TEXT,
    created_at    TIMESTAMP,
    expires_at    TIMESTAMP,
    invalid_fields TEXT
);

-- Table: schema_version
CREATE TABLE IF NOT EXISTS schema_version (
    version     INTEGER PRIMARY KEY,
    migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

// Database wraps a sql.DB connection and exposes schema management helpers.
type Database struct {
	db *sql.DB
}

// NewDatabase opens (or creates) the SQLite database at dbPath, applies the
// base schema, and runs any pending migrations.
func NewDatabase(dbPath string) (*Database, error) {
	// Ensure the parent directory exists.
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory %q: %w", dir, err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %q: %w", dbPath, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database at %q: %w", dbPath, err)
	}

	d := &Database{db: db}

	if err := d.initializeSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return d, nil
}

// initializeSchema creates all tables (idempotent) and runs pending migrations.
func (d *Database) initializeSchema() error {
	if _, err := d.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to execute base schema: %w", err)
	}

	currentVersion, err := d.getSchemaVersion()
	if err != nil {
		return fmt.Errorf("failed to read schema version: %w", err)
	}

	if currentVersion < currentSchemaVersion {
		if err := d.runMigrations(currentVersion); err != nil {
			return fmt.Errorf("migration from version %d failed: %w", currentVersion, err)
		}
	}

	return nil
}

// getSchemaVersion returns the highest version recorded in schema_version, or 0
// if the table is empty (fresh database).
func (d *Database) getSchemaVersion() (int, error) {
	var version int
	err := d.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("query failed: %w", err)
	}
	return version, nil
}

// runMigrations applies all migrations from currentVersion up to
// currentSchemaVersion and records each step in schema_version.
func (d *Database) runMigrations(currentVersion int) error {
	// migrations is an ordered list of (targetVersion, SQL) pairs.
	// Add new entries here as the schema evolves.
	migrations := []struct {
		version int
		sql     string
	}{
		{
			version: 1,
			// Version 1 is the initial schema; tables already exist from the
			// base schema above, so we only need to record the version.
			sql: "",
		},
	}

	for _, m := range migrations {
		if m.version <= currentVersion {
			continue
		}

		if m.sql != "" {
			if _, err := d.db.Exec(m.sql); err != nil {
				return fmt.Errorf("migration to version %d failed: %w", m.version, err)
			}
		}

		_, err := d.db.Exec(
			"INSERT OR REPLACE INTO schema_version (version, migrated_at) VALUES (?, ?)",
			m.version, time.Now().UTC(),
		)
		if err != nil {
			return fmt.Errorf("failed to record schema version %d: %w", m.version, err)
		}
	}

	return nil
}

// Close closes the underlying database connection.
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// GetDB returns the raw *sql.DB for use by repositories.
func (d *Database) GetDB() *sql.DB {
	return d.db
}
