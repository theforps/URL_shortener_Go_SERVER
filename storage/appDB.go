package storage

import (
	"database/sql"
	"fmt"
	"os"
	"url_shortner/config"

	_ "github.com/mattn/go-sqlite3"
)

// DbInit Initializing a storage connection
func DbInit(configuration *config.Config) (*sql.DB, error) {

	if !dbExists(configuration) {
		err := dbCreate(configuration)
		if err != nil {
			return nil, fmt.Errorf("couldn't create DB: %v", err)
		}
	}

	db, err := dbOpen(configuration)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// opens the sql connection
func dbOpen(configuration *config.Config) (*sql.DB, error) {
	db, err := sql.Open(configuration.DbConfig.Driver, configuration.DbConfig.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("couldn't open DB connection: %v", err)
	}
	return db, nil
}

// Checking the existence of the database
func dbExists(configuration *config.Config) bool {
	_, err := os.ReadFile(configuration.DbConfig.ConnectionString)
	if err != nil {
		return false
	}
	return true
}

// Create DB if not exists
func dbCreate(configuration *config.Config) error {

	_, err := os.Create(configuration.DbConfig.ConnectionString)
	if err != nil {
		return fmt.Errorf("couldn't create DB file %s: %v", configuration.DbConfig.ConnectionString, err)
	}

	db, err := sql.Open(configuration.DbConfig.Driver, configuration.DbConfig.ConnectionString)
	if err != nil {
		return fmt.Errorf("couldn't open DB connection: %v", err)
	}
	defer db.Close()

	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (Id INTEGER PRIMARY KEY AUTOINCREMENT, Code TEXT, UrlBase TEXT, FinallyDate TEXT)",
		configuration.DbConfig.TableName,
	)

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("couldn't create table %s: %v", configuration.DbConfig.TableName, err)
	}

	return nil
}
