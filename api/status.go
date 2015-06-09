package api

import (
	"encoding/json"
	"net/http"
)

type status struct {
	IsAlive bool `json:"is_alive"`
	IsHome  bool `json:"is_home"`
}

var db MemDB

func InitStatusAPI(dbVal MemDB) http.HandlerFunc {
	db = dbVal
	db.Set("IsHome", false)
	return serveStatus
}

func serveStatus(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		s := status{
			IsAlive: true,
		}
		raw, err := db.Get("IsHome")
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		s.IsHome = raw.(bool)

		json.NewEncoder(w).Encode(s)
		return
	}

	http.Error(w, http.StatusText(404), 404)
}
