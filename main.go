package main

import (
	"context"

	log "github.com/labstack/gommon/log"
)

func main() {
	l := log.New("HTTP Server")
	l.SetLevel(log.DEBUG)

	l.Info("Log Level: ", l.Level())

	ctx := context.Background()
	cfg := config{
		graceSeconds: 5,
		port:         8001,
		l:            l,
	}

	s := NewServer(cfg)
	s.Run(ctx)
}
