package config

import (
	"github.com/spf13/viper"
	"os"
)

const (
	defaultHTTPPort     = "5023"
	defaultPostgresPort = "5432"
	defaultPostgresHost = "localhost"
)

type Mode string

const (
	DEVELOPMENT Mode = "development"
	PRODUCTION  Mode = "production"
)

type (
	Config struct {
		HTTP     HTTPConfig     `mapstructure:"http"`
		Database DatabaseConfig `mapstructure:"database"`
		API      APIConfig      `mapstructure:"api"`
	}

	HTTPConfig struct {
		Port string `mapstructure:"port"`
	}

	APIConfig struct {
		Link string `mapstructure:"link"`
		Key  string `mapstructure:"key"`
	}

	DatabaseConfig struct {
		Postgres PostgresConfig `mapstructure:"postgres"`
	}

	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Name     string `mapstructure:"dbName"`
		Password string `mapstructure:"password"`
		SSLMode  string `mapstructure:"sslMode"`
	}
)

func New(cfgDir string, mode Mode) (*Config, error) {
	setupDefaultValues()

	if err := parseConfigFile(cfgDir, mode); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshall(&cfg); err != nil {
		return nil, err
	}

	parseEnvFile(&cfg)

	return &cfg, nil
}

func parseConfigFile(folder string, mode Mode) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName(string(mode))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func unmarshall(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("api", &cfg.API); err != nil {
		return err
	}
	return viper.UnmarshalKey("database", &cfg.Database)
}

func parseEnvFile(cfg *Config) {
	cfg.Database.Postgres.Password = os.Getenv("DB_POSTGRES_PASSWORD")
	cfg.API.Key = os.Getenv("EXCHANGE_API_KEY")
}

func setupDefaultValues() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("database.postgres.host", defaultPostgresHost)
	viper.SetDefault("database.postgres.port", defaultPostgresPort)
}
