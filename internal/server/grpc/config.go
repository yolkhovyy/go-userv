package grpc

type Config struct {
	Port       int  `yaml:"port" mapstructure:"Port"`
	Reflection bool `yaml:"reflection" mapstructure:"Reflection"`
}

func Defaults() map[string]any {
	return map[string]any{
		"GRPC.Port":            defaultGRPCPort,
		"GRPC.ShutdownTimeout": defaultRouterMode,
	}
}

const (
	defaultGRPCPort   = 50501
	defaultRouterMode = "release"
)
