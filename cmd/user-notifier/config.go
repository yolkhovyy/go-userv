package main

import (
	"fmt"

	"github.com/yolkhovyy/go-otelw/otelw"
	"github.com/yolkhovyy/go-userv/internal/notifier"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
	"github.com/yolkhovyy/go-utilities/viperx"
)

type Config struct {
	otelw.Config `yaml:",inline" mapstructure:",squash"`
	Kafka        notifier.Config `yaml:"kafka" mapstructure:"Kafka"`
	Postgres     postgres.Config `yaml:"postgres" mapstructure:"Postgres"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(configFile string, prefix string) error {
	vprx := viperx.New(configFile, prefix, nil)

	vprx.SetDefaults(otelw.Defaults())
	vprx.SetDefaults(notifier.Defaults())
	vprx.SetDefaults(postgres.Defaults())

	if err := vprx.Load(c); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}
