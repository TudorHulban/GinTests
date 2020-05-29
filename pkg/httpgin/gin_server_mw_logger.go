package httpgin

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/labstack/gommon/log"
)

// MLogger Middleware logger.
// See example of logger in https://stackoverflow.com/questions/50574796/gin-gonic-middleware-declaration#50575548 .
func MLogger(cfg MConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// executes pending handlers.
		c.Next()

		if cfg.Skipper() {
			return
		}

		log.Print("mw MLogger applied")

		t := time.Now()
		// before request

		// after request
		log.Print("Latency: ", time.Since(t))
		// access the status we are sending
		log.Print(c.Writer.Status())
	}
}
