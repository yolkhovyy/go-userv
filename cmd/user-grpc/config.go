package main

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/yolkhovyy/user/internal/config"
	"github.com/yolkhovyy/user/internal/notifier"
	grpcrouter "github.com/yolkhovyy/user/internal/router/grpc"
	grpcserver "github.com/yolkhovyy/user/internal/server/grpc"
	storage "github.com/yolkhovyy/user/internal/storage/postgres"
)

const (
	defaultGRPCPort   = 50501
	defaultRouterMode = "release"
)

type Config struct {
	GRPC     grpcserver.Config `yaml:"grpc" mapstructure:"GRPC"`
	Postgres storage.Config    `yaml:"postgres" mapstructure:"Postgres"`
	Kafka    notifier.Config   `yaml:"kafka" mapstructure:"Kafka"`
	Router   grpcrouter.Config `yaml:"router" mapstructure:"Router"`
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
