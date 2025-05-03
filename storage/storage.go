package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func NewPostgres(dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{DB: db}, nil
}

func (s *Storage) InitSchema() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS json_logs (
			id SERIAL PRIMARY KEY,
			data JSONB,
			timestamp TIMESTAMPTZ DEFAULT now()
		);
	`)
	return err
}

func (s *Storage) Insert(data []byte) error {
	_, err := s.DB.Exec("INSERT INTO json_logs (data) VALUES ($1)", string(data))
	return err
}
