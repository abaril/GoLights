package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/heatxsink/go-hue/src/lights"
	"log"
	"time"
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
			ll.SetLightState(light, lights.State{On: true, Bri: 200})
		}
	}
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
