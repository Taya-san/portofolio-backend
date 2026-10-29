package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
)

//go:embed sql/*.sql
var sqlFS embed.FS

func DB(ctx context.Context, db *sql.DB) error {
	sub, err := fs.Sub(sqlFS, "sql")
	if err != nil {
		return fmt.Errorf("Migrate - fs.Sub: %w", err)
	}

	if _, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXIST schema_migrations (
			name				TEXT PRIMARY KEY
			applied_at	TEXT NOT NULL
		)`); err != nil {
		return fmt.Errorf("Migrate - Create schema_migrations: %w", err)
	}

	entries, err := fs.ReadDir(sub, ".")
	if err != nil {
		return fmt.Errorf("Migrate - Read sql Dir: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		var applied int
		if err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM schema_migrations WHERE name = ?`,
			name,
		).Scan(&applied); err != nil {
			return fmt.Errorf("Migrate - Check %s: %w", name, err)
		}

		if applied > 0 {
			continue
		}

		body, err := fs.ReadFile(sub, name)
		if err != nil {
			return fmt.Errorf("Migrate - Read %s: %w", name, err)
		}

		if _, err := db.ExecContext(ctx, string(body)); err != nil {
			return fmt.Errorf("Migrate - Execute %s: %w", name, err)
		}

		if _, err := db.ExecContext(ctx,
			`INSERT INTO schema_migrations (name, applied_at)
			 VALUES (?, datetime('now'))`,
			name); err != nil {
			return fmt.Errorf("Migrate - Record %s: %w", name, err)
		}
	}

	return nil
}
