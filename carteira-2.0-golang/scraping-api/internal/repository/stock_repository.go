package repository

import (
	"time"

	"github.com/breinoso2006/scraping-api/internal/cache"
	"github.com/breinoso2006/scraping-api/internal/models"
	scraper "github.com/breinoso2006/scraping-api/internal/scraping"
)

type StockRepository struct {
	cacheRepo *cache.CacheRepository
}

var instance *StockRepository

// NewStockRepository creates a StockRepository with the given TTL and cache-enabled flag.
// This is the preferred constructor; it wires the ScraperManager into the CacheRepository
// so that cache misses trigger a fresh scrape automatically (Req 9.2, 9.3).
func NewStockRepository(ttlHours int, cacheEnabled bool) *StockRepository {
	manager := scraper.NewScraperManager()
	cacheRepo := cache.NewCacheRepositoryWithConfig(ttlHours, cacheEnabled)
	cacheRepo.SetScraper(manager)
	return &StockRepository{
		cacheRepo: cacheRepo,
	}
}

func GetStockRepository() *StockRepository {
	return instance
}

// GetStockData returns stock data for the given symbol.
// It delegates entirely to CacheRepository, which checks the cache first and
// falls back to a fresh scrape on a miss (Requirements 2.2, 6.1, 6.3).
func (r *StockRepository) GetStockData(symbol string) (*models.StockData, error) {
	return r.cacheRepo.GetStockData(symbol)
}

func (r *StockRepository) SetCacheTTL(hours int) {
	manager := scraper.NewScraperManager()
	r.cacheRepo = cache.NewCacheRepository(hours)
	r.cacheRepo.SetScraper(manager)
}

func (r *StockRepository) GetCacheStats() (int, time.Duration) {
	items := r.cacheRepo.GetStats()
	return items, r.cacheRepo.GetTTL()
}
