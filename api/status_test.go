package api_test

import (
	"github.com/abaril/GoLights/api"
	"github.com/cheekybits/is"
	"testing"
	"net/http/httptest"
	"net/http"
	"strings"
)

func TestServeHTTP(t *testing.T) {

	is := is.New(t)
	s := api.InitStatusAPI(api.UseMemDB)

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 200)
	is.Equal(strings.TrimSpace(w.Body.String()), "{\"is_alive\":true,\"is_home\":false}")

	r, _ = http.NewRequest("POST", "/", nil)
	w = httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 404)
}

