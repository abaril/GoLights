package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/heatxsink/go-hue/src/lights"
	"net/http"
	"time"
)

func main() {

	// TODO: move this data into config
	ll := lights.NewLights("192.168.1.105", "allanbaril")

	When(NewAlarmTrigger(10*time.Second), NewAlarmExpired(api.UseMemDB), NewAlarmHandler(ll, api.UseMemDB))

	http.HandleFunc("/api/v1/status", InitStatusAPI(api.UseMemDB))
	http.HandleFunc("/api/v1/config", InitConfigAPI(api.UseMemDB))

	http.ListenAndServe(":8080", http.DefaultServeMux)
}
