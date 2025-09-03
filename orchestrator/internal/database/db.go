package database

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func Connect(ConnStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", ConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to Connect Database; %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping Database : %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}
