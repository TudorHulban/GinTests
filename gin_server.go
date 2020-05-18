package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

const (
	k8    = "k8"
	logic = "logic"
)

type config struct {
	socket string
	l      *log.Logger
}

type route struct {
	group    string
	endpoint string
	method   string
	handler  gin.HandlerFunc
}

type Server struct {
	config
	engine *gin.Engine
}

func NewServer(cfg config) *Server {
	s := new(Server)

	s.config = cfg
	if !(s.l.Level() == 0) {
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine = gin.New()
	s.engine.RedirectTrailingSlash = true
	s.engine.HandleMethodNotAllowed = false

	return s
}

func (s *Server) registerRoute(r route) error {
	r.method = strings.ToTitle(r.method)

	var slash string
	if r.group != "" {
		slash = "/"
	}

	s.l.Debug("Adding Route: ", r.group+slash+r.endpoint, " method: ", r.method)

	switch r.method {
	case http.MethodGet:
		s.engine.GET(r.group+slash+r.endpoint, r.handler)
	case http.MethodPost:
		s.engine.POST(r.group+slash+r.endpoint, r.handler)
	case http.MethodPut:
		s.engine.PUT(r.group+slash+r.endpoint, r.handler)
	case http.MethodPatch:
		s.engine.PATCH(r.group+slash+r.endpoint, r.handler)
	case http.MethodDelete:
		s.engine.DELETE(r.group+slash+r.endpoint, r.handler)
	default:
		return errors.New("unsupported method: " + r.method)
	}
	return nil
}

func (s *Server) registerRoutes(routes []route) error {
	if len(routes) == 0 {
		return errors.New("no routes to add")
	}
	s.l.Debug("", routes)

	for _, route := range routes {
		if errReg := s.registerRoute(route); errReg != nil {
			return errReg
		}
	}
	return nil
}

// prepareRoutes Method helps with route preparation.
func (s *Server) prepareRoutes() []route {
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

	return []route{r1, r2}
}

func (s *Server) Run(ctx context.Context) error {
	// Register routes
	if errPrep := s.registerRoutes(s.prepareRoutes()); errPrep != nil {
		return errors.Wrap(errPrep, "route preparation failed")
	}

	s.l.Debug("Running Gin")
	s.engine.Run()
	s.l.Info("exit")
	return nil
}
