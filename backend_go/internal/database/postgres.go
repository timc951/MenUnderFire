package database

import (
	"database/sql"
	"embed"
	"fmt"

	"menunderfire/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

// RunMigrations runs all pending database migrations
// It uses golang-migrate with embedded SQL migration files
// Returns the number of migrations applied
func RunMigrations(db *sql.DB) (int, error) {
	// Create source driver from embedded filesystem
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return 0, fmt.Errorf("error creating migration source: %w", err)
	}

	// Create database driver
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return 0, fmt.Errorf("error creating database driver: %w", err)
	}

	// Create migrator
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return 0, fmt.Errorf("error creating migrator: %w", err)
	}

	// Get current version before migration
	versionBefore, _, _ := m.Version()

	// Run all pending migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return 0, fmt.Errorf("error running migrations: %w", err)
	}

	// Get current version after migration
	versionAfter, _, _ := m.Version()

	// Calculate migrations applied
	migrationsApplied := int(versionAfter) - int(versionBefore)
	if migrationsApplied < 0 {
		migrationsApplied = int(versionAfter)
	}

	return migrationsApplied, nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(db *sql.DB) (uint, bool, error) {
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("error creating database driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return 0, false, fmt.Errorf("error creating migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return 0, false, fmt.Errorf("error creating migrator: %w", err)
	}

	return m.Version()
}
