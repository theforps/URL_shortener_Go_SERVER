package tests

import (
	"database/sql"
	"fmt"
	"time"
)

func waitForDB() (*sql.DB, error) {
	dsn := "postgres://test:test@localhost:4321/test_db?sslmode=disable"

	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("couldn't connect to db after %d retries: %v", 10, err)
}
