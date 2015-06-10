package main

import (
	"github.com/abaril/GoLights/api"
	"net/http"
	"log"
	"time"
)

func main() {

//	ll := lights.NewLights("192.168.1.105", "allanbaril")
//	lightsArray := ll.GetAllLights()
//	for _, light := range lightsArray {
//		ll.SetLightState(light.Id, lights.State{On:false})
//	}

	condition := &Toggler{}
	When(timeTrigger, condition.toggle, firstAction(nextAction(nextAction(nil))))

	http.HandleFunc("/api/v1/status", api.InitStatusAPI(api.UseMemDB))
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func timeTrigger(events chan<- interface{}) {
	tc := time.Tick(1 * time.Second)
	for range tc {
		events <- true
	}
}

type Toggler struct {
	last bool
}

func (t *Toggler) toggle() bool {
	t.last = !t.last
	return t.last
}

func firstAction(a ActionFunc) ActionFunc {
	return func() {
		log.Println("Hello there")
		if a != nil {
			a()
		}
	}
}

func nextAction(a ActionFunc) ActionFunc {
	return func() {
		log.Println("  ... and there")
		if a != nil {
			a()
		}
	}
}
