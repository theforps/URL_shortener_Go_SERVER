package storage

import (
	"database/sql"
	"fmt"
	"os"
	"url_shortener/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

// DbInit Initializing a storage connection
func DbInit(configuration *config.Config) (db *sql.DB, err error) {

	if !dbExists(configuration) {
		err = dbCreate(configuration)
		if err != nil {
			return nil, fmt.Errorf("couldn't create DB: %v", err)
		}
	}

	db, err = dbOpen(configuration)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// opens the sql connection
func dbOpen(configuration *config.Config) (db *sql.DB, err error) {
	db, err = sql.Open(configuration.DbConfig.Driver, configuration.DbConfig.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("couldn't open DB connection: %v", err)
	}
	return db, nil
}

// Checking the existence of the database
func dbExists(configuration *config.Config) (isExists bool) {
	_, err := os.ReadFile(configuration.DbConfig.ConnectionString)
	return err == nil
}

// Create DB if not exists
func dbCreate(configuration *config.Config) (err error) {

	_, err = os.Create(configuration.DbConfig.ConnectionString)
	if err != nil {
		err = fmt.Errorf("couldn't create DB file %s: %v", configuration.DbConfig.ConnectionString, err)
		return
	}

	db, err := sql.Open(configuration.DbConfig.Driver, configuration.DbConfig.ConnectionString)
	if err != nil {
		err = fmt.Errorf("couldn't open DB connection: %v", err)
		return
	}
	defer func() {
		err = db.Close()
		if err != nil {
			err = fmt.Errorf("couldn't clode DB connection: %v", err)
			return
		}
	}()

	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (Id INTEGER PRIMARY KEY AUTOINCREMENT, Code TEXT, UrlBase TEXT, FinallyDate TEXT)",
		configuration.DbConfig.TableName,
	)

	_, err = db.Exec(query)
	if err != nil {
		err = fmt.Errorf("couldn't create table %s: %v", configuration.DbConfig.TableName, err)
		return
	}

	return nil
}
