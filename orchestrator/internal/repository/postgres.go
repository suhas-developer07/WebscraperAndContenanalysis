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
	Urls []string `json:"url"`
}

func (r *PostgresRepository) InsertJobs(payload InsertJobsPayload) (int64, []string, error) {
	query := `INSERT INTO urls (url) VALUES ($1) RETURNING id, url`
	var id int64
	var urls []string

	err := r.db.QueryRow(query, pq.Array(payload.Urls)).Scan(&id, pq.Array(&urls))

	if err != nil {
		return 0, nil, err
	}
	return id, urls, err
}
