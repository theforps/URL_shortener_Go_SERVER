package services

import "url_shortener/internal/data/models"

type ManagerService interface {
	GetUrl(code string) (string, error)

	GetStats(uniqCode string) *models.Response[models.URLStats]

	AddUrl(request *models.Request) *models.Response[models.URL]

	AddView(code string) error
}
