package repository

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/breinoso2006/carteira-api/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a temporary SQLite database for testing.
func setupTestDB(t *testing.T) *database.Database {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.NewDatabase(dbPath)
	require.NoError(t, err, "failed to create test database")
	t.Cleanup(func() { db.Close() })
	return db
}

func TestPortfolioRepository_Add(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	err := repo.Add("PETR4", 85.5)
	assert.NoError(t, err)

	// Duplicate should fail.
	err = repo.Add("PETR4", 85.5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestPortfolioRepository_GetAll(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	require.NoError(t, repo.Add("PETR4", 85.5))
	require.NoError(t, repo.Add("VALE3", 75.0))
	require.NoError(t, repo.Add("ITSA4", 90.0))

	entries, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, entries, 3)

	tickers := make(map[string]bool)
	for _, e := range entries {
		tickers[e.Ticker] = true
	}
	assert.True(t, tickers["PETR4"])
	assert.True(t, tickers["VALE3"])
	assert.True(t, tickers["ITSA4"])
}

func TestPortfolioRepository_Update(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	require.NoError(t, repo.Add("PETR4", 85.5))
	require.NoError(t, repo.Update("PETR4", 90.0))

	entries, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, 90.0, entries[0].FundamentalistGrade)
}

func TestPortfolioRepository_Remove(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	require.NoError(t, repo.Add("PETR4", 85.5))
	require.NoError(t, repo.Remove("PETR4"))

	entries, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, entries, 0)
}

func TestPortfolioRepository_RemoveNonExistent(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	err := repo.Remove("NONEXISTENT")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestPortfolioRepository_CalculateWeights(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	require.NoError(t, repo.Add("PETR4", 80.0))
	require.NoError(t, repo.Add("VALE3", 20.0))

	entries, err := repo.GetAll()
	require.NoError(t, err)

	require.NoError(t, repo.CalculateWeights(entries))

	var total float64
	for _, e := range entries {
		total += e.Weight
	}
	assert.InDelta(t, 100.0, total, 0.01)

	// Weights are quadratic: 80²=6400, 20²=400, total=6800
	// PETR4 = 6400/6800*100 ≈ 94.12%, VALE3 = 400/6800*100 ≈ 5.88%
	byTicker := make(map[string]float64)
	for _, e := range entries {
		byTicker[e.Ticker] = e.Weight
	}
	assert.InDelta(t, 94.12, byTicker["PETR4"], 0.01)
	assert.InDelta(t, 5.88, byTicker["VALE3"], 0.01)
}

func TestPortfolioRepository_EmptyDatabase(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	entries, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, entries, 0)
}

func TestPortfolioRepository_ConcurrentAdd(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			err := repo.Add("TICKER", 85.5)
			done <- err == nil
		}()
	}

	successes := 0
	for i := 0; i < 10; i++ {
		if <-done {
			successes++
		}
	}

	// Only one insert should succeed due to the UNIQUE constraint.
	assert.Equal(t, 1, successes)

	entries, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestPortfolioRepository_UpdatedAtTimestamp(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	require.NoError(t, repo.Add("PETR4", 85.5))

	entries, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, entries, 1)

	createdAt := entries[0].CreatedAt
	assert.False(t, createdAt.IsZero())

	// Wait a moment so the updated_at timestamp will differ.
	time.Sleep(10 * time.Millisecond)
	require.NoError(t, repo.Update("PETR4", 90.0))

	entries, err = repo.GetAll()
	require.NoError(t, err)
	require.Len(t, entries, 1)

	assert.False(t, entries[0].UpdatedAt.IsZero())
	assert.True(t, entries[0].UpdatedAt.After(createdAt) || entries[0].UpdatedAt.Equal(createdAt),
		"updated_at should be >= created_at")
}
