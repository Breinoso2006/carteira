# Implementation Plan: Portfolio Management Enhancements

## Overview

This feature enhances the portfolio management system by implementing SQLite persistence for portfolio data in `carteira-api` and caching with TTL for scraped stock data in `scraping-api`. The implementation follows a database-first approach, creating repositories with proper error handling and testing.

## Tasks

- [x] 1. Set up project structure and dependencies
  - [x] 1.1 Add SQLite dependency to carteira-api
    - Update go.mod to include github.com/mattn/go-sqlite3
    - Run go mod tidy to install dependency
    - _Requirements: 1.1, 1.6_
  
  - [x] 1.2 Add cache storage dependency to scraping-api
    - Update go.mod to include github.com/patrickmn/go-cache
    - Run go mod tidy to install dependency
    - _Requirements: 2.1, 2.2_

- [x] 2. Create database schema and migration system
  - [x] 2.1 Create database schema file
    - Define portfolio_entries table schema
    - Define stock_cache table schema
    - Define schema_version table schema
    - _Requirements: 3.1, 3.2_
  
  - [x] 2.2 Implement database connection and initialization
    - Create NewDatabase function in database package
    - Implement automatic database file creation if not exists
    - Implement schema version checking
    - _Requirements: 1.6, 1.7, 9.1_
  
  - [x] 2.3 Implement migration system
    - Create migration scripts for schema version upgrades
    - Implement version tracking in schema_version table
    - Add error handling for migration failures
    - _Requirements: 1.7, 8.1, 8.4_

- [x] 3. Implement Portfolio Repository
  - [x] 3.1 Create PortfolioEntry model
    - Define struct with ID, Ticker, FundamentalistGrade, CreatedAt, UpdatedAt
    - Implement validation methods
    - _Requirements: 3.4, 3.5_
  
  - [x] 3.2 Implement GetAll method
    - Query all portfolio entries from database
    - Calculate weights for returned entries
    - Handle database errors gracefully
    - _Requirements: 1.1, 1.5_
  
  - [x] 3.3 Implement Add method
    - Insert new portfolio entry with ticker and grade
    - Handle constraint violations
    - Return descriptive errors
    - _Requirements: 1.3, 1.8_
  
  - [x] 3.4 Implement Update method
    - Update existing portfolio entry by ticker
    - Update UpdatedAt timestamp
    - Handle missing entries
    - _Requirements: 1.3, 1.8_
  
  - [x] 3.5 Implement Remove method
    - Delete portfolio entry by ticker
    - Handle missing entries gracefully
    - _Requirements: 1.4, 1.8_
  
  - [ ]* 3.6 Write property test for Portfolio Entry Persistence Round Trip
    - **Property 1: Portfolio Entry Persistence Round Trip**
    - **Validates: Requirements 1.3, 1.4**
  
  - [ ]* 3.7 Write property test for Portfolio Entry Removal
    - **Property 2: Portfolio Entry Removal**
    - **Validates: Requirements 1.4**
  
  - [ ]* 3.8 Write property test for Portfolio Weight Calculation Consistency
    - **Property 3: Portfolio Weight Calculation Consistency**
    - **Validates: Requirements 1.5**

- [x] 4. Implement Cache Repository
  - [x] 4.1 Create StockData model
    - Define struct with all required fields
    - Implement invalidFields map for tracking invalid fields
    - _Requirements: 3.1, 3.5_
  
  - [x] 4.2 Implement GetStockData method
    - Check cache for valid non-expired entry
    - Trigger fresh scrape on cache miss
    - Return same format regardless of source
    - _Requirements: 2.2, 6.1, 6.3_
  
  - [x] 4.3 Implement StoreStockData method
    - Store valid StockData with 1-day TTL
    - Reject data with invalid fields
    - Handle storage errors gracefully
    - _Requirements: 2.1, 2.4, 2.6, 9.3_
  
  - [x] 4.4 Implement HasValidCache method
    - Check if cache entry exists and is not expired
    - Return true/false appropriately
    - _Requirements: 2.2, 2.3_
  
  - [x] 4.5 Implement Invalidate method
    - Remove cache entry for given symbol
    - Handle missing entries gracefully
    - _Requirements: 4.2, 4.3_
  
  - [x] 4.6 Implement Refresh method
    - Invalidate existing entry
    - Fetch fresh data from scrapers
    - Store new data in cache
    - _Requirements: 4.3_
  
  - [ ]* 4.7 Write property test for Cache Storage for Valid Data
    - **Property 6: Cache Storage for Valid Data**
    - **Validates: Requirements 2.1**
  
  - [ ]* 4.8 Write property test for Cache Rejection of Invalid Data
    - **Property 7: Cache Rejection of Invalid Data**
    - **Validates: Requirements 2.4**
  
  - [ ]* 4.9 Write property test for Cache Hit for Valid Non-Expired Data
    - **Property 8: Cache Hit for Valid Non-Expired Data**
    - **Validates: Requirements 2.2**
  
  - [ ]* 4.10 Write property test for Cache Miss for Expired Data
    - **Property 9: Cache Miss for Expired Data**
    - **Validates: Requirements 2.3**
  
  - [ ]* 4.11 Write property test for Cache Preservation on Complete Failure
    - **Property 10: Cache Preservation on Complete Failure**
    - **Validates: Requirements 4.1**
  
  - [ ]* 4.12 Write property test for Cache Invalidation on Expiration Access
    - **Property 11: Cache Invalidation on Expiration Access**
    - **Validates: Requirements 4.2**
  
  - [ ]* 4.13 Write property test for Manual Cache Refresh
    - **Property 12: Manual Cache Refresh**
    - **Validates: Requirements 4.3**

