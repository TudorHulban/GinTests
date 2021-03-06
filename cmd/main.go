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
		Version:      theMainVersion,
	}
}

func main() {
	log.Print("version:", theMainVersion)

	s := httpgin.NewServer(createConfig())

	/*

		mwLogger := httpgin.Middleware{
			MiddleW: httpgin.MLogger,
			Cfg: httpgin.MConfig{
				Skipper: func() bool {
					return false
				},
			},
		}
		s.RegisterMiddleware(mwLogger)
	*/

	ctx := context.Background()
	s.Run(ctx)
}
