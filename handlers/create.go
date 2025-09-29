package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"url_shortner/config"
	"url_shortner/handlers/entity"
	"url_shortner/service"
)

func Create(configuration *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			var url entity.URL
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

			resultCode, baseUrl, err := urlRouter.AddUrl(url.Url)
			if err != nil {
				log.Println(err)
				http.Error(w, "502 couldn't parse url", http.StatusBadGateway)
			}

			finallyUrl := GetFinnalyUrl(configuration, resultCode)

			responseUrl := &entity.URL{Url: finallyUrl}
			jsonData, err := json.Marshal(*responseUrl)
			if err != nil {
				http.Error(w, "502 couldn't convert url to json", http.StatusBadGateway)
			}

			userIp := ReadUserIP(r)
			log.Printf("%s create %s to %s",
				userIp,
				finallyUrl,
				baseUrl,
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
