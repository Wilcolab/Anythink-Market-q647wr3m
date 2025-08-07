package database

import (
	"database/sql"
	"fmt"
	"go-quiz-api/config"
	"go-quiz-api/migrations"
	"log"

	_ "github.com/lib/pq"
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// Connect establishes a database connection and runs migrations
func Connect() (*DB, error) {
	cfg := config.LoadDatabaseConfig()

	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Connected to PostgreSQL database")

	// Run migrations
	log.Println("🔄 Running database migrations...")
	if err = migrations.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
