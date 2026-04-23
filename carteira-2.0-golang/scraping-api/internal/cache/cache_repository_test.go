package cache

import (
	"errors"
	"testing"
	"time"

	gocache "github.com/patrickmn/go-cache"

	"github.com/breinoso2006/scraping-api/internal/models"
)

// mockScraper is a test double for StockScraper.
type mockScraper struct {
	data *models.StockData
	err  error
	// callCount tracks how many times SearchStockInformation was called.
	callCount int
}

func (m *mockScraper) SearchStockInformation(symbol string) (*models.StockData, error) {
	m.callCount++
	return m.data, m.err
}

// helpers

func validStockData(symbol string) *models.StockData {
	price := 10.0
	pe := 15.0
	pbv := 1.5
	psr := 2.0
	bvps := 5.0
	eps := 1.0
	dy := 0.05
	return &models.StockData{
		Symbol: symbol,
		Price:  &price,
		PE:     &pe,
		PBV:    &pbv,
		PSR:    &psr,
		BVps:   &bvps,
		EPS:    &eps,
		DY:     &dy,
		Source: "test",
	}
}

// newTestRepo creates a CacheRepository with a 1-hour TTL and cache enabled.
func newTestRepo() *CacheRepository {
	return NewCacheRepository(1)
}

// ─── StoreStockData ───────────────────────────────────────────────────────────

// TestStoreStockData_ValidData verifies that valid data (no invalid fields) is
// stored in cache and can be retrieved afterwards.
func TestStoreStockData_ValidData(t *testing.T) {
	repo := newTestRepo()
	data := validStockData("PETR4")

	if err := repo.StoreStockData(data); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	hit, err := repo.HasValidCache("PETR4")
	if err != nil {
		t.Fatalf("HasValidCache returned error: %v", err)
	}
	if !hit {
		t.Error("expected cache hit after storing valid data")
	}
}

// TestStoreStockData_InvalidFields verifies that data with at least one invalid
// field is NOT stored in cache (Req 2.4).
func TestStoreStockData_InvalidFields(t *testing.T) {
	repo := newTestRepo()
	data := validStockData("VALE3")
	data.MarkFieldInvalid("Price")

	if err := repo.StoreStockData(data); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	hit, err := repo.HasValidCache("VALE3")
	if err != nil {
		t.Fatalf("HasValidCache returned error: %v", err)
	}
	if hit {
		t.Error("expected cache miss for data with invalid fields")
	}
}

// TestStoreStockData_NilData verifies that storing nil data returns an error.
func TestStoreStockData_NilData(t *testing.T) {
	repo := newTestRepo()

	err := repo.StoreStockData(nil)
	if err == nil {
		t.Error("expected error when storing nil data, got nil")
	}
}

// TestStoreStockData_CacheDisabled verifies that StoreStockData is a no-op when
// cache is disabled (Req 9.3).
func TestStoreStockData_CacheDisabled(t *testing.T) {
	repo := NewCacheRepositoryWithConfig(1, false)
	data := validStockData("ITUB4")

	if err := repo.StoreStockData(data); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Even though we stored, cache is disabled so HasValidCache must return false.
	hit, err := repo.HasValidCache("ITUB4")
	if err != nil {
		t.Fatalf("HasValidCache returned error: %v", err)
	}
	if hit {
		t.Error("expected cache miss when cache is disabled")
	}
}

// ─── HasValidCache ────────────────────────────────────────────────────────────

// TestHasValidCache_Hit verifies that HasValidCache returns true after storing
// valid data.
func TestHasValidCache_Hit(t *testing.T) {
	repo := newTestRepo()
	data := validStockData("BBDC4")

	if err := repo.StoreStockData(data); err != nil {
		t.Fatalf("StoreStockData error: %v", err)
	}

	hit, err := repo.HasValidCache("BBDC4")
	if err != nil {
		t.Fatalf("HasValidCache error: %v", err)
	}
	if !hit {
		t.Error("expected cache hit")
	}
}

// TestHasValidCache_Miss verifies that HasValidCache returns false for a symbol
// that was never stored.
func TestHasValidCache_Miss(t *testing.T) {
	repo := newTestRepo()

	hit, err := repo.HasValidCache("UNKNOWN")
	if err != nil {
		t.Fatalf("HasValidCache error: %v", err)
	}
	if hit {
		t.Error("expected cache miss for unknown symbol")
	}
}

// TestHasValidCache_Disabled verifies that HasValidCache returns false when
// cache is disabled (Req 9.3).
func TestHasValidCache_Disabled(t *testing.T) {
	repo := NewCacheRepositoryWithConfig(1, false)

	hit, err := repo.HasValidCache("PETR4")
	if err != nil {
		t.Fatalf("HasValidCache error: %v", err)
	}
	if hit {
		t.Error("expected false when cache is disabled")
	}
}

