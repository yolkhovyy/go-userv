package grpc

type Config struct {
	Port       int  `yaml:"port" mapstructure:"Port"`
	Reflection bool `yaml:"reflection" mapstructure:"Reflection"`
}
