package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"url_shortener/internal/data/models"
	"url_shortener/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UrlHandler struct {
	routeService  services.ManagerService
	configuration *models.Config
	db            *sql.DB
	logger        *log.Logger
	validate      *validator.Validate
}

func Init(
	engine *gin.Engine,
	configuration *models.Config,
	db *sql.DB,
	logger *log.Logger) {

	uh := &UrlHandler{
		routeService:  services.NewUrlService(configuration, db, logger),
		configuration: configuration,
		db:            db,
		logger:        logger,
		validate:      validator.New(),
	}

	engine.POST("/shorten", uh.CreateUrl)
	engine.GET("/:short_id", uh.Redirect)
	engine.GET("/stats/:short_id", uh.Stats)
}

// принимает длинную ссылку, возвращает короткий идентификатор
func (uh *UrlHandler) CreateUrl(c *gin.Context) {

	var request models.Request
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		uh.logger.Println(err)

		cookResponse(c, uh.logger, &models.Response[string]{
			StatusCode:  400,
			Description: "bad request",
		})
	}

	err := uh.validate.Struct(request)
	if err != nil {
		cookResponse(c, uh.logger, &models.Response[string]{
			StatusCode:  400,
			Description: err.Error(),
		})
	} else {
		response := uh.routeService.AddUrl(&request)
		if response.StatusCode == 200 {
			userIp := readUserIP(c)
			uh.logger.Printf("user - %s create redirect %s -> %s",
				userIp,
				response.Data.UniqCode,
				request.BaseUrl,
			)
		}
		cookResponse(c, uh.logger, response)
	}
}

// возвращает количество переходов
func (uh *UrlHandler) Stats(c *gin.Context) {

	uniq_code := c.Param("short_id")
	if uniq_code == "" {
		uh.logger.Println("uniq code is empty")

		cookResponse(c, uh.logger, &models.Response[string]{
			StatusCode:  400,
			Description: "bad request",
		})
	} else {
		response := uh.routeService.GetStats(uniq_code)
		cookResponse(c, uh.logger, response)
	}
}

// редиректит на оригинальную ссылку
func (uh *UrlHandler) Redirect(c *gin.Context) {
	uniq_code := c.Param("short_id")
	if uniq_code == "" {
		uh.logger.Println("uniq code is empty")

		cookResponse(c, uh.logger, &models.Response[string]{
			StatusCode:  400,
			Description: "bad request",
		})
	} else {
		baseUrl, err := uh.routeService.GetUrl(uniq_code)
		if err != nil {
			uh.logger.Println(err)

			cookResponse(c, uh.logger, &models.Response[string]{
				StatusCode:  500,
				Description: "internal error",
			})
		} else {
			err := uh.routeService.AddView(uniq_code)
			if err != nil {
				uh.logger.Println(err)

				cookResponse(c, uh.logger, &models.Response[string]{
					StatusCode:  500,
					Description: "internal error",
				})
			} else {
				userIp := readUserIP(c)
				uh.logger.Printf(
					"user - %s redirected %s -> %s",
					userIp,
					uniq_code,
					baseUrl,
				)

				c.Redirect(http.StatusMovedPermanently, baseUrl)
			}
		}
	}
}
