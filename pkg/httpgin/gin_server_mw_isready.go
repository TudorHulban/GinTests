package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/labstack/gommon/log"
)

// MReady Middleware. 503 if middleware applied.
func MReady(cfg MConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Skipper() {
			// execute pending handlers
			c.Next()
			// can now exit middleware
			return
		}

		log.Print("mw MReady applied")
		c.Redirect(http.StatusSeeOther, EndPointGroupK8+"/noservice")
	}
}
