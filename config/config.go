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

type DomainConfig struct {
	Domain string
	Port   string
	Prefix string
}

type Config struct {
	DbConfig     *DbConfig
	DomainConfig *DomainConfig
	CodeLength   int
	UrlLifeDays  int
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
			ConnectionString: getEnvStr("DB", "urlDb"),
			Driver:           getEnvStr("DRIVER", "sqlite3"),
			TableName:        getEnvStr("TABLE_NAME", "url_short"),
		},
		DomainConfig: &DomainConfig{
			Domain: getEnvStr("DOMAIN", "localhost"),
			Port:   getEnvStr("PORT", "8080"),
			Prefix: getEnvStr("DOMAIN_PREFIX", "http"),
		},
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

func (configuration *Config) GetDomain() string {

	domain := fmt.Sprintf("%s://%s:%s",
		configuration.DomainConfig.Prefix,
		configuration.DomainConfig.Domain,
		configuration.DomainConfig.Port)
	return domain
}
