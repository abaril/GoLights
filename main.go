package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/heatxsink/go-hue/src/lights"
	"net/http"
	"time"
	"log"
	"net"
)

func main() {

	// TODO: move this data into config
	ll := lights.NewLights("192.168.1.105", "allanbaril")

	When(NewAlarmTrigger(10*time.Second), NewAlarmExpired(api.UseMemDB), NewAlarmHandler(ll, api.UseMemDB))
	When(DimLightsTrigger, nil, NewDimLightsAction(ll, api.UseMemDB))

	serverAddr, err := net.ResolveUDPAddr("udp",":8080")
	if err != nil {
		log.Fatalln("Unable to register UDP port")
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Fatalln("Unable to listen on UDP")
	}
	defer serverConn.Close()

	http.HandleFunc("/api/v1/status", InitStatusAPI(api.UseMemDB))
	http.HandleFunc("/api/v1/config", InitConfigAPI(api.UseMemDB))
	http.HandleFunc("/api/v1/dimlights", InitDimLightsAPI(api.UseMemDB))

	mqttHandleFunc("/dimlights", handleMQTTDimLights)
	mqttStart()

	http.ListenAndServe(":8080", http.DefaultServeMux)
}