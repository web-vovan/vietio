package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrationsFS embed.FS

func Up(db *sql.DB) error {
    if err := goose.SetDialect("postgres"); err != nil {
        return fmt.Errorf("set goose dialect: %w", err)
    }

    goose.SetBaseFS(migrationsFS)

    if err := goose.Up(db, "."); err != nil {
        return fmt.Errorf("run migrations: %w", err)
    }

    return nil
}

func Reset(db *sql.DB) error {
    if err := goose.SetDialect("postgres"); err != nil {
        return fmt.Errorf("set goose dialect: %w", err)
    }

    goose.SetBaseFS(migrationsFS)

    if err := goose.Reset(db, "."); err != nil {
        return fmt.Errorf("reset migrations: %w", err)
    }

    return nil
}
