# carteira-api

REST API for managing a stock portfolio with SQLite persistence. Part of the [Carteira 2.0](../README.md) system.

Runs on **port 3002**.

---

## Prerequisites

- Go 1.21+
- CGO enabled (required by `go-sqlite3`)
- GCC or compatible C compiler

---

## Running

```bash
go run ./cmd/main.go
```

The service will:
1. Load configuration from environment variables
2. Open (or create) the SQLite database at `DATABASE_PATH`
3. Apply schema migrations automatically
4. Seed the initial portfolio on first run
5. Start the HTTP server on `:3000`

---

## Configuration

| Variable | Default | Description |
|---|---|---|
| `DATABASE_PATH` | `./portfolio.db` | Path to the SQLite database file. Created automatically if it does not exist. |
| `CACHE_TTL_HOURS` | `24` | Unused by this service directly, but loaded by the shared config. |
| `CACHE_ENABLED` | `true` | Unused by this service directly, but loaded by the shared config. |

### Examples

```bash
# Development
DATABASE_PATH=./dev-portfolio.db go run ./cmd/main.go

# Production
DATABASE_PATH=/var/data/carteira/portfolio.db go run ./cmd/main.go

# Custom path with directory (created automatically)
DATABASE_PATH=/tmp/myapp/portfolio.db go run ./cmd/main.go
```

---

## Database

### SQLite setup

No manual setup is required. On startup the service:

- Creates the SQLite file at `DATABASE_PATH` if it does not exist (including any missing parent directories)
- Runs `CREATE TABLE IF NOT EXISTS` for all tables — safe to run on every restart
- Checks the `schema_version` table and applies any pending migrations
- Seeds the initial portfolio of 18 Brazilian stocks if the `portfolio_entries` table is empty

### Schema

```sql
-- Portfolio entries
CREATE TABLE IF NOT EXISTS portfolio_entries (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    ticker               TEXT    NOT NULL UNIQUE,
    fundamentalist_grade REAL    NOT NULL,
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Schema version tracking
CREATE TABLE IF NOT EXISTS schema_version (
    version     INTEGER PRIMARY KEY,
    migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Migration

Migrations run automatically on startup. The migration tool:

- Inserts each entry from the in-memory seed list into `portfolio_entries`
- Skips entries that already exist (idempotent)
- Logs skipped and migrated entries
- Verifies the final row count matches the expected seed count

To reset and reseed:

```bash
rm ./portfolio.db
go run ./cmd/main.go
```

---

## API Endpoints

### GET /portfolio

Returns all portfolio entries with calculated weights.

**Response 200**

```json
[
  {
    "id": 1,
    "ticker": "WEGE3",
    "fundamentalist_grade": 98.75,
    "weight": 0.065,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

Weight is calculated as each stock's grade divided by the sum of all grades.

---

### POST /portfolio

Adds a new stock to the portfolio.

**Request body**

```json
{
  "ticker": "PETR4",
  "fundamentalist_grade": 65.0
}
```

- `ticker` — required, non-empty string
- `fundamentalist_grade` — required, must be between 0 (exclusive) and 100 (inclusive)

**Response 201**

```json
{ "message": "stock added successfully" }
```

---

### PUT /portfolio

Updates the fundamentalist grade of an existing stock.

**Request body**

```json
{
  "ticker": "PETR4",
  "fundamentalist_grade": 70.0
}
```

**Response 200**

```json
{ "message": "stock updated successfully" }
```

---

### DELETE /portfolio/:ticker

Removes a stock from the portfolio.

```
DELETE /portfolio/PETR4
```

**Response 200**

```json
{ "message": "stock removed successfully" }
```

---

## Running Tests

```bash
go test ./...
```

Tests use an in-memory or temporary SQLite database and do not affect your `portfolio.db`.

---

## Project Structure

```
carteira-api/
├── cmd/
│   └── main.go              # Entry point — wires all components
└── internal/
    ├── config/
    │   └── config.go        # Loads DATABASE_PATH, CACHE_TTL_HOURS, CACHE_ENABLED
    ├── database/
    │   ├── database.go      # Opens SQLite, applies schema, runs migrations
    │   └── schema.sql       # Reference DDL
    ├── http/
    │   ├── http_client.go   # HTTP client helpers
    │   └── portfolio_handler.go  # Fiber route handlers
    ├── migration/
    │   └── migration_tool.go     # Seeds initial portfolio data
    ├── models/
    │   ├── portfolio_model.go    # PortfolioEntry model
    │   └── stock_model.go        # Stock / StockInPortfolio models
    └── repository/
        └── portfolio_repository.go  # CRUD operations against SQLite
```
