package main

import (
	"fmt"

	"github.com/yolkhovyy/go-userv/internal/config"
	grpcrouter "github.com/yolkhovyy/go-userv/internal/router/grpc"
	grpcserver "github.com/yolkhovyy/go-userv/internal/server/grpc"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
)

type Config struct {
	GRPC     grpcserver.Config `yaml:"grpc" mapstructure:"GRPC"`
	Router   grpcrouter.Config `yaml:"router" mapstructure:"Router"`
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
		defaultGRPCPort   = 50501
		defaultRouterMode = "release"
	)

	return map[string]any{
		"GRPC.Port":            defaultGRPCPort,
		"GRPC.ShutdownTimeout": defaultRouterMode,
	}
}
