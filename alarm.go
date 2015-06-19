package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/heatxsink/go-hue/src/lights"
	"log"
	"time"
)

const (
	RED_H = uint16(5482)
	RED_S = uint8(192)
	WHITE_H = uint16(35831)
	WHITE_S = uint8(254)
	BLUE_H = uint16(44538)
	BLUE_S = uint8(252)
	YELLOW_H = uint16(23806)
	YELLOW_S = uint8(208)
)

func NewAlarmTrigger(interval time.Duration) TriggerFunc {
	return func(events chan<- interface{}) {
		tc := time.Tick(interval)
		for range tc {
			events <- true
		}
	}
}

func NewAlarmExpired(db api.MemDB) ConditionFunc {
	return func() bool {
		raw, err := db.Get("NextAlarm")
		if err != nil {
			// likely the alarm hasn't been configured
			return false
		}
		alarm, ok := raw.(time.Time)
		if !ok {
			log.Fatalln("Alarm time misconfigured")
			return false
		}

		return time.Now().After(alarm)
	}
}

func NewAlarmHandler(ll *lights.Lights, db api.MemDB) ActionFunc {
	return func() {
		log.Println("Alarm triggered!")
		UpdateNextAlarm(db, time.Now())

		raw, err := db.Get("AlarmLights")
		if err != nil {
			log.Println("Unable to retrieve alarmLights", err)
			return
		}
		alarmLights := raw.([]int)
		for _, light := range alarmLights {
			log.Println("Turning on light", light)
			hue, sat := determineLightColour(db)
			ll.SetLightState(light, lights.State{
				On: true,
				Bri: 254,
				Hue: hue,
				Sat: sat,
				TransitionTime:100,
			})
		}
	}
}

func determineLightColour(db api.MemDB) (uint16, uint8) {
	raw, err := db.Get("Weather")
	if err == nil {
		w := raw.(WeatherForecast)

		if w.PrecipProbability > 0.2 {
			return BLUE_H, BLUE_S
		}
		if w.TemperatureMin < -5.0 {
			return WHITE_H, WHITE_S
		}
		if w.TemperatureMax > 25.0 {
			return RED_H, RED_S
		}
	}
	return YELLOW_H, YELLOW_S
}

func UpdateNextAlarm(db api.MemDB, now time.Time) {
	raw, err := db.Get("AlarmTime")
	if err != nil {
		log.Println("Unable to retrieve alarmTime. Ensure configuration is correct")
		return
	}
	alarm, ok := raw.(time.Time)
	if !ok {
		log.Fatalln("Invalid config value")
	}

	nextAlarm := time.Date(now.Year(), now.Month(), now.Day(), alarm.Hour(), alarm.Minute(), 0, 0, now.Location())
	if nextAlarm.Before(now) {
		nextAlarm = nextAlarm.AddDate(0, 0, 1)
		for nextAlarm.Weekday() == time.Saturday || nextAlarm.Weekday() == time.Sunday {
			nextAlarm = nextAlarm.AddDate(0, 0, 1)
		}
	}
	db.Set("NextAlarm", nextAlarm)
}
