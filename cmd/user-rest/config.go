package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/yolkhovyy/user/internal/notifier"
	router "github.com/yolkhovyy/user/internal/router/gin"
	"github.com/yolkhovyy/user/internal/server/http"
	"github.com/yolkhovyy/user/internal/storage"
)

type Config struct {
	HTTP     http.Config     `yaml:"http" mapstructure:"HTTP"`
	Postgres storage.Config  `yaml:"postgres" mapstructure:"Postgres"`
	Kafka    notifier.Config `yaml:"kafka" mapstructure:"Kafka"`
	Router   router.Config   `yaml:"router" mapstructure:"Router"`
}

func (c *Config) Load(configFile string, prefix string) error {
	c.initDefaults()

	viper.Reset()
	viper.SetConfigFile(configFile)
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config read error: %w", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("config unmarshal error: %w", err)
	}

	return nil
}

func (c *Config) initDefaults() {
	viper.SetDefault("HTTP.Port", defaultHTTPPort)
	viper.SetDefault("HTTP.ShutdownTimeout", defaultHTTPShutdownTimeout)
	viper.SetDefault("HTTP.readHeaderTimeout", defaultHTTPReadHeaderTimeout)
	viper.SetDefault("Router.Mode", defaultRouterMode)
}

const (
	defaultHTTPPort              = 8080
	defaultHTTPShutdownTimeout   = 5 * time.Second
	defaultHTTPReadHeaderTimeout = 1 * time.Second
	defaultRouterMode            = "release"
)
