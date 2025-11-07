package storage

import (
	"database/sql"
	"fmt"
	"os"
	"url_shortener/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

// DbInit initializes a storage connection
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

// dbOpen opens the sql connection
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

// Create a new database file if not exists
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS SHORT_TABLE (ID INTEGER PRIMARY KEY AUTOINCREMENT, CODE TEXT, URL_BASE TEXT, FINALLY_DATE TEXT)")
	if err != nil {
		err = fmt.Errorf("couldn't create table: %v", err)
		return
	}

	return nil
}
