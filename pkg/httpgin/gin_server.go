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
	k8           = "/k8"
	logic        = "/logic"
	endpoint_xxx = "/xxx"
)

type Config struct {
	GraceSeconds uint8
	// berkely sockets are still 16 bit
	Port uint16
	L    *log.Logger
}

type route struct {
	group    string
	endpoint string
	method   string
	handler  gin.HandlerFunc
}

type GinServer struct {
	Config
	engine *gin.Engine
	chStop chan struct{}
}

func NewServer(cfg Config) *GinServer {
	s := new(GinServer)

	s.Config = cfg
	// checking if log level is debug
	if !(s.L.Level() == 1) {
		s.L.Debug("Setting Gin Log Level to Release Mode")
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine = gin.New()
	s.engine.RedirectTrailingSlash = true
	s.engine.HandleMethodNotAllowed = false
	s.chStop = make(chan struct{})

	return s
}

func (s *GinServer) registerRoute(r route) error {
	r.method = strings.ToTitle(r.method)

	s.L.Debug("Adding Route: ", r.group+r.endpoint, " method: ", r.method)

	switch r.method {
	case http.MethodGet:
		s.engine.GET(r.group+r.endpoint, r.handler)
	case http.MethodPost:
		s.engine.POST(r.group+r.endpoint, r.handler)
	case http.MethodPut:
		s.engine.PUT(r.group+r.endpoint, r.handler)
	case http.MethodPatch:
		s.engine.PATCH(r.group+r.endpoint, r.handler)
	case http.MethodDelete:
		s.engine.DELETE(r.group+r.endpoint, r.handler)
	default:
		return errors.New("unsupported method: " + r.method)
	}
	return nil
}

func (s *GinServer) registerRoutes(routes []route) error {
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

// prepareRoutes Method helps with route preparation.
// Routes need to contain the starting slash ex. /route.
func (s *GinServer) prepareRoutes() []route {
	r1 := route{
		group:    k8,
		endpoint: endpoint_xxx,
		method:   "GET",
		handler:  func(c *gin.Context) { c.String(http.StatusOK, "xxx") },
	}

	r2 := route{
		group:    logic,
		endpoint: "/yyy",
		method:   "GET",
		handler:  s.handlerYYY,
	}

	r3 := route{
		group:    k8,
		endpoint: "/shut",
		method:   "GET",
		handler:  s.handlerShutdown,
	}

	return []route{r1, r2, r3}
}

func (s *GinServer) Run(ctx context.Context) error {
	// Register routes
	if errPrep := s.registerRoutes(s.prepareRoutes()); errPrep != nil {
		return errors.Wrap(errPrep, "route preparation failed")
	}

	gracefulServer := &http.Server{
		Addr:    ":" + strconv.FormatUint(uint64(s.Config.Port), 10),
		Handler: s.engine,
	}

	// non blocking starting Gin using standard HTTP server graceful shutdown.
	go func() {
		if errServe := gracefulServer.ListenAndServe(); errServe != nil && errServe != http.ErrServerClosed {
			s.L.Fatalf("listen: %s\n", errServe)
		}
	}()

	<-s.chStop
	s.shutdown(gracefulServer)

	return nil
}

func (s *GinServer) shutdown(serverHTTP *http.Server) {
	s.L.Print("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.GraceSeconds)*time.Second)
	defer cancel()

	if errShutdown := serverHTTP.Shutdown(ctx); errShutdown != nil {
		s.L.Printf("Error HTTP server shutdown: %v", errShutdown)
	}

}

func (s *GinServer) handlerYYY(c *gin.Context) {
	c.String(http.StatusOK, "yyy")
}

func (s *GinServer) handlerShutdown(c *gin.Context) {
	c.String(http.StatusOK, "shutting down in ", strconv.FormatUint(uint64(s.Config.GraceSeconds), 10), "...")
	s.chStop <- struct{}{}
}
