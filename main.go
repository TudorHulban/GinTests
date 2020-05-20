package main

import (
	"context"

	log "github.com/labstack/gommon/log"
)

func createConfig() config {
	l := log.New("HTTP Server")
	l.SetLevel(log.DEBUG)
	l.Info("Log Level: ", l.Level())

	return config{
		graceSeconds: 5,
		port:         8001,
		l:            l,
	}
}

func main() {
	s := NewServer(createConfig())

	ctx := context.Background()
	s.Run(ctx)
}
