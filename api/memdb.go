package api

import (
	"errors"
	"sync"
)

var ErrInvalidKey = errors.New("memdb: Invalid key provided")

type MemDB interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{})
}

var UseMemDB MemDB = &storage{
	data: make(map[string]interface{}),
	mutex: sync.RWMutex{},
}

type storage struct {
	data map[string]interface{}
	mutex   sync.RWMutex
}

func (s *storage) Get(key string) (interface{}, error) {

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if val, ok := s.data[key]; ok {
		return val, nil
	}
	return "", ErrInvalidKey
}

func (s *storage) Set(key string, value interface{}) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
}


