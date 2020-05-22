package httpgin

import (
	log "github.com/labstack/gommon/log"
)

func createConfig() Config {
	l := log.New("HTTP Server")
	l.SetLevel(log.DEBUG)
	l.Info("Log Level: ", l.Level())

	return Config{
		GraceSeconds: 5,
		Port:         8001,
		L:            l,
	}
}
