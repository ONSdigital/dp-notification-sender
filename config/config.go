package config

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config values for the application.
type Config struct {
	BindAddr                string        `envconfig:"BIND_ADDR"`
	KafkaAddr               []string      `envconfig:"KAFKA_ADDR"`
	FilterConsumerGroup     string        `envconfig:"FILTER_JOB_CONSUMER_GROUP"`
	FilterConsumerTopic     string        `envconfig:"FILTER_JOB_CONSUMER_TOPIC"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval     time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
}

// Get the configuration values from the environment or provide the defaults.
func Get() (*Config, error) {

	cfg := &Config{
		BindAddr:                ":23456",
		KafkaAddr:               []string{"localhost:9092"},
		FilterConsumerTopic:     "completed-jobs",
		FilterConsumerGroup:     "dp-dataset-exporter",
		GracefulShutdownTimeout: time.Second * 10,
		HealthCheckInterval:     time.Minute,
	}

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Config) String() string {
	json, _ := json.Marshal(config)
	return string(json)
}
