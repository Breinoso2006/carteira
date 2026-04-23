package cache

// Integration tests for CacheRepository covering hit/miss cycles, scraper
// fallback, cache invalidation, refresh, and response consistency.
//
// Validates: Requirements 10.3, 10.5

import (
	"errors"
	"testing"

	"github.com/breinoso2006/scraping-api/internal/models"
)

// TestCacheIntegration_HitMissCycle verifies the full hit → invalidate → miss
// cycle:
//  1. Store data → GetStockData returns cached data (scraper NOT called).
//  2. Invalidate → GetStockData triggers scraper (cache miss).
func TestCacheIntegration_HitMissCycle(t *testing.T) {
	repo := newTestRepo()
	cached := validStockData("PETR4")

	if err := repo.StoreStockData(cached); err != nil {
		t.Fatalf("StoreStockData: %v", err)
	}

	freshData := validStockData("PETR4")
	price := 99.0
	freshData.Price = &price
	scraper := &mockScraper{data: freshData}
	repo.SetScraper(scraper)

	// Cache HIT
	result, err := repo.GetStockData("PETR4")
	if err != nil {
		t.Fatalf("GetStockData (hit): %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result on cache hit")
	}
	if scraper.callCount != 0 {
		t.Errorf("scraper should not be called on cache hit; called %d time(s)", scraper.callCount)
	}

	// Invalidate
	if err := repo.Invalidate("PETR4"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}

	hit, err := repo.HasValidCache("PETR4")
	if err != nil {
		t.Fatalf("HasValidCache after invalidate: %v", err)
	}
	if hit {
		t.Error("expected cache miss after invalidation")
	}

	// Cache MISS → scraper called
	result, err = repo.GetStockData("PETR4")
	if err != nil {
		t.Fatalf("GetStockData (miss): %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result from scraper on cache miss")
	}
	if scraper.callCount != 1 {
		t.Errorf("expected scraper to be called once on cache miss; called %d time(s)", scraper.callCount)
	}
}

// TestCacheIntegration_ScraperFallback verifies that when there is no cache
// entry the scraper is called and the returned data is subsequently cached.
func TestCacheIntegration_ScraperFallback(t *testing.T) {
	repo := newTestRepo()
	freshData := validStockData("VALE3")
	scraper := &mockScraper{data: freshData}
	repo.SetScraper(scraper)

	hit, err := repo.HasValidCache("VALE3")
	if err != nil {
		t.Fatalf("HasValidCache (pre): %v", err)
	}
	if hit {
		t.Fatal("expected no cache entry before first request")
	}

	result, err := repo.GetStockData("VALE3")
	if err != nil {
		t.Fatalf("GetStockData: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result from scraper")
	}
	if scraper.callCount != 1 {
		t.Errorf("expected scraper called once; called %d time(s)", scraper.callCount)
	}

	// Data should now be cached.
	hit, err = repo.HasValidCache("VALE3")
	if err != nil {
		t.Fatalf("HasValidCache (post): %v", err)
	}
	if !hit {
		t.Error("expected data to be cached after scraper fallback")
	}

	// Second call should be served from cache.
	_, err = repo.GetStockData("VALE3")
	if err != nil {
		t.Fatalf("GetStockData (second call): %v", err)
	}
	if scraper.callCount != 1 {
		t.Errorf("expected scraper still called only once after cache population; called %d time(s)", scraper.callCount)
	}
}

// TestCacheIntegration_ScraperFallbackInvalidData verifies that when the
// scraper returns data with invalid fields the data is NOT stored in cache.
func TestCacheIntegration_ScraperFallbackInvalidData(t *testing.T) {
	repo := newTestRepo()

	invalidData := validStockData("MGLU3")
	invalidData.MarkFieldInvalid("Price")

	scraper := &mockScraper{data: invalidData}
	repo.SetScraper(scraper)

	result, err := repo.GetStockData("MGLU3")
	if err != nil {
		t.Fatalf("GetStockData: %v", err)
	}
	if result == nil {
		t.Fatal("expected scraper result to be returned even when invalid")
	}
	if scraper.callCount != 1 {
		t.Errorf("expected scraper called once; called %d time(s)", scraper.callCount)
	}

	// Invalid data must NOT be cached.
	hit, err := repo.HasValidCache("MGLU3")
	if err != nil {
		t.Fatalf("HasValidCache: %v", err)
	}
	if hit {
		t.Error("invalid data should NOT be cached")
	}

	// A subsequent call must hit the scraper again.
	_, err = repo.GetStockData("MGLU3")
	if err != nil {
		t.Fatalf("GetStockData (second call): %v", err)
	}
	if scraper.callCount != 2 {
		t.Errorf("expected scraper called twice (no cache); called %d time(s)", scraper.callCount)
	}
}

