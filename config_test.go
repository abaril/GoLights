package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/cheekybits/is"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeConfigHTTP(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()
	s := InitConfigAPI(db)

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 200)
	is.Equal(strings.TrimSpace(w.Body.String()), "{}")

	r, _ = http.NewRequest("PUT", "/", strings.NewReader("{\"hue_username\": \"username\", \"alarm_time\": \"10:00\", \"alarm_lights\": [3,5]}"))
	w = httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 202)
	raw, err := db.Get("AlarmLights")
	is.Equal(err, nil)
	is.Equal(raw, []int{3, 5})
	raw, err = db.Get("HueUsername")
	is.Equal(err, nil)
	is.Equal(raw, "username")

	r, _ = http.NewRequest("PATCH", "/", strings.NewReader("{\"hue_username\": \"username2\", \"mqtt_broker\": \"tcp://127.0.0.1:1883\"}"))
	w = httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 202)
	raw, err = db.Get("AlarmLights")
	is.Equal(err, nil)
	is.Equal(raw, []int{3, 5})
	raw, err = db.Get("HueUsername")
	is.Equal(err, nil)
	is.Equal(raw, "username2")
	raw, err = db.Get("MqttBrokerAddress")
	is.Equal(err, nil)
	is.Equal(raw, "tcp://127.0.0.1:1883")

	r, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	s(w, r)
	is.Equal(w.Code, 200)
	is.Equal(strings.TrimSpace(w.Body.String()), "{\"hue_username\":\"username2\",\"mqtt_broker\":\"tcp://127.0.0.1:1883\",\"alarm_time\":\"10:00\",\"alarm_lights\":[3,5]}")

}
