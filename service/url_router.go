package service

import (
	"fmt"
	"math/rand/v2"
	"net/url"
	"strings"

	"url_shortener/config"
	"url_shortener/storage"
	"url_shortener/storage/repository"
)

type UrlRouter struct {
	config *config.Config
	repo   repository.StorageRepository
}

func NewUrlRouter(configuration *config.Config) (router *UrlRouter, err error) {

	db, err := storage.DbInit(configuration)
	if err != nil {
		return nil, fmt.Errorf("couldn't init DB: %v", err)
	}

	var storageRepo repository.StorageRepository = repository.NewStorageRepository(db, configuration)

	return &UrlRouter{
		config: configuration,
		repo:   storageRepo,
	}, nil
}

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

func (urlRouter *UrlRouter) ClearOldUrls() (err error) {
	err = urlRouter.repo.ClearOld()
	if err != nil {
		return fmt.Errorf("couldn't clear old URLs: %v", err)
	}
	return nil
}

func generateString(length int, configuration *config.Config) (generatedCode string) {

	rangeSymbols := []rune(configuration.SymbolsBase)

	for i := 0; i < length; i++ {
		randIndex := rand.IntN(len(rangeSymbols))
		generatedCode += string(rangeSymbols[randIndex])
	}

	return generatedCode
}