// TestCacheIntegration_RefreshCycle verifies that Refresh invalidates the
// existing entry, calls the scraper, and stores the new data.
func TestCacheIntegration_RefreshCycle(t *testing.T) {
	repo := newTestRepo()

	initial := validStockData("WEGE3")
	if err := repo.StoreStockData(initial); err != nil {
		t.Fatalf("StoreStockData: %v", err)
	}

	fresh := validStockData("WEGE3")
	newPrice := 55.5
	fresh.Price = &newPrice
	scraper := &mockScraper{data: fresh}
	repo.SetScraper(scraper)

	refreshed, err := repo.Refresh("WEGE3")
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if refreshed == nil {
		t.Fatal("expected non-nil result from Refresh")
	}
	if scraper.callCount != 1 {
		t.Errorf("expected scraper called once during Refresh; called %d time(s)", scraper.callCount)
	}
	if refreshed.Price == nil || *refreshed.Price != newPrice {
		t.Errorf("expected refreshed price %.1f, got %v", newPrice, refreshed.Price)
	}

	hit, err := repo.HasValidCache("WEGE3")
	if err != nil {
		t.Fatalf("HasValidCache after Refresh: %v", err)
	}
	if !hit {
		t.Error("expected new data to be cached after Refresh")
	}

	// GetStockData should now return the refreshed data without calling the scraper.
	result, err := repo.GetStockData("WEGE3")
	if err != nil {
		t.Fatalf("GetStockData after Refresh: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result after Refresh")
	}
	if scraper.callCount != 1 {
		t.Errorf("scraper should not be called again after Refresh cached the data; called %d time(s)", scraper.callCount)
	}
}

// TestCacheIntegration_CachePreservationOnScraperFailure verifies that when
// the scraper fails the original cache entry is preserved (Req 4.1).
func TestCacheIntegration_CachePreservationOnScraperFailure(t *testing.T) {
	repo := newTestRepo()

	original := validStockData("BBDC4")
	if err := repo.StoreStockData(original); err != nil {
		t.Fatalf("StoreStockData: %v", err)
	}

	failingScraper := &mockScraper{err: errors.New("network error")}
	repo.SetScraper(failingScraper)

	// GetStockData should return the cached entry (not call the scraper).
	result, err := repo.GetStockData("BBDC4")
	if err != nil {
		t.Fatalf("GetStockData: %v", err)
	}
	if result == nil {
		t.Fatal("expected cached result to be returned")
	}
	if failingScraper.callCount != 0 {
		t.Errorf("scraper should not be called when cache hit exists; called %d time(s)", failingScraper.callCount)
	}

	hit, err := repo.HasValidCache("BBDC4")
	if err != nil {
		t.Fatalf("HasValidCache: %v", err)
	}
	if !hit {
		t.Error("original cache entry should still exist after scraper failure")
	}
}

// TestCacheIntegration_ResponseConsistency verifies that GetStockData returns
// the same structure whether data comes from cache or from the scraper.
func TestCacheIntegration_ResponseConsistency(t *testing.T) {
	assertStockDataShape := func(t *testing.T, label string, data *models.StockData) {
		t.Helper()
		if data == nil {
			t.Fatalf("[%s] expected non-nil StockData", label)
		}
		if data.Symbol == "" {
			t.Errorf("[%s] Symbol should not be empty", label)
		}
		if data.Price == nil {
			t.Errorf("[%s] Price should not be nil", label)
		}
		if data.PE == nil {
			t.Errorf("[%s] PE should not be nil", label)
		}
		if data.PBV == nil {
			t.Errorf("[%s] PBV should not be nil", label)
		}
		if data.PSR == nil {
			t.Errorf("[%s] PSR should not be nil", label)
		}
		if data.BVps == nil {
			t.Errorf("[%s] BVps should not be nil", label)
		}
		if data.EPS == nil {
			t.Errorf("[%s] EPS should not be nil", label)
		}
		if data.DY == nil {
			t.Errorf("[%s] DY should not be nil", label)
		}
		if data.Source == "" {
			t.Errorf("[%s] Source should not be empty", label)
		}
	}

	// From scraper (cache miss)
	repoMiss := newTestRepo()
	scraperData := validStockData("ITUB4")
	scraper := &mockScraper{data: scraperData}
	repoMiss.SetScraper(scraper)

	fromScraper, err := repoMiss.GetStockData("ITUB4")
	if err != nil {
		t.Fatalf("GetStockData (scraper): %v", err)
	}
	assertStockDataShape(t, "from scraper", fromScraper)

	// From cache (cache hit)
	repoHit := newTestRepo()
	cachedData := validStockData("ITUB4")
	if err := repoHit.StoreStockData(cachedData); err != nil {
		t.Fatalf("StoreStockData: %v", err)
	}

	fromCache, err := repoHit.GetStockData("ITUB4")
	if err != nil {
		t.Fatalf("GetStockData (cache): %v", err)
	}
	assertStockDataShape(t, "from cache", fromCache)

	if fromScraper.Symbol != fromCache.Symbol {
		t.Errorf("symbol mismatch: scraper=%s cache=%s", fromScraper.Symbol, fromCache.Symbol)
	}
}
