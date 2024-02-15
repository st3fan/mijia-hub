package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

type Sensor struct {
	accessory *accessory.A
	temperature service.TemperatureSensor
	humidity    service.HumiditySensor
	battery     service.BatteryService
}

func NewSensor(address string, data SensorData, pin string) (*Sensor, error) {
	info := accessory.Info{
		Name:         strings.ToUpper(fmt.Sprintf("ATC_%s%s%s", address[9:11], address[12:14], address[15:17])),
		Manufacturer: "Mijia",
		SerialNumber: address,
		Model:        "LYWSD03MMC",
	}

	sensor := &Sensor{}
	sensor.accessory = accessory.New(info, accessory.TypeSensor)
  // define temperature
  sensor.temperature.S = service.New(service.TypeTemperatureSensor)
  sensor.temperature.CurrentTemperature = characteristic.NewCurrentTemperature()
	sensor.temperature.AddC(sensor.temperature.CurrentTemperature.C)
  // add temperature service
  sensor.accessory.AddS(sensor.temperature.S)
  // set current temperature
  sensor.temperature.CurrentTemperature.SetValue(data.Temperature)
  // define humidity 
  sensor.humidity.S = service.New(service.TypeHumiditySensor)
  sensor.humidity.CurrentRelativeHumidity = characteristic.NewCurrentRelativeHumidity()
  sensor.humidity.AddC(sensor.humidity.CurrentRelativeHumidity.C)
  // add humidity service
  sensor.accessory.AddS(sensor.humidity.S)
  // set current humidity
  sensor.humidity.CurrentRelativeHumidity.SetValue(data.Humidity)
  // define battery
  sensor.battery.S = service.New(service.TypeBatteryService)
  sensor.battery.BatteryLevel = characteristic.NewBatteryLevel()
	sensor.battery.AddC(sensor.battery.BatteryLevel.C)
  // add battery service
  sensor.accessory.AddS(sensor.battery.S)
  // set current battery
  sensor.battery.BatteryLevel.SetValue(data.BatteryLevel)

	storageRoot := path.Join(defaultStateDirectory, "storage")

	fs := hap.NewFsStore(storageRoot)

  server, err := hap.NewServer(fs, sensor.accessory)
	if err != nil {
    return nil, err
	}
  server.Pin = pin
	return sensor, nil
}

func (s *Sensor) Update(data SensorData) error {
		s.temperature.CurrentTemperature.SetValue(data.Temperature)
		s.humidity.CurrentRelativeHumidity.SetValue(data.Humidity)
		s.battery.BatteryLevel.SetValue(data.BatteryLevel)
	return nil
}
