package database

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const (
	sqlDriver       = "pgx"
	gooseDialect    = "postgres"
	migrationFolder = "migrations/" // look at the embedMigrations path
)

func init() {
	goose.SetBaseFS(embedMigrations)
}

// Config represent postgres configuration.
type Config struct {
	DSN string
}

// Validate check if postgres configuration is valid.
func (c *Config) Validate() error {
	if c.DSN == "" {
		return fmt.Errorf("dsn can't be empty")
	}
	return nil
}

// Connect create a new postgres connection.
func Connect(conf Config) (*sql.DB, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	s, err := sql.Open(sqlDriver, conf.DSN)
	if err != nil {
		return nil, err
	}

	if err := s.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}
	return s, nil
}

// UpMigration run postgres migration.
func UpMigration(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("argument db can't be nil")
	}

	if err := goose.SetDialect(gooseDialect); err != nil {
		return fmt.Errorf("error setting goose dialect: %w", err)
	}
	if err := goose.Up(db, migrationFolder); err != nil {
		return fmt.Errorf("error running goose up: %w", err)
	}

	return nil
}
