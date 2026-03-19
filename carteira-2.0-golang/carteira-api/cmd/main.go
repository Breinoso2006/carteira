package main

import (
	"github.com/breinoso2006/carteira-api/internal/models"
)

func main() {
	stocks := []*models.StockInPortfolio{
		{Stock: models.NewStock("ALUP3", 77.5)},
		{Stock: models.NewStock("BBSE3", 77.5)},
		{Stock: models.NewStock("BMEB4", 75)},
		{Stock: models.NewStock("CAML3", 62.5)},
		{Stock: models.NewStock("CSAN3", 62.5)},
		{Stock: models.NewStock("EGIE3", 85)},
		{Stock: models.NewStock("FESA4", 70)},
		{Stock: models.NewStock("FLRY3", 70)},
		{Stock: models.NewStock("ITSA4", 90)},
		{Stock: models.NewStock("KLBN3", 70)},
		{Stock: models.NewStock("PRIO3", 75)},
		{Stock: models.NewStock("SUZB3", 72.5)},
		{Stock: models.NewStock("TAEE3", 70)},
		{Stock: models.NewStock("TUPY3", 70)},
		{Stock: models.NewStock("UNIP6", 80)},
		{Stock: models.NewStock("VIVT3", 70)},
		{Stock: models.NewStock("WEGE3", 98.75)},
		{Stock: models.NewStock("WIZC3", 75)},
	}

	portfolio := &models.Portfolio{Stocks: stocks}
	portfolio.CalculateWeights()
}
