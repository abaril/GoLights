package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/abaril/go-hue/src/lights"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"time"
)

const (
	MAX_PING_SILENCE      = 20 * time.Minute
	PING_RTT              = 2 * time.Second
	PING_POLLING_INTERVAL = 10 * time.Second
)

func NewUserTrigger(db api.MemDB) TriggerFunc {
	return func(events chan<- interface{}) {
		atHome := true
		db.Set("IsHome", atHome)
		lastResponse := time.Now()

		for {
			var ip string
			if raw, err := db.Get("DetectUserIP"); err != nil {
				log.Println("IP not set")
				time.Sleep(20 * time.Second)
			} else {
				ip = raw.(string)

				p := fastping.NewPinger()
				p.MaxRTT = PING_RTT
				p.Network("udp")
				ra, err := net.ResolveIPAddr("ip4:icmp", ip)
				if err != nil {
					log.Println("Unable to resolve IP:", ip, "err=", err)
					time.Sleep(20 * time.Second)
					continue
				}
				p.AddIPAddr(ra)

				p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
					lastResponse = time.Now()
				}
				p.OnIdle = func() {}

				err = p.Run()
				if err != nil {
					log.Println("Ping failed:", err)
				}
				if time.Now().Sub(lastResponse) >= MAX_PING_SILENCE {
					if atHome {
						atHome = false
						db.Set("IsHome", atHome)
						events <- true
					}
				} else {
					if !atHome {
						atHome = true
						db.Set("IsHome", atHome)
						events <- true
					}
				}

				time.Sleep(PING_POLLING_INTERVAL - PING_RTT)
			}
		}
	}
}

func NewAtHomeChangedHandler(ll *lights.Lights, db api.MemDB) func() {
	return func() {
		raw, err := db.Get("IsHome")
		if err != nil {
			log.Println("Unable to retrieve IsHome state", err)
			return
		}
		isHome := raw.(bool)
		if isHome {
			log.Println("Returned home")
			if !isCurrentlyDarkOutside(db, time.Now()) {
				return
			}

			raw, err = db.Get("LightsOnArrival")
			if err != nil {
				log.Println("No lights set for arrival", err)
				return
			}

			for _, light := range raw.([]int) {
				ll.SetLightState(light, lights.State{On: true})
			}
		} else {
			log.Println("No longer home")
			for _, light := range ll.GetAllLights() {
				ll.SetLightState(light.Id, lights.State{On: false})
			}
		}
	}
}

func isCurrentlyDarkOutside(db api.MemDB, now time.Time) bool {
	if raw, err := db.Get("Weather"); err == nil {
		weather := raw.(WeatherForecast)
		return now.Before(weather.SunriseTime.Time) || now.After(weather.SunsetTime.Time)
	}
	return true
}
