package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"url_shortener/internal/data/models"
	"url_shortener/internal/data/repository"
)

type UrlService struct {
	config     *models.Config
	logger     *log.Logger
	repository repository.StorageRepository
}

func NewUrlService(configuration *models.Config, db *sql.DB, logger *log.Logger) *UrlService {
	uh := &UrlService{
		config:     configuration,
		logger:     logger,
		repository: repository.NewStorageRepository(db),
	}

	go func() {
		ticker := time.NewTicker(time.Hour * 12)
		defer ticker.Stop()

		for range ticker.C {
			err := uh.repository.ClearOld()
			if err != nil {
				uh.logger.Printf("couldn't clear old URLs: %v", err)
			}
		}
	}()

	return uh
}

func (us *UrlService) GetUrl(code string) (string, error) {

	isExists, err := us.repository.IsExists(code)
	if err != nil {
		return "", fmt.Errorf("couldn't check url: %v", err)
	}
	if !isExists {
		return "", nil
	}

	baseUrl, err := us.repository.GetBaseUrl(code)
	if err != nil {
		return "", fmt.Errorf("couldn't get url: %v", err)
	}

	return baseUrl, nil
}

func (us *UrlService) AddView(code string) error {
	err := us.repository.IncreaseViews(code)
	if err != nil {
		return fmt.Errorf("coudn't add view: %v", err)
	}
	return nil
}

func (us *UrlService) AddUrl(request *models.Request) *models.Response[models.URL] {

	if !isValidURL(request.BaseUrl) {
		us.logger.Printf("wrong url '%s'", request.BaseUrl)

		return &models.Response[models.URL]{
			StatusCode:  400,
			Description: "bad request",
		}
	}

	isExists := true
	var generatedCode string
	var err error

	for isExists {
		generatedCode = genCode(request.CodeLength, us.config)
		isExists, err = us.repository.IsExists(generatedCode)

		if err != nil {
			us.logger.Printf("code generation error '%s': %v", request.BaseUrl, err)

			return &models.Response[models.URL]{
				StatusCode:  500,
				Description: "internal error",
			}
		}
	}

	date := time.Now().AddDate(0, 0, request.DayLife).UTC()
	dateFormat := date.Format("2006-01-02 15:04:05")
	err = us.repository.AddCode(generatedCode, request.BaseUrl, dateFormat)
	if err != nil {
		us.logger.Printf("add url error: %v", err)

		return &models.Response[models.URL]{
			StatusCode:  500,
			Description: "internal error",
		}
	}

	return &models.Response[models.URL]{
		StatusCode:  200,
		Description: "ok",
		Data: &models.URL{
			UniqCode:    generatedCode,
			FinallyDate: dateFormat,
		},
	}
}

func (us *UrlService) GetStats(uniqCode string) *models.Response[models.URLStats] {

	isExists, err := us.repository.IsExists(uniqCode)
	if err != nil {
		us.logger.Printf("check url error: %v", err)

		return &models.Response[models.URLStats]{
			StatusCode:  500,
			Description: "internal error",
		}
	}

	if !isExists {
		return &models.Response[models.URLStats]{
			StatusCode:  404,
			Description: "not found",
		}
	}

	stats, err := us.repository.GetStats(uniqCode)
	if err != nil {
		us.logger.Printf("get views error: %v", err)

		return &models.Response[models.URLStats]{
			StatusCode:  500,
			Description: "internal error",
		}
	}

	return &models.Response[models.URLStats]{
		StatusCode:  200,
		Description: "ok",
		Data:        stats,
	}
}
