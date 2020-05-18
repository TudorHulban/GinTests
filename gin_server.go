package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	k8    = "k8"
	logic = "logic"
)

type config struct {
	socket string
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
}

func NewServer(ctx context.Context, cfg config) *Server {
	s := new(Server)

	s.engine = gin.New()
	s.engine.RedirectTrailingSlash = true
	s.engine.HandleMethodNotAllowed = false
	s.config = cfg

	return s
}

func (s *Server) registerRoute(r route) error {
	r.method = strings.ToTitle(r.method)

	var slash string
	if r.group != "" {
		slash = "/"
	}

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

func (s *Server) registerRoutes(routes []route) {
	for _, route := range s.routes {
		s.registerRoute(route)
	}
}
