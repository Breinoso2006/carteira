package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/breinoso2006/scraping-api/internal/models"
	"github.com/patrickmn/go-cache"
)

// StockScraper is the interface for fetching fresh stock data.
// Defined here to avoid an import cycle between cache and scraping packages.
type StockScraper interface {
	SearchStockInformation(symbol string) (*models.StockData, error)
}

// CacheRepository handles caching of stock data with TTL support
type CacheRepository struct {
	cache        *cache.Cache
	ttl          time.Duration
	cacheEnabled bool
	scraper      StockScraper // used for fresh scrapes on cache miss; may be nil
}

// NewCacheRepository creates a new cache repository with the specified TTL.
// Pass a non-nil scraper to enable automatic fresh-scrape on cache miss.
// cacheEnabled defaults to true; use NewCacheRepositoryWithConfig to control it.
func NewCacheRepository(ttlHours int) *CacheRepository {
	ttl := time.Duration(ttlHours) * time.Hour
	return &CacheRepository{
		cache:        cache.New(ttl, 24*time.Hour),
		ttl:          ttl,
		cacheEnabled: true,
	}
}

// NewCacheRepositoryWithConfig creates a new cache repository with the specified TTL
// and cache-enabled flag (Req 9.3).
func NewCacheRepositoryWithConfig(ttlHours int, cacheEnabled bool) *CacheRepository {
	ttl := time.Duration(ttlHours) * time.Hour
	return &CacheRepository{
		cache:        cache.New(ttl, 24*time.Hour),
		ttl:          ttl,
		cacheEnabled: cacheEnabled,
	}
}

// SetScraper injects the scraper used to fetch fresh data on cache miss.
func (r *CacheRepository) SetScraper(s StockScraper) {
	r.scraper = s
}

// GetStockData returns stock data for the given symbol.
//
// Behaviour (Requirements 2.2, 6.1, 6.3, 9.3):
//  1. If cache is disabled, skip cache lookup and go straight to fresh scrape.
//  2. If a valid, non-expired cache entry exists it is returned immediately.
//  3. On a cache miss the method triggers a fresh scrape via the injected
//     StockScraper, stores the result (when valid), and returns it.
//  4. The returned *StockData has the same structure regardless of source.
func (r *CacheRepository) GetStockData(symbol string) (*models.StockData, error) {
	// When cache is disabled, bypass cache entirely (Req 9.3 / Req 2.7).
	if r.cacheEnabled {
		// 1. Check for a valid, non-expired cache entry.
		valid, err := r.HasValidCache(symbol)
		if err != nil {
			// Log and fall through to fresh scrape.
			fmt.Printf("cache validity check error for %s: %v\n", symbol, err)
		}

		if valid {
			item, found := r.cache.Get(symbol)
			if found {
				if data, ok := item.(*models.StockData); ok {
					// Cache hit – return as-is (Req 6.1, 6.3: same format, no cache indicator).
					return data, nil
				}
			}
		}
	}

	// 2. Cache miss (or cache disabled) – trigger a fresh scrape if a scraper is available.
	if r.scraper == nil {
		// No scraper configured; return nil so the caller can handle the miss.
		return nil, nil
	}

	freshData, scrapeErr := r.scraper.SearchStockInformation(symbol)
	if scrapeErr != nil {
		return nil, fmt.Errorf("cache miss and scrape failed for %s: %w", symbol, scrapeErr)
	}

	// 3. Store the fresh result in cache (StoreStockData skips invalid data per Req 2.4,
	//    and is a no-op when cache is disabled per Req 9.3).
	if storeErr := r.StoreStockData(freshData); storeErr != nil {
		fmt.Printf("cache store error for %s: %v\n", symbol, storeErr)
	}

	// Return fresh data in the same format (Req 6.1, 6.3).
	return freshData, nil
}

// StoreStockData stores stock data in cache with TTL (Req 2.1, 2.4, 9.3).
//
// Behaviour:
//   - Returns nil immediately when cache is disabled (Req 9.3).
//   - Returns an error for nil data.
//   - Silently skips (returns nil) when data contains any invalid fields (Req 2.4).
//   - Sets CreatedAt/UpdatedAt timestamps and stores with the configured TTL (Req 2.1).
func (r *CacheRepository) StoreStockData(data *models.StockData) error {
	// No-op when cache is disabled (Req 9.3).
	if !r.cacheEnabled {
		return nil
	}

	if data == nil {
		return fmt.Errorf("cannot store nil stock data")
	}

	// Don't cache data with invalid fields (Req 2.4).
	if len(data.GetInvalidFields()) > 0 {
		return nil
	}

	// Stamp timestamps before storing (Req 2.1).
	now := time.Now()
	if data.CreatedAt.IsZero() {
		data.CreatedAt = now
	}
	data.UpdatedAt = now

	r.cache.Set(data.Symbol, data, r.ttl)
	return nil
}

