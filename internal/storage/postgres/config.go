package postgres

type Config struct {
	Host     string `yaml:"host" mapstructure:"Host"`
	Port     int    `yaml:"port" mapstructure:"Port"`
	Database string `yaml:"database" mapstructure:"Database"`
	Username string `yaml:"username" mapstructure:"Username"`
	Password string `yaml:"password" mapstructure:"Password"`
}

func Defaults() map[string]any {
	return map[string]any{
		"Postgres.Host":     DefaultHost,
		"Postgres.Port":     DefaultPort,
		"Postgres.Database": DefaultDatabase,
		"Postgres.Username": DefaultUsername,
		"Postgres.Password": DefaultPassword,
	}
}

const (
	DefaultHost     = "localhost"
	DefaultPort     = 5432
	DefaultDatabase = "postgres"
	DefaultUsername = "postgres"
	DefaultPassword = "postgres"
)
