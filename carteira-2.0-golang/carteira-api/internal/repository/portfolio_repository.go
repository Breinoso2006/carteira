package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/breinoso2006/carteira-api/internal/database"
	"github.com/breinoso2006/carteira-api/internal/scoring"
)

// PortfolioRepository handles portfolio data persistence.
type PortfolioRepository struct {
	db *database.Database
}

// NewPortfolioRepository creates a new PortfolioRepository.
func NewPortfolioRepository(db *database.Database) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

// GetAll retrieves all portfolio entries ordered by ticker.
func (r *PortfolioRepository) GetAll() ([]*scoring.PortfolioEntry, error) {
	rows, err := r.db.GetDB().Query(
		"SELECT id, ticker, fundamentalist_grade, created_at, updated_at FROM portfolio_entries ORDER BY ticker",
	)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var entries []*scoring.PortfolioEntry
	for rows.Next() {
		e := &scoring.PortfolioEntry{}
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

// CalculateWeights delega o cálculo de pesos para scoring.CalculateWeights,
// mantendo a lógica centralizada no pacote scoring.
func (r *PortfolioRepository) CalculateWeights(entries []*scoring.PortfolioEntry) error {
	scoring.CalculateWeights(entries)
	return nil
}

// isConstraintViolation reports whether err is a SQLite unique/constraint error.
func isConstraintViolation(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "constraint failed")
}
