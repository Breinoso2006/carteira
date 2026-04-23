// Package scoring centraliza a lógica de amplificação e distribuição de pesos
// entre ativos do portfólio. Não possui dependências de outros pacotes internos,
// evitando ciclos de importação.
package scoring

import (
	"fmt"
	"time"
)

// PortfolioEntry represents a single stock entry in the portfolio database.
type PortfolioEntry struct {
	ID                  int64     `json:"id"`
	Ticker              string    `json:"ticker"`
	FundamentalistGrade float64   `json:"fundamentalist_grade"`
	Weight              float64   `json:"weight,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (e *PortfolioEntry) GetFundamentalistGrade() float64 { return e.FundamentalistGrade }
func (e *PortfolioEntry) SetWeight(w float64)             { e.Weight = w }

// Validate checks that the PortfolioEntry has a non-empty ticker and a
// fundamentalist grade in the valid range (0, 100].
func (e *PortfolioEntry) Validate() error {
	if e.Ticker == "" {
		return fmt.Errorf("portfolio entry: ticker must not be empty")
	}
	if e.FundamentalistGrade <= 0 || e.FundamentalistGrade > 100 {
		return fmt.Errorf("portfolio entry: fundamentalist_grade must be between 0 and 100, got %.2f", e.FundamentalistGrade)
	}
	return nil
}

// BoostedGrade amplifica a nota usando potenciação quadrática, aumentando a
// diferença relativa entre ativos de notas próximas de forma contínua.
// Ex: nota 80 → 6400, nota 60 → 3600 (diferença de ~78% vs ~22%).
func BoostedGrade(grade float64) float64 {
	return grade * grade
}

// CalculateWeights distribui os pesos entre as entradas usando amplificação
// quadrática sobre FundamentalistGrade. Entradas com total zero são deixadas em 0.
func CalculateWeights(entries []*PortfolioEntry) {
	var total float64
	for _, e := range entries {
		total += BoostedGrade(e.FundamentalistGrade)
	}

	if total == 0 {
		return
	}

	for _, e := range entries {
		e.Weight = BoostedGrade(e.FundamentalistGrade) / total * 100
	}
}

// CalculateWeightsFromGrades distribui os pesos usando notas finais externas
// (ex: nota fundamentalista + momento), em vez de FundamentalistGrade.
// grades deve ter o mesmo comprimento que entries.
func CalculateWeightsFromGrades(entries []*PortfolioEntry, grades []float64) {
	var total float64
	for _, g := range grades {
		total += BoostedGrade(g)
	}

	if total == 0 {
		return
	}

	for i, e := range entries {
		e.Weight = BoostedGrade(grades[i]) / total * 100
	}
}
