package repository

import (
	"UrlShorter/config"
	"database/sql"
	"fmt"
	"time"
)

type StorageRepositoryDB struct {
	DB     *sql.DB
	Config *config.Config
}

func NewStorageRepository(db *sql.DB, config *config.Config) *StorageRepositoryDB {
	return &StorageRepositoryDB{
		DB:     db,
		Config: config,
	}
}

func (storageRepo *StorageRepositoryDB) IsExists(code string) (bool, error) {

	query := fmt.Sprintf("SELECT COUNT(Code) FROM %s WHERE Code = '%s' LIMIT 1", storageRepo.Config.DbConfig.TableName, code)

	result := storageRepo.DB.QueryRow(query)

	var shorterCode *int

	err := result.Scan(&shorterCode)
	if err != nil {
		return false, fmt.Errorf("couldn't check code: %v", err)
	}

	if *shorterCode > 0 {
		return true, nil
	}

	return false, nil
}

func (storageRepo *StorageRepositoryDB) ClearOldCode() error {
	query := fmt.Sprintf("SELECT Id FROM %s WHERE datetime(FinallyDate) < DATETIME('now')", storageRepo.Config.DbConfig.TableName)

	resultSelect, err := storageRepo.DB.Query(query)
	if err != nil {
		return fmt.Errorf("couldn't select old Ids: %v", err)
	}

	var oldIds []*int

	for resultSelect.Next() {
		var oldId *int

		err = resultSelect.Scan(&oldId)
		if err != nil {
			return fmt.Errorf("couldn't get old Id: %v", err)
		}
		oldIds = append(oldIds, oldId)
	}

	err = resultSelect.Close()
	if err != nil {
		return fmt.Errorf("couldn't close connection to DB (old IDs): %v", err)
	}

	if len(oldIds) > 0 {
		for _, value := range oldIds {

			query = fmt.Sprintf("DELETE FROM %s WHERE Id = %d", storageRepo.Config.DbConfig.TableName, *value)

			_, err := storageRepo.DB.Exec(query)
			if err != nil {
				return fmt.Errorf("couldn't delete row with Id = %d: %v", *value, err)
			}
		}
	}

	return nil
}

func (storageRepo *StorageRepositoryDB) AddCode(code string, url string) error {

	date := time.Now().AddDate(0, 0, storageRepo.Config.UrlLifeDays).UTC()
	dateFormat := date.Format("2006-01-02 15:04:05")

	query := fmt.Sprintf("INSERT INTO %s (Code, UrlBase, FinallyDate) VALUES('%s', '%s', '%s')", storageRepo.Config.DbConfig.TableName, code, url, dateFormat)

	_, err := storageRepo.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("couldn't add new code '%s' to DB: %v", code, err)
	}
	return nil
}
