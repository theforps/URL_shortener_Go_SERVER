package repository

type StorageRepository interface {
	IsExists(code string) (bool, error)
	ClearOldCode() error
	AddCode(code string, url string) error
	GetBaseUrl(code string) (string, error)
}
