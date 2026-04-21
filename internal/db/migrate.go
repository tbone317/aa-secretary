package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Migrate(ctx context.Context, db *sql.DB, migrationsDir string) error {
	if err := ensureMigrationsTable(ctx, db); err != nil {
		return err
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		migrationName := entry.Name()
		if filepath.Ext(migrationName) != ".sql" {
			continue
		}

		alreadyApplied, err := isMigrationApplied(ctx, db, migrationName)
		if err != nil {
			return err
		}
		if alreadyApplied {
			continue
		}

		migrationPath := filepath.Join(migrationsDir, migrationName)
		sqlBytes, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationName, err)
		}

		if err := applyMigration(ctx, db, migrationName, string(sqlBytes)); err != nil {
			return err
		}
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, db *sql.DB) error {
	const q = "CREATE TABLE IF NOT EXISTS schema_migrations(id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL UNIQUE,applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)"

	if _, err := db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}
	return nil
}

func isMigrationApplied(ctx context.Context, db *sql.DB, name string) (bool, error) {
	const q = "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name = ?)"

	var exists bool
	if err := db.QueryRowContext(ctx, q, name).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check if migration %s is applied: %w", name, err)
	}
	return exists, nil
}

func applyMigration(ctx context.Context, db *sql.DB, name, sql string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for migration %s: %w", name, err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.ExecContext(ctx, strings.TrimSpace(sql)); err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", name, err)
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations(name) VALUES(?)", name); err != nil {
		return fmt.Errorf("failed to record applied migration %s: %w", name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration %s: %w", name, err)
	}
	committed = true
	return nil
}
