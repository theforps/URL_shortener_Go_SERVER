package main

import (
	"fmt"
	"log"
	"net/http"

	"url_shortner/config"
	"url_shortner/handlers"
)

func main() {

	configuration, err := config.Configuration()
	if err != nil {
		log.Println(err)
	}

	handler := newHandler(configuration)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", configuration.DomainConfig.PortDev),
		Handler: handler,
	}

	log.Printf("the server started on port %s", configuration.DomainConfig.PortDev)

	err = server.ListenAndServe()
	if err != nil {
		log.Printf("the server is stopped due to an error: %v", err)
	}
}

func newHandler(configuration *config.Config) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/create-url", handlers.Create(configuration))
	mux.HandleFunc("/{code}", handlers.Redirect(configuration))

	return mux
}
