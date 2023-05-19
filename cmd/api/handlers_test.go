package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPets(t *testing.T) {
	req, _ := http.NewRequest("GET", "/pets", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testPets.GetPets)

	handler.ServeHTTP(rr, req)

}
