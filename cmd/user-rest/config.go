package main

import (
	"fmt"

	"github.com/yolkhovyy/go-otelw/pkg/otelw"
	ginrouter "github.com/yolkhovyy/go-userv/internal/router/gin"
	httpserver "github.com/yolkhovyy/go-userv/internal/server/http"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
	"github.com/yolkhovyy/go-utilities/viperx"
)

type Config struct {
	otelw.Config `yaml:",inline" mapstructure:",squash"`
	HTTP         httpserver.Config `yaml:"http" mapstructure:"HTTP"`
	Router       ginrouter.Config  `yaml:"router" mapstructure:"Router"`
	Postgres     postgres.Config   `yaml:"postgres" mapstructure:"Postgres"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(configFile string, prefix string) error {
	vprx := viperx.New(configFile, prefix, nil)

	vprx.SetDefaults(otelw.Defaults())
	vprx.SetDefaults(httpserver.Defaults())
	vprx.SetDefaults(ginrouter.Defaults())
	vprx.SetDefaults(postgres.Defaults())

	if err := vprx.Load(c); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}