// ─── Invalidate ───────────────────────────────────────────────────────────────

// TestInvalidate_ExistingEntry verifies that after invalidation HasValidCache
// returns false.
func TestInvalidate_ExistingEntry(t *testing.T) {
	repo := newTestRepo()
	data := validStockData("WEGE3")

	if err := repo.StoreStockData(data); err != nil {
		t.Fatalf("StoreStockData error: %v", err)
	}

	if err := repo.Invalidate("WEGE3"); err != nil {
		t.Fatalf("Invalidate error: %v", err)
	}

	hit, err := repo.HasValidCache("WEGE3")
	if err != nil {
		t.Fatalf("HasValidCache error: %v", err)
	}
	if hit {
		t.Error("expected cache miss after invalidation")
	}
}

// TestInvalidate_NonExistentEntry verifies that invalidating a non-existent
// entry does not return an error.
func TestInvalidate_NonExistentEntry(t *testing.T) {
	repo := newTestRepo()

	if err := repo.Invalidate("DOESNOTEXIST"); err != nil {
		t.Errorf("expected no error invalidating non-existent entry, got: %v", err)
	}
}

// ─── GetStockData ─────────────────────────────────────────────────────────────

// TestGetStockData_CacheHit verifies that GetStockData returns cached data
// without calling the scraper.
func TestGetStockData_CacheHit(t *testing.T) {
	repo := newTestRepo()
	data := validStockData("RENT3")

	if err := repo.StoreStockData(data); err != nil {
		t.Fatalf("StoreStockData error: %v", err)
	}

	scraper := &mockScraper{}
	repo.SetScraper(scraper)

	result, err := repo.GetStockData("RENT3")
	if err != nil {
		t.Fatalf("GetStockData error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Symbol != "RENT3" {
		t.Errorf("expected symbol RENT3, got %s", result.Symbol)
	}
	if scraper.callCount != 0 {
		t.Errorf("expected scraper not to be called on cache hit, called %d times", scraper.callCount)
	}
}

// TestGetStockData_NoScraper verifies that GetStockData returns nil (no error)
// when there is no scraper configured and no cached data.
func TestGetStockData_NoScraper(t *testing.T) {
	repo := newTestRepo()
	// No scraper set.

	result, err := repo.GetStockData("MGLU3")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result when no scraper configured, got: %+v", result)
	}
}

// TestGetStockData_CacheDisabled verifies that when cache is disabled
// GetStockData always calls the scraper.
func TestGetStockData_CacheDisabled(t *testing.T) {
	repo := NewCacheRepositoryWithConfig(1, false)
	freshData := validStockData("ABEV3")
	scraper := &mockScraper{data: freshData}
	repo.SetScraper(scraper)

	result, err := repo.GetStockData("ABEV3")
	if err != nil {
		t.Fatalf("GetStockData error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result from scraper")
	}
	if scraper.callCount != 1 {
		t.Errorf("expected scraper to be called once, called %d times", scraper.callCount)
	}
}

// TestGetStockData_ScraperError verifies that a scraper error is propagated.
func TestGetStockData_ScraperError(t *testing.T) {
	repo := newTestRepo()
	scraper := &mockScraper{err: errors.New("scrape failed")}
	repo.SetScraper(scraper)

	_, err := repo.GetStockData("FAIL3")
	if err == nil {
		t.Error("expected error when scraper fails, got nil")
	}
}

// ─── TTL Expiration ───────────────────────────────────────────────────────────

// TestTTLExpiration verifies that after the TTL elapses the cache entry is
// treated as a miss (Req 10.4).
//
// We create a repository with a 1-millisecond TTL so the test runs quickly.
func TestTTLExpiration(t *testing.T) {
	// Build a repo with a very short TTL (1 ms).
	ttl := time.Millisecond
	repo := &CacheRepository{
		cache:        gocache.New(ttl, 10*time.Millisecond),
		ttl:          ttl,
		cacheEnabled: true,
	}

	data := validStockData("BBAS3")
	if err := repo.StoreStockData(data); err != nil {
		t.Fatalf("StoreStockData error: %v", err)
	}

	// Confirm it's in cache right after storing.
	hit, err := repo.HasValidCache("BBAS3")
	if err != nil {
		t.Fatalf("HasValidCache error: %v", err)
	}
	if !hit {
		t.Fatal("expected cache hit immediately after storing")
	}

	// Wait for the TTL to expire.
	time.Sleep(10 * time.Millisecond)

	// Now the entry should have expired.
	hit, err = repo.HasValidCache("BBAS3")
	if err != nil {
		t.Fatalf("HasValidCache error after expiry: %v", err)
	}
	if hit {
		t.Error("expected cache miss after TTL expiration")
	}
}
