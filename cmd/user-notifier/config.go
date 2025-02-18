package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/yolkhovyy/user/internal/notifier"
	"github.com/yolkhovyy/user/internal/storage"
)

const (
	defaultPort              = 8080
	defaultShutdownTimeout   = 5 * time.Second
	defaultReadHeaderTimeout = 1 * time.Second
	defaultGinMode           = "release"
)

type Config struct {
	Postgres storage.Config  `yaml:"postgres" mapstructure:"Postgres"`
	Kafka    notifier.Config `yaml:"kafka" mapstructure:"Kafka"`
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
	viper.SetDefault("Port", defaultPort)
	viper.SetDefault("ShutdownTimeout", defaultShutdownTimeout)
	viper.SetDefault("readHeaderTimeout", defaultReadHeaderTimeout)
	viper.SetDefault("GinMode", "release")
}
