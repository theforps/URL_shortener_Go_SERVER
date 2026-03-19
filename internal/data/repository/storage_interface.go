package repository

import "url_shortener/internal/data/models"

type StorageRepository interface {
	IsExists(code string) (bool, error)

	ClearOld() error

	AddCode(code string, url string, finDate string) error

	GetBaseUrl(code string) (string, error)

	GetStats(code string) (*models.URLStats, error)

	IncreaseViews(code string) error
}
