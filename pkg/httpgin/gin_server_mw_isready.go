package httpgin

import (
	"github.com/gin-gonic/gin"
	log "github.com/labstack/gommon/log"
)

// MReady Middleware. 200 if ready. 503 if not.
func MReady(cfg MConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// execute pending handlers
		c.Next()

		return

		log.Print("mw applied")
		c.Writer.Status()
	}
}
