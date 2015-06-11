package main

import (
	"encoding/json"
	"fmt"
	"github.com/abaril/GoLights/api"
	"log"
	"net/http"
	"time"
)

type config struct {
	AlarmTime   string `json:"alarm_time"`
	AlarmLights []int  `json:"alarm_lights"`
}

//var db MemDB

func InitConfigAPI(dbVal api.MemDB) http.HandlerFunc {
	db = dbVal
	return serveConfig
}

func serveConfig(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		performGET(w, r)
		return

	case "POST":
		performPOST(w, r)
		return
	}

	http.Error(w, http.StatusText(404), 404)
}

func performGET(w http.ResponseWriter, r *http.Request) {
	c := config{}
	raw, err := db.Get("AlarmTime")
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

func performPOST(w http.ResponseWriter, r *http.Request) {
	var newConfig config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		log.Println("Err", err)
		http.Error(w, http.StatusText(406), 406)
		return
	}
	r.Body.Close()

	if timeVal, err := time.Parse("15:04", newConfig.AlarmTime); err == nil {
		db.Set("AlarmTime", timeVal)
		UpdateNextAlarm(db)
	}
	db.Set("AlarmLights", newConfig.AlarmLights)

	w.Header().Set("Location", "api/v1/config")
	w.WriteHeader(http.StatusAccepted)
	performGET(w, r)
}
