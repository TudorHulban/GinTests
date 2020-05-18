package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	cfg := config{
		socket:    "127.0.0.1:8001",
		debugMode: true,
	}

	s := NewServer(cfg)
	s.Run(ctx)
}

func handlerYYY(c *gin.Context) {
	c.String(http.StatusOK, "yyy")
}
