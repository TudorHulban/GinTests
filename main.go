package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	cfg := config{
		socket: "127.0.0.1:8001",
	}

	r1 := route{
		group:    k8,
		endpoint: "/xxx",
		method:   "GET",
		handler:  func(c *gin.Context) { c.String(http.StatusOK, "xxx") },
	}

	r2 := route{
		group:    logic,
		endpoint: "/yyy",
		method:   "GET",
		handler:  handlerYYY,
	}

	s := NewServer(ctx, cfg, []route{r1, r2})
	s.engine.Run(s.socket)
}

func handlerYYY(c *gin.Context) {
	c.String(http.StatusOK, "yyy")
}
