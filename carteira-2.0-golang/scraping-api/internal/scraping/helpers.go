package scraper

import (
	"strconv"
)

// parseFloatPointer converte string para *float64, retorna nil se a conversão falhar
// A entrada pode estar em formatos como "34,05", "34.05", "1.234,56"
func parseFloatPointer(s string) *float64 {
	if s == "" {
		return nil
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &value
}

// ValidateNumericString verifica se uma string contém um número válido
// Retorna (isValid, cleanedValue)
func ValidateNumericString(s string) (bool, string) {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil, s
}
