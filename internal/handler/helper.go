package handler

import (
	"encoding/json"
	"log"
	"net"
	"strings"
	"url_shortener/internal/data/models"

	"github.com/gin-gonic/gin"
)

func readUserIP(c *gin.Context) (userIp string) {
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	xRealIP := c.GetHeader("X-Real-Ip")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

func cookResponse[T any](c *gin.Context, logger *log.Logger, response *models.Response[T]) {
	byteResponse, err := json.Marshal(response)
	if err != nil {
		logger.Println(err)
		c.Writer.WriteHeader(500)
	} else {
		c.Writer.WriteHeader(response.StatusCode)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(byteResponse)
	}

	c.Writer.WriteHeaderNow()
}
