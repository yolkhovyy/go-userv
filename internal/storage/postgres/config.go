package postgres

type Config struct {
	Host           string `yaml:"host" mapstructure:"Host"`
	Port           int    `yaml:"port" mapstructure:"Port"`
	Database       string `yaml:"database" mapstructure:"Database"`
	Username       string `yaml:"username" mapstructure:"Username"`
	Password       string `yaml:"password" mapstructure:"Password"`
	EnableNotifier bool   `yaml:"enableNotifier" mapstructure:"EnableNotifier"`
}
