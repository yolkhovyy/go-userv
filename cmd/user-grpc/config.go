package main

import (
	"fmt"

	"github.com/yolkhovyy/user/internal/config"
	router "github.com/yolkhovyy/user/internal/router/grpc"
	server "github.com/yolkhovyy/user/internal/server/grpc"
	storage "github.com/yolkhovyy/user/internal/storage/postgres"
)

type Config struct {
	GRPC     server.Config  `yaml:"grpc" mapstructure:"GRPC"`
	Postgres storage.Config `yaml:"postgres" mapstructure:"Postgres"`
	Router   router.Config  `yaml:"router" mapstructure:"Router"`
}

func (c *Config) Load(
	configFile string,
	prefix string,
) error {
	if err := config.Load(configFile, prefix, nil, defaults(), c); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

func defaults() map[string]any {
	return map[string]any{
		"GRPC.Port":            defaultGRPCPort,
		"GRPC.ShutdownTimeout": defaultRouterMode,
	}
}

const (
	defaultGRPCPort   = 50501
	defaultRouterMode = "release"
)
