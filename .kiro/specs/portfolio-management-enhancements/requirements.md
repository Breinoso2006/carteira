szaxx# Requirements Document

## Introduction

This feature enhances the portfolio management system by implementing two key improvements:

1. **SQLite Persistence for Portfolio Data**: Replace the current in-memory stock portfolio storage in `carteira-api` with a persistent SQLite database. This ensures portfolio data survives application restarts and enables data management capabilities.

2. **Caching with TTL for Scraped Data**: Add a caching layer to `scraping-api` that stores successfully scraped stock information with a 1-day Time-To-Live (TTL). This reduces redundant scraping operations, improves response times, and reduces load on external data sources.

The system currently scrapes stock data from multiple sources (Investidor10, Auvp, Fundamentus) and stores portfolio information manually in code. These enhancements will provide persistence and performance improvements while maintaining the existing functionality.

## Glossary

- **Portfolio**: A collection of stocks with assigned weights, used for investment management
- **Stock**: A financial instrument representing ownership in a company, identified by ticker symbol
- **SQLite**: A lightweight, disk-based database that doesn't require a separate server process
- **Cache**: A temporary storage layer for frequently accessed data to improve performance
- **TTL (Time-To-Live)**: The duration for which cached data remains valid before it must be refreshed
- **Scraper**: A component that extracts stock data from external websites
- **Invalid Fields**: Data fields that could not be successfully extracted or contain null/invalid values
- **Weight**: The percentage allocation of a stock within the portfolio, calculated based on its grade

## Requirements

### Requirement 1: Portfolio Data Persistence

**User Story:** As a portfolio manager, I want my stock portfolio to be persisted in a database, so that my investment data survives application restarts and can be managed over time.

#### Acceptance Criteria

1. WHEN the application starts, THE PortfolioRepository SHALL load all portfolio entries from the SQLite database
2. THE PortfolioRepository SHALL provide methods to add, update, and remove stocks from the portfolio
3. WHEN a stock is added to the portfolio, THE PortfolioRepository SHALL store it with its ticker and fundamentalist grade
4. WHEN a stock is removed from the portfolio, THE PortfolioRepository SHALL delete it from the database
5. WHEN the portfolio weights are calculated, THE PortfolioRepository SHALL retrieve all portfolio stocks from the database
6. IF the database file does not exist, THE PortfolioRepository SHALL create it automatically with the initial schema
7. IF the database schema is outdated, THE PortfolioRepository SHALL attempt to migrate to the latest schema version
8. IF a database operation fails, THE PortfolioRepository SHALL return a descriptive error

### Requirement 2: Stock Data Caching

**User Story:** As a user of the scraping API, I want previously scraped stock data to be cached, so that repeated requests for the same data are faster and reduce load on external sources.

#### Acceptance Criteria

1. WHEN valid stock data is successfully scraped (no null/invalid fields), THE CacheRepository SHALL store it with a 1-day TTL
2. WHEN a stock data request is made, THE CacheRepository SHALL check if valid cached data exists and return it if not expired
3. WHEN cached data has expired (older than 1 day), THE CacheRepository SHALL return a cache miss and trigger a fresh scrape
4. WHEN invalid stock data is scraped (contains null/invalid fields), THE CacheRepository SHALL NOT cache it
5. IF a cache operation fails, THE CacheRepository SHALL log the error and continue with fresh scraping
6. WHERE multiple scrapers are available, THE CacheRepository SHALL cache data from the first successful scraper that provides complete data
7. WHEN cache is disabled, THE CacheRepository SHALL bypass cache operations and always fetch fresh data

### Requirement 3: Data Model Consistency

**User Story:** As a developer, I want the data models to be consistent between the API and database, so that data integrity is maintained throughout the application.

#### Acceptance Criteria

1. THE StockData model SHALL include all fields from the database schema: Symbol, Price, PE, PBV, PSR, BVps, EPS, DY, Source, CreatedAt, and UpdatedAt
2. WHEN a StockData object is persisted, THE DatabaseRepository SHALL store all fields including timestamps
3. WHEN a StockData object is retrieved, THE DatabaseRepository SHALL populate all fields from the database
4. THE PortfolioEntry model SHALL include fields: ID, Ticker, FundamentalistGrade, CreatedAt, and UpdatedAt
5. IF a field is null in the database, THE DatabaseRepository SHALL represent it as a null pointer in the Go struct

