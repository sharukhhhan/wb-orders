package app

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(pgURL, migrationsPath string) error {
	if migrationsPath == "" {
		return nil
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		pgURL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
