package mappers

import (
	"exchanger/internal/config"
	"exchanger/pkg/database/postgres"
)

func MapPostgresConfig(pgConfig config.PostgresConfig) postgres.Config {
	return postgres.Config{
		Host:     pgConfig.Host,
		Port:     pgConfig.Port,
		DBName:   pgConfig.Name,
		User:     pgConfig.User,
		Password: pgConfig.Password,
	}
}
