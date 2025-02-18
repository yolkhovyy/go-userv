package main

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/yolkhovyy/user/internal/config"
	"github.com/yolkhovyy/user/internal/notifier"
	router "github.com/yolkhovyy/user/internal/router/grpc"
	server "github.com/yolkhovyy/user/internal/server/grpc"
	"github.com/yolkhovyy/user/internal/storage"
)

const (
	defaultGRPCPort   = 50501
	defaultRouterMode = "release"
)

type Config struct {
	GRPC     server.Config   `yaml:"grpc" mapstructure:"GRPC"`
	Postgres storage.Config  `yaml:"postgres" mapstructure:"Postgres"`
	Kafka    notifier.Config `yaml:"kafka" mapstructure:"Kafka"`
	Router   router.Config   `yaml:"router" mapstructure:"Router"`
}

func (c *Config) Load(prefix string) error {
	viper.SetDefault("GRPC.Port", defaultGRPCPort)
	viper.SetDefault("Router.Mode", defaultRouterMode)

	err := config.Load(prefix, c)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}
