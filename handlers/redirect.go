package handlers

import (
	"log"
	"net/http"
	"url_shortner/config"
	"url_shortner/service"
)

func Redirect(configuration *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlRouter, err := service.NewUrlRouter(configuration)
		if err != nil {
			log.Println(err)
			http.Error(w, "502 couldn't get url", http.StatusBadGateway)
		}

		code := r.PathValue("code")

		baseUrl, err := urlRouter.GetUrlByCode(code)
		if err != nil {
			log.Println(err)
			http.Error(w, "502 couldn't get url", http.StatusBadGateway)
		}

		userIp := ReadUserIP(r)
		finnalyUrl := GetFinnalyUrl(configuration, code)

		log.Printf(
			"%s clicked from %s to %s",
			userIp,
			finnalyUrl,
			baseUrl,
		)
		http.Redirect(w, r, baseUrl, http.StatusMovedPermanently)
	}
}
