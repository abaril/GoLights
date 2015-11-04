package main

import (
	"encoding/json"
	"github.com/abaril/GoLights/api"
	"net/http"
	"time"
	"math"
	"log"
)

type status struct {
	Version   string           `json:"version"`
	IsAlive   bool             `json:"is_alive"`
	IsHome    bool             `json:"is_home"`
	NextAlarm *time.Time       `json:"next_alarm,omitempty"`
	Weather   *WeatherForecast `json:"forecast,omitempty"`
	DeviceStatus *map[string]DeviceStatusReport `json:"device_status,omitempty"`
}

var db api.MemDB

func InitStatusAPI(dbVal api.MemDB) http.HandlerFunc {
	db = dbVal
	db.Set("IsHome", false)
	return serveStatus
}

func serveStatus(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		s := status{
			Version: VERSION_STRING,
			IsAlive: true,
		}
		raw, err := db.Get("IsHome")
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		s.IsHome = raw.(bool)

		raw, err = db.Get("NextAlarm")
		if err == nil {
			time := raw.(time.Time)
			s.NextAlarm = &time
		}

		raw, err = db.Get("Weather")
		if err == nil {
			weather := raw.(WeatherForecast)
			s.Weather = &weather
		}

		raw, err = db.Get("DeviceStatus")
		if err == nil {
			originalStatus := raw.(map[string]DeviceStatusReport)
			newStatus := make(map[string]DeviceStatusReport)
			for name, status := range originalStatus {
				delta := math.Abs(time.Since(status.Time.Time).Minutes())
				status.Alive = delta <= 5
				newStatus[name] = status
			}
			s.DeviceStatus = &newStatus
		}

		json.NewEncoder(w).Encode(s)
		return
	}

	http.Error(w, http.StatusText(404), 404)
}
