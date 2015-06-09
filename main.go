package main

import (
	"github.com/abaril/GoLights/api"
	"net/http"
	"github.com/heatxsink/go-hue/src/lights"
	"time"
)

func main() {

	ll := lights.NewLights("192.168.1.105", "allanbaril")
	lightsArray := ll.GetAllLights()
	for _, light := range lightsArray {
		ll.SetLightState(light.Id, lights.State{On:false})
	}

//	bridge := hue.NewBridge("192.168.1.105", "allanbaril").Debug()
//	log.Println("Bridge:", bridge)
//	lights, err := bridge.GetAllLights()
//	if err != nil {
//		log.Println("Unable to retrieve lights:", err)
//	}
//	log.Println("Lights:", lights)
//	for _, light := range lights {
//		log.Println("Adjusting light:", light)
//		light.ColorLoop()
//	}

	action := LogAction{}
	go action.When(timeAdapter(time.Tick(1*time.Second)))

	http.HandleFunc("/api/v1/status", api.InitStatusAPI(api.UseMemDB))
	http.ListenAndServe(":8080", http.DefaultServeMux)

}

func timeAdapter(timeChannel <-chan time.Time) <-chan bool {

	boolChannel := make(chan bool)
	go func() {
		for range timeChannel {
			boolChannel <- true;
		}
	}()
	return boolChannel
}
