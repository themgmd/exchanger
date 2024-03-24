package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"sync"
	"time"
)

type Config struct {
	HTTP        HTTPConfig
	CurrencyApi CurrencyApiConfig
	Postgre     PostgresConfig
}

type HTTPConfig struct {
	Host              string        `env:"HTTP_HOST"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT"`
}

type PostgresConfig struct {
	Name        string `env:"DATABASE_NAME"`
	User        string `env:"DATABASE_USER"`
	Password    string `env:"DATABASE_PASSWORD"`
	Port        int    `env:"DATABASE_PORT"`
	Host        string `env:"DATABASE_HOST"`
	MaxIdleConn int    `env:"MAX_IDLE_CONN"`
	MaxOpenConn int    `env:"MAX_OPEN_CONN"`
}

func (pc PostgresConfig) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		pc.User, pc.Password, pc.Host, pc.Port, pc.Name)
}

type CurrencyApiConfig struct {
	Link string `env:"CURRENCY_API_LINK"`
	Key  string `env:"CURRENCY_API_KEY"`
}

var (
	cfg  = &Config{}
	once = &sync.Once{}
)

func Get() *Config {
	once.Do(func() {
		cfg = &Config{}
		if err := env.Parse(cfg); err != nil {
			log.Fatalf("error occured while parse env: %s", err.Error())
		}
	})

	return cfg
}
