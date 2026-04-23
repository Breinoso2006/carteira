package http

import (
	"net/http"

	"github.com/breinoso2006/carteira-api/internal/repository"
	"github.com/gofiber/fiber/v2"
)

// PortfolioHandler handles portfolio HTTP requests.
type PortfolioHandler struct {
	repo *repository.PortfolioRepository
}

// NewPortfolioHandler creates a new PortfolioHandler.
func NewPortfolioHandler(repo *repository.PortfolioRepository) *PortfolioHandler {
	return &PortfolioHandler{repo: repo}
}

// GetAll returns all portfolio entries with calculated weights.
func (h *PortfolioHandler) GetAll(c *fiber.Ctx) error {
	entries, err := h.repo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.repo.CalculateWeights(entries); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(entries)
}

// Add adds a new stock to the portfolio.
func (h *PortfolioHandler) Add(c *fiber.Ctx) error {
	var req struct {
		Ticker              string  `json:"ticker"`
		FundamentalistGrade float64 `json:"fundamentalist_grade"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Ticker == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ticker is required"})
	}

	if req.FundamentalistGrade <= 0 || req.FundamentalistGrade > 100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "fundamentalist_grade must be between 0 and 100"})
	}

	if err := h.repo.Add(req.Ticker, req.FundamentalistGrade); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "stock added successfully"})
}

// Update updates an existing stock in the portfolio.
func (h *PortfolioHandler) Update(c *fiber.Ctx) error {
	var req struct {
		Ticker              string  `json:"ticker"`
		FundamentalistGrade float64 `json:"fundamentalist_grade"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Ticker == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ticker is required"})
	}

	if req.FundamentalistGrade <= 0 || req.FundamentalistGrade > 100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "fundamentalist_grade must be between 0 and 100"})
	}

	if err := h.repo.Update(req.Ticker, req.FundamentalistGrade); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "stock updated successfully"})
}

// Remove removes a stock from the portfolio.
func (h *PortfolioHandler) Remove(c *fiber.Ctx) error {
	ticker := c.Params("ticker")
	if ticker == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ticker is required"})
	}

	if err := h.repo.Remove(ticker); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "stock removed successfully"})
}
