package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

var (
	brokerURL = "mqtt://test.mosquitto.org:1883" // Use "ssl://test.mosquitto.org:8883" for encrypted
	clientID  = "go-client"
	topic1    = "sensors/esptemp04/temperature"

	influxDBAddress = "http://influxdb:8086"
	influxDBName    = "sensors"
	influxUsername  = "" // If authentication is required
	influxPassword  = ""
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

		// Parse the payload and log to InfluxDB
		err := logToInfluxDB(string(msg.Payload()))
		if err != nil {
			log.Printf("Error logging to InfluxDB: %v\n", err)
		}
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

func logToInfluxDB(temperatureStr string) error {
	// Convert the temperature string to a float64
	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse temperature '%s': %w", temperatureStr, err)
	}

	// Create a new InfluxDB client
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     influxDBAddress,
		Username: influxUsername,
		Password: influxPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to create InfluxDB client: %w", err)
	}
	defer func() {
		if cerr := c.Close(); cerr != nil {
			log.Printf("failed to close InfluxDB client: %v", cerr)
		}
	}()

	// Create a new point batch
	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  influxDBName,
		Precision: "s",
	})
	if err != nil {
		return fmt.Errorf("failed to create batch points: %w", err)
	}

	// Create a point and add to batch
	tags := map[string]string{"sensor": "esptemp04"}
	fields := map[string]interface{}{
		"temperature": temperature, // Use the float64 value here
	}

	pt, err := influxdb.NewPoint("temperature", tags, fields, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create point: %w", err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		return fmt.Errorf("failed to write batch points: %w", err)
	}

	return nil
}

