package repository

import (
	"database/sql"
	"fmt"
	"url_shortener/internal/data/models"
)

type StorageRepositoryDB struct {
	db *sql.DB
}

func NewStorageRepository(db *sql.DB) *StorageRepositoryDB {
	return &StorageRepositoryDB{
		db: db,
	}
}

func (sr *StorageRepositoryDB) IsExists(code string) (isExists bool, err error) {

	dbResponse := sr.db.QueryRow(
		"SELECT COUNT(uniq_code) FROM url_table WHERE uniq_code = $1;", code)

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

func (sr *StorageRepositoryDB) ClearOld() (err error) {

	resultSelect, err := sr.db.Query("SELECT id FROM url_table WHERE finally_date::timestamp < now()")
	if err != nil {
		return fmt.Errorf("couldn't select old Ids: %v", err)
	}

	var oldIds []int

	for resultSelect.Next() {
		var oldId int

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
			_, err = sr.db.Exec("DELETE FROM url_table WHERE id = $1", value)
			if err != nil {
				return fmt.Errorf("couldn't delete row with Id = %d: %v", value, err)
			}
		}
	}

	return nil
}

func (sr *StorageRepositoryDB) AddCode(code string, url string, finDate string) error {

	_, err := sr.db.Exec(
		"INSERT INTO url_table (uniq_code, url_base, views, finally_date) VALUES($1, $2, $3, $4);",
		code,
		url,
		0,
		finDate)
	if err != nil {
		return fmt.Errorf("couldn't add new code '%s' to DB: %v", code, err)
	}
	return nil
}

func (sr *StorageRepositoryDB) GetBaseUrl(code string) (string, error) {

	rows := sr.db.QueryRow(
		"SELECT url_base FROM url_table WHERE uniq_code = $1;", code)

	var baseUrl string
	err := rows.Scan(&baseUrl)
	if err != nil {
		return "", fmt.Errorf("couldn't get base url by code '%s': %v", code, err)
	}

	return baseUrl, nil
}

func (sr *StorageRepositoryDB) GetStats(code string) (*models.URLStats, error) {

	rows := sr.db.QueryRow(
		"SELECT views, finally_date FROM url_table WHERE uniq_code = $1;", code)

	var views int
	var finDate string
	err := rows.Scan(&views, &finDate)
	if err != nil {
		return nil, fmt.Errorf("couldn't get views by code '%s': %v", code, err)
	}

	return &models.URLStats{
		Views: views,
		FinallyDate: finDate,
	}, nil
}

func (sr *StorageRepositoryDB) IncreaseViews(code string) error {
	_, err := sr.db.Exec(
		"UPDATE url_table SET views = views + 1 WHERE uniq_code = $1", code)
	if err != nil {
		return fmt.Errorf("couldn't increase views '%s' to DB: %v", code, err)
	}
	return nil
}
