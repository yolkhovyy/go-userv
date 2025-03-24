package gin

type Config struct {
	Mode string `yaml:"mode" mapstructure:"Mode"`
}

func Defaults() map[string]any {
	return map[string]any{
		"Router.Mode": defaultRouterMode,
	}
}

const defaultRouterMode = "release"
