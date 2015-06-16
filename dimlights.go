package main
import (
	"github.com/abaril/GoLights/api"
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

func InitDimLightsAPI(dbVal api.MemDB) http.HandlerFunc {
	db = dbVal
	return serveDimLights
}

func serveDimLights(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		var settings DimLights
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			log.Println("Err", err)
			http.Error(w, http.StatusText(406), 406)
			return
		}
		r.Body.Close()

		db.Set("DimLights", settings)

		w.Header().Set("Location", "api/v1/dimlights")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(settings)

		return
	}

	http.Error(w, http.StatusText(404), 404)
}

func handleMQTTDimLights(client *mqtt.Client, message mqtt.Message) {
	var settings DimLights
	if err := json.NewDecoder(bytes.NewReader(message.Payload())).Decode(&settings); err != nil {
		log.Println("Err", err)
		return
	}

	db.Set("DimLights", settings)
}

func DimLightsTrigger(events chan<- interface{}) {
	changed := db.Notify("DimLights")
	for range changed {
		events <- true
	}
}

func NewDimLightsAction(ll *lights.Lights, db api.MemDB) ActionFunc {
	return func() {
		raw, err := db.Get("DimLights")
		if err != nil {
			log.Println("Unable to retrieve dimLights settings:", err)
			return
		}
		settings := raw.(DimLights)
		for _, light := range settings.Lights {
			if settings.Level == 0 {
				ll.SetLightState(light, lights.State{On: false})
			} else {
				var brightness uint8 = uint8(255.0 * settings.Level/10.0)
				log.Println("Dimming light", light, "to", brightness)
				ll.SetLightState(light, lights.State{On: true, Bri: brightness})
			}
		}
	}

}



