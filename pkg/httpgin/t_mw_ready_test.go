package httpgin

import (
	"log"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

// MW Ready moves to ready state after one request that is not available.
// This is the behavior tested.
func TestMWReadySwitch(t *testing.T) {
	s := NewServer(createConfig())

	mwReady := Middleware{
		MiddleW: MReady,
		Cfg:     MConfig{},
	}

	s.RegisterMiddleware(mwReady)

	if assert.Nil(t, s.RegisterRoutes(s.PrepareRoutes())) {
		checkedRoute := EndPointGroupK8 + Endpointxxx
		log.Println("Route: ", checkedRoute)

		apitest.New().
			Handler(s.Engine).
			Get(checkedRoute).
			Expect(t).
			Status(http.StatusSeeOther).
			End()

		isReady = func() bool { return true }

		// next request should reach desired route handler now that server is ready.
		apitest.New().
			Handler(s.Engine).
			Get(checkedRoute).
			Expect(t).
			Body(`xxx`).
			Status(http.StatusOK).
			End()
	}
}
