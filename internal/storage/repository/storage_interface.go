package repository

type StorageRepository interface {
	// IsExists verifies the existence of a record in the database using a unique code
	IsExists(code string) (bool, error)
	// ClearOld cleans the database of old records 
	ClearOld() error
	// AddCode add a new record to the database
	AddCode(code string, url string) error
	// GetBase read  the record from database
	GetBaseUrl(code string) (string, error)
}
