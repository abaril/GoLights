package main
import (
	"github.com/abaril/GoLights/api"
	"time"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
)

type WeatherSettings struct {
	Key string
	Lat float32
	Lon float32
}

type WeatherForecast struct {
	Time Timestamp `json:"time"`
	Summary string `json:"summary"`
	SunriseTime Timestamp `json:"sunriseTime"`
	SunsetTime Timestamp `json:"sunsetTime"`
	PrecipProbability float32 `json:"precipProbability"`
	TemperatureMin float32 `json:"temperatureMin"`
	TemperatureMax float32 `json:"temperatureMax"`
}

type weatherDaily struct {
	Summary string `json:"summary"`
	Icon string `json:"icon"`
	Data []WeatherForecast `json:"data"`
}

type weatherResponse struct {
	Latitude float32 `json:"latitude"`
	Daily weatherDaily `json:"daily"`
}

func NewCheckIfPollWeather(db api.MemDB) ConditionFunc {
	return func() bool {
		raw, err := db.Get("LastWeatherPoll")
		if err != nil {
			return true
		}
		lastPoll, ok := raw.(time.Time)
		if !ok {
			log.Println("Invalid last poll time")
			return true
		}

		// trigger once a day at or after 3am
		return time.Now().Day() != lastPoll.Day() && time.Now().Hour() >= 3
	}
}

func NewPerformWeatherPoll(db api.MemDB) ActionFunc {
	return func() {
		db.Set("LastWeatherPoll", time.Now())

		raw, err := db.Get("WeatherSettings")
		if err != nil || raw == nil {
			log.Println("Weather not configured")
			return
		}
		settings := raw.(*WeatherSettings)

		resp, err := http.Get(settings.determineRequestURL())
		if err != nil {
			log.Println("Unable to retrieve weather:", err)
			return
		}
		if resp.StatusCode != 200 {
			log.Println("Unable to retrieve weather. Status code =", resp.StatusCode)
			return
		}

		forecast := &weatherResponse{}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
			log.Println("Unable to parse weather response:", err)
			return
		}

		if len(forecast.Daily.Data) > 0 {
			db.Set("Weather", forecast.Daily.Data[0])
		}
	}
}

func (w *WeatherSettings)determineRequestURL() string {
	return fmt.Sprintf("https://api.forecast.io/forecast/%s/%f,%f", w.Key, w.Lat, w.Lon)
}
