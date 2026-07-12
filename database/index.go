package database

import (
	"adeeb_huma/config"
	"fmt"

	"time"

	"gorm.io/driver/postgres" // Or your specific driver
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Conn *gorm.DB

// NewDatabase creates a new GORM database instance with optimized settings
func NewDatabase() (*gorm.DB, error) {
	// Build the Data Source Name (DSN)
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		config.DB_HOST, config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_PORT,
	)

	var logger_mode = logger.Warn
	if config.APP_ENV == "dev" {
		logger_mode = logger.Info
	}

	// Configure GORM options
	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger_mode),
		SkipDefaultTransaction: true, // Improves performance
		PrepareStmt:            true, // Caches prepared statements
	}

	// Open the connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Access the underlying *sql.DB to configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)                  // Max idle connections
	sqlDB.SetMaxOpenConns(100)                 // Max open connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Max lifetime of a connection
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Max idle time

	Conn = db

	return Conn, nil
}
