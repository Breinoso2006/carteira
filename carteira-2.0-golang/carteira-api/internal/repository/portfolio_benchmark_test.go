package repository

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/breinoso2006/carteira-api/internal/database"
)

// setupBenchmarkDB creates a temporary SQLite database pre-populated with n
// portfolio entries for use in benchmarks.
func setupBenchmarkDB(b *testing.B, n int) *PortfolioRepository {
	b.Helper()
	dbPath := filepath.Join(b.TempDir(), "bench.db")
	db, err := database.NewDatabase(dbPath)
	if err != nil {
		b.Fatalf("failed to create benchmark database: %v", err)
	}
	b.Cleanup(func() { db.Close() })

	repo := NewPortfolioRepository(db)
	for i := 0; i < n; i++ {
		ticker := fmt.Sprintf("TK%04d", i)
		grade := float64(50 + (i % 50))
		if err := repo.Add(ticker, grade); err != nil {
			b.Fatalf("setup: Add(%s) failed: %v", ticker, err)
		}
	}
	return repo
}

// BenchmarkGetAll_100Stocks measures the latency of loading all portfolio
// entries when the database contains 100 rows.
//
// Requirement 7.3: WHEN portfolio data is loaded, THE DatabaseRepository SHALL
// load all entries within 1 second for up to 100 stocks.
//
// SQLite performs a full table scan on a 100-row table in microseconds, so
// each iteration should complete well under the 1-second budget.
func BenchmarkGetAll_100Stocks(b *testing.B) {
	repo := setupBenchmarkDB(b, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entries, err := repo.GetAll()
		if err != nil {
			b.Fatalf("GetAll error: %v", err)
		}
		if len(entries) != 100 {
			b.Fatalf("expected 100 entries, got %d", len(entries))
		}
	}
}

// BenchmarkAdd measures the latency of inserting a single portfolio entry.
//
// Each iteration inserts a unique ticker to avoid UNIQUE constraint violations.
func BenchmarkAdd(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "bench_add.db")
	db, err := database.NewDatabase(dbPath)
	if err != nil {
		b.Fatalf("failed to create benchmark database: %v", err)
	}
	b.Cleanup(func() { db.Close() })

	repo := NewPortfolioRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ticker := fmt.Sprintf("ADD%07d", i)
		if err := repo.Add(ticker, 75.0); err != nil {
			b.Fatalf("Add error: %v", err)
		}
	}
}

// BenchmarkGetAll_Empty measures the baseline latency of GetAll on an empty
// database (no rows to scan).
func BenchmarkGetAll_Empty(b *testing.B) {
	repo := setupBenchmarkDB(b, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entries, err := repo.GetAll()
		if err != nil {
			b.Fatalf("GetAll error: %v", err)
		}
		if len(entries) != 0 {
			b.Fatalf("expected 0 entries, got %d", len(entries))
		}
	}
}
