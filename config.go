package main

import (
	"os"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	ConnectionString string
	Driver           string
	TableName        string
}

type Config struct {
	DbConfig *DbConfig
	TgApi    string
	Domain   string
}

func Configuration() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DbConfig: &DbConfig{
			ConnectionString: getEnv("DB", ""),
			Driver:           getEnv("DRIVER", ""),
			TableName:        getEnv("TABLE_NAME", ""),
		},
		TgApi:  getEnv("TG_BOT_API", ""),
		Domain: getEnv("DOMAIN", ""),
	}, nil
}

func getEnv(key string, defaultValume string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValume
}
