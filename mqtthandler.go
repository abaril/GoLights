package main

import (
	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"log"
	"time"
)

const RETRY_INTERVAL time.Duration = 10 * time.Second
var handlers map[string]mqtt.MessageHandler = make(map[string]mqtt.MessageHandler)

func mqttStart() {
	options := &mqtt.ClientOptions{
		ClientID: "golights",
	}
	options.AddBroker("tcp://127.0.0.1:1883")
	options.SetConnectionLostHandler(mqttConnect)
	client := mqtt.NewClient(options)
	go mqttConnect(client, nil)
}

func mqttConnect(client *mqtt.Client, err error) {
	if err != nil {
		log.Println("MQTT connection lost", err)
	}

	for !client.IsConnected() {

		token := client.Connect();
		token.Wait();
		if token.Error() != nil {
			log.Println("Unable to connect to MQTT broker", token.Error())
		}

		if !client.IsConnected() {
			time.Sleep(RETRY_INTERVAL)
		}
	}

	mqttSubscribe(client)
}

func mqttSubscribe(client *mqtt.Client) {
	for topic, handler := range handlers {
		client.Subscribe(topic, 0, handler)
	}
}

func mqttHandleFunc(topic string, handler mqtt.MessageHandler) {
	handlers[topic] = handler
}