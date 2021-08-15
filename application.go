package main

import (
	"fmt"
	"log"
	"path"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pkg/errors"
)

type application struct {
	cfg      configuration
	sensors  map[string]*Sensor
	writeAPI api.WriteAPI
}

func newApplication(cfg configuration) (*application, error) {
	app := &application{
		cfg:     cfg,
		sensors: map[string]*Sensor{},
	}

	if cfg.InfluxDBServer != "" {
		client := influxdb2.NewClient(cfg.InfluxDBServer, cfg.InfluxDBToken)
		app.writeAPI = client.WriteAPI(cfg.InfluxDBOrg, cfg.InfluxDBBucket)
	}

	return app, nil
}

func (app *application) run() error {
	log.Printf("[I] Creating state directory <%s>", defaultStateDirectory)
	if err := createDirectory(defaultStateDirectory); err != nil {
		return errors.Wrapf(err, "Could not create state directory <%s>", defaultStateDirectory)
	}

	storageDirectory := path.Join(path.Join(defaultStateDirectory, "storage"))
	log.Printf("[I] Creating storage directory <%s>", storageDirectory)
	if err := createDirectory(storageDirectory); err != nil {
		return errors.Wrapf(err, "Could not create storage directory <%s>", storageDirectory)
	}

	// hc.OnTermination(func() {
	// 	log.Println("hc.OnTermination")
	// 	for address, sensor := range sensors {
	// 		<-sensor.transport.Stop()
	// 		delete(sensors, address)
	// 	}
	// 	time.Sleep(100 * time.Millisecond)
	// 	os.Exit(1)
	// })

	scanner, err := NewScanner()
	if err != nil {
		return errors.Wrap(err, "Could not create a Scanner")
	}

	if err := scanner.Start(); err != nil {
		return errors.Wrap(err, "Could not start scanner")
	}

	subscription := scanner.Subscribe()

	for event := range subscription.Events() {
		switch event.(type) {
		case EventDiscoveredSensor:
			log.Println("EventDiscoveredSensor")
			app.addSensor(event.(EventDiscoveredSensor).Address, event.(EventDiscoveredSensor).Data)
		case EventReceivedSensorData:
			log.Println("EventReceivedSensorData")
			app.updateSensor(event.(EventReceivedSensorData).Address, event.(EventReceivedSensorData).Data)
		case EventExpiredSensor:
			log.Println("EventExpiredSensor")
			app.removeSensor(event.(EventExpiredSensor).Address)
		}
	}

	return nil
}

func (app *application) addSensor(address string, data SensorData) {
	if _, found := app.sensors[address]; !found {
		sensor, err := NewSensor(address, data, app.cfg.Pin)
		if err != nil {
			log.Printf("[E] Could not create sensor <%s>: %s", address, err)
			return
		}
		app.sensors[address] = sensor
		sensor.Update(data)
	}
}

func (app *application) updateSensor(address string, data SensorData) {
	if sensor, found := app.sensors[address]; found {
		if err := sensor.Update(data); err != nil {
			log.Printf("[E] Failed to update Sensor <%s>: %s", address, err)
		}

		if app.writeAPI != nil {
			metrics := fmt.Sprintf("batteryLevel=%d,humidity=%f,temperature=%f", data.BatteryLevel, data.Humidity, data.Temperature)
			record := fmt.Sprintf("mijia,sensor=%s %s %d", address, metrics, time.Now().UnixNano())

			log.Printf("[D] Influx Record: %s", record)

			app.writeAPI.WriteRecord(record)
			app.writeAPI.Flush()
		}
	}
}

func (app *application) removeSensor(address string) {
	// TODO Do we just remove the sensor storage?
}
