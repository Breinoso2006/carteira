# scraping-api

Stock data scraping service with an in-memory cache layer. Part of the [Carteira 2.0](../README.md) system.

Scrapes real-time stock fundamentals from multiple Brazilian financial sources (Investidor10, Auvp, Fundamentus) and caches successful results to reduce redundant requests.

Runs on **port 3001**.

---

## Prerequisites

- Go 1.21+
- Internet access (the service scrapes external websites)

> No CGO or SQLite required — this service uses a pure-Go in-memory cache.

---

## Running

```bash
go run ./cmd/main.go
```

The service will:
1. Load configuration from environment variables
2. Initialize the in-memory cache with the configured TTL
3. Start the HTTP server on `:3001`

---

## Configuration

| Variable | Default | Description |
|---|---|---|
| `DATABASE_PATH` | `./portfolio.db` | Loaded by the shared config but not used by this service. |
| `CACHE_TTL_HOURS` | `24` | How long (in hours) scraped stock data is considered valid. Must be a positive integer; invalid values fall back to `24`. |
| `CACHE_ENABLED` | `true` | Set to `false` to bypass the cache and always scrape fresh data. Any value other than `false` is treated as `true`. |

### Examples

```bash
# Development — short TTL so data refreshes quickly
CACHE_TTL_HOURS=1 go run ./cmd/main.go

# Production — 24-hour cache (default)
CACHE_TTL_HOURS=24 CACHE_ENABLED=true go run ./cmd/main.go

# Disable cache entirely (always scrape fresh)
CACHE_ENABLED=false go run ./cmd/main.go

# Custom TTL via environment file
export CACHE_TTL_HOURS=6
export CACHE_ENABLED=true
go run ./cmd/main.go
```

---

## Cache Behaviour

- **Cache hit**: If valid (non-expired, no invalid fields) data exists for a ticker, it is returned immediately without scraping.
- **Cache miss**: The service tries each configured scraper in order. The first scraper that returns complete data (no null/invalid fields) wins, and the result is stored in the cache.
- **Partial data**: If a scraper returns data with invalid/null fields, that result is **not** cached. The next scraper is tried.
- **All scrapers fail**: The error is returned to the caller. The cache is not modified.
- **Cache disabled** (`CACHE_ENABLED=false`): Every request triggers a fresh scrape regardless of cached state.
- **Expired entry**: Treated as a cache miss; a fresh scrape is triggered.

---

## API Endpoints

### GET /:ticker

Returns stock fundamentals for the given ticker symbol.

```
GET /WEGE3
```

**Response 200**

```json
{
  "symbol": "WEGE3",
  "price": 35.50,
  "pe": 28.4,
  "pbv": 8.1,
  "psr": 4.2,
  "bvps": 4.38,
  "eps": 1.25,
  "dy": 1.8,
  "source": "Investidor10",
  "invalid_fields": []
}
```

**Response 404** — ticker not found or all scrapers failed

```json
{ "error": "failed to get stock data for XXXX: ..." }
```

### Fields

| Field | Type | Description |
|---|---|---|
| `symbol` | string | Ticker symbol |
| `price` | float | Current price |
| `pe` | float | Price-to-Earnings ratio |
| `pbv` | float | Price-to-Book Value ratio |
| `psr` | float | Price-to-Sales ratio |
| `bvps` | float | Book Value per Share |
| `eps` | float | Earnings per Share |
| `dy` | float | Dividend Yield (%) |
| `source` | string | Data source that provided the result |
| `invalid_fields` | array | List of fields that could not be scraped |

The response format is identical whether data came from cache or a fresh scrape.

---

## Running Tests

```bash
go test ./...
```

---

## Project Structure

```
scraping-api/
├── cmd/
│   └── main.go                  # Entry point — wires config, cache, and routes
└── internal/
    ├── cache/
    │   └── cache_repository.go  # In-memory cache with TTL (go-cache)
    ├── config/
    │   └── config.go            # Loads CACHE_TTL_HOURS, CACHE_ENABLED, DATABASE_PATH
    ├── http/
    │   └── http_client.go       # Shared HTTP client for scrapers
    ├── models/
    │   └── stock_model.go       # StockData model
    ├── repository/
    │   └── stock_repository.go  # Orchestrates cache + scraper
    └── scraping/
        ├── scraper.go           # Scraper interface
        ├── scraper_manager.go   # Tries scrapers in order, returns first complete result
        ├── scraper_rescraper.go # Re-scrape logic for partial results
        ├── scrapers_configs.go  # Per-source field selectors
        ├── sources_config.go    # Source URLs and priorities
        └── helpers.go           # Parsing utilities
```
