package main

import (
	"log"

	"url_shortener/internal/config"
	"url_shortener/internal/data"
	"url_shortener/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {

	// init logger
	logger := log.Default()

	// init config
	configuration, err := config.Default()
	if err != nil {
		logger.Println(err)
	}

	// init db
	db, err := data.Setup(configuration)
	if err != nil {
		logger.Fatalln(err)
	}

	if configuration.Mode == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	// init web
	engine := gin.Default()

	// init handler
	handler.Init(engine, configuration, db, logger)

	// check status
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "url-api",
		})
	})

	logger.Println("the server started on port 5050")

	// run server
	err = engine.Run(":5050")
	if err != nil {
		logger.Printf("the server is stopped due to an error: %v", err)
	}
}
