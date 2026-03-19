package scraper

import (
	"bytes"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/breinoso2006/scraping-api/internal/http"
	"github.com/breinoso2006/scraping-api/internal/models"
)

type StockScraper interface {
	SearchStockInformation(symbol string) (*models.StockData, error)
}

// SourceScraper é a implementação única para todas as fontes
type SourceScraper struct {
	source    string                         // "investidor10", "auvp", "fundamentus"
	urlPrefix string                         // URL base da fonte
	selectors map[string]string              // Seletores CSS por campo
	cleaners  map[string]func(string) string // Funções de limpeza por campo
}

// SourceConfig define a configuração de uma fonte
type SourceConfig struct {
	Source    string
	URLPrefix string
	Selectors map[string]string
	Cleaners  map[string]func(string) string
}

func NewSourceScraper(config SourceConfig) *SourceScraper {
	return &SourceScraper{
		source:    config.Source,
		urlPrefix: config.URLPrefix,
		selectors: config.Selectors,
		cleaners:  config.Cleaners,
	}
}

func (s *SourceScraper) SearchStockInformation(symbol string) (*models.StockData, error) {
	client := http.NewHTTPClient(5 * time.Second)
	body, err := client.Get(s.urlPrefix + symbol)
	if err != nil {
		return nil, err
	}

	return s.scrapeStockData(symbol, body)
}

func (s *SourceScraper) scrapeStockData(symbol string, body []byte) (*models.StockData, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	data := &models.StockData{
		Symbol: symbol,
		Source: s.source,
	}

	// Campos para validar
	fieldConfigs := []struct {
		name      string
		target    **float64
		fieldName string
	}{
		{"price", &data.Price, "price"},
		{"pe", &data.PE, "pe"},
		{"pbv", &data.PBV, "pbv"},
		{"psr", &data.PSR, "psr"},
		{"bvps", &data.BVps, "bvps"},
		{"eps", &data.EPS, "eps"},
		{"dy", &data.DY, "dy"},
	}

	// Extrai cada campo
	for _, fc := range fieldConfigs {
		selector, exists := s.selectors[fc.name]
		if !exists {
			continue
		}

		text := doc.Find(selector).First().Text()

		if cleaner, hasCleanup := s.cleaners[fc.name]; hasCleanup {
			text = cleaner(text)
		}

		value := parseFloatPointer(text)
		if value == nil {
			data.MarkFieldInvalid(fc.fieldName)
		} else {
			*fc.target = value
		}
	}

	return data, nil
}
