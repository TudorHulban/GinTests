package httpgin

import (
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestInternalRoutes(t *testing.T) {
	s := NewServer(createConfig())

	// middleware needs to be added prior to adding routes.
	s.RegisterMiddleware(Middleware{
		MiddleW: MLogger,
	})

	if assert.Nil(t, s.RegisterRoutes(s.PrepareRoutes())) {

		apitest.New().
			Handler(s.Engine).
			Get(EndPointGroupK8 + Endpointxxx).
			Expect(t).
			Body(`xxx`).
			Status(http.StatusOK).
			End()

		log.Println(isReady())

		apitest.New().
			Handler(s.Engine).
			Get(EndPointGroupK8 + EndPointNoService).
			Expect(t).
			Status(http.StatusServiceUnavailable).
			End()

		apitest.New("Test Echo Handler").
			Handler(s.Engine).
			Get(EndPointGroupK8 + "/echo/1234").
			Expect(t).
			Status(http.StatusOK).
			End()
	}
}

func TestRegisteredRoutes(t *testing.T) {
	s := NewServer(createConfig())
	s.RegisterRoutes(s.PrepareRoutes())

	testCases := []struct {
		endpoint string
	}{
		{"x1"},
		{"x"},
	}

	for _, tcase := range testCases {
		t.Run(tcase.endpoint, func(t *testing.T) {
			isContained := false

			for _, v := range s.Engine.Routes() {
				log.Println(v.Path)
				if strings.Contains(v.Path, tcase.endpoint) {
					isContained = true
					break
				}
			}

			assert.True(t, isContained)
		})
	}
}
