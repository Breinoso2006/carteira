package migration

import (
	"fmt"
	"log"
	"time"

	"github.com/breinoso2006/carteira-api/internal/database"
	"github.com/breinoso2006/carteira-api/internal/models"
)

// MigrationTool handles one-time migration of in-memory portfolio data into the
// SQLite database.
type MigrationTool struct {
	db *database.Database
}

// NewMigrationTool creates a new MigrationTool backed by the given database.
func NewMigrationTool(db *database.Database) *MigrationTool {
	return &MigrationTool{db: db}
}

// MigratePortfolio inserts every entry from inMemoryPortfolio into the database,
// skipping entries that already exist and logging (but not failing on) invalid
// entries.
func (m *MigrationTool) MigratePortfolio(inMemoryPortfolio []*models.StockInPortfolio) error {
	if len(inMemoryPortfolio) == 0 {
		log.Println("migration: no portfolio data to migrate")
		return nil
	}

	log.Printf("migration: migrating %d portfolio entries to database", len(inMemoryPortfolio))

	now := time.Now().UTC()

	for _, sip := range inMemoryPortfolio {
		if sip == nil || sip.Stock == nil {
			log.Println("migration: skipping nil portfolio entry")
			continue
		}

		ticker := sip.Stock.Ticker
		grade := sip.Stock.FundamentalistGrade

		// Check whether the entry already exists.
		var count int
		err := m.db.GetDB().QueryRow(
			"SELECT COUNT(*) FROM portfolio_entries WHERE ticker = ?", ticker,
		).Scan(&count)
		if err != nil {
			return fmt.Errorf("migration: failed to check existing entry for %s: %w", ticker, err)
		}

		if count > 0 {
			log.Printf("migration: entry for %s already exists, skipping", ticker)
			continue
		}

		_, err = m.db.GetDB().Exec(
			"INSERT INTO portfolio_entries (ticker, fundamentalist_grade, created_at, updated_at) VALUES (?, ?, ?, ?)",
			ticker, grade, now, now,
		)
		if err != nil {
			return fmt.Errorf("migration: failed to insert entry for %s: %w", ticker, err)
		}

		log.Printf("migration: migrated %s (grade: %.2f)", ticker, grade)
	}

	log.Println("migration: portfolio migration completed successfully")
	return nil
}

// VerifyMigration checks that the number of rows in portfolio_entries matches
// the number of valid entries in inMemoryPortfolio.
func (m *MigrationTool) VerifyMigration(inMemoryPortfolio []*models.StockInPortfolio) error {
	var dbCount int
	if err := m.db.GetDB().QueryRow("SELECT COUNT(*) FROM portfolio_entries").Scan(&dbCount); err != nil {
		return fmt.Errorf("migration: failed to count database entries: %w", err)
	}

	expectedCount := 0
	for _, sip := range inMemoryPortfolio {
		if sip != nil && sip.Stock != nil {
			expectedCount++
		}
	}

	if dbCount < expectedCount {
		return fmt.Errorf("migration: count mismatch — expected at least %d entries, got %d", expectedCount, dbCount)
	}

	log.Printf("migration: verification passed — %d entries in database", dbCount)
	return nil
}
