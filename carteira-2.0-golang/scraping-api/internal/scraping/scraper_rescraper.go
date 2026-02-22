package scraper

import (
	"bytes"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/breinoso2006/scraping-api/internal/http"
)

// RescrapeSingleField reextrai um campo específico de uma ação
// Útil para preencher campos que falharam em outras fontes
func (s *SourceScraper) RescrapeSingleField(symbol string, fieldName string) (*float64, error) {
	client := http.NewHTTPClient(5 * time.Second)
	body, err := client.Get(s.urlPrefix + symbol)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Busca o seletor e a função de limpeza
	selector, selectorExists := s.selectors[fieldName]
	if !selectorExists {
		return nil, fmt.Errorf("campo %s não existe na fonte %s", fieldName, s.source)
	}

	cleaner, cleanerExists := s.cleaners[fieldName]
	if !cleanerExists {
		return nil, fmt.Errorf("sem função de limpeza para %s", fieldName)
	}

	// Extrai o valor
	text := doc.Find(selector).First().Text()
	text = cleaner(text)

	// Converte para float64
	value := parseFloatPointer(text)
	if value == nil {
		return nil, fmt.Errorf("valor inválido para %s (texto: '%s')", fieldName, text)
	}

	return value, nil
}

