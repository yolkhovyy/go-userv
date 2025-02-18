package notifier

type Config struct {
	Brokers []string `yaml:"brokers" mapstructure:"Brokers"`
}
