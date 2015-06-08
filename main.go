package main

import (
	"github.com/abaril/GoLights/api"
	"net/http"
)

func main() {

	http.Handle("/api/v1/status", api.UseMemoryStatus)
	http.ListenAndServe(":8080", http.DefaultServeMux)
}
