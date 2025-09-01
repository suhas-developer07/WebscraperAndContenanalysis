package repository

import (
	"database/sql"

	"github.com/lib/pq"
)

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

// here i need to insert the jobs to the table created above whenever request comes from the api

type InsertJobsPayload struct {
	urls []string
}

func (r *PostgresRepository) InsertJobs(urls InsertJobsPayload) error {
	query := `INSERT INTO urls (url) VALUES ($1)`

	_, err := r.db.Exec(query, pq.Array(urls))
	return err
}
