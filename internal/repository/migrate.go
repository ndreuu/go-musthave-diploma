package repository

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

// RunMigrations executes all pending database migrations on the provided PostgreSQL database.
//
// This function uses the golang-migrate library to apply schema changes from the
// migrations/ directory. It creates the migration driver, initializes the migrator,
// and runs all up migrations until the schema is current.
//
// If no migrations need to be applied (database is already up-to-date), the function
// logs an informational message but does not return an error.
//
// Parameters:
//   - db: Active PostgreSQL database connection
//   - databaseURI: Connection string for the database (used by migrate library)
//   - logger: Structured logger for migration status messages
//
// Returns:
//   - error: Migration error if any occurred, nil if successful or if no changes needed
//
// Example usage:
//
//	db, err := sql.Open("postgres", dsn)
//	if err != nil {
//	    return err
//	}
//	if err := repository.RunMigrations(db, dsn, logger); err != nil {
//	    return fmt.Errorf("failed to run migrations: %w", err)
//	}
func RunMigrations(db *sql.DB, databaseURI string, logger *zap.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		databaseURI,
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("migrations: no changes")
	} else {
		logger.Info("migrations applied")
	}

	return nil
}
