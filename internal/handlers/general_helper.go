package handlers

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"url_shortener/internal/config"
)

// ReadUserIP is a method that helps to read the user's IP address
func ReadUserIP(r *http.Request) (userIp string) {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	xRealIP := r.Header.Get("X-Real-Ip")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// GetRedirectUrl is a method that creates a server url for the user
func GetRedirectUrl(configuration *config.Config, code string) (redirectUrl string, days int) {
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

	days = configuration.UrlLifeDays

	return
}
