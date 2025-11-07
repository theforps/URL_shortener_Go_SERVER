package service

import (
	"fmt"
	"math/rand/v2"
	"net/url"
	"strings"
	"url_shortener/internal/config"
	"url_shortener/internal/storage"
	repository2 "url_shortener/internal/storage/repository"
)

type UrlRouter struct {
	config *config.Config
	repo   repository2.StorageRepository
}

// NewUrlRouter creates new router
func NewUrlRouter(configuration *config.Config) (router *UrlRouter, err error) {

	db, err := storage.DbInit(configuration)
	if err != nil {
		return nil, fmt.Errorf("couldn't init DB: %v", err)
	}

	var storageRepo repository2.StorageRepository = repository2.NewStorageRepository(db, configuration)

	return &UrlRouter{
		config: configuration,
		repo:   storageRepo,
	}, nil
}

// GetUrlByCode uses a unique code to pull the user URL from the database
func (urlRouter *UrlRouter) GetUrlByCode(code string) (userUrl string, err error) {

	isExists, err := urlRouter.repo.IsExists(code)
	if err != nil {
		return "", fmt.Errorf("couldn't check url: %v", err)
	}
	if !isExists {
		return "", nil
	}

	userUrl, err = urlRouter.repo.GetBaseUrl(code)
	if err != nil {
		return "", fmt.Errorf("couldn't get url: %v", err)
	}

	return userUrl, nil
}

// AddUrl creates a unique code that binds to the user URL
func (urlRouter *UrlRouter) AddUrl(baseUrl string) (generatedCode string, userUrl string, err error) {
	if !strings.HasPrefix(baseUrl, "https://") && !strings.HasPrefix(baseUrl, "http://") {
		baseUrl = "https://" + baseUrl
	}

	_, err = url.ParseRequestURI(baseUrl)
	if err != nil {
		return "", "", fmt.Errorf("wrong url '%s': %v", baseUrl, err)
	}

	isExists := true

	for isExists {
		generatedCode = generateString(urlRouter.config.CodeLength, urlRouter.config)
		isExists, err = urlRouter.repo.IsExists(generatedCode)

		if err != nil {
			return "", "", err
		}
	}

	err = urlRouter.repo.AddCode(generatedCode, baseUrl)
	if err != nil {
		return "", "", err
	}

	return generatedCode, baseUrl, nil
}

// ClearOldUrls deletes old records in the database
func (urlRouter *UrlRouter) ClearOldUrls() (err error) {
	err = urlRouter.repo.ClearOld()
	if err != nil {
		return fmt.Errorf("couldn't clear old URLs: %v", err)
	}
	return nil
}

// generateString creates a unique code for the user
func generateString(length int, configuration *config.Config) (generatedCode string) {

	rangeSymbols := []rune(configuration.SymbolsBase)

	for i := 0; i < length; i++ {
		randIndex := rand.IntN(len(rangeSymbols))
		generatedCode += string(rangeSymbols[randIndex])
	}

	return generatedCode
}
