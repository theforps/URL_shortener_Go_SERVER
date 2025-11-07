package repository

import (
	"database/sql"
	"fmt"
	"time"
	"url_shortener/internal/config"
)

type StorageRepositoryDB struct {
	db     *sql.DB
	config *config.Config
}

// NewStorageRepository creates an object for accessing the database
func NewStorageRepository(db *sql.DB, config *config.Config) *StorageRepositoryDB {
	return &StorageRepositoryDB{
		db:     db,
		config: config,
	}
}

func (storageRepo *StorageRepositoryDB) IsExists(code string) (isExists bool, err error) {

	dbResponse := storageRepo.db.QueryRow("SELECT COUNT(CODE) FROM SHORT_TABLE WHERE Code = ? LIMIT 1", code)

	var shorterCode int
	err = dbResponse.Scan(&shorterCode)
	if err != nil {
		return false, fmt.Errorf("couldn't check code: %v", err)
	}
	if shorterCode > 0 {
		return true, nil
	}

	return false, nil
}

func (storageRepo *StorageRepositoryDB) ClearOld() (err error) {

	resultSelect, err := storageRepo.db.Query("SELECT ID FROM SHORT_TABLE WHERE datetime(FINALLY_DATE) < DATETIME('now')")
	if err != nil {
		return fmt.Errorf("couldn't select old Ids: %v", err)
	}

	var oldIds []*int

	for resultSelect.Next() {
		var oldId int

		err = resultSelect.Scan(&oldId)
		if err != nil {
			return fmt.Errorf("couldn't get old Id: %v", err)
		}
		oldIds = append(oldIds, &oldId)
	}

	err = resultSelect.Close()
	if err != nil {
		return fmt.Errorf("couldn't close connection to DB (old IDs): %v", err)
	}

	if len(oldIds) > 0 {
		for _, value := range oldIds {
			_, err = storageRepo.db.Exec("DELETE FROM SHORT_TABLE WHERE ID = ?", *value)
			if err != nil {
				return fmt.Errorf("couldn't delete row with Id = %d: %v", *value, err)
			}
		}
	}

	return nil
}

func (storageRepo *StorageRepositoryDB) AddCode(code string, url string) error {

	date := time.Now().AddDate(0, 0, storageRepo.config.UrlLifeDays).UTC()
	dateFormat := date.Format("2006-01-02 15:04:05")

	_, err := storageRepo.db.Exec("INSERT INTO SHORT_TABLE (CODE, URL_BASE, FINALLY_DATE) VALUES(?, ?, ?)", code, url, dateFormat)
	if err != nil {
		return fmt.Errorf("couldn't add new code '%s' to DB: %v", code, err)
	}
	return nil
}

func (storageRepo *StorageRepositoryDB) GetBaseUrl(code string) (string, error) {

	rows := storageRepo.db.QueryRow("SELECT URL_BASE FROM SHORT_TABLE WHERE CODE = ?", code)

	var baseUrl *string

	err := rows.Scan(&baseUrl)
	if err != nil {
		return "", fmt.Errorf("couldn't get base url by code '%s': %v", code, err)
	}

	return *baseUrl, nil
}
