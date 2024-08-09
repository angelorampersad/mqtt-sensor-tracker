# mqtt-sensor-tracker

## About
This is a small Go application that listen on the MQTT protocol for streaming data. The app uses the public topic `sensors/esptemp04/temperature` hosted on the `test.mosquitto.org` MQTT broker.

After sensor data is pushed into the pub/sub topic, the Go application inserts a temperature data into an InfluxDB timeseries database. Finally, a Grafana instance will read from InfluxDB for real-time visualisation of the metrics.

```text
On-site Sensor --> Public MQTT Broker --> Local Go Application --> InfluxDB --> Grafana
```

## Local development

In order to spin up the project locally, launch the services using Docker Compose:

```bash
docker-compose up
```

This starts running the Go application on the local machine (which connects to the MQTT broker) and starts up InfluxDB and Grafana instances. Check out the real-time stats in [Grafana](http://localhost:3000/)!
