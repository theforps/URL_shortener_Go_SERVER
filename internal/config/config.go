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
}

type DomainConfig struct {
	DomainDev  string
	PortDev    string
	PrefixDev  string
	DomainProd string
	PrefixProd string
}

type Config struct {
	DbConfig     *DbConfig
	DomainConfig *DomainConfig
	CodeLength   int
	UrlLifeDays  int
	Lvl          string
	SymbolsBase  string
}

// NewConfiguration create new object from .env file
func NewConfiguration() (config *Config, err error) {

	err = godotenv.Load("../.env")
	if err != nil {
		return nil, fmt.Errorf("couldn't load .env: %v", err)
	}

	codeLength, err := getEnvInt("CODE_LENGTH", 6)
	if err != nil {
		return nil, fmt.Errorf("couldn't convert environment %#v: %v", codeLength, err)
	}

	urlLifeDays, err := getEnvInt("URL_LIFE_DAYS", 7)
	if err != nil {
		return nil, fmt.Errorf("couldn't convert environment %#v: %v", urlLifeDays, err)
	}

	return &Config{
		DbConfig: &DbConfig{
			ConnectionString: getEnvStr("DB", "urlDb"),
			Driver:           getEnvStr("DRIVER", "sqlite3"),
		},
		DomainConfig: &DomainConfig{
			DomainDev:  getEnvStr("DOMAIN_DEV", "localhost"),
			PortDev:    getEnvStr("PORT_DEV", "5050"),
			PrefixDev:  getEnvStr("PREFIX_DEV", "http"),
			DomainProd: getEnvStr("DOMAIN_PROD", "example.com"),
			PrefixProd: getEnvStr("PREFIX_PROD", "https"),
		},
		CodeLength:  codeLength,
		UrlLifeDays: urlLifeDays,
		Lvl:         getEnvStr("LVL", "DEV"),
		SymbolsBase: getEnvStr("SYMBOLS_BASE", "123"),
	}, nil
}

// getEnvStr read string from file
func getEnvStr(key string, defaultValue string) (defValue string) {
	defValue, isExists := os.LookupEnv(key)
	if isExists {
		return
	}

	return defaultValue
}

// getEnvInt read int from file
func getEnvInt(key string, defaultValue int) (defValue int, err error) {
	if value, exists := os.LookupEnv(key); exists {
		defValue, err = strconv.Atoi(value)
		return
	}
	return defaultValue, nil
}
