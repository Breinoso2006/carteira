package scraper

import (
	"fmt"

	"github.com/breinoso2006/scraping-api/internal/models"
)

type ScraperManager struct {
	scrapers map[string]StockScraper
	order    []string // Ordem de prioridade para fallback
}

func NewScraperManager() *ScraperManager {
	configs := GetSourceConfigs()

	// Cria scrapers e mantém ordem de fallback
	scrapers := make(map[string]StockScraper)
	order := []string{"investidor10", "auvp", "fundamentus"}

	for _, source := range order {
		if config, exists := configs[source]; exists {
			scrapers[source] = NewSourceScraper(config)
		}
	}

	return &ScraperManager{
		scrapers: scrapers,
		order:    order,
	}
}

// SearchStockInformation busca dados da ação usando fallback inteligente
// Tenta na ordem: Investidor10 → Auvp → Fundamentus
// Para campos que falharam, reexperimenta apenas desse campo em outra fonte
func (sm *ScraperManager) SearchStockInformation(symbol string) (*models.StockData, error) {
	var lastError error
	var bestData *models.StockData

	for _, source := range sm.order {
		scraper, exists := sm.scrapers[source]
		if !exists {
			continue
		}

		data, err := scraper.SearchStockInformation(symbol)
		if err != nil {
			lastError = err
			continue
		}

		// Primeira fonte com algum dado válido
		if bestData == nil {
			bestData = data
		} else {
			// Para fontes subsequentes, tenta preencher apenas campos que faltam
			sm.mergeData(bestData, data)
		}

		// Se conseguiu dados completos, retorna
		if sm.isCompleteData(bestData) {
			bestData.ClearInvalidFields() // Remove marcas internas
			return bestData, nil
		}
	}

	// Se tem dados parciais, tenta preencher os campos inválidos em outras fontes
	if bestData != nil && len(bestData.GetInvalidFields()) > 0 {
		sm.fillMissingFieldsFromOtherSources(bestData, symbol)
		bestData.ClearInvalidFields() // Remove marcas internas antes de retornar
		return bestData, nil
	}

	// Se tem algum dado, retorna mesmo que incompleto
	if bestData != nil {
		bestData.ClearInvalidFields()
		return bestData, nil
	}

	// Sem dados de nenhuma fonte
	if lastError != nil {
		return nil, fmt.Errorf("falha em todas as fontes ao buscar %s: %w", symbol, lastError)
	}

	return nil, fmt.Errorf("nenhuma fonte conseguiu extrair dados para %s", symbol)
}

// mergeData mescla dados de uma nova fonte na melhor fonte encontrada
// Preenche apenas campos vazios, mantendo dados da fonte anterior
func (sm *ScraperManager) mergeData(target *models.StockData, source *models.StockData) {
	if target.Price == nil && source.Price != nil {
		target.Price = source.Price
	}
	if target.PE == nil && source.PE != nil {
		target.PE = source.PE
	}
	if target.PBV == nil && source.PBV != nil {
		target.PBV = source.PBV
	}
	if target.PSR == nil && source.PSR != nil {
		target.PSR = source.PSR
	}
	if target.BVps == nil && source.BVps != nil {
		target.BVps = source.BVps
	}
	if target.EPS == nil && source.EPS != nil {
		target.EPS = source.EPS
	}
	if target.DY == nil && source.DY != nil {
		target.DY = source.DY
	}
}

// fillMissingFieldsFromOtherSources tenta preencher apenas os campos inválidos
// usando o mecanismo de re-scraping de campo único
func (sm *ScraperManager) fillMissingFieldsFromOtherSources(data *models.StockData, symbol string) {
	invalidFields := data.GetInvalidFields()
	if len(invalidFields) == 0 {
		return
	}

	// Para cada campo que falhou, tenta em outras fontes
	for _, invalidField := range invalidFields {
		// Pula a fonte que falhou originalmente
		for _, source := range sm.order {
			if source == data.Source {
				continue
			}

			scraper, exists := sm.scrapers[source]
			if !exists {
				continue
			}

			// Cast para SourceScraper para acessar RescrapeSingleField
			if ss, ok := scraper.(*SourceScraper); ok {
				value, err := ss.RescrapeSingleField(symbol, invalidField)
				if err == nil && value != nil {
					// Preenche o campo baseado no nome
					switch invalidField {
					case "price":
						data.Price = value
					case "pe":
						data.PE = value
					case "pbv":
						data.PBV = value
					case "psr":
						data.PSR = value
					case "bvps":
						data.BVps = value
					case "eps":
						data.EPS = value
					case "dy":
						data.DY = value
					}
					// Campo foi preenchido, não precisa tentar nas outras fontes
					break
				}
			}
		}
	}
}

// SearchStockInformationFromSource busca dados de uma fonte específica
func (sm *ScraperManager) SearchStockInformationFromSource(symbol string, sourceName string) (*models.StockData, error) {
	scraper, exists := sm.scrapers[sourceName]
	if !exists {
		return nil, fmt.Errorf("fonte %s não encontrada", sourceName)
	}
	return scraper.SearchStockInformation(symbol)
}

// GetAvailableSources retorna lista de fontes disponíveis
func (sm *ScraperManager) GetAvailableSources() []string {
	return sm.order
}

func (sm *ScraperManager) isCompleteData(data *models.StockData) bool {
	return data.Price != nil &&
		data.PE != nil &&
		data.PBV != nil &&
		data.PSR != nil &&
		data.BVps != nil &&
		data.EPS != nil
}
