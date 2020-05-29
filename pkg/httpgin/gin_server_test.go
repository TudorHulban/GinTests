package httpgin

import (
	"log"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestInternalRoutes(t *testing.T) {
	s := NewServer(createConfig())

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
