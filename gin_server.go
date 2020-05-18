package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	k8    = "k8"
	logic = "logic"
)

type config struct {
	socket    string
	debugMode bool
}

type route struct {
	group    string
	endpoint string
	method   string
	handler  gin.HandlerFunc
}

type Server struct {
	config
	routes []route
	engine *gin.Engine
	// TODO: assert if need a logger
}

func NewServer(cfg config) *Server {
	s := new(Server)

	s.engine = gin.New()
	s.engine.RedirectTrailingSlash = true
	s.engine.HandleMethodNotAllowed = false
	s.config = cfg

	if s.config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	return s
}

func (s *Server) registerRoute(r route) error {
	r.method = strings.ToTitle(r.method)

	var slash string
	if r.group != "" {
		slash = "/"
	}

	log.Println("Route:", r.method, r.group+slash+r.endpoint)

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

	for _, route := range s.routes {
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

	log.Println("run")
	s.engine.Run()
	log.Println("exit")
	return nil
}
