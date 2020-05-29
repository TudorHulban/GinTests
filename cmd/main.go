package main

import (
	"context"

	"github.com/TudorHulban/GinTests/pkg/httpgin"
	log "github.com/labstack/gommon/log"
)

func createConfig() httpgin.Config {
	l := log.New("HTTP Server")
	l.SetLevel(log.DEBUG)
	l.Info("Log Level: ", l.Level())

	return httpgin.Config{
		GraceSeconds: 5,
		Port:         8001,
		L:            l,
	}
}

func main() {
	s := httpgin.NewServer(createConfig())
	mwLogger := httpgin.Middleware{
		MiddleW: httpgin.MLogger,
		Cfg: httpgin.MConfig{
			Skipper: func() bool {
				return true
			},
		},
	}
	s.RegisterMiddleware(mwLogger)

	mwReady := httpgin.Middleware{
		MiddleW: httpgin.MReady,
	}
	s.RegisterMiddleware(mwReady)

	ctx := context.Background()
	s.Run(ctx)
}
