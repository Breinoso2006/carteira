package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPortfolioIntegration_FullCRUDCycle tests the complete lifecycle of portfolio
// entries: Add → GetAll → Update → GetAll → Remove → GetAll.
func TestPortfolioIntegration_FullCRUDCycle(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	// Step 1: Add entries.
	require.NoError(t, repo.Add("PETR4", 85.5))
	require.NoError(t, repo.Add("VALE3", 70.0))
	require.NoError(t, repo.Add("ITSA4", 90.0))

	// Step 2: GetAll – verify all three are present.
	entries, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, entries, 3)

	tickerGrades := make(map[string]float64)
	for _, e := range entries {
		tickerGrades[e.Ticker] = e.FundamentalistGrade
	}
	assert.Equal(t, 85.5, tickerGrades["PETR4"])
	assert.Equal(t, 70.0, tickerGrades["VALE3"])
	assert.Equal(t, 90.0, tickerGrades["ITSA4"])

	// Step 3: Update one entry.
	require.NoError(t, repo.Update("VALE3", 95.0))

	// Step 4: GetAll – verify the update is reflected.
	entries, err = repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, entries, 3)

	tickerGrades = make(map[string]float64)
	for _, e := range entries {
		tickerGrades[e.Ticker] = e.FundamentalistGrade
	}
	assert.Equal(t, 95.0, tickerGrades["VALE3"], "VALE3 grade should be updated to 95.0")
	assert.Equal(t, 85.5, tickerGrades["PETR4"], "PETR4 grade should remain unchanged")

	// Step 5: Remove one entry.
	require.NoError(t, repo.Remove("ITSA4"))

	// Step 6: GetAll – verify removal.
	entries, err = repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, entries, 2)

	remaining := make(map[string]bool)
	for _, e := range entries {
		remaining[e.Ticker] = true
	}
	assert.True(t, remaining["PETR4"])
	assert.True(t, remaining["VALE3"])
	assert.False(t, remaining["ITSA4"], "ITSA4 should have been removed")
}

// TestPortfolioIntegration_WeightCalculation adds multiple stocks, then verifies
// that weights sum to 100% and are proportional to the fundamentalist grades.
func TestPortfolioIntegration_WeightCalculation(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	// Add stocks with known grades so we can verify exact proportions.
	require.NoError(t, repo.Add("AAPL", 40.0))
	require.NoError(t, repo.Add("GOOG", 30.0))
	require.NoError(t, repo.Add("MSFT", 20.0))
	require.NoError(t, repo.Add("AMZN", 10.0))

	entries, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, entries, 4)

	require.NoError(t, repo.CalculateWeights(entries))

	// Weights must sum to 100%.
	var totalWeight float64
	for _, e := range entries {
		totalWeight += e.Weight
	}
	assert.InDelta(t, 100.0, totalWeight, 0.01, "weights must sum to 100%%")

	// Verify each weight is proportional to its grade (total grade = 100).
	byTicker := make(map[string]float64)
	for _, e := range entries {
		byTicker[e.Ticker] = e.Weight
	}
	assert.InDelta(t, 40.0, byTicker["AAPL"], 0.01)
	assert.InDelta(t, 30.0, byTicker["GOOG"], 0.01)
	assert.InDelta(t, 20.0, byTicker["MSFT"], 0.01)
	assert.InDelta(t, 10.0, byTicker["AMZN"], 0.01)
}

// TestPortfolioIntegration_UpdateNonExistent verifies that updating a ticker that
// does not exist in the portfolio returns an appropriate error.
func TestPortfolioIntegration_UpdateNonExistent(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	err := repo.Update("NONEXISTENT", 75.0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found", "error should indicate the ticker was not found")
}

// TestPortfolioIntegration_LargePortfolio adds 20 stocks and verifies all are
// retrieved correctly with the right grades.
func TestPortfolioIntegration_LargePortfolio(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	const count = 20
	expectedGrades := make(map[string]float64, count)

	for i := range count {
		ticker := fmt.Sprintf("TICK%02d", i)
		grade := float64(i+1) * 5.0 // grades: 5, 10, 15, ..., 100
		expectedGrades[ticker] = grade
		require.NoError(t, repo.Add(ticker, grade), "failed to add ticker %s", ticker)
	}

	entries, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, entries, count, "all %d stocks should be retrieved", count)

	// Verify every entry has the correct grade.
	for _, e := range entries {
		expected, ok := expectedGrades[e.Ticker]
		require.True(t, ok, "unexpected ticker %s in results", e.Ticker)
		assert.Equal(t, expected, e.FundamentalistGrade,
			"grade mismatch for ticker %s", e.Ticker)
	}
}

// TestPortfolioIntegration_WeightCalculationEmpty verifies that calling
// CalculateWeights on an empty portfolio does not return an error and leaves
// the (empty) slice unchanged.
func TestPortfolioIntegration_WeightCalculationEmpty(t *testing.T) {
	repo := NewPortfolioRepository(setupTestDB(t))

	entries, err := repo.GetAll()
	require.NoError(t, err)
	assert.Empty(t, entries)

	// CalculateWeights on an empty slice must not error.
	err = repo.CalculateWeights(entries)
	assert.NoError(t, err)
	assert.Empty(t, entries)
}
