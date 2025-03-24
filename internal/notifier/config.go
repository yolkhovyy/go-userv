package notifier

type Config struct {
	Brokers []string `yaml:"brokers" mapstructure:"Brokers"`
}

func Defaults() map[string]any {
	return map[string]any{
		"Kafka.Brokers": defaultKafkaBrokers,
	}
}

const defaultKafkaBrokers = "localhost:9092"
