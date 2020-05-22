package main

import (
	"log"
	"net/http"
	"testing"

	"github.com/TudorHulban/GinTests"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestXXX(t *testing.T) {
	s := NewServer(createConfig())
	if assert.Nil(t, s.registerRoutes(s.prepareRoutes())) {
		checkedRoute := k8 + endpoint_xxx
		log.Println("Route: ", checkedRoute)

		apitest.New().
			Handler(s.engine).
			Get(checkedRoute).
			Expect(t).
			Body(`xxx`).
			Status(http.StatusOK).
			End()
	}
}