### Requirement 4: Cache Invalidation

**User Story:** As a system administrator, I want the cache to be properly invalidated when data changes, so that stale data is not served to users.

#### Acceptance Criteria

1. WHEN a stock data request fails completely (all scrapers fail), THE CacheRepository SHALL NOT invalidate the cache entry
2. WHEN a cache entry is accessed after expiration, THE CacheRepository SHALL invalidate it after returning the miss
3. WHERE a manual cache refresh is requested, THE CacheRepository SHALL invalidate the existing entry and fetch fresh data
4. IF the cache storage fails to write, THE CacheRepository SHALL invalidate the stale entry and continue operation

### Requirement 5: Error Handling and Logging

**User Story:** As a developer, I want errors to be properly logged and handled, so that issues can be diagnosed and resolved quickly.

#### Acceptance Criteria

1. WHEN a database connection fails, THE DatabaseRepository SHALL log the error and return a descriptive error to the caller
2. WHEN a cache operation fails, THE CacheRepository SHALL log the error and continue with fresh scraping
3. IF a migration fails, THE DatabaseRepository SHALL log the error and return an error without modifying the database
4. WHEN invalid data is encountered, THE System SHALL log the specific fields that are invalid
5. ALL error messages SHALL include context about the operation that failed and the relevant identifiers

### Requirement 6: API Compatibility

**User Story:** As a client of the scraping API, I want the API responses to remain consistent, so that my application doesn't break when the underlying implementation changes.

#### Acceptance Criteria

1. WHEN a stock data request is made, THE API SHALL return the same response format regardless of whether data came from cache or fresh scraping
2. THE API response SHALL include all stock fields: Symbol, Price, PE, PBV, PSR, BVps, EPS, DY, Source, and invalid fields indicator
3. IF cached data is returned, THE API response SHALL NOT indicate that it came from cache
4. WHEN portfolio data is requested, THE API SHALL return the same format as before the persistence changes

### Requirement 7: Performance Requirements

**User Story:** As a user, I want the system to respond quickly, so that I can get stock information without unnecessary delays.

#### Acceptance Criteria

1. WHEN valid cached data is available, THE CacheRepository SHALL return it within 50ms
2. WHEN fresh scraping is required, THE System SHALL complete within 5 seconds for a single stock
3. WHEN portfolio data is loaded, THE DatabaseRepository SHALL load all entries within 1 second for up to 100 stocks
4. IF cache is disabled, THE System SHALL still meet the 5-second response time for stock data

### Requirement 8: Data Migration

**User Story:** As a developer, I want existing portfolio data to be migrated to the new database, so that no data is lost during the transition.

#### Acceptance Criteria

1. WHEN the application starts with an existing in-memory portfolio, THE MigrationTool SHALL create the database and populate it with existing portfolio data
2. THE MigrationTool SHALL preserve all stock tickers and fundamentalist grades during migration
3. IF migration encounters invalid data, THE MigrationTool SHALL log the issue and skip that entry
4. AFTER migration completes, THE System SHALL verify data integrity by comparing counts and sample entries

### Requirement 9: Configuration

**User Story:** As a system administrator, I want to configure database and cache settings, so that I can adapt the system to different environments.

#### Acceptance Criteria

1. WHERE environment variables are available, THE System SHALL read database path from `DATABASE_PATH` with default `./portfolio.db`
2. WHERE environment variables are available, THE System SHALL read cache TTL in hours from `CACHE_TTL_HOURS` with default `24`
3. WHERE environment variables are available, THE System SHALL read cache enabled flag from `CACHE_ENABLED` with default `true`
4. IF configuration values are invalid, THE System SHALL log warnings and use defaults

### Requirement 10: Testing and Validation

**User Story:** As a developer, I want comprehensive tests for the persistence and caching features, so that I can ensure they work correctly and don't introduce regressions.

#### Acceptance Criteria

1. FOR ALL database operations, THE TestSuite SHALL include unit tests covering success and failure cases
2. FOR all cache operations, THE TestSuite SHALL include property-based tests for cache hit/miss behavior
3. FOR the migration tool, THE TestSuite SHALL include tests that verify data integrity after migration
4. FOR the cache TTL functionality, THE TestSuite SHALL include tests that verify expiration behavior
5. FOR error handling, THE TestSuite SHALL include tests that verify proper error propagation and logging
