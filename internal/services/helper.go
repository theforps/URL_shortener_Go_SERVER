package services

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"url_shortener/internal/data/models"
)

func genCode(length int, configuration *models.Config) string {
	rangeSymbols := []rune(configuration.Symbols)
	result := make([]rune, length)

	for i := 0; i < length; i++ {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(rangeSymbols))))
		result[i] = rangeSymbols[nBig.Int64()]
	}

	return string(result)
}

func isValidURL(link string) bool {
	parsed, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}
	return parsed.Scheme != "" && parsed.Host != ""
}
