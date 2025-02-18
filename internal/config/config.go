package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Load(prefix string, config any) error {
	configFile := flag.String("config", "config.yml",
		"Path to the configuration file (default: config.yml)")

	flag.Parse()

	viper.Reset()
	viper.SetConfigFile(*configFile)
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config read error: %w", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("config unmarshal error: %w", err)
	}

	return nil
}
