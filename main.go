package main

import (
	"flag"
	"log"
	"os"
	"path"
)

var sensors map[string]*Sensor = map[string]*Sensor{}

func AddSensor(address string, data SensorData, pin string) {
	if _, found := sensors[address]; !found {
		sensor, err := NewSensor(address, data, pin)
		if err != nil {
			log.Printf("Could not create sensor <%s>: %s", address, err)
			return
		}
		sensors[address] = sensor
		sensor.Update(data)
	}
}

func UpdateSensor(address string, data SensorData) {
	if sensor, found := sensors[address]; found {
		if err := sensor.Update(data); err != nil {
			log.Printf("Failed to update Sensor <%s>: %s", address, err)
		}
	}
}

func RemoveSensor(address string) {
	// TODO
}

func createDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}

func main() {
	pin := flag.String("pin", defaultPin, "pin used to pair new sensors")
	flag.Parse()

	log.Printf("Creating directory <%s>", defaultStateDirectory)
	if err := createDirectory(defaultStateDirectory); err != nil {
		log.Fatalf("Could not create <%s>: %s", defaultStateDirectory, err)
	}

	storageDirectory := path.Join(path.Join(defaultStateDirectory, "storage"))
	log.Printf("Creating directory <%s>", storageDirectory)
	if err := createDirectory(storageDirectory); err != nil {
		log.Fatalf("Could not create <%s>: %s", storageDirectory, err)
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
		log.Fatalf("Could not create a Scanner: %s", err)
	}

	if err := scanner.Start(); err != nil {
		log.Fatalf("Could not start scanning: %s", err)
	}

	log.Printf("Scanning...")

	subscription := scanner.Subscribe()

	for event := range subscription.Events() {
		switch event.(type) {
		case EventDiscoveredSensor:
			log.Println("EventDiscoveredSensor")
			AddSensor(event.(EventDiscoveredSensor).Address, event.(EventDiscoveredSensor).Data, *pin)
		case EventReceivedSensorData:
			log.Println("EventReceivedSensorData")
			UpdateSensor(event.(EventReceivedSensorData).Address, event.(EventReceivedSensorData).Data)
		case EventExpiredSensor:
			log.Println("EventExpiredSensor")
			RemoveSensor(event.(EventExpiredSensor).Address)
		}
	}
}
