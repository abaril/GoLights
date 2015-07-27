package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/cheekybits/is"
	"testing"
	"time"
)

func TestNewCheckIfPollWeather(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()
	f := NewCheckIfPollWeather(db)

	is.Equal(f(), true)

	db.Set("LastWeatherPoll", time.Now().Add(-24*time.Hour))
	is.Equal(f(), true)

	db.Set("LastWeatherPoll", time.Now())
	is.Equal(f(), false)
}

func TestNewPerformWeatherPoll(t *testing.T) {

	db := api.NewMemDB()
	f := NewPerformWeatherPoll(db)

	validKey := false
	db.Set("WeatherSettings", &WeatherSettings{
		Key: "<key needs to be provided to test>",
		Lat: 43.6314075,
		Lon: -79.3941305,
	})
	f()
	_, err := db.Get("Weather")
	if err != nil && validKey {
		t.Error("Weather fetch failed", err)
	}
}
