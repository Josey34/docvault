package factory

import (
	"database/sql"
	"docvault/config"
	"docvault/database"
	"fmt"
)

type Factory struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Factory, error) {
	db, err := database.NewSQLite(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to inititalize sqlite: %w", err)
	}
	if err := database.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	return &Factory{
		DB: db,
	}, nil
}
