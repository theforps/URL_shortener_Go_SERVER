package handlers

import (
	"fmt"
	"net/http"
	"url_shortner/config"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func GetFinnalyUrl(configuration *config.Config, code string) string {
	finallyUrl := ""

	if configuration.Lvl == "PROD" {
		finallyUrl = fmt.Sprintf(
			"%s://%s/%s",
			configuration.DomainConfig.PrefixProd,
			configuration.DomainConfig.DomainProd,
			code,
		)
	} else {
		finallyUrl = fmt.Sprintf(
			"%s://%s:%s/%s",
			configuration.DomainConfig.PrefixDev,
			configuration.DomainConfig.DomainDev,
			configuration.DomainConfig.PortDev,
			code,
		)
	}

	return finallyUrl
}
