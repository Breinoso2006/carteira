package main

import (
	"fmt"
	"log"

	"github.com/breinoso2006/carteira-api/internal/config"
	"github.com/breinoso2006/carteira-api/internal/database"
	"github.com/breinoso2006/carteira-api/internal/http"
	"github.com/breinoso2006/carteira-api/internal/migration"
	"github.com/breinoso2006/carteira-api/internal/models"
	"github.com/breinoso2006/carteira-api/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration from environment variables.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Open (or create) the SQLite database and apply migrations.
	db, err := database.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Seed the database with the initial in-memory portfolio on first run.
	initialPortfolio := []*models.StockInPortfolio{
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

	migrationTool := migration.NewMigrationTool(db)

	if err := migrationTool.MigratePortfolio(initialPortfolio); err != nil {
		log.Fatalf("Failed to migrate portfolio: %v", err)
	}

	if err := migrationTool.VerifyMigration(initialPortfolio); err != nil {
		log.Fatalf("Migration verification failed: %v", err)
	}

	// Wire up the repository and HTTP handler.
	portfolioRepo := repository.NewPortfolioRepository(db)
	handler := http.NewPortfolioHandler(portfolioRepo)

	app := fiber.New()
	app.Get("/portfolio", handler.GetAll)
	app.Post("/portfolio", handler.Add)
	app.Put("/portfolio", handler.Update)
	app.Delete("/portfolio/:ticker", handler.Remove)

	fmt.Println("Server starting on :3000")
	if err := app.Listen(":3002"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
