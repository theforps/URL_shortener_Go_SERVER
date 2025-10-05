package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"url_shortener/internal/config"
	"url_shortener/internal/handlers/entity"
	"url_shortener/internal/service"
)

func Create(configuration *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			var url entity.UrlDTO
			body := json.NewDecoder(r.Body)

			err := body.Decode(&url)
			if err != nil {
				log.Println(err)
				http.Error(w, "400 bad request", http.StatusBadRequest)
			}

			urlRouter, err := service.NewUrlRouter(configuration)
			if err != nil {
				log.Println(err)
				http.Error(w, "502 couldn't get url", http.StatusBadGateway)
			}

			currentDate := time.Now().Weekday()
			if currentDate == 0 {
				err = urlRouter.ClearOldUrls()
				if err != nil {
					log.Printf("couldn't clear old URLs: %v\n", err)
				}
			}

			generatedCode, userUrl, err := urlRouter.AddUrl(url.Url)
			if err != nil {
				log.Println(err)
				http.Error(w, "502 couldn't parse url", http.StatusBadGateway)
			}

			redirectUrl, days := GetRedirectUrl(configuration, generatedCode)
			responseUrlEntity := &entity.UrlDTO{Url: redirectUrl, DayLife: days}

			jsonData, err := json.Marshal(responseUrlEntity)
			if err != nil {
				http.Error(w, "502 couldn't convert url to json", http.StatusBadGateway)
			}

			userIp := ReadUserIP(r)
			log.Printf("user - %s create redirect %s -> %s",
				userIp,
				redirectUrl,
				userUrl,
			)

			if generatedCode != "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonData)
			}
		} else {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
