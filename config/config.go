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

func Configuration() (config *Config, err error) {

	err = godotenv.Load()
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

func getEnvStr(key string, defaultValue string) (defValue string) {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) (defValue int, err error) {
	if value, exists := os.LookupEnv(key); exists {
		resConv, err := strconv.Atoi(value)
		return resConv, err
	}
	return defaultValue, nil
}
