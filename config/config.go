package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	ConnectionString string
	Driver           string
	TableName        string
}

type Config struct {
	DbConfig    *DbConfig
	TgApi       string
	Domain      string
	CodeLength  int
	UrlLifeDays int
}

func Configuration() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("couldn't load .env: %v", err)
	}

	codeLength, err := getEnvInt("CODE_LENGTH", 6)
	if err != nil {
		return nil, fmt.Errorf("couldn't convert environment %#v: %v", codeLength, err)
	}

	urlLifeDays, err := getEnvInt("URL_LIFE_DAYS", 6)
	if err != nil {
		return nil, fmt.Errorf("couldn't convert environment %#v: %v", urlLifeDays, err)
	}

	return &Config{
		DbConfig: &DbConfig{
			ConnectionString: getEnvStr("DB", ""),
			Driver:           getEnvStr("DRIVER", ""),
			TableName:        getEnvStr("TABLE_NAME", ""),
		},
		TgApi:       getEnvStr("TG_BOT_API", ""),
		Domain:      getEnvStr("DOMAIN", ""),
		CodeLength:  codeLength,
		UrlLifeDays: urlLifeDays,
	}, nil
}

func getEnvStr(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		resConv, err := strconv.Atoi(value)
		return resConv, err
	}
	return defaultValue, nil
}
