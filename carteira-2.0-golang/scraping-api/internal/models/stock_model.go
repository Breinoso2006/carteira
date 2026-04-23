package models

import "time"

type StockData struct {
	Symbol    string    `json:"Symbol"`
	Price     *float64  `json:"Price"`
	PE        *float64  `json:"PE"`
	PBV       *float64  `json:"PBV"`
	PSR       *float64  `json:"PSR"`
	BVps      *float64  `json:"BVps"`
	EPS       *float64  `json:"EPS"`
	DY        *float64  `json:"DY"`
	Source    string    `json:"Source"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`

	// Interno: quais campos estão vazios/inválidos
	invalidFields map[string]bool
}

// IsFieldInvalid returns true if the given field has been marked as invalid.
func (s *StockData) IsFieldInvalid(field string) bool {
	if s.invalidFields == nil {
		return false
	}
	return s.invalidFields[field]
}

// SetFieldInvalid marks the given field as invalid (true) or valid (false).
func (s *StockData) SetFieldInvalid(field string, invalid bool) {
	if s.invalidFields == nil {
		s.invalidFields = make(map[string]bool)
	}
	s.invalidFields[field] = invalid
}

// MarkFieldInvalid marks the given field as invalid.
func (s *StockData) MarkFieldInvalid(field string) {
	if s.invalidFields == nil {
		s.invalidFields = make(map[string]bool)
	}
	s.invalidFields[field] = true
}

// GetInvalidFields returns a slice of field names that are marked as invalid.
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

// GetInvalidFieldsMap returns the raw invalidFields map (a copy).
func (s *StockData) GetInvalidFieldsMap() map[string]bool {
	if s.invalidFields == nil {
		return map[string]bool{}
	}
	copy := make(map[string]bool, len(s.invalidFields))
	for k, v := range s.invalidFields {
		copy[k] = v
	}
	return copy
}

// ClearInvalidFields resets the invalid fields tracking map.
func (s *StockData) ClearInvalidFields() {
	s.invalidFields = make(map[string]bool)
}
