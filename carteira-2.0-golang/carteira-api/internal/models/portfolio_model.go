package models

import (
	"fmt"
	"sync"
)

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
