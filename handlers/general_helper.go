package handlers

import (
	"fmt"
	"net/http"
	"url_shortner/config"
)

func ReadUserIP(r *http.Request) (userIp string) {
	userIp = r.Header.Get("X-Real-Ip")
	if userIp == "" {
		userIp = r.Header.Get("X-Forwarded-For")
	}
	if userIp == "" {
		userIp = r.RemoteAddr
	}
	return
}

func GetRedirectUrl(configuration *config.Config, code string) (redirectUrl string) {
	if configuration.Lvl == "PROD" {
		redirectUrl = fmt.Sprintf(
			"%s://%s/%s",
			configuration.DomainConfig.PrefixProd,
			configuration.DomainConfig.DomainProd,
			code,
		)
	} else {
		redirectUrl = fmt.Sprintf(
			"%s://%s:%s/%s",
			configuration.DomainConfig.PrefixDev,
			configuration.DomainConfig.DomainDev,
			configuration.DomainConfig.PortDev,
			code,
		)
	}

	return
}
