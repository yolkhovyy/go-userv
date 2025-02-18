package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/yolkhovyy/user/internal/config"
	router "github.com/yolkhovyy/user/internal/router/gin"
	"github.com/yolkhovyy/user/internal/server/http"
	storage "github.com/yolkhovyy/user/internal/storage/postgres"
)

const (
	defaultHTTPPort              = 8080
	defaultHTTPShutdownTimeout   = 5 * time.Second
	defaultHTTPReadHeaderTimeout = 1 * time.Second
	defaultRouterMode            = "release"
)

type Config struct {
	HTTP     http.Config    `yaml:"http" mapstructure:"HTTP"`
	Postgres storage.Config `yaml:"postgres" mapstructure:"Postgres"`
	Router   router.Config  `yaml:"router" mapstructure:"Router"`
}

func (c *Config) Load(prefix string) error {
	viper.SetDefault("HTTP.Port", defaultHTTPPort)
	viper.SetDefault("HTTP.ShutdownTimeout", defaultHTTPShutdownTimeout)
	viper.SetDefault("HTTP.readHeaderTimeout", defaultHTTPReadHeaderTimeout)
	viper.SetDefault("Router.Mode", defaultRouterMode)

	err := config.Load(prefix, c)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}
