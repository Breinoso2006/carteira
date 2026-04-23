package migration

import (
	"path/filepath"
	"testing"

	"github.com/breinoso2006/carteira-api/internal/database"
	"github.com/breinoso2006/carteira-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a temporary SQLite database for migration tests.
func setupTestDB(t *testing.T) *database.Database {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "migration_test.db")
	db, err := database.NewDatabase(dbPath)
	require.NoError(t, err, "failed to create test database")
	t.Cleanup(func() { db.Close() })
	return db
}

// makeStock is a helper that builds a *models.StockInPortfolio.
func makeStock(ticker string, grade float64) *models.StockInPortfolio {
	return &models.StockInPortfolio{
		Stock: models.NewStock(ticker, grade),
	}
}

// countRows returns the number of rows in portfolio_entries.
func countRows(t *testing.T, db *database.Database) int {
	t.Helper()
	var n int
	err := db.GetDB().QueryRow("SELECT COUNT(*) FROM portfolio_entries").Scan(&n)
	require.NoError(t, err)
	return n
}

// TestMigratePortfolio_Success verifies that valid portfolio entries are
// inserted into the database.
func TestMigratePortfolio_Success(t *testing.T) {
	db := setupTestDB(t)
	tool := NewMigrationTool(db)

	portfolio := []*models.StockInPortfolio{
		makeStock("PETR4", 85.5),
		makeStock("VALE3", 75.0),
	}

	err := tool.MigratePortfolio(portfolio)
	require.NoError(t, err)

	assert.Equal(t, 2, countRows(t, db))
}

// TestMigratePortfolio_Empty verifies that migrating an empty slice returns no
// error and leaves the database unchanged.
func TestMigratePortfolio_Empty(t *testing.T) {
	db := setupTestDB(t)
	tool := NewMigrationTool(db)

	err := tool.MigratePortfolio([]*models.StockInPortfolio{})
	require.NoError(t, err)

	assert.Equal(t, 0, countRows(t, db))
}

// TestMigratePortfolio_SkipExisting verifies that migrating the same entries
// twice does not create duplicates.
func TestMigratePortfolio_SkipExisting(t *testing.T) {
	db := setupTestDB(t)
	tool := NewMigrationTool(db)

	portfolio := []*models.StockInPortfolio{
		makeStock("PETR4", 85.5),
		makeStock("VALE3", 75.0),
	}

	require.NoError(t, tool.MigratePortfolio(portfolio))
	// Second migration of the same data should be a no-op.
	require.NoError(t, tool.MigratePortfolio(portfolio))

	assert.Equal(t, 2, countRows(t, db), "no duplicate rows should be created")
}

// TestMigratePortfolio_NilEntry verifies that nil entries in the portfolio
// slice are silently skipped.
func TestMigratePortfolio_NilEntry(t *testing.T) {
	db := setupTestDB(t)
	tool := NewMigrationTool(db)

	portfolio := []*models.StockInPortfolio{
		makeStock("PETR4", 85.5),
		nil,
		makeStock("VALE3", 75.0),
	}

	err := tool.MigratePortfolio(portfolio)
	require.NoError(t, err)

	// Only the two valid entries should be inserted.
	assert.Equal(t, 2, countRows(t, db))
}

// TestVerifyMigration_Success verifies that VerifyMigration passes when the
// database contains at least as many entries as the in-memory portfolio.
func TestVerifyMigration_Success(t *testing.T) {
	db := setupTestDB(t)
	tool := NewMigrationTool(db)

	portfolio := []*models.StockInPortfolio{
		makeStock("PETR4", 85.5),
		makeStock("VALE3", 75.0),
	}

	require.NoError(t, tool.MigratePortfolio(portfolio))

	err := tool.VerifyMigration(portfolio)
	assert.NoError(t, err)
}

// TestVerifyMigration_CountMismatch verifies that VerifyMigration returns an
// error when the database has fewer entries than the in-memory portfolio.
func TestVerifyMigration_CountMismatch(t *testing.T) {
	db := setupTestDB(t)
	tool := NewMigrationTool(db)

	// Insert only one entry into the database.
	require.NoError(t, tool.MigratePortfolio([]*models.StockInPortfolio{
		makeStock("PETR4", 85.5),
	}))

	// But verify against a portfolio with two entries.
	portfolio := []*models.StockInPortfolio{
		makeStock("PETR4", 85.5),
		makeStock("VALE3", 75.0),
	}

	err := tool.VerifyMigration(portfolio)
	assert.Error(t, err, "VerifyMigration should fail when DB has fewer entries")
	assert.Contains(t, err.Error(), "count mismatch")
}
