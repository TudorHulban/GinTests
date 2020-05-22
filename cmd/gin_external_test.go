package main

import (
	"log"
	"net/http"
	"testing"

	"github.com/TudorHulban/GinTests/pkg/httpgin"
	"github.com/gin-gonic/gin"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestXXX(t *testing.T) {
	s := httpgin.NewServer(createConfig())

	if assert.Nil(t, s.RegisterRoutes(s.PrepareRoutes())) {
		checkedRoute := httpgin.EndPointGroupK8 + httpgin.Endpointxxx
		log.Println("Route: ", checkedRoute)

		apitest.New().
			Handler(s.Engine).
			Get(checkedRoute).
			Expect(t).
			Body(`xxx`).
			Status(http.StatusOK).
			End()
	}
}

func TestEcho(t *testing.T) {
	s := httpgin.NewServer(createConfig())
	testRoute := httpgin.Route{Endpoint: "/echo/:echo", Method: "GET", Handler: echoRoute}

	if assert.Nil(t, s.RegisterRoutes([]httpgin.Route{testRoute})) {
		log.Println("Route: ", testRoute.Group+testRoute.Endpoint)

		apitest.New().
			Handler(s.Engine).
			Get("/echo/xxx").
			Expect(t).
			Body(`xxx`).
			Status(http.StatusOK).
			End()
	}
}

func echoRoute(c *gin.Context) {
	echo := c.Params.ByName("echo")
	c.String(http.StatusOK, echo)
}
