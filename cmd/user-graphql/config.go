package main

import (
	"fmt"
	"time"

	"github.com/yolkhovyy/go-userv/internal/config"
	gqlrouter "github.com/yolkhovyy/go-userv/internal/router/graphql"
	httpserver "github.com/yolkhovyy/go-userv/internal/server/http"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
)

type Config struct {
	HTTP     httpserver.Config `yaml:"http" mapstructure:"HTTP"`
	Router   gqlrouter.Config  `yaml:"graphql" mapstructure:"GraphQL"`
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
