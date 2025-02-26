package main

import (
	"fmt"

	"github.com/yolkhovyy/user/internal/config"
	"github.com/yolkhovyy/user/internal/notifier"
	storage "github.com/yolkhovyy/user/internal/storage/postgres"
)

type Config struct {
	Postgres storage.Config  `yaml:"postgres" mapstructure:"Postgres"`
	Kafka    notifier.Config `yaml:"kafka" mapstructure:"Kafka"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(
	configFile string,
	prefix string,
) error {
	if err := config.Load(configFile, prefix, defaults(), c); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

func defaults() map[string]any {
	const (
		defaultPostgresHost     = "localhost"
		defaultPostgresPort     = 5432
		defaultPostgresDatabase = "postgres"
		defaultPostgresUsername = "postgres"
		defaultPostgresPassword = "postgres"
		defaultKafkaBrokers     = "localhost:9092"
	)

	return map[string]any{
		"Postrgres.Host":     defaultPostgresHost,
		"Postrgres.Port":     defaultPostgresPort,
		"Postrgres.Database": defaultPostgresDatabase,
		"Postrgres.Username": defaultPostgresUsername,
		"Postrgres.Password": defaultPostgresPassword,
		"Kafka.Brokers":      defaultKafkaBrokers,
	}
}
