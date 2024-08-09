package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

var (
	brokerURL     = getEnv("MQTT_BROKER_URL", "mqtt://test.mosquitto.org:1883")
	clientID      = "go-client"
	topic1        = "sensors/esptemp04/temperature"
	influxDBAddress = getEnv("INFLUXDB_ADDRESS", "http://influxdb:8086")
	influxDBName  = "sensors"
	influxUsername = "" // If authentication is required
	influxPassword = ""
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func connectToBroker(client mqtt.Client) error {
    maxRetries := 5
    for retries := 0; retries < maxRetries; retries++ {
        if token := client.Connect(); token.Wait() && token.Error() != nil {
            log.Printf("Error connecting to broker (attempt %d/%d): %v\n", retries+1, maxRetries, token.Error())
            time.Sleep(2 * time.Second) // Wait before retrying
        } else {
            log.Println("Connected to broker")
            return nil
        }
    }
    return fmt.Errorf("failed to connect to broker after %d attempts", maxRetries)
}

func main() {
    // Create a new MQTT client options
    opts := mqtt.NewClientOptions()
    opts.AddBroker(brokerURL)
    opts.SetClientID(clientID)
    opts.SetCleanSession(true)

    // Set up the message handler
    opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), msg.Payload())

        // Parse the payload and log to InfluxDB
        err := logToInfluxDB(string(msg.Payload()))
        if err != nil {
            log.Printf("Error logging to InfluxDB: %v\n", err)
        }
    })

    // Create a new MQTT client
    client := mqtt.NewClient(opts)

    // Connect to the broker with retries
    if err := connectToBroker(client); err != nil {
        log.Fatal(err)
    }

    // Subscribe to topics
    if token := client.Subscribe(topic1, 0, nil); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
    }
    fmt.Printf("Subscribed to topic %s\n", topic1)

    // Wait forever
    select {}
}


func logToInfluxDB(temperature string) error {
	// Create a new InfluxDB client
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     influxDBAddress,
		Username: influxUsername,
		Password: influxPassword,
	})
	if err != nil {
		return err
	}
	defer c.Close()

	// Create a new point batch
	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  influxDBName,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	// Create a point and add to batch
	tags := map[string]string{"sensor": "esptemp04"}
	fields := map[string]interface{}{
		"temperature": temperature,
	}

	pt, err := influxdb.NewPoint("temperature", tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		return err
	}

	return nil
}
