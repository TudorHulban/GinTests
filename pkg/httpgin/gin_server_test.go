package httpgin

import (
	"log"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestXXX(t *testing.T) {
	s := NewServer(createConfig())

	if assert.Nil(t, s.RegisterRoutes(s.PrepareRoutes())) {
		checkedRoute := EndPointGroupK8 + Endpointxxx
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
