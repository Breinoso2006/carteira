package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/breinoso2006/carteira-api/internal/database"
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

// PortfolioRepository handles portfolio data persistence.
type PortfolioRepository struct {
	db *database.Database
}

// NewPortfolioRepository creates a new PortfolioRepository.
func NewPortfolioRepository(db *database.Database) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

// GetAll retrieves all portfolio entries ordered by ticker.
func (r *PortfolioRepository) GetAll() ([]*PortfolioEntry, error) {
	rows, err := r.db.GetDB().Query(
		"SELECT id, ticker, fundamentalist_grade, created_at, updated_at FROM portfolio_entries ORDER BY ticker",
	)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var entries []*PortfolioEntry
	for rows.Next() {
		e := &PortfolioEntry{}
		if err := rows.Scan(&e.ID, &e.Ticker, &e.FundamentalistGrade, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan portfolio row: %w", err)
		}
		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating portfolio rows: %w", err)
	}

	return entries, nil
}

// Add inserts a new stock into the portfolio.
func (r *PortfolioRepository) Add(ticker string, fundamentalistGrade float64) error {
	now := time.Now().UTC()
	_, err := r.db.GetDB().Exec(
		"INSERT INTO portfolio_entries (ticker, fundamentalist_grade, created_at, updated_at) VALUES (?, ?, ?, ?)",
		ticker, fundamentalistGrade, now, now,
	)
	if err != nil {
		if isConstraintViolation(err) {
			return fmt.Errorf("stock %s already exists in portfolio: %w", ticker, err)
		}
		return fmt.Errorf("database insert failed for ticker %s: %w", ticker, err)
	}
	return nil
}

// Update modifies the fundamentalist grade of an existing portfolio entry.
func (r *PortfolioRepository) Update(ticker string, fundamentalistGrade float64) error {
	result, err := r.db.GetDB().Exec(
		"UPDATE portfolio_entries SET fundamentalist_grade = ?, updated_at = ? WHERE ticker = ?",
		fundamentalistGrade, time.Now().UTC(), ticker,
	)
	if err != nil {
		return fmt.Errorf("database update failed for ticker %s: %w", ticker, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("stock %s not found in portfolio", ticker)
	}
	return nil
}

// Remove deletes a portfolio entry by ticker.
func (r *PortfolioRepository) Remove(ticker string) error {
	result, err := r.db.GetDB().Exec(
		"DELETE FROM portfolio_entries WHERE ticker = ?", ticker,
	)
	if err != nil {
		return fmt.Errorf("database delete failed for ticker %s: %w", ticker, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("stock %s not found in portfolio", ticker)
	}
	return nil
}

// CalculateWeights sets the Weight field on each entry proportional to its
// FundamentalistGrade. Entries with a zero total grade are left at 0.
func (r *PortfolioRepository) CalculateWeights(entries []*PortfolioEntry) error {
	var total float64
	for _, e := range entries {
		total += e.FundamentalistGrade
	}

	if total == 0 {
		return nil
	}

	for _, e := range entries {
		e.Weight = e.FundamentalistGrade / total * 100
	}
	return nil
}

// isConstraintViolation reports whether err is a SQLite unique/constraint error.
func isConstraintViolation(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "constraint failed")
}
