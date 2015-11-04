package main

import (
	"log"
	"bytes"
	"encoding/json"
	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/abaril/GoLights/api"
	"time"
)

type DeviceStatusReport struct {
	Name  string `json:"name"`
	Timestamp int32   `json:"time"`
	Alive bool `json:"alive"`
}

func NewHandleMQTTDeviceStatusReport(db api.MemDB) mqtt.MessageHandler {
	return func(client *mqtt.Client, message mqtt.Message) {
		var status DeviceStatusReport
		if err := json.NewDecoder(bytes.NewReader(message.Payload())).Decode(&status); err != nil {
		log.Println("Err", err)
		return
		}

		log.Println("Received ", status)

		raw := db.GetOrDefault("DeviceStatus", make(map[string]DeviceStatusReport))
		deviceStatus, ok := raw.(map[string]DeviceStatusReport)
		if !ok {
			log.Println("Unable to store device status")
			return
		}

		deviceStatus[status.Name] = status
		db.Set("DeviceStatus", deviceStatus)
	}
}

func NewDeviceStatusPoll(db api.MemDB, client *mqtt.Client) ActionFunc {
	return func() {
		db.Set("LastDeviceStatusPoll", time.Now())

		client.Publish("/devicestatus", 0, false, "{'request': 'hello'}");
	}
}


