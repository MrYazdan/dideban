// Package storage provides a GORM-based database layer for the Dideban monitoring system.
//
// It supports both SQLite (for development) and PostgreSQL (for production) with
// automatic schema migration, connection pooling, and clean resource management.
// The design eliminates custom ORM layers in favor of GORM's fluent, type-safe API.
package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"dideban/internal/config"
)

// Storage wraps the GORM database instance and provides access to it.
type Storage struct {
	db *gorm.DB
}

// New initializes a new Storage instance using GORM based on the provided configuration.
//
// Supported drivers:
//   - "sqlite": for development and single-node deployments
//   - "postgres": for production, high-availability setups
//
// Connection pooling and timeouts are configured according to config.StorageConfig.
// All models are auto-migrated on startup.
func New(cfg config.StorageConfig) (*Storage, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "sqlite":
		// Enable WAL, foreign keys, and performance pragmas via DSN
		dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=on", cfg.DSN)
		dialector = sqlite.Open(dsn)
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	// Configure GORM logger (silent by default; can be overridden via config if needed)
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sql.DB from GORM: %w", err)
	}

	// Apply connection pool settings from config
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Auto-migrate all models
	models := []interface{}{
		&Check{},
		&CheckHistory{},
		&Agent{},
		&AgentHistory{},
		&Alert{},
		&AlertHistory{},
		&Admin{},
	}
	if err := db.AutoMigrate(models...); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	return &Storage{db: db}, nil
}

// DB returns the underlying GORM database instance.
// This is the primary interface for all database operations (queries, transactions, etc.).
func (s *Storage) DB() *gorm.DB {
	return s.db
}

// Close closes the underlying database connection.
func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve sql.DB for closing: %w", err)
	}
	return sqlDB.Close()
}
