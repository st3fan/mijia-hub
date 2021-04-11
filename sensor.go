package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Sensor struct {
	*accessory.Accessory
	temperature *service.TemperatureSensor
	humidity    *service.HumiditySensor
	battery     *service.BatteryService
	transport   hc.Transport
}

func NewSensor(address string, data SensorData) (*Sensor, error) {
	info := accessory.Info{
		Name:         strings.ToUpper(fmt.Sprintf("ATC_%s%s%s", address[9:11], address[12:14], address[15:17])),
		Manufacturer: "Mija",
		SerialNumber: address,
		Model:        "LYWSD03MMC",
	}

	sensor := &Sensor{}
	sensor.Accessory = accessory.New(info, accessory.TypeSensor)

	sensor.temperature = service.NewTemperatureSensor()
	sensor.temperature.CurrentTemperature.SetValue(data.Temperature)
	sensor.Accessory.AddService(sensor.temperature.Service)

	sensor.humidity = service.NewHumiditySensor()
	sensor.humidity.CurrentRelativeHumidity.SetValue(float64(data.Humidity))
	sensor.Accessory.AddService(sensor.humidity.Service)

	sensor.battery = service.NewBatteryService()
	sensor.battery.BatteryLevel.SetValue(data.BatteryLevel)
	sensor.Accessory.AddService(sensor.battery.Service)

	storageRoot := path.Join(defaultStateDirectory, "storage")

	config := hc.Config{
		StoragePath: path.Join(storageRoot, info.Name),
		Pin:         "11223344",
	}

	transport, err := hc.NewIPTransport(config, sensor.Accessory)
	if err != nil {
		return nil, err
	}

	go transport.Start()

	return sensor, nil
}

func (s *Sensor) Update(data SensorData) error {
	if s.temperature.CurrentTemperature.GetValue() != data.Temperature {
		s.temperature.CurrentTemperature.SetValue(data.Temperature)
	}
	if s.humidity.CurrentRelativeHumidity.GetValue() != data.Humidity {
		s.humidity.CurrentRelativeHumidity.SetValue(data.Humidity)
	}
	if s.battery.BatteryLevel.GetValue() != data.BatteryLevel {
		s.battery.BatteryLevel.SetValue(data.BatteryLevel)
	}
	return nil
}
