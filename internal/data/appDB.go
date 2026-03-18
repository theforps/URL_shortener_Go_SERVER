package data

import (
	"database/sql"
	"url_shortener/internal/data/models"

	_ "github.com/lib/pq"
)

func Setup(conf *models.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.ConnectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
