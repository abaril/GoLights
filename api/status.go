package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

var ErrInvalidKey = errors.New("status: Invalid key provided")
var ErrInvalidType = errors.New("status: Invalid type provided for value")

type Status interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type memoryStatus struct {
	IsAlive bool `json:"is_alive"`
	IsHome  bool `json:"is_home"`
	mutex   sync.RWMutex
}

var UseMemoryStatus *memoryStatus = &memoryStatus{
	IsAlive: true,
	IsHome:  false,
	mutex:   sync.RWMutex{},
}

func (s *memoryStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(s)
		return
	}

	http.Error(w, http.StatusText(404), 404)
}

func (s *memoryStatus) Get(key string) (interface{}, error) {

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	switch key {
	case "IsAlive":
		return s.IsAlive, nil
	case "IsHome":
		return s.IsHome, nil
	}
	return "", ErrInvalidKey
}

func (s *memoryStatus) Set(key string, value interface{}) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch key {
	case "IsHome":
		if cast, ok := value.(bool); ok {
			s.IsHome = cast
			return nil
		}
		return ErrInvalidType
	}
	return ErrInvalidKey
}
