package api_test

import (
	"github.com/abaril/GoLights/api"
	"github.com/cheekybits/is"
	"testing"
	"time"
)

func TestGet(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()
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

func TestSet(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()
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

func TestNotify(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()
	notifyChan := db.Notify("IsHome")
	notifyCount := 0
	go func() {
		for range notifyChan {
			notifyCount += 1
		}
	}()

	is.Equal(notifyCount, 0)

	db.Set("IsHome", true)
	time.Sleep(10 * time.Millisecond)
	is.Equal(notifyCount, 1)

	db.Set("IsHome", true)
	time.Sleep(10 * time.Millisecond)
	is.Equal(notifyCount, 1)

	db.Set("IsHome", false)
	time.Sleep(10 * time.Millisecond)
	is.Equal(notifyCount, 2)

}
