package httpgin

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	s := NewServer(createConfig())

	if assert.Nil(t, s.RegisterRoutes(s.PrepareRoutes())) {
		w := httptest.NewRecorder()
		testURL := EndPointGroupK8 + EndpointVersion
		log.Println("TEST URL: ", testURL)

		req, _ := http.NewRequest(http.MethodGet, testURL, nil)
		s.Engine.ServeHTTP(w, req)

		log.Println("  ")
		for i, v := range s.Engine.Routes() {
			log.Println("Route: ", i, v)
		}
		assert.Equal(t, 200, w.Code)
	}
}
