package main
import (
	"testing"
	"github.com/cheekybits/is"
	"github.com/abaril/GoLights/api"
	"net/http"
	"net/http/httptest"
	"strings"
)

func TestServeDimLights(t *testing.T) {

	is := is.New(t)
	s := InitDimLightsAPI(api.NewMemDB())

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 404)

	r, _ = http.NewRequest("POST", "/", strings.NewReader("invalid"))
	w = httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 406)

	r, _ = http.NewRequest("POST", "/", strings.NewReader("{\"level\": 1, \"lights\": [3]}"))
	w = httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 202)
	is.Equal(strings.TrimSpace(w.Body.String()), "{\"level\":1,\"lights\":[3]}")
}
