package main

import (
	"fmt"
	"time"

	"github.com/yolkhovyy/go-userv/internal/config"
	ginrouter "github.com/yolkhovyy/go-userv/internal/router/gin"
	httpserver "github.com/yolkhovyy/go-userv/internal/server/http"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
)

type Config struct {
	HTTP     httpserver.Config `yaml:"http" mapstructure:"HTTP"`
	Router   ginrouter.Config  `yaml:"router" mapstructure:"Router"`
	Postgres postgres.Config   `yaml:"postgres" mapstructure:"Postgres"`
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
		defaultHTTPPort              = 8080
		defaultHTTPShutdownTimeout   = 5 * time.Second
		defaultHTTPReadHeaderTimeout = 1 * time.Second
		defaultRouterMode            = "release"
	)

	return map[string]any{
		"HTTP.Port":              defaultHTTPPort,
		"HTTP.ShutdownTimeout":   defaultHTTPShutdownTimeout,
		"HTTP.ReadHeaderTimeout": defaultHTTPReadHeaderTimeout,
		"Router.Mode":            defaultRouterMode,
	}
}
