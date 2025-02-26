package config

import (
	"fmt"
	"strings"

	"github.com/yolkhovyy/go-utilities/viperx"
)

func Load(
	configFile string,
	prefix string,
	replacer *strings.Replacer,
	defaults map[string]any,
	config any,
) error {
	vprx := viperx.New(configFile, prefix, replacer)

	vprx.SetDefaults(defaults)

	if err := vprx.Load(config); err != nil {
		return fmt.Errorf("load config error: %w", err)
	}

	return nil
}
