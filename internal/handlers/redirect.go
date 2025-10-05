package handlers

import (
	"log"
	"net/http"
	"url_shortener/internal/config"
	"url_shortener/internal/service"
)

func Redirect(configuration *config.Config) (handler http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		urlRouter, err := service.NewUrlRouter(configuration)
		if err != nil {
			log.Println(err)
			http.Error(w, "502 couldn't get url", http.StatusBadGateway)
		}

		code := r.PathValue("code")

		userUrl, err := urlRouter.GetUrlByCode(code)
		if err != nil {
			log.Println(err)
			http.Error(w, "502 couldn't get url", http.StatusBadGateway)
		}

		userIp := ReadUserIP(r)
		redirectUrl, _ := GetRedirectUrl(configuration, code)

		log.Printf(
			"user - %s redirected %s -> %s",
			userIp,
			redirectUrl,
			userUrl,
		)
		http.Redirect(w, r, userUrl, http.StatusMovedPermanently)
	}
}
