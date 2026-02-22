package repository

import (
	"fmt"

	"github.com/breinoso2006/scraping-api/internal/models"
	scraper "github.com/breinoso2006/scraping-api/internal/scraping"
)

type StockRepository struct {
	scrapers []scraper.StockScraper
}

var instance *StockRepository

func init() {
	instance = &StockRepository{
		scrapers: scraper.GetScrapers(),
	}
}

func GetStockRepository() *StockRepository {
	return instance
}

func (r *StockRepository) GetStockData(symbol string) (*models.StockData, error) {
	var result *models.StockData
	var errors []string

	for _, s := range r.scrapers {
		data, err := s.SearchStockInformation(symbol)
		if err != nil {
			errors = append(errors, fmt.Sprintf("scraper failed: %v", err))
			continue
		}

		result = r.mergeStockData(result, data)

		if r.isComplete(result) {
			return result, nil
		}
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("falha em todos os scrapers: %v", errors)
}

func (r *StockRepository) mergeStockData(existing, new *models.StockData) *models.StockData {
	if existing == nil {
		return new
	}

	if existing.Price == nil && new.Price != nil {
		existing.Price = new.Price
	}
	if existing.PE == nil && new.PE != nil {
		existing.PE = new.PE
	}
	if existing.PSR == nil && new.PSR != nil {
		existing.PSR = new.PSR
	}
	if existing.BVps == nil && new.BVps != nil {
		existing.BVps = new.BVps
	}
	if existing.EPS == nil && new.EPS != nil {
		existing.EPS = new.EPS
	}

	return existing
}

func (r *StockRepository) isComplete(data *models.StockData) bool {
	return data.Price != nil && data.PE != nil && data.PSR != nil && data.BVps != nil && data.EPS != nil
}
