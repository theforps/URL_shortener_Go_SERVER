package main

import (
	"UrlShorter/config"
	"UrlShorter/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type URL struct {
	Url string `json:"url"`
}

func main() {

	configuration, err := config.Configuration()
	if err != nil {
		log.Println(err)
	}

	handler := newHandler(configuration)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", configuration.DomainConfig.Domain, configuration.DomainConfig.Port),
		Handler: handler,
	}

	log.Printf("server is starting using %s", configuration.GetDomain())

	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func newHandler(configuration *config.Config) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/create-url", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			var url URL
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

			finallyUrl := fmt.Sprintf(
				"%s/%s",
				configuration.GetDomain(),
				resultCode,
			)
			responseBody := fmt.Sprintf("{\"url\":\"%s\"}", finallyUrl)

			userIp := readUserIP(r)
			log.Printf("%s create %s to %s",
				userIp,
				finallyUrl,
				baseUrl,
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(responseBody))
		} else {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		}

	})

	mux.HandleFunc("/{code}", func(w http.ResponseWriter, r *http.Request) {

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

		userIp := readUserIP(r)

		log.Printf(
			"%s clicked from %s/%s to %s",
			userIp,
			configuration.GetDomain(),
			code,
			baseUrl,
		)
		http.Redirect(w, r, baseUrl, http.StatusMovedPermanently)
	})

	return mux
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
