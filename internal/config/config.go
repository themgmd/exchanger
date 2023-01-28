package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	FolderPath          = "configs"
	defaultHTTPPort     = "5023"
	defaultPostgresPort = "5432"
	defaultPostgresHost = "localhost"
)

type (
	Config struct {
		HTTP     HTTPConfig     `mapstructure:"http"`
		Database DatabaseConfig `mapstructure:"database"`
		API      APIConfig      `mapstructure:"api"`
	}

	HTTPConfig struct {
		Port    string        `mapstructure:"port"`
		Timeout TimeoutConfig `mapstructure:"timeout"`
	}

	TimeoutConfig struct {
		Read  time.Duration `mapstructure:"read"`
		Write time.Duration `mapstructure:"write"`
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

func New(cfgFile string) (*Config, error) {
	setupDefaultValues()

	if err := parseConfigFile(cfgFile); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshall(&cfg); err != nil {
		return nil, err
	}

	parseEnvFile(&cfg)

	return &cfg, nil
}

func parseConfigFile(file string) error {
	viper.AddConfigPath(FolderPath)
	viper.SetConfigName(file)

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
