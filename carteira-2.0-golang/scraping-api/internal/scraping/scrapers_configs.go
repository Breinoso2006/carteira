package scraper

// GetScrapers retorna lista de todos os scrapers disponíveis
// Mantida aqui para compatibilidade com código antigo
func GetScrapers() []StockScraper {
	configs := GetSourceConfigs()
	scrapers := make([]StockScraper, 0, len(configs))

	for _, source := range []string{"investidor10", "auvp", "fundamentus"} {
		if config, exists := configs[source]; exists {
			scrapers = append(scrapers, NewSourceScraper(config))
		}
	}

	return scrapers
}
