package main
import (
	"testing"
	"github.com/cheekybits/is"
	"net/http"
	"net/http/httptest"
	"strings"
)

func TestServeDimLights(t *testing.T) {

	is := is.New(t)
	s := InitDimLightsAPI(nil)

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

func TestRestoreLastLevel(t *testing.T) {

	is := is.New(t)
	InitDimLightsAPI(nil)

	var level float64
	level = restoreLastLevel(4, 5.0)
	is.Equal(level, 5.0)

	saveLastLevel(4, 3.0)
	level = restoreLastLevel(4, 5.0)
	is.Equal(level, 3.0)

	level = restoreLastLevel(4, 7.5)
	is.Equal(level, 7.5)
}
