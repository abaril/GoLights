package api_test
import (
	"testing"
	"github.com/cheekybits/is"
	"github.com/abaril/GoLights/api"
)

func TestGetStatus(t *testing.T) {

	is := is.New(t)
	db := api.UseMemDB
	_, err := db.Get("IsAlive")
	if err == nil {
		t.Error("IsAlive should return err")
	}
	is.Equal(err, api.ErrInvalidKey)

	db.Set("IsAlive", true)
	val, err := db.Get("IsAlive")
	if err != nil {
		t.Error("IsAlive should return a value")
	}
	is.Equal(val, true)
}

func TestSetStatus(t *testing.T) {

	is := is.New(t)
	db := api.UseMemDB
	db.Set("IsHome", "value")
	val, err := db.Get("IsHome")
	if err != nil {
		t.Error("IsHome should return a value")
	}
	is.Equal(val, "value")

	db.Set("IsHome", true)
	val, err = db.Get("IsHome")
	if err != nil {
		t.Error("IsHome should return a value")
	}
	is.Equal(val, true)
}
