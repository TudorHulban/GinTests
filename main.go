package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type config struct {
	socket string
}

type route struct {
	endpoint string
	method   string
	handler  gin.HandlerFunc
}

type Server struct {
	engine *gin.Engine
	routes []route
	config
}

func NewServer(ctx context.Context, cfg config, routes []route) *Server {
	result := new(Server)

	result.engine = gin.New()
	result.engine.RedirectTrailingSlash = true
	result.engine.HandleMethodNotAllowed = false
	result.config = cfg
	result.routes = routes

	return result
}

func (s *Server) registerRoutes() {
	for _, v := range s.routes {
		method := strings.ToTitle(v.method)

		switch method {
		case http.MethodGet:
			s.engine.GET(v.endpoint, v.handler)
		case http.MethodPost:
			r.cfg.Engine.POST(path, handler)
		case http.MethodPut:
			r.cfg.Engine.PUT(path, handler)
		case http.MethodPatch:
			r.cfg.Engine.PATCH(path, handler)
		case http.MethodDelete:
			r.cfg.Engine.DELETE(path, handler)
		default:
			r.cfg.Logger.Error("Unsupported method", method)
		}
	}
}

func main() {
	ctx := context.Background()
	cfg := config{
		socket: "127.0.0.1:8001",
	}

	r1 := route{
		endpoint: "/",
		method:   "GET",
		handler:  func(c *gin.Context) { c.String(http.StatusOK, "xxx") },
	}

	s := NewServer(ctx, cfg, []route{r1})
	s.engine.Run(s.socket)
}
