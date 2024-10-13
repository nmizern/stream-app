package middleware

import (
	"stream-app/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.SugaredLogger.Infof("Status: %d | Latency: %v | IP: %s | Method: %s | Path: %s",
			status, latency, clientIP, method, path)
	}
}
