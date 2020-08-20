package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MReady Middleware. 503 if middleware applied.
// curl -L http://localhost:8001/k8/xxx
func MReady(cfg MConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Skipper() || isReady() {
			// execute pending handlers
			c.Next()
			// can now exit middleware
			return
		}

		// check if no service route so we do not rediect again
		if c.Request.URL.String() == EndPointGroupK8+EndPointNoService {
			return
		}

		// can now do the redirection to desired route
		c.Redirect(http.StatusSeeOther, EndPointGroupK8+EndPointNoService)
	}
}
