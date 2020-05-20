package main

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
)

func beforeEach() {

}

func TestGetUser_Success(t *testing.T) {

	apitest.New().
		Handler(newApp().Router).
		Get("/user/1234").
		Expect(t).
		Body(`{"id": "1234", "name": "Andy"}`).
		Status(http.StatusOK).
		End()
}
