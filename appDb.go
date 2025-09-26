package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DbInit Initializing a database connection
func DbInit() (*sql.DB, error) {

	config, err := Configuration()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(config.DbConfig.Driver, config.DbConfig.ConnectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// DbCreate Create DB if not exists
func DbCreate() error {

	config, err := Configuration()
	if err != nil {
		return err
	}

	_, err = os.ReadFile(config.DbConfig.ConnectionString)
	if err != nil {
		_, err = os.Create(config.DbConfig.ConnectionString)
		if err != nil {
			return err
		}
	}

	db, err := sql.Open(config.DbConfig.Driver, config.DbConfig.ConnectionString)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS ?  (Id INTEGER PRIMARY KEY AUTOINCREMENT, UrlBase TEXT, UrlCustom TEXT, FinallyDate INTEGER)",
		config.DbConfig.TableName)
	if err != nil {
		return err
	}

	return nil
}
