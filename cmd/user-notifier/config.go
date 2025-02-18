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

func (c *Config) Load(prefix string) error {
	err := config.Load(prefix, c)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}
