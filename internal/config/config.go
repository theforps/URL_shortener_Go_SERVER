package config

import (
	"fmt"
	"os"

	"url_shortener/internal/data/models"

	"github.com/joho/godotenv"
)

func Default() (*models.Config, error) {

	err := godotenv.Load("../.env")
	if err != nil {
		err = fmt.Errorf("couldn't load .env: %v", err)
	}

	return &models.Config{
		ConnectionString: getEnvStr("CONNECTION_STRING", "postgres://user:password@postgres:5432/db?sslmode=disable"),
		Mode:             getEnvStr("MODE", "DEV"),
		Symbols:          getEnvStr("SYMBOLS", "123"),
	}, err
}

func getEnvStr(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultVal
}
