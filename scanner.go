package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
)

type SensorData struct {
	ReceivedAt   time.Time
	BatteryLevel int
	Humidity     float64
	Temperature  float64
}

func ParseSensorData(data []byte) (SensorData, error) {
	if len(data) != 13 {
		return SensorData{}, errors.New("invalid data; too short")
	}

	return SensorData{
		ReceivedAt:   time.Now(),
		BatteryLevel: int(data[9]),
		Humidity:     float64(data[8]),
		Temperature:  float64(int(data[6])*256|int(data[7])) / 10.0,
	}, nil
}

//

type EventDiscoveredSensor struct {
	Address string
	Data    SensorData
}

type EventReceivedSensorData struct {
	Address string
	Data    SensorData
}

type EventExpiredSensor struct {
	Address string
}

type Scanner struct {
	device  ble.Device
	sensors map[string]SensorData
	SubscriptionProvider
	sync.RWMutex
}

func NewScanner() (*Scanner, error) {
	device, err := dev.NewDevice("default")
	if err != nil {
		return nil, err
	}

	ble.SetDefaultDevice(device)

	return &Scanner{
		device:  device,
		sensors: map[string]SensorData{},
	}, nil
}

func (s *Scanner) handler(a ble.Advertisement) bool {
	s.Lock()
	defer s.Unlock()

	if strings.HasPrefix(a.Address().String(), "a4:c1:38") && len(a.ServiceData()) == 1 {
		if sensorData, err := ParseSensorData(a.ServiceData()[0].Data); err == nil {
			//log.Printf("%s %+v", a.Address(), sensorData)
			if sd, ok := s.sensors[a.Address().String()]; !ok {
				s.sensors[a.Address().String()] = sensorData
				s.Notify(EventDiscoveredSensor{
					Address: a.Address().String(),
					Data:    sensorData,
				})
			} else {
				if sd.ReceivedAt.Before(time.Now().Add(-60 * time.Second)) {
					s.sensors[a.Address().String()] = sensorData
					s.Notify(EventReceivedSensorData{
						Address: a.Address().String(),
						Data:    sensorData,
					})
				}
			}
		}
	}
	return false
}

func (s *Scanner) start() {
	if _, err := ble.Connect(context.Background(), s.handler); err != nil {
		log.Printf("Could not scan: %s", err)
	}

	// if err := s.device.Scan(context.Background(), true, s.handler); err != nil {
	// 	// What to do? Error channel?
	// 	log.Printf("Could not scan: %s", err)
	// }
}

func (s *Scanner) Start() error {
	go s.start()
	return nil
}
