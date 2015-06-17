package main

import (
	"encoding/json"
	"fmt"
	"github.com/abaril/GoLights/api"
	"log"
	"net/http"
	"time"
	"io"
)

type configOut struct {
	HueAddress  string `json:"hue_address,omitempty"`
	HueUsername string `json:"hue_username,omitempty"`
	HttpBindAddress string `json:"http_address,omitempty"`
	MqttBrokerAddress string `json:"mqtt_broker,omitempty"`
	AlarmTime   string `json:"alarm_time,omitempty"`
	AlarmLights []int  `json:"alarm_lights,omitempty"`
}

type configIn struct {
	HueAddress  *string `json:"hue_address"`
	HueUsername *string `json:"hue_username"`
	HttpBindAddress *string `json:"http_address"`
	MqttBrokerAddress *string `json:"mqtt_broker"`
	AlarmTime   string `json:"alarm_time"`
	AlarmLights *[]int  `json:"alarm_lights"`
}


func InitConfigAPI(dbVal api.MemDB) http.HandlerFunc {
	db = dbVal
	return serveConfig
}

func serveConfig(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		performGET(w, r)
		return

	case "PUT":
		performUpdate(w, r, false)
		return

	case "PATCH":
		performUpdate(w, r, true)
		return
	}

	http.Error(w, http.StatusText(404), 404)
}

func performGET(w http.ResponseWriter, r *http.Request) {
	c := configOut{}

	raw, err := db.Get("HueAddress")
	if err == nil {
		c.HueAddress = raw.(string)
	}
	raw, err = db.Get("HueUsername")
	if err == nil {
		c.HueUsername = raw.(string)
	}
	raw, err = db.Get("HttpBindAddress")
	if err == nil {
		c.HttpBindAddress = raw.(string)
	}
	raw, err = db.Get("MqttBrokerAddress")
	if err == nil {
		c.MqttBrokerAddress = raw.(string)
	}
	raw, err = db.Get("AlarmTime")
	if err == nil {
		time := raw.(time.Time)
		c.AlarmTime = fmt.Sprintf("%02d:%02d", time.Hour(), time.Minute())
	}
	raw, err = db.Get("AlarmLights")
	if err == nil {
		c.AlarmLights = raw.([]int)
	}

	json.NewEncoder(w).Encode(c)
}

func performUpdate(w http.ResponseWriter, r *http.Request, skipEmpty bool) {

	defer r.Body.Close()
	if err := updateFromReader(r.Body, skipEmpty); err != nil {
		log.Println("Err", err)
		http.Error(w, http.StatusText(406), 406)
	}

	w.Header().Set("Location", "api/v1/config")
	w.WriteHeader(http.StatusAccepted)
	performGET(w, r)
}

func updateFromReader(reader io.Reader, skipEmpty bool) error {
	var newConfig configIn
	if err := json.NewDecoder(reader).Decode(&newConfig); err != nil {
		return err
	}

	updateDatabase("HueAddress", newConfig.HueAddress, skipEmpty)
	updateDatabase("HueUsername", newConfig.HueUsername, skipEmpty)
	updateDatabase("HttpBindAddress", newConfig.HttpBindAddress, skipEmpty)
	updateDatabase("MqttBrokerAddress", newConfig.MqttBrokerAddress, skipEmpty)

	if timeVal, err := time.Parse("15:04", newConfig.AlarmTime); err == nil {
		updateDatabase("AlarmTime", timeVal, skipEmpty)
		UpdateNextAlarm(db)
	}
	updateDatabase("AlarmLights", newConfig.AlarmLights, skipEmpty)
	return nil
}

func updateDatabase(field string, value interface{}, skipEmpty bool) {

	switch t := value.(type) {
	case string, time.Time:
		db.Set(field, value)
	case *string:
		if t != nil {
			db.Set(field, *t)
		} else if !skipEmpty {
			db.Remove(field)
		}
	case *[]int:
		if t != nil {
			db.Set(field, *t)
		} else if !skipEmpty {
			db.Remove(field)
		}
	}
}