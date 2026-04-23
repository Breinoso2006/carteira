package cache

import (
	"fmt"
	"testing"
	"time"

	gocache "github.com/patrickmn/go-cache"

	"github.com/breinoso2006/scraping-api/internal/models"
)

// newBenchmarkRepo creates a CacheRepository backed by a real in-memory go-cache
// instance, identical to what production code uses.
func newBenchmarkRepo() *CacheRepository {
	ttl := time.Hour
	return &CacheRepository{
		cache:        gocache.New(ttl, 24*time.Hour),
		ttl:          ttl,
		cacheEnabled: true,
	}
}

// benchmarkStockData returns a fully-populated StockData with no invalid fields.
func benchmarkStockData(symbol string) *models.StockData {
	price := 42.50
	pe := 12.0
	pbv := 1.8
	psr := 2.5
	bvps := 23.0
	eps := 3.5
	dy := 0.06
	return &models.StockData{
		Symbol: symbol,
		Price:  &price,
		PE:     &pe,
		PBV:    &pbv,
		PSR:    &psr,
		BVps:   &bvps,
		EPS:    &eps,
		DY:     &dy,
		Source: "benchmark",
	}
}

// BenchmarkGetStockData_CacheHit measures the latency of a cache hit.
//
// Requirement 7.1: WHEN valid cached data is available, THE CacheRepository
// SHALL return it within 50ms.
//
// go-cache.Get() is an O(1) in-memory map lookup protected by a read-lock,
// so each iteration should complete in well under 1 µs – orders of magnitude
// below the 50 ms budget.
func BenchmarkGetStockData_CacheHit(b *testing.B) {
	repo := newBenchmarkRepo()
	data := benchmarkStockData("PETR4")

	if err := repo.StoreStockData(data); err != nil {
		b.Fatalf("setup: StoreStockData failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := repo.GetStockData("PETR4")
		if err != nil {
			b.Fatalf("GetStockData error: %v", err)
		}
		if result == nil {
			b.Fatal("expected non-nil result on cache hit")
		}
	}
}

// BenchmarkStoreStockData measures the latency of writing a valid entry to cache.
//
// go-cache.Set() is an O(1) in-memory map write protected by a write-lock,
// so each iteration should complete in well under 1 µs.
func BenchmarkStoreStockData(b *testing.B) {
	repo := newBenchmarkRepo()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use a unique symbol per iteration to avoid overwriting the same key
		// and to exercise the full write path each time.
		symbol := fmt.Sprintf("TICK%d", i)
		data := benchmarkStockData(symbol)
		if err := repo.StoreStockData(data); err != nil {
			b.Fatalf("StoreStockData error: %v", err)
		}
	}
}

// BenchmarkHasValidCache measures the latency of a cache-presence check.
//
// This is the hot path called before every GetStockData; it must be fast.
func BenchmarkHasValidCache(b *testing.B) {
	repo := newBenchmarkRepo()
	data := benchmarkStockData("VALE3")

	if err := repo.StoreStockData(data); err != nil {
		b.Fatalf("setup: StoreStockData failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hit, err := repo.HasValidCache("VALE3")
		if err != nil {
			b.Fatalf("HasValidCache error: %v", err)
		}
		if !hit {
			b.Fatal("expected cache hit")
		}
	}
}
