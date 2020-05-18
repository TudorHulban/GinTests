package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/labstack/gommon/log"
)

func main() {
	l := log.New("HTTP Server")
	l.SetLevel(log.DEBUG)

	ctx := context.Background()
	cfg := config{
		socket: "127.0.0.1:8001",
		l:      l,
	}

	s := NewServer(cfg)
	s.Run(ctx)
}

func handlerYYY(c *gin.Context) {
	c.String(http.StatusOK, "yyy")
}
