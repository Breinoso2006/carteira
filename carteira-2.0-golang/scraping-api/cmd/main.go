package main

import (
	"log"

	"github.com/breinoso2006/scraping-api/internal/config"
	"github.com/breinoso2006/scraping-api/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app := fiber.New()

	stockRepo := repository.NewStockRepository(cfg.CacheTTlHours, cfg.CacheEnabled)

	app.Get("/:ticker", func(c *fiber.Ctx) error {
		ticker := c.Params("ticker")
		data, err := stockRepo.GetStockData(ticker)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	app.Listen(":3001")
}
