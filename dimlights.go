package main
import (
	"net/http"
	"encoding/json"
	"log"
	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"bytes"
	"github.com/abaril/go-hue/src/lights"
	"strings"
)

const (
	ACTION_SET = iota
	ACTION_PUSH
	ACTION_POP
)

type DimLights struct {
	Level   float64 `json:"level"`
	Lights []int  `json:"lights"`
	State string `json:"state,omitempty"`
}

var lightsService *lights.Lights
var pastLevels map[int][]float64

func InitDimLightsAPI(ll *lights.Lights) http.HandlerFunc {
	lightsService = ll
	pastLevels = make(map[int][]float64)
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

func saveLastLevel(lightId int, level float64) {
	pastLevels[lightId] = append(pastLevels[lightId], level)
}

func restoreLastLevel(lightId int, defaultLevel float64) float64 {
	level := defaultLevel
	if pastLevels[lightId] != nil {
		stack := pastLevels[lightId]
		if len(stack) > 0 {
			level, pastLevels[lightId] = stack[len(stack)-1], stack[:len(stack)-1]
		}
	}
	return level
}

func dimLightsAction(ll *lights.Lights, request DimLights)  {
	action := ACTION_SET
	var allLights []lights.Light

	if len(request.State) > 0 {
		switch strings.ToLower(request.State) {
		case "push":
			allLights = ll.GetAllLights()
			action = ACTION_PUSH
		case "pop":
			action = ACTION_POP
		}
	}

	for _, light := range request.Lights {
		level := request.Level
		if action == ACTION_PUSH {
			for _, lightState := range allLights {
				if lightState.Id == light {
					currentLevel := 0.0
					if lightState.State.On {
						currentLevel = (float64(lightState.State.Bri) / 255.0) * 10.0
					}
					saveLastLevel(lightState.Id, currentLevel)
					break
				}
			}
		} else if action == ACTION_POP {
			level = restoreLastLevel(light, request.Level)
		}

		if level == 0 {
			if ll != nil {
				ll.SetLightState(light, lights.State{On: false})
			}
		} else {
			var brightness uint8 = uint8(255.0 * level/10.0)
			log.Println("Dimming light", light, "to", brightness)
			if ll != nil {
				ll.SetLightState(light, lights.State{On: true, Bri: brightness})
			}
		}
	}
}



