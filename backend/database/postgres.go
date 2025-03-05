package database

import (
	"database/sql"
	"fmt"
	"log"
	"run-goals/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenPG(config *config.Config, logger *log.Logger) *sql.DB {
	db, err := sql.Open(
		"pgx",
		String(config.Database),
	)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}

func String(cfg config.Database) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
}
