package main

import (
	"fmt"
	"time"

	"github.com/yolkhovyy/user/internal/config"
	router "github.com/yolkhovyy/user/internal/router/gin"
	"github.com/yolkhovyy/user/internal/server/http"
	storage "github.com/yolkhovyy/user/internal/storage/postgres"
)

type Config struct {
	HTTP     http.Config    `yaml:"http" mapstructure:"HTTP"`
	Postgres storage.Config `yaml:"postgres" mapstructure:"Postgres"`
	Router   router.Config  `yaml:"router" mapstructure:"Router"`
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
	return map[string]any{
		"HTTP.Port":              defaultHTTPPort,
		"HTTP.ShutdownTimeout":   defaultHTTPShutdownTimeout,
		"HTTP.ReadHeaderTimeout": defaultHTTPReadHeaderTimeout,
		"Router.Mode":            defaultRouterMode,
	}
}

const (
	defaultHTTPPort              = 8080
	defaultHTTPShutdownTimeout   = 5 * time.Second
	defaultHTTPReadHeaderTimeout = 1 * time.Second
	defaultRouterMode            = "release"
)
