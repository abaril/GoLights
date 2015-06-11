package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/cheekybits/is"
	"testing"
	"time"
)

func TestAlarmTrigger(t *testing.T) {

	is := is.New(t)
	count := 0
	f := NewAlarmTrigger(1 * time.Second)
	events := make(chan interface{})

	go func() {
		for range events {
			count += 1
		}
	}()
	go f(events)

	time.Sleep(2*time.Second + 200*time.Millisecond)
	close(events)

	is.Equal(count, 2)
}

func TestAlarmExpired(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()
	f := NewAlarmExpired(db)

	is.Equal(f(), false)

	db.Set("NextAlarm", time.Now().Add(+1*time.Minute))
	is.Equal(f(), false)

	db.Set("NextAlarm", time.Now().Add(-1*time.Minute))
	is.Equal(f(), true)
}

func TestUpdateNextAlarm(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())

	alarm := now.Add(-1 * time.Hour)
	db.Set("AlarmTime", alarm)
	UpdateNextAlarm(db)
	raw, err := db.Get("NextAlarm")
	if err != nil {
		t.Fatal("Unable to retrieve nextAlarm")
	}
	is.Equal(raw.(time.Time), alarm.Add(24*time.Hour))

	alarm = now.Add(1 * time.Hour)
	db.Set("AlarmTime", alarm)
	UpdateNextAlarm(db)
	raw, err = db.Get("NextAlarm")
	if err != nil {
		t.Fatal("Unable to retrieve nextAlarm")
	}
	is.Equal(raw.(time.Time), alarm)
}
