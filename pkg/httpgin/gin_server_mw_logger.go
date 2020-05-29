package httpgin

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/labstack/gommon/log"
)

// MLogger Middleware logger.
func MLogger(cfg MConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request

		c.Next()

		// after request
		log.Print("Latency: ", time.Since(t))
		// access the status we are sending
		log.Print(c.Writer.Status())
	}
}
