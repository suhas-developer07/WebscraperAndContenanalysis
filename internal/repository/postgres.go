package repository

import "database/sql"

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) InitTable() error {
	query := `CREATE TABLE IF NOT EXISTS urls (
	id BIGSERIAL PRIMARY KEY,
	url TEXT[] 
	)`
	_, err := r.db.Exec(query)
	return err
}
