package api_test

import (
	"github.com/abaril/GoLights/api"
	"github.com/cheekybits/is"
	"testing"
)

func TestGetStatus(t *testing.T) {

	is := is.New(t)
	status := api.UseMemoryStatus
	value, err := status.Get("IsAlive")
	if err != nil {
		t.Error("IsAlive should return a value")
	}
	is.Equal(value, true)

	value, err = status.Get("InvalidKey")
	if err == nil {
		t.Error("InvalidKey should result in an error")
	}
	is.Equal(err, api.ErrInvalidKey)
}

func TestSetStatus(t *testing.T) {

	is := is.New(t)
	status := api.UseMemoryStatus
	err := status.Set("IsHome", "worng type")
	if err == nil {
		t.Error("Should result in a wrong type error")
	}
	is.Equal(err, api.ErrInvalidType)

	err = status.Set("IsHome", true)
	if err != nil {
		t.Error("Set IsHome should work")
	}
	value, _ := status.Get("IsHome")
	is.Equal(value, true)

	err = status.Set("InvalidKey", "hello")
	if err == nil {
		t.Error("InvalidKey should result in an error")
	}
	is.Equal(err, api.ErrInvalidKey)
}
