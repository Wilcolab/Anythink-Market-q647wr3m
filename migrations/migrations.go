package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

// GetMigrations returns all available migrations
func GetMigrations() []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "create_questions_table",
			Up: `
				CREATE TABLE IF NOT EXISTS questions (
					id SERIAL PRIMARY KEY,
					question TEXT NOT NULL,
					created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
					updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
				);
			`,
			Down: `DROP TABLE IF EXISTS questions;`,
		},
		{
			Version: 2,
			Name:    "create_options_table",
			Up: `
				CREATE TABLE IF NOT EXISTS options (
					id SERIAL PRIMARY KEY,
					question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
					text TEXT NOT NULL,
					is_correct BOOLEAN NOT NULL DEFAULT FALSE,
					created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
				);
			`,
			Down: `DROP TABLE IF EXISTS options;`,
		},
		{
			Version: 3,
			Name:    "create_quiz_sessions_table",
			Up: `
				CREATE TABLE IF NOT EXISTS quiz_sessions (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					user_id VARCHAR(255) NOT NULL,
					score INTEGER DEFAULT 0,
					total_count INTEGER DEFAULT 0,
					started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
					completed_at TIMESTAMP WITH TIME ZONE
				);
			`,
			Down: `DROP TABLE IF EXISTS quiz_sessions;`,
		},
		{
			Version: 4,
			Name:    "create_quiz_results_table",
			Up: `
				CREATE TABLE IF NOT EXISTS quiz_results (
					id SERIAL PRIMARY KEY,
					user_id VARCHAR(255) NOT NULL,
					question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
					selected_id INTEGER NOT NULL REFERENCES options(id) ON DELETE CASCADE,
					is_correct BOOLEAN NOT NULL,
					completed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
				);
			`,
			Down: `DROP TABLE IF EXISTS quiz_results;`,
		},
		{
			Version: 5,
			Name:    "create_migrations_table",
			Up: `
				CREATE TABLE IF NOT EXISTS schema_migrations (
					version INTEGER PRIMARY KEY,
					applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
				);
			`,
			Down: `DROP TABLE IF EXISTS schema_migrations;`,
		},
	}
}

// RunMigrations executes all pending migrations
func RunMigrations(db *sql.DB) error {
	// First ensure the migrations table exists
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations := GetMigrations()

	for _, migration := range migrations {
		applied, err := isMigrationApplied(db, migration.Version)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if !applied {
			log.Printf("Running migration %d: %s", migration.Version, migration.Name)

			if err := runMigration(db, migration); err != nil {
				return fmt.Errorf("failed to run migration %d: %w", migration.Version, err)
			}

			if err := recordMigration(db, migration.Version); err != nil {
				return fmt.Errorf("failed to record migration: %w", err)
			}

			log.Printf("✅ Migration %d completed", migration.Version)
		}
	}

	return nil
}

func ensureMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	_, err := db.Exec(query)
	return err
}

func isMigrationApplied(db *sql.DB, version int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func runMigration(db *sql.DB, migration Migration) error {
	_, err := db.Exec(migration.Up)
	return err
}

func recordMigration(db *sql.DB, version int) error {
	_, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
	return err
}
