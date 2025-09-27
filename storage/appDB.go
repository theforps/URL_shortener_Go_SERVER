package storage

import (
	"UrlShorter/config"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DbInit Initializing a storage connection
func DbInit() (*sql.DB, error) {

	err := dbCreate()
	if err != nil {
		return nil, fmt.Errorf("couldn't create DB: %v", err)
	}

	configuration, err := config.Configuration()
	if err != nil {
		return nil, fmt.Errorf("couldn't get config: %v", err)
	}

	db, err := sql.Open(configuration.DbConfig.Driver, configuration.DbConfig.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("couldn't open DB connection: %v", err)
	}

	return db, nil
}

// DbCreate Create DB if not exists
func dbCreate() error {

	configuration, err := config.Configuration()
	if err != nil {
		return fmt.Errorf("couldn't get config: %v", err)
	}

	_, err = os.ReadFile(configuration.DbConfig.ConnectionString)
	if err != nil {
		_, err = os.Create(configuration.DbConfig.ConnectionString)
		if err != nil {
			return fmt.Errorf("couldn't create DB file %s: %v", configuration.DbConfig.ConnectionString, err)
		}
	}

	db, err := sql.Open(configuration.DbConfig.Driver, configuration.DbConfig.ConnectionString)
	if err != nil {
		return fmt.Errorf("couldn't open DB connection: %v", err)
	}

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
