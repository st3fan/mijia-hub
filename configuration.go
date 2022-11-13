package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type configuration struct {
	Pin            string `envconfig:"PIN"`
	InfluxDBServer string `envconfig:"INFLUXDB_SERVER"`
	InfluxDBToken  string `envconfig:"INFLUXDB_TOKEN"`
	InfluxDBBucket string `envconfig:"INFLUXDB_BUCKET"`
	InfluxDBOrg    string `envconfig:"INFLUXDB_ORG"`
	Verbose        bool   `envconfig:"VERBOSE"`
}

func newConfigurationFromEnvironment() (configuration, error) {
	var cfg configuration
	if err := envconfig.Process("", &cfg); err != nil {
		return configuration{}, errors.Wrap(err, "can't to parse configuration from environment")
	}
	return cfg, nil
}
