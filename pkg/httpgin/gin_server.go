package httpgin

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

const (
	// EndPointGroupK8 Is the endpoint group for k8 related endpoints.
	EndPointGroupK8 = "/k8"
	// EndPointGroupLogic Is the endpoint group for business logic related endpoints.
	EndPointGroupLogic = "/logic"
	// Endpointxxx Is test endpoint.
	Endpointxxx = "/xxx"
)

// Config Concentrates attributes for starting a Gin server.
type Config struct {
	GraceSeconds uint8
	// berkeley sockets are still 16 bit
	Port uint16
	// logger to use with Gin
	L *log.Logger
}

// Route Concentrates information to define a Gin route.
type Route struct {
	Group    string
	Endpoint string
	Method   string
	Handler  gin.HandlerFunc
}

// MConfig Is middleware config.
type MConfig struct {
	// Skipper Provides a way to skip middleware. If true skip the middleware.
	Skipper func() bool
	// MParams passes middleware specifc configuration. TODO: Refactor under a more concrete example.
	MParams interface{}
}

// Middleware Type defined for injecting middlewares.
type Middleware struct {
	// MConfig Is middleware config.
	Cfg     MConfig
	MiddleW func(MConfig) gin.HandlerFunc
}

// GinServer type has everything for starting a Gin server.
type GinServer struct {
	Config
	Engine *gin.Engine
	chStop chan struct{}
	// Server could be alive but not ready to take requests.
	isReady bool
}

// NewServer Is constructor for Gin server. Returns a pointer to the created instance.
func NewServer(cfg Config) *GinServer {
	s := new(GinServer)

	s.Config = cfg
	// checking if log level is debug
	if !(s.L.Level() == 1) {
		s.L.Debug("Setting Gin Log Level to Release Mode")
		gin.SetMode(gin.ReleaseMode)
	}

	// New does not insert any middleware.
	s.Engine = gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	s.Engine.Use(gin.Recovery())
	s.Engine.RedirectTrailingSlash = true
	s.Engine.HandleMethodNotAllowed = false
	s.chStop = make(chan struct{})

	return s
}

// RegisterMiddleware Public method to add middleware to the Gin server.
func (s *GinServer) RegisterMiddleware(m Middleware) {
	s.Engine.Use(m.MiddleW(m.Cfg))
}

func (s *GinServer) registerRoute(r Route) error {
	r.Method = strings.ToTitle(r.Method)

	s.L.Debug("Adding Route: ", r.Group+r.Endpoint, " method: ", r.Method)

	switch r.Method {
	case http.MethodGet:
		s.Engine.GET(r.Group+r.Endpoint, r.Handler)
	case http.MethodPost:
		s.Engine.POST(r.Group+r.Endpoint, r.Handler)
	case http.MethodPut:
		s.Engine.PUT(r.Group+r.Endpoint, r.Handler)
	case http.MethodPatch:
		s.Engine.PATCH(r.Group+r.Endpoint, r.Handler)
	case http.MethodDelete:
		s.Engine.DELETE(r.Group+r.Endpoint, r.Handler)
	default:
		return errors.New("unsupported method: " + r.Method)
	}
	return nil
}

// RegisterRoutes Registeres the routes passed.
func (s *GinServer) RegisterRoutes(routes []Route) error {
	if len(routes) == 0 {
		return errors.New("no routes to add")
	}
	s.L.Debug("Routes to add: ", routes)

	for _, route := range routes {
		if errReg := s.registerRoute(route); errReg != nil {
			return errReg
		}
	}
	return nil
}

// PrepareRoutes Method helps with route preparation.
// Routes need to contain the starting slash ex. /route.
func (s *GinServer) PrepareRoutes() []Route {
	r1 := Route{
		Group:    EndPointGroupK8,
		Endpoint: Endpointxxx,
		Method:   "GET",
		Handler:  func(c *gin.Context) { c.String(http.StatusOK, "xxx") },
	}

	r2 := Route{
		Group:    EndPointGroupLogic,
		Endpoint: "/echo/:echo",
		Method:   "GET",
		Handler:  s.handlerEcho,
	}

	r3 := Route{
		Group:    EndPointGroupK8,
		Endpoint: "/shut",
		Method:   "GET",
		Handler:  s.handlerShutdown,
	}

	return []Route{r1, r2, r3}
}

// Run Method used to start server.
func (s *GinServer) Run(ctx context.Context) error {
	// Register routes
	if errPrep := s.RegisterRoutes(s.PrepareRoutes()); errPrep != nil {
		return errors.Wrap(errPrep, "route preparation failed")
	}

	gracefulServer := &http.Server{
		Addr:    ":" + strconv.FormatUint(uint64(s.Config.Port), 10),
		Handler: s.Engine,
	}

	// non blocking starting Gin using standard HTTP server graceful shutdown.
	go func() {
		s.L.Print("Listening on: ", s.Config.Port)
		if errServe := gracefulServer.ListenAndServe(); errServe != nil && errServe != http.ErrServerClosed {
			s.L.Fatalf("listen: %s\n", errServe)
		}
	}()

	<-s.chStop
	s.shutdown(gracefulServer)

	return nil
}

// shutdown Method providing gracefull shutdown.
func (s *GinServer) shutdown(serverHTTP *http.Server) {
	s.L.Print("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.GraceSeconds)*time.Second)
	defer cancel()

	if errShutdown := serverHTTP.Shutdown(ctx); errShutdown != nil {
		s.L.Printf("Error HTTP server shutdown: %v", errShutdown)
	}
}

func (s *GinServer) handlerEcho(c *gin.Context) {
	echo := c.Params.ByName("echo")
	c.String(http.StatusOK, echo)
}

func (s *GinServer) handlerShutdown(c *gin.Context) {
	c.String(http.StatusOK, "shutting down in ", strconv.FormatUint(uint64(s.Config.GraceSeconds), 10), "...")
	s.chStop <- struct{}{}
}
