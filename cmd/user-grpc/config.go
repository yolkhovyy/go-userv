package main

import (
	"fmt"

	"github.com/yolkhovyy/go-otelw/pkg/otelw"
	grpcrouter "github.com/yolkhovyy/go-userv/internal/router/grpc"
	grpcserver "github.com/yolkhovyy/go-userv/internal/server/grpc"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
	"github.com/yolkhovyy/go-utilities/viperx"
)

type Config struct {
	otelw.Config `yaml:",inline" mapstructure:",squash"`
	GRPC         grpcserver.Config `yaml:"grpc" mapstructure:"GRPC"`
	Router       grpcrouter.Config `yaml:"router" mapstructure:"Router"`
	Postgres     postgres.Config   `yaml:"postgres" mapstructure:"Postgres"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(configFile string, prefix string) error {
	vprx := viperx.New(configFile, prefix, nil)

	vprx.SetDefaults(otelw.Defaults())
	vprx.SetDefaults(grpcserver.Defaults())
	vprx.SetDefaults(postgres.Defaults())

	if err := vprx.Load(c); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}
