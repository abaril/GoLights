package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/cheekybits/is"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeHTTP(t *testing.T) {

	is := is.New(t)
	s := InitStatusAPI(api.NewMemDB())

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 200)
	is.Equal(strings.TrimSpace(w.Body.String()), "{\"version\":\"1.0.2\",\"is_alive\":true,\"is_home\":false}")

	r, _ = http.NewRequest("POST", "/", nil)
	w = httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 404)
}
