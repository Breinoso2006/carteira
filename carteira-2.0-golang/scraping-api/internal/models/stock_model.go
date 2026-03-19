package models

type StockData struct {
	Symbol string
	Price  *float64 // Preço da Ação
	PE     *float64 // Price to Earnings (Preço/Lucro)
	PBV		*float64 // Price to Book Value (Preço/Valor Patrimonial)
	PSR    *float64 // Price to Sales Ratio (Preço/Vendas)
	BVps   *float64 // Book Value Per Share (Valor Patrimonial Por Ação)
	EPS    *float64 // Earnings Per Share (Lucro Por Ação)
	DY     *float64 // Dividend Yield (Rendimento de Dividendos) em %
	Source string   // Fonte dos dados (investidor10, auvp, fundamentus)

	// Interno: quais campos estão vazios/inválidos
	invalidFields map[string]bool
}

func (s *StockData) GetInvalidFields() []string {
	if s.invalidFields == nil {
		return []string{}
	}
	invalid := []string{}
	for field := range s.invalidFields {
		invalid = append(invalid, field)
	}
	return invalid
}

func (s *StockData) MarkFieldInvalid(field string) {
	if s.invalidFields == nil {
		s.invalidFields = make(map[string]bool)
	}
	s.invalidFields[field] = true
}

func (s *StockData) ClearInvalidFields() {
	s.invalidFields = make(map[string]bool)
}
