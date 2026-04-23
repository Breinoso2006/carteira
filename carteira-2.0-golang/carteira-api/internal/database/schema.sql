-- Database Schema for Portfolio Management
-- Version: 1

-- Table: portfolio_entries
-- Stores the user's stock portfolio with fundamentalist grades
CREATE TABLE IF NOT EXISTS portfolio_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticker TEXT NOT NULL UNIQUE,
    fundamentalist_grade REAL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: stock_cache
-- Stores cached stock data with TTL support
CREATE TABLE IF NOT EXISTS stock_cache (
    symbol TEXT PRIMARY KEY,
    price REAL,
    pe REAL,
    pbv REAL,
    psr REAL,
    bvps REAL,
    eps REAL,
    dy REAL,
    source TEXT,
    created_at TIMESTAMP,
    expires_at TIMESTAMP,
    invalid_fields TEXT
);

-- Table: schema_version
-- Tracks the current schema version for migrations
CREATE TABLE IF NOT EXISTS schema_version (
    version INTEGER PRIMARY KEY,
    migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
