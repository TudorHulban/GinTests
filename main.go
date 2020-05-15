package main

import (
	"context"
	"log"
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
	result.registerRoutes()

	return result
}

func (s *Server) registerRoutes() {
	for _, v := range s.routes {
		method := strings.ToTitle(v.method)

		var slash string
		if v.group != "" {
			slash = "/"
		}

		switch method {
		case http.MethodGet:
			s.engine.GET(v.group+slash+v.endpoint, v.handler)
		case http.MethodPost:
			s.engine.POST(v.group+slash+v.endpoint, v.handler)
		case http.MethodPut:
			s.engine.PUT(v.group+slash+v.endpoint, v.handler)
		case http.MethodPatch:
			s.engine.PATCH(v.group+slash+v.endpoint, v.handler)
		case http.MethodDelete:
			s.engine.DELETE(v.group+slash+v.endpoint, v.handler)
		default:
			log.Println("Unsupported method", method)
		}
	}
}

func main() {
	ctx := context.Background()
	cfg := config{
		socket: "127.0.0.1:8001",
	}

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

	s := NewServer(ctx, cfg, []route{r1, r2})
	s.engine.Run(s.socket)
}

func handlerYYY(c *gin.Context) {
	c.String(http.StatusOK, "yyy")
}
