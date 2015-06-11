package api

import (
	"errors"
	"reflect"
	"sync"
)

var ErrInvalidKey = errors.New("memdb: Invalid key provided")

type MemDB interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{})
	Notify(key string) chan bool
}

func NewMemDB() *storage {
	return &storage{
		data:          make(map[string]interface{}),
		notifications: make(map[string]chan bool),
		mutex:         sync.RWMutex{},
	}
}

var UseMemDB MemDB = &storage{
	data:          make(map[string]interface{}),
	notifications: make(map[string]chan bool),
	mutex:         sync.RWMutex{},
}

type storage struct {
	data          map[string]interface{}
	notifications map[string]chan bool

	mutex sync.RWMutex
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

	if reflect.TypeOf(value).Comparable() && s.data[key] == value {
		return
	}
	s.data[key] = value
	if notify, ok := s.notifications[key]; ok {
		notify <- true
	}
}

func (s *storage) Notify(key string) chan bool {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.notifications[key] = make(chan bool)
	return s.notifications[key]
}
