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

	db.Set("IsHome", true)
	db.Set("NextAlarm", time.Now().Add(-1*time.Minute))
	is.Equal(f(), true)

	db.Set("IsHome", false)
	db.Set("NextAlarm", time.Now().Add(-1*time.Minute))
	is.Equal(f(), false)
}

func TestUpdateNextAlarm(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()

	now := time.Date(2015, 6, 16, 8, 0, 0, 0, time.Local)

	alarm := now.Add(-1 * time.Hour)
	db.Set("AlarmTime", alarm)
	UpdateNextAlarm(db, now)
	raw, err := db.Get("NextAlarm")
	if err != nil {
		t.Fatal("Unable to retrieve nextAlarm")
	}
	is.Equal(raw.(time.Time), alarm.AddDate(0, 0, 1))

	alarm = now.Add(1 * time.Hour)
	db.Set("AlarmTime", alarm)
	UpdateNextAlarm(db, now)
	raw, err = db.Get("NextAlarm")
	if err != nil {
		t.Fatal("Unable to retrieve nextAlarm")
	}
	is.Equal(raw.(time.Time), alarm)

	now = time.Date(2015, 6, 19, 8, 0, 0, 0, time.Local)
	alarm = now.Add(-1 * time.Hour)
	db.Set("AlarmTime", alarm)
	UpdateNextAlarm(db, now)
	raw, err = db.Get("NextAlarm")
	if err != nil {
		t.Fatal("Unable to retrieve nextAlarm")
	}
	is.Equal(raw.(time.Time), alarm.AddDate(0, 0, 3))
}

func TestDetermineLightColour(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()

	w := WeatherForecast{
		PrecipProbability: 0.5,
	}
	db.Set("Weather", w)
	hue, sat := determineLightColour(db)
	is.Equal(hue, BLUE_H)
	is.Equal(sat, BLUE_S)

	w = WeatherForecast{
		PrecipProbability: 0.0,
		TemperatureMin: -20,
		TemperatureMax: -10,
	}
	db.Set("Weather", w)
	hue, sat = determineLightColour(db)
	is.Equal(hue, WHITE_H)
	is.Equal(sat, WHITE_S)

	w = WeatherForecast{
		PrecipProbability: 0.0,
		TemperatureMin: 10,
		TemperatureMax: 30,
	}
	db.Set("Weather", w)
	hue, sat = determineLightColour(db)
	is.Equal(hue, RED_H)
	is.Equal(sat, RED_S)

	w = WeatherForecast{
		PrecipProbability: 0.0,
		TemperatureMin: 10,
		TemperatureMax: 20,
	}
	db.Set("Weather", w)
	hue, sat = determineLightColour(db)
	is.Equal(hue, YELLOW_H)
	is.Equal(sat, YELLOW_S)

}
