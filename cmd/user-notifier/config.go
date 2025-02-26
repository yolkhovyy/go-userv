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

func (c *Config) Load(
	configFile string,
	prefix string,
) error {
	if err := config.Load(configFile, prefix, nil, defaults(), c); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

func defaults() map[string]any {
	return map[string]any{}
}
