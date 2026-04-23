package models

import (
	"fmt"
	"sync"

	"github.com/breinoso2006/carteira-api/internal/scoring"
	"github.com/breinoso2006/carteira-api/internal/stock"
)

// PortfolioEntry é um alias para scoring.PortfolioEntry, mantendo
// compatibilidade com código que importa o tipo via models.
type PortfolioEntry = scoring.PortfolioEntry

type StockInPortfolio struct {
	Stock  *stock.Stock
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
		totalGrade += scoring.BoostedGrade(s.Stock.FinalGrade)
	}

	if totalGrade == 0 {
		return
	}

	for _, s := range p.Stocks {
		s.Weight = scoring.BoostedGrade(s.Stock.FinalGrade) / totalGrade * 100
	}

	for _, s := range p.Stocks {
		fmt.Printf("Stock: %s, Final Grade: %.2f, Weight: %.2f%%\n",
			s.Stock.Ticker, s.Stock.FinalGrade, s.Weight)
	}
}
