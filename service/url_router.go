package service

import (
	"UrlShorter/config"
	"UrlShorter/storage"
	"UrlShorter/storage/repository"
	"fmt"
	"math/rand/v2"
	"net/url"
	"strings"
)

type UrlRouter struct {
	config *config.Config
	repo   *repository.StorageRepository
}

func NewUrlRouter(configuration *config.Config) (*UrlRouter, error) {

	db, err := storage.DbInit(configuration)
	if err != nil {
		return nil, fmt.Errorf("couldn't init DB: %v", err)
	}

	var storageRepo repository.StorageRepository = repository.NewStorageRepository(db, configuration)

	return &UrlRouter{
		config: configuration,
		repo:   &storageRepo,
	}, nil
}

func (urlRouter *UrlRouter) GetUrlByCode(code string) (string, error) {

	check, err := (*urlRouter.repo).IsExists(code)
	if err != nil {
		return "", err
	}

	if !check {
		return "", nil
	}

	baseUrl, err := (*urlRouter.repo).GetBaseUrl(code)
	if err != nil {
		return "", fmt.Errorf("couldn't get url: %v", err)
	}

	return baseUrl, nil
}

func (urlRouter *UrlRouter) AddUrl(baseUrl string) (string, string, error) {
	err := (*urlRouter.repo).ClearOldCode()
	if err != nil {
		return "", "", err
	}

	if !strings.HasPrefix(baseUrl, "https://") && !strings.HasPrefix(baseUrl, "http://") {
		baseUrl = "https://" + baseUrl
	}

	_, err = url.ParseRequestURI(baseUrl)
	if err != nil {
		return "", "", fmt.Errorf("wrong url '%s': %v", baseUrl, err)
	}

	exists := true
	var code string

	for exists {
		code = generateString(urlRouter.config.CodeLength)

		exists, err = (*urlRouter.repo).IsExists(code)

		if err != nil {
			return "", "", err
		}
	}

	err = (*urlRouter.repo).AddCode(code, baseUrl)
	if err != nil {
		return "", "", err
	}

	return code, baseUrl, nil
}

func generateString(length int) string {

	result := ""

	rangeSymbols := []rune("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890")

	for i := 0; i < length; i++ {
		randIndex := rand.IntN(len(rangeSymbols))
		result += string(rangeSymbols[randIndex])
	}

	return result
}
