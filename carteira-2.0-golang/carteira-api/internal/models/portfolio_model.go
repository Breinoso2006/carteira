package models

import (
	"fmt"
	"sync"
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

type StockInPortfolio struct {
	Stock  *Stock
	Weight float64
}

type Portfolio struct {
	Stocks []*StockInPortfolio
}

func (p *Portfolio) CalculateWeights() {
	var wg sync.WaitGroup

	for _, stockInPortfolio := range p.Stocks {
		wg.Add(1)

		go func(sip *StockInPortfolio) {
			defer wg.Done()
			sip.Stock.SetFinalGrade()
		}(stockInPortfolio)
	}

	wg.Wait()

	totalGrade := 0.0
	for _, s := range p.Stocks {
		totalGrade += s.Stock.FinalGrade
	}

	for _, s := range p.Stocks {
		s.Weight = s.Stock.FinalGrade / totalGrade * 100
	}

	for _, s := range p.Stocks {
		fmt.Printf("Stock: %s, Final Grade: %.2f, Weight: %.2f%%\n", s.Stock.Ticker, s.Stock.FinalGrade, s.Weight)
	}
}
