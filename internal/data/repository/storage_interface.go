package repository

type StorageRepository interface {
	IsExists(code string) (bool, error)

	ClearOld() error

	AddCode(code string, url string, finDate string) error

	GetBaseUrl(code string) (string, error)

	GetViews(code string) (int, error)

	IncreaseViews(code string) error
}
