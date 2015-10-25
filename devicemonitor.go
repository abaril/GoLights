package main

import (
	"log"
	"bytes"
	"encoding/json"
	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/abaril/GoLights/api"
	"time"
)

type DeviceStatusReport struct {
	Name  string `json:"name"`
	Timestamp int32   `json:"time"`
}

func handleMQTTDeviceStatusReport(client *mqtt.Client, message mqtt.Message) {
	var status DeviceStatusReport
	if err := json.NewDecoder(bytes.NewReader(message.Payload())).Decode(&status); err != nil {
		log.Println("Err", err)
		return
	}

	log.Println("Received ", status);
}

func NewDeviceStatusPoll(db api.MemDB, client *mqtt.Client) ActionFunc {
	return func() {
		db.Set("LastDeviceStatusPoll", time.Now())

		client.Publish("/devicestatus", 0, false, "{'request': 'hello'}");

//		raw, err := db.Get("WeatherSettings")
//		if err != nil || raw == nil {
//			log.Println("Weather not configured")
//			return
//		}
//		settings := raw.(*WeatherSettings)
//
//		resp, err := http.Get(settings.determineRequestURL())
//		if err != nil {
//			log.Println("Unable to retrieve weather:", err)
//			return
//		}
//		if resp.StatusCode != 200 {
//			log.Println("Unable to retrieve weather. Status code =", resp.StatusCode)
//			return
//		}
//
//		forecast := &weatherResponse{}
//		defer resp.Body.Close()
//		if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
//			log.Println("Unable to parse weather response:", err)
//			return
//		}
//
//		if len(forecast.Daily.Data) > 0 {
//			forecast.Daily.Data[0].TemperatureMin = fahrenheitToCelcius(forecast.Daily.Data[0].TemperatureMin)
//			forecast.Daily.Data[0].TemperatureMax = fahrenheitToCelcius(forecast.Daily.Data[0].TemperatureMax)
//			db.Set("Weather", forecast.Daily.Data[0])
//		}
	}
}


