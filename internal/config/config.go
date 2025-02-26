package config

import (
	"fmt"

	"github.com/yolkhovyy/go-utilities/viperx"
)

func Load(
	configFile string,
	prefix string,
	defaults map[string]any,
	config any,
) error {
	vprx := viperx.New(configFile, prefix, nil)

	vprx.SetDefaults(defaults)

	if err := vprx.Load(config); err != nil {
		return fmt.Errorf("load config error: %w", err)
	}

	return nil
}
