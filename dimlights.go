package main
import (
	"net/http"
	"encoding/json"
	"log"
	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"bytes"
	"github.com/heatxsink/go-hue/src/lights"
)

type DimLights struct {
	Level   int `json:"level"`
	Lights []int  `json:"lights"`
}

var lightsService *lights.Lights

func InitDimLightsAPI(ll *lights.Lights) http.HandlerFunc {
	lightsService = ll
	return serveDimLights
}

func serveDimLights(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		var request DimLights
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Println("Err", err)
			http.Error(w, http.StatusText(406), 406)
			return
		}
		r.Body.Close()

		dimLightsAction(lightsService, request)

		w.Header().Set("Location", "api/v1/dimlights")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(request)

		return
	}

	http.Error(w, http.StatusText(404), 404)
}

func handleMQTTDimLights(client *mqtt.Client, message mqtt.Message) {
	var request DimLights
	if err := json.NewDecoder(bytes.NewReader(message.Payload())).Decode(&request); err != nil {
		log.Println("Err", err)
		return
	}

	dimLightsAction(lightsService, request)
}

func dimLightsAction(ll *lights.Lights, request DimLights)  {
	for _, light := range request.Lights {
		if request.Level == 0 {
			if ll != nil {
				ll.SetLightState(light, lights.State{On: false})
			}
		} else {
			var brightness uint8 = uint8(255.0 * request.Level/10.0)
			log.Println("Dimming light", light, "to", brightness)
			if ll != nil {
				ll.SetLightState(light, lights.State{On: true, Bri: brightness})
			}
		}
	}
}



