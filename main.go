package main

import (
	"flag"
	"github.com/abaril/GoLights/api"
	"github.com/abaril/go-hue/src/lights"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	VERSION_STRING = "1.0.4"
)

func main() {

	db := api.UseMemDB
	configHandler := retrieveBaseConfiguration(db)

	log.Println("Starting GoLights " + VERSION_STRING)

	var err error
	var hueAddress interface{}
	var hueUsername interface{}
	if hueAddress, err = db.Get("HueAddress"); err != nil {
		log.Fatalln("Unable to retrieve configuration")
	}
	if hueUsername, err = db.Get("HueUsername"); err != nil {
		log.Fatalln("Unable to retrieve configuration")
	}

	ll := lights.NewLights(hueAddress.(string), hueUsername.(string))

	mqttBroker := db.GetOrDefault("MqttBrokerAddress", "tcp://127.0.0.1:1883")
	mqttHandleFunc("/dimlights", handleMQTTDimLights)
	mqttHandleFunc("/devicestatus", NewHandleMQTTDeviceStatusReport(db));
	var mqttClient = mqttStart(mqttBroker.(string), "golights")

	When(NewAlarmTrigger(10*time.Second), NewAlarmExpired(db), NewAlarmHandler(ll, db))
	When(NewAlarmTrigger(10*time.Minute), NewCheckIfPollWeather(db), NewPerformWeatherPoll(db))
	When(NewUserTrigger(db), nil, NewAtHomeChangedHandler(ll, db))
	When(NewAlarmTrigger(1*time.Minute), nil, NewDeviceStatusPoll(db, mqttClient))

	httpBindAddress := db.GetOrDefault("HttpBindAddress", ":8080")
	http.HandleFunc("/api/v1/status", InitStatusAPI(db))
	http.HandleFunc("/api/v1/config", configHandler)
	http.HandleFunc("/api/v1/dimlights", InitDimLightsAPI(ll))
	http.ListenAndServe(httpBindAddress.(string), http.DefaultServeMux)
}

// TODO: create a config "object" and move this functionality there
func retrieveBaseConfiguration(db api.MemDB) http.HandlerFunc {

	configHandler := InitConfigAPI(db)

	env := os.Getenv("GOLIGHTS_CONFIG")
	if len(env) > 0 {
		if err := updateFromReader(strings.NewReader(env), true); err != nil {
			log.Println("Unable to parse configuration from environment variable GOLIGHTS_CONFIG")
		}
	}

	hueAddress := flag.String("ha", "", "Hue address")
	hueUsername := flag.String("hu", "", "Hue username")
	mqttBroker := flag.String("mqtt", "", "MQTT broker")
	flag.Parse()

	if len(*hueAddress) > 0 {
		db.Set("HueAddress", *hueAddress)
	}
	if len(*hueUsername) > 0 {
		db.Set("HueUsername", *hueUsername)
	}
	if len(*mqttBroker) > 0 {
		db.Set("MqttBrokerAddress", *mqttBroker)
	}

	return configHandler
}
