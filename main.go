package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
    brokerURL  = "mqtt://test.mosquitto.org:1883" // Use "ssl://test.mosquitto.org:8883" for encrypted
    clientID   = "go-client"
    topic1      = "sensors/esptemp04/temperature"
)

func main() {
    // Create a new MQTT client options
    opts := mqtt.NewClientOptions()
    opts.AddBroker(brokerURL)
    opts.SetClientID(clientID)
    opts.SetCleanSession(true)
    
    // Set up the message handler
    opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), msg.Payload())
    })
    
    // Create a new MQTT client
    client := mqtt.NewClient(opts)
    
    // Connect to the broker
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    
    fmt.Println("Connected to broker")

    // Subscribe to topics
    if token := client.Subscribe(topic1, 0, nil); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    fmt.Printf("Subscribed to topic %s\n", topic1)

    // Wait forever
    select {}
}