// HasValidCache checks if valid cached data exists for the given symbol.
// Returns false immediately when cache is disabled (Req 9.3).
func (r *CacheRepository) HasValidCache(symbol string) (bool, error) {
	if !r.cacheEnabled {
		return false, nil
	}

	item, found := r.cache.Get(symbol)
	if !found {
		return false, nil
	}

	_, ok := item.(*models.StockData)
	if !ok {
		return false, fmt.Errorf("invalid cache item type for symbol %s", symbol)
	}

	return true, nil
}

// Invalidate removes a cache entry for the given symbol
func (r *CacheRepository) Invalidate(symbol string) error {
	r.cache.Delete(symbol)
	return nil
}

// Refresh invalidates the existing cache entry, fetches fresh data from the scraper,
// stores it in cache, and returns it (Req 4.3).
//
// Behaviour:
//   - Always invalidates the existing cache entry first.
//   - If no scraper is configured, returns nil data with no error.
//   - If the scraper returns an error, the error is propagated to the caller.
//   - Fresh data is stored in cache when valid (same rules as StoreStockData).
func (r *CacheRepository) Refresh(symbol string) (*models.StockData, error) {
	// Step 1: Invalidate the existing entry (Req 4.3).
	if err := r.Invalidate(symbol); err != nil {
		return nil, fmt.Errorf("refresh: failed to invalidate cache for %s: %w", symbol, err)
	}

	// Step 2: If no scraper is available, return gracefully.
	if r.scraper == nil {
		return nil, nil
	}

	// Step 3: Fetch fresh data from the scraper.
	freshData, err := r.scraper.SearchStockInformation(symbol)
	if err != nil {
		return nil, fmt.Errorf("refresh: scrape failed for %s: %w", symbol, err)
	}

	// Step 4: Store the fresh data in cache (skips invalid data per Req 2.4).
	if storeErr := r.StoreStockData(freshData); storeErr != nil {
		fmt.Printf("refresh: cache store error for %s: %v\n", symbol, storeErr)
	}

	return freshData, nil
}

// GetInvalidFieldsFromCache retrieves invalid fields for a symbol from cache
func (r *CacheRepository) GetInvalidFieldsFromCache(symbol string) ([]string, error) {
	item, found := r.cache.Get(symbol)
	if !found {
		return nil, nil
	}

	data, ok := item.(*models.StockData)
	if !ok {
		return nil, fmt.Errorf("invalid cache item type for symbol %s", symbol)
	}

	return data.GetInvalidFields(), nil
}

// StoreWithInvalidFields stores stock data with invalid fields information
func (r *CacheRepository) StoreWithInvalidFields(data *models.StockData) error {
	if data == nil {
		return fmt.Errorf("cannot store nil stock data")
	}

	// Store even if there are invalid fields (for tracking purposes)
	r.cache.Set(data.Symbol, data, r.ttl)
	return nil
}

// MarshalInvalidFields converts invalid fields map to JSON string
func MarshalInvalidFields(fields map[string]bool) (string, error) {
	if len(fields) == 0 {
		return "[]", nil
	}

	data, err := json.Marshal(fields)
	if err != nil {
		return "", fmt.Errorf("failed to marshal invalid fields: %w", err)
	}
	return string(data), nil
}

// UnmarshalInvalidFields converts JSON string to invalid fields map
func UnmarshalInvalidFields(jsonStr string) (map[string]bool, error) {
	if jsonStr == "" || jsonStr == "[]" {
		return make(map[string]bool), nil
	}

	fields := make(map[string]bool)
	if err := json.Unmarshal([]byte(jsonStr), &fields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal invalid fields: %w", err)
	}
	return fields, nil
}

// GetStats returns the number of items in cache
func (r *CacheRepository) GetStats() int {
	return r.cache.ItemCount()
}

// GetTTL returns the current TTL
func (r *CacheRepository) GetTTL() time.Duration {
	return r.ttl
}
