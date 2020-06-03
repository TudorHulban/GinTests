package httpgin

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OSVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Env []OSVar

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
	c.JSON(http.StatusOK, gin.H{"staus": echo})
}

func (s *GinServer) handlerShutdown(c *gin.Context) {
	c.String(http.StatusOK, "shutting down in ", strconv.FormatUint(uint64(s.Config.GraceSeconds), 10), "...")
	s.chStop <- struct{}{}
}

func (s *GinServer) handlerServiceNotOperational(c *gin.Context) {
	s.L.Debug("endpoint service not operational")
	c.JSON(http.StatusServiceUnavailable, gin.H{"status": "temporary no service"})
}

// environmentHandler Method is common handler that fetched environment variables.
func (s *GinServer) environmentHandler(c *gin.Context) {
	buf4Sort := []string{}
	for _, osvar := range os.Environ() {
		buf4Sort = append(buf4Sort, osvar+"\n")
	}

	// sort list now for humans.
	sort.Strings(buf4Sort)

	// transform to bytes for sending over HTTP
	buf := &bytes.Buffer{}
	for _, sortedosvar := range buf4Sort {
		buf.WriteString(sortedosvar)
	}

	/*
		var env Env
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			buf.WriteString(e + "\n")
			env = append(env, OSVar{
				Name:  pair[0],
				Value: pair[1],
			})
		}
	*/
	//c.JSON(http.StatusOK, env)
	c.String(http.StatusOK, buf.String())
}