- [x] 5. Implement Stock Repository (scraping-api)
  - [x] 5.1 Update StockRepository to use CacheRepository
    - Modify GetStockInformation to check cache first
    - Store successful results in cache
    - Handle cache failures gracefully
    - _Requirements: 2.1, 2.2, 2.5, 2.6_
  
  - [x] 5.2 Implement cache TTL configuration
    - Read CACHE_TTL_HOURS from environment
    - Default to 24 hours if not set
    - Apply TTL to cache storage
    - _Requirements: 9.2_

- [x] 6. Update carteira-api to use database
  - [x] 6.1 Update main.go to initialize database
    - Load configuration from environment
    - Initialize database connection
    - Run migrations if needed
    - _Requirements: 9.1, 9.4_
  
  - [x] 6.2 Update portfolio handler to use PortfolioRepository
    - Replace in-memory portfolio with database calls
    - Implement GetAll endpoint
    - Implement Add/Update/Delete endpoints
    - _Requirements: 1.1, 1.2, 1.5, 6.4_
  
  - [x] 6.3 Implement migration tool
    - Create migration tool to populate database from in-memory portfolio
    - Preserve all tickers and grades during migration
    - Verify data integrity after migration
    - _Requirements: 8.1, 8.2, 8.3, 8.4_

- [x] 7. Implement configuration system
  - [x] 7.1 Create Config struct
    - Define DatabasePath, CacheTTlHours, CacheEnabled fields
    - Add validation methods
    - _Requirements: 9.1, 9.2, 9.3, 9.4_
  
  - [x] 7.2 Implement environment variable loading
    - Load DATABASE_PATH, CACHE_TTL_HOURS, CACHE_ENABLED
    - Apply defaults for missing values
    - Log warnings for invalid values
    - _Requirements: 9.1, 9.2, 9.3, 9.4_

- [x] 8. Write comprehensive tests
  - [x] 8.1 Write unit tests for database operations
    - Test connection success and failure
    - Test migration success and failure
    - Test constraint violations
    - _Requirements: 10.1_
  
  - [x] 8.2 Write unit tests for cache operations
    - Test cache write success and failure
    - Test cache read success and failure
    - Test TTL expiration edge cases
    - _Requirements: 10.2_
  
  - [x] 8.3 Write integration tests for portfolio operations
    - Test end-to-end portfolio CRUD operations
    - Test weight calculation with database
    - _Requirements: 10.4_
  
  - [x] 8.4 Write integration tests for caching
    - Test cache hit/miss behavior
    - Test cache invalidation
    - Test API response consistency
    - _Requirements: 10.3, 10.5_
  
  - [ ]* 8.5 Write property test for Stock Data Persistence Round Trip
    - **Property 4: Stock Data Persistence Round Trip**
    - **Validates: Requirements 3.2, 3.3**
  
  - [ ]* 8.6 Write property test for Null Field Preservation
    - **Property 5: Null Field Preservation**
    - **Validates: Requirements 3.5**
  
  - [ ]* 8.7 Write property test for API Response Format Consistency
    - **Property 13: API Response Format Consistency**
    - **Validates: Requirements 6.1, 6.3**
  
  - [ ]* 8.8 Write property test for API Response Field Completeness
    - **Property 14: API Response Field Completeness**
    - **Validates: Requirements 6.2**
  
  - [ ]* 8.9 Write property test for Migration Data Preservation
    - **Property 15: Migration Data Preservation**
    - **Validates: Requirements 8.2**

- [x] 9. Update documentation
  - [x] 9.1 Update README with database setup
    - Document SQLite database requirements
    - Add migration instructions
    - _Requirements: 10.6_
  
  - [x] 9.2 Document configuration options
    - Document DATABASE_PATH, CACHE_TTL_HOURS, CACHE_ENABLED
    - Add examples for different environments
    - _Requirements: 9.1, 9.2, 9.3_

- [x] 10. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 11. Final checkpoint - Verify performance requirements
  - Verify cache returns data within 50ms for valid entries
  - Verify database queries complete within 1 second for 100 stocks
  - Verify API response times meet requirements
  - _Requirements: 7.1, 7.2, 7.3, 7.4_

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases
- The implementation follows Go best practices and the repository pattern
- All error handling includes descriptive error messages with context