package main

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Configuration struct {
	DatabasePath   string `env:"THERMON_DB_PATH" envDefault:"thermon.db"`
	SampleInterval int    `env:"THERMON_SAMPLE_INTERVAL" envDefault:"60"`
	Retention      int    `env:"THERMON_RETENTION" envDefault:"604800"`
	Port           int    `env:"THERMON_PORT" envDefault:"8080"`
}

var configuration *Configuration

func loadConfiguration() {
	var cfg Configuration
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error parsing configuration: %s", err.Error())
	}

	configuration = &cfg
}

func GetConfig() *Configuration {
	if configuration == nil {
		loadConfiguration()
	}

	return configuration
}
